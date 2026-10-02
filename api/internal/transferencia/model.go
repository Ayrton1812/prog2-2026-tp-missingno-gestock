package transferencia

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	EstadoSolicitada = "solicitada"
	EstadoAprobada   = "aprobada"
	EstadoRechazada  = "rechazada"
	EstadoCompletada = "completada"
)

type Transferencia struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	ProductoID bson.ObjectID `bson:"producto_id"`

	DepositoOrigenID  bson.ObjectID `bson:"deposito_origen_id"`
	DepositoDestinoID bson.ObjectID `bson:"deposito_destino_id"`

	Cantidad int `bson:"cantidad"`

	Estado string `bson:"estado"`

	Motivo string `bson:"motivo"`

	SolicitadaPor string `bson:"solicitada_por"`

	AprobadaPor string `bson:"aprobada_por,omitempty"`

	CompletadaPor string `bson:"completada_por,omitempty"`

	MotivoRechazo string `bson:"motivo_rechazo,omitempty"`

	FechaSolicitud  *time.Time `bson:"fecha_solicitud"`
	FechaAprobacion *time.Time `bson:"fecha_aprobacion,omitempty"`
	FechaRechazo    *time.Time `bson:"fecha_rechazo,omitempty"`
	FechaCompletada *time.Time `bson:"fecha_completada,omitempty"`
}
