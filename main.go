package main

import (
	"net/http"
	"os"
	"fmt"
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
	readLine("/tmp/ip")

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

func readLine(filename string) {
  f, err := os.Open(filename)
  if err != nil {
      fmt.Println(err)
      return
  }
  defer f.Close()
  r := bufio.NewReaderSize(f, 4*1024)
  line, isPrefix, err := r.ReadLine()
  for err == nil && !isPrefix {
      s := string(line)
      fmt.Println(s)
      line, isPrefix, err = r.ReadLine()
  }
  if isPrefix {
      fmt.Println("buffer size to small")
      return
  }
  if err != io.EOF {
      fmt.Println(err)
      return
  }
}