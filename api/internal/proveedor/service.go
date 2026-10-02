package proveedor

import (
	"context"
	"errors"
	"regexp"
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

func (s *Service) FindAll(
	ctx context.Context,
) ([]ProveedorDTO, error) {

	proveedores, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	resultado := make([]ProveedorDTO, 0, len(proveedores))

	for _, proveedor := range proveedores {
		resultado = append(
			resultado,
			convertirDTO(proveedor),
		)
	}

	return resultado, nil
}

func (s *Service) FindByID(
	ctx context.Context,
	id string,
) (ProveedorDTO, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ProveedorDTO{}, errors.New(
			"id de proveedor inválido",
		)
	}

	proveedor, err := s.repository.FindByID(
		ctx,
		objectID,
	)

	if err != nil {
		return ProveedorDTO{}, err
	}

	return convertirDTO(proveedor), nil
}

func (s *Service) Create(
	ctx context.Context,
	request CrearProveedorRequest,
	usuarioID string,
) (ProveedorDTO, error) {

	if err := validarProveedor(
		request.RazonSocial,
		request.CUIT,
		request.Productos,
	); err != nil {
		return ProveedorDTO{}, err
	}

	productos, err := convertirProductos(
		request.Productos,
	)

	if err != nil {
		return ProveedorDTO{}, err
	}

	now := time.Now()

	proveedor := Proveedor{
		ID:                 bson.NewObjectID(),
		RazonSocial:        strings.TrimSpace(request.RazonSocial),
		CUIT:               limpiarCUIT(request.CUIT),
		Contacto:           strings.TrimSpace(request.Contacto),
		Productos:          productos,
		CreadoPor:          usuarioID,
		ModificadoPor:      usuarioID,
		FechaCreacion:      now,
		FechaActualizacion: now,
	}

	if err := s.repository.Create(
		ctx,
		proveedor,
	); err != nil {
		return ProveedorDTO{}, err
	}

	return convertirDTO(proveedor), nil
}

func (s *Service) Update(
	ctx context.Context,
	id string,
	request ActualizarProveedorRequest,
	usuarioID string,
) (ProveedorDTO, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ProveedorDTO{}, errors.New(
			"id de proveedor inválido",
		)
	}

	if err := validarProveedor(
		request.RazonSocial,
		request.CUIT,
		request.Productos,
	); err != nil {
		return ProveedorDTO{}, err
	}

	productos, err := convertirProductos(
		request.Productos,
	)

	if err != nil {
		return ProveedorDTO{}, err
	}

	proveedor := Proveedor{
		ID:                 objectID,
		RazonSocial:        strings.TrimSpace(request.RazonSocial),
		CUIT:               limpiarCUIT(request.CUIT),
		Contacto:           strings.TrimSpace(request.Contacto),
		Productos:          productos,
		ModificadoPor:      usuarioID,
		FechaActualizacion: time.Now(),
	}

	if err := s.repository.Update(
		ctx,
		proveedor,
	); err != nil {
		return ProveedorDTO{}, err
	}

	return convertirDTO(proveedor), nil
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
) error {

	objectID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return errors.New(
			"id de proveedor inválido",
		)
	}

	return s.repository.Delete(ctx, objectID)
}

func validarProveedor(
	razonSocial string,
	cuit string,
	productos []ProductoProveedorRequest,
) error {

	if strings.TrimSpace(razonSocial) == "" {
		return errors.New(
			"la razón social es obligatoria",
		)
	}

	if !validarCUIT(cuit) {
		return errors.New(
			"CUIT inválido",
		)
	}

	productosUsados := make(map[string]bool)

	for _, producto := range productos {

		if _, err := bson.ObjectIDFromHex(
			producto.ProductoID,
		); err != nil {
			return errors.New(
				"producto_id inválido",
			)
		}

		if producto.Costo < 0 {
			return errors.New(
				"el costo no puede ser negativo",
			)
		}

		if producto.TiempoEntrega < 0 {
			return errors.New(
				"el tiempo de entrega no puede ser negativo",
			)
		}

		if productosUsados[producto.ProductoID] {
			return errors.New(
				"un producto no puede repetirse dentro del proveedor",
			)
		}

		productosUsados[producto.ProductoID] = true
	}

	return nil
}

func convertirProductos(
	productos []ProductoProveedorRequest,
) ([]ProductoProveedor, error) {

	resultado := make(
		[]ProductoProveedor,
		0,
		len(productos),
	)

	for _, producto := range productos {

		objectID, err := bson.ObjectIDFromHex(
			producto.ProductoID,
		)

		if err != nil {
			return nil, errors.New(
				"producto_id inválido",
			)
		}

		resultado = append(
			resultado,
			ProductoProveedor{
				ProductoID:    objectID,
				Costo:         producto.Costo,
				TiempoEntrega: producto.TiempoEntrega,
			},
		)
	}

	return resultado, nil
}

func convertirDTO(
	proveedor Proveedor,
) ProveedorDTO {

	productos := make(
		[]ProductoProveedorDTO,
		0,
		len(proveedor.Productos),
	)

	for _, producto := range proveedor.Productos {

		productos = append(
			productos,
			ProductoProveedorDTO{
				ProductoID:    producto.ProductoID.Hex(),
				Costo:         producto.Costo,
				TiempoEntrega: producto.TiempoEntrega,
			},
		)
	}

	return ProveedorDTO{
		ID:          proveedor.ID.Hex(),
		RazonSocial: proveedor.RazonSocial,
		CUIT:        proveedor.CUIT,
		Contacto:    proveedor.Contacto,
		Productos:   productos,
	}
}

func limpiarCUIT(cuit string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.TrimSpace(cuit),
			"-",
			"",
		),
		" ",
		"",
	)
}

func validarCUIT(cuit string) bool {

	cuit = limpiarCUIT(cuit)

	match, _ := regexp.MatchString(
		`^[0-9]{11}$`,
		cuit,
	)

	if !match {
		return false
	}

	digitos := make([]int, 11)

	for i, caracter := range cuit {
		digitos[i] = int(caracter - '0')
	}

	multiplicadores := []int{
		5, 4, 3, 2, 7, 6, 5, 4, 3, 2,
	}

	suma := 0

	for i := 0; i < 10; i++ {
		suma += digitos[i] * multiplicadores[i]
	}

	resto := suma % 11
	digitoVerificador := 11 - resto

	if digitoVerificador == 11 {
		digitoVerificador = 0
	}

	if digitoVerificador == 10 {
		digitoVerificador = 9
	}

	return digitoVerificador == digitos[10]
}
