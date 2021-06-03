package main
import (
    "github.com/labstack/echo-contrib/jaegertracing"
    "github.com/labstack/echo/v4"
    "net/http"
    "time"
)
func main() {
    e := echo.New()
    // Enable tracing middleware
    c := jaegertracing.New(e, nil)
    defer c.Close()
    e.GET("/", func(c echo.Context) error {
        // Wrap slowFunc on a new span to trace it's execution passing the function arguments
		jaegertracing.TraceFunction(c, slowFunc, "Test String")
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":1323"))
}

// A function to be wrapped. No need to change it's arguments due to tracing
func slowFunc(s string) {
	time.Sleep(200 * time.Millisecond)
	return
}
