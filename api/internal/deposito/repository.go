package deposito

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Deposito, error)
	FindByID(ctx context.Context, id bson.ObjectID) (Deposito, error)
	Create(ctx context.Context, deposito Deposito) error
	Update(ctx context.Context, deposito Deposito) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{
		collection: collection,
	}
}

func (r *MongoRepository) FindAll(ctx context.Context) ([]Deposito, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var depositos []Deposito

	if err := cursor.All(ctx, &depositos); err != nil {
		return nil, err
	}

	return depositos, nil
}

func (r *MongoRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (Deposito, error) {

	var deposito Deposito

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&deposito)

	return deposito, err
}

func (r *MongoRepository) Create(
	ctx context.Context,
	deposito Deposito,
) error {
	_, err := r.collection.InsertOne(ctx, deposito)
	return err
}

func (r *MongoRepository) Update(
	ctx context.Context,
	deposito Deposito,
) error {

	update := bson.M{
		"$set": bson.M{
			"nombre":              deposito.Nombre,
			"localidad":           deposito.Localidad,
			"provincia":           deposito.Provincia,
			"responsable_id":      deposito.ResponsableID,
			"modificado_por":      deposito.ModificadoPor,
			"fecha_actualizacion": deposito.FechaActualizacion,
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": deposito.ID},
		update,
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
