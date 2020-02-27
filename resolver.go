package golang_graphql_authentication

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/machariamuguku/golang-graphql-authentication/db"
	"github.com/machariamuguku/golang-graphql-authentication/emails"
	"github.com/machariamuguku/golang-graphql-authentication/models"
	"github.com/machariamuguku/golang-graphql-authentication/sms"
	"github.com/machariamuguku/golang-graphql-authentication/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// validate and universal translate instances
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

type Resolver struct {
	DB *db.DB
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) RegisterUser(ctx context.Context, input RegisterUserInput) (*RegisterUserPayload, error) {

	// validate input fields

	// english locale
	en := en.New()
	// universal english translator
	uni = ut.New(en, en)

	// translator for english
	// this is usually know/extracted from http 'Accept-Language' header
	trans, _ := uni.GetTranslator("en")

	// initialize validate v10 instance
	validate = validator.New()

	en_translations.RegisterDefaultTranslations(validate, trans)

	// validate against the validation struct
	// returns nil or ValidationErrors ( []FieldError )
	ValidationErr := validate.Struct(&models.GormUser{
		FirstName:                    input.FirstName,
		LastName:                     input.LastName,
		Email:                        input.Email,
		PhoneNumber:                  input.PhoneNumber,
		Password:                     input.Password,
		EmailVerificationCallBackURL: input.EmailVerificationCallBackURL,
	})

	// if validation errors
	if ValidationErr != nil {

		// init a slice of field errors
		var errorsSlice []*FieldErrors

		errs := ValidationErr.(validator.ValidationErrors)

		//  translate each error at a time.
		for _, e := range errs {
			// model the resultant errors to the expected (fieldErrors struct)
			errors := &FieldErrors{
				Field: e.Field(),
				Error: e.Translate(trans),
			}
			// append to the errorsSlice
			errorsSlice = append(errorsSlice, errors)
		}

		// return with validation error
		return &RegisterUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  400,
			Message:     "Input validation errors!",
			FieldErrors: errorsSlice,
		}, nil

	}

	// if no validation errors
	// initialize the db instance
	db := r.DB

	// pointer to base user model
	user := &models.GormUser{}

	// check if user already exists
	if !db.Where("email = ?", input.Email).First(&user).RecordNotFound() {
		// if they do return user already exists
		return &RegisterUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  400,
			Message:     "A user with that email already exists!",
			FieldErrors: nil,
		}, nil
	}

	// if user doesn't already exist

	// generate random uuid
	uuid, UIDGenerationErr := uuid.NewRandom()

	// if there's an error generating uuid
	if UIDGenerationErr != nil {
		// log the error for the backend
		log.Printf("ResolveRegisterUser: uid generation error: %v", UIDGenerationErr)

		// return an error
		return &RegisterUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  500,
			Message:     "Server error, try again!",
			FieldErrors: nil,
		}, nil
	}

	// generate random email verification token
	EmailVerificationToken := utils.RandEmailCodeGenerator()

	if EmailVerificationToken == "" {
		// log for the backend
		log.Println("ResolveRegisterUser: Anonymous func: error generating email string")

		// return an error
		return &RegisterUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  500,
			Message:     "Server error, try again!",
			FieldErrors: nil,
		}, nil

	}

	// format object to save to db
	newUser := &models.GormUser{
		ID:                           uuid.String(),
		FirstName:                    input.FirstName,
		LastName:                     input.LastName,
		Email:                        input.Email,
		PhoneNumber:                  input.PhoneNumber,
		Password:                     input.Password,
		EmailVerificationCallBackURL: input.EmailVerificationCallBackURL,
		EmailVerificationToken:       EmailVerificationToken,
	}

	// hash password
	hashed, hashPassErr := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// if error hashing pass
	if hashPassErr != nil {
		log.Printf("ResolveRegisterUser: password hashing error: %v", hashPassErr)
		// return an error
		return &RegisterUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  500,
			Message:     "Server error, try again!",
			FieldErrors: nil,
		}, nil
	}

	// re-assign pass to hashed
	newUser.Password = string(hashed)

	// try to save the object to the db
	err := db.Create(&newUser).Error

	// if there's an error saving
	if err != nil {
		// log the error for the backend
		log.Printf("ResolveRegisterUser: error saving user: %v", err)

		// return an error
		return &RegisterUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  500,
			Message:     "Server error, try again!",
			FieldErrors: nil,
		}, nil
	}

	// if save successful
	// generate jwt token

	// get JWT Secret from .env
	jwtSecret := os.Getenv("JWT_SECRET")

	// if it returns an empty key
	if jwtSecret == "" {
		// log for the backend
		log.Printf("ResolveRegisterUser: jwt secret returned empty")

		// return but with missing keys
		return &RegisterUserPayload{
			User: &User{
				ID:              newUser.ID,
				FirstName:       newUser.FirstName,
				LastName:        newUser.LastName,
				Email:           newUser.Email,
				PhoneNumber:     newUser.PhoneNumber,
				IsEmailVerified: newUser.IsEmailVerified,
				IsPhoneVerified: newUser.IsPhoneVerified,
				CreatedAt:       newUser.CreatedAt,
				UpdatedAt:       newUser.UpdatedAt,
			},
			JwtToken:    nil,
			StatusCode:  200,
			Message:     "User successfully registered!",
			FieldErrors: nil,
		}, nil

	}

	// Create the JWT key used to create the jwt signature
	var jwtKey = []byte(jwtSecret)

	// Declare the expiration time of the token
	// 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)

	// Create the JWT claims (jwt body)
	// with username, issued at time and expiry time in unix milliseconds
	claims := &models.Claims{
		UserID: newUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string (sign)
	jwtToken, jwtErr := token.SignedString(jwtKey)
	if jwtErr != nil {
		// If there is an error creating the JWT
		// log the error for the backend
		log.Printf("ResolveRegisterUser: error saving user: %v", jwtErr)

		// return but with missing keys
		return &RegisterUserPayload{
			User: &User{
				ID:              newUser.ID,
				FirstName:       newUser.FirstName,
				LastName:        newUser.LastName,
				Email:           newUser.Email,
				PhoneNumber:     newUser.PhoneNumber,
				IsEmailVerified: newUser.IsEmailVerified,
				IsPhoneVerified: newUser.IsPhoneVerified,
				CreatedAt:       newUser.CreatedAt,
				UpdatedAt:       newUser.UpdatedAt,
			},
			JwtToken:    nil,
			StatusCode:  200,
			Message:     "User successfully registered!",
			FieldErrors: nil,
		}, nil

	}

	// Todo:
	// wait for `email verification string` in the email routine with channels
	// phone number validation using reflection
	// format html in html template
	// send the send email fn to the anonymous fn
	// instead of direct access. pinter maybe?

	// try to send user successfully created
	// and email verification
	// in a different go routine (concurrency)
	go func(callBackURL, EmailVerificationToken, firstName string) {

		// composed email verification token
		var EmailVerificationLink string

		// check if callback url ends with a forward slash
		// e.g localhost:3000/
		// and append if not
		if strings.HasSuffix(callBackURL, "/") {
			EmailVerificationLink = callBackURL + EmailVerificationToken
		} else {
			EmailVerificationLink = callBackURL + "/" + EmailVerificationToken
		}

		// basic verification link
		// to be used if html content fails
		plainTextContent := fmt.Sprintf(`You're on your way!. Let's confirm your email address. Copy-paste this link on your browser's tab to verify your email: %v`, EmailVerificationLink)

		// compose clickable verification link
		verifyLink := fmt.Sprintf(`<p><strong><a href="%v" target="_blank" rel="noopener noreferrer"> Click here to verify your email address.</a></strong></p> <p>Or copy-paste the following link on your browser's tab</p> <p><strong><code> %v </code></strong></p>`, EmailVerificationLink, EmailVerificationLink)

		// unified html body content
		// first small case the name then title case
		htmlEmailContent := fmt.Sprintf(`<p>You're on your way!</p> <p>Welcome to our system <strong>%v</strong>. </p> <p>Click the following link to verify your email.</p> %v`, firstName, verifyLink)

		// email subject
		subject := "Welcome to www.muguku.co.ke! Confirm Your Email"

		// try to send the email in a different go routine
		go emails.SendEmail(newUser.Email, subject, plainTextContent, htmlEmailContent)

	}(newUser.EmailVerificationCallBackURL, newUser.EmailVerificationToken, strings.Title(strings.ToLower(newUser.FirstName))) // self invoke

	// if everything goes right return created object
	return &RegisterUserPayload{
		User: &User{
			ID:              newUser.ID,
			FirstName:       newUser.FirstName,
			LastName:        newUser.LastName,
			Email:           newUser.Email,
			PhoneNumber:     newUser.PhoneNumber,
			IsEmailVerified: newUser.IsEmailVerified,
			IsPhoneVerified: newUser.IsPhoneVerified,
			CreatedAt:       newUser.CreatedAt,
			UpdatedAt:       newUser.UpdatedAt,
		},
		JwtToken:    &jwtToken,
		StatusCode:  200,
		Message:     "User successfully registered!",
		FieldErrors: nil,
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) LoginUser(ctx context.Context, input LoginUserInput) (*LoginUserPayload, error) {
	// Todo: modularise the validation function

	// validate input fields

	// english locale
	en := en.New()
	// universal english translator
	uni = ut.New(en, en)

	// translator for english
	// this is usually know/extracted from http 'Accept-Language' header
	trans, _ := uni.GetTranslator("en")

	// initialize validate v10 instance
	validate = validator.New()

	en_translations.RegisterDefaultTranslations(validate, trans)

	// validate against the validation struct
	// returns nil or ValidationErrors ( []FieldError )
	ValidationErr := validate.Struct(&models.ValidateLoginInput{
		Email:    input.Email,
		Password: input.Password,
	})

	// if validation errors
	if ValidationErr != nil {

		// init a slice of field errors
		var errorsSlice []*FieldErrors

		errs := ValidationErr.(validator.ValidationErrors)

		//  translate each error at a time.
		for _, e := range errs {
			// model the resultant errors to the expected (fieldErrors struct)
			errors := &FieldErrors{
				Field: e.Field(),
				Error: e.Translate(trans),
			}
			// append to the errorsSlice
			errorsSlice = append(errorsSlice, errors)
		}

		// return with validation error
		return &LoginUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  400,
			Message:     "Input validation errors!",
			FieldErrors: errorsSlice,
		}, nil

	}

	// if no validation errors
	// initialize the db instance
	db := r.DB

	user := &models.GormUser{}

	// check if user exists
	if db.Where("email = ?", input.Email).First(&user).RecordNotFound() {
		// if they do not return error
		return &LoginUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  400,
			Message:     "No user with those credentials exist!",
			FieldErrors: nil,
		}, nil
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		// If the two passwords don't match, return a 404 status and error
		return &LoginUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  400,
			Message:     "No user with those credentials exist!",
			FieldErrors: nil,
		}, nil
	}

	// if passwords match
	// generate jwt token

	// get JWT Secret from .env
	jwtSecret := os.Getenv("JWT_SECRET")

	// if it returns an empty key
	if jwtSecret == "" {
		// log for the backend
		log.Printf("ResolveLoginUser: jwt secret returned empty")

		// return an error
		return &LoginUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  500,
			Message:     "Server error, try again!",
			FieldErrors: nil,
		}, nil

	}

	// Create the JWT key used to create the jwt signature
	var jwtKey = []byte(jwtSecret)

	// Declare the expiration time of the token
	// 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)

	// Create the JWT claims (jwt body)
	// with username, issued at time and expiry time in unix milliseconds
	claims := &models.Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string (sign)
	jwtToken, jwtErr := token.SignedString(jwtKey)
	if jwtErr != nil {
		// If there is an error creating the JWT
		// log the error for the backend
		log.Printf("ResolveLoginUser: error saving user: %v", jwtErr)
		// return an error
		return &LoginUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  500,
			Message:     "Server error, try again!",
			FieldErrors: nil,
		}, nil

	}

	// if everything goes right return created object
	return &LoginUserPayload{
		User: &User{
			ID:              user.ID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			PhoneNumber:     user.PhoneNumber,
			IsEmailVerified: user.IsEmailVerified,
			IsPhoneVerified: user.IsPhoneVerified,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		},
		JwtToken:    &jwtToken,
		StatusCode:  200,
		Message:     "User login successful!",
		FieldErrors: nil,
	}, nil

}

func (r *mutationResolver) VerifyEmail(ctx context.Context, emailVerificationToken string) (*VerifyEmailPayload, error) {

	// validate for empty or random string of code
	// the email verification token func generates a 62 characters long string
	if len([]rune(emailVerificationToken)) != 62 {
		return &VerifyEmailPayload{
			StatusCode: 400,
			Message:    "bad request, check your input!",
		}, nil
	}

	// if no validation errors
	// initialize the db instance
	db := r.DB

	user := &models.GormUser{}

	// check if user with that verification token exists
	if db.Where("email_verification_token = ?", emailVerificationToken).First(&user).RecordNotFound() {
		// if they don't return error
		return &VerifyEmailPayload{
			StatusCode: 400,
			Message:    "bad request, check your input!",
		}, nil
	}

	// check if user is already verified
	if user.IsEmailVerified == true {
		// if yes return message
		return &VerifyEmailPayload{
			StatusCode: 400,
			Message:    "This email is already verified!",
		}, nil
	}

	// if not verify them (update is email verified flag to true)
	if err := db.Model(&user).Where("email_verification_token = ?", emailVerificationToken).Update("is_email_verified", true).Error; err != nil {
		// error handling
		log.Println("ResolveVerifyEmail: error changing email verification to true")
		// return an error
		return &VerifyEmailPayload{
			StatusCode: 500,
			Message:    "Server error, try again!",
		}, nil
	}

	// try to send phone verification code
	// in a different go routine (concurrency)
	go func() {

		// generate random 6 digit phone verification token
		PhoneVerificationToken := utils.RandPhoneCodeGenerator()

		if PhoneVerificationToken == 0 {
			// log for the backend
			log.Println("ResolveRegisterUser: Anonymous func: error generating phone token")
		}

		// save this token to the db (update)
		if err := db.Model(&user).Where("email = ?", user.Email).Update("phone_verification_token", PhoneVerificationToken).Error; err != nil {
			// error handling
			log.Println("ResolveRegisterUser: Anonymous func: error saving email_verification_token string")
		}

		// composed phone verification sms message
		message := fmt.Sprintf("Your www.muguku.co.ke verification token is: %d", PhoneVerificationToken)

		// receiver
		receiver := user.PhoneNumber

		// try to send the code in a different go routine
		go sms.SendSms(receiver, message)

	}() // self invoke

	return &VerifyEmailPayload{
		StatusCode: 200,
		Message:    "Email successfully verified and Phone verification sent!",
	}, nil
}

func (r *mutationResolver) VerifyPhone(ctx context.Context, phoneVerificationToken string) (*VerifyPhonePayload, error) {
	panic("not implemented")
}
