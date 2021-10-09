package main

import ("fmt"
		"net/http"
		"io/ioutil"
		"os"
		"bufio"
		"strings"
		)

func fetchProxyscrape(protocol string) []Proxy{
	var proxies []Proxy
	var ssl string
	
	if protocol == "http"{
		ssl = "no"
	}else if protocol == "https"{
		ssl = "yes"
	}
	request, _ := http.NewRequest("GET", fmt.Sprintf("https://api.proxyscrape.com/v2/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=%s&anonymity=all&simplified=true", ssl), nil)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil{
		fmt.Println("Failed making request")
		os.Exit(0)
	}
	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan(){
		proxies = append(proxies, Proxy{scanner.Text(), protocol, "N/A"})
	}
	return proxies
}