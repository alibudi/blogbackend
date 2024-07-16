package controller

import (
	"fmt"

	"github.com/alibudi/blogbackend/database"
	"github.com/alibudi/blogbackend/models"
	"github.com/gofiber/fiber/v2"
)

func CreatePost(c *fiber.Ctx) error {
	var blogpost models.Blog
	if err:=c.BodyParser(&blogpost);err!=nil {
		fmt.Println("Unable to parse Body")
	}
	if err:=database.DB.Create(&blogpost).Error; err!=nil{
		c.Status(400)
		return c.JSON(fiber.Map{
			"Message": "Invalid Payload",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Congralutation, Your post live",
	})
}