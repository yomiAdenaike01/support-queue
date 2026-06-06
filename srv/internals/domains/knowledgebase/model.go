package knowledgebase

import "encoding/json"

type CreateRequest struct {
	SourceId        string          `json:"source_id" binding:"required"`
	SourceEventType string          `json:"source_event_type" binding:"required"`
	Content         string          `json:"content" binding:"required"`
	Metadata        json.RawMessage `json:"metadata"`
	Embedding       []float32       `json:"embedding" binding:"required"`
}

type DbCreateKnowledgeBaseRow struct {
	Id         string          `db:"id" json:"id"`
	Pk         int             `db:"pk" json:"-"`
	SourceId   string          `db:"source_id" json:"source_id"`
	SourceType string          `db:"source_type" json:"source_type"`
	Metadata   json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	Content    string          `db:"content" json:"content"`
}
