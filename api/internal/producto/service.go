package producto

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

func (s *Service) FindAll(ctx context.Context) ([]ProductoDTO, error) {
	productos, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	resultado := make([]ProductoDTO, 0, len(productos))

	for _, producto := range productos {
		resultado = append(resultado, convertirDTO(producto))
	}

	return resultado, nil
}

func (s *Service) FindByID(
	ctx context.Context,
	id string,
) (ProductoDTO, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ProductoDTO{}, errors.New("id de producto inválido")
	}

	producto, err := s.repository.FindByID(ctx, objectID)
	if err != nil {
		return ProductoDTO{}, err
	}

	return convertirDTO(producto), nil
}

func (s *Service) Create(
	ctx context.Context,
	request CrearProductoRequest,
	usuarioID string,
) (ProductoDTO, error) {

	if err := validarProducto(
		request.Nombre,
		request.CategoriaID,
		request.UnidadMedida,
		request.CostoUnitario,
		request.StockMinimoGeneral,
		request.StockMinimoPorDeposito,
	); err != nil {
		return ProductoDTO{}, err
	}

	categoriaID, err := bson.ObjectIDFromHex(request.CategoriaID)
	if err != nil {
		return ProductoDTO{}, errors.New("categoria_id inválido")
	}

	fechaVencimiento, err := parseFecha(request.FechaVencimiento)
	if err != nil {
		return ProductoDTO{}, err
	}

	now := time.Now()

	producto := Producto{
		ID:                     bson.NewObjectID(),
		Nombre:                 strings.TrimSpace(request.Nombre),
		Descripcion:            strings.TrimSpace(request.Descripcion),
		CategoriaID:            categoriaID,
		UnidadMedida:           request.UnidadMedida,
		CostoUnitario:          request.CostoUnitario,
		StockMinimoGeneral:     request.StockMinimoGeneral,
		StockMinimoPorDeposito: request.StockMinimoPorDeposito,
		FechaVencimiento:       fechaVencimiento,
		CreadoPor:              usuarioID,
		ModificadoPor:          usuarioID,
		FechaCreacion:          now,
		FechaActualizacion:     now,
	}

	if err := s.repository.Create(ctx, producto); err != nil {
		return ProductoDTO{}, err
	}

	return convertirDTO(producto), nil
}

func (s *Service) Update(
	ctx context.Context,
	id string,
	request ActualizarProductoRequest,
	usuarioID string,
) (ProductoDTO, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ProductoDTO{}, errors.New("id de producto inválido")
	}

	if err := validarProducto(
		request.Nombre,
		request.CategoriaID,
		request.UnidadMedida,
		request.CostoUnitario,
		request.StockMinimoGeneral,
		request.StockMinimoPorDeposito,
	); err != nil {
		return ProductoDTO{}, err
	}

	categoriaID, err := bson.ObjectIDFromHex(request.CategoriaID)
	if err != nil {
		return ProductoDTO{}, errors.New("categoria_id inválido")
	}

	fechaVencimiento, err := parseFecha(request.FechaVencimiento)
	if err != nil {
		return ProductoDTO{}, err
	}

	producto := Producto{
		ID:                     objectID,
		Nombre:                 strings.TrimSpace(request.Nombre),
		Descripcion:            strings.TrimSpace(request.Descripcion),
		CategoriaID:            categoriaID,
		UnidadMedida:           request.UnidadMedida,
		CostoUnitario:          request.CostoUnitario,
		StockMinimoGeneral:     request.StockMinimoGeneral,
		StockMinimoPorDeposito: request.StockMinimoPorDeposito,
		FechaVencimiento:       fechaVencimiento,
		ModificadoPor:          usuarioID,
		FechaActualizacion:     time.Now(),
	}

	if err := s.repository.Update(ctx, producto); err != nil {
		return ProductoDTO{}, err
	}

	return convertirDTO(producto), nil
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
) error {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id de producto inválido")
	}

	return s.repository.Delete(ctx, objectID)
}

func validarProducto(
	nombre string,
	categoriaID string,
	unidad string,
	costo float64,
	stockMinimo int,
	stockMinimoDepositos map[string]int,
) error {

	if strings.TrimSpace(nombre) == "" {
		return errors.New("el nombre es obligatorio")
	}

	if strings.TrimSpace(categoriaID) == "" {
		return errors.New("la categoría es obligatoria")
	}

	if _, err := bson.ObjectIDFromHex(categoriaID); err != nil {
		return errors.New("categoria_id inválido")
	}

	if !UnidadesMedidaValidas[unidad] {
		return errors.New(
			"unidad_medida inválida: debe ser unidad, kilogramo, litro o metro",
		)
	}

	if costo < 0 {
		return errors.New("el costo unitario no puede ser negativo")
	}

	if stockMinimo < 0 {
		return errors.New("el stock mínimo no puede ser negativo")
	}

	for depositoID, minimo := range stockMinimoDepositos {
		if _, err := bson.ObjectIDFromHex(depositoID); err != nil {
			return errors.New("hay un deposito_id inválido en stock_minimo_por_deposito")
		}

		if minimo < 0 {
			return errors.New("los stocks mínimos no pueden ser negativos")
		}
	}

	return nil
}

func parseFecha(fecha *string) (*time.Time, error) {
	if fecha == nil || strings.TrimSpace(*fecha) == "" {
		return nil, nil
	}

	resultado, err := time.Parse("2006-01-02", *fecha)
	if err != nil {
		return nil, errors.New(
			"fecha_vencimiento inválida, usar formato YYYY-MM-DD",
		)
	}

	return &resultado, nil
}

func convertirDTO(producto Producto) ProductoDTO {
	dto := ProductoDTO{
		ID:                     producto.ID.Hex(),
		Nombre:                 producto.Nombre,
		Descripcion:            producto.Descripcion,
		CategoriaID:            producto.CategoriaID.Hex(),
		UnidadMedida:           producto.UnidadMedida,
		CostoUnitario:          producto.CostoUnitario,
		StockMinimoGeneral:     producto.StockMinimoGeneral,
		StockMinimoPorDeposito: producto.StockMinimoPorDeposito,
	}

	if producto.FechaVencimiento != nil {
		fecha := producto.FechaVencimiento.Format("2006-01-02")
		dto.FechaVencimiento = &fecha
	}

	return dto
}
