package controller

import (
	"context"
	"errors"
	"fmt"

	// "fmt"
	"net/http"
	"spurtcms-graphql/graph/model"
	"spurtcms-graphql/storage"

	// "spurtcms-graphql/storage"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spurtcms/ecommerce"
	"gorm.io/gorm"
)

func EcommerceProductList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.ProductFilter, sort *model.ProductSort) (*model.EcommerceProducts, error) {

	// c, _ := ctx.Value(ContextKey).(*gin.Context)

	var (
		final_ecomProducts []model.EcommerceProduct
		product            model.EcommerceProduct
		count              int64
		ecomFilter         ecommerce.ProductFilter
		ecomSort           ecommerce.ProductSort
	)

	if filter != nil {

		if filter.CategoryName.IsSet() {

			ecomFilter.CategoryName = *filter.CategoryName.Value()
		}

		if filter.CategoryID.IsSet() {

			ecomFilter.CategoryId = *filter.CategoryID.Value()
		}

		if filter.ReleaseDate.IsSet() {

			ecomFilter.ReleaseDate = *filter.ReleaseDate.Value()
		}

		if filter.StartingPrice.IsSet() {

			ecomFilter.StartingPrice = *filter.StartingPrice.Value()
		}

		if filter.EndingPrice.IsSet() {

			ecomFilter.EndingPrice = *filter.EndingPrice.Value()
		}

		if filter.SearchKeyword.IsSet() {

			ecomFilter.SearchKeyword = *filter.SearchKeyword.Value()
		}
	}

	if sort != nil {

		if sort.Date.IsSet() {

			ecomSort.Date = *sort.Date.Value()
		}

		if sort.Price.IsSet() {

			ecomSort.Price = *sort.Price.Value()
		}

		if sort.ViewCount.IsSet() {

			ecomSort.ViewCount = *sort.ViewCount.Value()
		}
	}

	dbType := db.Config.Dialector.Name()

	productsList, count, err := EcomInstance.GetProductListAndCount(limit, offset, ecomFilter, ecomSort, dbType)
	if err != nil {

		return &model.EcommerceProducts{}, err
	}

	for _, productList := range productsList {

		product.CategoriesID = productList.CategoriesID
		product.CreatedBy = productList.CreatedBy
		product.CreatedOn = productList.CreatedOn
		product.DefaultPrice = productList.DefaultPrice
		product.DeletedBy = &productList.DeletedBy
		product.DeletedOn = &productList.DeletedOn
		product.DiscountPrice = &productList.DiscountPrice
		product.EcommerceCart.CreatedOn = productList.CreatedOn
		product.EcommerceCart.CustomerID = productList.TblEcomCart.CustomerID
		product.EcommerceCart.DeletedOn = &productList.DeletedOn
		product.EcommerceCart.ID = productList.TblEcomCart.ID
		product.EcommerceCart.IsDeleted = productList.TblEcomCart.IsDeleted
		product.EcommerceCart.ModifiedOn = &productList.TblEcomCart.ModifiedOn
		product.EcommerceCart.ProductID = productList.TblEcomCart.ProductID
		product.EcommerceCart.Quantity = productList.TblEcomCart.Quantity
		product.ID = productList.ID
		product.IsActive = productList.IsActive
		product.IsDeleted = productList.IsDeleted
		product.ModifiedBy = &productList.ModifiedBy
		product.ModifiedOn = &productList.ModifiedOn
		product.OrderCustomer = &productList.OrderCustomer
		product.OrderID = &productList.OrderID
		product.OrderPrice = &productList.OrderPrice
		product.OrderQuantity = &productList.OrderQuantity
		product.OrderStatus = &productList.OrderStatus
		product.OrderTax = &productList.OrderTax
		product.OrderTime = &productList.OrderTime
		product.OrderUniqueID = &productList.OrderUniqueID
		product.PaymentMode = &productList.PaymentMode
		product.ProductDescription = productList.ProductDescription
		product.ProductImageArray = productList.ProductImageArray
		product.ProductImagePath = productList.ProductImagePath
		product.ProductName = productList.ProductName
		product.ProductSlug = productList.ProductSlug
		product.ProductVimeoPath = &productList.ProductVimeoPath
		product.ProductYoutubePath = &productList.ProductYoutubePath
		product.ShippingDetails = &productList.ShippingDetails
		product.Sku = productList.Sku
		product.SpecialPrice = &productList.SpecialPrice
		product.Tax = productList.Tax
		product.Totalcost = productList.Totalcost
		product.ViewCount = &productList.ViewCount

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

	// c, _ := ctx.Value(ContextKey).(*gin.Context)

	var (
		productdtl model.EcommerceProduct
		id         int
		slug       string
	)

	// currentTime := time.Now().In(TimeZone).Format("2006-01-02 15:04:05"

	if productId != nil {

		id = *productId
	}
	if productSlug != nil {

		slug = *productSlug
	}

	product, err := EcomInstance.GetProductdetailsById(id, slug)
	if err != nil {

		return &model.EcommerceProduct{}, err
	}

	productdtl.CategoriesID = product.CategoriesID
	productdtl.CreatedBy = product.CreatedBy
	productdtl.CreatedOn = product.CreatedOn
	productdtl.DefaultPrice = product.DefaultPrice
	productdtl.DeletedBy = &product.DeletedBy
	productdtl.DeletedOn = &product.DeletedOn
	productdtl.DiscountPrice = &product.DiscountPrice
	productdtl.EcommerceCart.CreatedOn = product.CreatedOn
	productdtl.EcommerceCart.CustomerID = product.TblEcomCart.CustomerID
	productdtl.EcommerceCart.DeletedOn = &product.DeletedOn
	productdtl.EcommerceCart.ID = product.TblEcomCart.ID
	productdtl.EcommerceCart.IsDeleted = product.TblEcomCart.IsDeleted
	productdtl.EcommerceCart.ModifiedOn = &product.TblEcomCart.ModifiedOn
	productdtl.EcommerceCart.ProductID = product.TblEcomCart.ProductID
	productdtl.EcommerceCart.Quantity = product.TblEcomCart.Quantity
	productdtl.ID = product.ID
	productdtl.IsActive = product.IsActive
	productdtl.IsDeleted = product.IsDeleted
	productdtl.ModifiedBy = &product.ModifiedBy
	productdtl.ModifiedOn = &product.ModifiedOn
	productdtl.OrderCustomer = &product.OrderCustomer
	productdtl.OrderID = &product.OrderID
	productdtl.OrderPrice = &product.OrderPrice
	productdtl.OrderQuantity = &product.OrderQuantity
	productdtl.OrderStatus = &product.OrderStatus
	productdtl.OrderTax = &product.OrderTax
	productdtl.OrderTime = &product.OrderTime
	productdtl.OrderUniqueID = &product.OrderUniqueID
	productdtl.PaymentMode = &product.PaymentMode
	productdtl.ProductDescription = product.ProductDescription
	productdtl.ProductImageArray = product.ProductImageArray
	productdtl.ProductImagePath = product.ProductImagePath
	productdtl.ProductName = product.ProductName
	productdtl.ProductSlug = product.ProductSlug
	productdtl.ProductVimeoPath = &product.ProductVimeoPath
	productdtl.ProductYoutubePath = &product.ProductYoutubePath
	productdtl.ShippingDetails = &product.ShippingDetails
	productdtl.Sku = product.Sku
	productdtl.SpecialPrice = &product.SpecialPrice
	productdtl.Tax = product.Tax
	productdtl.Totalcost = product.Totalcost
	productdtl.ViewCount = &product.ViewCount

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

// func EcommerceAddToCart(db *gorm.DB, ctx context.Context, productID *int, productSlug *string, quantity int) (bool, error) {

// 	c, _ := ctx.Value(ContextKey).(*gin.Context)

// 	memberid := c.GetInt("memberid")

// 	if memberid == 0 {

// 		err := errors.New("unauthorized access")

// 		c.AbortWithError(http.StatusUnauthorized, err)

// 		return false, err

// 	}

// 	var cart model.EcommerceCart

// 	var productId int

// 	if productID != nil {

// 		productId = *productID

// 	} else if productSlug != nil {

// 		if err := db.Table("tbl_ecom_products").Select("id").Where("is_deleted = 0 and product_slug = ?", *productSlug).Scan(&productId).Error; err != nil {

// 			c.AbortWithError(500, err)

// 			return false, err
// 		}
// 	}

// 	var customer_id int

// 	if err := db.Table("tbl_ecom_customers").Select("tbl_ecom_customers.id").Where("tbl_ecom_customers.is_deleted = 0 and tbl_ecom_customers.member_id = ?", memberid).Scan(&customer_id).Error; err != nil {

// 		c.AbortWithError(500, err)

// 		return false, err
// 	}

// 	if customer_id == 0 {

// 		err := errors.New("customer id not found")

// 		c.AbortWithError(500, err)

// 		return false, err
// 	}

// 	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

// 	cart.ProductID = productId

// 	cart.CustomerID = customer_id

// 	cart.Quantity = quantity

// 	cart.IsDeleted = 0

// 	cart.CreatedOn = currentTime

// 	var count int64

// 	if err := db.Debug().Table("tbl_ecom_carts").Where("is_deleted = 0 and customer_id = ? and product_id = ?", customer_id, productId).Count(&count).Error; err != nil {

// 		c.AbortWithError(http.StatusInternalServerError, err)

// 		return false, err
// 	}

// 	query := db.Table("tbl_ecom_carts")

// 	if count > 0 {

// 		query = query.Where("is_deleted = 0 and customer_id = ? and product_id = ?", customer_id, productId).UpdateColumns(map[string]interface{}{"quantity": gorm.Expr("quantity + ?", cart.Quantity), "modified_on": currentTime})

// 	} else {

// 		query = query.Create(&cart)
// 	}

// 	if err := query.Error; err != nil {

// 		c.AbortWithError(http.StatusInternalServerError, err)

// 		return false, err
// 	}

// 	return true, nil
// }

func EcommerceAddToCart(db *gorm.DB, ctx context.Context, productID *int, productSlug *string, quantity int) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return false, err

	}

	var cart ecommerce.TblEcomCart

	var product ecommerce.EcommerceProduct

	var productId int

	var err error

	if productID != nil {

		productId = *productID

	} else if productSlug != nil {

		product, err = EcomInstance.GetProduct(*productID, *productSlug)
		if err != nil {

			c.AbortWithError(500, err)

			return false, err
		}

		productId = product.ID
	}

	var customer ecommerce.TblEcomCustomers

	customer, err = EcomInstance.GetCustomer(memberid)
	if err != nil {

		c.AbortWithError(500, err)

		return false, err
	}

	if customer.Id == 0 {

		err := errors.New("customer id not found")

		c.AbortWithError(500, err)

		return false, err
	}

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	cart.ProductID = productId

	cart.CustomerID = customer.Id

	cart.Quantity = quantity

	cart.IsDeleted = 0

	cart.CreatedOn = currentTime

	err = EcomAuthInstance.AddToCart(cart)
	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	return true, nil
}

// func EcommerceCartList(db *gorm.DB, ctx context.Context, limit, offset int) (*model.EcommerceCartDetails, error) {

// 	c, _ := ctx.Value(ContextKey).(*gin.Context)

// 	memberid := c.GetInt("memberid")

// 	if memberid == 0 {

// 		err := errors.New("unauthorized access")

// 		c.AbortWithError(http.StatusUnauthorized, err)

// 		return &model.EcommerceCartDetails{}, err

// 	}

// 	var customer_id int

// 	if err := db.Table("tbl_ecom_customers").Select("tbl_ecom_customers.id").Where("tbl_ecom_customers.is_deleted = 0 and tbl_ecom_customers.member_id = ?", memberid).Scan(&customer_id).Error; err != nil {

// 		c.AbortWithError(500, err)

// 		return &model.EcommerceCartDetails{}, err
// 	}

// 	if customer_id == 0 {

// 		err := errors.New("customer id not found")

// 		c.AbortWithError(500, err)

// 		return &model.EcommerceCartDetails{}, err
// 	}

// 	var cartList []model.EcommerceProduct

// 	var count int64

// 	if err := db.Debug().Table("tbl_ecom_products").Select("tbl_ecom_products.*,rp.price AS discount_price ,rs.price AS special_price,tbl_ecom_carts.*").Joins("inner join tbl_ecom_carts on tbl_ecom_carts.product_id =  tbl_ecom_products.id ").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='discount' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rp on rp.product_id = tbl_ecom_products.id").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='special' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rs on rs.product_id = tbl_ecom_products.id").Joins("inner join tbl_ecom_customers on tbl_ecom_customers.id = tbl_ecom_carts.customer_id").
// 		Where("tbl_ecom_carts.is_deleted = 0 and tbl_ecom_products.is_deleted = 0 and tbl_ecom_customers.is_deleted = 0 and tbl_ecom_products.is_active = 1 and tbl_ecom_customers.id = ?", customer_id).Preload("EcommerceCart").Limit(limit).Offset(offset).Order("tbl_ecom_carts.id desc").Find(&cartList).Error; err != nil {

// 		c.AbortWithError(http.StatusInternalServerError, err)

// 		return &model.EcommerceCartDetails{}, err
// 	}

// 	if err := db.Table("tbl_ecom_carts").Joins("inner join tbl_ecom_products on tbl_ecom_products.id = tbl_ecom_carts.product_id").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='discount' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rp on rp.product_id = tbl_ecom_products.id").Joins("left join (select *, ROW_NUMBER() OVER (PARTITION BY tbl_ecom_product_pricings.id, tbl_ecom_product_pricings.type ORDER BY tbl_ecom_product_pricings.priority,tbl_ecom_product_pricings.start_date desc) AS rn from tbl_ecom_product_pricings where tbl_ecom_product_pricings.type ='special' and tbl_ecom_product_pricings.start_date <= now() and tbl_ecom_product_pricings.end_date >= now()) rs on rs.product_id = tbl_ecom_products.id").Joins("inner join tbl_ecom_customers on tbl_ecom_customers.id = tbl_ecom_carts.customer_id").
// 		Where("tbl_ecom_carts.is_deleted = 0 and tbl_ecom_products.is_deleted = 0 and tbl_ecom_customers.is_deleted = 0 and tbl_ecom_products.is_active = 1 and tbl_ecom_customers.id = ?", customer_id).Count(&count).Error; err != nil {

// 		c.AbortWithError(http.StatusInternalServerError, err)

// 		return &model.EcommerceCartDetails{}, err
// 	}

// 	var final_cartList []model.EcommerceProduct

// 	var subtotal, totalTax int64

// 	var totalQuantity int

// 	for _, cartProduct := range cartList {

// 		if cartProduct.ProductImagePath != "" {

// 			imagePaths := strings.Split(cartProduct.ProductImagePath, ",")

// 			for index, path := range imagePaths {

// 				modified_path := PathUrl + strings.TrimPrefix(path, "/")

// 				imagePaths[index] = modified_path
// 			}

// 			cartProduct.ProductImageArray = imagePaths
// 		}

// 		if cartProduct.EcommerceCart != nil {

// 			var priceByQuantity int64

// 			if cartProduct.SpecialPrice != nil {

// 				reductionPrice := cartProduct.DefaultPrice - *cartProduct.SpecialPrice

// 				priceByQuantity = int64(cartProduct.EcommerceCart.Quantity) * int64(reductionPrice)

// 				subtotal = subtotal + priceByQuantity

// 			} else if cartProduct.DiscountPrice != nil {

// 				priceByQuantity = int64(cartProduct.EcommerceCart.Quantity) * int64(*cartProduct.DiscountPrice)

// 				subtotal = subtotal + priceByQuantity

// 			} else {

// 				priceByQuantity = int64(cartProduct.EcommerceCart.Quantity) * int64(cartProduct.DefaultPrice)

// 				subtotal = subtotal + priceByQuantity
// 			}

// 			var taxByQuantity = int64(cartProduct.EcommerceCart.Quantity) * int64(cartProduct.Tax)

// 			totalTax = totalTax + taxByQuantity

// 			totalQuantity = totalQuantity + cartProduct.EcommerceCart.Quantity

// 		}

// 		final_cartList = append(final_cartList, cartProduct)

// 	}

// 	conv_totalCost := strconv.Itoa(int(subtotal) + int(totalTax))

// 	cartSummary := model.CartSummary{SubTotal: strconv.Itoa(int(subtotal)), TotalTax: strconv.Itoa(int(totalTax)), TotalCost: conv_totalCost, TotalQuantity: totalQuantity}

// 	return &model.EcommerceCartDetails{CartList: final_cartList, CartSummary: cartSummary, Count: int(count)}, nil
// }

func EcommerceCartList(db *gorm.DB, ctx context.Context, limit, offset int) (*model.EcommerceCartDetails, error) {

	var (
		cartList           []ecommerce.EcommerceProduct
		count              int64
		final_cartList     []model.EcommerceProduct
		cartProductLocal   model.EcommerceProduct
		subtotal, totalTax int64
		totalQuantity      int
		customer           ecommerce.TblEcomCustomers
	)

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	var err error

	if memberid == 0 {

		err = errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return &model.EcommerceCartDetails{}, err

	}

	customer, err = EcomInstance.GetCustomer(memberid)
	if err != nil {

		c.AbortWithError(500, err)

		return &model.EcommerceCartDetails{}, err
	}

	if customer.Id == 0 {

		err := errors.New("customer id not found")

		c.AbortWithError(500, err)

		return &model.EcommerceCartDetails{}, err
	}

	cartList, err = EcomAuthInstance.GetCartListById(customer.Id, limit, offset)
	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.EcommerceCartDetails{}, err

	}

	count, err = EcomAuthInstance.GetCartListCountById(customer.Id)
	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return &model.EcommerceCartDetails{}, err

	}

	for _, cartProduct := range cartList {

		cartProductLocal.CategoriesID = cartProduct.CategoriesID
		cartProductLocal.CreatedBy = cartProduct.CreatedBy
		cartProductLocal.CreatedOn = cartProduct.CreatedOn
		cartProductLocal.DefaultPrice = cartProduct.DefaultPrice
		cartProductLocal.DeletedBy = &cartProduct.DeletedBy
		cartProductLocal.DeletedOn = &cartProduct.DeletedOn
		cartProductLocal.DiscountPrice = &cartProduct.DiscountPrice
		cartProductLocal.EcommerceCart.CreatedOn = cartProduct.TblEcomCart.CreatedOn
		cartProductLocal.EcommerceCart.CustomerID = cartProduct.TblEcomCart.CustomerID
		cartProductLocal.EcommerceCart.DeletedOn = &cartProduct.TblEcomCart.DeletedOn
		cartProductLocal.EcommerceCart.ID = cartProduct.TblEcomCart.ID
		cartProductLocal.EcommerceCart.IsDeleted = cartProduct.TblEcomCart.IsDeleted
		cartProductLocal.EcommerceCart.ModifiedOn = &cartProduct.TblEcomCart.ModifiedOn
		cartProductLocal.EcommerceCart.ProductID = cartProduct.TblEcomCart.ProductID
		cartProductLocal.EcommerceCart.Quantity = cartProduct.TblEcomCart.Quantity
		cartProductLocal.ID = cartProduct.ID
		cartProductLocal.IsActive = cartProduct.IsActive
		cartProductLocal.IsDeleted = cartProduct.IsActive
		cartProductLocal.ModifiedBy = &cartProduct.ModifiedBy
		cartProductLocal.ModifiedOn = &cartProduct.ModifiedOn
		cartProductLocal.OrderCustomer = &cartProduct.OrderCustomer
		cartProductLocal.OrderID = &cartProduct.OrderID
		cartProductLocal.OrderPrice = &cartProduct.OrderPrice
		cartProductLocal.OrderQuantity = &cartProduct.OrderQuantity
		cartProductLocal.OrderStatus = &cartProduct.OrderStatus
		cartProductLocal.OrderTax = &cartProduct.OrderTax
		cartProductLocal.OrderTime = &cartProduct.OrderTime
		cartProductLocal.OrderUniqueID = &cartProduct.OrderUniqueID
		cartProductLocal.PaymentMode = &cartProduct.PaymentMode
		cartProductLocal.ProductDescription = cartProduct.ProductDescription
		cartProductLocal.ProductImageArray = cartProduct.ProductImageArray
		cartProductLocal.ProductImagePath = cartProduct.ProductImagePath
		cartProductLocal.ProductName = cartProduct.ProductName
		cartProductLocal.ProductSlug = cartProduct.ProductSlug
		cartProductLocal.ProductVimeoPath = &cartProduct.ProductVimeoPath
		cartProductLocal.ProductYoutubePath = &cartProduct.ProductYoutubePath
		cartProductLocal.ShippingDetails = &cartProduct.ShippingDetails
		cartProductLocal.Sku = cartProduct.Sku
		cartProductLocal.SpecialPrice = &cartProduct.SpecialPrice
		cartProductLocal.Tax = cartProduct.Tax
		cartProductLocal.Totalcost = cartProduct.Totalcost
		cartProductLocal.EcommerceCart.CustomerID = cartProduct.TblEcomCart.CustomerID
		cartProductLocal.EcommerceCart.ID = cartProduct.TblEcomCart.ID
		cartProductLocal.EcommerceCart.ProductID = cartProduct.TblEcomCart.ProductID

		if cartProductLocal.ProductImagePath != "" {

			imagePaths := strings.Split(cartProductLocal.ProductImagePath, ",")

			for index, path := range imagePaths {

				modified_path := PathUrl + strings.TrimPrefix(path, "/")

				imagePaths[index] = modified_path
			}

			cartProductLocal.ProductImageArray = imagePaths
		}

		if cartProductLocal.EcommerceCart != nil {

			var priceByQuantity int64

			if cartProductLocal.SpecialPrice != nil {

				reductionPrice := cartProductLocal.DefaultPrice - *cartProductLocal.SpecialPrice

				priceByQuantity = int64(cartProductLocal.EcommerceCart.Quantity) * int64(reductionPrice)

				subtotal = subtotal + priceByQuantity

			} else if cartProductLocal.DiscountPrice != nil {

				priceByQuantity = int64(cartProductLocal.EcommerceCart.Quantity) * int64(*cartProductLocal.DiscountPrice)

				subtotal = subtotal + priceByQuantity

			} else {

				priceByQuantity = int64(cartProductLocal.EcommerceCart.Quantity) * int64(cartProductLocal.DefaultPrice)

				subtotal = subtotal + priceByQuantity
			}

			var taxByQuantity = int64(cartProductLocal.EcommerceCart.Quantity) * int64(cartProductLocal.Tax)

			totalTax = totalTax + taxByQuantity

			totalQuantity = totalQuantity + cartProductLocal.EcommerceCart.Quantity

		}

		final_cartList = append(final_cartList, cartProductLocal)

	}

	conv_totalCost := strconv.Itoa(int(subtotal) + int(totalTax))

	cartSummary := model.CartSummary{SubTotal: strconv.Itoa(int(subtotal)), TotalTax: strconv.Itoa(int(totalTax)), TotalCost: conv_totalCost, TotalQuantity: totalQuantity}

	return &model.EcommerceCartDetails{CartList: final_cartList, CartSummary: cartSummary, Count: int(count)}, nil
}

// func RemoveProductFromCartlist(db *gorm.DB, ctx context.Context, productID int) (bool, error) {

// 	c, _ := ctx.Value(ContextKey).(*gin.Context)

// 	memberid := c.GetInt("memberid")

// 	if memberid == 0 {

// 		err := errors.New("unauthorized access")

// 		c.AbortWithError(http.StatusUnauthorized, err)

// 		return false, err

// 	}

// 	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

// 	subquery := db.Table("tbl_ecom_customers").Select("id").Where("is_deleted = 0 and member_id = ?", memberid)

// 	if err := db.Debug().Table("tbl_ecom_carts").Where("tbl_ecom_carts.is_deleted = 0 and tbl_ecom_carts.product_id = ? and tbl_ecom_carts.customer_id = (?)", productID, subquery).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": currentTime}).Error; err != nil {

// 		c.AbortWithError(500, err)

// 		return false, err
// 	}

// 	return true, nil
// }

func RemoveProductFromCartlist(db *gorm.DB, ctx context.Context, productId int) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberId := c.GetInt("memberid")

	if memberId == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return false, err
	}

	err := EcomInstance.RemoveProductFromCartlist(productId, memberId)
	if err != nil {
		c.AbortWithError(500, err)

		return false, err
	}

	return true, nil
}

// func EcommerceProductOrdersList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.OrderFilter, sort *model.OrderSort) (*model.EcommerceProducts, error) {

// 	c, _ := ctx.Value(ContextKey).(*gin.Context)

// 	memberid := c.GetInt("memberid")

// 	if memberid == 0 {

// 		err := errors.New("unauthorized access")

// 		c.AbortWithError(http.StatusUnauthorized, err)

// 		return &model.EcommerceProducts{}, err

// 	}

// 	var orderedProducts []model.EcommerceProduct

// 	var count int64

// 	var customerId int

// 	if err := db.Table("tbl_ecom_customers").Select("id").Where("is_deleted = 0 and member_id = ?", memberid).Scan(&customerId).Error; err != nil {

// 		return &model.EcommerceProducts{}, err
// 	}

// 	query := db.Debug().Table("tbl_ecom_products as p").Joins("inner join tbl_ecom_product_order_details d on d.product_id = p.id").Joins("inner join tbl_ecom_product_orders o on o.id = d.order_id").Joins("inner join tbl_ecom_order_payments op on op.order_id = o.id").Where("p.is_deleted = 0 and o.is_deleted = 0 and o.customer_id = ?", customerId)

// 	var (
// 		status, searchKeyword, orderId, startingDate, endingDate string

// 		startingPrice, endingPrice, orderHistory, upcomingOrders int
// 	)

// 	if filter != nil {

// 		if filter.Status.IsSet() {

// 			status = *filter.Status.Value()
// 		}

// 		if filter.StartingPrice.IsSet() {

// 			startingPrice = *filter.StartingPrice.Value()
// 		}

// 		if filter.EndingPrice.IsSet() {

// 			endingPrice = *filter.EndingPrice.Value()
// 		}

// 		if filter.SearchKeyword.IsSet() {

// 			searchKeyword = *filter.SearchKeyword.Value()
// 		}

// 		if filter.StartingDate.IsSet() {

// 			startingDate = *filter.StartingDate.Value()
// 		}

// 		if filter.EndingDate.IsSet() {

// 			endingDate = *filter.EndingDate.Value()
// 		}

// 		if filter.OrderID.IsSet() {

// 			orderId = *filter.OrderID.Value()
// 		}

// 		if filter.OrderHistory.IsSet() {

// 			orderHistory = *filter.OrderHistory.Value()
// 		}

// 		if filter.UpcomingOrders.IsSet() {

// 			upcomingOrders = *filter.UpcomingOrders.Value()
// 		}

// 	}

// 	if upcomingOrders == 1 {

// 		query = query.Where("o.status in (?)", []string{"placed", "outofdelivery", "shipped"})

// 	} else if orderHistory == 1 {

// 		query = query.Where("o.status in (?)", []string{"delivered", "cancelled"})

// 	} else if status != "" {

// 		query = query.Where("o.status = ?", status)
// 	}

// 	if startingPrice != 0 && endingPrice != 0 {

// 		query = query.Where("d.price between ? and ?", startingPrice, endingPrice)

// 	} else if startingPrice != 0 {

// 		query = query.Where("d.price >= ?", startingPrice)

// 	} else if endingPrice != 0 {

// 		query = query.Where("d.price <= ?", endingPrice)

// 	}

// 	if searchKeyword != "" {

// 		query = query.Where("LOWER(TRIM(p.product_name)) ILIKE LOWER(TRIM(?))", "%"+searchKeyword+"%")
// 	}

// 	if startingDate != "" && endingDate != "" {

// 		query = query.Where("o.created_on between ? and ?", startingDate, endingDate)

// 	} else if startingDate != "" {

// 		query = query.Where("o.created_on >= ?", startingDate)

// 	} else if endingDate != "" {

// 		query = query.Where("o.created_on <= ?", endingDate)
// 	}

// 	if orderId != "" {

// 		query = query.Where("o.uuid = ?", orderId)
// 	}

// 	if err := query.Count(&count).Error; err != nil {

// 		return &model.EcommerceProducts{}, err
// 	}

// 	if sort != nil && sort.Date.Value() != nil && *sort.Date.Value() != -1 {

// 		if *sort.Date.Value() == 1 {

// 			query = query.Order("o.id desc")

// 		} else if *sort.Date.Value() == 0 {

// 			query = query.Order("o.id")

// 		}

// 	} else if sort != nil && sort.Price.Value() != nil && *sort.Price.Value() != -1 {

// 		if *sort.Price.Value() == 1 {

// 			query = query.Order("d.price desc")

// 		} else if *sort.Price.Value() == 0 {

// 			query = query.Order("d.price")

// 		}

// 	} else {

// 		query = query.Order("o.id desc")
// 	}

// 	if err := query.Select("p.*,o.id,o.uuid,o.status,o.customer_id,o.created_on,o.shipping_address,d.quantity,d.price,d.tax,op.payment_mode").Limit(limit).Offset(offset).Find(&orderedProducts).Error; err != nil {

// 		return &model.EcommerceProducts{}, err
// 	}

// 	var final_OrderedProductList []model.EcommerceProduct

// 	for _, product := range orderedProducts {

// 		if product.ProductImagePath != "" {

// 			imagePaths := strings.Split(product.ProductImagePath, ",")

// 			for index, path := range imagePaths {

// 				modified_path := PathUrl + strings.TrimPrefix(path, "/")

// 				imagePaths[index] = modified_path
// 			}

// 			product.ProductImageArray = imagePaths

// 		}

// 		final_OrderedProductList = append(final_OrderedProductList, product)
// 	}

// 	return &model.EcommerceProducts{ProductList: final_OrderedProductList, Count: int(count)}, nil
// }

func EcommerceProductOrdersList(db *gorm.DB, ctx context.Context, limit int, offset int, filter *model.OrderFilter, sort *model.OrderSort) (*model.EcommerceProducts, error) {

	var (
		count                    int64
		err                      error
		productFilter            ecommerce.ProductFilter
		productSort              ecommerce.ProductSort
		final_OrderedProductList []model.EcommerceProduct
		productLocal             model.EcommerceProduct
		customer                 ecommerce.TblEcomCustomers
	)

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberId := c.GetInt("memberid")

	if memberId == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return &model.EcommerceProducts{}, err

	}

	customer, err = EcomInstance.GetCustomer(memberId)
	if err != nil {
		return &model.EcommerceProducts{}, err
	}

	if filter != nil {

		if filter.Status.IsSet() {

			productFilter.Status = *filter.Status.Value()
		}

		if filter.StartingPrice.IsSet() {

			productFilter.StartingPrice = *filter.StartingPrice.Value()
		}

		if filter.EndingPrice.IsSet() {

			productFilter.EndingPrice = *filter.EndingPrice.Value()
		}

		if filter.SearchKeyword.IsSet() {

			productFilter.SearchKeyword = *filter.SearchKeyword.Value()
		}

		if filter.StartingDate.IsSet() {

			productFilter.StartingDate = *filter.StartingDate.Value()
		}

		if filter.EndingDate.IsSet() {

			productFilter.EndingDate = *filter.EndingDate.Value()
		}

		if filter.OrderID.IsSet() {

			productFilter.OrderId = *filter.OrderID.Value()
		}

		if filter.OrderHistory.IsSet() {

			productFilter.OrderHistory = *filter.OrderHistory.Value()
		}

		if filter.UpcomingOrders.IsSet() {

			productFilter.UpcomingOrders = *filter.UpcomingOrders.Value()
		}

	}

	if sort != nil {

		if sort.Price.IsSet() {

			productSort.Date = *sort.Date.Value()

		}
		if sort.Price.IsSet() {

			productSort.Price = *sort.Price.Value()

		}
	}

	productOrderedList, count, err := EcomAuthInstance.GetProductOrdersList(productFilter, productSort, customer.Id, limit, offset)
	if err != nil {

		return &model.EcommerceProducts{}, err
	}

	// need to be changed

	for _, product := range productOrderedList {

		productLocal.CategoriesID = product.CategoriesID
		productLocal.CreatedBy = product.CreatedBy
		productLocal.CreatedOn = product.CreatedOn
		productLocal.DefaultPrice = product.DefaultPrice
		productLocal.DeletedBy = &product.DeletedBy
		productLocal.DeletedOn = &product.DeletedOn
		productLocal.DiscountPrice = &product.DiscountPrice
		productLocal.EcommerceCart.CreatedOn = product.TblEcomCart.CreatedOn
		productLocal.EcommerceCart.CustomerID = product.TblEcomCart.CustomerID
		productLocal.EcommerceCart.DeletedOn = &product.TblEcomCart.DeletedOn
		productLocal.EcommerceCart.ID = product.TblEcomCart.ID
		productLocal.EcommerceCart.IsDeleted = product.TblEcomCart.IsDeleted
		productLocal.EcommerceCart.ModifiedOn = &product.TblEcomCart.ModifiedOn
		productLocal.EcommerceCart.ProductID = product.TblEcomCart.ProductID
		productLocal.EcommerceCart.Quantity = product.TblEcomCart.Quantity
		productLocal.ID = product.ID
		productLocal.IsActive = product.IsActive
		productLocal.IsDeleted = product.IsDeleted
		productLocal.ModifiedBy = &product.ModifiedBy
		productLocal.ModifiedOn = &product.ModifiedOn
		productLocal.OrderCustomer = &product.OrderCustomer
		productLocal.OrderID = &product.OrderID
		productLocal.OrderPrice = &product.OrderPrice
		productLocal.OrderQuantity = &product.OrderQuantity
		productLocal.OrderStatus = &product.OrderStatus
		productLocal.OrderTax = &product.OrderTax
		productLocal.OrderTime = &product.OrderTime
		productLocal.OrderUniqueID = &product.OrderUniqueID
		productLocal.PaymentMode = &product.PaymentMode
		productLocal.ProductDescription = product.ProductDescription
		productLocal.ProductImageArray = product.ProductImageArray
		productLocal.ProductImagePath = product.ProductImagePath
		productLocal.ProductName = product.ProductName
		productLocal.ProductSlug = product.ProductSlug
		productLocal.ProductVimeoPath = &product.ProductVimeoPath
		productLocal.ProductYoutubePath = &product.ProductYoutubePath
		productLocal.ShippingDetails = &product.ShippingDetails
		productLocal.Sku = product.Sku
		productLocal.ShippingDetails = &product.ShippingDetails
		productLocal.Tax = product.Tax
		productLocal.Totalcost = product.Totalcost
		productLocal.ViewCount = &product.ViewCount

		if product.ProductImagePath != "" {

			imagePaths := strings.Split(productLocal.ProductImagePath, ",")

			for index, path := range imagePaths {

				modified_path := PathUrl + strings.TrimPrefix(path, "/")

				imagePaths[index] = modified_path
			}

			productLocal.ProductImageArray = imagePaths

		}

		final_OrderedProductList = append(final_OrderedProductList, productLocal)
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

	customer, err := EcomInstance.GetCustomer(memberid)
	if err != nil {

		return &model.EcomOrderedProductDetails{}, err
	}

	var (
		orderedProduct model.EcommerceProduct
		id             int
		slug           string
	)

	if productID != nil {

		id = *productID
	} else if productSlug != nil {

		slug = *productSlug
	}

	product, statuses, err := EcomAuthInstance.GetProductOrderDetailsById(id, slug, customer.Id, orderId)
	if err != nil {

		return &model.EcomOrderedProductDetails{}, err
	}

	orderedProduct.CategoriesID = product.CategoriesID
	orderedProduct.CreatedBy = product.CreatedBy
	orderedProduct.CreatedOn = product.CreatedOn
	orderedProduct.DefaultPrice = product.DefaultPrice
	orderedProduct.DeletedBy = &product.DeletedBy
	orderedProduct.DeletedOn = &product.DeletedOn
	orderedProduct.DiscountPrice = &product.DiscountPrice
	orderedProduct.EcommerceCart.CreatedOn = product.TblEcomCart.CreatedOn
	orderedProduct.EcommerceCart.CustomerID = product.TblEcomCart.CustomerID
	orderedProduct.EcommerceCart.DeletedOn = &product.TblEcomCart.DeletedOn
	orderedProduct.EcommerceCart.ID = product.TblEcomCart.ID
	orderedProduct.EcommerceCart.IsDeleted = product.TblEcomCart.IsDeleted
	orderedProduct.EcommerceCart.ModifiedOn = &product.TblEcomCart.ModifiedOn
	orderedProduct.EcommerceCart.ProductID = product.TblEcomCart.ProductID
	orderedProduct.EcommerceCart.Quantity = product.TblEcomCart.Quantity
	orderedProduct.ID = product.ID
	orderedProduct.IsActive = product.IsActive
	orderedProduct.IsDeleted = product.IsDeleted
	orderedProduct.ModifiedBy = &product.ModifiedBy
	orderedProduct.ModifiedOn = &product.ModifiedOn
	orderedProduct.OrderCustomer = &product.OrderCustomer
	orderedProduct.OrderID = &product.OrderID
	orderedProduct.OrderPrice = &product.OrderPrice
	orderedProduct.OrderQuantity = &product.OrderQuantity
	orderedProduct.OrderStatus = &product.OrderStatus
	orderedProduct.OrderTax = &product.OrderTax
	orderedProduct.OrderTime = &product.OrderTime
	orderedProduct.OrderUniqueID = &product.OrderUniqueID
	orderedProduct.PaymentMode = &product.PaymentMode
	orderedProduct.ProductDescription = product.ProductDescription
	orderedProduct.ProductImageArray = product.ProductImageArray
	orderedProduct.ProductImagePath = product.ProductImagePath
	orderedProduct.ProductName = product.ProductName
	orderedProduct.ProductSlug = product.ProductSlug
	orderedProduct.ProductVimeoPath = &product.ProductVimeoPath
	orderedProduct.ProductYoutubePath = &product.ProductYoutubePath
	orderedProduct.ShippingDetails = &product.ShippingDetails
	orderedProduct.Sku = product.Sku
	orderedProduct.ShippingDetails = &product.ShippingDetails
	orderedProduct.Tax = product.Tax
	orderedProduct.Totalcost = product.Totalcost
	orderedProduct.ViewCount = &product.ViewCount

	if orderedProduct.ProductImagePath != "" {

		imagePaths := strings.Split(orderedProduct.ProductImagePath, ",")

		for index, path := range imagePaths {

			modified_path := PathUrl + strings.TrimPrefix(path, "/")

			imagePaths[index] = modified_path
		}

		orderedProduct.ProductImageArray = imagePaths

	}

	var (
		productOrderStatuses []model.OrderStatus
		orderStatus          model.OrderStatus
	)

	for _, status := range statuses {

		orderStatus.CreatedBy = status.CreatedBy
		orderStatus.CreatedOn = status.CreatedOn
		orderStatus.ID = status.Id
		orderStatus.OrderStatus = status.OrderStatus
		orderStatus.OrderID = status.OrderId

		productOrderStatuses = append(productOrderStatuses, orderStatus)
	}

	return &model.EcomOrderedProductDetails{EcommerceProduct: orderedProduct, OrderStatuses: productOrderStatuses}, nil
}

// func EcommerceOrderPlacement(db *gorm.DB, ctx context.Context, paymentMode string, shippingAddress string, orderProducts []model.OrderProduct, orderSummary *model.OrderSummary) (bool, error) {

// 	c, _ := ctx.Value(ContextKey).(*gin.Context)

// 	memberid := c.GetInt("memberid")

// 	if memberid == 0 {

// 		err := errors.New("unauthorized access")

// 		c.AbortWithError(http.StatusUnauthorized, err)

// 		return false, err

// 	}

// 	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

// 	var customerId int

// 	if err := db.Table("tbl_ecom_customers").Select("id").Where("is_deleted = 0 and member_id = ?", memberid).Scan(&customerId).Error; err != nil {

// 		return false, err
// 	}

// 	unixTime := time.Now().Unix()

// 	var orderplaced model.EcommerceOrder

// 	orderId := "SP" + strconv.Itoa(int(unixTime))

// 	orderplaced.OrderID = orderId

// 	orderplaced.ShippingAddress = shippingAddress

// 	orderplaced.CustomerID = customerId

// 	orderplaced.Status = "placed"

// 	orderplaced.IsDeleted = 0

// 	orderplaced.CreatedOn = currentTime

// 	var totalPrice, totalTax, totalCost int

// 	if orderSummary != nil {

// 		totalPrice, _ = strconv.Atoi(orderSummary.SubTotal)

// 		totalTax, _ = strconv.Atoi(orderSummary.TotalTax)

// 		totalCost, _ = strconv.Atoi(orderSummary.TotalCost)

// 	} else {

// 		for _, product := range orderProducts {

// 			sum := product.Price * product.Quantity

// 			totalPrice += sum

// 			tax := product.Tax * product.Quantity

// 			totalTax += tax

// 		}

// 		totalCost = totalPrice + totalTax
// 	}

// 	orderplaced.Price = totalPrice

// 	orderplaced.Tax = totalTax

// 	orderplaced.TotalCost = totalCost

// 	if err := db.Table("tbl_ecom_product_orders").Create(&orderplaced).Error; err != nil {

// 		c.AbortWithError(http.StatusInternalServerError, err)

// 		return false, err
// 	}

// 	var createorder model.EcommerceOrder

// 	if err := db.Table("tbl_ecom_product_orders").Where("uuid = ?", orderId).First(&createorder).Error; err != nil {

// 		c.AbortWithError(http.StatusInternalServerError, err)

// 		return false, err
// 	}

// 	var orderedProductIds []int

// 	for _, value := range orderProducts {

// 		var productDetails = model.OrderProductDetails{OrderID: createorder.ID, ProductID: value.ProductID, Quantity: value.Quantity, Price: value.Price, Tax: value.Tax}

// 		if err := db.Table("tbl_ecom_product_order_details").Create(&productDetails).Error; err != nil {

// 			c.AbortWithError(http.StatusInternalServerError, err)

// 			return false, err
// 		}

// 		if err := db.Debug().Table("tbl_ecom_products").Where("is_deleted = 0 and is_active = 1 and id = ?", value.ProductID).Update("stock", gorm.Expr("stock - ?", value.Quantity)).Error; err != nil {

// 			return false, err
// 		}

// 		orderedProductIds = append(orderedProductIds, value.ProductID)
// 	}

// 	var orderstatus model.OrderStatus

// 	orderstatus.OrderID = createorder.ID

// 	orderstatus.OrderStatus = "placed"

// 	orderstatus.CreatedBy = customerId

// 	orderstatus.CreatedOn = currentTime

// 	if err := db.Table("tbl_ecom_order_statuses").Create(&orderstatus).Error; err != nil {

// 		c.AbortWithError(http.StatusInternalServerError, err)

// 		return false, err
// 	}

// 	var orderPayment model.OrderPayment

// 	orderPayment.OrderID = createorder.ID

// 	orderPayment.PaymentMode = paymentMode

// 	if err := db.Table("tbl_ecom_order_payments").Create(&orderPayment).Error; err != nil {

// 		c.AbortWithError(http.StatusInternalServerError, err)

// 		return false, err
// 	}

// 	if err := db.Table("tbl_ecom_carts").Where("is_deleted = 0 and product_id in (?) and customer_id = ?", orderedProductIds, customerId).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": currentTime}).Error; err != nil {

// 		c.AbortWithError(http.StatusInternalServerError, err)

// 		return false, err
// 	}

// 	return true, nil
// }

func EcommerceOrderPlacement(db *gorm.DB, ctx context.Context, paymentMode string, shippingAddress string, orderProducts []model.OrderProduct, orderSummary *model.OrderSummary) (bool, error) {

	var (
		customer          ecommerce.TblEcomCustomers
		err               error
		createorder       model.EcommerceOrder
		orderedProductIds []int
	)

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return false, err

	}

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	customer, err = EcomInstance.GetCustomer(memberid)
	if err != nil {

		return false, err
	}

	unixTime := time.Now().Unix()

	var orderplaced ecommerce.EcommerceOrder

	orderId := "SP" + strconv.Itoa(int(unixTime))

	// orderplaced.OrderId = orderId

	orderplaced.ShippingAddress = shippingAddress

	orderplaced.CustomerId = customer.Id

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

	err = EcomAuthInstance.PlaceOrder(orderplaced)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	order, err := EcomAuthInstance.GetOrderByOrderId(orderId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	createorder.CreatedOn = order.CreatedOn
	createorder.CustomerID = order.CustomerId
	createorder.ID = order.Id
	createorder.IsDeleted = order.IsDeleted
	createorder.ModifiedOn = &order.ModifiedOn
	// createorder.OrderID = order.OrderId
	createorder.Price = order.Price
	createorder.ShippingAddress = order.ShippingAddress
	createorder.Status = order.Status
	createorder.Tax = order.Tax
	createorder.TotalCost = order.TotalCost

	for _, value := range orderProducts {

		var productDetails = ecommerce.OrderProduct{OrderId: createorder.ID, ProductId: value.ProductID, Quantity: value.Quantity, Price: value.Price, Tax: value.Tax}

		err = EcomAuthInstance.CreateOrderDetails(productDetails)
		if err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)

			return false, err
		}

		err = EcomAuthInstance.UpdateStock(value.ProductID, value.Quantity)
		if err != nil {

			return false, err

		}

		orderedProductIds = append(orderedProductIds, value.ProductID)
	}

	var orderstatus ecommerce.TblEcomOrderStatuses

	orderstatus.OrderId = createorder.ID

	orderstatus.OrderStatus = "placed"

	orderstatus.CreatedBy = customer.Id

	orderstatus.CreatedOn = currentTime

	err = EcomAuthInstance.CreateOrderStatus(orderstatus)
	if err != nil {

		return false, err
	}

	var orderPayment ecommerce.OrderPayment

	orderPayment.OrderId = createorder.ID

	orderPayment.PaymentMode = paymentMode

	err = EcomAuthInstance.CreateOrderPayment(orderPayment)
	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	err = EcomAuthInstance.DeleteFromCartAfterOrder(orderedProductIds, customer.Id)

	if err != nil {

		c.AbortWithError(http.StatusInternalServerError, err)

		return false, err
	}

	return true, nil
}

func EcommerceCustomerDetails(db *gorm.DB, ctx context.Context) (*model.CustomerDetails, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return &model.CustomerDetails{}, err
	}

	var customerDetails model.CustomerDetails

	customer, err := EcomAuthInstance.GetCustomerDetailsById(memberid)
	if err != nil {

		return &model.CustomerDetails{}, err
	}

	customerDetails.Area = &customer.Area
	customerDetails.City = &customer.City
	customerDetails.Country = &customer.Country
	customerDetails.CreatedBy = customer.CreatedBy
	customerDetails.CreatedOn = customer.CreatedOn
	customerDetails.DeletedOn = &customer.DeletedOn
	customerDetails.Email = customer.Email
	customerDetails.FirstName = customer.FirstName
	customerDetails.HouseNo = &customer.HouseNo
	customerDetails.ID = customer.Id
	customerDetails.IsActive = customer.IsActive
	customerDetails.IsDeleted = &customer.IsDeleted
	customerDetails.LastName = &customer.LastName
	customerDetails.MemberID = &customer.MemberID
	customerDetails.MobileNo = customer.MobileNo
	customerDetails.ModifiedBy = &customer.ModifiedBy
	customerDetails.ModifiedOn = &customer.ModifiedOn
	customerDetails.Password = customer.Password
	customerDetails.ProfileImage = &customer.ProfileImage
	customerDetails.ProfileImagePath = &customer.ProfileImagePath
	customerDetails.State = &customer.State
	customerDetails.StreetAddress = &customer.StreetAddress
	customerDetails.Username = customer.Username
	customerDetails.ZipCode = &customer.ZipCode

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

// func CustomerProfileUpdate(db *gorm.DB, ctx context.Context, customerInput model.CustomerInput) (bool, error) {

// 	// c, _ := ctx.Value(ContextKey).(*gin.Context)

// 	// memberid := c.GetInt("memberid")

// 	// if memberid == 0 {

// 	// 	err := errors.New("unauthorized access")

// 	// 	c.AbortWithError(http.StatusUnauthorized, err)

// 	// 	return false, err
// 	// }

// 	// customerDetails := make(map[string]interface{})

// 	// memberDetails := make(map[string]interface{})

// 	// if customerInput.ProfileImage.IsSet() && customerInput.ProfileImage.Value() != nil {

// 	// 	var fileName, filePath string

// 	// 	storageType, _ := GetStorageType(db)

// 	// 	fileName = customerInput.ProfileImage.Value().Filename

// 	// 	file := customerInput.ProfileImage.Value().File

// 	// 	if storageType.SelectedType == "aws" {

// 	// 		fmt.Printf("aws-S3 storage selected\n")

// 	// 		filePath = "member/" + fileName

// 	// 		err := storage.UploadFileS3(storageType.Aws, customerInput.ProfileImage.Value(), filePath)

// 	// 		if err != nil {

// 	// 			fmt.Printf("image upload failed %v\n", err)

// 	// 			return false, ErrUpload

// 	// 		}

// 	// 	} else if storageType.SelectedType == "local" {

// 	// 		fmt.Printf("local storage selected\n")

// 	// 		b64Data, err := IoReadSeekerToBase64(file)

// 	// 		if err != nil {

// 	// 			return false, err
// 	// 		}

// 	// 		endpoint := "gqlSaveLocal"

// 	// 		url := PathUrl + endpoint

// 	// 		filePath, err = storage.UploadImageToAdminLocal(b64Data, fileName, url)

// 	// 		if err != nil {

// 	// 			return false, ErrUpload
// 	// 		}

// 	// 		fmt.Printf("local stored path: %v\n", filePath)

// 	// 	} else if storageType.SelectedType == "azure" {

// 	// 		fmt.Printf("azure storage selected")

// 	// 	} else if storageType.SelectedType == "drive" {

// 	// 		fmt.Println("drive storage selected")
// 	// 	}

// 	// 	customerDetails["profile_image"] = fileName

// 	// 	memberDetails["profile_image"] = fileName

// 	// 	customerDetails["profile_image_path"] = filePath

// 	// 	memberDetails["profile_image_path"] = filePath

// 	// }

// 	// customerDetails["first_name"] = customerInput.FirstName

// 	// memberDetails["first_name"] = customerInput.FirstName

// 	// customerDetails["email"] = customerInput.Email

// 	// memberDetails["email"] = customerInput.Email

// 	// if customerInput.LastName.IsSet() && customerInput.LastName.Value() != nil {

// 	// 	customerDetails["last_name"] = *customerInput.LastName.Value()

// 	// 	memberDetails["last_name"] = *customerInput.LastName.Value()
// 	// }

// 	// if customerInput.MobileNo.IsSet() && customerInput.MobileNo.Value() != nil {

// 	// 	customerDetails["mobile_no"] = *customerInput.MobileNo.Value()

// 	// 	memberDetails["mobile_no"] = *customerInput.MobileNo.Value()
// 	// }

// 	// if customerInput.Username.IsSet() && customerInput.Username.Value() != nil {

// 	// 	customerDetails["username"] = *customerInput.Username.Value()

// 	// 	memberDetails["username"] = *customerInput.Username.Value()
// 	// }

// 	// if customerInput.IsActive.IsSet() && customerInput.IsActive.Value() != nil {

// 	// 	customerDetails["is_active"] = *customerInput.IsActive.Value()

// 	// 	memberDetails["is_active"] = *customerInput.IsActive.Value()
// 	// }

// 	// if customerInput.StreetAddress.IsSet() && customerInput.StreetAddress.Value() != nil {

// 	// 	customerDetails["street_address"] = *customerInput.StreetAddress.Value()
// 	// }

// 	// if customerInput.City.IsSet() && customerInput.City.Value() != nil {

// 	// 	customerDetails["city"] = *customerInput.City.Value()
// 	// }

// 	// if customerInput.Country.IsSet() && customerInput.Country.Value() != nil {

// 	// 	customerDetails["country"] = *customerInput.Country.Value()
// 	// }

// 	// if customerInput.State.IsSet() && customerInput.State.Value() != nil {

// 	// 	customerDetails["state"] = *customerInput.State.Value()
// 	// }

// 	// if customerInput.ZipCode.IsSet() && customerInput.ZipCode.Value() != nil {

// 	// 	customerDetails["zip_code"] = *customerInput.ZipCode.Value()
// 	// }

// 	// if customerInput.Password.IsSet() && customerInput.Password.Value() != nil && *customerInput.Password.Value() != "" {

// 	// 	hashpass, err := HashingPassword(*customerInput.Password.Value())

// 	// 	if err != nil {

// 	// 		return false, ErrPassHash
// 	// 	}

// 	// 	customerDetails["password"] = hashpass

// 	// 	memberDetails["password"] = hashpass
// 	// }

// 	// currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

// 	// customerDetails["modified_on"] = currentTime

// 	// memberDetails["modified_on"] = currentTime

// 	// customerDetails["modified_by"] = memberid

// 	// memberDetails["modified_by"] = memberid

// 	// if err := db.Debug().Table("tbl_ecom_customers").Where("is_deleted = 0 and member_id = ?", memberid).UpdateColumns(&customerDetails).Error; err != nil {

// 	// 	return false, err
// 	// }

// 	// if err := db.Debug().Table("tbl_members").Where("is_deleted = 0 and id = ?", memberid).UpdateColumns(&memberDetails).Error; err != nil {

// 	// 	return false, err
// 	// }

// 	return true, nil
// }

func CustomerProfileUpdate(db *gorm.DB, ctx context.Context, customerInput model.CustomerInput) (bool, error) {

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return false, err
	}

	customerDetails := make(map[string]interface{})

	memberDetails := make(map[string]interface{})

	if customerInput.ProfileImage.IsSet() && customerInput.ProfileImage.Value() != nil {

		var (
			fileName, filePath string
			storageType        StorageType
			err                error
		)

		var imageData = *customerInput.ProfileImage.Value()

		storageType, err = GetStorageType(db)
		if err != nil {

			return false, err
		}

		if imageData != "" {

			isValidBase64, base64Data, extension := IsValidBase64(imageData)

			if isValidBase64 && base64Data != "" {

				randNum := strconv.Itoa(int(time.Now().Unix()))

				fileName = "IMG-" + randNum + "." + extension

				if storageType.SelectedType == "aws" {

					fmt.Printf("aws-S3 storage selected\n")

					filePath = "member/" + fileName

					err = storage.UploadFileS3(storageType.Aws, nil, base64Data, filePath)
					if err != nil {

						fmt.Printf("image upload failed %v\n", err)

						return false, ErrUpload

					}
				} else if storageType.SelectedType == "azure" {

					fmt.Printf("azure storage selected")

				} else if storageType.SelectedType == "drive" {

					fmt.Println("drive storage selected")
				}

			} else if strings.Contains(imageData, "image-resize?name") {

				filePath = strings.ReplaceAll(imageData, "image-resize?name=", "")

				indexOf := strings.Index(filePath, "/")

				fileName = filePath[indexOf+1:]
			}
			customerDetails["profile_image"] = fileName

			memberDetails["profile_image"] = fileName

			customerDetails["profile_image_path"] = filePath

			memberDetails["profile_image_path"] = filePath
		}

	}

	customerDetails["first_name"] = customerInput.FirstName

	memberDetails["first_name"] = customerInput.FirstName

	customerDetails["email"] = customerInput.Email

	memberDetails["email"] = customerInput.Email

	if customerInput.LastName.IsSet() && customerInput.LastName.Value() != nil {

		customerDetails["last_name"] = *customerInput.LastName.Value()

		memberDetails["last_name"] = *customerInput.LastName.Value()
	}

	if customerInput.MobileNo.IsSet() && customerInput.MobileNo.Value() != nil {

		customerDetails["mobile_no"] = *customerInput.MobileNo.Value()

		memberDetails["mobile_no"] = *customerInput.MobileNo.Value()
	}

	if customerInput.Username.IsSet() && customerInput.Username.Value() != nil {

		customerDetails["username"] = *customerInput.Username.Value()

		memberDetails["username"] = *customerInput.Username.Value()
	}

	if customerInput.IsActive.IsSet() && customerInput.IsActive.Value() != nil {

		customerDetails["is_active"] = *customerInput.IsActive.Value()

		memberDetails["is_active"] = *customerInput.IsActive.Value()
	}

	if customerInput.StreetAddress.IsSet() && customerInput.StreetAddress.Value() != nil {

		customerDetails["street_address"] = *customerInput.StreetAddress.Value()
	}

	if customerInput.City.IsSet() && customerInput.City.Value() != nil {

		customerDetails["city"] = *customerInput.City.Value()
	}

	if customerInput.Country.IsSet() && customerInput.Country.Value() != nil {

		customerDetails["country"] = *customerInput.Country.Value()
	}

	if customerInput.State.IsSet() && customerInput.State.Value() != nil {

		customerDetails["state"] = *customerInput.State.Value()
	}

	if customerInput.ZipCode.IsSet() && customerInput.ZipCode.Value() != nil {

		customerDetails["zip_code"] = *customerInput.ZipCode.Value()
	}

	if customerInput.Password.IsSet() && customerInput.Password.Value() != nil && *customerInput.Password.Value() != "" {

		hashpass, err := HashingPassword(*customerInput.Password.Value())

		if err != nil {

			return false, ErrPassHash
		}

		customerDetails["password"] = hashpass

		memberDetails["password"] = hashpass
	}

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	customerDetails["modified_on"] = currentTime

	memberDetails["modified_on"] = currentTime

	customerDetails["modified_by"] = memberid

	memberDetails["modified_by"] = memberid

	err := EcomAuthInstance.UpdateCustomerAndMemberDetails(memberid, memberDetails, customerDetails)
	if err != nil {

		return false, err
	}

	return true, nil
}

// func UpdateProductViewCount(db *gorm.DB, ctx context.Context, productID *int, productSlug *string) (bool, error) {

// 	c, _ := ctx.Value(ContextKey).(*gin.Context)

// 	if productID == nil && productSlug == nil {

// 		return false, ErrMandatory
// 	}

// 	query := db.Debug().Table("tbl_ecom_products").Where("is_deleted = 0 and is_active = 1")

// 	if productID != nil {

// 		query = query.Where("id = ?", *productID)

// 	} else if productSlug != nil {

// 		query = query.Where("product_slug = ?", *productSlug)
// 	}

// 	err := query.Update("view_count", gorm.Expr("view_count + 1")).Error

// 	if err != nil {

// 		c.AbortWithError(500, err)

// 		return false, err
// 	}

// 	return true, nil
// }

func UpdateProductViewCount(db *gorm.DB, ctx context.Context, productID *int, productSlug *string) (bool, error) {

	var (
		productId    int
		productSlugg string
	)

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	if productID == nil && productSlug == nil {

		return false, ErrMandatory
	}

	if productID != nil {

		productId = *productID

	} else if productSlug != nil {

		productSlugg = *productSlug
	}

	err := EcomAuthInstance.UpdateProductViewCount(productId, productSlugg)

	if err != nil {

		c.AbortWithError(500, err)

		return false, err
	}

	return true, nil
}

// func EcommerceOrderStatusNames(db *gorm.DB, ctx context.Context) ([]model.OrderStatusNames, error) {

// 	c, _ := ctx.Value(ContextKey).(*gin.Context)

// 	memberid := c.GetInt("memberid")

// 	if memberid == 0 {

// 		err := errors.New("unauthorized access")

// 		c.AbortWithError(http.StatusUnauthorized, err)

// 		return []model.OrderStatusNames{}, err
// 	}

// 	var orderStatus []model.OrderStatusNames

// 	if err := db.Debug().Table("tbl_ecom_statuses").Find(&orderStatus).Error; err != nil {

// 		return []model.OrderStatusNames{}, err
// 	}

// 	return orderStatus, nil
// }

func EcommerceOrderStatusNames(db *gorm.DB, ctx context.Context) ([]model.OrderStatusNames, error) {

	var (
		orderStatusNames []model.OrderStatusNames
		orderStatusLocal model.OrderStatusNames
	)

	c, _ := ctx.Value(ContextKey).(*gin.Context)

	memberid := c.GetInt("memberid")

	if memberid == 0 {

		err := errors.New("unauthorized access")

		c.AbortWithError(http.StatusUnauthorized, err)

		return []model.OrderStatusNames{}, err
	}

	orderStatuses, err := EcomAuthInstance.GetOrderStatusNames()
	if err != nil {

		return []model.OrderStatusNames{}, err
	}

	for _, orderStatus := range orderStatuses {

		orderStatusLocal.CreatedBy = orderStatus.CreatedBy
		orderStatusLocal.CreatedOn = orderStatus.CreatedOn
		orderStatusLocal.Description = &orderStatus.Description
		orderStatusLocal.ID = orderStatus.Id
		orderStatusLocal.IsActive = orderStatus.IsActive
		orderStatusLocal.IsDeleted = orderStatus.IsDeleted
		orderStatusLocal.ModifiedBy = &orderStatus.ModifiedBy
		orderStatusLocal.ModifiedOn = &orderStatus.ModifiedOn
		orderStatusLocal.Priority = orderStatus.Priority
		orderStatusLocal.Status = orderStatus.Status

		orderStatusNames = append(orderStatusNames, orderStatusLocal)

	}

	return orderStatusNames, nil
}
