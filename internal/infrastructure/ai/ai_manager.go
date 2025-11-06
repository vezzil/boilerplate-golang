package ai

import (
	"context"
	"fmt"
	"sync"

	"boilerplate-golang/internal/infrastructure/config"
	"boilerplate-golang/internal/infrastructure/logger"
)

var (
	aiClient *AIManager
	once     sync.Once
)

// AIManager manages AI operations and provider selection
type AIManager struct {
	config   config.AppConfig
	provider AIProvider
	enabled  bool
}

// AIProvider defines the interface that all AI providers must implement
type AIProvider interface {
	AnalyzeResume(ctx context.Context, jobDescription, resumeText string) (*AnalysisResult, error)
	GenerateInterviewQuestions(ctx context.Context, jobDescription, candidateProfile string) ([]string, error)
	ScoreCandidate(ctx context.Context, jobDescription, resumeText string) (*CandidateScore, error)
	GenerateJobSummary(ctx context.Context, jobDescription string) (string, error)
	ExtractSkills(ctx context.Context, resumeText string) ([]string, error)
}

// AnalysisResult represents the result of AI resume analysis
type AnalysisResult struct {
	SimilarityScore float64  `json:"similarity_score"`
	MatchReasons    []string `json:"match_reasons"`
	MismatchReasons []string `json:"mismatch_reasons"`
	OverallFit      string   `json:"overall_fit"`
	KeyStrengths    []string `json:"key_strengths"`
	AreasOfConcern  []string `json:"areas_of_concern"`
}

// CandidateScore represents a candidate's AI-calculated score
type CandidateScore struct {
	Score           float64  `json:"score"`
	Confidence      float64  `json:"confidence"`
	Explanation     string   `json:"explanation"`
	KeySkills       []string `json:"key_skills"`
	MissingSkills   []string `json:"missing_skills"`
	ExperienceMatch string   `json:"experience_match"`
}

// Init initializes the AI manager with OpenAI provider
func Init() error {
	var initErr error
	once.Do(func() {
		cfg := config.Get()

		if !cfg.AI.Enabled {
			logger.Info("AI features are disabled in configuration")
			aiClient = &AIManager{
				config:  cfg,
				enabled: false,
			}
			return
		}

		logger.Info("Initializing AI manager with OpenAI provider")

		// Only OpenAI is supported for now
		provider, providerErr := NewOpenAIProvider(cfg)
		if providerErr != nil {
			initErr = fmt.Errorf("failed to initialize OpenAI provider: %w", providerErr)
			return
		}

		aiClient = &AIManager{
			config:   cfg,
			provider: provider,
			enabled:  true,
		}

		logger.Info("AI manager initialized successfully")
	})

	return initErr
}

// GetAIManager returns the singleton AI manager instance
func GetAIManager() (*AIManager, error) {
	if aiClient == nil {
		return nil, fmt.Errorf("AI manager not initialized. Call Init() first")
	}
	return aiClient, nil
}

// IsEnabled returns whether AI features are enabled
func (m *AIManager) IsEnabled() bool {
	return m.enabled
}

// AnalyzeResume analyzes a resume against a job description
func (m *AIManager) AnalyzeResume(jobDescription, resumeText string) (*AnalysisResult, error) {
	if !m.enabled {
		return nil, fmt.Errorf("AI features are disabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.config.AI.Timeout)
	defer cancel()

	result, err := m.provider.AnalyzeResume(ctx, jobDescription, resumeText)
	if err != nil {
		logger.Error("Failed to analyze resume: %v", err)
		return nil, fmt.Errorf("failed to analyze resume: %w", err)
	}

	return result, nil
}

// GenerateInterviewQuestions generates tailored interview questions
func (m *AIManager) GenerateInterviewQuestions(jobDescription, candidateProfile string) ([]string, error) {
	if !m.enabled {
		return nil, fmt.Errorf("AI features are disabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.config.AI.Timeout)
	defer cancel()

	questions, err := m.provider.GenerateInterviewQuestions(ctx, jobDescription, candidateProfile)
	if err != nil {
		logger.Error("Failed to generate interview questions: %v", err)
		return nil, fmt.Errorf("failed to generate interview questions: %w", err)
	}

	return questions, nil
}

// ScoreCandidate provides an AI-powered candidate score
func (m *AIManager) ScoreCandidate(jobDescription, resumeText string) (*CandidateScore, error) {
	if !m.enabled {
		return nil, fmt.Errorf("AI features are disabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.config.AI.Timeout)
	defer cancel()

	score, err := m.provider.ScoreCandidate(ctx, jobDescription, resumeText)
	if err != nil {
		logger.Error("Failed to score candidate: %v", err)
		return nil, fmt.Errorf("failed to score candidate: %w", err)
	}

	return score, nil
}

// GenerateJobSummary creates a concise summary of a job description
func (m *AIManager) GenerateJobSummary(jobDescription string) (string, error) {
	if !m.enabled {
		return "", fmt.Errorf("AI features are disabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.config.AI.Timeout)
	defer cancel()

	summary, err := m.provider.GenerateJobSummary(ctx, jobDescription)
	if err != nil {
		logger.Error("Failed to generate job summary: %v", err)
		return "", fmt.Errorf("failed to generate job summary: %w", err)
	}

	return summary, nil
}

// ExtractSkills extracts skills from a resume using AI
func (m *AIManager) ExtractSkills(resumeText string) ([]string, error) {
	if !m.enabled {
		return nil, fmt.Errorf("AI features are disabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.config.AI.Timeout)
	defer cancel()

	skills, err := m.provider.ExtractSkills(ctx, resumeText)
	if err != nil {
		logger.Error("Failed to extract skills: %v", err)
		return nil, fmt.Errorf("failed to extract skills: %w", err)
	}

	return skills, nil
}
