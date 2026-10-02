package producto

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Producto struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Nombre      string        `bson:"nombre"`
	Descripcion string        `bson:"descripcion"`

	CategoriaID bson.ObjectID `bson:"categoria_id"`

	UnidadMedida  string  `bson:"unidad_medida"`
	CostoUnitario float64 `bson:"costo_unitario"`

	StockMinimoGeneral int `bson:"stock_minimo_general"`

	// Permite sobrescribir el mínimo general para determinados depósitos.
	// La clave es el ID del depósito en formato string.
	StockMinimoPorDeposito map[string]int `bson:"stock_minimo_por_deposito,omitempty"`

	FechaVencimiento *time.Time `bson:"fecha_vencimiento,omitempty"`

	CreadoPor          string    `bson:"creado_por"`
	ModificadoPor      string    `bson:"modificado_por"`
	FechaCreacion      time.Time `bson:"fecha_creacion"`
	FechaActualizacion time.Time `bson:"fecha_actualizacion"`
}
