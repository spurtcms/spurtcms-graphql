package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"spurtcms-graphql/graph/model"
)

// MemberLogin is the resolver for the memberLogin field.
func (r *mutationResolver) MemberLogin(ctx context.Context, email string) (bool, error) {
	return MemberLogin(r.DB, ctx, email)
}

// VerifyMemberOtp is the resolver for the verifyMemberOtp field.
func (r *mutationResolver) VerifyMemberOtp(ctx context.Context, email string, otp int) (*model.LoginDetails, error) {
	return VerifyMemberOtp(r.DB, ctx, email, otp)
}

// MemberProfileUpdate is the resolver for the memberProfileUpdate field.
func (r *mutationResolver) MemberProfileUpdate(ctx context.Context, profiledata model.ProfileData, entryID int) (bool, error) {
	return MemberProfileUpdate(r.DB, ctx, profiledata, entryID)
}

// Memberclaimnow is the resolver for the memberclaimnow field.
func (r *mutationResolver) Memberclaimnow(ctx context.Context, input model.ClaimData, entryID int) (bool, error) {
	return Memberclaimnow(r.DB, ctx, input, entryID)
}

// ProfileNameVerification is the resolver for the profileNameVerification field.
func (r *mutationResolver) ProfileNameVerification(ctx context.Context, profileName string) (bool, error) {
	return VerifyProfileName(r.DB, ctx, profileName)
}

// UpdateChannelEntryViewCount is the resolver for the updateChannelEntryViewCount field.
func (r *mutationResolver) UpdateChannelEntryViewCount(ctx context.Context, entryID *int, slug *string) (bool, error) {
	return UpdateChannelEntryViewCount(r.DB, ctx,entryID,slug)
}

// Channellist is the resolver for the channellist field.
func (r *queryResolver) ChannelList(ctx context.Context, limit int, offset int) (*model.ChannelDetails, error) {
	return Channellist(r.DB, ctx, limit, offset)
}

// ChannelDetail is the resolver for the channelDetail field.
func (r *queryResolver) ChannelDetail(ctx context.Context, channelID int) (*model.Channel, error) {
	return ChannelDetail(r.DB, ctx, channelID)
}

// ChannelEntriesList is the resolver for the channelEntriesList field.
func (r *queryResolver) ChannelEntriesList(ctx context.Context, channelID *int, categoryID *int, limit int, offset int, title *string, categoryChildID *int, categorySlug *string, categoryChildSlug *string) (*model.ChannelEntriesDetails, error) {
	return ChannelEntriesList(r.DB, ctx, channelID, categoryID, limit, offset, title, categoryChildID, categorySlug, categoryChildSlug)
}

// ChannelEntryDetail is the resolver for the channelEntryDetail field.
func (r *queryResolver) ChannelEntryDetail(ctx context.Context, categoryID *int, channelID *int, channelEntryID *int, slug *string, categoryChildID *int, profileSlug *string) (*model.ChannelEntries, error) {
	return ChannelEntryDetail(r.DB, ctx, channelID, categoryID, channelEntryID, slug, profileSlug)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
