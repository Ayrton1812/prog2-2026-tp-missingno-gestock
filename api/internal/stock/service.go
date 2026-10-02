package stock

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CrearMovimiento(
	ctx context.Context,
	request CrearMovimientoRequest,
	usuarioID string,
) error {

	productoID, err := bson.ObjectIDFromHex(request.ProductoID)
	if err != nil {
		return errors.New("producto_id inválido")
	}

	depositoID, err := bson.ObjectIDFromHex(request.DepositoID)
	if err != nil {
		return errors.New("deposito_id inválido")
	}

	if request.Cantidad <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}

	tipo := TipoMovimiento(request.Tipo)

	switch tipo {
	case MovimientoEntrada:
	case MovimientoSalida:
	case MovimientoAjuste:
	default:
		return errors.New(
			"tipo de movimiento inválido: usar entrada, salida o ajuste",
		)
	}

	if strings.TrimSpace(request.Motivo) == "" {
		return errors.New("el motivo es obligatorio")
	}

	stock, err := s.repository.FindStock(
		ctx,
		productoID,
		depositoID,
	)

	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}

		stock = Stock{
			ID:                 bson.NewObjectID(),
			ProductoID:         productoID,
			DepositoID:         depositoID,
			Cantidad:           0,
			CreadoPor:          usuarioID,
			ModificadoPor:      usuarioID,
			FechaCreacion:      time.Now(),
			FechaActualizacion: time.Now(),
		}
	}

	nuevaCantidad := stock.Cantidad

	switch tipo {

	case MovimientoEntrada:
		nuevaCantidad += request.Cantidad

	case MovimientoSalida:
		if stock.Cantidad < request.Cantidad {
			return errors.New(
				"stock insuficiente: la operación dejaría el stock en negativo",
			)
		}

		nuevaCantidad -= request.Cantidad

	case MovimientoAjuste:
		nuevaCantidad = request.Cantidad
	}

	if nuevaCantidad < 0 {
		return errors.New("el stock no puede ser negativo")
	}

	ahora := time.Now()

	if stock.ID.IsZero() {
		stock.ID = bson.NewObjectID()
		stock.CreadoPor = usuarioID
		stock.FechaCreacion = ahora
	}

	stock.Cantidad = nuevaCantidad
	stock.ModificadoPor = usuarioID
	stock.FechaActualizacion = ahora

	if stock.FechaCreacion.IsZero() {
		stock.FechaCreacion = ahora
	}

	if err := s.guardarStock(
		ctx,
		stock,
	); err != nil {
		return err
	}

	movimiento := Movimiento{
		ID:         bson.NewObjectID(),
		ProductoID: productoID,
		DepositoID: depositoID,
		Cantidad:   request.Cantidad,
		Tipo:       tipo,
		UsuarioID:  usuarioID,
		Motivo:     strings.TrimSpace(request.Motivo),
		FechaHora:  ahora,
	}

	return s.repository.CreateMovimiento(
		ctx,
		movimiento,
	)
}

func (s *Service) guardarStock(
	ctx context.Context,
	stock Stock,
) error {

	existente, err := s.repository.FindStock(
		ctx,
		stock.ProductoID,
		stock.DepositoID,
	)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return s.repository.CreateStock(ctx, stock)
		}

		return err
	}

	stock.ID = existente.ID

	return s.repository.UpdateStock(ctx, stock)
}

func (s *Service) FindMovimientos(
	ctx context.Context,
) ([]MovimientoDTO, error) {

	movimientos, err := s.repository.FindMovimientos(ctx)

	if err != nil {
		return nil, err
	}

	resultado := make([]MovimientoDTO, 0, len(movimientos))

	for _, movimiento := range movimientos {
		resultado = append(
			resultado,
			convertirMovimientoDTO(movimiento),
		)
	}

	return resultado, nil
}

func convertirMovimientoDTO(
	movimiento Movimiento,
) MovimientoDTO {

	return MovimientoDTO{
		ID:         movimiento.ID.Hex(),
		ProductoID: movimiento.ProductoID.Hex(),
		DepositoID: movimiento.DepositoID.Hex(),
		Cantidad:   movimiento.Cantidad,
		Tipo:       string(movimiento.Tipo),
		UsuarioID:  movimiento.UsuarioID,
		Motivo:     movimiento.Motivo,
		FechaHora:  movimiento.FechaHora.Format(time.RFC3339),
	}
}
