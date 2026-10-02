package proveedor

type ProductoProveedorDTO struct {
	ProductoID    string  `json:"producto_id"`
	Costo         float64 `json:"costo"`
	TiempoEntrega int     `json:"tiempo_entrega_dias"`
}

type ProveedorDTO struct {
	ID          string                 `json:"id"`
	RazonSocial string                 `json:"razon_social"`
	CUIT        string                 `json:"cuit"`
	Contacto    string                 `json:"contacto"`
	Productos   []ProductoProveedorDTO `json:"productos,omitempty"`
}

type ProductoProveedorRequest struct {
	ProductoID    string  `json:"producto_id" binding:"required"`
	Costo         float64 `json:"costo"`
	TiempoEntrega int     `json:"tiempo_entrega_dias"`
}

type CrearProveedorRequest struct {
	RazonSocial string `json:"razon_social" binding:"required"`
	CUIT        string `json:"cuit" binding:"required"`
	Contacto    string `json:"contacto"`

	Productos []ProductoProveedorRequest `json:"productos"`
}

type ActualizarProveedorRequest struct {
	RazonSocial string `json:"razon_social" binding:"required"`
	CUIT        string `json:"cuit" binding:"required"`
	Contacto    string `json:"contacto"`

	Productos []ProductoProveedorRequest `json:"productos"`
}
