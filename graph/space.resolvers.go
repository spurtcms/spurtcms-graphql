package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"spurtcms-graphql/graph/model"
)

// SpaceList is the resolver for the spaceList field.
func (r *queryResolver) SpaceList(ctx context.Context, limit int, offset int, categoriesID *int) (*model.SpaceDetails, error) {
	return SpaceList(r.DB, ctx, limit, offset, categoriesID)
}

// SpaceDetails is the resolver for the spaceDetails field.
func (r *queryResolver) SpaceDetails(ctx context.Context, spaceID *int, spaceSlug *string) (*model.Space, error) {
	return SpaceDetails(r.DB, ctx, spaceID, spaceSlug)
}

// PagesAndPageGroupsUnderSpace is the resolver for the PagesAndPageGroupsUnderSpace field.
func (r *queryResolver) PagesAndPageGroupsUnderSpace(ctx context.Context, spaceID int) (*model.PageAndPageGroups, error) {
	return PagesAndPageGroupsUnderSpace(r.DB, ctx, spaceID)
}
