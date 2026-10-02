package categoria

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Categoria, error)
	FindByID(ctx context.Context, id string) (Categoria, error)
	Create(ctx context.Context, categoria Categoria) (Categoria, error)
	Update(ctx context.Context, id string, categoria Categoria) (Categoria, error)
	Delete(ctx context.Context, id string) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{
		collection: collection,
	}
}

func (r *MongoRepository) FindAll(ctx context.Context) ([]Categoria, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var categorias []Categoria

	if err := cursor.All(ctx, &categorias); err != nil {
		return nil, err
	}

	return categorias, nil
}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (Categoria, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return Categoria{}, err
	}

	var categoria Categoria

	err = r.collection.FindOne(
		ctx,
		bson.M{"_id": objectID},
	).Decode(&categoria)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Categoria{}, errors.New("categoria no encontrada")
		}

		return Categoria{}, err
	}

	return categoria, nil
}

func (r *MongoRepository) Create(ctx context.Context, categoria Categoria) (Categoria, error) {
	result, err := r.collection.InsertOne(ctx, categoria)
	if err != nil {
		return Categoria{}, err
	}

	categoria.ID = result.InsertedID.(bson.ObjectID)

	return categoria, nil
}

func (r *MongoRepository) Update(
	ctx context.Context,
	id string,
	categoria Categoria,
) (Categoria, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return Categoria{}, err
	}

	update := bson.M{
		"$set": bson.M{
			"nombre":              categoria.Nombre,
			"modificado_por":      categoria.ModificadoPor,
			"fecha_actualizacion": categoria.FechaActualizacion,
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)

	if err != nil {
		return Categoria{}, err
	}

	if result.MatchedCount == 0 {
		return Categoria{}, errors.New("categoria no encontrada")
	}

	categoria.ID = objectID

	return categoria, nil
}

func (r *MongoRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(
		ctx,
		bson.M{"_id": objectID},
	)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("categoria no encontrada")
	}

	return nil
}

var _ Repository = (*MongoRepository)(nil)
