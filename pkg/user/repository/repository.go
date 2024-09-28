package repository

import (
	"database/sql"
	"log"
	"oirakif/simple-blog-service/pkg/user/model"

	"github.com/Masterminds/squirrel"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) InsertUser(newUser model.User) (newID int, err error) {
	builder := squirrel.
		Insert("users").
		Columns(
			"name",
			"email",
			"password_hash",
			"created_at",
			"updated_at",
		).
		Values(
			newUser.Name,
			newUser.Email,
			newUser.PasswordHash,
			newUser.CreatedAt,
			newUser.UpdatedAt,
		)

	sql, args, err := builder.ToSql()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	result, err := r.db.Exec(sql, args...) // ? = placeholder
	if err != nil {
		log.Println(err)
		return 0, err
	}

	newIDInt64, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return int(newIDInt64), nil
}

func (r *UserRepository) FindUser(filterQuery model.FindUserFilterQuery) (retrievedUser *model.User, err error) {
	builder := squirrel.
		Select("id,name,email,created_at,updated_at").
		From("users")

	if filterQuery.Email != nil {
		builder = builder.Where("email=?", filterQuery.Email)
	}
	if filterQuery.PasswordHash != nil {
		builder = builder.Where("password_hash=?", filterQuery.PasswordHash)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var user model.User
	err = r.db.QueryRow(query, args...).
		Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}

	return &user, nil
}
