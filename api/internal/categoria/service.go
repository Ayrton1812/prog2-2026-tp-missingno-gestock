package categoria

import (
	"context"
	"errors"
	"strings"
	"time"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) ListarTodas(ctx context.Context) ([]Categoria, error) {
	return s.repository.FindAll(ctx)
}

func (s *Service) BuscarPorID(ctx context.Context, id string) (Categoria, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Crear(ctx context.Context, categoria Categoria) (Categoria, error) {
	categoria.Nombre = strings.TrimSpace(categoria.Nombre)

	if categoria.Nombre == "" {
		return Categoria{}, errors.New("el nombre es obligatorio")
	}

	ahora := time.Now()

	categoria.CreadoPor = "sistema"
	categoria.ModificadoPor = "sistema"
	categoria.FechaCreacion = ahora
	categoria.FechaActualizacion = ahora

	return s.repository.Create(ctx, categoria)
}

func (s *Service) Actualizar(
	ctx context.Context,
	id string,
	categoria Categoria,
) (Categoria, error) {

	categoria.Nombre = strings.TrimSpace(categoria.Nombre)

	if categoria.Nombre == "" {
		return Categoria{}, errors.New("el nombre es obligatorio")
	}

	categoria.ModificadoPor = "sistema"
	categoria.FechaActualizacion = time.Now()

	return s.repository.Update(ctx, id, categoria)
}

func (s *Service) Eliminar(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
