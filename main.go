package main

import (
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"io"
  "bufio"
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

	downloadFile("/tmp/ip", ipListUrl)
	linByLine("/tmp/ip")

	fs := http.FileServer(http.Dir("/tmp/"))
	http.Handle("/", http.StripPrefix("/", fs))

	err := http.ListenAndServe(":" + httpPort, nil)
	if err != nil {
		fmt.Errorf("Webserver startfailed")
	} else {
		fmt.Println("Webserver running on: http://localhost:" + httpPort)
	}
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

function linByLine(filepath string) (err error) {
  file, err := os.Open(filepath)
  if err != nil {
          panic(err)
  }
  defer file.Close()

  reader := bufio.NewReader(file)

  for {
          line, _, err := reader.ReadLine()

          if err == io.EOF {
                  break
          }

          fmt.Printf("%s \n", line)
  }
}