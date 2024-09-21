package model

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/sickodev/floqer-backend/helper"
	"github.com/xuri/excelize/v2"
	"google.golang.org/api/option"
)

type RequestBody struct {
	Message string `json:"message"`
}

func GenerateResponse(c *fiber.Ctx) error {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": "No client made",
		})
	}

	model := client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(1)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	response, _ := helper.ReadCSV()
	responses, err:= helper.ConvertResponsesToCSV(response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":"Error reading CSV File",
		})
	}

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(responses),
		},
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	session := model.StartChat()

	resp, err := session.SendMessage(ctx, genai.Text(body.Message))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": "Not a valid resp",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Value": resp.Candidates[0].Content.Parts,
	})
}

type Response struct {
	WorkYear          string
	ExperienceLevel   string
	EmploymentType    string
	JobTitle          string
	Salary            string
	SalaryCurrency    string
	SalaryInUSD       string
	EmployeeResidence string
	RemoteRatio       string
	CompanyLocation   string
	CompanySize       string
}

// Function to read the Excel file and convert it to []byte
func excelToBytes(filePath string) ([]byte, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	// Get all rows in the first sheet
	rows, err := f.GetRows(f.GetSheetName(1))
	if err != nil {
		return nil, err
	}

	// Assuming the first row is the header
	for i, row := range rows {
		if i == 0 {
			continue // Skip header
		}
		if len(row) < 11 { // Check for enough columns
			continue
		}

		// Create a Response struct (optional)
		response := Response{
			WorkYear:          row[0],
			ExperienceLevel:   row[1],
			EmploymentType:    row[2],
			JobTitle:          row[3],
			Salary:            row[4],
			SalaryCurrency:    row[5],
			SalaryInUSD:       row[6],
			EmployeeResidence: row[7],
			RemoteRatio:       row[8],
			CompanyLocation:   row[9],
			CompanySize:       row[10],
		}

		// Append the struct data to the buffer as a string (or in any desired format)
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			response.WorkYear,
			response.ExperienceLevel,
			response.EmploymentType,
			response.JobTitle,
			response.Salary,
			response.SalaryCurrency,
			response.SalaryInUSD,
			response.EmployeeResidence,
			response.RemoteRatio,
			response.CompanyLocation,
			response.CompanySize,
		)
		buffer.WriteString(line)
	}

	return buffer.Bytes(), nil
}
