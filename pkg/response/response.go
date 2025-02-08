package response

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Success(c *fiber.Ctx, data interface{}, statusCode ...int) error {
	code := fiber.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	return c.Status(code).JSON(fiber.Map{
		"status": "success",
		"data":   data,
		"code":   code, 
	})
}

func Error(c *fiber.Ctx, message string, statusCode ...int) error {
	code := fiber.StatusInternalServerError
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"code":    code,
	})
}

func ValidationError(c *fiber.Ctx, errors interface{}) error {
	code := fiber.StatusBadRequest
	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"errors":  errors,
		"code":    code, 
	})
}

var validate = validator.New()

func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)

	if err == nil {
		return nil
	}

	errs := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errs[err.Field()] = err.Tag()
	}

	return errs
}

func HandleValidationError(c *fiber.Ctx, err string) error {
	return ValidationError(c, err)
}

func ParseAndValidateRequest(c *fiber.Ctx, req interface{}) error {
	if err := c.BodyParser(req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return ValidationError(c, "Invalid request format")
	}
	if err := validate.Struct(req); err != nil {
		log.Printf("Validation failed: %v", err)
		return ValidationError(c, err.Error())
	}
	return nil
}