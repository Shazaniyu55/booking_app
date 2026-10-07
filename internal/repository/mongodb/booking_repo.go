package mongodb

import (
	"context"
	"errors"
	"time"

	"booking-app/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const bookingsColl = "bookings"

type bookingDoc struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	ResourceID    bson.ObjectID `bson:"resource_id"`
	CustomerName  string        `bson:"customer_name"`
	CustomerEmail string        `bson:"customer_email"`
	StartTime     time.Time     `bson:"start_time"`
	EndTime       time.Time     `bson:"end_time"`
	Status        string        `bson:"status"`
	CreatedAt     time.Time     `bson:"created_at"`
}

func (d bookingDoc) toDomain() domain.Booking {
	return domain.Booking{
		ID:            d.ID.Hex(),
		ResourceID:    d.ResourceID.Hex(),
		CustomerName:  d.CustomerName,
		CustomerEmail: d.CustomerEmail,
		StartTime:     d.StartTime,
		EndTime:       d.EndTime,
		Status:        domain.BookingStatus(d.Status),
		CreatedAt:     d.CreatedAt,
	}
}

type BookingRepo struct {
	coll *mongo.Collection
}

func NewBookingRepo(db *mongo.Database) *BookingRepo {
	return &BookingRepo{coll: db.Collection(bookingsColl)}
}

func (r *BookingRepo) Create(ctx context.Context, b *domain.Booking) error {
	rid, err := bson.ObjectIDFromHex(b.ResourceID)
	if err != nil {
		return domain.ErrNotFound
	}
	doc := bookingDoc{
		ResourceID:    rid,
		CustomerName:  b.CustomerName,
		CustomerEmail: b.CustomerEmail,
		StartTime:     b.StartTime,
		EndTime:       b.EndTime,
		Status:        string(b.Status),
		CreatedAt:     b.CreatedAt,
	}
	out, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	b.ID = out.InsertedID.(bson.ObjectID).Hex()
	return nil
}

func (r *BookingRepo) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	var doc bookingDoc
	if err := r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	b := doc.toDomain()
	return &b, nil
}

func (r *BookingRepo) List(ctx context.Context, resourceID string) ([]domain.Booking, error) {
	filter := bson.M{}
	if resourceID != "" {
		rid, err := bson.ObjectIDFromHex(resourceID)
		if err != nil {
			return []domain.Booking{}, nil
		}
		filter["resource_id"] = rid
	}
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: 1}})
	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var docs []bookingDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	out := make([]domain.Booking, 0, len(docs))
	for _, d := range docs {
		out = append(out, d.toDomain())
	}
	return out, nil
}

// HasOverlap: two ranges overlap when existing.start < new.end AND existing.end > new.start.
// Cancelled bookings don't block the slot.
func (r *BookingRepo) HasOverlap(ctx context.Context, resourceID string, start, end time.Time) (bool, error) {
	rid, err := bson.ObjectIDFromHex(resourceID)
	if err != nil {
		return false, domain.ErrNotFound
	}
	filter := bson.M{
		"resource_id": rid,
		"status":      string(domain.StatusConfirmed),
		"start_time":  bson.M{"$lt": end},
		"end_time":    bson.M{"$gt": start},
	}
	n, err := r.coll.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *BookingRepo) UpdateStatus(ctx context.Context, id string, status domain.BookingStatus) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrNotFound
	}
	res, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"status": string(status)}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}
