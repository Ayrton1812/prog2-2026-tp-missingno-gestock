package usuario

type UsuarioDTO struct {
	ID         string `json:"id,omitempty"`
	Nombre     string `json:"nombre"`
	Email      string `json:"email"`
	Rol        string `json:"rol"`
	DepositoID string `json:"deposito_id,omitempty"`
	Activo     bool   `json:"activo"`
}

type CrearUsuarioRequest struct {
	Nombre     string `json:"nombre" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Rol        string `json:"rol" binding:"required"`
	DepositoID string `json:"deposito_id,omitempty"`
}

type ActualizarUsuarioRequest struct {
	Nombre     string `json:"nombre" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Rol        string `json:"rol" binding:"required"`
	DepositoID string `json:"deposito_id,omitempty"`
	Activo     bool   `json:"activo"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token   string     `json:"token"`
	Usuario UsuarioDTO `json:"usuario"`
}

func (u Usuario) ToDTO() UsuarioDTO {
	dto := UsuarioDTO{
		ID:     u.ID.Hex(),
		Nombre: u.Nombre,
		Email:  u.Email,
		Rol:    u.Rol,
		Activo: u.Activo,
	}

	if u.DepositoID != nil {
		dto.DepositoID = u.DepositoID.Hex()
	}

	return dto
}
