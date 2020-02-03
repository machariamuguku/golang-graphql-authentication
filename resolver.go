package golang_graphql_authentication

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/machariamuguku/golang-graphql-authentication/db"
	"github.com/machariamuguku/golang-graphql-authentication/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

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

	// initialize the db instance
	db := r.DB

	// generate random uuid
	uuid, UIDGenerationErr := uuid.NewRandom()

	// if there's an error generating uuid log error
	if UIDGenerationErr != nil {
		// log the error for the backend
		log.Printf("ResolveRegisterUser: uid generation error: %v", UIDGenerationErr)

		// return an error
		return &RegisterUserPayload{
			User:       nil,
			JwtToken:   nil,
			StatusCode: "500",
			Message:    "Server error, try again!",
		}, nil
	}

	// generate jwt Token
	jwtToken := "001t0o3nmate"

	// format object to save to db
	newUser := &models.GormUser{
		ID:          uuid.String(),
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Password:    input.Password,
	}

	// hash the password
	hashed, hashPassErr := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// error hashing pass
	if hashPassErr != nil {
		log.Printf("ResolveRegisterUser: password hashing error: %v", hashPassErr)
	}

	// re-assign pass to hashed
	newUser.Password = string(hashed)

	// check if user already exists
	if !db.Where("email = ?", input.Email).First(&models.GormUser{}).RecordNotFound() {
		// return user already exists
		return &RegisterUserPayload{
			User:       nil,
			JwtToken:   nil,
			StatusCode: "400",
			Message:    "user with that email already exists!",
		}, nil
	}

	// if not save the object to the db
	err := db.Create(&newUser).Error

	// if there's an error saving return an error
	if err != nil {
		// log the error for the backend
		log.Printf("ResolveRegisterUser: error saving user: %v", err)

		// return an error
		return &RegisterUserPayload{
			User:       nil,
			JwtToken:   nil,
			StatusCode: "500",
			Message:    "Server error, try again!",
		}, nil
	}

	// format return object
	Register := &RegisterUserPayload{
		User: &User{
			ID:          newUser.ID,
			FirstName:   newUser.FirstName,
			LastName:    newUser.LastName,
			Email:       newUser.Email,
			PhoneNumber: newUser.PhoneNumber,
		},
		JwtToken:   &jwtToken,
		StatusCode: "200",
		Message:    "this shit succeeded!",
	}

	// return created object
	return Register, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) LoginUser(ctx context.Context, input LoginUserInput) (*LoginUserPayload, error) {
	// the jwtToken
	jwtToken := "001t0o3nmate"

	Login := &LoginUserPayload{
		User: &User{
			ID:          "001",
			FirstName:   "njoroge",
			LastName:    "kaihu",
			Email:       input.Email,
			PhoneNumber: input.PhoneNumber,
		},
		JwtToken:   &jwtToken,
		StatusCode: "200",
		Message:    "this shit succeeded",
	}

	err := errors.New("Sorry, the user couldn't be logged in")

	// if error: log error to std-io
	// return status code and message for front end
	if err != nil {
		// log the error for use in the backend
		// starting with the function name
		log.Printf("ResolveLoginUser: %v", err)

		// return nulls
		failed := &LoginUserPayload{
			StatusCode: "400",
			Message:    "Sorry, the user couldn't be logged in",
		}

		return failed, nil
	}

	return Login, nil
}
func (r *queryResolver) Users(ctx context.Context) ([]*UserPayload, error) {
	panic("not implemented")
}
func (r *queryResolver) User(ctx context.Context, id string) (*UserPayload, error) {
	panic("not implemented")
}
