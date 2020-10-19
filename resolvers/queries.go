package resolvers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	golang_graphql_authentication "github.com/machariamuguku/golang-graphql-authentication"
	"github.com/machariamuguku/golang-graphql-authentication/models"
	"golang.org/x/crypto/bcrypt"
)

// LoginUserQuery : Validate existing user
func LoginUserQuery(ctx context.Context, input golang_graphql_authentication.LoginUserInput, r *queryResolver) (*golang_graphql_authentication.LoginUserPayload, error) {

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
		var errorsSlice []*golang_graphql_authentication.FieldErrors

		errs := ValidationErr.(validator.ValidationErrors)

		//  translate each error at a time.
		for _, e := range errs {
			// model the resultant errors to the expected (golang_graphql_authentication.FieldErrors struct)
			errors := &golang_graphql_authentication.FieldErrors{
				Field: e.Field(),
				Error: e.Translate(trans),
			}
			// append to the errorsSlice
			errorsSlice = append(errorsSlice, errors)
		}

		// return with validation error
		return &golang_graphql_authentication.LoginUserPayload{
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
		return &golang_graphql_authentication.LoginUserPayload{
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
		return &golang_graphql_authentication.LoginUserPayload{
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
		return &golang_graphql_authentication.LoginUserPayload{
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
		return &golang_graphql_authentication.LoginUserPayload{
			User:        nil,
			JwtToken:    nil,
			StatusCode:  500,
			Message:     "Server error, try again!",
			FieldErrors: nil,
		}, nil

	}

	// if everything goes right return created object
	return &golang_graphql_authentication.LoginUserPayload{
		User: &golang_graphql_authentication.User{
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
