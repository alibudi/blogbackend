package controller

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alibudi/blogbackend/database"
	"github.com/alibudi/blogbackend/models"
	"github.com/alibudi/blogbackend/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`[a-z0-9.%\-]+@[a-z0-9.%\-]+\.[a-z0-9.%\-]+`)
	return Re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body:", err)
		c.Status(400)
		return c.JSON(fiber.Map{
			"Message": "Unable to parse request body",
		})
	}

	// Debug log to check the received data
	fmt.Println("Received data:", data)

	password, ok := data["password"].(string)
	if !ok {
		fmt.Println("Password field missing or not a string")
		c.Status(400)
		return c.JSON(fiber.Map{
			"Message": "Password is required and must be a string",
		})
	}

	// Debug log to check the password length
	fmt.Println("Password length:", len(password))

	if len(password) <= 6 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"Message": "Password must be greater than 6 characters",
		})
	}

	email, ok := data["email"].(string)
	if !ok || !validateEmail(strings.TrimSpace(email)) {
		fmt.Println("Invalid or missing email field")
		c.Status(400)
		return c.JSON(fiber.Map{
			"Message": "Invalid Email Address",
		})
	}

	var userData models.User
	database.DB.Where("email = ?", strings.TrimSpace(email)).First(&userData)
	if userData.Id != 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"Message": "Email already exists",
		})
	}

	firstName, _ := data["first_name"].(string)
	lastName, _ := data["last_name"].(string)
	phone, _ := data["phone"].(string)

	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Email:     strings.TrimSpace(email),
	}
	user.SetPassword(password)

	if err := database.DB.Create(&user).Error; err != nil {
		log.Println("Error creating user:", err)
		c.Status(500)
		return c.JSON(fiber.Map{
			"Message": "Error creating user",
		})
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"user":    user,
		"Message": "Account created successfully",
	})
}

func Login(c *fiber.Ctx)error{
	var data map[string]string
	if err:=c.BodyParser(&data);err!=nil {
		fmt.Println("Unable to parse Body")
		
	}
	var user models.User
	database.DB.Where("email=?", data["email"]).First(&user)
	if user.Id ==0{
		c.Status(404)
		return c.JSON(fiber.Map{
			"Message": "Email address doesn't exit, kindly create account",
		})
	}
	if err:=user.ComparePassword(data["password"]); err !=nil{
		c.Status(400)
		return c.JSON(fiber.Map{
			"Message": "Incorrect password",
		})
	}
	token,err:=util.GenerateJWT(strconv.Itoa(int(user.Id)),)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"Message": "You have successfully login",
		"user": user,
	})
} 

type Claims struct{
	jwt.StandardClaims
}