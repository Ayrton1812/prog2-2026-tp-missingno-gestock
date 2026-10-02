package transferencia

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"gestock/api/internal/stock"
	"gestock/api/internal/usuario"
)

type Service struct {
	repository   Repository
	stockService *stock.Service
}

func NewService(
	repository Repository,
	stockService *stock.Service,
) *Service {

	return &Service{
		repository:   repository,
		stockService: stockService,
	}
}

func (s *Service) Crear(
	ctx context.Context,
	request CrearTransferenciaRequest,
	usuarioID string,
	rol string,
	depositoUsuarioID string,
) (TransferenciaDTO, error) {

	productoID, err := bson.ObjectIDFromHex(request.ProductoID)
	if err != nil {
		return TransferenciaDTO{}, errors.New(
			"producto_id inválido",
		)
	}

	origenID, err := bson.ObjectIDFromHex(
		request.DepositoOrigenID,
	)
	if err != nil {
		return TransferenciaDTO{}, errors.New(
			"deposito_origen_id inválido",
		)
	}

	destinoID, err := bson.ObjectIDFromHex(
		request.DepositoDestinoID,
	)
	if err != nil {
		return TransferenciaDTO{}, errors.New(
			"deposito_destino_id inválido",
		)
	}

	if origenID == destinoID {
		return TransferenciaDTO{}, errors.New(
			"el depósito origen y destino deben ser diferentes",
		)
	}

	if request.Cantidad <= 0 {
		return TransferenciaDTO{}, errors.New(
			"la cantidad debe ser mayor a cero",
		)
	}

	if strings.TrimSpace(request.Motivo) == "" {
		return TransferenciaDTO{}, errors.New(
			"el motivo es obligatorio",
		)
	}

	// Gerente y operario solamente pueden solicitar
	// desde su depósito asignado.
	if rol != usuario.RolAdministrador &&
		rol != usuario.RolAuditor {

		if depositoUsuarioID == "" ||
			depositoUsuarioID != request.DepositoOrigenID {

			return TransferenciaDTO{}, errors.New(
				"no puede solicitar una transferencia desde este depósito",
			)
		}
	}

	ahora := time.Now()

	transferencia := Transferencia{
		ID:                bson.NewObjectID(),
		ProductoID:        productoID,
		DepositoOrigenID:  origenID,
		DepositoDestinoID: destinoID,
		Cantidad:          request.Cantidad,
		Estado:            EstadoSolicitada,
		Motivo:            strings.TrimSpace(request.Motivo),
		SolicitadaPor:     usuarioID,
		FechaSolicitud:    &ahora,
	}

	if err := s.repository.Create(
		ctx,
		transferencia,
	); err != nil {
		return TransferenciaDTO{}, err
	}

	return convertirDTO(transferencia), nil
}

func (s *Service) Aprobar(
	ctx context.Context,
	id string,
	usuarioID string,
	rol string,
	depositoUsuarioID string,
) (TransferenciaDTO, error) {

	if rol != usuario.RolAdministrador &&
		rol != usuario.RolGerenteDeposito {

		return TransferenciaDTO{}, errors.New(
			"solo un administrador o gerente de depósito puede aprobar transferencias",
		)
	}

	transferencia, err := s.buscar(ctx, id)

	if err != nil {
		return TransferenciaDTO{}, err
	}

	if transferencia.Estado != EstadoSolicitada {
		return TransferenciaDTO{}, errors.New(
			"la transferencia no está en estado solicitada",
		)
	}

	// Un gerente solamente puede aprobar
	// transferencias destinadas a su depósito.
	if rol == usuario.RolGerenteDeposito &&
		depositoUsuarioID != transferencia.DepositoDestinoID.Hex() {

		return TransferenciaDTO{}, errors.New(
			"solo puede aprobar transferencias destinadas a su depósito",
		)
	}

	// Descontamos del origen.
	movimiento := stock.CrearMovimientoRequest{
		ProductoID: transferencia.ProductoID.Hex(),
		DepositoID: transferencia.DepositoOrigenID.Hex(),
		Cantidad:   transferencia.Cantidad,
		Tipo:       string(stock.MovimientoSalida),
		Motivo: "Transferencia aprobada hacia depósito " +
			transferencia.DepositoDestinoID.Hex(),
	}

	if err := s.stockService.CrearMovimiento(
		ctx,
		movimiento,
		usuarioID,
	); err != nil {
		return TransferenciaDTO{}, err
	}

	ahora := time.Now()

	transferencia.Estado = EstadoAprobada
	transferencia.AprobadaPor = usuarioID
	transferencia.FechaAprobacion = &ahora

	if err := s.repository.Update(
		ctx,
		transferencia,
	); err != nil {
		return TransferenciaDTO{}, err
	}

	return convertirDTO(transferencia), nil
}

func (s *Service) Rechazar(
	ctx context.Context,
	id string,
	motivo string,
	usuarioID string,
	rol string,
	depositoUsuarioID string,
) (TransferenciaDTO, error) {

	if rol != usuario.RolAdministrador &&
		rol != usuario.RolGerenteDeposito {

		return TransferenciaDTO{}, errors.New(
			"solo un administrador o gerente de depósito puede rechazar transferencias",
		)
	}

	transferencia, err := s.buscar(ctx, id)

	if err != nil {
		return TransferenciaDTO{}, err
	}

	if transferencia.Estado != EstadoSolicitada {
		return TransferenciaDTO{}, errors.New(
			"la transferencia no está en estado solicitada",
		)
	}

	if rol == usuario.RolGerenteDeposito &&
		depositoUsuarioID != transferencia.DepositoDestinoID.Hex() {

		return TransferenciaDTO{}, errors.New(
			"solo puede rechazar transferencias destinadas a su depósito",
		)
	}

	if strings.TrimSpace(motivo) == "" {
		return TransferenciaDTO{}, errors.New(
			"el motivo del rechazo es obligatorio",
		)
	}

	ahora := time.Now()

	transferencia.Estado = EstadoRechazada
	transferencia.MotivoRechazo = strings.TrimSpace(motivo)
	transferencia.AprobadaPor = usuarioID
	transferencia.FechaRechazo = &ahora

	if err := s.repository.Update(
		ctx,
		transferencia,
	); err != nil {
		return TransferenciaDTO{}, err
	}

	return convertirDTO(transferencia), nil
}

func (s *Service) Completar(
	ctx context.Context,
	id string,
	usuarioID string,
	rol string,
	depositoUsuarioID string,
) (TransferenciaDTO, error) {

	if rol != usuario.RolAdministrador &&
		rol != usuario.RolGerenteDeposito &&
		rol != usuario.RolOperarioDeposito {

		return TransferenciaDTO{}, errors.New(
			"no tiene permisos para completar transferencias",
		)
	}

	transferencia, err := s.buscar(ctx, id)

	if err != nil {
		return TransferenciaDTO{}, err
	}

	if transferencia.Estado != EstadoAprobada {
		return TransferenciaDTO{}, errors.New(
			"la transferencia debe estar aprobada para poder completarse",
		)
	}

	if rol != usuario.RolAdministrador &&
		depositoUsuarioID != transferencia.DepositoDestinoID.Hex() {

		return TransferenciaDTO{}, errors.New(
			"solo puede completar transferencias destinadas a su depósito",
		)
	}

	// Acreditamos la mercadería en el destino.
	movimiento := stock.CrearMovimientoRequest{
		ProductoID: transferencia.ProductoID.Hex(),
		DepositoID: transferencia.DepositoDestinoID.Hex(),
		Cantidad:   transferencia.Cantidad,
		Tipo:       string(stock.MovimientoEntrada),
		Motivo: "Recepción de transferencia desde depósito " +
			transferencia.DepositoOrigenID.Hex(),
	}

	if err := s.stockService.CrearMovimiento(
		ctx,
		movimiento,
		usuarioID,
	); err != nil {
		return TransferenciaDTO{}, err
	}

	ahora := time.Now()

	transferencia.Estado = EstadoCompletada
	transferencia.CompletadaPor = usuarioID
	transferencia.FechaCompletada = &ahora

	if err := s.repository.Update(
		ctx,
		transferencia,
	); err != nil {
		return TransferenciaDTO{}, err
	}

	return convertirDTO(transferencia), nil
}

func (s *Service) FindAll(
	ctx context.Context,
) ([]TransferenciaDTO, error) {

	transferencias, err := s.repository.FindAll(ctx)

	if err != nil {
		return nil, err
	}

	resultado := make(
		[]TransferenciaDTO,
		0,
		len(transferencias),
	)

	for _, transferencia := range transferencias {
		resultado = append(
			resultado,
			convertirDTO(transferencia),
		)
	}

	return resultado, nil
}

func (s *Service) buscar(
	ctx context.Context,
	id string,
) (Transferencia, error) {

	objectID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return Transferencia{}, errors.New(
			"id de transferencia inválido",
		)
	}

	transferencia, err := s.repository.FindByID(
		ctx,
		objectID,
	)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Transferencia{}, errors.New(
				"transferencia no encontrada",
			)
		}

		return Transferencia{}, err
	}

	return transferencia, nil
}

func convertirDTO(
	transferencia Transferencia,
) TransferenciaDTO {

	dto := TransferenciaDTO{
		ID:                transferencia.ID.Hex(),
		ProductoID:        transferencia.ProductoID.Hex(),
		DepositoOrigenID:  transferencia.DepositoOrigenID.Hex(),
		DepositoDestinoID: transferencia.DepositoDestinoID.Hex(),
		Cantidad:          transferencia.Cantidad,
		Estado:            transferencia.Estado,
		Motivo:            transferencia.Motivo,
		SolicitadaPor:     transferencia.SolicitadaPor,
		AprobadaPor:       transferencia.AprobadaPor,
		CompletadaPor:     transferencia.CompletadaPor,
		MotivoRechazo:     transferencia.MotivoRechazo,
	}

	if transferencia.FechaSolicitud != nil {
		dto.FechaSolicitud = transferencia.FechaSolicitud.Format(time.RFC3339)
	}

	if transferencia.FechaAprobacion != nil {
		dto.FechaAprobacion = transferencia.FechaAprobacion.Format(time.RFC3339)
	}

	if transferencia.FechaRechazo != nil {
		dto.FechaRechazo = transferencia.FechaRechazo.Format(time.RFC3339)
	}

	if transferencia.FechaCompletada != nil {
		dto.FechaCompletada = transferencia.FechaCompletada.Format(time.RFC3339)
	}

	return dto
}
