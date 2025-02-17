package user

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/repository/dto"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) (string, error) {
	_, err := r.db.Exec("INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)", user.ID, user.Username, user.Email, user.Password)
	return user.ID, err
}

func (r *userRepository) GetById(id string) (*entity.User, error) {
	var user dto.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	return &entity.User{ID: user.ID, Username: user.Username, Email: user.Email, Password: user.Password}, err
}
