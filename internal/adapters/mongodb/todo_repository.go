package mongodb

import (
	"context"
	"fmt"

	"github.com/tipananchakr/todo-list/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type todoDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Completed bool               `bson:"completed"`
	Body      string             `bson:"body"`
}

type TodoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(collection *mongo.Collection) *TodoRepository {
	return &TodoRepository{collection: collection}
}

func (r *TodoRepository) FindAll(ctx context.Context) ([]domain.Todo, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("find todos: %w", err)
	}
	defer cursor.Close(ctx)

	var documents []todoDocument
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("decode todos: %w", err)
	}

	todos := make([]domain.Todo, 0, len(documents))
	for _, document := range documents {
		todos = append(todos, document.toDomain())
	}

	return todos, nil
}

func (r *TodoRepository) Create(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	document := todoDocument{
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
	return document.toDomain(), nil
}

func (r *TodoRepository) MarkCompleted(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidTodoID
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{"completed": true},
	})
	if err != nil {
		return fmt.Errorf("mark todo completed: %w", err)
	}

	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidTodoID
	}

	if _, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID}); err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}

	return nil
}

func (d todoDocument) toDomain() domain.Todo {
	return domain.Todo{
		ID:        d.ID.Hex(),
		Completed: d.Completed,
		Body:      d.Body,
	}
}
