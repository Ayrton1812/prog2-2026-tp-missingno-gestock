package stock

type StockDTO struct {
	ID string `json:"id"`

	ProductoID string `json:"producto_id"`
	DepositoID string `json:"deposito_id"`

	Cantidad int `json:"cantidad"`
}

type MovimientoDTO struct {
	ID string `json:"id"`

	ProductoID string `json:"producto_id"`
	DepositoID string `json:"deposito_id"`

	Cantidad int `json:"cantidad"`

	Tipo string `json:"tipo"`

	UsuarioID string `json:"usuario_id"`

	Motivo string `json:"motivo"`

	FechaHora string `json:"fecha_hora"`
}

type CrearMovimientoRequest struct {
	ProductoID string `json:"producto_id" binding:"required"`
	DepositoID string `json:"deposito_id" binding:"required"`

	Cantidad int `json:"cantidad" binding:"required"`

	Tipo string `json:"tipo" binding:"required"`

	Motivo string `json:"motivo" binding:"required"`
}
