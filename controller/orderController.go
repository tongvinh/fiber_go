package controller

import (
	db "fiber_rest_api/config"
	"fiber_rest_api/middleware"
	"fiber_rest_api/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func CreateOrder(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]any{},
		})
	}
	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]any{},
		})
	}
	//Token authenticate

	type products struct {
		ProductId int `json:"productId"`
		Quantity  int `json:"qty"`
	}

	body := struct {
		PaymentId int        `json:"paymentId"`
		TotalPaid int        `json:"totalPaid"`
		Products  []products `json:"products"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Empty Body",
			"error":   map[string]any{},
		})
	}

	Prodresponse := make([]*model.ProductResponseOrder, 0)

	var TotalInvoicePrice = struct {
		ttprice int
	}{}

	productsIds := ""
	quantities := ""
	for _, v := range body.Products {
		totalPrice := 0
		productsIds = productsIds + "," + strconv.Itoa(v.ProductId)
		quantities = quantities + "," + strconv.Itoa(v.Quantity)

		prods := model.ProductOrder{}
		var discount model.Discount
		db.DB.Take("products").Where("id=?", v.ProductId).Find(&prods)
		db.DB.Where("id = ?", prods.DiscountId).Find(&discount)
		discCount := 0

		if discount.Type == "BUY_N" {
			totalPrice = prods.Price * v.Quantity

			discCount = totalPrice - discount.Result
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + discCount
		}

		if discount.Type == "PERCENT" {
			totalPrice = prods.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discCount = totalPrice - percentage
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + discCount
		}

		Prodresponse = append(Prodresponse,
			&model.ProductResponseOrder{
				ProductId:        prods.Id,
				Name:             prods.Name,
				Price:            prods.Price,
				Discount:         discount,
				Qty:              v.Quantity,
				TotalNormalPrice: prods.Price,
				TotalFinalPrice:  discCount,
			})
	}
	orderResp := model.Order{
		CashierID:      1,
		PaymentTypesId: body.PaymentId,
		TotalPrice:     TotalInvoicePrice.ttprice,
		TotalPaid:      body.TotalPaid,
		TotalReturn:    body.TotalPaid - TotalInvoicePrice.ttprice,
		ReceiptId:      "R000" + strconv.Itoa(rand.Intn(1000)),
		ProductId:      productsIds,
		Quantities:     quantities,
		UpdatedAt:      time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
	}
	db.DB.Create(&orderResp)

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"success": true,
		"data": map[string]any{
			"order":    orderResp,
			"products": Prodresponse,
		},
	})
}

func SubTotalOrder(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("authorization")
	if headerToken == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]any{},
		})
	}
	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]any{},
		})
	}
	//Token authenticate
	type products struct {
		ProductId int `json:"productId"`
		Quantity  int `json:"qty"`
	}

	body := struct {
		Products []products `json:"products"`
	}{}

	if err := c.BodyParser(&body.Products); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Empty Body",
		})
	}

	Prodresponse := make([]*model.ProductResponseOrder, 0)

	var TotalInvoicePrice = struct {
		ttprice int
	}{}

	for _, v := range body.Products {
		totalPrice := 0

		prods := model.ProductOrder{}
		var discount model.Discount
		db.DB.Take("products").Where("id=?", v.ProductId).First(&prods)
		db.DB.Where("id = ?", prods.DiscountId).Find(&discount)

		disc := 0
		if discount.Type == "PERCENT" {
			totalPrice = prods.Price * v.Quantity // 5000*3=15000
			percentage := totalPrice * discount.Result / 100
			disc = totalPrice - percentage
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + disc
		}
		if discount.Type == "BUY_N" {
			totalPrice = prods.Price * v.Quantity //5000*3=15000
			disc = totalPrice - discount.Result
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + disc
		}

		Prodresponse = append(Prodresponse,
			&model.ProductResponseOrder{
				ProductId:        prods.Id,
				Name:             prods.Name,
				Price:            prods.Price,
				Discount:         discount,
				Qty:              v.Quantity,
				TotalNormalPrice: prods.Price,
				TotalFinalPrice:  disc,
			})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"success": true,
		"data": map[string]any{
			"subTotal": TotalInvoicePrice.ttprice,
			"products": Prodresponse,
		},
	})
}

func CheckOrder(c *fiber.Ctx) error {
	param := c.Params("orderId")

	var order model.Order
	db.DB.Where("id=?", param).First(&order)
	if order.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Order does not exist",
		})
	}

	if order.IsDownload == 0 {
		return c.Status(200).JSON(fiber.Map{
			"status":  true,
			"message": "success",
			"data": map[string]any{
				"isDownloaded": false,
			},
		})
	}

	if order.IsDownload == 1 {
		return c.Status(200).JSON(fiber.Map{
			"status":  true,
			"message": "success",
			"data": map[string]any{
				"isDownloaded": true,
			},
		})
	}
	return nil
}

func OrderDetail(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]any{},
		})
	}

	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]any{},
		})
	}
	//Token authenticate

	param := c.Params("orderId")

	var order model.Order
	db.DB.Where("id=?", param).First(&order)

	if order.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Not Found",
			"error":   map[string]any{},
		})
	}
	productIds := strings.Split(order.ProductId, ",")
	TotalProducts := make([]*model.Product, 0)

	for i := 1; i < len(productIds); i++ {
		prods := model.Product{}
		db.DB.Where("id=?", productIds[i]).Find(&prods)
		TotalProducts = append(TotalProducts, &prods)
	}

	cashier := model.Cashier{}
	db.DB.Where("id =?", order.CashierID).Find(&cashier)

	paymentType := model.PaymentType{}
	db.DB.Where("id=?", order.PaymentTypesId).Find(&paymentType)

	orderTable := model.Order{}
	db.DB.Where("id=?", order.Id).Find(&orderTable)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": map[string]any{
			"order": map[string]any{
				"orderId":        order.Id,
				"cashiersId":     order.CashierID,
				"paymentTypesId": order.PaymentTypesId,
				"totalPrice":     order.TotalPrice,
				"totalPaid":      order.TotalPaid,
				"totalReturn":    order.TotalReturn,
				"receiptId":      order.ReceiptId,
				"createdAt":      order.CreatedAt,
				"cashier":        cashier,
				"payment_type":   paymentType,
			},
			"products": TotalProducts,
		},
		"message": "Success",
	})
}

func OrdersList(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var order []model.Order

	db.DB.Select("*").Limit(limit).Offset(skip).Find(&order).Count(&count)

	type OrderList struct {
		OrderId        int               `json:"orderId"`
		CashierID      int               `json:"cashiersId"`
		PaymentTypesId int               `json:"paymentTypesId"`
		TotalPrice     int               `json:"totalPrice"`
		TotalPaid      int               `json:"totalPaid"`
		TotalReturn    int               `json:"totalReturn"`
		ReceiptId      string            `json:"receiptId"`
		CreatedAt      time.Time         `json:"createdAt"`
		Payments       model.PaymentType `json:"payment_type"`
		Cashiers       model.Cashier     `json:"cashier"`
	}
	OrderResponse := make([]*OrderList, 0)

	for _, v := range order {
		cashier := model.Cashier{}
		db.DB.Where("id = ?", v.CashierID).Find(&cashier)
		paymentType := model.PaymentType{}
		db.DB.Where("id = ?", v.PaymentTypesId).Find(&paymentType)

		OrderResponse = append(OrderResponse, &OrderList{
			OrderId:        v.Id,
			CashierID:      v.CashierID,
			PaymentTypesId: v.PaymentTypesId,
			TotalPrice:     v.TotalPrice,
			TotalPaid:      v.TotalPaid,
			TotalReturn:    v.TotalReturn,
			ReceiptId:      v.ReceiptId,
			CreatedAt:      v.CreatedAt,
			Payments:       paymentType,
			Cashiers:       cashier,
		})
	}

	return c.Status(404).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    OrderResponse,
		"meta": map[string]any{
			"total": count,
			"limit": limit,
			"skip":  skip,
		},
	})
}

func DownloadOrder(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Token expired or invalid",
		})
	}

	//Token authenticate
	param := c.Params("orderId")

	var order model.Order
	db.DB.Where("id =?", param).First(&order)

	if order.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Order not found",
			"error":   map[string]any{},
		})
	}
	productIds := strings.Split(order.ProductId, ",")

	TotalProducts := make([]*model.Product, 0)

	for i := 1; i < len(productIds); i++ {
		prods := model.Product{}
		db.DB.Where("id=?", productIds[i]).Find(&prods)
		TotalProducts = append(TotalProducts, &prods)
	}
	cashier := model.Cashier{}
	db.DB.Where("id=?", order.CashierID).Find(&cashier)
	paymentType := model.PaymentType{}

	db.DB.Where("id=?", order.PaymentTypesId).Find(&paymentType)
	orderTable := model.Order{}
	db.DB.Where("id = ?", order.Id).Find(&orderTable)

	//pdf Generating
	twoDarray := [][]string{{}}
	quantities := strings.Split(order.Quantities, ",")
	quantities = quantities[1:]
	for i := 0; i < len(TotalProducts); i++ {

		s_array := []string{}
		s_array = append(s_array, TotalProducts[i].Sku)

		s_array = append(s_array, TotalProducts[i].Name)
		s_array = append(s_array, quantities[i])
		s_array = append(s_array, strconv.Itoa(TotalProducts[i].Price))
		twoDarray = append(twoDarray, s_array)

	}
	begin := time.Now()
	grayColor := getGrayColor()
	whiteColor := color.NewWhite()
	header := getHeader()
	contents := twoDarray

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	//m.SetBorder(true)

	//Top Heading
	m.SetBackgroundColor(grayColor)
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Order Invoice #"+strconv.Itoa(order.Id), props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})
	m.SetBackgroundColor(whiteColor)

	//Table setting
	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 4, 2, 3},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 4, 2, 3},
		},
		Align:                consts.Center,
		AlternatedBackground: &grayColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})
	//Total price
	m.Row(20, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalPrice), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})
	m.Row(21, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total Paid:", props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalPaid), props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})

	m.Row(22, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total Return", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalReturn), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})

	//Invoice creation
	currentTime := time.Now()
	pdfFileName := "invoice-" + currentTime.Format("2006-Jan-02")
	err := m.OutputFileAndClose(pdfFileName + ".pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))

	//update recepit is downloaded to 1 means true
	db.DB.Table("orders").Where("id=?", order.Id).Update("is_download", 1)
	return c.Status(200).JSON(fiber.Map{
		"Success": true,
		"Message": "Success",
	})
}

func getHeader() []string {
	return []string{"Product Sku", "Name", "Qty", "Price"}
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}
