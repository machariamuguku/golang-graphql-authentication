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
	users []*User
	DB    *db.DB
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *Resolver) RegisterUser(ctx context.Context, input RegisterUserInput) (*RegisterUserPayload, error) {

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

	// check if user already exists
	if !db.Where("email = ?", input.Email).First(&models.GormUser{}).RecordNotFound() {
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

	// format object to save to db
	newUser := &models.GormUser{
		ID:                           uuid.String(),
		FirstName:                    input.FirstName,
		LastName:                     input.LastName,
		Email:                        input.Email,
		PhoneNumber:                  input.PhoneNumber,
		Password:                     input.Password,
		EmailVerificationCallBackURL: input.EmailVerificationCallBackURL,
	}

	// hash password
	hashed, hashPassErr := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// if error hashing pass
	if hashPassErr != nil {
		log.Printf("ResolveRegisterUser: password hashing error: %v", hashPassErr)
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
		// then set another placeholder jwt secret
		jwtSecret = "a_very_!@#$%^&_secret"
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

	}

	// Todo:
	// wait for `random verification string` in the email routine with channels
	// phone number validation using reflection
	// format html in html template

	// try to send user successfully created
	// and email verification
	// in a different go routine (concurrency)
	go func(callBackURL string) {
		// random email verification token
		randomString := utils.RandGenerator()

		if randomString == "" {
			// log for the backend
			log.Println("ResolveRegisterUser: Anonymous func: error generating random string")
		}

		// composed random verification token
		var randomVerifyToken string

		// check if callback url ends with a forward slash
		// e.g localhost:3000/
		// and append if not
		if strings.HasSuffix(callBackURL, "/") {
			randomVerifyToken = callBackURL + randomString
		} else {
			randomVerifyToken = callBackURL + "/" + randomString
		}

		// clickable verification link
		verifyLink := fmt.Sprintf(`<a href="%v"> Click here to verify your email address.</a> <p> Or copy-paste this link on your browser tab <strong> %v </strong>`, randomVerifyToken, randomVerifyToken)

		// unified html body content
		emailContent := fmt.Sprintf(`<p>You're on your way! Let's confirm your email address. By clicking on the following link, you are confirming your email address.</p> %v`, verifyLink)

		// email subject
		subject := "Welcome to www.muguku.co.ke! Confirm Your Email"

		// try to send the email in a different go routine
		go emails.SendEmail(newUser.Email, subject, emailContent)

	}(newUser.EmailVerificationCallBackURL) // self invoke

	// if everything goes right return created object
	return &RegisterUserPayload{
		User: &User{
			ID:          newUser.ID,
			FirstName:   newUser.FirstName,
			LastName:    newUser.LastName,
			Email:       newUser.Email,
			PhoneNumber: newUser.PhoneNumber,
		},
		JwtToken:    &jwtToken,
		StatusCode:  200,
		Message:     "User successfully registered!",
		FieldErrors: nil,
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) LoginUser(ctx context.Context, input LoginUserInput) (*LoginUserPayload, error) {
	panic("not implemented")
}
func (r *queryResolver) Users(ctx context.Context) ([]*UserPayload, error) {
	panic("not implemented")
}
func (r *queryResolver) User(ctx context.Context, id string) (*UserPayload, error) {
	panic("not implemented")
}
