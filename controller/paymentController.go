package controller

import (
	db "fiber_rest_api/config"
	"fiber_rest_api/middleware"
	"fiber_rest_api/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

type Payment struct {
	Id            uint   `json:"paymentId"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	PaymentTypeId int    `json:"payment_type_id"`
	Logo          string `json:"logo"`
}

func CreatePayment(c *fiber.Ctx) error {
	var data map[string]string
	paymentError := c.BodyParser(&data)
	if paymentError != nil {
		log.Fatalf("payment error in post request %v", paymentError)
	}

	if data["name"] == "" || data["type"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Payment name and type is required",
			"error":   map[string]any{},
		})
	}

	var paymentTypes model.PaymentType
	db.DB.Where("name", data["type"]).First(&paymentTypes)
	payment := model.Payment{
		Name:          data["name"],
		Type:          data["type"],
		PaymentTypeId: int(paymentTypes.Id),
		Logo:          data["logo"],
	}
	db.DB.Create(&payment)

	/*result := db.DB.Create(&payment)

	if result.RowsAffected == 0{
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "payment insertion failed",
		})
	}*/

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    payment,
	})
}

func PaymentList(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Token not found",
		})
	}
	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	//Token authenticate

	//subtotal, _ := strconv.Atoi(c.Query("subtotal"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var payment []Payment
	db.DB.Select("id, name, payment_type_id, logo, created_at, updated_at").Limit(limit).Offset(skip).Find(&payment).Count(&count)
	metaMap := map[string]any{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}
	categoriesData := map[string]any{
		"payments": payment,
		"meta":     metaMap,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    categoriesData,
	})
}

func GetPaymentDetails(c *fiber.Ctx) error {
	paymentId := c.Params("paymentId")

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

	var payment model.Payment
	db.DB.Where("id=?", paymentId).First(&payment)

	if payment.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment not found",
			"error":   map[string]any{},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    payment,
	})
}

func DeletePayment(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Token not found",
		})
	}
	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	//Token authenticate

	paymentId := c.Params("paymentId")
	var payment model.Payment

	db.DB.First(&payment, paymentId)
	if payment.Name == "" {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"Message": "No payment found against this payment id",
		})
	}

	result := db.DB.Delete(&payment)
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "payment removing failed",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
	})
}

func UpdatePayment(c *fiber.Ctx) error {
	paymentId := c.Params("paymentId")
	fmt.Println("-----------------------------------")
	fmt.Println("---------------Params payment id--------", paymentId)
	fmt.Println("-----------------------------------")
	var totalPayment model.Payment
	db.DB.Find(&totalPayment)

	fmt.Println("-----------------------------------")
	fmt.Println("---------------All payments--------------------", totalPayment)
	fmt.Println("-----------------------------------")
	var payment model.Payment

	db.DB.Find(&payment, "id = ?", paymentId)

	//if payment.Id == 0 {
	//	return c.Status(404).JSON(fiber.Map{
	//		"success": false,
	//		"Message": "Payment not exist against this id",
	//		"error":   map[string]interface{}{},
	//	})
	//}

	var updatePaymentData model.Payment
	c.BodyParser(&updatePaymentData)
	if updatePaymentData.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Payment name is required",
			"error":   map[string]any{},
		})
	}

	var paymentTypeId int
	if updatePaymentData.Type == "CASH" {
		paymentTypeId = 1
	}
	if updatePaymentData.Type == "E-WALLET" {
		paymentTypeId = 2
	}
	if updatePaymentData.Type == "EDC" {
		paymentTypeId = 3
	}
	fmt.Println(paymentTypeId)
	payment.Name = updatePaymentData.Name
	payment.Type = updatePaymentData.Type
	payment.PaymentTypeId = paymentTypeId
	payment.Logo = updatePaymentData.Logo

	db.DB.Save(&payment)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    payment,
	})
}
