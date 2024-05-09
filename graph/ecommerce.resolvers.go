package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"spurtcms-graphql/graph/model"
)

// EcommerceAddToCart is the resolver for the ecommerceAddToCart field.
func (r *mutationResolver) EcommerceAddToCart(ctx context.Context, productID *int, productSlug *string, quantity int) (bool, error) {
	return EcommerceAddToCart(r.DB, ctx, productID, productSlug, quantity)
}

// EcommerceOrderPlacement is the resolver for the ecommerceOrderPlacement field.
func (r *mutationResolver) EcommerceOrderPlacement(ctx context.Context, paymentMode string, shippingAddress string, orderProducts []model.OrderProduct, orderSummary *model.OrderSummary) (bool, error) {
	return EcommerceOrderPlacement(r.DB, ctx, paymentMode, shippingAddress, orderProducts, orderSummary)
}

// RemoveProductFromCartlist is the resolver for the removeProductFromCartlist field.
func (r *mutationResolver) RemoveProductFromCartlist(ctx context.Context, productID int) (bool, error) {
	return RemoveProductFromCartlist(r.DB, ctx, productID)
}

// EcommerceProductList is the resolver for the ecommerceProductList field.
func (r *queryResolver) EcommerceProductList(ctx context.Context, limit int, offset int, filter *model.ProductFilter, sort *model.ProductSort) (*model.EcommerceProducts, error) {
	return EcommerceProductList(r.DB, ctx, limit, offset, filter, sort)
}

// EcommerceProductDetails is the resolver for the ecommerceProductDetails field.
func (r *queryResolver) EcommerceProductDetails(ctx context.Context, productID *int, productSlug *string) (*model.EcommerceProduct, error) {
	return EcommerceProductDetails(r.DB, ctx, productID, productSlug)
}

// EcommerceCartList is the resolver for the ecommerceCartList field.
func (r *queryResolver) EcommerceCartList(ctx context.Context, limit int, offset int) (*model.EcommerceCartDetails, error) {
	return EcommerceCartList(r.DB, ctx, limit, offset)
}

// EcommerceProductOrdersList is the resolver for the ecommerceProductOrdersList field.
func (r *queryResolver) EcommerceProductOrdersList(ctx context.Context, limit int, offset int, filter *model.OrderFilter, sort *model.OrderSort) (*model.EcommerceProducts, error) {
	return EcommerceProductOrdersList(r.DB, ctx, limit, offset, filter, sort)
}

// EcommerceProductOrderDetails is the resolver for the ecommerceProductOrderDetails field.
func (r *queryResolver) EcommerceProductOrderDetails(ctx context.Context, productID *int, productSlug *string) (*model.EcommerceProduct, error) {
	return EcommerceProductOrderDetails(r.DB, ctx, productID, productSlug)
}

// EcommerceCustomerDetails is the resolver for the ecommerceCustomerDetails field.
func (r *queryResolver) EcommerceCustomerDetails(ctx context.Context) (*model.CustomerDetails, error) {
	return EcommerceCustomerDetails(r.DB, ctx)
}
