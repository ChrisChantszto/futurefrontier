package models

import "time"

type I18nString map[string]string

type Background struct {
	Type  string      `bson:"type" json:"type"`  // none|solid|gradient|image|video
	Value interface{} `bson:"value" json:"value"` // color string or {mediaId} or {from,to,angle}
}

type SectionElement struct {
	Title   I18nString `bson:"title,omitempty" json:"title,omitempty"`
	Content I18nString `bson:"content,omitempty" json:"content,omitempty"`
	Color   string     `bson:"color,omitempty" json:"color,omitempty"`
	MediaID string     `bson:"mediaId,omitempty" json:"mediaId,omitempty"`
}

type Section struct {
	Identifier string           `bson:"identifier" json:"identifier"`
	Enabled    bool             `bson:"enabled" json:"enabled"`
	Background Background       `bson:"background" json:"background"`
	Data       map[string]any   `bson:"data" json:"data"` // includes title/content/elements
}

type Page struct {
	ID         string      `bson:"_id,omitempty" json:"id"`
	Identifier string      `bson:"identifier" json:"identifier"`
	Enabled    bool        `bson:"enabled" json:"enabled"`
	Meta       struct {
		Title       map[string]string `bson:"title" json:"title"`
		Description map[string]string `bson:"description" json:"description"`
		Keywords    []string          `bson:"keywords,omitempty" json:"keywords,omitempty"`
		OGImageID   string            `bson:"ogImageId,omitempty" json:"ogImageId,omitempty"`
	} `bson:"meta" json:"meta"`
	Grouping struct {
		ParentIdentifier string            `bson:"parentIdentifier,omitempty" json:"parentIdentifier,omitempty"`
		Label            map[string]string `bson:"label" json:"label"`
	} `bson:"grouping" json:"grouping"`
	Header struct {
		NoticeBar struct {
			Enabled bool              `bson:"enabled" json:"enabled"`
			Text    map[string]string `bson:"text,omitempty" json:"text,omitempty"`
		} `bson:"noticeBar" json:"noticeBar"`
	} `bson:"header" json:"header"`
	Sections  []Section  `bson:"sections" json:"sections"`
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updatedAt" json:"updatedAt"`
}