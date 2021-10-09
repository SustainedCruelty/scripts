package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
	"regexp"
)

func fetchFreeProxyList(protocol string) []Proxy{
	var proxies []Proxy

	response, err := http.Get("https://free-proxy-list.net/")
	if err != nil{
		log.Fatal(err)
	}
	defer response.Body.Close()
	
	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)
	re := regexp.MustCompile("(?s)<tr>(.*?)</tr>")
	rows := re.FindAllStringSubmatch(pageContent, -1)
	if rows == nil{
		fmt.Println("No match")
		os.Exit(0)
	}
	for _, r := range rows{
	
		re = regexp.MustCompile("<td>(.*?)</td>")
		proxyInfo := re.FindAllStringSubmatch(r[1], 3)
		re = regexp.MustCompile("class='hx'>(.*?)</td>")
		protoInfo := re.FindAllStringSubmatch(r[1],1)
		
		if proxyInfo != nil && protoInfo != nil{
			
			var ip = proxyInfo[0][1]
			var port = proxyInfo[1][1]
			var adress = ip+":"+port			
			var country = proxyInfo[2][1]
			var proxyProto string
			
			if protoInfo[0][1] == "yes"{
				proxyProto = "https"
			}else{
				proxyProto = "http"
			}
			if proxyProto == protocol{
				proxies = append(proxies, Proxy{adress, protocol, country})
			}
		}
	}
	return proxies
}	