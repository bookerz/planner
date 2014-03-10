package systemtests

import (
	"testing"
	//"bytes"
	"database/sql"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	_ "github.com/lib/pq"
	. "launchpad.net/gocheck"
	"net/http"
	"os"
)

func Test(t *testing.T) { TestingT(t) }

type EmployeeSuite struct{}

var _ = Suite(&EmployeeSuite{})

func (s *EmployeeSuite) SetUpSuite(c *C) {
	id := 42
	firstName := fmt.Sprintf("Firstname_%v", id)
	lastName := fmt.Sprintf("Lastname_%v", id)

	db, err := sql.Open("postgres", fmt.Sprintf("user=%v sslmode=disable", test_db_user))

	if err != nil {
		c.Fatalf("Unable to connect to test database, error: %v", err)
	}

	_, err = db.Exec("INSERT INTO employee (id,first_name,last_name) VALUES ($1,$2,$3)", id, firstName, lastName)

	if err != nil {
		c.Fatalf("Unable to insert database test emplyee, error: %v", err)
	}
}

func (s *EmployeeSuite) TearDownSuite(c *C) {
	id := 42

	db, err := sql.Open("postgres", fmt.Sprintf("user=%v sslmode=disable", test_db_user))

	if err != nil {
		c.Fatalf("Unable to connect to test database, error: %v", err)
	}

	_, err = db.Exec("DELETE FROM employee WHERE id = $1", id)

	if err != nil {
		c.Fatalf("Unable to insert database test emplyee, error: %v", err)
	}
}

func (s *EmployeeSuite) TestEployeeBasicGET(c *C) {

	client := &http.Client{}

	uri := getBaseURI() + "/employee/42"
	req, err := http.NewRequest("GET", uri, nil)
	resp, err := client.Do(req)

	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, 200, Commentf("Expected status 200 for basic GET employee", 200, resp.StatusCode))
}

func getBaseURI() string {
	host := os.Getenv("PLANNER_TEST_HOST")
	if len(host) == 0 {
		return defaultURI
	}
	return host
}
