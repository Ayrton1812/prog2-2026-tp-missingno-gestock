package proveedor

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Proveedor, error)
	FindByID(ctx context.Context, id bson.ObjectID) (Proveedor, error)
	Create(ctx context.Context, proveedor Proveedor) error
	Update(ctx context.Context, proveedor Proveedor) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(
	collection *mongo.Collection,
) *MongoRepository {
	return &MongoRepository{
		collection: collection,
	}
}

func (r *MongoRepository) FindAll(
	ctx context.Context,
) ([]Proveedor, error) {

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var proveedores []Proveedor

	if err := cursor.All(ctx, &proveedores); err != nil {
		return nil, err
	}

	return proveedores, nil
}

func (r *MongoRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (Proveedor, error) {

	var proveedor Proveedor

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&proveedor)

	return proveedor, err
}

func (r *MongoRepository) Create(
	ctx context.Context,
	proveedor Proveedor,
) error {

	_, err := r.collection.InsertOne(
		ctx,
		proveedor,
	)

	return err
}

func (r *MongoRepository) Update(
	ctx context.Context,
	proveedor Proveedor,
) error {

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": proveedor.ID},
		bson.M{
			"$set": bson.M{
				"razon_social":        proveedor.RazonSocial,
				"cuit":                proveedor.CUIT,
				"contacto":            proveedor.Contacto,
				"productos":           proveedor.Productos,
				"modificado_por":      proveedor.ModificadoPor,
				"fecha_actualizacion": proveedor.FechaActualizacion,
			},
		},
	)

	return err
}

func (r *MongoRepository) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {

	_, err := r.collection.DeleteOne(
		ctx,
		bson.M{"_id": id},
	)

	return err
}
