package main

import (
	"log"

	"github.com/jcorrea-videoamp/GoPostgreCrud/bazel-GoPostgreCrud/project/repository"
)

func main() {
	driver := "postgres"
	url := "postgres://frtzcnqy:pYvsWxUKNQhG6xtFFqAj6sdTZdoc0lvB@chunee.db.elephantsql.com/frtzcnqy"
	db, err := repository.ConnectDb(driver, url)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.Close()
	_, err = repository.NewRepository(db)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
