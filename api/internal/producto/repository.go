package producto

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Producto, error)
	FindByID(ctx context.Context, id bson.ObjectID) (Producto, error)
	Create(ctx context.Context, producto Producto) error
	Update(ctx context.Context, producto Producto) error
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

func (r *MongoRepository) FindAll(ctx context.Context) ([]Producto, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var productos []Producto

	if err := cursor.All(ctx, &productos); err != nil {
		return nil, err
	}

	return productos, nil
}

func (r *MongoRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (Producto, error) {

	var producto Producto

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&producto)

	return producto, err
}

func (r *MongoRepository) Create(
	ctx context.Context,
	producto Producto,
) error {

	_, err := r.collection.InsertOne(ctx, producto)

	return err
}

func (r *MongoRepository) Update(
	ctx context.Context,
	producto Producto,
) error {

	update := bson.M{
		"$set": bson.M{
			"nombre":                    producto.Nombre,
			"descripcion":               producto.Descripcion,
			"categoria_id":              producto.CategoriaID,
			"unidad_medida":             producto.UnidadMedida,
			"costo_unitario":            producto.CostoUnitario,
			"stock_minimo_general":      producto.StockMinimoGeneral,
			"stock_minimo_por_deposito": producto.StockMinimoPorDeposito,
			"fecha_vencimiento":         producto.FechaVencimiento,
			"modificado_por":            producto.ModificadoPor,
			"fecha_actualizacion":       producto.FechaActualizacion,
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": producto.ID},
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
