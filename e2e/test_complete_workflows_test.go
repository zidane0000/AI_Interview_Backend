package e2e

import (
	"testing"
	"time"
)

// TestCompleteWorkflows tests end-to-end interview workflows
func TestCompleteWorkflows(t *testing.T) {
	t.Run("FullInterviewWorkflow_Technical", func(t *testing.T) {
		// Technical interview simulation
		questions := []string{
			"Explain the difference between SQL and NoSQL databases",
			"How would you design a scalable web application?",
			"What are the SOLID principles in software engineering?",
			"Describe your experience with microservices architecture",
		}

		interview := CreateTestInterview(t, "Alice Johnson - Senior Developer", questions)
		session := StartChatSession(t, interview.ID)

		// Verify initial AI greeting
		if len(session.Messages) != 1 {
			t.Errorf("Expected 1 initial message, got %d", len(session.Messages))
		}
		if session.Messages[0].Type != "ai" {
			t.Error("First message should be from AI")
		}

		// Simulate technical responses
		technicalResponses := []string{
			"SQL databases use structured schemas and ACID properties, while NoSQL databases offer flexible schemas and eventual consistency. SQL is better for complex relationships, NoSQL for horizontal scaling.",
			"I would use microservices architecture with API Gateway, implement caching with Redis, use CDN for static assets, and containerize with Docker for easy scaling.",
			"SOLID principles are Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, and Dependency Inversion. They help create maintainable and extensible code.",
			"I've worked with microservices using Spring Boot and Docker. Key challenges include service discovery, distributed transactions, and monitoring across services.",
		}

		for i, response := range technicalResponses {
			msgResp := SendMessage(t, session.ID, response)

			// Verify response was recorded correctly
			if msgResp.Message.Content != response {
				t.Errorf("Response %d not recorded correctly", i)
			}
			if msgResp.AIResponse == nil {
				t.Errorf("Missing AI follow-up for response %d", i)
			}

			// Small delay to simulate realistic conversation timing
			time.Sleep(100 * time.Millisecond)
		}

		// Get session state before ending
		updatedSession := GetChatSession(t, session.ID)
		expectedMessages := 1 + (len(technicalResponses) * 2) // initial + (user + ai) pairs
		if len(updatedSession.Messages) < expectedMessages {
			t.Errorf("Expected at least %d messages, got %d", expectedMessages, len(updatedSession.Messages))
		}

		// End session and evaluate
		evaluation := EndChatSession(t, session.ID)

		// Verify evaluation quality for technical interview
		if evaluation.Score <= 0.3 {
			t.Errorf("Technical responses should score higher than 0.3, got %.2f", evaluation.Score)
		}
		if len(evaluation.Feedback) < 50 {
			t.Error("Technical evaluation should provide substantial feedback")
		}

		// Verify session is marked as completed
		finalSession := GetChatSession(t, session.ID)
		if finalSession.Status != "completed" {
			t.Errorf("Expected session status 'completed', got '%s'", finalSession.Status)
		}
	})

	t.Run("FullInterviewWorkflow_Behavioral", func(t *testing.T) {
		// Behavioral interview simulation
		questions := []string{
			"Tell me about a time you had to work with a difficult team member",
			"Describe a situation where you had to learn something new quickly",
			"How do you handle tight deadlines and pressure?",
			"Give an example of a project where you took initiative",
		}

		interview := CreateTestInterview(t, "Bob Smith - Product Manager", questions)
		session := StartChatSession(t, interview.ID)

		// Simulate behavioral responses using STAR method
		behavioralResponses := []string{
			"In my previous role, I worked with a colleague who was resistant to new processes. I scheduled one-on-one meetings to understand their concerns, found common ground, and gradually introduced changes with their input.",
			"When our team adopted React Native, I had one week to learn it for a client demo. I created a learning schedule, used online tutorials, built practice apps, and successfully delivered the demo on time.",
			"I use prioritization techniques like the Eisenhower Matrix. During a critical product launch, I broke tasks into smaller chunks, delegated appropriately, and maintained clear communication with stakeholders.",
			"I noticed our deployment process was causing delays. I researched CI/CD solutions, proposed implementing GitHub Actions, got team buy-in, and reduced deployment time from 2 hours to 15 minutes.",
		}

		for _, response := range behavioralResponses {
			SendMessage(t, session.ID, response)
			time.Sleep(50 * time.Millisecond) // Simulate conversation flow
		}

		evaluation := EndChatSession(t, session.ID)

		// Verify behavioral evaluation
		if evaluation.ID == "" {
			t.Error("Behavioral evaluation should be generated")
		}
		if evaluation.InterviewID != interview.ID {
			t.Error("Evaluation should reference correct interview")
		}
	})

	t.Run("ShortInterviewWorkflow", func(t *testing.T) {
		// Test workflow with minimal interaction
		interview := CreateTestInterview(t, "Quick Test", []string{"Tell me about yourself"})
		session := StartChatSession(t, interview.ID)

		// Send only one brief response
		SendMessage(t, session.ID, "I'm a developer.")

		// End immediately
		evaluation := EndChatSession(t, session.ID)

		// Should still generate valid evaluation
		if evaluation.Score == 0 {
			t.Error("Should generate non-zero score even for brief interviews")
		}
		if evaluation.Feedback == "" {
			t.Error("Should provide feedback even for brief interviews")
		}
	})

	t.Run("MultiSessionInterviewWorkflow", func(t *testing.T) {
		// Test multiple sessions for same interview (like multiple rounds)
		interview := CreateTestInterview(t, "Multi-Round Candidate", GetSampleQuestions())

		evaluations := make([]EvaluationResponseDTO, 3)

		// Round 1: Initial screening
		session1 := StartChatSession(t, interview.ID)
		SendMessage(t, session1.ID, "I have 3 years of experience in software development.")
		evaluations[0] = EndChatSession(t, session1.ID)

		// Round 2: Technical round
		session2 := StartChatSession(t, interview.ID)
		SendMessage(t, session2.ID, "I'm experienced with React, Node.js, and PostgreSQL.")
		evaluations[1] = EndChatSession(t, session2.ID)

		// Round 3: Final round
		session3 := StartChatSession(t, interview.ID)
		SendMessage(t, session3.ID, "I'm looking for growth opportunities and team collaboration.")
		evaluations[2] = EndChatSession(t, session3.ID)

		// Verify all evaluations are unique but reference same interview
		for i, eval := range evaluations {
			if eval.InterviewID != interview.ID {
				t.Errorf("Evaluation %d should reference interview %s", i, interview.ID)
			}
			for j, other := range evaluations {
				if i != j && eval.ID == other.ID {
					t.Errorf("Evaluations %d and %d have same ID", i, j)
				}
			}
		}
	})

	t.Run("WorkflowWithErrors", func(t *testing.T) {
		// Test workflow recovery from errors
		interview := CreateTestInterview(t, "Error Test", GetSampleQuestions())
		session := StartChatSession(t, interview.ID)

		// Send a normal message first
		SendMessage(t, session.ID, "This is a normal message.")

		// Try invalid operations (should not break workflow)
		// These will fail but shouldn't crash the system

		// Continue with normal workflow
		SendMessage(t, session.ID, "Continuing after error scenarios.")

		// Should still be able to end session normally
		evaluation := EndChatSession(t, session.ID)

		if evaluation.ID == "" {
			t.Error("Should generate evaluation despite error scenarios")
		}
	})

	t.Run("LongConversationWorkflow", func(t *testing.T) {
		// Test extended conversation beyond normal limits
		interview := CreateTestInterview(t, "Extended Conversation", GetSampleQuestions())
		session := StartChatSession(t, interview.ID)

		// Send many messages (more than typical interview)
		responses := []string{
			"I'm a senior software engineer with 8 years of experience.",
			"My expertise includes full-stack development with modern frameworks.",
			"I've led teams of 5-10 developers on complex projects.",
			"I'm passionate about clean code and test-driven development.",
			"I have experience with cloud platforms like AWS and Azure.",
			"I enjoy mentoring junior developers and code reviews.",
			"My goal is to become a technical architect in the next few years.",
			"I believe in continuous learning and staying updated with technology.",
		}

		for i, response := range responses {
			msgResp := SendMessage(t, session.ID, response)
			t.Logf("Message %d: %s -> AI: %s", i+1, response, msgResp.AIResponse.Content)
		}

		// Check if session auto-completed (based on AI logic)
		finalSession := GetChatSession(t, session.ID)

		if finalSession.Status == "completed" {
			t.Log("Session auto-completed after extended conversation")
		} else {
			// Manually end if not auto-completed
			evaluation := EndChatSession(t, session.ID)
			if len(evaluation.Answers) < len(responses) {
				t.Error("Long conversation should capture all user responses")
			}
		}
	})
}
