package ordencompra

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	EstadoBorrador   = "borrador"
	EstadoConfirmada = "confirmada"
	EstadoParcial    = "parcial"
	EstadoCompletada = "completada"
)

type OrdenCompra struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	DepositoID  bson.ObjectID `bson:"deposito_id"`
	ProveedorID bson.ObjectID `bson:"proveedor_id"`

	Estado string `bson:"estado"`

	Items []ItemOrdenCompra `bson:"items"`

	CreadoPor     string `bson:"creado_por"`
	ModificadoPor string `bson:"modificado_por"`
	ConfirmadoPor string `bson:"confirmado_por,omitempty"`

	FechaCreacion      time.Time  `bson:"fecha_creacion"`
	FechaActualizacion time.Time  `bson:"fecha_actualizacion"`
	FechaConfirmacion  *time.Time `bson:"fecha_confirmacion,omitempty"`
}

type ItemOrdenCompra struct {
	ProductoID bson.ObjectID `bson:"producto_id"`

	CantidadSolicitada int `bson:"cantidad_solicitada"`
	CantidadRecibida   int `bson:"cantidad_recibida"`

	CostoUnitario     float64 `bson:"costo_unitario"`
	TiempoEntregaDias int     `bson:"tiempo_entrega_dias"`
}
