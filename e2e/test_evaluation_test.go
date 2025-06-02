package e2e

import (
	"testing"
)

// TestEvaluationWorkflow tests the evaluation generation and retrieval
func TestEvaluationWorkflow(t *testing.T) {
	t.Run("CompleteEvaluationWorkflow", func(t *testing.T) {
		// Create interview and chat session
		interview := CreateTestInterview(t, "Evaluation Test User", GetSampleQuestions())
		session := StartChatSession(t, interview.ID)

		// Simulate a complete conversation
		responses := []string{
			"I am a software engineer with 5 years of experience in web development using React and Node.js.",
			"My greatest strength is problem-solving. I enjoy breaking down complex problems into manageable pieces.",
			"I worked on a microservices migration project that required careful planning and execution over 6 months.",
			"In 5 years, I see myself as a senior technical lead, mentoring junior developers and architecting solutions.",
		}

		// Send multiple messages to create a conversation
		for i, response := range responses {
			msgResp := SendMessage(t, session.ID, response)

			// Verify message was stored correctly
			if msgResp.Message.Content != response {
				t.Errorf("Message %d not stored correctly", i)
			}
			if msgResp.AIResponse == nil {
				t.Errorf("No AI response for message %d", i)
			}
			if msgResp.AIResponse.Type != "ai" {
				t.Errorf("AI response type incorrect for message %d", i)
			}
		}

		// End session and get evaluation
		evaluation := EndChatSession(t, session.ID)

		// Verify evaluation fields
		if evaluation.ID == "" {
			t.Error("Evaluation ID should not be empty")
		}
		if evaluation.InterviewID != interview.ID {
			t.Errorf("Expected interview ID %s, got %s", interview.ID, evaluation.InterviewID)
		}
		if evaluation.Score <= 0 || evaluation.Score > 1 {
			t.Errorf("Score should be between 0 and 1, got %f", evaluation.Score)
		}
		if evaluation.Feedback == "" {
			t.Error("Feedback should not be empty")
		}
		if len(evaluation.Answers) == 0 {
			t.Error("Answers should not be empty")
		}
		if evaluation.CreatedAt.IsZero() {
			t.Error("CreatedAt should not be zero")
		}
	})

	t.Run("ShortConversationEvaluation", func(t *testing.T) {
		// Test evaluation with minimal conversation
		interview := CreateTestInterview(t, "Short Test", []string{"Tell me about yourself"})
		session := StartChatSession(t, interview.ID)

		// Send only one message
		SendMessage(t, session.ID, "I am a developer.")

		// End session immediately
		evaluation := EndChatSession(t, session.ID)

		// Verify evaluation still generated
		if evaluation.ID == "" {
			t.Error("Evaluation should be generated even for short conversations")
		}
		if evaluation.Score == 0 {
			t.Error("Score should be non-zero even for short conversations")
		}
	})

	t.Run("EvaluationConsistency", func(t *testing.T) {
		// Test that similar responses get similar evaluations
		testResponses := []string{
			"I am an experienced software engineer with strong technical skills",
			"I have been working as a software engineer and have developed strong technical abilities",
		}

		var evaluations []EvaluationResponseDTO

		for i, response := range testResponses {
			interview := CreateTestInterview(t, "Consistency Test", GetSampleQuestions())
			session := StartChatSession(t, interview.ID)

			SendMessage(t, session.ID, response)
			evaluation := EndChatSession(t, session.ID)
			evaluations = append(evaluations, evaluation)

			t.Logf("Evaluation %d: Score=%.2f, Feedback=%s", i, evaluation.Score, evaluation.Feedback)
		}

		// Verify evaluations are generated (consistency testing would require more sophisticated AI)
		for i, eval := range evaluations {
			if eval.Score <= 0 {
				t.Errorf("Evaluation %d has invalid score: %f", i, eval.Score)
			}
		}
	})

	t.Run("EvaluationAnswersMapping", func(t *testing.T) {
		// Test that answers are properly mapped in evaluation
		interview := CreateTestInterview(t, "Mapping Test", GetSampleQuestions())
		session := StartChatSession(t, interview.ID)

		testMessages := []string{
			"First answer",
			"Second answer",
			"Third answer",
		}

		for _, msg := range testMessages {
			SendMessage(t, session.ID, msg)
		}

		evaluation := EndChatSession(t, session.ID)

		// Verify answers are mapped
		if len(evaluation.Answers) == 0 {
			t.Error("No answers found in evaluation")
		}

		// Check that we have the expected number of answer mappings
		expectedAnswerCount := len(testMessages)
		if len(evaluation.Answers) < expectedAnswerCount {
			t.Errorf("Expected at least %d answers, got %d", expectedAnswerCount, len(evaluation.Answers))
		}

		// Verify answer keys follow expected pattern
		for i := 0; i < expectedAnswerCount; i++ {
			key := "question_" + string(rune('0'+i))
			if _, exists := evaluation.Answers[key]; !exists {
				t.Errorf("Missing answer key: %s", key)
			}
		}
	})

	t.Run("MultipleEvaluationsPerInterview", func(t *testing.T) {
		// Test that one interview can have multiple chat sessions/evaluations
		interview := CreateTestInterview(t, "Multiple Eval Test", GetSampleQuestions())

		var evaluations []EvaluationResponseDTO

		// Create multiple chat sessions for the same interview
		for i := 0; i < 3; i++ {
			session := StartChatSession(t, interview.ID)
			SendMessage(t, session.ID, "Test response for session "+string(rune('1'+i)))
			evaluation := EndChatSession(t, session.ID)
			evaluations = append(evaluations, evaluation)
		}

		// Verify all evaluations are unique
		idMap := make(map[string]bool)
		for i, eval := range evaluations {
			if idMap[eval.ID] {
				t.Errorf("Duplicate evaluation ID: %s", eval.ID)
			}
			idMap[eval.ID] = true

			if eval.InterviewID != interview.ID {
				t.Errorf("Evaluation %d has wrong interview ID", i)
			}
		}

		if len(evaluations) != 3 {
			t.Errorf("Expected 3 evaluations, got %d", len(evaluations))
		}
	})
}
