package main

import (
	"net/http"
	"os"
	"fmt"
	"io"
  "io/ioutil"
  "strings"
  "time"

  "github.com/go-co-op/gocron"
)


func runCronJobs(ipListUrl string) {
  s := gocron.NewScheduler(time.UTC)
 
  s.Every(10).Seconds().Do(func() {
    downloadFile("/tmp/ip", ipListUrl)
    updateFile("/tmp/ip", "/tmp/ip.rsc")
  })
 
  s.StartBlocking()
 }

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
    fmt.Errorf("Webserver start failed")
  }

  runCronJobs(ipListUrl)

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
  
  for i, line := range lines {
    if strings.Contains(line, ".") {
      lines[i] = "add list=firehol-blacklist address=" + lines[i]
    }
  }
  currentTime := time.Now()
  output := strings.Join(lines, "\n")
  output = "/ip firewall address-list\n" + output
  output = "# Generated on: " + currentTime.Format("2006.01.02 15:04:05") +"\n" + output
  err = ioutil.WriteFile(outputFile, []byte(output), 0644)
  if err != nil {
    fmt.Errorf("ERROR")
  }

  fmt.Println("IP.rsc created")
}


