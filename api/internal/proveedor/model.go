package proveedor

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Proveedor struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	RazonSocial string `bson:"razon_social"`
	CUIT        string `bson:"cuit"`
	Contacto    string `bson:"contacto"`

	Productos []ProductoProveedor `bson:"productos,omitempty"`

	CreadoPor          string    `bson:"creado_por"`
	ModificadoPor      string    `bson:"modificado_por"`
	FechaCreacion      time.Time `bson:"fecha_creacion"`
	FechaActualizacion time.Time `bson:"fecha_actualizacion"`
}

type ProductoProveedor struct {
	ProductoID    bson.ObjectID `bson:"producto_id"`
	Costo         float64       `bson:"costo"`
	TiempoEntrega int           `bson:"tiempo_entrega_dias"`
}
