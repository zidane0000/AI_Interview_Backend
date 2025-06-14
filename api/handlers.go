// HTTP handler functions for each endpoint
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zidane0000/AI_Interview_Backend/ai"
	"github.com/zidane0000/AI_Interview_Backend/data"
)

// Helper: parse integer query parameter with default value
func parseIntQuery(r *http.Request, key string, defaultValue int) int {
	if str := r.URL.Query().Get(key); str != "" {
		if val, err := strconv.Atoi(str); err == nil && val >= 0 {
			return val
		}
	}
	return defaultValue
}

// Helper: write JSON response
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to encode JSON: %v", err)
	}
}

// Helper: write JSON error response
func writeJSONError(w http.ResponseWriter, status int, msg string, details ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errResp := ErrorResponseDTO{Error: msg}
	if len(details) > 0 {
		errResp.Details = details[0]
	}

	if err := json.NewEncoder(w).Encode(errResp); err != nil {
		log.Printf("failed to encode JSON: %v", err)
	}
}

// CreateInterviewHandler handles POST /interviews
func CreateInterviewHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateInterviewRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}
	if req.CandidateName == "" || len(req.Questions) == 0 {
		writeJSONError(w, http.StatusBadRequest, "Missing candidate_name or questions")
		return
	}

	// Validate language if provided
	if req.InterviewLanguage != "" && !data.ValidateLanguage(req.InterviewLanguage) {
		writeJSONError(w, http.StatusBadRequest, "Invalid language code. Supported languages: en, zh-TW")
		return
	}

	// Process language parameter with default fallback
	language := data.GetValidatedLanguage(req.InterviewLanguage)

	// Generate unique ID and create interview record
	interviewID := data.GenerateID()
	interview := &data.Interview{
		ID:            interviewID,
		CandidateName: req.CandidateName,
		Questions:     req.Questions,
		Language:      language,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// Store interview in hybrid store
	err := data.GlobalStore.CreateInterview(interview)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to create interview", err.Error())
		return
	}

	resp := InterviewResponseDTO{
		ID:            interview.ID,
		CandidateName: interview.CandidateName,
		Questions:     interview.Questions,
		Language:      interview.Language,
		CreatedAt:     interview.CreatedAt,
	}
	writeJSON(w, http.StatusCreated, resp)
}

// ListInterviewsHandler handles GET /interviews
func ListInterviewsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for pagination, filtering, and sorting
	opts := data.ListInterviewsOptions{
		Limit:  parseIntQuery(r, "limit", 10),
		Offset: parseIntQuery(r, "offset", 0),
		Page:   parseIntQuery(r, "page", 0),
	}

	// Parse filtering parameters
	if candidateName := r.URL.Query().Get("candidate_name"); candidateName != "" {
		opts.CandidateName = candidateName
	}
	if status := r.URL.Query().Get("status"); status != "" {
		opts.Status = status
	}
	if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			opts.DateFrom = parsed
		}
	}
	if dateTo := r.URL.Query().Get("date_to"); dateTo != "" {
		if parsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			opts.DateTo = parsed
		}
	}

	// Parse sorting parameters
	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		opts.SortBy = sortBy
	}
	if sortOrder := r.URL.Query().Get("sort_order"); sortOrder != "" {
		opts.SortOrder = sortOrder
	}
	// Fetch interviews from memory store with options
	result, err := data.GlobalStore.GetInterviewsWithOptions(opts)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch interviews", err.Error())
		return
	}

	// Convert to DTOs
	interviewDTOs := make([]InterviewResponseDTO, len(result.Interviews))
	for i, interview := range result.Interviews {
		interviewDTOs[i] = InterviewResponseDTO{
			ID:            interview.ID,
			CandidateName: interview.CandidateName,
			Questions:     interview.Questions,
			Language:      interview.Language,
			CreatedAt:     interview.CreatedAt,
		}
	}

	resp := ListInterviewsResponseDTO{
		Interviews: interviewDTOs,
		Total:      result.Total,
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetInterviewHandler handles GET /interviews/{id}
func GetInterviewHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSONError(w, ErrCodeBadRequest, ErrMsgMissingInterviewID)
		return
	} // Get interview from memory store
	interview, err := data.GlobalStore.GetInterview(id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Interview not found")
		return
	}

	resp := InterviewResponseDTO{
		ID:            interview.ID,
		CandidateName: interview.CandidateName,
		Questions:     interview.Questions,
		Language:      interview.Language,
		CreatedAt:     interview.CreatedAt,
	}
	writeJSON(w, http.StatusOK, resp)
}

// SubmitEvaluationHandler handles POST /evaluation
func SubmitEvaluationHandler(w http.ResponseWriter, r *http.Request) {
	var req SubmitEvaluationRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}
	if req.InterviewID == "" || len(req.Answers) == 0 {
		writeJSONError(w, http.StatusBadRequest, "Missing interview_id or answers")
		return
	}
	// Validate interview exists before creating evaluation
	interview, err := data.GlobalStore.GetInterview(req.InterviewID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Interview not found")
		return
	}

	// Convert answers map to arrays for AI evaluation
	questions := interview.Questions
	answers := make([]string, len(questions))

	// Map answers from the request to the questions order
	for i := range questions {
		answerKey := fmt.Sprintf("question_%d", i)
		if answer, exists := req.Answers[answerKey]; exists {
			answers[i] = answer
		} else {
			answers[i] = "" // Empty answer if not provided
		}
	}

	// Generate AI evaluation using the same method as chat evaluation
	jobTitle := "Software Engineer" // Default job title
	jobDesc := fmt.Sprintf("Interview for %s position", interview.CandidateName)
	language := interview.Language // Use interview language for evaluation

	score, feedback, err := ai.Client.EvaluateAnswersWithContext(questions, answers, jobTitle, jobDesc, language)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to generate evaluation")
		return
	}

	// Create evaluation record
	evaluationID := data.GenerateID()
	evaluation := &data.Evaluation{
		ID:          evaluationID,
		InterviewID: req.InterviewID,
		Answers:     req.Answers,
		Score:       score,
		Feedback:    feedback,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = data.GlobalStore.CreateEvaluation(evaluation)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to save evaluation")
		return
	}

	resp := EvaluationResponseDTO{
		ID:          evaluationID,
		InterviewID: req.InterviewID,
		Answers:     req.Answers,
		Score:       score,
		Feedback:    feedback,
		CreatedAt:   evaluation.CreatedAt,
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetEvaluationHandler handles GET /evaluation/{id}
func GetEvaluationHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSONError(w, ErrCodeBadRequest, ErrMsgMissingEvaluationID)
		return
	}
	// Get evaluation from database
	evaluation, err := data.GlobalStore.GetEvaluation(id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Evaluation not found")
		return
	}

	resp := EvaluationResponseDTO{
		ID:          evaluation.ID,
		InterviewID: evaluation.InterviewID,
		Answers:     evaluation.Answers,
		Score:       evaluation.Score,
		Feedback:    evaluation.Feedback,
		CreatedAt:   evaluation.CreatedAt,
	}
	writeJSON(w, http.StatusOK, resp)
}

// StartChatSessionHandler handles POST /interviews/{id}/chat/start
func StartChatSessionHandler(w http.ResponseWriter, r *http.Request) {
	interviewID := chi.URLParam(r, "id")
	if interviewID == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing interview ID")
		return
	}

	// Validate interview exists and get it for language inheritance
	interview, err := data.GlobalStore.GetInterview(interviewID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Interview not found")
		return
	}

	// Parse optional request body for language preference
	var req StartChatSessionRequestDTO
	if r.ContentLength > 0 {
		// Ignore decode errors for optional body - use interview language as fallback
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	// Determine language: use request language if provided, otherwise inherit from interview
	sessionLanguage := interview.Language // Default to interview language
	if req.InterviewLanguage != "" {
		sessionLanguage = data.GetValidatedLanguage(req.InterviewLanguage)
	}

	// Create chat session
	sessionID := data.GenerateID()
	session := &data.ChatSession{
		ID:          sessionID,
		InterviewID: interviewID,
		Language:    sessionLanguage,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = data.GlobalStore.CreateChatSession(session)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to create chat session")
		return
	}

	// Generate initial AI greeting message
	aiResponse, err := ai.Client.GenerateChatResponseWithLanguage(sessionID, []map[string]string{}, "", sessionLanguage)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to generate AI response")
		return
	}

	// Create initial AI message
	messageID := data.GenerateID()
	aiMessage := &data.ChatMessage{
		ID:        messageID,
		SessionID: sessionID,
		Type:      "ai",
		Content:   aiResponse, Timestamp: time.Now(), CreatedAt: time.Now(),
	}

	err = data.GlobalStore.AddChatMessage(sessionID, aiMessage)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to save AI message")
		return
	}

	// Convert to DTO format
	messages, _ := data.GlobalStore.GetChatMessages(sessionID)
	messageDTOs := make([]ChatMessageDTO, len(messages))
	for i, msg := range messages {
		messageDTOs[i] = ChatMessageDTO{
			ID:        msg.ID,
			Type:      msg.Type,
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
		}
	}

	response := ChatInterviewSessionDTO{
		ID:          session.ID,
		InterviewID: session.InterviewID,
		Language:    session.Language,
		Messages:    messageDTOs,
		Status:      session.Status,
		CreatedAt:   session.CreatedAt,
	}

	writeJSON(w, http.StatusCreated, response)
}

// SendMessageHandler handles POST /chat/{sessionId}/message
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing session ID")
		return
	}

	// Parse request body
	var req SendMessageRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	if req.Message == "" {
		writeJSONError(w, http.StatusBadRequest, "Message cannot be empty")
		return
	}
	// Validate chat session exists and is active
	session, err := data.GlobalStore.GetChatSession(sessionID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Chat session not found")
		return
	}

	if session.Status != "active" {
		writeJSONError(w, http.StatusBadRequest, "Chat session is not active")
		return
	}

	// Create user message
	userMessageID := data.GenerateID()
	userMessage := &data.ChatMessage{
		ID:        userMessageID,
		SessionID: sessionID,
		Type:      "user", Content: req.Message,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}
	err = data.GlobalStore.AddChatMessage(sessionID, userMessage)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to save user message")
		return
	} // Get conversation history for AI context (excluding the current message)
	messages, err := data.GlobalStore.GetChatMessages(sessionID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get chat history")
		return
	}

	// Check if interview should end BEFORE generating AI response
	userMessageCount := 0
	for _, msg := range messages {
		if msg.Type == "user" {
			userMessageCount++
		}
	}

	shouldEndInterview := ai.Client.ShouldEndInterview(userMessageCount)

	// Build structured conversation history excluding the current user message
	conversationHistory := make([]map[string]string, 0)
	for _, msg := range messages {
		// Skip the current user message we just added
		if msg.ID != userMessage.ID {
			conversationHistory = append(conversationHistory, map[string]string{
				"role":    msg.Type,
				"content": msg.Content,
			})
		}
	} // Generate AI response - use closing context if interview should end
	var aiResponse string
	if shouldEndInterview {
		aiResponse, err = ai.Client.GenerateClosingMessageWithLanguage(sessionID, conversationHistory, req.Message, session.Language)
	} else {
		aiResponse, err = ai.Client.GenerateChatResponseWithLanguage(sessionID, conversationHistory, req.Message, session.Language)
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to generate AI response")
		return
	}

	// Create AI message
	aiMessageID := data.GenerateID()
	aiMessage := &data.ChatMessage{
		ID:        aiMessageID,
		SessionID: sessionID,
		Type:      "ai",
		Content:   aiResponse, Timestamp: time.Now(),
		CreatedAt: time.Now()}

	err = data.GlobalStore.AddChatMessage(sessionID, aiMessage)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to save AI message")
		return
	}

	// Update session status if interview should end
	if shouldEndInterview {
		session.Status = "completed"
		session.UpdatedAt = time.Now()
		endedAt := time.Now()
		session.EndedAt = &endedAt
		if err := data.GlobalStore.UpdateChatSession(session); err != nil {
			log.Printf("Failed to update chat session: %v", err)
		}
	}

	// Convert to DTO format
	userMessageDTO := ChatMessageDTO{
		ID:        userMessage.ID,
		Type:      userMessage.Type,
		Content:   userMessage.Content,
		Timestamp: userMessage.Timestamp,
	}

	aiMessageDTO := ChatMessageDTO{
		ID:        aiMessage.ID,
		Type:      aiMessage.Type,
		Content:   aiMessage.Content,
		Timestamp: aiMessage.Timestamp,
	}
	response := SendMessageResponseDTO{
		Message:       userMessageDTO,
		AIResponse:    &aiMessageDTO,
		SessionStatus: session.Status,
	}

	writeJSON(w, http.StatusOK, response)
}

// GetChatSessionHandler handles GET /chat/{sessionId}
func GetChatSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing session ID")
		return
	}
	// Get chat session
	session, err := data.GlobalStore.GetChatSession(sessionID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Chat session not found")
		return
	}

	// Get all messages for the session
	messages, err := data.GlobalStore.GetChatMessages(sessionID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get chat messages")
		return
	}

	// Convert to DTO format
	messageDTOs := make([]ChatMessageDTO, len(messages))
	for i, msg := range messages {
		messageDTOs[i] = ChatMessageDTO{
			ID:        msg.ID,
			Type:      msg.Type,
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
		}
	}

	response := ChatInterviewSessionDTO{
		ID:          session.ID,
		InterviewID: session.InterviewID,
		Language:    session.Language,
		Messages:    messageDTOs,
		Status:      session.Status,
		CreatedAt:   session.CreatedAt,
	}

	writeJSON(w, http.StatusOK, response)
}

// EndChatSessionHandler handles POST /chat/{sessionId}/end
func EndChatSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing session ID")
		return
	}

	// Get chat session
	session, err := data.GlobalStore.GetChatSession(sessionID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Chat session not found")
		return
	}

	// Mark session as completed
	session.Status = "completed"
	session.UpdatedAt = time.Now()
	endedAt := time.Now()
	session.EndedAt = &endedAt

	err = data.GlobalStore.UpdateChatSession(session)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update session")
		return
	}

	// Get all messages for evaluation
	messages, err := data.GlobalStore.GetChatMessages(sessionID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get chat messages")
		return
	}

	// Get interview details for context
	interview, err := data.GlobalStore.GetInterview(session.InterviewID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get interview details")
		return
	}

	// Convert chat messages to evaluation format
	answers := make(map[string]string)
	questions := make([]string, 0)
	userAnswers := make([]string, 0)

	for _, msg := range messages {
		if msg.Type == "ai" {
			questions = append(questions, msg.Content)
		} else if msg.Type == "user" {
			userAnswers = append(userAnswers, msg.Content)
			// Map answers to question indices
			questionIndex := len(userAnswers) - 1
			answers[fmt.Sprintf("question_%d", questionIndex)] = msg.Content
		}
	}

	// Generate evaluation using AI service with interview context
	jobTitle := "Software Engineer" // Default job title
	jobDesc := fmt.Sprintf("Interview for %s position", interview.CandidateName)
	language := session.Language // Use session language for evaluation

	score, feedback, err := ai.Client.EvaluateAnswersWithContext(questions, userAnswers, jobTitle, jobDesc, language)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to generate evaluation")
		return
	}

	// Create evaluation record
	evaluationID := data.GenerateID()
	evaluation := &data.Evaluation{
		ID:          evaluationID,
		InterviewID: session.InterviewID, Answers: answers,
		Score:     score,
		Feedback:  feedback,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = data.GlobalStore.CreateEvaluation(evaluation)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to save evaluation")
		return
	}

	// Convert to DTO format
	response := EvaluationResponseDTO{
		ID:          evaluation.ID,
		InterviewID: evaluation.InterviewID,
		Answers:     evaluation.Answers,
		Score:       evaluation.Score,
		Feedback:    evaluation.Feedback,
		CreatedAt:   evaluation.CreatedAt,
	}

	writeJSON(w, http.StatusOK, response)
}
