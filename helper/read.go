package helper

import (
	"bytes"
	"encoding/csv"
	"os"
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

func ReadCSV() ([]Response, error){
	csvFile, err := os.Open("data.csv")
	if err != nil {
		return []Response{}, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		return []Response{}, err
	}

	var responses []Response

	// Skip header row
	if len(csvData) > 0 {
		headers := csvData[0]
		for _, row := range csvData[1:] {
			if len(row) < len(headers) {
				continue // Skip rows with fewer columns than headers
			}

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

			responses = append(responses, response)
		}
	}

	return responses, err
}

func ConvertResponsesToCSV(responses []Response) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{
		"work_year",
		"experience_level",
		"employment_type",
		"job_title",
		"salary",
		"salary_currency",
		"salary_in_usd",
		"employee_residence",
		"remote_ratio",
		"company_location",
		"company_size",
	}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	// Write each response as a row
	for _, response := range responses {
		row := []string{
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
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	// Flush the writer to ensure all data is written
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}