package categoria

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Categoria struct {
	ID                 bson.ObjectID `bson:"_id,omitempty"`
	Nombre             string        `bson:"nombre"`
	CreadoPor          string        `bson:"creado_por"`
	ModificadoPor      string        `bson:"modificado_por"`
	FechaCreacion      time.Time     `bson:"fecha_creacion"`
	FechaActualizacion time.Time     `bson:"fecha_actualizacion"`
}
