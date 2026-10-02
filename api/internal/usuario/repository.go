package usuario

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Usuario, error)
	FindByID(ctx context.Context, id string) (Usuario, error)
	FindByEmail(ctx context.Context, email string) (Usuario, error)
	Create(ctx context.Context, usuario Usuario) (Usuario, error)
	Update(ctx context.Context, id string, usuario Usuario) (Usuario, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{
		collection: collection,
	}
}

func (r *MongoRepository) FindAll(
	ctx context.Context,
) ([]Usuario, error) {

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var usuarios []Usuario

	if err := cursor.All(ctx, &usuarios); err != nil {
		return nil, err
	}

	return usuarios, nil
}

func (r *MongoRepository) FindByID(
	ctx context.Context,
	id string,
) (Usuario, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return Usuario{}, err
	}

	var usuario Usuario

	err = r.collection.FindOne(
		ctx,
		bson.M{"_id": objectID},
	).Decode(&usuario)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Usuario{}, errors.New("usuario no encontrado")
		}

		return Usuario{}, err
	}

	return usuario, nil
}

func (r *MongoRepository) FindByEmail(
	ctx context.Context,
	email string,
) (Usuario, error) {

	var usuario Usuario

	err := r.collection.FindOne(
		ctx,
		bson.M{"email": email},
	).Decode(&usuario)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Usuario{}, errors.New("usuario no encontrado")
		}

		return Usuario{}, err
	}

	return usuario, nil
}

func (r *MongoRepository) Create(
	ctx context.Context,
	usuario Usuario,
) (Usuario, error) {

	result, err := r.collection.InsertOne(ctx, usuario)
	if err != nil {
		return Usuario{}, err
	}

	usuario.ID = result.InsertedID.(bson.ObjectID)

	return usuario, nil
}

func (r *MongoRepository) Update(
	ctx context.Context,
	id string,
	usuario Usuario,
) (Usuario, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return Usuario{}, err
	}

	update := bson.M{
		"$set": bson.M{
			"nombre":              usuario.Nombre,
			"email":               usuario.Email,
			"rol":                 usuario.Rol,
			"deposito_id":         usuario.DepositoID,
			"activo":              usuario.Activo,
			"modificado_por":      usuario.ModificadoPor,
			"fecha_actualizacion": usuario.FechaActualizacion,
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)

	if err != nil {
		return Usuario{}, err
	}

	if result.MatchedCount == 0 {
		return Usuario{}, errors.New("usuario no encontrado")
	}

	usuario.ID = objectID

	return usuario, nil
}

var _ Repository = (*MongoRepository)(nil)
