package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"spurtcms-graphql/graph/model"
)

// CategoriesList is the resolver for the categoriesList field.
func (r *queryResolver) CategoriesList(ctx context.Context, limit *int, offset *int, categoryGroupID *int, hierarchyLevel *int, checkEntriesPresence *int) (*model.CategoriesList, error) {
	return CategoriesList(r.DB, ctx, limit, offset, categoryGroupID, hierarchyLevel, checkEntriesPresence)
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }