package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"gqlserver/graph/model"
)

// MemberLogin is the resolver for the memberLogin field.
func (r *mutationResolver) MemberLogin(ctx context.Context, email string) (bool, error) {
	return MemberLogin(r.DB, ctx, email)
}

// VerifyMemberOtp is the resolver for the verifyMemberOtp field.
func (r *mutationResolver) VerifyMemberOtp(ctx context.Context, otp int) (string, error) {
	return VerifyMemberOtp(r.DB, ctx, otp)
}

// MemberRegister is the resolver for the memberRegister field.
func (r *mutationResolver) MemberRegister(ctx context.Context, input model.MemberDetails) (bool, error) {
	return MemberRegister(r.DB, input)
}

// MemberUpdate is the resolver for the memberUpdate field.
func (r *mutationResolver) MemberUpdate(ctx context.Context, memberdata model.MemberDetails) (bool, error) {
	return MemberUpdate(r.DB, ctx, memberdata)
}

// MemberProfileUpdate is the resolver for the memberProfileUpdate field.
func (r *mutationResolver) MemberProfileUpdate(ctx context.Context, profiledata model.ProfileData, entryID int, updateExactMemberProfileOnly bool) (bool, error) {
	return MemberProfileUpdate(r.DB, ctx, profiledata, entryID, updateExactMemberProfileOnly)
}

// Memberclaimnow is the resolver for the memberclaimnow field.
func (r *mutationResolver) Memberclaimnow(ctx context.Context, input model.ClaimData, entryID int) (bool, error) {
	return Memberclaimnow(r.DB, ctx, input, entryID)
}

// Channellist is the resolver for the channellist field.
func (r *queryResolver) ChannelList(ctx context.Context, limit int, offset int) (model.ChannelDetails, error) {
	return Channellist(r.DB, ctx, limit, offset)
}

// ChannelDetail is the resolver for the channelDetail field.
func (r *queryResolver) ChannelDetail(ctx context.Context, channelID int) (model.Channel, error) {
	return ChannelDetail(r.DB, ctx, channelID)
}

// ChannelEntriesList is the resolver for the channelEntriesList field.
func (r *queryResolver) ChannelEntriesList(ctx context.Context, channelID *int, categoryID *int, limit int, offset int, title *string) (model.ChannelEntriesDetails, error) {
	return ChannelEntriesList(r.DB, ctx, channelID, categoryID, limit, offset, title)
}

// ChannelEntryDetail is the resolver for the channelEntryDetail field.
func (r *queryResolver) ChannelEntryDetail(ctx context.Context, categoryID *int, channelID *int, channelEntryID *int, slug *string) (model.ChannelEntries, error) {
	return ChannelEntryDetail(r.DB, ctx, channelID, categoryID, channelEntryID, slug)
}

// SpaceList is the resolver for the spaceList field.
func (r *queryResolver) SpaceList(ctx context.Context, limit int, offset int) (model.SpaceDetails, error) {
	return SpaceList(r.DB, ctx, limit, offset)
}

// SpaceDetails is the resolver for the spaceDetails field.
func (r *queryResolver) SpaceDetails(ctx context.Context, spaceID int) (model.Space, error) {
	return SpaceDetails(r.DB, ctx, spaceID)
}

// PagesAndPageGroupsUnderSpace is the resolver for the PagesAndPageGroupsUnderSpace field.
func (r *queryResolver) PagesAndPageGroupsUnderSpace(ctx context.Context, spaceID int) (model.PageAndPageGroups, error) {
	return PagesAndPageGroupsUnderSpace(r.DB, ctx, spaceID)
}

// CategoriesList is the resolver for the categoriesList field.
func (r *queryResolver) CategoriesList(ctx context.Context, limit *int, offset *int, categoryGroupID *int, hierarchyLevel *int) (model.CategoriesList, error) {
	return CategoriesList(r.DB, ctx, limit, offset, categoryGroupID, hierarchyLevel)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
