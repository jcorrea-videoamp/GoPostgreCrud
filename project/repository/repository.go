package repository

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jcorrea-videoamp/GoPostgreCrud/project/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository interface {
	CreateOrder(*models.Order) (string, error)
	ListOrders() ([]*models.Order, error)
	GetOrder(ID string) (*models.Order, error)
	UpdateOrder(map[string]any, string) (string, error)
	DeleteOrder(ID string) (string, error)
}

type repository struct {
	repo *sqlx.DB
}

func NewRepository(r *sqlx.DB) (*repository, error) {
	if r == nil {
		return &repository{}, fmt.Errorf("a sqlx.DB instance is strictly needed to create a repository")
	}
	return &repository{
		repo: r,
	}, nil
}

func (r *repository) CreateOrder(order *models.Order) (string, error) {
	query := "insert into orders (status, customer, quantity, price, created_at, updated_at) values ($1,$2,$3,$4,$5,$6)"
	ret, err := r.repo.Exec(query, order.Status, order.Customer, order.Quantity, order.Price, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return fmt.Sprintf("insert failed, err:%v\n", err), err
	}
	rows, err := ret.RowsAffected()
	if err != nil {
		return fmt.Sprintf("insert failed, err:%v\n", err), err
	}
	return fmt.Sprintf("Order was successfully created. Rows affected: %#v", rows), nil
}

func (r *repository) ListOrders() ([]*models.Order, error) {
	query := "select * from orders"
	rows, err := r.repo.Queryx(query)
	if err != nil {
		return nil, err
	}
	orders := []*models.Order{}
	for rows.Next() {
		var order = new(models.Order)
		err = rows.StructScan(order)
		if err != nil {
			log.Fatalf("error %#v", err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *repository) GetOrder(ID string) (*models.Order, error) {
	query := "select * from orders where id=$1"
	var order = new(models.Order)
	err := r.repo.QueryRowx(query, ID).StructScan(order)
	if err != nil {
		return order, fmt.Errorf("failed to fectch order from db: %#v", err)
	}
	return order, nil
}

func (r *repository) UpdateOrder(values map[string]any, ID string) (string, error) {
	queryKeys := []string{}
	queryValues := []any{}

	for k, v := range values {
		queryKeys = append(queryKeys, k)
		queryValues = append(queryValues, v)
	}
	var query strings.Builder
	args := []string{}
	query.WriteString("update orders set ")
	for i, v := range queryKeys {
		s := fmt.Sprintf("%s=$%d", v, i+1)
		args = append(args, s)
	}
	args = append(args, fmt.Sprintf("updated_at=$%s", strconv.Itoa(len(queryKeys)+1)))
	queryValues = append(queryValues, time.Now())
	formatArgs := strings.Join(args, ", ")
	query.WriteString(formatArgs)
	formatId := fmt.Sprintf(" where id=$%s", strconv.Itoa(len(queryKeys)+2))
	query.WriteString(formatId)
	queryValues = append(queryValues, ID)

	ret, err := r.repo.Exec(query.String(), queryValues...)
	if err != nil {
		return fmt.Sprintf("update failed, err:%v\n", err), err
	}
	rows, err := ret.RowsAffected()
	if err != nil {
		return fmt.Sprintf("update failed, err:%v\n", err), err
	}
	return fmt.Sprintf("Order was successfully updated. Rows affected: %#v", rows), nil
}

func (r *repository) DeleteOrder(ID string) (string, error) {
	query := "delete from orders where id = $1"
	ret, err := r.repo.Exec(query, ID)
	if err != nil {
		return fmt.Sprintf("delete failed, err:%v\n", err), err
	}
	rows, err := ret.RowsAffected()
	if err != nil {
		return fmt.Sprintf("delete failed, err:%v\n", err), err
	}
	return fmt.Sprintf("Order was successfully deleted. Rows affected: %#v", rows), nil
}

func ConnectDb(driver, url string) (*sqlx.DB, error) {
	//url := "postgres://frtzcnqy:pYvsWxUKNQhG6xtFFqAj6sdTZdoc0lvB@chunee.db.elephantsql.com/frtzcnqy"
	pgUrl, _ := pq.ParseURL(url)
	db, err := sqlx.Connect(driver, pgUrl) // postgres
	if err != nil {
		log.Fatalln("Failed to connect to Db with error:", err)
		return nil, err
	}
	log.Printf("Connected to db: %#v", db)
	return db, nil
	/*repo, err := NewRepository(db)
	if err != nil {
		log.Fatalln(err)
		return
	}*/
	/*order := &models.Order{
		Status:   "payed",
		Customer: "JC",
		Quantity: 1,
		Price:    134.8,
		OrderHistory: models.OrderHistory{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	resp, err := repo.CreateOrder(order)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(resp)*/
	/*resp, err := repo.ListOrders()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(resp))*/
	/*resp, err := repo.GetOrder("1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*resp)*/
	/*changes := map[string]any{
		"customer": "Juan",
		"quantity": 100,
		"price":    1500.25,
	}
	resp, err := repo.UpdateOrder(changes, "1")
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(resp)*/
	/*resp, err := repo.DeleteOrder("1")
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(resp)*/
}

/*
CREATE TABLE orders (
	id SERIAL PRIMARY KEY,
	status TEXT,
	customer TEXT,
	quantity INT,
	price FLOAT,
	created_at TEXT,
	updated_at TEXT
);
*/
