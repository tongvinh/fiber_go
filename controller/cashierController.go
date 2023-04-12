package controller

import (
	db "fiber_rest_api/config"
	"fiber_rest_api/model"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

func CreateCashier(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("registeration error in post request %v", err)
	}

	if data["name"] == "" || data["passcode"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Cashier name is required",
			"error":   map[string]any{},
		})
	}
	//passCode := strconv.Itoa(rand.Intn(1000000))
	//fmt.Println("passCode:::", passCode)
	cashier := model.Cashier{
		Name:     data["name"],
		Passcode: data["passcode"],
	}
	db.DB.Create(&cashier)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    cashier,
	})
}

func GetCashierDetail(c *fiber.Ctx) error {
	cashierId := c.Params("cashierId")

	var cashier model.Cashier
	db.DB.Select("id, name").Where("id=?", cashierId).First(&cashier)
	cashierData := make(map[string]any)
	cashierData["cashierId"] = cashier.Id
	cashierData["name"] = cashier.Name

	if cashier.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier not found",
			"error":   map[string]any{},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    cashierData,
	})
}

func DeleteCashier(c *fiber.Ctx) error {
	cashierId := c.Params("cashierId")
	var cashier model.Cashier
	db.DB.Where("id=?", cashierId).First(&cashier)

	if cashier.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier not found",
			"error":   map[string]any{},
		})
	}
	db.DB.Where("id=?", cashierId).Delete(&cashier)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func UpdateCashier(c *fiber.Ctx) error {
	cashierId := c.Params("cashierId")
	var cashier model.Cashier

	db.DB.First(&cashier, "id=?", cashierId)
	if cashier.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier not found",
		})
	}

	var updateCashierData model.Cashier
	c.BodyParser(&updateCashierData)
	if updateCashierData.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Cashier name is required",
			"error":   map[string]any{},
		})
	}

	cashier.Name = updateCashierData.Name
	db.DB.Save(&cashier)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    cashier,
	})
}

// Cashiers struct with tow values
type Cashiers struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func CashiersList(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var cashier []Cashiers
	db.DB.Select("*").Limit(limit).Offset(skip).Find(&cashier).Count(&count)
	metaMap := map[string]any{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}
	cashiersData := map[string]any{
		"cashiers": cashier,
		"meta":     metaMap,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    cashiersData,
	})
}
