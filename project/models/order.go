package models

import "time"

type Order struct {
	ID       int     `db:"id"`
	Status   string  `db:"status"`
	Customer string  `db:"customer"`
	Quantity int     `db:"quantity"`
	Price    float64 `db:"price"`
	OrderHistory
}

type OrderHistory struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
