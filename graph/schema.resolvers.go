package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	//"fmt"
	"strconv"

	"7linternational.com/gql-server/graph/generated"
	"7linternational.com/gql-server/graph/model"
	"7linternational.com/gql-server/graph/db"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.User, error) {
	user, _, err := db.Login(input.Username, input.Password)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *queryResolver) GetUserRecipes(ctx context.Context, id string) (*model.UserRecipes, error) {
	userId, _ := strconv.ParseInt(id, 10, 64)

	_, err, recipes := db.GetUserRecipes(userId)

	if err != nil {
		return nil, err
	}

	userRecipes := &model.UserRecipes{}

	userRecipes.Hits = recipes

	return userRecipes, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
