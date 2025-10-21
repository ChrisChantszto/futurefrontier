package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ChrisChantszto/futurefrontier/internal/models"
)

type PagesService struct {
	Coll *mongo.Collection
}

func NewPagesService(db *mongo.Database) *PagesService {
	return &PagesService{Coll: db.Collection("pages")}
}

func (s *PagesService) List(ctx context.Context, limit int64) ([]models.Page, error) {
	opts := options.Find().SetLimit(limit).SetSort(bson.D{{Key: "identifier", Value: 1}})
	cur, err := s.Coll.Find(ctx, bson.M{}, opts)
	if err != nil { return nil, err }
	defer cur.Close(ctx)
	var pages []models.Page
	for cur.Next(ctx) {
		var p models.Page
		if err := cur.Decode(&p); err == nil { pages = append(pages, p) }
	}
	return pages, cur.Err()
}

func (s *PagesService) GetByIdentifier(ctx context.Context, identifier string) (*models.Page, error) {
	var p models.Page
	err := s.Coll.FindOne(ctx, bson.M{"identifier": identifier}).Decode(&p)
	return &p, err
}

func (s *PagesService) Create(ctx context.Context, p *models.Page) (*models.Page, error) {
	now := time.Now().UTC()
	p.CreatedAt, p.UpdatedAt = now, now
	_, err := s.Coll.InsertOne(ctx, p)
	return p, err
}

func (s *PagesService) UpdateByIdentifier(ctx context.Context, identifier string, update bson.M) (*models.Page, error) {
	if update == nil { update = bson.M{} }
	// never allow updating identifier via patch
	delete(update, "identifier")
	update["updatedAt"] = time.Now().UTC()
	_, err := s.Coll.UpdateOne(ctx, bson.M{"identifier": identifier}, bson.M{"$set": update})
	if err != nil { return nil, err }
	return s.GetByIdentifier(ctx, identifier)
}

func (s *PagesService) DeleteByIdentifier(ctx context.Context, identifier string) error {
	_, err := s.Coll.DeleteOne(ctx, bson.M{"identifier": identifier})
	return err
}

// Import upserts by identifier
func (s *PagesService) Import(ctx context.Context, pages []models.Page) (int, error) {
	count := 0
	for _, p := range pages {
		p.UpdatedAt = time.Now().UTC()
		if p.CreatedAt.IsZero() {
			p.CreatedAt = p.UpdatedAt
		}
		_, err := s.Coll.UpdateOne(ctx,
			bson.M{"identifier": p.Identifier},
			bson.M{"$set": p},
			options.Update().SetUpsert(true),
		)
		if err != nil { return count, err }
		count++
	}
	return count, nil
}

// GenerateTemplate returns a default page skeleton
func (s *PagesService) GenerateTemplate(ctx context.Context, identifier string) *models.Page {
	now := time.Now().UTC()
	return &models.Page{
		Identifier: identifier,
		Enabled:    true,
		Meta: struct {
			Title       map[string]string `bson:"title" json:"title"`
			Description map[string]string `bson:"description" json:"description"`
			Keywords    []string          `bson:"keywords,omitempty" json:"keywords,omitempty"`
			OGImageID   string            `bson:"ogImageId,omitempty" json:"ogImageId,omitempty"`
		}{
			Title:       map[string]string{"en": identifier, "zh-Hant": identifier, "zh-Hans": identifier},
			Description: map[string]string{"en": "", "zh-Hant": "", "zh-Hans": ""},
			Keywords:    []string{},
		},
		Grouping: struct {
			ParentIdentifier string            `bson:"parentIdentifier,omitempty" json:"parentIdentifier,omitempty"`
			Label            map[string]string `bson:"label" json:"label"`
		}{
			Label: map[string]string{"en": "Default", "zh-Hant": "默認", "zh-Hans": "默认"},
		},
		Header: struct {
			NoticeBar struct {
				Enabled bool              `bson:"enabled" json:"enabled"`
				Text    map[string]string `bson:"text,omitempty" json:"text,omitempty"`
			} `bson:"noticeBar" json:"noticeBar"`
		}{
			NoticeBar: struct {
				Enabled bool              `bson:"enabled" json:"enabled"`
				Text    map[string]string `bson:"text,omitempty" json:"text,omitempty"`
			}{Enabled: false, Text: map[string]string{}},
		},
		Sections:  []models.Section{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
