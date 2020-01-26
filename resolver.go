package golang_graphql_authentication

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	users []*User
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) RegisterUser(ctx context.Context, input RegisterUserInput) (*RegisterUserPayload, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) LoginUser(ctx context.Context, input LoginUserInput) (*LoginUserPayload, error) {
	// resolve log in user
	return ResolveLoginUser(ctx, input)
}
func (r *queryResolver) Users(ctx context.Context) ([]*UserPayload, error) {
	panic("not implemented")
}
func (r *queryResolver) User(ctx context.Context, id string) (*UserPayload, error) {
	panic("not implemented")
}
