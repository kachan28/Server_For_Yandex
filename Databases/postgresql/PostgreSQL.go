package postgresql

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx"
)

func setConnection() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "postgresql://localhost/users")
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
	} else {
		return conn
	}
	return nil
}

type runQueryParams struct { //Struct for running SQL Queries
	sql    string        //sql string
	params []interface{} //params for sql
}

//NewParams function that returns new parameters object
func NewParams(sql string, params []interface{}) *runQueryParams {
	return &runQueryParams{sql, params}
}

//RunQuery function to run SQL query
func RunQuery(params *runQueryParams) (pgx.Rows, error) {
	var rows pgx.Rows
	var err error
	conn := setConnection()
	defer conn.Close(context.Background())
	if conn == nil {
		log.Println("Empty connection")
		os.Exit(1)
	} else {
		rows, err = conn.Query(context.Background(), params.sql, params.params...)
		if err != nil {
			log.Println(err)
		}
	}
	return rows, err
}
