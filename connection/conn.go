package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)


var Conn *pgx.Conn

func DataBaseConnect(){

	databaseUrl := "postgres://postgres:harlequin01@localhost:5432/b47-s1"

	var err error
	Conn, err = pgx.Connect(context.Background(), databaseUrl)

	if err != nil{
		fmt.Fprintf(os.Stderr,"unable to connect: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to database")
}