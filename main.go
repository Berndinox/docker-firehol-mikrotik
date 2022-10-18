package main

import (
	"net/http"
	"os"
	"fmt"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func main() {

	// IP List to Block
	ipListUrl := os.Getenv("IP_LIST_URL")
	if ipListUrl == "" {
        ipListUrl = "https://iplists.firehol.org/files/firehol_level1.netset"
    }

	// Port the Service should listen to default 8080
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// Get the data
	var client http.Client
	resp, err := client.Get(ipListUrl)
	if err != nil {
		fmt.Printf("Cant Get IP List via HTTP")
	}
	defer resp.Body.Close()
  
	var bodyString string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Cant read resp Body")
		}
		bodyString = string(bodyBytes)
		fmt.Printf(bodyString)
	} else {
		bodyString = "Something went wrong"
	}



	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Docker FireHOL Mikrotik Converter WEB")
	})

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.GET("/ip", func(c echo.Context) error {
		return c.HTML(http.StatusOK, bodyString)
	})

	e.Logger.Fatal(e.Start(":" + httpPort))
}

