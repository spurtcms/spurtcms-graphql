package controller

import (
	"context"
	"log"
	"net/http"
	"spurtcms-graphql/graph/model"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func EcommerceProductList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.ProductFilter, sort *model.ProductSort) (*model.EcommerceProducts, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	var ecom_products []model.EcommerceProduct

	var count int64

	// currentTime := time.Now().In(TimeZone).Format("2006-01-02 15:04:05")

	listQuery := db.Debug().Table("tbl_ecom_products").Select("tbl_ecom_products.*, rp.price AS discount_price ,rs.price AS special_price").Joins("inner join tbl_categories on tbl_categories.id = ANY(STRING_TO_ARRAY(tbl_ecom_products.categories_id," + "','" + ")::INTEGER[])").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='discount' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rp on rp.product_id = tbl_ecom_products.id").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='special' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rs on rs.product_id = tbl_ecom_products.id").Where("tbl_ecom_products.is_deleted = 0 and tbl_ecom_products.is_active = 1")

	if filter != nil {

		if filter.CategoryName.IsSet() {

			listQuery = listQuery.Where("tbl_categories.category_name = ?", filter.CategoryName.Value())

		} else if filter.CategoryID.IsSet() {

			listQuery = listQuery.Where("tbl_categories.id = ?", filter.CategoryID.Value())

		}

		if filter.ReleaseDate.IsSet() {

			listQuery = listQuery.Where("tbl_ecom_products.created_on >= ?", filter.ReleaseDate.Value())

		}

		if filter.StartingPrice.IsSet() && filter.EndingPrice.IsSet() {

			listQuery = listQuery.Where("tbl_ecom_products.product_price >= ?", filter.StartingPrice.Value())

		} else if filter.StartingPrice.IsSet() && filter.EndingPrice.IsSet() {

			listQuery = listQuery.Where("tbl_ecom_products.product_price <= ?", filter.StartingPrice.Value())

		} else if filter.StartingPrice.IsSet() && filter.EndingPrice.IsSet() {

			listQuery = listQuery.Where("tbl_ecom_products.product_price between (?) and (?)", filter.StartingPrice, filter.EndingPrice)
		}
	}

	var orderBy string

	if sort != nil {

		if sort.Date.IsSet() {

			if *sort.Date.Value() == 1 {

				orderBy = "tbl_ecom_products.id desc"

			} else {

				orderBy = ""
			}

		} else if sort.Price.IsSet() {

			if *sort.Price.Value() == 1 {

				orderBy = "tbl_ecom_products.product_price desc"

			} else {

				orderBy = ""
			}

		}

	}

	countQuery := listQuery.Count(&count)

	if err := countQuery.Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.EcommerceProducts{}, err
	}

	listQuery = listQuery.Order(orderBy).Limit(limit).Offset(offset).Find(&ecom_products)

	if err := listQuery.Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.EcommerceProducts{}, err
	}

	return &model.EcommerceProducts{ProductList: ecom_products, Count: int(count)}, nil
}

func EcommerceProductDetails(db *gorm.DB, ctx context.Context, productId int) (*model.EcommerceProduct, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	var productdtl model.EcommerceProduct

	currentTime := time.Now().In(TimeZone).Format("2006-01-02 15:04:05")

	if err := db.Debug().Table("tbl_ecom_products").Select("tbl_ecom_products.*,rp.price AS discount_price ,rs.price AS special_price").Joins("inner join tbl_ecom_product_pricings on tbl_ecom_product_pricings.product_id = tbl_ecom_products.id").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='discount' and tbl_ecom_product_pricings.start_date <= ? and tbl_ecom_product_pricings.end_date >= ?) rp on rp.product_id = tbl_ecom_products.id", currentTime, currentTime).Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='special' and tbl_ecom_product_pricings.start_date <= ? and tbl_ecom_product_pricings.end_date >= ?) rs on rs.product_id = tbl_ecom_products.id", currentTime, currentTime).Where("tbl_ecom_products.is_deleted = 0 and tbl_ecom_products.is_active = 1").Where("tbl_ecom_products.id = ?", productId).First(&productdtl).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.EcommerceProduct{}, err
	}

	return &productdtl, nil

}

func EcommerceAddToCart(db *gorm.DB, ctx context.Context, productID int, customerID int, quantity int) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	var cart model.EcommerceCart

	cart.ProductID = productID

	cart.CustomerID = customerID

	cart.Quantity = quantity

	cart.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	if err := db.Debug().Table("tbl_ecom_carts").Create(&cart).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	return true, nil
}

func EcommerceCartList(db *gorm.DB, ctx context.Context, customerId int) (*model.EcommerceCartDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	log.Println("customer_id", customerId)

	var cartProductList []model.EcommerceProduct

	if err := db.Debug().Table("tbl_ecom_carts").Select("tbl_ecom_carts.*,tbl_ecom_products.*,SUM(tbl_ecom_carts.id)").Joins("inner join tbl_ecom_products on tbl_ecom_products.id = tbl_ecom_carts.product_id").Joins("inner join tbl_ecom_customers on tbl_ecom_customers.id = tbl_ecom_carts.customer_id").
		Where("tbl_ecom_carts.is_deleted = 0 and tbl_ecom_products.is_deleted = 0 and tbl_ecom_customers.is_deleted = 0").Preload("EcommerceCart").Find(&cartProductList).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.EcommerceCartDetails{}, err
	}

	log.Println("cartList", cartProductList)


	return &model.EcommerceCartDetails{CartList: cartProductList, OrderSummary: model.OrderSummary{}, Count: len(cartProductList)}, nil
}
