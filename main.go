package main

import (
	"net/http"
	"os"
	"fmt"
	"io"
  "io/ioutil"
  "strings"
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
	updateFile("/tmp/ip", "/tmp/ip.rsc")

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

func RemoveIndex(s []int, index int) []int {
  ret := make([]int, 0)
  ret = append(ret, s[:index]...)
  return append(ret, s[index+1:]...)
}

func updateFile(filepath string, outputFile string) {
  input, err := ioutil.ReadFile(filepath)
  if err != nil {
    fmt.Errorf("ERROR")
  }

  lines := strings.Split(string(input), "\n")
  
  newLength := 0
  for i, line := range lines {
          if !(strings.Contains(line, "#")) {
            lines[newLength] = lines[i]
            newLength++
          }
  }
  lines = lines[:newLength]

  output := strings.Join(lines, "\n")
  err = ioutil.WriteFile(outputFile, []byte(output), 0644)
  if err != nil {
    fmt.Errorf("ERROR")
  }
}