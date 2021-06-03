package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	// Enable tracing middleware
	c := jaegertracing.New(e, nil)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	defer c.Close()
	e.GET("/", func(c echo.Context) error {
		// Wrap slowFunc on a new span to trace it's execution passing the function arguments
		jaegertracing.TraceFunction(c, slowFunc, "Test String")
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/child", func(c echo.Context) error {
		// Do something before creating the child span
		time.Sleep(40 * time.Millisecond)
		sp := jaegertracing.CreateChildSpan(c, "Child span for additional processing")
		defer sp.Finish()
		sp.LogEvent("Test log")
		sp.SetBaggageItem("Test baggage", "baggage")
		sp.SetTag("test tag", "Istio Study")
		time.Sleep(100 * time.Millisecond)
		return c.String(http.StatusOK, "Hello, Child!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

// A function to be wrapped. No need to change it's arguments due to tracing
func slowFunc(s string) {
	time.Sleep(200 * time.Millisecond)
	return
}
