package deposito

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindAll(ctx context.Context) ([]DepositoDTO, error) {
	depositos, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	resultado := make([]DepositoDTO, 0, len(depositos))

	for _, deposito := range depositos {
		resultado = append(resultado, convertirDTO(deposito))
	}

	return resultado, nil
}

func (s *Service) FindByID(ctx context.Context, id string) (DepositoDTO, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return DepositoDTO{}, errors.New("id de depósito inválido")
	}

	deposito, err := s.repository.FindByID(ctx, objectID)
	if err != nil {
		return DepositoDTO{}, err
	}

	return convertirDTO(deposito), nil
}

func (s *Service) Create(
	ctx context.Context,
	request CrearDepositoRequest,
	usuarioID string,
) (DepositoDTO, error) {

	if strings.TrimSpace(request.Nombre) == "" {
		return DepositoDTO{}, errors.New("el nombre es obligatorio")
	}

	if strings.TrimSpace(request.Localidad) == "" {
		return DepositoDTO{}, errors.New("la localidad es obligatoria")
	}

	if strings.TrimSpace(request.Provincia) == "" {
		return DepositoDTO{}, errors.New("la provincia es obligatoria")
	}

	var responsableID *bson.ObjectID

	if request.ResponsableID != "" {
		id, err := bson.ObjectIDFromHex(request.ResponsableID)
		if err != nil {
			return DepositoDTO{}, errors.New("responsable_id inválido")
		}

		responsableID = &id
	}

	now := time.Now()

	deposito := Deposito{
		ID:                 bson.NewObjectID(),
		Nombre:             strings.TrimSpace(request.Nombre),
		Localidad:          strings.TrimSpace(request.Localidad),
		Provincia:          strings.TrimSpace(request.Provincia),
		ResponsableID:      responsableID,
		CreadoPor:          usuarioID,
		ModificadoPor:      usuarioID,
		FechaCreacion:      now,
		FechaActualizacion: now,
	}

	err := s.repository.Create(ctx, deposito)
	if err != nil {
		return DepositoDTO{}, err
	}

	return convertirDTO(deposito), nil
}

func (s *Service) Update(
	ctx context.Context,
	id string,
	request ActualizarDepositoRequest,
	usuarioID string,
) (DepositoDTO, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return DepositoDTO{}, errors.New("id de depósito inválido")
	}

	if strings.TrimSpace(request.Nombre) == "" {
		return DepositoDTO{}, errors.New("el nombre es obligatorio")
	}

	if strings.TrimSpace(request.Localidad) == "" {
		return DepositoDTO{}, errors.New("la localidad es obligatoria")
	}

	if strings.TrimSpace(request.Provincia) == "" {
		return DepositoDTO{}, errors.New("la provincia es obligatoria")
	}

	var responsableID *bson.ObjectID

	if request.ResponsableID != "" {
		idResponsable, err := bson.ObjectIDFromHex(request.ResponsableID)
		if err != nil {
			return DepositoDTO{}, errors.New("responsable_id inválido")
		}

		responsableID = &idResponsable
	}

	deposito := Deposito{
		ID:                 objectID,
		Nombre:             strings.TrimSpace(request.Nombre),
		Localidad:          strings.TrimSpace(request.Localidad),
		Provincia:          strings.TrimSpace(request.Provincia),
		ResponsableID:      responsableID,
		ModificadoPor:      usuarioID,
		FechaActualizacion: time.Now(),
	}

	err = s.repository.Update(ctx, deposito)
	if err != nil {
		return DepositoDTO{}, err
	}

	return convertirDTO(deposito), nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id de depósito inválido")
	}

	return s.repository.Delete(ctx, objectID)
}

func convertirDTO(deposito Deposito) DepositoDTO {
	dto := DepositoDTO{
		ID:        deposito.ID.Hex(),
		Nombre:    deposito.Nombre,
		Localidad: deposito.Localidad,
		Provincia: deposito.Provincia,
	}

	if deposito.ResponsableID != nil {
		dto.ResponsableID = deposito.ResponsableID.Hex()
	}

	return dto
}
