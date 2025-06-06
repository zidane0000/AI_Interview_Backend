// Data models (structs for DB tables)
package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// StringArray is a custom type for handling PostgreSQL arrays with GORM
type StringArray []string

// Scan implements the Scanner interface for database/sql
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}
}

// Value implements the Valuer interface for database/sql
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// StringMap is a custom type for handling JSON maps with GORM
type StringMap map[string]string

// Scan implements the Scanner interface for database/sql
func (s *StringMap) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return fmt.Errorf("cannot scan %T into StringMap", value)
	}
}

// Value implements the Valuer interface for database/sql
func (s StringMap) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Interview model with proper GORM tags
type Interview struct {
	ID            string      `gorm:"primaryKey;type:varchar(255)" json:"id"`
	CandidateName string      `gorm:"type:varchar(255);not null" json:"candidate_name"`
	Questions     StringArray `gorm:"type:jsonb" json:"questions"`
	Status        string      `gorm:"type:varchar(50);not null;default:'draft'" json:"status"` // "draft", "active", "completed"
	Type          string      `gorm:"type:varchar(50);not null" json:"type"`                   // "general", "technical", "behavioral"
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

// Evaluation model with proper GORM tags
type Evaluation struct {
	ID          string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	InterviewID string    `gorm:"type:varchar(255);not null;index" json:"interview_id"`
	Answers     StringMap `gorm:"type:jsonb" json:"answers"`
	Score       float64   `gorm:"type:decimal(5,2)" json:"score"`
	Feedback    string    `gorm:"type:text" json:"feedback"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// ChatSession model for conversational interviews with proper GORM tags
type ChatSession struct {
	ID          string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	InterviewID string     `gorm:"type:varchar(255);not null;index" json:"interview_id"`
	Status      string     `gorm:"type:varchar(50);not null;default:'active'" json:"status"` // "active", "completed", "abandoned"
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	EndedAt     *time.Time `gorm:"type:timestamp" json:"ended_at,omitempty"`
}

// ChatMessage model with proper GORM tags
type ChatMessage struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	SessionID string    `gorm:"type:varchar(255);not null;index" json:"session_id"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"` // "user", "ai"
	Content   string    `gorm:"type:text;not null" json:"content"`
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
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
