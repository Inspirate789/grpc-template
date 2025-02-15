package repository

import "github.com/Inspirate789/grpc-template/internal/models"

type UserDTO struct {
	ID   uint64 `db:"id"`
	Name string `db:"name"`
}

func (dto UserDTO) ToModel() models.User {
	return models.User{
		ID:   dto.ID,
		Name: dto.Name,
	}
}

type CountedUserDTO struct {
	ID         uint64 `db:"id"`
	Name       string `db:"name"`
	TotalCount uint64 `db:"total_count"`
}

func (dto CountedUserDTO) ToModel() models.User {
	return models.User{
		ID:   dto.ID,
		Name: dto.Name,
	}
}

type UsersDTO []CountedUserDTO

func (dto UsersDTO) ToModel() ([]models.User, uint64) {
	res := make([]models.User, 0, len(dto))

	for _, user := range dto {
		res = append(res, user.ToModel())
	}

	if len(dto) != 0 {
		return res, dto[0].TotalCount
	}

	return res, 0
}
