package main

import (
	"net/http"
	"os"
	"fmt"
	"io"
  "io/ioutil"
  "strings"
  "time"
)

// Static Vars
var webroot = "/tmp/"
var rawIPFile = webroot + "ip"
var modIPFile = webroot + "ip.rsc"
var refreshInt int = 30

// Get Settings from ENV
var ipListUrl = os.Getenv("IP_LIST_URL")
// TODO: GetENV for INT is fuck* broken
// var refreshInt = os.Getenv("IP_REFRESH_INTERVAL")
var httpPort = os.Getenv("HTTP_PORT")

// Scheduler for IP Liste create & update
func schedule(f func(), interval time.Duration) *time.Ticker {
  ticker := time.NewTicker(interval)
  go func() {
      for range ticker.C {
          f()
      }
  }()
  return ticker
}

// Main
func main() {

	// IP List to load default
	if ipListUrl == "" {
    ipListUrl = "https://iplists.firehol.org/files/firehol_level1.netset"
  }

  // Refresh Intervall in s default (30)
  // Not Needed until GETENV is fixed
  // if refreshInt == "" {
  //   refreshInt = 30
  // }

	// Port the Service should listen to default 8080
	if httpPort == "" {
		httpPort = "8080"
	}
  
  // Create Timespan for Schedule
  refreshTime := time.Second*time.Duration(refreshInt)

  // Initial get and modifie the list
  createFiles()

  // Schedule the Job
  schedule(createFiles, refreshTime)

  // Init Web-Fileserver
  fs := http.FileServer(http.Dir(webroot))
  http.Handle("/", http.StripPrefix("/", fs))

  err := http.ListenAndServe(":" + httpPort, nil)
  if err != nil {
    fmt.Errorf("Webserver start failed")
  }

}

// Call download & modify Function
func createFiles() {
  downloadFile(rawIPFile, ipListUrl)
  updateFile(rawIPFile, modIPFile)
}

// Download IP-List file
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

// modify the File for Mikrotik (.rsc)
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
