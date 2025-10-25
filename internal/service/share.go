package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"survey-system/internal/dto/request"
	"survey-system/internal/dto/response"
	"survey-system/internal/model"
	"survey-system/internal/repository"
	"survey-system/pkg/errors"
)

// ShareService defines the interface for share link business logic
type ShareService interface {
	GenerateShareLink(ctx context.Context, userID, surveyID uint, req *request.GenerateShareLinkRequest) (*response.ShareLinkResponse, error)
	ValidateAndGetSurvey(ctx context.Context, token string) (*response.SurveyWithPrefillResponse, error)
}

// shareService implements ShareService interface
type shareService struct {
	surveyRepo    repository.SurveyRepository
	questionRepo  repository.QuestionRepository
	oneLinkRepo   repository.OneLinkRepository
	encryptionSvc EncryptionService
	cache         Cache
	baseURL       string
	defaultExpiry time.Duration
	maxExpiry     time.Duration
}

// NewShareService creates a new share service instance
func NewShareService(
	surveyRepo repository.SurveyRepository,
	questionRepo repository.QuestionRepository,
	oneLinkRepo repository.OneLinkRepository,
	encryptionSvc EncryptionService,
	cache Cache,
	baseURL string,
	defaultExpiry time.Duration,
	maxExpiry time.Duration,
) ShareService {
	return &shareService{
		surveyRepo:    surveyRepo,
		questionRepo:  questionRepo,
		oneLinkRepo:   oneLinkRepo,
		encryptionSvc: encryptionSvc,
		cache:         cache,
		baseURL:       baseURL,
		defaultExpiry: defaultExpiry,
		maxExpiry:     maxExpiry,
	}
}

// GenerateShareLink generates an encrypted share link with prefill data
func (s *shareService) GenerateShareLink(ctx context.Context, userID, surveyID uint, req *request.GenerateShareLinkRequest) (*response.ShareLinkResponse, error) {
	// Find the survey and verify ownership
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "failed to find survey")
	}

	// Verify ownership
	if survey.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Get all questions for the survey to validate prefill keys
	questions, err := s.questionRepo.FindBySurveyID(surveyID)
	if err != nil {
		return nil, errors.WrapError(err, "failed to find questions")
	}

	// Validate prefill data - ensure all prefill keys match question prefill_key fields
	if req.PrefillData != nil && len(req.PrefillData) > 0 {
		validPrefillKeys := make(map[string]bool)
		for _, q := range questions {
			if q.PrefillKey != "" {
				validPrefillKeys[q.PrefillKey] = true
			}
		}

		for key := range req.PrefillData {
			if !validPrefillKeys[key] {
				return nil, errors.NewValidationError("prefill_data", fmt.Sprintf("invalid prefill key '%s' - no matching question found", key))
			}
		}
	}

	// Determine expiration time
	var expiresAt time.Time
	if req.ExpiresAt != nil {
		expiresAt = *req.ExpiresAt
		
		// Validate expiration is in the future
		if expiresAt.Before(time.Now()) {
			return nil, errors.NewValidationError("expires_at", "expiration time must be in the future")
		}
		
		// Validate expiration doesn't exceed max expiry
		maxExpiresAt := time.Now().Add(s.maxExpiry)
		if expiresAt.After(maxExpiresAt) {
			return nil, errors.NewValidationError("expires_at", fmt.Sprintf("expiration time exceeds maximum allowed duration of %v", s.maxExpiry))
		}
	} else {
		// Use default expiration
		expiresAt = time.Now().Add(s.defaultExpiry)
	}

	// Generate unique ID for this link
	uniqueID := uuid.New().String()

	// Build TokenData
	tokenData := &TokenData{
		SurveyID:    surveyID,
		PrefillData: req.PrefillData,
		ExpiresAt:   expiresAt.Unix(),
		UniqueID:    uniqueID,
	}

	// Encrypt the token
	encryptedToken, err := s.encryptionSvc.EncryptToken(tokenData)
	if err != nil {
		return nil, errors.WrapError(err, "failed to encrypt token")
	}

	// Create OneLink record in database
	oneLink := &model.OneLink{
		SurveyID:    surveyID,
		Token:       encryptedToken,
		PrefillData: req.PrefillData,
		ExpiresAt:   expiresAt,
		Used:        false,
	}

	if err := s.oneLinkRepo.Create(oneLink); err != nil {
		return nil, errors.WrapError(err, "failed to create one-time link")
	}

	// Build the complete share URL
	shareURL := fmt.Sprintf("%s/surveys/%d?token=%s", s.baseURL, surveyID, encryptedToken)

	return &response.ShareLinkResponse{
		Token:     encryptedToken,
		URL:       shareURL,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateAndGetSurvey validates a token and returns the survey with prefilled values
func (s *shareService) ValidateAndGetSurvey(ctx context.Context, token string) (*response.SurveyWithPrefillResponse, error) {
	// Step 1: Decrypt the token to get TokenData
	tokenData, err := s.encryptionSvc.DecryptToken(token)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// Step 2: Validate expiration time
	if time.Now().Unix() > tokenData.ExpiresAt {
		return nil, errors.ErrTokenExpired
	}

	// Step 3: Check Redis cache for link status first to avoid database query
	cachedUsed, err := s.cache.GetOneLinkStatus(ctx, token)
	if err != nil {
		// Log error but continue to database check
		fmt.Printf("failed to get onelink status from cache: %v\n", err)
	} else if cachedUsed {
		// Link is marked as used in cache
		return nil, errors.ErrLinkUsed
	}

	// Step 4: Find the OneLink record in database
	oneLink, err := s.oneLinkRepo.FindByToken(token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrInvalidToken
		}
		return nil, errors.WrapError(err, "failed to find one-time link")
	}

	// Step 5: Check if link has been used
	if oneLink.Used {
		// Update cache with used status
		expiresAt := time.Unix(tokenData.ExpiresAt, 0)
		cacheTTL := time.Until(expiresAt)
		if cacheTTL > 0 {
			if err := s.cache.SetOneLinkStatus(ctx, token, true, cacheTTL); err != nil {
				fmt.Printf("failed to cache onelink used status: %v\n", err)
			}
		}
		return nil, errors.ErrLinkUsed
	}

	// Step 6: Check if link has expired (double check with database record)
	if oneLink.IsExpired() {
		return nil, errors.ErrTokenExpired
	}

	// Step 7: Cache the unused status to avoid repeated database queries
	expiresAt := time.Unix(tokenData.ExpiresAt, 0)
	cacheTTL := time.Until(expiresAt)
	if cacheTTL > 0 {
		if err := s.cache.SetOneLinkStatus(ctx, token, false, cacheTTL); err != nil {
			fmt.Printf("failed to cache onelink unused status: %v\n", err)
		}
	}

	// Step 8: Mark link as accessed (first time viewing)
	if oneLink.AccessedAt == nil {
		if err := s.oneLinkRepo.MarkAsAccessed(oneLink.ID); err != nil {
			// Log error but don't fail the request
			fmt.Printf("failed to mark link as accessed: %v\n", err)
		}
	}

	// Step 9: Get the survey with questions
	survey, err := s.surveyRepo.FindByIDWithQuestions(tokenData.SurveyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, "failed to find survey")
	}

	// Step 10: Build response with prefilled values
	questionsWithPrefill := make([]response.QuestionWithPrefill, len(survey.Questions))
	for i, q := range survey.Questions {
		questionResp := response.QuestionWithPrefill{
			QuestionResponse: response.QuestionResponse{
				ID:          q.ID,
				SurveyID:    q.SurveyID,
				Type:        q.Type,
				Title:       q.Title,
				Description: q.Description,
				Required:    q.Required,
				Order:       q.Order,
				Config:      q.Config,
				PrefillKey:  q.PrefillKey,
			},
		}

		// Add prefill value if available
		if q.PrefillKey != "" && tokenData.PrefillData != nil {
			if prefillValue, exists := tokenData.PrefillData[q.PrefillKey]; exists {
				questionResp.PrefillValue = prefillValue
			}
		}

		questionsWithPrefill[i] = questionResp
	}

	return &response.SurveyWithPrefillResponse{
		ID:          survey.ID,
		Title:       survey.Title,
		Description: survey.Description,
		Questions:   questionsWithPrefill,
		PrefillData: tokenData.PrefillData,
	}, nil
}
