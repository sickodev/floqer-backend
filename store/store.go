package store

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sickodev/floqer-backend/helper"
)

func GetData(c *fiber.Ctx) error {

	responses, err := helper.ReadCSV()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": "Error reading file",
		})
	}

	c.Set("Content-Type", "application/json")

	// Send the JSON data as the response
	return c.JSON(responses)
}
