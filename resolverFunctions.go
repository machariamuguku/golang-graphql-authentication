package golang_graphql_authentication

import (
	"context"
	// "fmt"
	"errors"
	"log"
	"os"
)

// Logger function
func Logger() {
	// log file name
	fileName := "webrequests.log"

	// open the file and give write permissions
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	// if there's error opening/creating file
	if err != nil {
		panic(err)
	}

	// defer closing the file
	defer logFile.Close()

	// direct all log messages to webrequests.log
	log.SetOutput(logFile)
}

//ResolveLoginUser used toresolve user login
func ResolveLoginUser(ctx context.Context, input LoginUserInput) (*LoginUserPayload, error) {

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
