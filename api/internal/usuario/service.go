package usuario

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	repository Repository
	jwtSecret  string
}

func NewService(
	repository Repository,
	jwtSecret string,
) *Service {
	return &Service{
		repository: repository,
		jwtSecret:  jwtSecret,
	}
}

func (s *Service) ListarTodas(
	ctx context.Context,
) ([]Usuario, error) {
	return s.repository.FindAll(ctx)
}

func (s *Service) BuscarPorID(
	ctx context.Context,
	id string,
) (Usuario, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Crear(
	ctx context.Context,
	request CrearUsuarioRequest,
) (Usuario, error) {

	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Nombre = strings.TrimSpace(request.Nombre)

	if !rolValido(request.Rol) {
		return Usuario{}, errors.New("rol invalido")
	}

	if requiereDeposito(request.Rol) && request.DepositoID == "" {
		return Usuario{}, errors.New(
			"el gerente u operario debe tener un deposito asignado",
		)
	}

	if !requiereDeposito(request.Rol) {
		request.DepositoID = ""
	}

	_, err := s.repository.FindByEmail(
		ctx,
		request.Email,
	)

	if err == nil {
		return Usuario{}, errors.New(
			"ya existe un usuario con ese email",
		)
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return Usuario{}, err
	}

	usuario := Usuario{
		Nombre:        request.Nombre,
		Email:         request.Email,
		PasswordHash:  string(passwordHash),
		Rol:           request.Rol,
		Activo:        true,
		CreadoPor:     "sistema",
		ModificadoPor: "sistema",
	}

	if request.DepositoID != "" {
		depositoID, err := bson.ObjectIDFromHex(
			request.DepositoID,
		)

		if err != nil {
			return Usuario{}, errors.New(
				"deposito_id invalido",
			)
		}

		usuario.DepositoID = &depositoID
	}

	ahora := time.Now()

	usuario.FechaCreacion = ahora
	usuario.FechaActualizacion = ahora

	return s.repository.Create(ctx, usuario)
}

func (s *Service) Actualizar(
	ctx context.Context,
	id string,
	request ActualizarUsuarioRequest,
) (Usuario, error) {

	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Nombre = strings.TrimSpace(request.Nombre)

	if !rolValido(request.Rol) {
		return Usuario{}, errors.New("rol invalido")
	}

	if requiereDeposito(request.Rol) && request.DepositoID == "" {
		return Usuario{}, errors.New(
			"el gerente u operario debe tener un deposito asignado",
		)
	}

	usuario, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return Usuario{}, err
	}

	usuario.Nombre = request.Nombre
	usuario.Email = request.Email
	usuario.Rol = request.Rol
	usuario.Activo = request.Activo

	usuario.DepositoID = nil

	if request.DepositoID != "" {
		depositoID, err := bson.ObjectIDFromHex(
			request.DepositoID,
		)

		if err != nil {
			return Usuario{}, errors.New(
				"deposito_id invalido",
			)
		}

		usuario.DepositoID = &depositoID
	}

	usuario.ModificadoPor = "sistema"
	usuario.FechaActualizacion = time.Now()

	return s.repository.Update(ctx, id, usuario)
}

func (s *Service) Login(
	ctx context.Context,
	request LoginRequest,
) (LoginResponse, error) {

	email := strings.ToLower(strings.TrimSpace(request.Email))

	usuario, err := s.repository.FindByEmail(
		ctx,
		email,
	)

	if err != nil {
		return LoginResponse{}, errors.New(
			"credenciales invalidas",
		)
	}

	if !usuario.Activo {
		return LoginResponse{}, errors.New(
			"el usuario esta desactivado",
		)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(usuario.PasswordHash),
		[]byte(request.Password),
	)

	if err != nil {
		return LoginResponse{}, errors.New(
			"credenciales invalidas",
		)
	}

	claims := jwt.MapClaims{
		"user_id": usuario.ID.Hex(),
		"rol":     usuario.Rol,
	}

	if usuario.DepositoID != nil {
		claims["deposito_id"] = usuario.DepositoID.Hex()
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		[]byte(s.jwtSecret),
	)

	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Token:   tokenString,
		Usuario: usuario.ToDTO(),
	}, nil
}

func rolValido(rol string) bool {
	switch rol {
	case RolAdministrador,
		RolGerenteDeposito,
		RolOperarioDeposito,
		RolAuditor:
		return true
	default:
		return false
	}
}

func requiereDeposito(rol string) bool {
	return rol == RolGerenteDeposito ||
		rol == RolOperarioDeposito
}
