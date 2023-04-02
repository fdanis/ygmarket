package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/fdanis/yg-loyalsys/internal/db/entities"
)

type UserRepository struct {
	db         *sql.DB
	insertStmt *sql.Stmt
}

func NewUserRepository(d *sql.DB) UserRepository {
	insertStmt, err := d.PrepareContext(context.TODO(), "insert into public.user (login,password,withdrawn,balance) values ($1,$2,$3,$4) RETURNING id")
	if err != nil {
		log.Fatal("can not create statement for UserRepository")
	}
	return UserRepository{db: d, insertStmt: insertStmt}
}

func (r *UserRepository) GetByLogin(login string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, "select id, login, password,withdrawn,balance FROM public.user where login = $1 limit 1", login)
	m := entities.User{}
	err := row.Scan(&m.ID, &m.Login, &m.Password, &m.Withdrawn, &m.Balance)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println("can not get user by login:", err)
		return nil, err
	}
	return &m, nil
}

func (r *UserRepository) GetByID(id int) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := r.db.QueryRowContext(ctx, "select id, login, password,withdrawn,balance FROM public.user where id = $1 limit 1", id)
	m := entities.User{}
	err := row.Scan(&m.ID, &m.Login, &m.Password, &m.Withdrawn, &m.Balance)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println("can not get user by login:", err)
		return nil, err
	}
	return &m, nil
}

func (r *UserRepository) Add(data *entities.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.insertStmt.QueryRowContext(ctx, data.Login, data.Password, data.Withdrawn, data.Balance).Scan(&data.ID)
	if err != nil {
		return err
	}
	return nil
}
