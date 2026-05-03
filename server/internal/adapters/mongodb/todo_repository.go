package mongodb

import (
	"context"
	"fmt"

	"github.com/tipananchakr/todo-list/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(collection *mongo.Collection) *TodoRepository {
	return &TodoRepository{collection: collection}
}

func (r *TodoRepository) FindAllByUser(ctx context.Context, userID string) ([]domain.Todo, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID, "is_deleted": bson.M{"$ne": true}})
	if err != nil {
		return nil, fmt.Errorf("find todos: %w", err)
	}
	defer cursor.Close(ctx)

	var documents []domain.Todo
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("decode todos: %w", err)
	}

	todos := make([]domain.Todo, 0, len(documents))
	for _, document := range documents {
		todos = append(todos, document)
	}

	return todos, nil
}

func (r *TodoRepository) Create(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	document := domain.Todo{
		UserID:    todo.UserID,
		Body:      todo.Body,
		Completed: todo.Completed,
	}

	result, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		return domain.Todo{}, fmt.Errorf("insert todo: %w", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return domain.Todo{}, fmt.Errorf("insert todo: unexpected inserted id type")
	}

	document.ID = insertedID
	return document, nil
}

func (r *TodoRepository) MarkCompleted(ctx context.Context, userID string, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidTodoID
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID, "userId": userID}, bson.M{
		"$set": bson.M{"completed": true},
	})
	if err != nil {
		return fmt.Errorf("mark todo completed: %w", err)
	}

	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, userID string, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidTodoID
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID, "userId": userID}, bson.M{
		"$set": bson.M{"is_deleted": true},
	})
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}

	return nil
}
