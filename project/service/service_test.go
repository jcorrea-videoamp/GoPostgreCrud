package service

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/jcorrea-videoamp/GoPostgreCrud/project/proto"
	"github.com/jcorrea-videoamp/GoPostgreCrud/project/repository"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"google.golang.org/protobuf/types/known/timestamppb"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetOrder(t *testing.T) {
	Convey("Fetch order from database", t, func() {
		db, mock, err := sqlxmock.Newx()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		repo, err := repository.NewRepository(db)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		srv, err := NewOrderService(repo)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when creating the service", err)
		}

		Convey("The value should be equal", func() {
			date := time.Now()
			rows := sqlxmock.NewRows([]string{"id", "status", "customer", "quantity", "price", "created_at", "updated_at"}).
				AddRow(1, "payed", "John Doe", 25, 150.34, date, date)
			mock.ExpectQuery(regexp.QuoteMeta(`select * from orders where id=$1`)).WithArgs("1").WillReturnRows(rows)
			resp, err := srv.GetOrder(context.Background(), &proto.GetRequest{Id: int32(1)})
			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &proto.OrderResponse{
				Order: &proto.Order{
					Id:        int32(1),
					Status:    "payed",
					Customer:  "John Doe",
					Quantity:  25,
					Price:     150.34,
					CreatedAt: timestamppb.New(date),
					UpdatedAt: timestamppb.New(date),
				},
			})
		})
	})
}
