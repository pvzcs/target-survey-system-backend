package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	"survey-system/internal/model"
	"survey-system/internal/repository"
	"survey-system/pkg/errors"

	"github.com/xuri/excelize/v2"
)

// ExportService handles data export functionality
type ExportService struct {
	surveyRepo   repository.SurveyRepository
	questionRepo repository.QuestionRepository
	responseRepo repository.ResponseRepository
}

// NewExportService creates a new ExportService
func NewExportService(
	surveyRepo repository.SurveyRepository,
	questionRepo repository.QuestionRepository,
	responseRepo repository.ResponseRepository,
) *ExportService {
	return &ExportService{
		surveyRepo:   surveyRepo,
		questionRepo: questionRepo,
		responseRepo: responseRepo,
	}
}

// ExportResponses exports survey responses in the specified format
func (s *ExportService) ExportResponses(userID, surveyID uint, format string) ([]byte, string, error) {
	// Verify survey ownership
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, "", errors.ErrNotFound
	}

	if survey.UserID != userID {
		return nil, "", errors.ErrForbidden
	}

	// Get all questions for the survey
	questions, err := s.questionRepo.FindBySurveyID(surveyID)
	if err != nil {
		return nil, "", &errors.AppError{
			Code:    "INTERNAL_ERROR",
			Message: "获取问卷题目失败",
			Status:  500,
		}
	}

	// Get all responses (no pagination for export)
	responses, _, err := s.responseRepo.FindBySurveyID(surveyID, 1, 999999)
	if err != nil {
		return nil, "", &errors.AppError{
			Code:    "INTERNAL_ERROR",
			Message: "获取填答记录失败",
			Status:  500,
		}
	}

	switch format {
	case "csv":
		return s.exportCSV(survey, questions, responses)
	case "excel":
		return s.exportExcel(survey, questions, responses)
	default:
		return nil, "", &errors.AppError{
			Code:    "INVALID_FORMAT",
			Message: "不支持的导出格式",
			Status:  400,
		}
	}
}

// exportCSV exports responses as CSV format
func (s *ExportService) exportCSV(survey *model.Survey, questions []model.Question, responses []model.Response) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Build header row
	header := s.buildCSVHeader(questions)
	if err := writer.Write(header); err != nil {
		return nil, "", &errors.AppError{
			Code:    "EXPORT_ERROR",
			Message: "生成 CSV 表头失败",
			Status:  500,
		}
	}

	// Write data rows
	for _, response := range responses {
		rows := s.buildCSVRows(questions, response)
		for _, row := range rows {
			if err := writer.Write(row); err != nil {
				return nil, "", &errors.AppError{
					Code:    "EXPORT_ERROR",
					Message: "写入 CSV 数据失败",
					Status:  500,
				}
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", &errors.AppError{
			Code:    "EXPORT_ERROR",
			Message: "生成 CSV 文件失败",
			Status:  500,
		}
	}

	filename := fmt.Sprintf("%s_responses.csv", survey.Title)
	return buf.Bytes(), filename, nil
}

// buildCSVHeader builds the CSV header row from questions
func (s *ExportService) buildCSVHeader(questions []model.Question) []string {
	header := []string{"Response ID", "Submitted At", "IP Address"}

	for _, question := range questions {
		if question.Type == model.QuestionTypeTable {
			// For table questions, add columns for each table column
			for _, col := range question.Config.Columns {
				header = append(header, fmt.Sprintf("%s - %s", question.Title, col.Label))
			}
		} else {
			header = append(header, question.Title)
		}
	}

	return header
}

// buildCSVRows builds CSV data rows from a response
// Returns multiple rows if there are table questions with multiple rows
func (s *ExportService) buildCSVRows(questions []model.Question, response model.Response) [][]string {
	// Create answer map for quick lookup
	answerMap := make(map[uint]interface{})
	for _, answer := range response.Data.Answers {
		answerMap[answer.QuestionID] = answer.Value
	}

	// Find the maximum number of rows needed (for table questions)
	maxRows := 1
	for _, question := range questions {
		if question.Type == model.QuestionTypeTable {
			if value, exists := answerMap[question.ID]; exists {
				if rows, ok := value.([]interface{}); ok {
					if len(rows) > maxRows {
						maxRows = len(rows)
					}
				}
			}
		}
	}

	// Build rows
	result := make([][]string, maxRows)
	for rowIdx := 0; rowIdx < maxRows; rowIdx++ {
		row := []string{}

		// Add response metadata only in the first row
		if rowIdx == 0 {
			row = append(row, strconv.FormatUint(uint64(response.ID), 10))
			row = append(row, response.SubmittedAt.Format("2006-01-02 15:04:05"))
			row = append(row, response.IPAddress)
		} else {
			row = append(row, "", "", "")
		}

		// Add answer values
		for _, question := range questions {
			value, exists := answerMap[question.ID]
			if !exists {
				// Add empty cells for missing answers
				if question.Type == model.QuestionTypeTable {
					for range question.Config.Columns {
						row = append(row, "")
					}
				} else {
					row = append(row, "")
				}
				continue
			}

			switch question.Type {
			case model.QuestionTypeText:
				if rowIdx == 0 {
					row = append(row, s.formatTextValue(value))
				} else {
					row = append(row, "")
				}

			case model.QuestionTypeSingle:
				if rowIdx == 0 {
					row = append(row, s.formatTextValue(value))
				} else {
					row = append(row, "")
				}

			case model.QuestionTypeMultiple:
				if rowIdx == 0 {
					row = append(row, s.formatMultipleChoiceValue(value))
				} else {
					row = append(row, "")
				}

			case model.QuestionTypeTable:
				row = append(row, s.formatTableRow(value, question.Config.Columns, rowIdx)...)
			}
		}

		result[rowIdx] = row
	}

	return result
}

// formatTextValue formats a text value for CSV
func (s *ExportService) formatTextValue(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}

// formatMultipleChoiceValue formats multiple choice values for CSV
func (s *ExportService) formatMultipleChoiceValue(value interface{}) string {
	switch v := value.(type) {
	case []interface{}:
		result := ""
		for i, item := range v {
			if i > 0 {
				result += "; "
			}
			result += fmt.Sprintf("%v", item)
		}
		return result
	case []string:
		result := ""
		for i, item := range v {
			if i > 0 {
				result += "; "
			}
			result += item
		}
		return result
	default:
		return fmt.Sprintf("%v", value)
	}
}

// formatTableRow formats a single row of table data for CSV
func (s *ExportService) formatTableRow(value interface{}, columns []model.TableColumn, rowIdx int) []string {
	rows, ok := value.([]interface{})
	if !ok {
		// Return empty cells if format is incorrect
		result := make([]string, len(columns))
		return result
	}

	// If this row index doesn't exist in the data, return empty cells
	if rowIdx >= len(rows) {
		result := make([]string, len(columns))
		return result
	}

	// Each row is an array of values
	rowData, ok := rows[rowIdx].([]interface{})
	if !ok {
		result := make([]string, len(columns))
		return result
	}

	// Extract values for each column (by index)
	result := make([]string, len(columns))
	for i := range columns {
		if i < len(rowData) {
			result[i] = fmt.Sprintf("%v", rowData[i])
		} else {
			result[i] = ""
		}
	}

	return result
}

// exportExcel exports responses as Excel format
func (s *ExportService) exportExcel(survey *model.Survey, questions []model.Question, responses []model.Response) ([]byte, string, error) {
	// Create a new Excel file
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Responses"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, "", &errors.AppError{
			Code:    "EXPORT_ERROR",
			Message: "创建 Excel 工作表失败",
			Status:  500,
		}
	}

	// Set active sheet
	f.SetActiveSheet(index)

	// Build and write header row
	header := s.buildCSVHeader(questions)
	for colIdx, headerValue := range header {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheetName, cell, headerValue)
	}

	// Apply header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E0E0E0"},
			Pattern: 1,
		},
	})
	if err == nil {
		endCol, _ := excelize.CoordinatesToCellName(len(header), 1)
		f.SetCellStyle(sheetName, "A1", endCol, headerStyle)
	}

	// Write data rows
	currentRow := 2
	for _, response := range responses {
		rows := s.buildCSVRows(questions, response)
		for _, row := range rows {
			for colIdx, cellValue := range row {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)
				f.SetCellValue(sheetName, cell, cellValue)
			}
			currentRow++
		}
	}

	// Auto-fit column widths
	for colIdx := range header {
		colName, _ := excelize.ColumnNumberToName(colIdx + 1)
		f.SetColWidth(sheetName, colName, colName, 15)
	}

	// Delete default Sheet1 if it exists and is not our sheet
	if sheetName != "Sheet1" {
		f.DeleteSheet("Sheet1")
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", &errors.AppError{
			Code:    "EXPORT_ERROR",
			Message: "生成 Excel 文件失败",
			Status:  500,
		}
	}

	filename := fmt.Sprintf("%s_responses.xlsx", survey.Title)
	return buf.Bytes(), filename, nil
}
