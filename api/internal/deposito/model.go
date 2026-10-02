package deposito

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Deposito struct {
	ID            bson.ObjectID  `bson:"_id,omitempty"`
	Nombre        string         `bson:"nombre"`
	Localidad     string         `bson:"localidad"`
	Provincia     string         `bson:"provincia"`
	ResponsableID *bson.ObjectID `bson:"responsable_id,omitempty"`

	CreadoPor          string    `bson:"creado_por"`
	ModificadoPor      string    `bson:"modificado_por"`
	FechaCreacion      time.Time `bson:"fecha_creacion"`
	FechaActualizacion time.Time `bson:"fecha_actualizacion"`
}
