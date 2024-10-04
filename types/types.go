package types

import (
	"encoding/json"
	"time"

	"github.com/tmc/langchaingo/llms"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
}

type ChatMessage struct {
	gorm.Model
	Content      string    `json:"content"`
	Role         string    `json:"role"`
	Date         time.Time `json:"date"`
	UserID       uint      `json:"user_id"`
	User         User      `json:"user" gorm:"constraint:OnDelete:CASCADE;"`
	HideFromUser bool      `json:"hide_from_user"`
}

type ChatMessagePacket struct {
	UserInputStatus string      `json:"user_input_status"`
	Message         ChatMessage `json:"message"`
	Type            string      `json:"type"`
}

type ChatMessagesPacket struct {
	UserInputStatus string        `json:"user_input_status"`
	Type            string        `json:"type"`
	Messages        []ChatMessage `json:"messages"`
}

type AgentDescription struct {
	Name           string
	Caller         string
	SystemMessage  string
	PromptTemplate string
	Model          string
}

type FunctionCall interface {
	Name() string
	Execute(args json.RawMessage) (string, error)
	ToolDefinition() llms.Tool
}

type SearchableDocument struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Summary      string    `json:"summary"`
	Signature    string    `json:"signature"`
	Path         string    `json:"path"`
	ShortSummary string    `json:"short_summary"`
	Authors      []string  `json:"authors"`
	Dbid         int       `json:"dbid"` // This field should not be indexed
	Type         string    `json:"type"`
	LatestCommit time.Time `json:"latest_commit"`
	RepositoryID int       `json:"repository_id"`
	PackageID    int       `json:"package_id"`
	FileID       int       `json:"file_id"`
}

type UserSearchResults struct {
	gorm.Model
	UserID  uint   `json:"user_id" gorm:"constraint:OnDelete:CASCADE;"`
	Results string `json:"results" gorm:"type:json"`
}

type UserAgentType struct {
	gorm.Model
	UserID             uint   `json:"user_id" gorm:"constraint:OnDelete:CASCADE;"`
	AgentType          string `json:"type"`
	SearchableDocument string `json:"searchable_document"`
}

type CtxKey string
