package ai

// Example usage of AI Manager in your application
// This file shows how to integrate AI features into your services

/*
// Example 1: Analyze Resume in Application Service
func (s *ApplicationService) CreateApplicationWithAI(candidateID, jobID string) error {
    // Get candidate and job from database
    candidate := // ... fetch candidate
    job := // ... fetch job

    // Get AI manager
    aiManager, err := ai.GetAIManager()
    if err != nil || !aiManager.IsEnabled() {
        // AI not available, continue without AI analysis
        return s.CreateApplication(candidateID, jobID)
    }

    // Analyze resume compatibility
    analysis, err := aiManager.AnalyzeResume(job.Description, candidate.ResumeText)
    if err != nil {
        log.Printf("AI analysis failed: %v", err)
        // Continue without AI analysis
    } else {
        // Use AI analysis results
        application.AIScore = analysis.SimilarityScore
        // Store analysis results in database
    }

    return nil
}

// Example 2: Generate Interview Questions
func (s *JobService) GenerateInterviewQuestions(jobID, candidateID string) ([]string, error) {
    aiManager, err := ai.GetAIManager()
    if err != nil {
        return nil, fmt.Errorf("AI not available: %w", err)
    }

    job := // ... fetch job
    candidate := // ... fetch candidate

    questions, err := aiManager.GenerateInterviewQuestions(
        job.Description,
        candidate.ResumeText,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to generate questions: %w", err)
    }

    return questions, nil
}

// Example 3: Auto-score Candidates
func (s *ApplicationService) ScoreAllCandidates(jobID string) error {
    aiManager, err := ai.GetAIManager()
    if err != nil {
        return fmt.Errorf("AI not available: %w", err)
    }

    job := // ... fetch job
    applications := // ... fetch all applications for this job

    for _, app := range applications {
        candidate := // ... fetch candidate

        score, err := aiManager.ScoreCandidate(job.Description, candidate.ResumeText)
        if err != nil {
            log.Printf("Failed to score candidate %s: %v", candidate.ID, err)
            continue
        }

        // Update application with AI score
        app.AIScore = score.Score
        // Save to database
    }

    return nil
}

// Example 4: Extract Skills from Resume
func (s *CandidateService) EnrichCandidateProfile(candidateID string) error {
    aiManager, err := ai.GetAIManager()
    if err != nil {
        return fmt.Errorf("AI not available: %w", err)
    }

    candidate := // ... fetch candidate

    skills, err := aiManager.ExtractSkills(candidate.ResumeText)
    if err != nil {
        return fmt.Errorf("failed to extract skills: %w", err)
    }

    // Update candidate skills
    candidate.Skills = skills
    // Save to database

    return nil
}

// Example 5: Batch Processing with Error Handling
func (s *AIService) BatchAnalyzeCandidates(jobID string) error {
    aiManager, err := ai.GetAIManager()
    if err != nil {
        return fmt.Errorf("AI not available: %w", err)
    }

    job := // ... fetch job
    applications := // ... fetch applications

    successCount := 0
    failureCount := 0

    for _, app := range applications {
        candidate := // ... fetch candidate

        result, err := aiManager.AnalyzeResume(job.Description, candidate.ResumeText)
        if err != nil {
            log.Printf("Failed to analyze candidate %s: %v", candidate.ID, err)
            failureCount++
            continue
        }

        // Store results
        app.AIScore = result.SimilarityScore
        // ... save to database

        successCount++

        // Rate limiting: sleep between requests if needed
        time.Sleep(100 * time.Millisecond)
    }

    log.Printf("Batch analysis complete: %d succeeded, %d failed", successCount, failureCount)
    return nil
}
*/
