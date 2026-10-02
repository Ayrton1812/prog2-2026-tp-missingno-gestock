package usuario

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	RolAdministrador    = "administrador"
	RolGerenteDeposito  = "gerente_deposito"
	RolOperarioDeposito = "operario_deposito"
	RolAuditor          = "auditor"
)

type Usuario struct {
	ID           bson.ObjectID  `bson:"_id,omitempty"`
	Nombre       string         `bson:"nombre"`
	Email        string         `bson:"email"`
	PasswordHash string         `bson:"password_hash"`
	Rol          string         `bson:"rol"`
	DepositoID   *bson.ObjectID `bson:"deposito_id,omitempty"`
	Activo       bool           `bson:"activo"`

	CreadoPor          string    `bson:"creado_por"`
	ModificadoPor      string    `bson:"modificado_por"`
	FechaCreacion      time.Time `bson:"fecha_creacion"`
	FechaActualizacion time.Time `bson:"fecha_actualizacion"`
}
