package data

import (
	"fmt"
	"sync"
	"time"
)

// MemoryStore provides in-memory storage for development and testing
// TODO: Replace with proper database implementation
type MemoryStore struct {
	interviews   map[string]*Interview
	evaluations  map[string]*Evaluation
	chatSessions map[string]*ChatSession
	chatMessages map[string][]*ChatMessage
	mu           sync.RWMutex
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		interviews:   make(map[string]*Interview),
		evaluations:  make(map[string]*Evaluation),
		chatSessions: make(map[string]*ChatSession),
		chatMessages: make(map[string][]*ChatMessage),
	}
}

// Global memory store instance
var Store = NewMemoryStore()

// Interview operations
func (ms *MemoryStore) CreateInterview(interview *Interview) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.interviews[interview.ID] = interview
	return nil
}

func (ms *MemoryStore) GetInterview(id string) (*Interview, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	interview, exists := ms.interviews[id]
	if !exists {
		return nil, fmt.Errorf("interview not found")
	}
	return interview, nil
}

func (ms *MemoryStore) GetInterviews() ([]*Interview, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	interviews := make([]*Interview, 0, len(ms.interviews))
	for _, interview := range ms.interviews {
		interviews = append(interviews, interview)
	}
	return interviews, nil
}

// Evaluation operations
func (ms *MemoryStore) CreateEvaluation(evaluation *Evaluation) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.evaluations[evaluation.ID] = evaluation
	return nil
}

func (ms *MemoryStore) GetEvaluation(id string) (*Evaluation, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	evaluation, exists := ms.evaluations[id]
	if !exists {
		return nil, fmt.Errorf("evaluation not found")
	}
	return evaluation, nil
}

// Chat session operations
func (ms *MemoryStore) CreateChatSession(session *ChatSession) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.chatSessions[session.ID] = session
	ms.chatMessages[session.ID] = []*ChatMessage{}
	return nil
}

func (ms *MemoryStore) GetChatSession(id string) (*ChatSession, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	session, exists := ms.chatSessions[id]
	if !exists {
		return nil, fmt.Errorf("chat session not found")
	}
	return session, nil
}

func (ms *MemoryStore) UpdateChatSession(session *ChatSession) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.chatSessions[session.ID]; !exists {
		return fmt.Errorf("chat session not found")
	}
	session.UpdatedAt = time.Now()
	ms.chatSessions[session.ID] = session
	return nil
}

// Chat message operations
func (ms *MemoryStore) AddChatMessage(message *ChatMessage) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.chatMessages[message.SessionID]; !exists {
		return fmt.Errorf("chat session not found")
	}
	ms.chatMessages[message.SessionID] = append(ms.chatMessages[message.SessionID], message)
	return nil
}

func (ms *MemoryStore) GetChatMessages(sessionID string) ([]*ChatMessage, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	messages, exists := ms.chatMessages[sessionID]
	if !exists {
		return nil, fmt.Errorf("chat session not found")
	}
	return messages, nil
}
