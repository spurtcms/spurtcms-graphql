package controller

import (
	"context"
	"errors"
	"net/http"
	"spurtcms-graphql/graph/model"
	"strconv"
	"strings"
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

		if filter.StartingPrice.IsSet() && !filter.EndingPrice.IsSet() {

			listQuery = listQuery.Where("tbl_ecom_products.product_price >= ?", filter.StartingPrice.Value())

		} else if !filter.StartingPrice.IsSet() && filter.EndingPrice.IsSet() {

			listQuery = listQuery.Where("tbl_ecom_products.product_price <= ?", filter.EndingPrice.Value())

		} else if filter.StartingPrice.IsSet() && filter.EndingPrice.IsSet() {

			listQuery = listQuery.Where("tbl_ecom_products.product_price between (?) and (?)", filter.StartingPrice.Value(), filter.EndingPrice.Value())
		}

		if filter.SearchKeyword.IsSet() {

			listQuery = listQuery.Where("LOWER(TRIM(tbl_ecom_products.product_name)) ILIKE LOWER(TRIM(?))", "%"+*filter.SearchKeyword.Value()+"%")
		}
	}

	if sort != nil {

		if sort.Date.IsSet() && !sort.Price.IsSet() {

			if *sort.Date.Value() == 1 {

				listQuery = listQuery.Order("tbl_ecom_products.id desc")

			} else {

				listQuery = listQuery.Order("tbl_ecom_products.id ")
			}

		} else if sort.Price.IsSet() && !sort.Date.IsSet() {

			if *sort.Price.Value() == 1 {

				listQuery = listQuery.Order("tbl_ecom_products.product_price desc")

			} else {

				listQuery = listQuery.Order("tbl_ecom_products.product_price")

			}

		}
	} else {

		listQuery = listQuery.Order("tbl_ecom_products.id desc")
	}

	countQuery := listQuery.Count(&count)

	if err := countQuery.Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.EcommerceProducts{}, err
	}

	listQuery = listQuery.Limit(limit).Offset(offset).Find(&ecom_products)

	if err := listQuery.Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.EcommerceProducts{}, err
	}

	var final_ecomProducts []model.EcommerceProduct

	for _,product := range ecom_products{

		if product.ProductImagePath != "" {

			imagePaths := strings.Split(product.ProductImagePath, ",")

			for index, path := range imagePaths {

				modified_path := PathUrl + strings.TrimPrefix(path, "/")

				imagePaths[index] = modified_path
			}

			product.ProductImageArray = imagePaths

		}

		if product.ProductVideoPath != "" {

			modified_path := PathUrl + strings.TrimPrefix(product.ProductVideoPath, "/")

			product.ProductVideoPath = modified_path
		}

		final_ecomProducts = append(final_ecomProducts,product)
	}

	return &model.EcommerceProducts{ProductList: final_ecomProducts, Count: int(count)}, nil
}

func EcommerceProductDetails(db *gorm.DB, ctx context.Context, productId *int, productSlug *string) (*model.EcommerceProduct, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	var productdtl model.EcommerceProduct

	// currentTime := time.Now().In(TimeZone).Format("2006-01-02 15:04:05")

	query := db.Debug().Table("tbl_ecom_products").Select("tbl_ecom_products.*,rp.price AS discount_price ,rs.price AS special_price").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='discount' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rp on rp.product_id = tbl_ecom_products.id").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='special' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rs on rs.product_id = tbl_ecom_products.id").Where("tbl_ecom_products.is_deleted = 0 and tbl_ecom_products.is_active = 1")

	if productId != nil {

		query = query.Where("tbl_ecom_products.id = ?", *productId)

	} else if productSlug != nil {

		query = query.Where("tbl_ecom_products.product_slug = ?", *productSlug)
	}

	if err := query.First(&productdtl).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.EcommerceProduct{}, err
	}

	if productdtl.ProductImagePath != "" {

		imagePaths := strings.Split(productdtl.ProductImagePath, ",")

		for index, path := range imagePaths {

			modified_path := PathUrl + strings.TrimPrefix(path, "/")

			imagePaths[index] = modified_path
		}

		productdtl.ProductImageArray = imagePaths

	}

	if productdtl.ProductVideoPath != "" {

		modified_path := PathUrl + strings.TrimPrefix(productdtl.ProductVideoPath, "/")

		productdtl.ProductVideoPath = modified_path
	}

	return &productdtl, nil

}

func EcommerceAddToCart(db *gorm.DB, ctx context.Context, productID *int, productSlug *string, quantity int) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return false, err

	}

	var cart model.EcommerceCart

	var productId int

	if productID != nil {

		productId = *productID

	} else if productSlug != nil {

		if err := db.Table("tbl_ecom_products").Select("id").Where("is_deleted = 0 and product_slug = ?", *productSlug).Scan(&productId).Error; err != nil {

			c.AbortWithError(500, err)

			return false, err
		}
	}

	var customer_id int

	if err := db.Table("tbl_ecom_customers").Select("tbl_ecom_customers.id").Where("tbl_ecom_customers.is_deleted = 0 and tbl_ecom_customers.member_id = ?", memberid).Scan(&customer_id).Error; err != nil {

		c.AbortWithError(500, err)

		return false, err
	}

	if customer_id == 0 {

		err := errors.New("customer id not found")

		c.AbortWithError(500, err)

		return false, err
	}

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	cart.ProductID = productId

	cart.CustomerID = customer_id

	cart.Quantity = quantity

	cart.IsDeleted = 0

	cart.CreatedOn = currentTime

	var count int64

	if err := db.Table("tbl_ecom_carts").Where("is_deleted = 0 and customer_id = ? and product_id = ?", customer_id, productId).Count(&count).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	query := db.Table("tbl_ecom_carts")

	if count != 0 {

		query = query.Where("is_deleted = 0 and customer_id = ? and product_id = ?", customer_id, productId).UpdateColumns(map[string]interface{}{"quantity": gorm.Expr("quantity + ?",cart.Quantity), "modified_on": currentTime})

	} else {

		query = query.Create(&cart)
	}

	if err := query.Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	return true, nil
}

func EcommerceCartList(db *gorm.DB, ctx context.Context, limit, offset int) (*model.EcommerceCartDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return &model.EcommerceCartDetails{}, err

	}

	var customer_id int

	if err := db.Table("tbl_ecom_customers").Select("tbl_ecom_customers.id").Where("tbl_ecom_customers.is_deleted = 0 and tbl_ecom_customers.member_id = ?", memberid).Scan(&customer_id).Error; err != nil {

		c.AbortWithError(500, err)

		return &model.EcommerceCartDetails{}, err
	}

	if customer_id == 0 {

		err := errors.New("customer id not found")

		c.AbortWithError(500, err)

		return &model.EcommerceCartDetails{}, err
	}

	var cartList []model.EcommerceProduct

	var count int64

	if err := db.Table("tbl_ecom_products").Select("tbl_ecom_carts.*,tbl_ecom_products.*,rp.price AS discount_price ,rs.price AS special_price").Joins("inner join tbl_ecom_carts on tbl_ecom_carts.product_id =  tbl_ecom_products.id ").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='discount' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rp on rp.product_id = tbl_ecom_products.id").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='special' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rs on rs.product_id = tbl_ecom_products.id").Joins("inner join tbl_ecom_customers on tbl_ecom_customers.id = tbl_ecom_carts.customer_id").
		Where("tbl_ecom_carts.is_deleted = 0 and tbl_ecom_products.is_deleted = 0 and tbl_ecom_customers.is_deleted = 0 and tbl_ecom_products.is_active = 1 and tbl_ecom_customers.id = ?", customer_id).Preload("EcommerceCart").Limit(limit).Offset(offset).Order("tbl_ecom_carts.id desc").Find(&cartList).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.EcommerceCartDetails{}, err
	}

	if err := db.Table("tbl_ecom_carts").Joins("inner join tbl_ecom_products on tbl_ecom_products.id = tbl_ecom_carts.product_id").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='discount' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rp on rp.product_id = tbl_ecom_products.id").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='special' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rs on rs.product_id = tbl_ecom_products.id").Joins("inner join tbl_ecom_customers on tbl_ecom_customers.id = tbl_ecom_carts.customer_id").
		Where("tbl_ecom_carts.is_deleted = 0 and tbl_ecom_products.is_deleted = 0 and tbl_ecom_customers.is_deleted = 0 and tbl_ecom_products.is_active = 1 and tbl_ecom_customers.id = ?", customer_id).Count(&count).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.EcommerceCartDetails{}, err
	}

	var final_cartList []model.EcommerceProduct

	var subtotal, totalTax int64

	var totalQuantity int

	for _, cartProduct := range cartList {

		if cartProduct.ProductImagePath != "" {

			imagePaths := strings.Split(cartProduct.ProductImagePath, ",")

			for index, path := range imagePaths {

				modified_path := PathUrl + strings.TrimPrefix(path, "/")

				imagePaths[index] = modified_path
			}

			cartProduct.ProductImageArray = imagePaths

		}

		if cartProduct.ProductVideoPath != "" {

			modified_path := PathUrl + strings.TrimPrefix(cartProduct.ProductVideoPath, "/")

			cartProduct.ProductVideoPath = modified_path
		}

		var priceByQuantity int64

		if cartProduct.SpecialPrice != nil {

			reductionPrice := cartProduct.DefaultPrice - *cartProduct.SpecialPrice

			priceByQuantity = int64(cartProduct.EcommerceCart.Quantity) * int64(reductionPrice)

			subtotal = subtotal + priceByQuantity

		} else if cartProduct.DiscountPrice != nil {

			priceByQuantity = int64(cartProduct.EcommerceCart.Quantity) * int64(*cartProduct.DiscountPrice)

			subtotal = subtotal + priceByQuantity

		} else {

			priceByQuantity = int64(cartProduct.EcommerceCart.Quantity) * int64(cartProduct.DefaultPrice)

			subtotal = subtotal + priceByQuantity
		}

		var taxByQuantity int64 = int64(cartProduct.EcommerceCart.Quantity) * int64(cartProduct.Tax)

		totalTax = totalTax + taxByQuantity

		totalQuantity = totalQuantity + cartProduct.EcommerceCart.Quantity

		final_cartList = append(final_cartList, cartProduct)

	}

	conv_totalCost := strconv.Itoa(int(subtotal) + int(totalTax))

	cartSummary := model.CartSummary{SubTotal: strconv.Itoa(int(subtotal)), TotalTax: strconv.Itoa(int(totalTax)), TotalCost: conv_totalCost, TotalQuantity: totalQuantity}

	return &model.EcommerceCartDetails{CartList: final_cartList, CartSummary: cartSummary, Count: int(count)}, nil
}

func RemoveProductFromCartlist(db *gorm.DB, ctx context.Context, productID int) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return false, err

	}

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	subquery := db.Table("tbl_ecom_customers").Select("id").Where("is_deleted = 0 and member_id = ?", memberid)

	if err := db.Table("tbl_ecom_carts").Where("tbl_ecom_carts.is_deleted = 0 and tbl_ecom_carts.product_id = ? and tbl_ecom_carts.customer_id = (?)", productID, subquery).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": currentTime}).Error; err != nil {

		c.AbortWithError(500, err)

		return false, err
	}

	return true, nil
}
