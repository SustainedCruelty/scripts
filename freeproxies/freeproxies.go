package main

import ("fmt"
		"bufio"
		"os"
		"flag"
		)

type Proxy struct {

	Adress string //IP:Port
	Protocol string //http or https
	Country string
}

func main(){
	var protocol string
	flag.StringVar(&protocol, "p", "http", "whether to fetch http or https proxies")
	
	var outputfile string
	flag.StringVar(&outputfile, "o", "proxies.txt", "what file to save the proxies to")
	
	flag.Parse()
	
	var proxies []Proxy
	
	proxies = append(proxies, fetchFreeProxyList(protocol)...)
	proxies = append(proxies, fetchGeonode(protocol)...)
	proxies = append(proxies, fetchProxyscrape(protocol)...)
	
	file, err := os.OpenFile(outputfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil{
		fmt.Println("Failed creating file")
	}
	
	datawriter := bufio.NewWriter(file)
	
	allKeys := make(map[string]bool) // for removing duplicates
	
	for _, p := range proxies{
		if _, value := allKeys[p.Adress]; !value{
			allKeys[p.Adress] = true
			_, _ = datawriter.WriteString(p.Adress+"\n")
			fmt.Println(p.Adress)
		}
	}
	datawriter.Flush()
	file.Close()
}