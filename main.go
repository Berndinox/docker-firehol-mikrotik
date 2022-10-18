package main

import (
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"io"
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

	writeIndex()
	downloadFile("/tmp/ip", ipListUrl)
	
	fs := http.FileServer(http.Dir("/tmp/"))
	http.Handle("/", http.StripPrefix("/", fs))

	err := http.ListenAndServe(":" + httpPort, nil)
	if err != nil {
		fmt.Errorf("Webserver startfailed")
	} else {
		fmt.Println("Webserver running on: http://localhost:" + httpPort)
	}


}

func writeIndex() {
    val := "Docker FireHOL Mikrotik"
    data := []byte(val)

    err := ioutil.WriteFile("/tmp/index.html", data, 0)
	if err != nil {
		fmt.Println("Index not created")
	}
    fmt.Println("Index created")
}

func downloadFile(filepath string, url string) (err error) {

  // Create the file
  out, err := os.Create(filepath)
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

  fmt.Println("IPs downloaded")

  return nil
}
