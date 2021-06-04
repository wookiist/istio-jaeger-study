package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
)

var mainLogger *log.Logger

func init() {
	mainLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

type Employee struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age,omitempty"`
	Email string `json:"email,omitempty"`
}

var (
	EmployeeMap map[string]string
)

func main() {
	e := echo.New()
	c := jaegertracing.New(e, nil)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	defer c.Close()
	e.GET("/employee/:name", EmployeePageHandler)
	mainLogger.Fatal(e.Start(":31116"))
}

// EmployeePageHandler is deprecated
func EmployeePageHandler(e echo.Context) error {
	sp := jaegertracing.CreateChildSpan(e, "employee page handler")
	defer sp.Finish()
	name := e.Param("name")
	ID := ""
	switch name {
	case "alice", "Alice":
		ID = "1"
	case "bob", "Bob":
		ID = "2"
	case "cathy", "Cathy":
		ID = "3"
	case "david", "David":
		ID = "4"
	default:
		return e.String(http.StatusBadRequest, "could not find an employee named "+name)
	}
	headers := make(map[string]string)
	for k, v := range e.Response().Header() {
		fmt.Println(k, v)
		headers[k] = v[0]
	}
	employee, err := GetEmployeeDetail(ID, headers)
	if err != nil {
		return e.String(http.StatusInternalServerError, "internal server error")
	}
	return e.JSONPretty(http.StatusOK, *employee, "  ")
}

func GetEmployeeDetail(ID string, headers map[string]string) (*Employee, error) {
	client := resty.New()
	client.SetHostURL("http://details:31117")
	result := new(Employee)
	fmt.Println("current ID", ID)
	resp, err := client.R().
		SetHeaders(headers).
		SetPathParams(map[string]string{
			"id": ID,
		}).
		SetResult(result).
		Get("/details/{id}")
	fmt.Println(resp)
	fmt.Println(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
