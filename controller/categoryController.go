package controller

import (
	db "fiber_rest_api/config"
	"fiber_rest_api/middleware"
	"fiber_rest_api/model"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

// Category struct with two values
type Category struct {
	Id   uint   `json:"categoryId"`
	Name string `json:"name"`
}

func CreateCategory(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)

	if err != nil {
		log.Fatalf("category error in post request %v", err)
	}

	if data["name"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Category name is required",
			"error":   map[string]any{},
		})
	}

	category := model.Category{
		Name: data["name"],
	}
	db.DB.Create(&category)
	//result := db.DB.Create(&category)

	/*if result.RowsAffected == 0{
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category insertion failed",
		})
	}*/

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    category,
	})
}

func GetCategoryDetails(c *fiber.Ctx) error {
	categoryId := c.Params("categoryId")

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
	var category model.Category
	db.DB.Select("id, name").Where("id=?", categoryId).First(&category)
	categoryData := make(map[string]any)
	categoryData["categoryId"] = category.Id
	categoryData["name"] = category.Name

	if category.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "No category found",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    categoryData,
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	categoryId := c.Params("categoryId")
	var category model.Category
	db.DB.Where("id=?", categoryId).First(&category)

	if category.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"error":   map[string]any{},
		})
	}

	db.DB.Where("id = ?", categoryId).Delete(&category)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	categoryId := c.Params("categoryId")
	var category model.Category

	db.DB.Find(&category, " id = ?", categoryId)

	if category.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category not exist against this id",
		})
	}

	var updateCashierData model.Category
	c.BodyParser(&updateCashierData)
	if updateCashierData.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Category name is required",
			"error":   map[string]any{},
		})
	}

	category.Name = updateCashierData.Name
	db.DB.Save(&category)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    category,
	})
}

func CategoryList(c *fiber.Ctx) error {
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

	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var category []Category
	db.DB.Select("id, name").Limit(limit).Offset(skip).Find(&category).Count(&count)
	metaMap := map[string]any{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}
	categoriesData := map[string]any{
		"categories": category,
		"meta":       metaMap,
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    categoriesData,
	})
}

/*func CatL(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == ""{
		return  c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error": map[string]any{},
		})
	}

	if err := middleware.AuthenticateToken(middleware.SplitToken(auth)); err != nil{
		return c.Status(401).JSON(fiber.Map{
			"status": "error",
			"message": "Token expired or invalid",
			"error": map[string]any{},
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var countes int64

	var category []model.Category
	db.DB.Select([]string{"category_id, name"}).Limit(limit).Offset(skip).Find(&category).Count(&countes)

	return c.Status(404).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": category,
		"meta": map[string]any{
			"total": countes,
			"limit": limit,
			"skip": skip,
		},
	})
}*/
