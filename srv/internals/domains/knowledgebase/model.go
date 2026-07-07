package knowledgebase

import "encoding/json"

type CreateRequest struct {
	SourceId        string          `json:"source_id" binding:"required"`
	SourceEventType string          `json:"source_event_type" binding:"required"`
	Content         json.RawMessage `json:"content" binding:"required"`
	Metadata        json.RawMessage `json:"metadata"`
	Embedding       []float32       `json:"embedding" binding:"required"`
}

type DbCreateKnowledgeBaseRow struct {
	Id         string          `db:"id" json:"id"`
	Pk         int             `db:"pk" json:"-"`
	SourceId   string          `db:"source_id" json:"source_id"`
	SourceType string          `db:"source_type" json:"source_type"`
	Metadata   json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	Content    json.RawMessage `db:"content" json:"content"`
}

type searchInput struct {
	Id         *string
	Embedding  []float32
	SourceType *string
	Page       int
	Limit      int
}

func (s searchInput) getPage() int {
	return max(s.Page, 1)
}
func (s searchInput) getLimit() int {
	if s.Limit == 0 {
		return 20
	}
	return s.Limit
}
func (s searchInput) getOffset() int {
	return (s.getPage() - 1) * s.getLimit()
}
