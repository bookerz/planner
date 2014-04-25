package systemtests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
	"os"
	"strings"
	"testing"
)

// Break out into separate package to avoid duplication and mismatch
type Employee struct {
	Id        int
	FirstName string
	LastName  string
}

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

func (s *EmployeeSuite) TestEmployeeBasicGET(c *C) {

	client := &http.Client{}

	uri := getBaseURI() + "/employee/42"
	req, err := http.NewRequest("GET", uri, nil)
	resp, err := client.Do(req)

	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusOK, Commentf("Expected status 200 for basic GET employee", 200, resp.StatusCode))

	payload := make(map[string]interface{})

	json.NewDecoder(resp.Body).Decode(&payload)

	c.Assert(payload["Id"], Equals, float64(42))
	c.Assert(payload["FirstName"], Equals, "Firstname_42")
	c.Assert(payload["LastName"], Equals, "Lastname_42")
}

func (s *EmployeeSuite) TestEmployeeBasicGETNoNoEmployee(c *C) {

	client := &http.Client{}

	uri := getBaseURI() + "/employee/43"
	req, err := http.NewRequest("GET", uri, nil)
	resp, err := client.Do(req)

	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusNotFound, Commentf("Expected status 200 for basic GET employee", 200, resp.StatusCode))

	b, err := ioutil.ReadAll(resp.Body)

	c.Assert(err, IsNil)
	c.Assert(strings.Contains(resp.Header["Content-Type"][0], "text/plain"), Equals, true)
	c.Assert(strings.TrimSpace(string(b)), Equals, "Employee not found")
}

func (s *EmployeeSuite) TestEmployeeBasicCreateEmployee(c *C) {

	empl := &Employee{
		Id:        0,
		FirstName: "A first name",
		LastName:  "A Last name",
	}

	body, err := json.Marshal(empl)
	c.Assert(err, IsNil)

	client := &http.Client{}

	uri := getBaseURI() + "/employee"
	req, err := http.NewRequest("POST", uri, bytes.NewReader(body))
	resp, err := client.Do(req)

	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, http.StatusOK, Commentf("Expected status 200 for basic POST to create employee", 200, resp.StatusCode))

	b, err := ioutil.ReadAll(resp.Body)

	c.Assert(err, IsNil)
	c.Assert(strings.Contains(resp.Header["Content-Type"][0], "text/plain"), Equals, true)
	c.Assert(strings.TrimSpace(string(b)), Equals, "Employee not found")
}

// Benchmarks

func (s *EmployeeSuite) BenchmarkEmployeeBasicGET(c *C) {

	c.StopTimer()

	client := &http.Client{}
	uri := getBaseURI() + "/employee/42"
	req, _ := http.NewRequest("GET", uri, nil)

	c.StartTimer()

	for i := 0; i < c.N; i++ {
		resp, err := client.Do(req)
		c.Assert(err, IsNil)
		c.Assert(resp.StatusCode, Equals, http.StatusOK, Commentf("Expected status 200 for basic GET employee", 200, resp.StatusCode))
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}

func getBaseURI() string {
	host := os.Getenv("PLANNER_TEST_HOST")
	if len(host) == 0 {
		return defaultURI
	}
	return host
}
