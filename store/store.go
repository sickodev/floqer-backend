package store

import (
	"encoding/csv"
	"os"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	WorkYear          string `json:"work_year"`
	ExperienceLevel   string `json:"experience_level"`
	EmploymentType    string `json:"employment_type"`
	JobTitle          string `json:"job_title"`
	Salary            string `json:"salary"`
	SalaryCurrency    string `json:"salary_currency"`
	SalaryInUSD       string `json:"salary_in_usd"`
	EmployeeResidence string `json:"employee_residence"`
	RemoteRatio       string `json:"remote_ratio"`
	CompanyLocation   string `json:"company_location"`
	CompanySize       string `json:"company_size"`
}

func GetData(c *fiber.Ctx) error {
	// Open the CSV file
	csvFile, err := os.Open("data.csv")
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Error reading data.csv",
		})
	}
	defer csvFile.Close()

	// Create a new CSV reader
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	// Read all data from CSV
	csvData, err := reader.ReadAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error reading data from csv.",
		})
	}

	// Check if the CSV has headers and rows
	if len(csvData) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "CSV file is empty or has no data rows.",
		})
	}

	// Initialize the slice to hold all response objects
	var responses []Response

	// Loop through each CSV row and map it to the Response struct
	for _, each := range csvData {
		if len(each) < 11 {
			// Handle rows that don't have enough columns
			continue
		}
		response := Response{
			WorkYear:          each[0],
			ExperienceLevel:   each[1],
			EmploymentType:    each[2],
			JobTitle:          each[3],
			Salary:            each[4],
			SalaryCurrency:    each[5],
			SalaryInUSD:       each[6],
			EmployeeResidence: each[7],
			RemoteRatio:       each[8],
			CompanyLocation:   each[9],
			CompanySize:       each[10],
		}

		// Append each response to the slice
		responses = append(responses, response)
	}

	// Write the responses as JSON to the client
	err = c.JSON(responses)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error sending JSON response.",
		})
	}

	return nil
}
