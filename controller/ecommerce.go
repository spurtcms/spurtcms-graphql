package controller

import (
	"context"
	"gqlserver/graph/model"
	// "log"

	"gorm.io/gorm"
)

func EcommerceProductList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.ProductFilter, sort *model.ProductSort) (model.EcommerceProducts, error) {

	var ecom_products []model.EcommerceProduct

	var count int64

	listQuery := db.Debug().Table("tbl_ecom_products").Joins("inner join tbl_categories on tbl_categories.id = ANY(STRING_TO_ARRAY(tbl_ecom_products.categories_id," + "','" + ")::INTEGER[])").Where("tbl_ecom_products.is_deleted = 0 and tbl_ecom_products.is_active = 1 and tbl_categories.parent_id = (select id from tbl_categories where tbl_categories.is_deleted =0 and tbl_categories.parent_id = 0 and tbl_categories.category_slug = 'ecommerce_default_group')")

	if filter != nil {

		if filter.CategoryName != nil {

			listQuery = listQuery.Where("tbl_categories.category_name = ?", *filter.CategoryName)

		} else if filter.CategoryID != nil {

			listQuery = listQuery.Where("tbl_categories.id = ?", *filter.CategoryID)

		}

		if filter.ReleaseDate != nil {

			listQuery = listQuery.Where("tbl_ecom_products.created_on >= ?", *filter.ReleaseDate)

		}

		if filter.StartingPrice != nil && filter.EndingPrice == nil{

			listQuery = listQuery.Where("tbl_ecom_products.product_price >= ?", *filter.StartingPrice)

		}else if filter.StartingPrice == nil && filter.EndingPrice != nil{

			listQuery = listQuery.Where("tbl_ecom_products.product_price <= ?", *filter.StartingPrice)

		}else if filter.StartingPrice != nil && filter.EndingPrice != nil {

			listQuery = listQuery.Where("tbl_ecom_products.product_price between (?) and (?)", *filter.StartingPrice, *filter.EndingPrice)
		}
	}

	var orderBy string

	if sort != nil {

		if sort.Date != nil {
			
			if *sort.Date == 1{

				orderBy = "tbl_ecom_products.id desc"

			}else{

				orderBy = ""
			}

		} else if sort.Price != nil {

			if *sort.Price == 1{

				orderBy = "tbl_ecom_products.product_price desc"

			}else{

				orderBy = ""
			}

		}

	}

	countQuery := listQuery.Count(&count)

	if err := countQuery.Error; err != nil {

		return model.EcommerceProducts{}, err
	}

	listQuery = listQuery.Order(orderBy).Limit(limit).Offset(offset).Find(&ecom_products)

	if err := listQuery.Error; err != nil {

		return model.EcommerceProducts{}, err
	}

	return model.EcommerceProducts{ProductList: ecom_products, Count: int(count)}, nil
}
