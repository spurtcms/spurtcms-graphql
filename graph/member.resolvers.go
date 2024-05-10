package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"spurtcms-graphql/graph/model"
)

// TemplateMemberLogin is the resolver for the templateMemberLogin field.
func (r *mutationResolver) TemplateMemberLogin(ctx context.Context, username *string, email *string, password string) (string, error) {
	return TemplateMemberLogin(r.DB, ctx, username, email, password)
}

// MemberRegister is the resolver for the memberRegister field.
func (r *mutationResolver) MemberRegister(ctx context.Context, input model.MemberDetails, ecomModule *int) (bool, error) {
	return MemberRegister(r.DB, ctx, input, ecomModule)
}

// MemberUpdate is the resolver for the memberUpdate field.
func (r *mutationResolver) MemberUpdate(ctx context.Context, memberdata model.MemberDetails) (bool, error) {
	return MemberUpdate(r.DB, ctx, memberdata)
}

// MemberProfileDetails is the resolver for the memberProfileDetails field.
func (r *queryResolver) MemberProfileDetails(ctx context.Context) (*model.MemberProfile, error) {
	return MemberProfileDetails(r.DB, ctx)
}
