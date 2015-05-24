package main

import (
	"fmt"
	"net/url"
	"net/http"
	"bufio"
	"strings"
	"strconv"
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

func main() {

//	resp, err := http.Get("http://www.prezydent.pl/robots.txt")
	resp, err := http.Get("http://allegro.pl/robots.txt")
	
	if err != nil {
		panic(err.Error())
	}
	
	defer resp.Body.Close()
	
	r := new(Robots)
	r.NewRobots(GetScannerFromHttpReponse(resp))
	
	fmt.Printf("%+v", r)
	
}
