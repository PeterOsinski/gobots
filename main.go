package main

import (
	"net/http"
	"strings"
	"time"
	"bufio"
	"os"
)

func GetScannerFromHttpReponse(r *http.Response) *bufio.Scanner {
	scanner := bufio.NewScanner(r.Body)
	scanner.Split(bufio.ScanLines) 
	return scanner
}

func getRobots(host string) *Robots{
	
	if has := strings.Contains(host, "http://"); !has {
		host = "http://" + host
	}
	
	resp, err := http.Get(host + "/robots.txt")
	
	if err != nil {
		panic(err.Error())
	}
	
	Logger.Print(host, http.DetectContentType(bufio.NewScanner(resp.Body).Bytes()))
	
	defer resp.Body.Close()

	r := new(Robots)
	r.Address = host
	r.NewRobots(GetScannerFromHttpReponse(resp))
	return r
}

func getHosts(filepath string) []string {
	file, err := os.Open(filepath)
	
	if err != nil {
		panic(err.Error())
	}
	
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines) 
	
	var lines []string
	
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	return lines
}

func main(){
	
	cf := parseCli()
	
	domains := getHosts(cf.file)
	
	start := time.Now()
	
	master(domains, getRobots, cf.workerNum)
	
	elapsed := time.Since(start)
	Logger.Printf("took %s to process %d addresses\n", elapsed, len(domains))
	
}

