package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/fdanis/yg-loyalsys/internal/db/entities"
)

type OrderRepository struct {
	db         *sql.DB
	insertStmt *sql.Stmt
}

func NewOrderRepository(d *sql.DB) OrderRepository {
	insertStmt, err := d.PrepareContext(context.TODO(), "insert into public.order (userid,ordernumber,status,accrual) values ($1,$2,$3,$4)")
	if err != nil {
		log.Fatal("can not create statement for OrderRepository")
	}
	return OrderRepository{db: d, insertStmt: insertStmt}
}

func (r *OrderRepository) GetByNumber(number string) (*entities.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, "select id, userid, ordernumber, status, accrual, created FROM public.order where ordernumber = $1 limit 1", number)
	m := entities.Order{}
	err := row.Scan(&m.ID, &m.UserID, &m.Number, &m.Status, &m.Accrual, &m.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println("can not get order by number:", err)
		return nil, err
	}
	return &m, nil
}

func (r *OrderRepository) GetAllByUser(id int) ([]*entities.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row, err := r.db.QueryContext(ctx, "select id, userid, ordernumber, status, accrual, created FROM public.order where userid = $1", id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	res := make([]*entities.Order, 0)
	for row.Next() {
		o := &entities.Order{}
		err = row.Scan(&o.ID, &o.UserID, &o.Number, &o.Status, &o.Accrual, &o.Created)
		if err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	err = row.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *OrderRepository) Add(data entities.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s, err := r.insertStmt.ExecContext(ctx, data.UserID, data.Number, data.Status, data.Accrual)
	if err != nil {
		return err
	}
	i, err := s.RowsAffected()
	if err != nil {
		return err
	}
	if i == 0 {
		return errors.New("order was not save")
	}
	return nil
}
