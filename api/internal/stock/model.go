package stock

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Stock struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	ProductoID bson.ObjectID `bson:"producto_id"`
	DepositoID bson.ObjectID `bson:"deposito_id"`

	Cantidad int `bson:"cantidad"`

	CreadoPor          string    `bson:"creado_por"`
	ModificadoPor      string    `bson:"modificado_por"`
	FechaCreacion      time.Time `bson:"fecha_creacion"`
	FechaActualizacion time.Time `bson:"fecha_actualizacion"`
}

type TipoMovimiento string

const (
	MovimientoEntrada       TipoMovimiento = "entrada"
	MovimientoSalida        TipoMovimiento = "salida"
	MovimientoAjuste        TipoMovimiento = "ajuste"
	MovimientoTransferencia TipoMovimiento = "transferencia"
)

type Movimiento struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	ProductoID bson.ObjectID `bson:"producto_id"`
	DepositoID bson.ObjectID `bson:"deposito_id"`

	Cantidad int `bson:"cantidad"`

	Tipo TipoMovimiento `bson:"tipo"`

	UsuarioID string `bson:"usuario_id"`

	Motivo string `bson:"motivo"`

	FechaHora time.Time `bson:"fecha_hora"`
}
