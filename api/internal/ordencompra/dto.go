package ordencompra

type OrdenCompraDTO struct {
	ID          string               `json:"id"`
	DepositoID  string               `json:"deposito_id"`
	ProveedorID string               `json:"proveedor_id"`
	Estado      string               `json:"estado"`
	Items       []ItemOrdenCompraDTO `json:"items"`

	CreadoPor     string `json:"creado_por"`
	ModificadoPor string `json:"modificado_por"`
	ConfirmadoPor string `json:"confirmado_por,omitempty"`

	FechaCreacion      string `json:"fecha_creacion"`
	FechaActualizacion string `json:"fecha_actualizacion"`
	FechaConfirmacion  string `json:"fecha_confirmacion,omitempty"`
}

type ItemOrdenCompraDTO struct {
	ProductoID         string  `json:"producto_id"`
	CantidadSolicitada int     `json:"cantidad_solicitada"`
	CantidadRecibida   int     `json:"cantidad_recibida"`
	CostoUnitario      float64 `json:"costo_unitario"`
	TiempoEntregaDias  int     `json:"tiempo_entrega_dias"`
}

type CrearOrdenCompraRequest struct {
	DepositoID  string                        `json:"deposito_id" binding:"required"`
	ProveedorID string                        `json:"proveedor_id" binding:"required"`
	Items       []CrearItemOrdenCompraRequest `json:"items" binding:"required,min=1"`
}

type CrearItemOrdenCompraRequest struct {
	ProductoID string `json:"producto_id" binding:"required"`
	Cantidad   int    `json:"cantidad" binding:"required,min=1"`
}

type ActualizarOrdenCompraRequest struct {
	Items []CrearItemOrdenCompraRequest `json:"items" binding:"required,min=1"`
}

type RecibirOrdenCompraRequest struct {
	Items []RecibirItemOrdenCompraRequest `json:"items" binding:"required,min=1"`
}

type RecibirItemOrdenCompraRequest struct {
	ProductoID string `json:"producto_id" binding:"required"`
	Cantidad   int    `json:"cantidad" binding:"required,min=1"`
}
