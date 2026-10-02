package producto

const (
	UnidadUnidad    = "unidad"
	UnidadKilogramo = "kilogramo"
	UnidadLitro     = "litro"
	UnidadMetro     = "metro"
)

var UnidadesMedidaValidas = map[string]bool{
	UnidadUnidad:    true,
	UnidadKilogramo: true,
	UnidadLitro:     true,
	UnidadMetro:     true,
}

type ProductoDTO struct {
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`

	CategoriaID string `json:"categoria_id"`

	UnidadMedida  string  `json:"unidad_medida"`
	CostoUnitario float64 `json:"costo_unitario"`

	StockMinimoGeneral     int            `json:"stock_minimo_general"`
	StockMinimoPorDeposito map[string]int `json:"stock_minimo_por_deposito,omitempty"`

	FechaVencimiento *string `json:"fecha_vencimiento,omitempty"`
}

type CrearProductoRequest struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion"`

	CategoriaID string `json:"categoria_id" binding:"required"`

	UnidadMedida  string  `json:"unidad_medida" binding:"required"`
	CostoUnitario float64 `json:"costo_unitario"`

	StockMinimoGeneral int `json:"stock_minimo_general"`

	StockMinimoPorDeposito map[string]int `json:"stock_minimo_por_deposito"`

	FechaVencimiento *string `json:"fecha_vencimiento"`
}

type ActualizarProductoRequest struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion"`

	CategoriaID string `json:"categoria_id" binding:"required"`

	UnidadMedida  string  `json:"unidad_medida" binding:"required"`
	CostoUnitario float64 `json:"costo_unitario"`

	StockMinimoGeneral int `json:"stock_minimo_general"`

	StockMinimoPorDeposito map[string]int `json:"stock_minimo_por_deposito"`

	FechaVencimiento *string `json:"fecha_vencimiento"`
}
