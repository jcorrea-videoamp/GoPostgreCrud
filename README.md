# GoPostgreCrud
Simple example of CRUD operations in a PostgreSQL database

For generating .pb.go and grpc.pb.go files using the command line:
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative project/proto/order.proto


