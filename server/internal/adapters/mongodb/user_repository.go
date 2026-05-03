package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tipananchakr/todo-list/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"passwordHash"`
	CreatedAt    time.Time          `bson:"createdAt"`
}

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{collection: collection}
}

func (r *UserRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
		Options: options.Index().
			SetName("unique_email").
			SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create user indexes: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, domain.ErrUserNotFound
	}

	return r.findOne(ctx, bson.M{"_id": objectID})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.findOne(ctx, bson.M{"email": email})
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	document := userDocument{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
	}

	result, err := r.collection.InsertOne(ctx, document)
	if mongo.IsDuplicateKeyError(err) {
		return domain.User{}, domain.ErrEmailAlreadyExists
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return domain.User{}, fmt.Errorf("insert user: unexpected inserted id type")
	}

	document.ID = insertedID
	return document.toDomain(), nil
}

func (r *UserRepository) findOne(ctx context.Context, filter bson.M) (domain.User, error) {
	var document userDocument
	err := r.collection.FindOne(ctx, filter).Decode(&document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("find user: %w", err)
	}

	return document.toDomain(), nil
}

func (d userDocument) toDomain() domain.User {
	return domain.User{
		ID:           d.ID.Hex(),
		Email:        d.Email,
		PasswordHash: d.PasswordHash,
		CreatedAt:    d.CreatedAt,
	}
}
