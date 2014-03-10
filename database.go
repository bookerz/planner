package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Must adapt to HttpRouter library?
func RunInTransaction(f func(w http.ResponseWriter, r *http.Request, tx *sql.Tx) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		tx, err := db.Begin()

		if err != nil {
			log.Printf("[TRANSACTION]: Unable to create database transaction. error: '%v'", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = f(w, r, tx)

		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				log.Printf("TRANSACTION: Unable to rollback transaction, error: {%v}", err2)
			}
			log.Printf("TRANSACTION: Rolling back transaction, error: {%v}", err)
			return
		}

		err = tx.Commit()

		if err != nil {
			log.Printf("TRANSACTION: Unable to commit transaction, error: {%v}", err)
		}
	}
}
