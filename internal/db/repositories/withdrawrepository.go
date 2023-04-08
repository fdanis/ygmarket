package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/fdanis/yg-loyalsys/internal/db/entities"
)

type WithdrawRepository struct {
	db *sql.DB
}

func NewWithdrawRepository(d *sql.DB) WithdrawRepository {
	return WithdrawRepository{db: d}
}

func (r *WithdrawRepository) GetAllByUser(id int) ([]*entities.Withdraw, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row, err := r.db.QueryContext(ctx, "SELECT id, userid, ordernumber, sum, created FROM public.withdraw WHERE userid = $1 order by created", id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	res := make([]*entities.Withdraw, 0)
	for row.Next() {
		o := &entities.Withdraw{}
		err = row.Scan(&o.ID, &o.UserID, &o.Number, &o.Sum, &o.Created)
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

func (r *WithdrawRepository) Add(data entities.Withdraw) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = r.db.ExecContext(ctx, `
				insert into public.withdraw
			(userid, ordernumber, sum)
			values
			($1,$3,$2);`, data.UserID, data.Sum, data.Number)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	update public.user
	set 
		balance = (SELECT coalesce(sum(accrual),0) FROM public.order where userid = $1) - (SELECT coalesce(sum(sum),0) FROM public.withdraw where userid = $1),
		withdrawn = (SELECT coalesce(sum(sum),0) FROM public.withdraw where userid = $1);
	
	`, data.UserID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
