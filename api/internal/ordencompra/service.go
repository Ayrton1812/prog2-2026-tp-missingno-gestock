package ordencompra

import (
	"context"
	"errors"
	"strings"
	"time"

	"gestock/api/internal/producto"
	"gestock/api/internal/proveedor"
	"gestock/api/internal/stock"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	repository          Repository
	proveedorRepository proveedor.Repository
	productoRepository  producto.Repository
	stockService        *stock.Service
}

func NewService(
	repository Repository,
	proveedorRepository proveedor.Repository,
	productoRepository producto.Repository,
	stockService *stock.Service,
) *Service {
	return &Service{
		repository:          repository,
		proveedorRepository: proveedorRepository,
		productoRepository:  productoRepository,
		stockService:        stockService,
	}
}

func (s *Service) Crear(
	ctx context.Context,
	req CrearOrdenCompraRequest,
	usuarioID string,
	rol string,
) (*OrdenCompra, error) {

	if rol != "administrador" && rol != "gerente_deposito" {
		return nil, errors.New("solo un administrador o gerente puede crear órdenes de compra")
	}

	depositoID, err := bson.ObjectIDFromHex(req.DepositoID)
	if err != nil {
		return nil, errors.New("deposito_id inválido")
	}

	proveedorID, err := bson.ObjectIDFromHex(req.ProveedorID)
	if err != nil {
		return nil, errors.New("proveedor_id inválido")
	}

	if len(req.Items) == 0 {
		return nil, errors.New("la orden debe tener al menos un producto")
	}

	prov, err := s.proveedorRepository.BuscarPorID(ctx, proveedorID)
	if err != nil {
		return nil, errors.New("proveedor no encontrado")
	}

	items := make([]ItemOrdenCompra, 0, len(req.Items))
	productosUsados := make(map[bson.ObjectID]bool)

	for _, itemReq := range req.Items {
		if itemReq.Cantidad <= 0 {
			return nil, errors.New("la cantidad debe ser mayor a cero")
		}

		productoID, err := bson.ObjectIDFromHex(itemReq.ProductoID)
		if err != nil {
			return nil, errors.New("producto_id inválido")
		}

		if productosUsados[productoID] {
			return nil, errors.New("no se puede repetir un producto en la orden")
		}

		productosUsados[productoID] = true

		_, err = s.productoRepository.BuscarPorID(ctx, productoID)
		if err != nil {
			return nil, errors.New("producto no encontrado")
		}

		var relacion *proveedor.ProductoProveedor

		for i := range prov.Productos {
			if prov.Productos[i].ProductoID == productoID {
				relacion = &prov.Productos[i]
				break
			}
		}

		if relacion == nil {
			return nil, errors.New(
				"el proveedor no ofrece uno de los productos seleccionados",
			)
		}

		items = append(items, ItemOrdenCompra{
			ProductoID:         productoID,
			CantidadSolicitada: itemReq.Cantidad,
			CantidadRecibida:   0,
			CostoUnitario:      relacion.Costo,
			TiempoEntregaDias:  relacion.TiempoEntrega,
		})
	}

	ahora := time.Now()

	orden := &OrdenCompra{
		ID:          bson.NewObjectID(),
		DepositoID:  depositoID,
		ProveedorID: proveedorID,
		Estado:      EstadoBorrador,
		Items:       items,

		CreadoPor:     usuarioID,
		ModificadoPor: usuarioID,

		FechaCreacion:      ahora,
		FechaActualizacion: ahora,
	}

	if err := s.repository.Crear(ctx, orden); err != nil {
		return nil, err
	}

	return orden, nil
}

func (s *Service) Actualizar(
	ctx context.Context,
	id string,
	req ActualizarOrdenCompraRequest,
	usuarioID string,
) error {

	ordenID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id de orden inválido")
	}

	orden, err := s.repository.BuscarPorID(ctx, ordenID)
	if err != nil {
		return errors.New("orden de compra no encontrada")
	}

	if orden.Estado != EstadoBorrador {
		return errors.New(
			"solo se puede modificar una orden en estado borrador",
		)
	}

	if len(req.Items) == 0 {
		return errors.New("la orden debe tener al menos un producto")
	}

	prov, err := s.proveedorRepository.BuscarPorID(
		ctx,
		orden.ProveedorID,
	)
	if err != nil {
		return errors.New("proveedor no encontrado")
	}

	items := make([]ItemOrdenCompra, 0, len(req.Items))
	productosUsados := make(map[bson.ObjectID]bool)

	for _, itemReq := range req.Items {
		productoID, err := bson.ObjectIDFromHex(itemReq.ProductoID)
		if err != nil {
			return errors.New("producto_id inválido")
		}

		if itemReq.Cantidad <= 0 {
			return errors.New("la cantidad debe ser mayor a cero")
		}

		if productosUsados[productoID] {
			return errors.New("no se puede repetir un producto")
		}

		productosUsados[productoID] = true

		_, err = s.productoRepository.BuscarPorID(
			ctx,
			productoID,
		)
		if err != nil {
			return errors.New("producto no encontrado")
		}

		var relacion *proveedor.ProductoProveedor

		for i := range prov.Productos {
			if prov.Productos[i].ProductoID == productoID {
				relacion = &prov.Productos[i]
				break
			}
		}

		if relacion == nil {
			return errors.New(
				"el proveedor no ofrece uno de los productos",
			)
		}

		items = append(items, ItemOrdenCompra{
			ProductoID:         productoID,
			CantidadSolicitada: itemReq.Cantidad,
			CantidadRecibida:   0,
			CostoUnitario:      relacion.Costo,
			TiempoEntregaDias:  relacion.TiempoEntrega,
		})
	}

	return s.repository.Actualizar(
		ctx,
		ordenID,
		bson.M{
			"items":               items,
			"modificado_por":      usuarioID,
			"fecha_actualizacion": time.Now(),
		},
	)
}

func (s *Service) Confirmar(
	ctx context.Context,
	id string,
	usuarioID string,
) error {

	ordenID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id de orden inválido")
	}

	orden, err := s.repository.BuscarPorID(ctx, ordenID)
	if err != nil {
		return errors.New("orden de compra no encontrada")
	}

	if orden.Estado != EstadoBorrador {
		return errors.New(
			"solo se puede confirmar una orden en estado borrador",
		)
	}

	ahora := time.Now()

	return s.repository.Actualizar(
		ctx,
		ordenID,
		bson.M{
			"estado":              EstadoConfirmada,
			"confirmado_por":      usuarioID,
			"fecha_confirmacion":  ahora,
			"modificado_por":      usuarioID,
			"fecha_actualizacion": ahora,
		},
	)
}

func (s *Service) Recibir(
	ctx context.Context,
	id string,
	req RecibirOrdenCompraRequest,
	usuarioID string,
) error {

	ordenID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id de orden inválido")
	}

	orden, err := s.repository.BuscarPorID(ctx, ordenID)
	if err != nil {
		return errors.New("orden de compra no encontrada")
	}

	if orden.Estado != EstadoConfirmada &&
		orden.Estado != EstadoParcial {
		return errors.New(
			"la orden no está disponible para recibir mercadería",
		)
	}

	for _, recibido := range req.Items {
		productoID, err := bson.ObjectIDFromHex(
			recibido.ProductoID,
		)

		if err != nil {
			return errors.New("producto_id inválido")
		}

		if recibido.Cantidad <= 0 {
			return errors.New("la cantidad recibida debe ser mayor a cero")
		}

		encontrado := false

		for i := range orden.Items {
			item := &orden.Items[i]

			if item.ProductoID != productoID {
				continue
			}

			encontrado = true

			restante := item.CantidadSolicitada -
				item.CantidadRecibida

			if recibido.Cantidad > restante {
				return errors.New(
					"la cantidad recibida supera la cantidad solicitada",
				)
			}

			if err := s.stockService.CrearMovimiento(
				ctx,
				stock.CrearMovimientoRequest{
					ProductoID: productoID.Hex(),
					DepositoID: orden.DepositoID.Hex(),
					Cantidad:   recibido.Cantidad,
					Tipo:       string(stock.MovimientoEntrada),
					Motivo:     "Recepción de orden de compra",
				},
				usuarioID,
			); err != nil {
				return err
			}

			item.CantidadRecibida += recibido.Cantidad
			break
		}

		if !encontrado {
			return errors.New(
				"el producto no pertenece a la orden",
			)
		}
	}

	completa := true
	huboRecepcion := false

	for _, item := range orden.Items {
		if item.CantidadRecibida < item.CantidadSolicitada {
			completa = false
		}

		if item.CantidadRecibida > 0 {
			huboRecepcion = true
		}
	}

	nuevoEstado := EstadoConfirmada

	if completa {
		nuevoEstado = EstadoCompletada
	} else if huboRecepcion {
		nuevoEstado = EstadoParcial
	}

	return s.repository.Actualizar(
		ctx,
		ordenID,
		bson.M{
			"items":               orden.Items,
			"estado":              nuevoEstado,
			"modificado_por":      usuarioID,
			"fecha_actualizacion": time.Now(),
		},
	)
}

func (s *Service) BuscarPorID(
	ctx context.Context,
	id string,
) (*OrdenCompra, error) {

	ordenID, err := bson.ObjectIDFromHex(strings.TrimSpace(id))
	if err != nil {
		return nil, errors.New("id inválido")
	}

	return s.repository.BuscarPorID(ctx, ordenID)
}

func (s *Service) Listar(
	ctx context.Context,
) ([]OrdenCompra, error) {
	return s.repository.Listar(ctx)
}
