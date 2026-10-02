package stock

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	FindStock(
		ctx context.Context,
		productoID bson.ObjectID,
		depositoID bson.ObjectID,
	) (Stock, error)

	CreateStock(
		ctx context.Context,
		stock Stock,
	) error

	UpdateStock(
		ctx context.Context,
		stock Stock,
	) error

	CreateMovimiento(
		ctx context.Context,
		movimiento Movimiento,
	) error

	FindMovimientos(
		ctx context.Context,
	) ([]Movimiento, error)
}

type MongoRepository struct {
	stockCollection       *mongo.Collection
	movimientosCollection *mongo.Collection
}

func NewMongoRepository(
	database *mongo.Database,
) *MongoRepository {

	return &MongoRepository{
		stockCollection:       database.Collection("stock"),
		movimientosCollection: database.Collection("movimientos"),
	}
}

func (r *MongoRepository) FindStock(
	ctx context.Context,
	productoID bson.ObjectID,
	depositoID bson.ObjectID,
) (Stock, error) {

	var stock Stock

	err := r.stockCollection.FindOne(
		ctx,
		bson.M{
			"producto_id": productoID,
			"deposito_id": depositoID,
		},
	).Decode(&stock)

	return stock, err
}

func (r *MongoRepository) CreateStock(
	ctx context.Context,
	stock Stock,
) error {

	_, err := r.stockCollection.InsertOne(
		ctx,
		stock,
	)

	return err
}

func (r *MongoRepository) UpdateStock(
	ctx context.Context,
	stock Stock,
) error {

	_, err := r.stockCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": stock.ID,
		},
		bson.M{
			"$set": bson.M{
				"cantidad":            stock.Cantidad,
				"modificado_por":      stock.ModificadoPor,
				"fecha_actualizacion": stock.FechaActualizacion,
			},
		},
	)

	return err
}

func (r *MongoRepository) CreateMovimiento(
	ctx context.Context,
	movimiento Movimiento,
) error {

	_, err := r.movimientosCollection.InsertOne(
		ctx,
		movimiento,
	)

	return err
}

func (r *MongoRepository) FindMovimientos(
	ctx context.Context,
) ([]Movimiento, error) {

	cursor, err := r.movimientosCollection.Find(
		ctx,
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var movimientos []Movimiento

	if err := cursor.All(ctx, &movimientos); err != nil {
		return nil, err
	}

	return movimientos, nil
}
