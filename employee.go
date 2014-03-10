package main

import (
	"database/sql"
)

type Employee struct {
	Id        int
	FirstName string
	LastName  string
}

func (e *Employee) Load(tx *sql.Tx) error {
	e.FirstName = "Kalle"
	e.LastName = "Persson"
	return nil
}
