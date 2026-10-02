package transferencia

type CrearTransferenciaRequest struct {
	ProductoID string `json:"producto_id" binding:"required"`

	DepositoOrigenID  string `json:"deposito_origen_id" binding:"required"`
	DepositoDestinoID string `json:"deposito_destino_id" binding:"required"`

	Cantidad int `json:"cantidad" binding:"required"`

	Motivo string `json:"motivo" binding:"required"`
}

type RechazarTransferenciaRequest struct {
	Motivo string `json:"motivo" binding:"required"`
}

type TransferenciaDTO struct {
	ID string `json:"id"`

	ProductoID string `json:"producto_id"`

	DepositoOrigenID  string `json:"deposito_origen_id"`
	DepositoDestinoID string `json:"deposito_destino_id"`

	Cantidad int `json:"cantidad"`

	Estado string `json:"estado"`

	Motivo string `json:"motivo"`

	SolicitadaPor string `json:"solicitada_por"`

	AprobadaPor string `json:"aprobada_por,omitempty"`

	CompletadaPor string `json:"completada_por,omitempty"`

	MotivoRechazo string `json:"motivo_rechazo,omitempty"`

	FechaSolicitud string `json:"fecha_solicitud"`

	FechaAprobacion string `json:"fecha_aprobacion,omitempty"`

	FechaRechazo string `json:"fecha_rechazo,omitempty"`

	FechaCompletada string `json:"fecha_completada,omitempty"`
}
