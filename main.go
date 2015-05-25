package main

import (
	"fmt"
	"net/url"
	"net/http"
	"bufio"
	"strings"
	"strconv"
	"time"
	"sync"
)

type Robots struct {
	UserAgent 	[]string
	Disallow 	[]*url.URL
	Allow 		[]*url.URL
	Sitemap		[]*url.URL
	CrawlDelay	int64
	Comments		[]string
}

const UA_DEF 			= "User-agent:"
const ALLOW_DEF 			= "Allow:"
const DISALLOW_DEF 		= "Disallow:"
const CRAWL_DELAY_DEF	= "Crawl-Delay:"
const SITEMAP_DEF		= "Sitemap:"
const COMMENT_DEF		= "#"

const MAX_WORKERS		= 250

func (robots *Robots) NewRobots(s *bufio.Scanner) {
	
	var check = func(line string, def string) (string, bool){
		if strings.Index(line, def) >= 0 {
			l := strings.Replace(line, def, "", 1)
			return strings.Trim(l, " "), true
		}
		
		return "", false
	}
	
	for s.Scan() {
		line := s.Text()
		
		if v, ok := check(line, UA_DEF); ok {
			robots.UserAgent = append(robots.UserAgent, v)
		}else if v, ok := check(line, SITEMAP_DEF); ok {
			u, _ := url.Parse(v)
			robots.Sitemap = append(robots.Sitemap, u)
		}else if v, ok := check(line, ALLOW_DEF); ok {
			u, _ := url.Parse(v)
			robots.Allow = append(robots.Allow, u)
		}else if v, ok := check(line, DISALLOW_DEF); ok {
			u, _ := url.Parse(v)
			robots.Disallow = append(robots.Disallow, u)
		}else if v, ok := check(line, COMMENT_DEF); ok {
			robots.Comments = append(robots.Comments, v)
		}else if v, ok := check(line, CRAWL_DELAY_DEF); ok {
			robots.CrawlDelay, _ = strconv.ParseInt(v, 0, 0)
		}
	}
}

func GetScannerFromHttpReponse(r *http.Response) *bufio.Scanner {
	scanner := bufio.NewScanner(r.Body)
	scanner.Split(bufio.ScanLines) 
	return scanner
}

func log1(i interface{}){
	fmt.Println(i)
}

func log2(i interface{}){
	fmt.Printf("%+v\n", i)
}

func getRobots(host string) *Robots{
	
	if has := strings.Contains(host, "http://"); !has {
		host = "http://" + host
	}
	
	resp, err := http.Get(host + "/robots.txt")
	
	if err != nil {
		panic(err.Error())
	}
	
	defer resp.Body.Close()

	r := new(Robots)
	r.NewRobots(GetScannerFromHttpReponse(resp))
	log2(r)
	return r
}

var domains = []string{"www.prezydent.pl","allegro.pl","stooq.pl","gazeta.pl","tvn24.pl","tvnwarszwa.pl","pudelek.pl","platforma.org","pis.org.pl","nowaprawicajkm.pl"}

func worker(linkChan chan string, wg *sync.WaitGroup) {
   defer wg.Done()

   for url := range linkChan {
     getRobots(url)
   }
}

func paraler() {
	lCh := make(chan string)
    wg := new(sync.WaitGroup)

    // Adding routines to workgroup and running then
    for i := 0; i < MAX_WORKERS; i++ {
        wg.Add(1)
        go worker(lCh, wg)
    }

    // Processing all links by spreading them to `free` goroutines
    for _, link := range domains {
        lCh <- link
    }

    // Closing channel (waiting in goroutines won't continue any more)
    close(lCh)

    // Waiting for all goroutines to finish (otherwise they die as main routine dies)
    wg.Wait()
}

func main(){
	start := time.Now()
	synchronized()
	elapsed := time.Since(start)
    fmt.Printf("took %s", elapsed)
}

