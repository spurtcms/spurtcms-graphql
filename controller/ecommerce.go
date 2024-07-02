package controller

import (
	"context"
	"errors"
	// "fmt"
	"net/http"
	"spurtcms-graphql/graph/model"
	// "spurtcms-graphql/storage"
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

	var (
		categoryName, releaseDate, searchKeyword string

		categoryId, startingPrice, endingPrice int
	)

	if filter != nil {

		if filter.CategoryName.IsSet() {

			categoryName = *filter.CategoryName.Value()
		}

		if filter.CategoryID.IsSet() {

			categoryId = *filter.CategoryID.Value()
		}

		if filter.ReleaseDate.IsSet() {

			releaseDate = *filter.ReleaseDate.Value()
		}

		if filter.StartingPrice.IsSet() {

			startingPrice = *filter.StartingPrice.Value()
		}

		if filter.EndingPrice.IsSet() {

			endingPrice = *filter.EndingPrice.Value()
		}

		if filter.SearchKeyword.IsSet() {

			searchKeyword = *filter.SearchKeyword.Value()
		}
	}

	if categoryName != "" {

		listQuery = listQuery.Where("tbl_categories.category_name = ?", categoryName)

	} else if categoryId != 0 {

		listQuery = listQuery.Where("tbl_categories.id = ?", categoryId)
	}

	if releaseDate != "" {

		listQuery = listQuery.Where("tbl_ecom_products.created_on >= ?", releaseDate)
	}

	if startingPrice != 0 && endingPrice != 0 {

		listQuery = listQuery.Where("tbl_ecom_products.product_price between (?) and (?)", startingPrice, endingPrice)

	} else if startingPrice != 0 {

		listQuery = listQuery.Where("tbl_ecom_products.product_price >= ?", startingPrice)

	} else if endingPrice != 0 {

		listQuery = listQuery.Where("tbl_ecom_products.product_price <= ?", endingPrice)
	}

	if searchKeyword != "" {

		listQuery = listQuery.Where("LOWER(TRIM(tbl_ecom_products.product_name)) ILIKE LOWER(TRIM(?))", "%"+searchKeyword+"%")
	}

	if sort != nil && sort.Date.IsSet() && *sort.Date.Value() != -1 {

		if *sort.Date.Value() == 1 {

			listQuery = listQuery.Order("tbl_ecom_products.id desc")

		} else if *sort.Date.Value() == 0 {

			listQuery = listQuery.Order("tbl_ecom_products.id ")
		}

	} else if sort != nil && sort.Price.IsSet() && *sort.Price.Value() != -1 {

		if *sort.Price.Value() == 1 {

			listQuery = listQuery.Order("tbl_ecom_products.product_price desc")

		} else if *sort.Price.Value() == 0 {

			listQuery = listQuery.Order("tbl_ecom_products.product_price")

		}

	} else if sort != nil && sort.ViewCount.IsSet() && *sort.ViewCount.Value() != -1 {

		if *sort.ViewCount.Value() == 1 {

			listQuery = listQuery.Order("tbl_ecom_products.view_count desc")

		} else if *sort.ViewCount.Value() == 0 {

			listQuery = listQuery.Order("tbl_ecom_products.view_count")
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

	for _, product := range ecom_products {

		if product.ProductImagePath != "" {

			imagePaths := strings.Split(product.ProductImagePath, ",")

			for index, path := range imagePaths {

				modified_path := PathUrl + strings.TrimPrefix(path, "/")

				imagePaths[index] = modified_path
			}

			product.ProductImageArray = imagePaths

		}

		final_ecomProducts = append(final_ecomProducts, product)
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

	if err := db.Debug().Table("tbl_ecom_carts").Where("is_deleted = 0 and customer_id = ? and product_id = ?", customer_id, productId).Count(&count).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	query := db.Table("tbl_ecom_carts")

	if count > 0 {

		query = query.Where("is_deleted = 0 and customer_id = ? and product_id = ?", customer_id, productId).UpdateColumns(map[string]interface{}{"quantity": gorm.Expr("quantity + ?", cart.Quantity), "modified_on": currentTime})

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

	if err := db.Debug().Table("tbl_ecom_products").Select("tbl_ecom_products.*,rp.price AS discount_price ,rs.price AS special_price,tbl_ecom_carts.*").Joins("inner join tbl_ecom_carts on tbl_ecom_carts.product_id =  tbl_ecom_products.id ").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='discount' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rp on rp.product_id = tbl_ecom_products.id").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='special' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rs on rs.product_id = tbl_ecom_products.id").Joins("inner join tbl_ecom_customers on tbl_ecom_customers.id = tbl_ecom_carts.customer_id").
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

		if cartProduct.EcommerceCart != nil {

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

			var taxByQuantity = int64(cartProduct.EcommerceCart.Quantity) * int64(cartProduct.Tax)

			totalTax = totalTax + taxByQuantity

			totalQuantity = totalQuantity + cartProduct.EcommerceCart.Quantity

		}

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

	if err := db.Debug().Table("tbl_ecom_carts").Where("tbl_ecom_carts.is_deleted = 0 and tbl_ecom_carts.product_id = ? and tbl_ecom_carts.customer_id = (?)", productID, subquery).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": currentTime}).Error; err != nil {

		c.AbortWithError(500, err)

		return false, err
	}

	return true, nil
}

func EcommerceProductOrdersList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.OrderFilter, sort *model.OrderSort) (*model.EcommerceProducts, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return &model.EcommerceProducts{}, err

	}

	var orderedProducts []model.EcommerceProduct

	var count int64

	var customerId int

	if err := db.Table("tbl_ecom_customers").Select("id").Where("is_deleted = 0 and member_id = ?", memberid).Scan(&customerId).Error; err != nil {

		return &model.EcommerceProducts{}, err
	}

	query := db.Debug().Table("tbl_ecom_products as p").Joins("inner join tbl_ecom_product_order_details d on d.product_id = p.id").Joins("inner join tbl_ecom_product_orders o on o.id = d.order_id").Joins("inner join tbl_ecom_order_payments op on op.order_id = o.id").Where("p.is_deleted = 0 and o.is_deleted = 0 and o.customer_id = ?", customerId)

	var (
		status, searchKeyword, orderId, startingDate, endingDate string

		startingPrice, endingPrice, orderHistory, upcomingOrders int
	)

	if filter != nil {

		if filter.Status.IsSet() {

			status = *filter.Status.Value()
		}

		if filter.StartingPrice.IsSet() {

			startingPrice = *filter.StartingPrice.Value()
		}

		if filter.EndingPrice.IsSet() {

			endingPrice = *filter.EndingPrice.Value()
		}

		if filter.SearchKeyword.IsSet() {

			searchKeyword = *filter.SearchKeyword.Value()
		}

		if filter.StartingDate.IsSet() {

			startingDate = *filter.StartingDate.Value()
		}

		if filter.EndingDate.IsSet() {

			endingDate = *filter.EndingDate.Value()
		}

		if filter.OrderID.IsSet() {

			orderId = *filter.OrderID.Value()
		}

		if filter.OrderHistory.IsSet() {

			orderHistory = *filter.OrderHistory.Value()
		}

		if filter.UpcomingOrders.IsSet() {

			upcomingOrders = *filter.UpcomingOrders.Value()
		}

	}

	if upcomingOrders == 1 {

		query = query.Where("o.status in (?)", []string{"placed", "outofdelivery", "shipped"})

	} else if orderHistory == 1 {

		query = query.Where("o.status in (?)", []string{"delivered", "cancelled"})

	} else if status != "" {

		query = query.Where("o.status = ?", status)
	}

	if startingPrice != 0 && endingPrice != 0 {

		query = query.Where("d.price between ? and ?", startingPrice, endingPrice)

	} else if startingPrice != 0 {

		query = query.Where("d.price >= ?", startingPrice)

	} else if endingPrice != 0 {

		query = query.Where("d.price <= ?", endingPrice)

	}

	if searchKeyword != "" {

		query = query.Where("LOWER(TRIM(p.product_name)) ILIKE LOWER(TRIM(?))", "%"+searchKeyword+"%")
	}

	if startingDate != "" && endingDate != "" {

		query = query.Where("o.created_on between ? and ?", startingDate, endingDate)

	} else if startingDate != "" {

		query = query.Where("o.created_on >= ?", startingDate)

	} else if endingDate != "" {

		query = query.Where("o.created_on <= ?", endingDate)
	}

	if orderId != "" {

		query = query.Where("o.uuid = ?", orderId)
	}

	if err := query.Count(&count).Error; err != nil {

		return &model.EcommerceProducts{}, err
	}

	if sort != nil && sort.Date.Value() != nil && *sort.Date.Value() != -1 {

		if *sort.Date.Value() == 1 {

			query = query.Order("o.id desc")

		} else if *sort.Date.Value() == 0 {

			query = query.Order("o.id")

		}

	} else if sort != nil && sort.Price.Value() != nil && *sort.Price.Value() != -1 {

		if *sort.Price.Value() == 1 {

			query = query.Order("d.price desc")

		} else if *sort.Price.Value() == 0 {

			query = query.Order("d.price")

		}

	} else {

		query = query.Order("o.id desc")
	}

	if err := query.Select("p.*,o.id,o.uuid,o.status,o.customer_id,o.created_on,o.shipping_address,d.quantity,d.price,d.tax,op.payment_mode").Limit(limit).Offset(offset).Find(&orderedProducts).Error; err != nil {

		return &model.EcommerceProducts{}, err
	}

	var final_OrderedProductList []model.EcommerceProduct

	for _, product := range orderedProducts {

		if product.ProductImagePath != "" {

			imagePaths := strings.Split(product.ProductImagePath, ",")

			for index, path := range imagePaths {

				modified_path := PathUrl + strings.TrimPrefix(path, "/")

				imagePaths[index] = modified_path
			}

			product.ProductImageArray = imagePaths

		}

		final_OrderedProductList = append(final_OrderedProductList, product)
	}

	return &model.EcommerceProducts{ProductList: final_OrderedProductList, Count: int(count)}, nil
}

func EcommerceProductOrderDetails(db *gorm.DB, ctx context.Context, productID *int, productSlug *string, orderId int) (*model.EcomOrderedProductDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return &model.EcomOrderedProductDetails{}, err

	}

	var customerId int

	if err := db.Table("tbl_ecom_customers").Select("id").Where("is_deleted = 0 and member_id = ?", memberid).Scan(&customerId).Error; err != nil {

		return &model.EcomOrderedProductDetails{}, err
	}

	var orderedProduct model.EcommerceProduct

	query := db.Debug().Table("tbl_ecom_products as p").Joins("inner join tbl_ecom_product_order_details d on d.product_id = p.id").Joins("inner join tbl_ecom_product_orders o on o.id = d.order_id").Joins("inner join tbl_ecom_order_payments op on op.order_id = o.id").Where("p.is_deleted = 0 and o.is_deleted = 0 and o.customer_id = ? and o.id = ?", customerId, orderId)

	if productID != nil {

		query = query.Where("p.id = ?", *productID)

	} else if productSlug != nil {

		query = query.Where("p.product_slug = ?", *productSlug)
	}

	if err := query.Select("p.*,o.id,o.uuid,o.status,o.customer_id,o.created_on,o.shipping_address,d.quantity,d.price,d.tax,op.payment_mode").First(&orderedProduct).Error; err != nil {

		return &model.EcomOrderedProductDetails{}, err
	}

	if orderedProduct.ProductImagePath != "" {

		imagePaths := strings.Split(orderedProduct.ProductImagePath, ",")

		for index, path := range imagePaths {

			modified_path := PathUrl + strings.TrimPrefix(path, "/")

			imagePaths[index] = modified_path
		}

		orderedProduct.ProductImageArray = imagePaths

	}

	var productOrderStatuses []model.OrderStatus

	if err := db.Debug().Table("tbl_ecom_order_statuses").Where("order_id = ?", orderedProduct.OrderID).Find(&productOrderStatuses).Error; err != nil {

		return &model.EcomOrderedProductDetails{}, err
	}

	return &model.EcomOrderedProductDetails{EcommerceProduct: orderedProduct, OrderStatuses: productOrderStatuses}, nil
}

func EcommerceOrderPlacement(db *gorm.DB, ctx context.Context, paymentMode string, shippingAddress string, orderProducts []model.OrderProduct, orderSummary *model.OrderSummary) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return false, err

	}

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	var customerId int

	if err := db.Table("tbl_ecom_customers").Select("id").Where("is_deleted = 0 and member_id = ?", memberid).Scan(&customerId).Error; err != nil {

		return false, err
	}

	unixTime := time.Now().Unix()

	var orderplaced model.EcommerceOrder

	orderId := "SP" + strconv.Itoa(int(unixTime))

	orderplaced.OrderID = orderId

	orderplaced.ShippingAddress = shippingAddress

	orderplaced.CustomerID = customerId

	orderplaced.Status = "placed"

	orderplaced.IsDeleted = 0

	orderplaced.CreatedOn = currentTime

	var totalPrice, totalTax, totalCost int

	if orderSummary != nil {

		totalPrice, _ = strconv.Atoi(orderSummary.SubTotal)

		totalTax, _ = strconv.Atoi(orderSummary.TotalTax)

		totalCost, _ = strconv.Atoi(orderSummary.TotalCost)

	} else {

		for _, product := range orderProducts {

			sum := product.Price * product.Quantity

			totalPrice += sum

			tax := product.Tax * product.Quantity

			totalTax += tax

		}

		totalCost = totalPrice + totalTax
	}

	orderplaced.Price = totalPrice

	orderplaced.Tax = totalTax

	orderplaced.TotalCost = totalCost

	if err := db.Table("tbl_ecom_product_orders").Create(&orderplaced).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	var createorder model.EcommerceOrder

	if err := db.Table("tbl_ecom_product_orders").Where("uuid = ?", orderId).First(&createorder).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	var orderedProductIds []int

	for _, value := range orderProducts {

		var productDetails = model.OrderProductDetails{OrderID: createorder.ID, ProductID: value.ProductID, Quantity: value.Quantity, Price: value.Price, Tax: value.Tax}

		if err := db.Table("tbl_ecom_product_order_details").Create(&productDetails).Error; err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)

			return false, err
		}

		if err := db.Debug().Table("tbl_ecom_products").Where("is_deleted = 0 and is_active = 1 and id = ?", value.ProductID).Update("stock", gorm.Expr("stock - ?", value.Quantity)).Error; err != nil {

			return false, err
		}

		orderedProductIds = append(orderedProductIds, value.ProductID)
	}

	var orderstatus model.OrderStatus

	orderstatus.OrderID = createorder.ID

	orderstatus.OrderStatus = "placed"

	orderstatus.CreatedBy = customerId

	orderstatus.CreatedOn = currentTime

	if err := db.Table("tbl_ecom_order_statuses").Create(&orderstatus).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	var orderPayment model.OrderPayment

	orderPayment.OrderID = createorder.ID

	orderPayment.PaymentMode = paymentMode

	if err := db.Table("tbl_ecom_order_payments").Create(&orderPayment).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	if err := db.Table("tbl_ecom_carts").Where("is_deleted = 0 and product_id in (?) and customer_id = ?", orderedProductIds, customerId).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": currentTime}).Error; err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	return true, nil
}

func EcommerceCustomerDetails(db *gorm.DB, ctx context.Context) (*model.CustomerDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid, _ := c.Get("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return &model.CustomerDetails{}, err
	}

	var customerDetails model.CustomerDetails

	if err := db.Table("tbl_ecom_customers").Where("is_deleted = 0 and member_id = ?", memberid).First(&customerDetails).Error; err != nil {

		return &model.CustomerDetails{}, err
	}

	if customerDetails.ProfileImagePath != nil {

		modified_path := PathUrl + strings.TrimPrefix(*customerDetails.ProfileImagePath, "/")

		customerDetails.ProfileImagePath = &modified_path
	}

	if customerDetails.StreetAddress != nil {

		houseDetails := strings.Split(*customerDetails.StreetAddress, ",")

		customerDetails.HouseNo = &houseDetails[0]

		var area string

		for index, cut := range houseDetails {

			if index == 1 {

				area = area + cut

			} else if index > 1 {

				area = area + "," + cut
			}
		}

		customerDetails.Area = &area
	}

	return &customerDetails, nil
}

func CustomerProfileUpdate(db *gorm.DB, ctx context.Context, customerInput model.CustomerInput) (bool, error) {

	// c, _ := ctx.Value(ContextKey).(*gin.Context)

	// memberid := c.GetInt("memberid")

	// if memberid == 0 {

	// 	err := errors.New("unauthorized access")

	// 	c.AbortWithError(http.StatusUnauthorized, err)

	// 	return false, err
	// }

	// customerDetails := make(map[string]interface{})

	// memberDetails := make(map[string]interface{})

	// if customerInput.ProfileImage.IsSet() && customerInput.ProfileImage.Value() != nil {

	// 	var fileName, filePath string

	// 	storageType, _ := GetStorageType(db)

	// 	fileName = customerInput.ProfileImage.Value().Filename

	// 	file := customerInput.ProfileImage.Value().File

	// 	if storageType.SelectedType == "aws" {

	// 		fmt.Printf("aws-S3 storage selected\n")

	// 		filePath = "member/" + fileName

	// 		err := storage.UploadFileS3(storageType.Aws, customerInput.ProfileImage.Value(), filePath)

	// 		if err != nil {

	// 			fmt.Printf("image upload failed %v\n", err)

	// 			return false, ErrUpload

	// 		}

	// 	} else if storageType.SelectedType == "local" {

	// 		fmt.Printf("local storage selected\n")

	// 		b64Data, err := IoReadSeekerToBase64(file)

	// 		if err != nil {

	// 			return false, err
	// 		}

	// 		endpoint := "gqlSaveLocal"

	// 		url := PathUrl + endpoint

	// 		filePath, err = storage.UploadImageToAdminLocal(b64Data, fileName, url)

	// 		if err != nil {

	// 			return false, ErrUpload
	// 		}

	// 		fmt.Printf("local stored path: %v\n", filePath)

	// 	} else if storageType.SelectedType == "azure" {

	// 		fmt.Printf("azure storage selected")

	// 	} else if storageType.SelectedType == "drive" {

	// 		fmt.Println("drive storage selected")
	// 	}

	// 	customerDetails["profile_image"] = fileName

	// 	memberDetails["profile_image"] = fileName

	// 	customerDetails["profile_image_path"] = filePath

	// 	memberDetails["profile_image_path"] = filePath

	// }

	// customerDetails["first_name"] = customerInput.FirstName

	// memberDetails["first_name"] = customerInput.FirstName

	// customerDetails["email"] = customerInput.Email

	// memberDetails["email"] = customerInput.Email

	// if customerInput.LastName.IsSet() && customerInput.LastName.Value() != nil {

	// 	customerDetails["last_name"] = *customerInput.LastName.Value()

	// 	memberDetails["last_name"] = *customerInput.LastName.Value()
	// }

	// if customerInput.MobileNo.IsSet() && customerInput.MobileNo.Value() != nil {

	// 	customerDetails["mobile_no"] = *customerInput.MobileNo.Value()

	// 	memberDetails["mobile_no"] = *customerInput.MobileNo.Value()
	// }

	// if customerInput.Username.IsSet() && customerInput.Username.Value() != nil {

	// 	customerDetails["username"] = *customerInput.Username.Value()

	// 	memberDetails["username"] = *customerInput.Username.Value()
	// }

	// if customerInput.IsActive.IsSet() && customerInput.IsActive.Value() != nil {

	// 	customerDetails["is_active"] = *customerInput.IsActive.Value()

	// 	memberDetails["is_active"] = *customerInput.IsActive.Value()
	// }

	// if customerInput.StreetAddress.IsSet() && customerInput.StreetAddress.Value() != nil {

	// 	customerDetails["street_address"] = *customerInput.StreetAddress.Value()
	// }

	// if customerInput.City.IsSet() && customerInput.City.Value() != nil {

	// 	customerDetails["city"] = *customerInput.City.Value()
	// }

	// if customerInput.Country.IsSet() && customerInput.Country.Value() != nil {

	// 	customerDetails["country"] = *customerInput.Country.Value()
	// }

	// if customerInput.State.IsSet() && customerInput.State.Value() != nil {

	// 	customerDetails["state"] = *customerInput.State.Value()
	// }

	// if customerInput.ZipCode.IsSet() && customerInput.ZipCode.Value() != nil {

	// 	customerDetails["zip_code"] = *customerInput.ZipCode.Value()
	// }

	// if customerInput.Password.IsSet() && customerInput.Password.Value() != nil && *customerInput.Password.Value() != "" {

	// 	hashpass, err := HashingPassword(*customerInput.Password.Value())

	// 	if err != nil {

	// 		return false, ErrPassHash
	// 	}

	// 	customerDetails["password"] = hashpass

	// 	memberDetails["password"] = hashpass
	// }

	// currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	// customerDetails["modified_on"] = currentTime

	// memberDetails["modified_on"] = currentTime

	// customerDetails["modified_by"] = memberid

	// memberDetails["modified_by"] = memberid

	// if err := db.Debug().Table("tbl_ecom_customers").Where("is_deleted = 0 and member_id = ?", memberid).UpdateColumns(&customerDetails).Error; err != nil {

	// 	return false, err
	// }

	// if err := db.Debug().Table("tbl_members").Where("is_deleted = 0 and id = ?", memberid).UpdateColumns(&memberDetails).Error; err != nil {

	// 	return false, err
	// }

	return true, nil
}

func UpdateProductViewCount(db *gorm.DB, ctx context.Context, productID *int, productSlug *string) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	if productID == nil && productSlug == nil {

		return false, ErrMandatory
	}

	query := db.Debug().Table("tbl_ecom_products").Where("is_deleted = 0 and is_active = 1")

	if productID != nil {

		query = query.Where("id = ?", *productID)

	} else if productSlug != nil {

		query = query.Where("product_slug = ?", *productSlug)
	}

	err := query.Update("view_count", gorm.Expr("view_count + 1")).Error

	if err != nil {

		c.AbortWithError(500, err)

		return false, err
	}

	return true, nil
}

func EcommerceOrderStatusNames(db *gorm.DB,ctx context.Context) ([]model.OrderStatusNames, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return []model.OrderStatusNames{}, err
	}

	var orderStatus []model.OrderStatusNames

	if err := db.Debug().Table("tbl_ecom_statuses").Find(&orderStatus).Error;err!= nil{

		return []model.OrderStatusNames{}, err
	}

	return orderStatus, nil
}
