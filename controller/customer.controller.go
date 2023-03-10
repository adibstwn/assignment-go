package controller

import (
	"strconv"
	"time"

	"github.com/adibSetiawann/transaction-api-go/model"
	"github.com/adibSetiawann/transaction-api-go/service/customer"
	"github.com/gofiber/fiber/v2"
)

type CustomerController struct {
	customerService service.CustomerService
}

func NewCustomerController(customerService *service.CustomerService) CustomerController {
	return CustomerController{customerService: *customerService}
}

func (mc *CustomerController) Login(c *fiber.Ctx) error {
	customerReq := new(model.LoginForm)
	var status int

	if err := c.BodyParser(customerReq); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": "request can't go on",
		})
	}

	isErrorValidation := mc.customerService.Validation(*customerReq)
	if isErrorValidation != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": isErrorValidation.Error(),
		})
	}

	token, errToken := mc.customerService.Login(customerReq)
	if errToken != nil {
		c.JSON(fiber.Map{
			"error":   errToken,
			"message": token,
		})
	}

	if token == "please input correct password" {
		status = 404
	} else {
		status = 200
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(status).JSON(fiber.Map{
		"token": "success login",
	})
}

func (mc *CustomerController) Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func (mc *CustomerController) Create(c *fiber.Ctx) error {
	customerReq := new(model.CreateCustomer)

	if err := c.BodyParser(customerReq); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": "request can't go on",
		})
	}

	isErrorValidation := mc.customerService.Validation(*customerReq)
	if isErrorValidation != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": isErrorValidation.Error(),
		})
	}

	customer, errCreate := mc.customerService.Create(*customerReq)
	if errCreate != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "create data successfully",
		"data":    customer,
	})
}

func (mc *CustomerController) Update(c *fiber.Ctx) error {
	customerId := c.Params("id")
	customerReq := new(model.UpdateCustomer)

	if err := c.BodyParser(customerReq); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": "some field is wrong",
		})
	}

	isErrorValidation := mc.customerService.Validation(*customerReq)
	if isErrorValidation != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": isErrorValidation.Error(),
		})
	}

	intId, _ := strconv.ParseInt(customerId, 10, 64)
	customer, errCreate := mc.customerService.Update(intId, *customerReq)

	if errCreate != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "customer not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "update data successfully",
		"data":    customer,
	})
}

func (mc *CustomerController) GetById(c *fiber.Ctx) error {
	customerId := c.Params("id")

	customers, err := mc.customerService.GetById(customerId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "customer not found in database",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": customers,
	})
}

func (mc *CustomerController) GetAll(c *fiber.Ctx) error {
	customers, _ := mc.customerService.GetAllData()

	return c.Status(200).JSON(fiber.Map{
		"data": customers,
	})
}

func (mc *CustomerController) Remove(c *fiber.Ctx) error {
	customerId := c.Params("id")
	err := mc.customerService.Remove(customerId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "customer not found in database",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "delete success",
	})
}
