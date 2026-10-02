package categoria

import "go.mongodb.org/mongo-driver/v2/bson"

type CategoriaDTO struct {
	ID     string `json:"id,omitempty"`
	Nombre string `json:"nombre" binding:"required"`
}

func (c Categoria) ToDTO() CategoriaDTO {
	return CategoriaDTO{
		ID:     c.ID.Hex(),
		Nombre: c.Nombre,
	}
}

func (dto CategoriaDTO) ToModel() (Categoria, error) {
	categoria := Categoria{
		Nombre: dto.Nombre,
	}

	if dto.ID == "" {
		return categoria, nil
	}

	id, err := bson.ObjectIDFromHex(dto.ID)
	if err != nil {
		return Categoria{}, err
	}

	categoria.ID = id

	return categoria, nil
}
