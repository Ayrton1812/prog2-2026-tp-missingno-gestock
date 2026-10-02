package deposito

type DepositoDTO struct {
	ID            string `json:"id"`
	Nombre        string `json:"nombre"`
	Localidad     string `json:"localidad"`
	Provincia     string `json:"provincia"`
	ResponsableID string `json:"responsable_id,omitempty"`
}

type CrearDepositoRequest struct {
	Nombre        string `json:"nombre" binding:"required"`
	Localidad     string `json:"localidad" binding:"required"`
	Provincia     string `json:"provincia" binding:"required"`
	ResponsableID string `json:"responsable_id"`
}

type ActualizarDepositoRequest struct {
	Nombre        string `json:"nombre" binding:"required"`
	Localidad     string `json:"localidad" binding:"required"`
	Provincia     string `json:"provincia" binding:"required"`
	ResponsableID string `json:"responsable_id"`
}
