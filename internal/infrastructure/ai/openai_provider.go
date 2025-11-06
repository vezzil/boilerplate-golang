package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"boilerplate-golang/internal/infrastructure/config"
	"boilerplate-golang/internal/infrastructure/logger"

	"github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements AIProvider interface for OpenAI
type OpenAIProvider struct {
	client      *openai.Client
	model       string
	maxTokens   int
	temperature float32
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(cfg config.AppConfig) (*OpenAIProvider, error) {
	if cfg.AI.OpenAI.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	clientConfig := openai.DefaultConfig(cfg.AI.OpenAI.APIKey)
	if cfg.AI.OpenAI.BaseURL != "" {
		clientConfig.BaseURL = cfg.AI.OpenAI.BaseURL
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &OpenAIProvider{
		client:      client,
		model:       cfg.AI.OpenAI.Model,
		maxTokens:   cfg.AI.OpenAI.MaxTokens,
		temperature: float32(cfg.AI.OpenAI.Temperature),
	}, nil
}

// AnalyzeResume analyzes resume compatibility with job description
func (p *OpenAIProvider) AnalyzeResume(ctx context.Context, jobDescription, resumeText string) (*AnalysisResult, error) {
	prompt := fmt.Sprintf(`You are an expert ATS (Applicant Tracking System) analyzer. Analyze the compatibility between this job description and resume.

Job Description:
%s

Resume:
%s

Provide a detailed analysis in the following JSON format:
{
  "similarity_score": <number 0-100>,
  "match_reasons": [<list of reasons why this candidate is a good match>],
  "mismatch_reasons": [<list of concerns or gaps>],
  "overall_fit": "<brief summary of overall fit>",
  "key_strengths": [<candidate's key strengths for this role>],
  "areas_of_concern": [<areas where candidate may need development>]
}

Be objective and professional in your analysis.`, jobDescription, resumeText)

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       p.model,
		MaxTokens:   p.maxTokens,
		Temperature: p.temperature,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert ATS analyzer. Always respond with valid JSON.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	var result AnalysisResult
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		logger.Error("Failed to parse OpenAI response: %v. Response: %s", err, content)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &result, nil
}

// GenerateInterviewQuestions generates tailored interview questions
func (p *OpenAIProvider) GenerateInterviewQuestions(ctx context.Context, jobDescription, candidateProfile string) ([]string, error) {
	prompt := fmt.Sprintf(`Generate 5 targeted interview questions for this candidate.

Job Description:
%s

Candidate Profile:
%s

Generate specific, insightful questions that:
1. Assess technical skills mentioned in the job
2. Evaluate relevant experience
3. Test problem-solving abilities
4. Gauge cultural fit
5. Explore areas of concern

Respond with a JSON array of questions: ["question 1", "question 2", ...]`, jobDescription, candidateProfile)

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       p.model,
		MaxTokens:   p.maxTokens,
		Temperature: p.temperature,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert technical interviewer. Always respond with valid JSON array.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	var questions []string
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &questions); err != nil {
		logger.Error("Failed to parse OpenAI response: %v. Response: %s", err, content)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return questions, nil
}

// ScoreCandidate provides an AI-powered candidate score
func (p *OpenAIProvider) ScoreCandidate(ctx context.Context, jobDescription, resumeText string) (*CandidateScore, error) {
	prompt := fmt.Sprintf(`Score this candidate for the job position.

Job Description:
%s

Resume:
%s

Provide a comprehensive score in the following JSON format:
{
  "score": <number 0-100>,
  "confidence": <number 0.0-1.0>,
  "explanation": "<detailed explanation of the score>",
  "key_skills": [<skills the candidate possesses that match the job>],
  "missing_skills": [<important skills the candidate lacks>],
  "experience_match": "<assessment of experience level match>"
}

Be thorough and objective in your scoring.`, jobDescription, resumeText)

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       p.model,
		MaxTokens:   p.maxTokens,
		Temperature: p.temperature,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert recruiter. Always respond with valid JSON.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	var score CandidateScore
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &score); err != nil {
		logger.Error("Failed to parse OpenAI response: %v. Response: %s", err, content)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &score, nil
}

// GenerateJobSummary creates a concise summary of a job description
func (p *OpenAIProvider) GenerateJobSummary(ctx context.Context, jobDescription string) (string, error) {
	prompt := fmt.Sprintf(`Create a concise, professional summary of this job description in 2-3 sentences.

Job Description:
%s

Focus on the key role, main responsibilities, and required qualifications.`, jobDescription)

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       p.model,
		MaxTokens:   200,
		Temperature: p.temperature,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a professional job description summarizer.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// ExtractSkills extracts skills from a resume using AI
func (p *OpenAIProvider) ExtractSkills(ctx context.Context, resumeText string) ([]string, error) {
	prompt := fmt.Sprintf(`Extract all technical and professional skills from this resume.

Resume:
%s

Return a JSON array of skills: ["skill 1", "skill 2", ...]
Include programming languages, frameworks, tools, soft skills, and domain expertise.
Be comprehensive but avoid duplicates.`, resumeText)

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       p.model,
		MaxTokens:   500,
		Temperature: 0.3, // Lower temperature for more consistent extraction
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a professional resume parser. Always respond with valid JSON array.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	var skills []string
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &skills); err != nil {
		logger.Error("Failed to parse OpenAI response: %v. Response: %s", err, content)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return skills, nil
}
