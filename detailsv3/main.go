package main

import (
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"strconv"
)

var mainLogger *log.Logger

func init() {
	mainLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

type Employee struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
	Email string `json:"email"`
}
var (
	Alice = Employee{
		ID:    1,
		Name:  "Alice",
		Age:   25,
		Email: "alice@alice.com",
	}
	Bob = Employee{
		ID: 2,
		Name: "Bob",
		Age: 30,
		Email: "bob@bob.com",
	}
	Cathy = Employee{
		ID: 3,
		Name: "Cathy",
		Age: 35,
		Email: "cathy@cathy.com",
	}
	David = Employee{
		ID:4,
		Name: "David",
		Age: 40,
		Email: "david@david.com",
	}

	EmployeeMap map[int]Employee
)

func init() {
	EmployeeMap = make(map[int]Employee)
	EmployeeMap[1] = Alice
	EmployeeMap[2] = Bob
	EmployeeMap[3] = Cathy
	EmployeeMap[4] = David
}

func main() {
	e := echo.New()
	c := jaegertracing.New(e, nil)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	defer c.Close()
	e.GET("/details/:id",  DetailHandler)
	mainLogger.Fatal(e.Start(":31117"))
}

func DetailHandler(e echo.Context) error {
	sp := jaegertracing.CreateChildSpan(e, "detail handler")
	defer sp.Finish()
	ID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid ID value")
	}
	if _, ok := EmployeeMap[ID]; !ok {
		return e.String(http.StatusNotFound, "ID " + string(ID) + " not found")
	}
	return e.JSONPretty(http.StatusOK, EmployeeMap[ID], "  ")
}