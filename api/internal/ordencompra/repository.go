package ordencompra

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	Crear(ctx context.Context, orden *OrdenCompra) error
	BuscarPorID(ctx context.Context, id bson.ObjectID) (*OrdenCompra, error)
	Listar(ctx context.Context) ([]OrdenCompra, error)
	Actualizar(ctx context.Context, id bson.ObjectID, campos bson.M) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{
		collection: collection,
	}
}

func (r *MongoRepository) Crear(
	ctx context.Context,
	orden *OrdenCompra,
) error {
	_, err := r.collection.InsertOne(ctx, orden)
	return err
}

func (r *MongoRepository) BuscarPorID(
	ctx context.Context,
	id bson.ObjectID,
) (*OrdenCompra, error) {
	var orden OrdenCompra

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&orden)

	if err != nil {
		return nil, err
	}

	return &orden, nil
}

func (r *MongoRepository) Listar(
	ctx context.Context,
) ([]OrdenCompra, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var ordenes []OrdenCompra

	if err := cursor.All(ctx, &ordenes); err != nil {
		return nil, err
	}

	return ordenes, nil
}

func (r *MongoRepository) Actualizar(
	ctx context.Context,
	id bson.ObjectID,
	campos bson.M,
) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": campos},
	)

	return err
}
