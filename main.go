package main

import (
	"net/http"
	"os"
	"io"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func main() {
    localFileLocation := "ip.txt"

	ipListUrl := os.Getenv("IP_LIST_URL")
	if ipListUrl = "" {
        ipListUrl := "https://iplists.firehol.org/files/firehol_level1.netset"
    }

	downloadFile(localFileLocation, ipListUrl)
	
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
		return c.HTML(http.StatusOK, openFile(localFileLocation))
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func openFile(filePath string) {
    body, err := iotuil.ReadFile(filePath)
    if err != nil {
        return err
    }
    return string(body)
}

func downloadFile(filePath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filePath)
	if err != nil  {
	  return err
	}
	defer out.Close()
  
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
	  return err
	}
	defer resp.Body.Close()
  
	// Check server response
	if resp.StatusCode != http.StatusOK {
	  return fmt.Errorf("bad status: %s", resp.Status)
	}
  
	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
	  return err
	}
  
	return nil
}