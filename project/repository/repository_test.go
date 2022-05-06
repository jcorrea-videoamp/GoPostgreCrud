package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/jcorrea-videoamp/GoPostgreCrud/project/models"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetOrder(t *testing.T) {
	Convey("Get order from repository", t, func() {
		db, mock, err := sqlxmock.Newx()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		repo, err := NewRepository(db)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		Convey("The value should be equal", func() {
			date := time.Now()
			rows := sqlxmock.NewRows([]string{"id", "status", "customer", "quantity", "price", "created_at", "updated_at"}).
				AddRow(1, "payed", "John Doe", 25, 150.34, date, date)
			mock.ExpectQuery(regexp.QuoteMeta(`select * from orders where id=$1`)).WithArgs("1").WillReturnRows(rows)
			order, err := repo.GetOrder("1")
			So(err, ShouldBeNil)
			So(order, ShouldResemble, &models.Order{
				ID:       1,
				Status:   "payed",
				Customer: "John Doe",
				Quantity: 25,
				Price:    150.34,
				OrderHistory: models.OrderHistory{
					CreatedAt: date,
					UpdatedAt: date,
				},
			})
		})
	})
}
