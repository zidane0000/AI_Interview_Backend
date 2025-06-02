// Data models (structs for DB tables)
package data

import (
	"time"
)

// Define Go structs for interviews, interview_logs, and evaluations

// Interview model with proper database tags
type Interview struct {
	ID            string    `db:"id" json:"id"`
	CandidateName string    `db:"candidate_name" json:"candidate_name"`
	Questions     []string  `db:"questions" json:"questions"` // Consider JSON column or separate table
	Status        string    `db:"status" json:"status"`       // "draft", "active", "completed"
	Type          string    `db:"type" json:"type"`           // "general", "technical", "behavioral"
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// Evaluation model
type Evaluation struct {
	ID          string            `db:"id" json:"id"`
	InterviewID string            `db:"interview_id" json:"interview_id"`
	Answers     map[string]string `db:"answers" json:"answers"` // JSON column
	Score       float64           `db:"score" json:"score"`
	Feedback    string            `db:"feedback" json:"feedback"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
}

// ChatSession model for conversational interviews
type ChatSession struct {
	ID          string     `db:"id" json:"id"`
	InterviewID string     `db:"interview_id" json:"interview_id"`
	Status      string     `db:"status" json:"status"` // "active", "completed", "abandoned"
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	EndedAt     *time.Time `db:"ended_at" json:"ended_at,omitempty"`
}

// ChatMessage model
type ChatMessage struct {
	ID        string    `db:"id" json:"id"`
	SessionID string    `db:"session_id" json:"session_id"`
	Type      string    `db:"type" json:"type"` // "user", "ai"
	Content   string    `db:"content" json:"content"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// TODO: Implement File model for resume uploads
// type File struct {
//     ID           string    `db:"id" json:"id"`
//     OriginalName string    `db:"original_name" json:"original_name"`
//     FileName     string    `db:"file_name" json:"file_name"`
//     FilePath     string    `db:"file_path" json:"file_path"`
//     FileSize     int64     `db:"file_size" json:"file_size"`
//     ContentType  string    `db:"content_type" json:"content_type"`
//     InterviewID  *string   `db:"interview_id" json:"interview_id,omitempty"`
//     CreatedAt    time.Time `db:"created_at" json:"created_at"`
// }

// TODO: Add database migration scripts
// TODO: Add indexes for performance optimization
// TODO: Add foreign key constraints
// TODO: Add validation tags for input validation
// TODO: Consider soft delete functionality (deleted_at fields)
// TODO: Add audit trail fields (created_by, updated_by)
// TODO: Add support for database transactions
// TODO: Add model conversion methods (ToDTO, FromDTO)
