package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/golang/glog"
	"net/http"
	"strconv"
)

type Employee struct {
	Id        int
	FirstName string
	LastName  string
}

type EmployeeList struct {
	Employees []*Employee
}

func (e *Employee) Load(tx Transaction) error {
	err := tx.QueryRow("SELECT first_name,last_name FROM employee WHERE id=$1", e.Id).Scan(&e.FirstName, &e.LastName)

	switch {
	case err == sql.ErrNoRows:
		log.Warningln("[EMPLOYEE]: No user with that ID.")
	case err != nil:
		log.Warningf("[EMPLOYEE]: Unknown error reading from database, error: '%v'", err)
	}

	return err
}

func (e *Employee) Delete(tx Transaction) error {
	result, err := tx.Exec("DELETE FROM employee WHERE id = $1", e.Id)

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unable to delete employee '%v', error: '%v'", e.Id, err)
		return err
	}

	cnt, err := result.RowsAffected()

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unable to get rows affected count, error: '%v'", err)
		return err
	}

	if cnt != 1 {
		log.Warningf("Deleted an invalid number of employees '%v'", cnt)
		return errors.New(fmt.Sprintf("invalid number of employees deleted, '%v'", cnt))
	}

	log.V(INFO).Infof("Deleted '%v' employees", cnt)

	return nil
}

func (e *Employee) Save(tx Transaction) error {
	result, err := tx.Exec("UPDATE employee SET first_name=$1,last_name=$2 WHERE id = $3", e.FirstName, e.LastName, e.Id)

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unable to save employee '%v', error: '%v'", e.Id, err)
		return err
	}

	cnt, err := result.RowsAffected()

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unable to get rows affected count, error: '%v'", err)
		return err
	}

	if cnt != 1 {
		log.Warningf("Updated an invalid number of employees '%v'", cnt)
		return errors.New(fmt.Sprintf("invalid number of employees updated, '%v'", cnt))
	}

	log.V(INFO).Infof("Updated '%v' employees", cnt)

	return nil
}

func (e *Employee) Insert(tx Transaction) error {
	err := tx.QueryRow("INSERT INTO employee (id,first_name,last_name) values (nextval('employee_seq'),$1,$2) returning id", e.FirstName, e.LastName).Scan(&e.Id)

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unable to save employee '%v', error: '%v'", e.Id, err)
		return err
	}

	return nil
}

func (e *EmployeeList) Load(tx Transaction, offset, limit int) error {
	rows, err := tx.Query("SELECT id,first_name,last_name FROM employee offset $1 limit $2", offset, limit)

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unknown error reading from database, error: '%v'", err)
		return err
	}

	for rows.Next() {
		employee := Employee{}
		if err := rows.Scan(&employee.Id, &employee.FirstName, &employee.LastName); err != nil {
			log.Warningf("[EMPLOYEE]: error reading data from database, error: ", err)
		}
		e.Employees = append(e.Employees, &employee)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func EmployeeSaveHandler(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unable to format input parameter. error: '%v'", err)
		http.Error(w, "The id have to be a number", http.StatusBadRequest)
		return err
	}

	e := &Employee{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(e); err != nil {
		log.Warningf("[EMPLOYEE]: Unable to read json, error: '%v'", err)
		http.Error(w, "Bad json, "+err.Error(), http.StatusBadRequest)
		return err
	}

	// Disregard any id that may be in the payload
	e.Id = id

	if err := e.Save(tx); err != nil {
		switch {
		case err == sql.ErrNoRows:
			log.Warningf("[EMPLOYEE]: Unable to find employee with id '%v', error: '%v'", id, err)
			http.Error(w, "Employee not found", http.StatusNotFound)
		default:
			log.Warningf("[EMPLOYEE]: Unable to load data from database, error: '%v'", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	return nil
}

func EmployeeCreateHandler(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	e := &Employee{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(e); err != nil {
		log.Warningf("[EMPLOYEE]: Unable to read json, error: '%v'", err)
		return err
	}

	if err := e.Insert(tx); err != nil {
		log.Warningf("[EMPLOYEE]: Unable to create employee, error: '%v'", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func EmployeeHandler(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unable to format input parameter. error: '%v'", err)
		http.Error(w, "The id have to be a number", http.StatusBadRequest)
		return err
	}

	log.V(INFO).Infof("Employee id '%v'", id)

	e := &Employee{
		Id: id,
	}

	if err := e.Load(tx); err != nil {
		switch {
		case err == sql.ErrNoRows:
			log.Warningf("[EMPLOYEE]: Unable to find employee with id '%v', error: '%v'", id, err)
			http.Error(w, "Employee not found", http.StatusNotFound)
		default:
			log.Warningf("[EMPLOYEE]: Unable to load data from database, error: '%v'", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	b, err := json.Marshal(e)

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unable to marshal json data, error: '%v'", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	return nil
}

func EmployeeDeleteHandler(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unable to format input parameter. error: '%v'", err)
		http.Error(w, "The id have to be a number", http.StatusBadRequest)
		return err
	}

	e := &Employee{
		Id: id,
	}

	if err := e.Delete(tx); err != nil {
		switch {
		case err == sql.ErrNoRows:
			log.Warningf("[EMPLOYEE]: Unable to find employee with id '%v', error: '%v'", id, err)
			http.Error(w, "Employee not found", http.StatusNotFound)
		default:
			log.Warningf("[EMPLOYEE]: Unable to load data from database, error: '%v'", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	return nil
}

func EmployeeListHandler(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	offset := getInt(r, "offset", 0)
	limit := getInt(r, "limit", 10)

	log.V(INFO).Infof("Offset = %v, Limit = %v", offset, limit)

	e := &EmployeeList{}

	if err := e.Load(tx, offset, limit); err != nil {
		log.Warningf("[EMPLOYEE]: Unable to load data from database, error: '%v'", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	b, err := json.Marshal(e)

	if err != nil {
		log.Warningf("[EMPLOYEE]: Unable to marshal json data, error: '%v'", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	return nil
}
