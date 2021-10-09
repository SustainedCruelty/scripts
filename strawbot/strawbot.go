package main

import ("fmt"
		"net/http"
		"net/url"
		"io/ioutil"
		"os"
		"regexp"
		"strings"
		"flag"
		"crypto/tls"
		"time"
		"sync"
		"bufio")
		
func main(){

	var url string
	flag.StringVar(&url, "u", "", "link to the strawpoll")
	
	var option string
	flag.StringVar(&option, "o", "", "what option to vote for")
	
	var concurrency int
	flag.IntVar(&concurrency, "c", 10, "what level of concurrency to use")
	
	var proxies string
	flag.StringVar(&proxies, "p", "proxies.txt", "file containing a list of https proxies")
	
	var timeout int
	flag.IntVar(&timeout, "t", 20, "how many seconds until timeout for each request")
	
	flag.Parse()
	
	ch := make(chan string)
	var wg sync.WaitGroup
	var options = getOptions(url)
	
    f, _ := os.Open(proxies)
	defer f.Close()
	var scanner = bufio.NewScanner(f)
	
	for i := 0; i < concurrency; i++{
		wg.Add(1)
		
		go func(){
			defer wg.Done()
			
			for proxy := range ch{
				makeVote(proxy, url, options[option], timeout)
			}
		}()
	}
	for scanner.Scan(){
		ch <- scanner.Text()
	}
	close(ch)
    wg.Wait()
}

func makeVote(proxy string, straw_url string, oids string, timeout int){

	form := url.Values{}
	form.Add("pid", strings.Split(straw_url, "/")[3])
	form.Add("oids", oids)
	
	var proxyURL = url.URL{Host: proxy}
	
	transport := &http.Transport{Proxy: http.ProxyURL(&proxyURL), TLSClientConfig: &tls.Config{},}
	client := &http.Client{Transport: transport, Timeout: time.Duration(timeout)*time.Second}
	
	request, _ := http.NewRequest("POST", "https://strawpoll.de/vote", strings.NewReader(form.Encode()))
	
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("Referer", straw_url)
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	
	resp, err := client.Do(request)
	if err != nil{
		fmt.Printf("[-] Failed making request with proxy %s\n", proxy)
	}else{
		fmt.Printf("[+] Successfully casted vote with proxy %s\n", proxy)
		defer resp.Body.Close()
	}
}

func getOptions(url string) map[string]string{
	var options map[string]string = make(map[string]string)
	
	request, _ := http.NewRequest("GET", url, nil)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil{
		fmt.Printf("Failed making request: %s", err)
		os.Exit(0)
	}
	defer response.Body.Close()
	
	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)
	
	re := regexp.MustCompile("(?s)<div class=\"checkbox checkbox-danger\">(.*?)</div>")
	checkboxes := re.FindAllStringSubmatch(pageContent, -1)
	if checkboxes == nil{
		fmt.Println("No match")
		os.Exit(0)
	}
	for _, box := range checkboxes{
	
		re = regexp.MustCompile(`name="(.*?)" id="`)
		box_name := re.FindStringSubmatch(box[1])[1]
		
		re = regexp.MustCompile(`id="(.*?)" class`)
		box_id := re.FindStringSubmatch(box[1])[1]
		
		re = regexp.MustCompile(fmt.Sprintf(`for="%s">(?s)(.*?)(?s)</label>`, box_id))
		box_option:= strings.TrimSpace(re.FindStringSubmatch(box[1])[1])
		
		options[box_option] = box_name
	}
	
	return options
}