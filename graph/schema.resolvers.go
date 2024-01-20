package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"gqlserver/graph/model"
)

// MemberLogin is the resolver for the memberLogin field.
func (r *mutationResolver) MemberLogin(ctx context.Context, input model.LoginCredentials) (string, error) {
	return MemberLogin(r.DB, input)
}

// MemberRegister is the resolver for the memberRegister field.
func (r *mutationResolver) MemberRegister(ctx context.Context, input model.MemberDetails) (bool, error) {
	return MemberRegister(r.DB, input)
}

// MemberUpdate is the resolver for the memberUpdate field.
func (r *mutationResolver) MemberUpdate(ctx context.Context, memberdata model.MemberDetails) (bool, error) {
	return MemberUpdate(r.DB, ctx, memberdata)
}

// SendOtpToMail is the resolver for the sendOtpToMail field.
func (r *mutationResolver) SendOtpToMail(ctx context.Context, email string) (bool, error) {
	return SendOtpToMail(r.DB, ctx, email)
}

// ResetPassword is the resolver for the resetPassword field.
func (r *mutationResolver) ResetPassword(ctx context.Context, otp int, newPassword string, email string) (bool, error) {
	return ResetPassword(r.DB, otp, newPassword, email)
}

// Channellist is the resolver for the channellist field.
func (r *queryResolver) ChannelList(ctx context.Context, limit int, offset int) (model.ChannelDetails, error) {
	return Channellist(r.DB, ctx, limit, offset)
}

// ChannelEntriesList is the resolver for the channelEntriesList field.
func (r *queryResolver) ChannelEntriesList(ctx context.Context, channelID *int, channelEntryID *int, limit *int, offset *int) (model.ChannelEntryDetails, error) {
	return ChannelEntriesList(r.DB, ctx, channelID, channelEntryID, limit, offset)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
