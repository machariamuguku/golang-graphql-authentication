package resolvers

//go:generate go run github.com/99designs/gqlgen

import (
	"context"

	"github.com/machariamuguku/golang-graphql-authentication/db"

	golang_graphql_authentication "github.com/machariamuguku/golang-graphql-authentication"
)

// Resolver : pointer to db
type Resolver struct {
	DB *db.DB
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

// Mutation : does something
func (r *Resolver) Mutation() golang_graphql_authentication.MutationResolver {
	return &mutationResolver{r}
}

// Query : does something
func (r *Resolver) Query() golang_graphql_authentication.QueryResolver {
	return &queryResolver{r}
}

func (r *queryResolver) LoginUser(ctx context.Context, input golang_graphql_authentication.LoginUserInput) (*golang_graphql_authentication.LoginUserPayload, error) {

	return LoginUserQuery(ctx, input, r)
}

func (r *mutationResolver) RegisterUser(ctx context.Context, input golang_graphql_authentication.RegisterUserInput) (*golang_graphql_authentication.RegisterUserPayload, error) {
	return RegisterUserMutation(ctx, input, r)

}

func (r *mutationResolver) VerifyEmail(ctx context.Context, emailVerificationToken string) (*golang_graphql_authentication.VerifyEmailPayload, error) {
	return VerifyEmailMutation(ctx, emailVerificationToken, r)
}

func (r *mutationResolver) VerifyPhone(ctx context.Context, phoneVerificationToken int) (*golang_graphql_authentication.VerifyPhonePayload, error) {
	return VerifyPhoneMutation(ctx, phoneVerificationToken, r)
}
