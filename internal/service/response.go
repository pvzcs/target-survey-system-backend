package service

import (
	"context"
	"fmt"
	"time"

	"survey-system/internal/cache"
	"survey-system/internal/dto/request"
	"survey-system/internal/dto/response"
	"survey-system/internal/model"
	"survey-system/internal/repository"
	"survey-system/pkg/errors"
)

// ResponseService handles response-related business logic
type ResponseService struct {
	responseRepo  repository.ResponseRepository
	surveyRepo    repository.SurveyRepository
	questionRepo  repository.QuestionRepository
	oneLinkRepo   repository.OneLinkRepository
	encryptionSvc EncryptionService
	cache         cache.Cache
	exportSvc     *ExportService
}

// NewResponseService creates a new ResponseService
func NewResponseService(
	responseRepo repository.ResponseRepository,
	surveyRepo repository.SurveyRepository,
	questionRepo repository.QuestionRepository,
	oneLinkRepo repository.OneLinkRepository,
	encryptionSvc EncryptionService,
	cache cache.Cache,
	exportSvc *ExportService,
) *ResponseService {
	return &ResponseService{
		responseRepo:  responseRepo,
		surveyRepo:    surveyRepo,
		questionRepo:  questionRepo,
		oneLinkRepo:   oneLinkRepo,
		encryptionSvc: encryptionSvc,
		cache:         cache,
		exportSvc:     exportSvc,
	}
}

// validateResponseData validates the response data against question configurations
func (s *ResponseService) validateResponseData(questions []model.Question, answers []request.AnswerRequest) error {
	// Create a map of question ID to question for easy lookup
	questionMap := make(map[uint]*model.Question)
	for i := range questions {
		questionMap[questions[i].ID] = &questions[i]
	}
	
	// Create a map of answered question IDs
	answeredQuestions := make(map[uint]bool)
	for _, answer := range answers {
		answeredQuestions[answer.QuestionID] = true
	}
	
	// Check all required questions are answered
	for _, question := range questions {
		if question.Required && !answeredQuestions[question.ID] {
			return &errors.AppError{
				Code:    "VALIDATION_FAILED",
				Message: fmt.Sprintf("必填题目 '%s' 未回答", question.Title),
				Status:  400,
			}
		}
	}
	
	// Validate each answer
	for _, answer := range answers {
		question, exists := questionMap[answer.QuestionID]
		if !exists {
			return &errors.AppError{
				Code:    "VALIDATION_FAILED",
				Message: fmt.Sprintf("题目 ID %d 不存在", answer.QuestionID),
				Status:  400,
			}
		}
		
		if err := s.validateAnswer(question, answer.Value); err != nil {
			return err
		}
	}
	
	return nil
}

// validateAnswer validates a single answer based on question type and configuration
func (s *ResponseService) validateAnswer(question *model.Question, value interface{}) error {
	switch question.Type {
	case model.QuestionTypeText:
		return s.validateTextAnswer(question, value)
	case model.QuestionTypeSingle:
		return s.validateSingleChoiceAnswer(question, value)
	case model.QuestionTypeMultiple:
		return s.validateMultipleChoiceAnswer(question, value)
	case model.QuestionTypeTable:
		return s.validateTableAnswer(question, value)
	default:
		return &errors.AppError{
			Code:    "VALIDATION_FAILED",
			Message: fmt.Sprintf("不支持的题目类型: %s", question.Type),
			Status:  400,
		}
	}
}

// validateTextAnswer validates text question answer
func (s *ResponseService) validateTextAnswer(question *model.Question, value interface{}) error {
	_, ok := value.(string)
	if !ok {
		return &errors.AppError{
			Code:    "VALIDATION_FAILED",
			Message: fmt.Sprintf("题目 '%s' 的答案必须是字符串", question.Title),
			Status:  400,
		}
	}
	return nil
}

// validateSingleChoiceAnswer validates single choice question answer
func (s *ResponseService) validateSingleChoiceAnswer(question *model.Question, value interface{}) error {
	answer, ok := value.(string)
	if !ok {
		return &errors.AppError{
			Code:    "VALIDATION_FAILED",
			Message: fmt.Sprintf("题目 '%s' 的答案必须是字符串", question.Title),
			Status:  400,
		}
	}
	
	// Check if the answer is in the options
	validOption := false
	for _, option := range question.Config.Options {
		if option == answer {
			validOption = true
			break
		}
	}
	
	if !validOption {
		return &errors.AppError{
			Code:    "VALIDATION_FAILED",
			Message: fmt.Sprintf("题目 '%s' 的答案 '%s' 不在选项中", question.Title, answer),
			Status:  400,
		}
	}
	
	return nil
}

// validateMultipleChoiceAnswer validates multiple choice question answer
func (s *ResponseService) validateMultipleChoiceAnswer(question *model.Question, value interface{}) error {
	// Value can be []interface{} or []string
	var answers []string
	
	switch v := value.(type) {
	case []interface{}:
		answers = make([]string, len(v))
		for i, item := range v {
			str, ok := item.(string)
			if !ok {
				return &errors.AppError{
					Code:    "VALIDATION_FAILED",
					Message: fmt.Sprintf("题目 '%s' 的答案必须是字符串数组", question.Title),
					Status:  400,
				}
			}
			answers[i] = str
		}
	case []string:
		answers = v
	default:
		return &errors.AppError{
			Code:    "VALIDATION_FAILED",
			Message: fmt.Sprintf("题目 '%s' 的答案必须是字符串数组", question.Title),
			Status:  400,
		}
	}
	
	// Check if all answers are in the options
	optionMap := make(map[string]bool)
	for _, option := range question.Config.Options {
		optionMap[option] = true
	}
	
	for _, answer := range answers {
		if !optionMap[answer] {
			return &errors.AppError{
				Code:    "VALIDATION_FAILED",
				Message: fmt.Sprintf("题目 '%s' 的答案 '%s' 不在选项中", question.Title, answer),
				Status:  400,
			}
		}
	}
	
	return nil
}

// validateTableAnswer validates table question answer
func (s *ResponseService) validateTableAnswer(question *model.Question, value interface{}) error {
	// Value should be []interface{} where each item is a map[string]interface{}
	rows, ok := value.([]interface{})
	if !ok {
		return &errors.AppError{
			Code:    "VALIDATION_FAILED",
			Message: fmt.Sprintf("题目 '%s' 的答案必须是数组", question.Title),
			Status:  400,
		}
	}
	
	// Check row count constraints
	rowCount := len(rows)
	if question.Config.MinRows > 0 && rowCount < question.Config.MinRows {
		return &errors.AppError{
			Code:    "VALIDATION_FAILED",
			Message: fmt.Sprintf("题目 '%s' 至少需要 %d 行，当前只有 %d 行", question.Title, question.Config.MinRows, rowCount),
			Status:  400,
		}
	}
	if question.Config.MaxRows > 0 && rowCount > question.Config.MaxRows {
		return &errors.AppError{
			Code:    "VALIDATION_FAILED",
			Message: fmt.Sprintf("题目 '%s' 最多允许 %d 行，当前有 %d 行", question.Title, question.Config.MaxRows, rowCount),
			Status:  400,
		}
	}
	
	// Create column map for validation
	columnMap := make(map[string]*model.TableColumn)
	for i := range question.Config.Columns {
		columnMap[question.Config.Columns[i].ID] = &question.Config.Columns[i]
	}
	
	// Validate each row
	for rowIdx, rowInterface := range rows {
		row, ok := rowInterface.(map[string]interface{})
		if !ok {
			return &errors.AppError{
				Code:    "VALIDATION_FAILED",
				Message: fmt.Sprintf("题目 '%s' 第 %d 行格式错误", question.Title, rowIdx+1),
				Status:  400,
			}
		}
		
		// Validate each cell in the row
		for colID, cellValue := range row {
			column, exists := columnMap[colID]
			if !exists {
				return &errors.AppError{
					Code:    "VALIDATION_FAILED",
					Message: fmt.Sprintf("题目 '%s' 列 '%s' 不存在", question.Title, colID),
					Status:  400,
				}
			}
			
			if err := s.validateTableCell(question.Title, rowIdx+1, column, cellValue); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// validateTableCell validates a single cell in a table question
func (s *ResponseService) validateTableCell(questionTitle string, rowNum int, column *model.TableColumn, value interface{}) error {
	switch column.Type {
	case "text":
		if _, ok := value.(string); !ok {
			return &errors.AppError{
				Code:    "VALIDATION_FAILED",
				Message: fmt.Sprintf("题目 '%s' 第 %d 行列 '%s' 必须是文本", questionTitle, rowNum, column.Label),
				Status:  400,
			}
		}
	case "number":
		// JSON unmarshaling can give us float64 for numbers
		switch value.(type) {
		case float64, int, int64:
			// Valid number types
		default:
			return &errors.AppError{
				Code:    "VALIDATION_FAILED",
				Message: fmt.Sprintf("题目 '%s' 第 %d 行列 '%s' 必须是数字", questionTitle, rowNum, column.Label),
				Status:  400,
			}
		}
	case "select":
		strValue, ok := value.(string)
		if !ok {
			return &errors.AppError{
				Code:    "VALIDATION_FAILED",
				Message: fmt.Sprintf("题目 '%s' 第 %d 行列 '%s' 必须是字符串", questionTitle, rowNum, column.Label),
				Status:  400,
			}
		}
		
		// Check if value is in options
		validOption := false
		for _, option := range column.Options {
			if option == strValue {
				validOption = true
				break
			}
		}
		
		if !validOption {
			return &errors.AppError{
				Code:    "VALIDATION_FAILED",
				Message: fmt.Sprintf("题目 '%s' 第 %d 行列 '%s' 的值 '%s' 不在选项中", questionTitle, rowNum, column.Label, strValue),
				Status:  400,
			}
		}
	}
	
	return nil
}

// SubmitResponse handles the submission of a survey response
func (s *ResponseService) SubmitResponse(req *request.SubmitResponseRequest, ipAddress, userAgent string) (*response.SubmitResponseResponse, error) {
	ctx := context.Background()
	
	// Decrypt and validate token
	tokenData, err := s.encryptionSvc.DecryptToken(req.Token)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}
	
	// Check if token is expired
	if time.Now().Unix() > tokenData.ExpiresAt {
		return nil, errors.ErrTokenExpired
	}
	
	// Check one-time link status in cache first
	used, err := s.cache.GetOneLinkStatus(ctx, req.Token)
	if err == nil && used {
		return nil, errors.ErrLinkUsed
	}
	
	// Acquire distributed lock to prevent concurrent submissions
	lockKey := fmt.Sprintf("response:%s", req.Token)
	acquired, err := s.cache.AcquireLock(ctx, lockKey, 10*time.Second)
	if err != nil || !acquired {
		return nil, &errors.AppError{
			Code:    "CONCURRENT_SUBMISSION",
			Message: "请勿重复提交",
			Status:  409,
		}
	}
	defer s.cache.ReleaseLock(ctx, lockKey)
	
	// Verify one-time link in database
	oneLink, err := s.oneLinkRepo.FindByToken(req.Token)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}
	
	if oneLink.Used {
		// Update cache
		s.cache.SetOneLinkStatus(ctx, req.Token, true, time.Until(time.Unix(tokenData.ExpiresAt, 0)))
		return nil, errors.ErrLinkUsed
	}
	
	// Get survey with questions
	survey, err := s.surveyRepo.FindByID(tokenData.SurveyID)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	
	// Check if survey is published
	if survey.Status != "published" {
		return nil, errors.ErrSurveyNotPublished
	}
	
	// Get all questions for the survey
	questions, err := s.questionRepo.FindBySurveyID(survey.ID)
	if err != nil {
		return nil, &errors.AppError{
			Code:    "INTERNAL_ERROR",
			Message: "获取问卷题目失败",
			Status:  500,
		}
	}
	
	// Validate response data
	if err := s.validateResponseData(questions, req.Answers); err != nil {
		return nil, err
	}
	
	// Convert request answers to model answers
	answers := make([]model.Answer, len(req.Answers))
	for i, ans := range req.Answers {
		answers[i] = model.Answer{
			QuestionID: ans.QuestionID,
			Value:      ans.Value,
		}
	}
	
	// Create response record
	responseModel := &model.Response{
		SurveyID:    survey.ID,
		OneLinkID:   oneLink.ID,
		Data: model.ResponseData{
			Answers: answers,
		},
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		SubmittedAt: time.Now(),
	}
	
	if err := s.responseRepo.Create(responseModel); err != nil {
		return nil, &errors.AppError{
			Code:    "INTERNAL_ERROR",
			Message: "保存填答记录失败",
			Status:  500,
		}
	}
	
	// Mark one-time link as used
	if err := s.oneLinkRepo.MarkAsUsed(oneLink.ID); err != nil {
		// Log error but don't fail the request since response is already saved
		// In production, this should be logged properly
	}
	
	// Update cache
	s.cache.SetOneLinkStatus(ctx, req.Token, true, time.Until(time.Unix(tokenData.ExpiresAt, 0)))
	
	return &response.SubmitResponseResponse{
		ID:          responseModel.ID,
		SurveyID:    responseModel.SurveyID,
		SubmittedAt: responseModel.SubmittedAt,
		Message:     "提交成功",
	}, nil
}

// GetResponses retrieves paginated responses for a survey
func (s *ResponseService) GetResponses(userID, surveyID uint, page, pageSize int) (*response.PaginatedResponseResponse, error) {
	// Verify survey ownership
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	
	if survey.UserID != userID {
		return nil, errors.ErrForbidden
	}
	
	// Get responses with pagination
	responses, total, err := s.responseRepo.FindBySurveyID(surveyID, page, pageSize)
	if err != nil {
		return nil, &errors.AppError{
			Code:    "INTERNAL_ERROR",
			Message: "获取填答记录失败",
			Status:  500,
		}
	}
	
	// Convert to response DTOs
	responseList := make([]response.ResponseListItem, len(responses))
	for i, resp := range responses {
		responseList[i] = response.ResponseListItem{
			ID:          resp.ID,
			SurveyID:    resp.SurveyID,
			IPAddress:   resp.IPAddress,
			SubmittedAt: resp.SubmittedAt,
		}
	}
	
	return &response.PaginatedResponseResponse{
		Responses: responseList,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}

// GetStatistics retrieves statistics for a survey
func (s *ResponseService) GetStatistics(userID, surveyID uint) (*response.StatisticsResponse, error) {
	// Verify survey ownership
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	
	if survey.UserID != userID {
		return nil, errors.ErrForbidden
	}
	
	// Count total responses
	count, err := s.responseRepo.CountBySurveyID(surveyID)
	if err != nil {
		return nil, &errors.AppError{
			Code:    "INTERNAL_ERROR",
			Message: "获取统计信息失败",
			Status:  500,
		}
	}
	
	// Calculate completion rate (assuming all submitted responses are complete)
	completionRate := 100.0
	if count == 0 {
		completionRate = 0.0
	}
	
	return &response.StatisticsResponse{
		SurveyID:       surveyID,
		TotalResponses: count,
		CompletionRate: completionRate,
	}, nil
}

// ExportResponses exports survey responses in the specified format
func (s *ResponseService) ExportResponses(userID, surveyID uint, format string) ([]byte, string, error) {
	return s.exportSvc.ExportResponses(userID, surveyID, format)
}
