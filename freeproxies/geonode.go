package main

import ("fmt"
		"net/http"
		"io/ioutil"
		"encoding/json"
		"os"
		"math"
		)
		
func fetchGeonode(protocol string) []Proxy{
	var proxies []Proxy
	
	request, _ := http.NewRequest("GET", fmt.Sprintf("https://proxylist.geonode.com/api/proxy-list?limit=200&page=1&sort_by=lastChecked&sort_type=desc&protocols=%s", protocol), nil)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil{
		fmt.Println("Failed making request")
		os.Exit(0)
	}
	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	
	var proxyData map[string][]map[string]string
	//var protocolData map[string][]map[string][]string
	var total map[string]int
	
	addProxies := func(data []byte){
		
		json.Unmarshal([]byte(data), &proxyData)
		//json.Unmarshal([]byte(data), &protocolData)
		
		for _, proxy := range proxyData["data"]{
	
			var adress = proxy["ip"]+":"+proxy["port"]
			//var protocol = protocolData["data"][i]["protocols"][0]
			var country = proxy["country"]
			
			proxies = append(proxies, Proxy{adress, protocol, country})
		}
	}
	addProxies(data)
	
	json.Unmarshal([]byte(data), &total)
	var pages = int(math.Ceil(float64(total["total"])/float64(200)))
	for i := 2; i<=pages; i++{
		request, _ = http.NewRequest("GET", fmt.Sprintf("https://proxylist.geonode.com/api/proxy-list?limit=200&page=%d&sort_by=lastChecked&sort_type=desc&protocols=%s", i, protocol), nil)
		response, err = client.Do(request)
		if err != nil{
			fmt.Println("Failed making request")
			os.Exit(0)
		}
		data, _ = ioutil.ReadAll(response.Body)
		defer response.Body.Close()
			
		addProxies(data)
	}
	return proxies
}