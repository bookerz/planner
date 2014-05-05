package main

import (
	"database/sql"
	"testing"
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	_ "github.com/lib/pq"
	. "launchpad.net/gocheck"
)

func Test(t *testing.T) { TestingT(t) }

type CustomerSuite struct{}

var _ = Suite(&CustomerSuite{})

type testResultCustomer struct {
	lastInsertId func() (int64, error)
	rowsAffected func() (int64, error)
}

func (r testResultCustomer) LastInsertId() (int64, error) {
	return r.lastInsertId()
}

func (r testResultCustomer) RowsAffected() (int64, error) {
	return r.rowsAffected()
}

type testTransaction struct {
	result      testResultCustomer
	execCnt     int
	rollbackCnt int
	queryRowCnt int
	queryCnt    int
	prepareCnt  int
	commitCnt   int
	stmtCnt     int
}

func (m *testTransaction) Commit() error {
	m.commitCnt += 1
	return nil
}

func (m *testTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.execCnt += 1
	return m.result, nil
}

func (m *testTransaction) Prepare(query string) (*sql.Stmt, error) {
	m.prepareCnt += 1
	return nil, nil
}
func (m *testTransaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	m.queryCnt += 1
	return nil, nil
}

func (m *testTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	m.queryRowCnt += 1
	return nil
}

func (m *testTransaction) Rollback() error {
	m.rollbackCnt += 1
	return nil
}

func (m *testTransaction) Stmt(stmt *sql.Stmt) *sql.Stmt {
	m.stmtCnt += 1
	return nil
}

func (s *CustomerSuite) TestCustomerDelete(c *C) {
	e := &Customer{}

	r := testResultCustomer{
		lastInsertId: func() (int64, error) {
			return 1, nil
		},
		rowsAffected: func() (int64, error) {
			return 1, nil
		},
	}

	tx := &testTransaction{result: r}

	err := e.Delete(tx)

	c.Assert(err, IsNil)
	c.Assert(tx.execCnt, Equals, 1)

	c.Assert(tx.queryRowCnt, Equals, 0)
	c.Assert(tx.queryCnt, Equals, 0)
	c.Assert(tx.prepareCnt, Equals, 0)
	c.Assert(tx.stmtCnt, Equals, 0)
	c.Assert(tx.commitCnt, Equals, 0)
	c.Assert(tx.rollbackCnt, Equals, 0)
}

func (s *CustomerSuite) SetUpSuite(c *C) {

}

func (s *CustomerSuite) TearDownSuite(c *C) {

}
