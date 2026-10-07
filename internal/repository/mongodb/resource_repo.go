package mongodb

import (
	"context"
	"errors"
	"time"

	"booking-app/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const resourcesColl = "resources"

// resourceDoc is the storage shape; it never leaks outside this package.
type resourceDoc struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	CreatedAt   time.Time     `bson:"created_at"`
}

func (d resourceDoc) toDomain() domain.Resource {
	return domain.Resource{
		ID:          d.ID.Hex(),
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
	}
}

type ResourceRepo struct {
	coll *mongo.Collection
}

func NewResourceRepo(db *mongo.Database) *ResourceRepo {
	return &ResourceRepo{coll: db.Collection(resourcesColl)}
}

func (r *ResourceRepo) Create(ctx context.Context, res *domain.Resource) error {
	doc := resourceDoc{Name: res.Name, Description: res.Description, CreatedAt: res.CreatedAt}
	out, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	res.ID = out.InsertedID.(bson.ObjectID).Hex()
	return nil
}

func (r *ResourceRepo) GetByID(ctx context.Context, id string) (*domain.Resource, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	var doc resourceDoc
	if err := r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	res := doc.toDomain()
	return &res, nil
}

func (r *ResourceRepo) List(ctx context.Context) ([]domain.Resource, error) {
	cur, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var docs []resourceDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	out := make([]domain.Resource, 0, len(docs))
	for _, d := range docs {
		out = append(out, d.toDomain())
	}
	return out, nil
}
