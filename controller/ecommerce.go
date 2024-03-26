package controller

import (
	"context"
	"gqlserver/graph/model"

	"gorm.io/gorm"
)

func EcommerceProductList(db *gorm.DB,ctx context.Context, limit int, offset int, filter *model.ProductFilter, sort *model.ProductSort) (model.EcommerceProducts, error) {

	var ecom_products []model.EcommerceProduct

	listQuery := db.Debug().Table("tbl_ecom_products").Joins("inner join tbl_categories on tbl_categories.id = tbl_ecom_products.categories_id").Where("tbl_ecom_products.is_deleted = 0 and tbl_categories.parent_id = (select id from tbl_categories where tbl_categories.is_deleted =0 and tbl_categories.parent_id = 0 and tbl_categories.category_slug = 'ecommerce_default_group')")

	if filter.CategoryName != nil{

		listQuery = listQuery.Where("tbl_categories.category_name = ?",filter.CategoryName)

	}else if filter.CategoryID !=nil{
		
		listQuery = listQuery.Where("tbl_categories.id = ?",filter.CategoryID)

	}

	if filter.ReleaseDate!=nil{

		listQuery = listQuery.Where("tbl_ecom_products.created_on >= ?",filter.ReleaseDate)

	}

	if filter.StartingPrice!=nil && filter.EndingPrice!=nil{

		listQuery = listQuery.Where("tbl_ecom_products.price between (?) and (?) ",filter.StartingPrice,filter.EndingPrice)
	}

	var orderBy string

	if *sort.Date == 1{

		orderBy = "tbl_ecom_products.id desc"

	}else if *sort.Price == 1{

		orderBy = "tbl_ecom_products.price desc"
	}

	listQuery = listQuery.Order(orderBy).Limit(limit).Offset(offset).Find(&ecom_products)

	if err := listQuery.Error;err!=nil{

		return model.EcommerceProducts{},err
	}


	return  model.EcommerceProducts{ProductList: ecom_products,Count: len(ecom_products)},nil
}







