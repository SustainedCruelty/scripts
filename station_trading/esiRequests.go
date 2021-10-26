//go build -buildmode=c-shared -o Libraries/esiRequests.dll esiRequests.go

package main

import (
   "fmt"
   "io/ioutil"
   "net/http"
   "strings"
   "net/url"
   "os"
   "encoding/json"
   "bufio"
   "sync"
   "strconv"
   "C"
   "path/filepath"
   //"crypto/tls"
)

type Order struct {
   Duration int `json:"duration"`
   Is_Buy bool `json:"is_buy_order"`
   Issued string `json:"issued"`
   LocationID int `json:"location_id"`
   Min_Volume int `json:"min_volume"`
   OrderID int `json:"order_id"`
   Price float64 `json:"price"`
   Range string `json:"range"`
   TypeID int `json:"type_id"`
   VolumeRemain int `json:"volume_remain"`
   VolumeTotal int `json:"volume_total"`
}

type MarketDay struct{
   Average float64
   Date string
   Highest float64
   Lowest float64
   OrderCount int
   Volume int
}

type MarketItem struct {
   TypeID int
   MaxBuy float64 
   MinSell float64
   BuyQuant int
   SellQuant int
   BuyOrders int
   SellOrders int
   AvgVolume float64
   Bought int
   Sold int
}

func main(){
}

//export pullMarketData
func pullMarketData(inputfilePtr *C.char, outputfilePtr *C.char, concurrency int, region_id int){
   inputfile := filepath.FromSlash(C.GoString(inputfilePtr))
   outputfile := filepath.FromSlash(C.GoString(outputfilePtr))

   var items []MarketItem

   ch := make(chan int)
	var wg sync.WaitGroup

   f, _ := os.Open(inputfile)
   defer f.Close()

   scanner := bufio.NewScanner(f)

   for i := 0; i < concurrency; i++{
      wg.Add(1)

      go func(){
         defer wg.Done()

         for type_id := range ch{
            item := getMarketItem(type_id, region_id)
            items = append(items, item)
         }
      }()
   }

   for scanner.Scan(){
      id, _ := strconv.Atoi(scanner.Text())
      ch <- id
   }
   close(ch)
   wg.Wait()
   j, _ := json.MarshalIndent(items, "", "   ")
   ioutil.WriteFile(outputfile, j, 0644)
}

func getMarketItem(type_id int, region_id int) MarketItem{

   // Get price information

   fmt.Printf("%d\n", type_id)
   response, err := http.Get(fmt.Sprintf("https://esi.evetech.net/latest/markets/%s/orders/?datasource=tranquility&order_type=all&page=1&type_id=%s", strconv.Itoa(region_id), strconv.Itoa(type_id)))
   if err != nil{
      fmt.Println("Failed making request")
      os.Exit(1)
   }
   responseData, err := ioutil.ReadAll(response.Body)
   if err != nil{
      fmt.Println("Failed reading response data")
      os.Exit(1)
   }
   defer response.Body.Close()

   orders := make([]Order, 0)
   json.Unmarshal(responseData, &orders)

   minSell := float64(0)
   maxBuy := float64(0)
   sellQuant := 0
   buyQuant := 0
   sellOrders := 0
   buyOrders := 0
   bought := 0
   sold := 0

   if len(orders) > 0{
      minSell = orders[0].Price
   }

   for i := 0; i < len(orders); i++{

      if (orders[i].Is_Buy){
         buyQuant += orders[i].VolumeRemain
         buyOrders += 1
         bought += (orders[i].VolumeTotal-orders[i].VolumeRemain)
         if (orders[i].Price > maxBuy){
            maxBuy = orders[i].Price
         }
      }else{
         sellQuant += orders[i].VolumeRemain
         sellOrders += 1
         sold += (orders[i].VolumeTotal-orders[i].VolumeRemain)
         if (orders[i].Price < minSell){
            minSell = orders[i].Price
         }
      }
   }

   // Get volume information

   response, err = http.Get(fmt.Sprintf("https://esi.evetech.net/latest/markets/%s/history/?datasource=tranquility&type_id=%s", strconv.Itoa(region_id), strconv.Itoa(type_id)))
   if err != nil{
      fmt.Println("Failed making request")
      os.Exit(1)
   }
   responseData, err = ioutil.ReadAll(response.Body)
   if err != nil{
      fmt.Println("Failed reading response data")
      os.Exit(1)
   }
   defer response.Body.Close()

   history := make([]MarketDay, 0)
   json.Unmarshal(responseData, &history)

   totalVol := 0

   for i := 0; i < len(history); i++{
      totalVol += history[i].Volume
   }
   avgVol := float64(0)

   if len(history) > 0{
      avgVol = float64(totalVol) / float64(400)
   }
   

   return MarketItem{
      TypeID: type_id,
      MaxBuy: maxBuy,
      MinSell: minSell,
      BuyQuant: buyQuant,
      SellQuant: sellQuant,
      BuyOrders: buyOrders,
      SellOrders: sellOrders,
      AvgVolume: avgVol,
      Bought: bought,
      Sold: sold,
   }
}

//export pullStructureOrders
func pullStructureOrders(refresh_tokenPtr *C.char, client_idPtr *C.char, structure_id int, concurrency int, outputfilePtr *C.char){
   refresh_token := C.GoString(refresh_tokenPtr)
   client_id := C.GoString(client_idPtr)
   outputfile := filepath.FromSlash(C.GoString(outputfilePtr))

   fmt.Printf("CLIENT_ID: '%s'\n", client_id)
   fmt.Printf("REFRESH_TOKEN: '%s'\n", refresh_token)

	// Used later for making concurrent requests

	ch := make(chan int)
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	// Fetch an access token for the authenticated ESI request

	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refresh_token) 
	form.Add("client_id", client_id) 

   /*var proxyURL = url.URL{Host: "127.0.0.1:8080"}
   transport := &http.Transport{Proxy: http.ProxyURL(&proxyURL), TLSClientConfig: &tls.Config{},}
	client := &http.Client{Transport: transport}*/

	client := &http.Client{}
	request, _ := http.NewRequest("POST", "https://login.eveonline.com/v2/oauth/token", strings.NewReader(form.Encode()))

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Host", "login.eveonline.com")

	response, err := client.Do(request)
	if err != nil{
		fmt.Println("Error making request")
		os.Exit(1)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var jsonResp map[string]string

	json.Unmarshal([]byte(data), &jsonResp)

	access_token := jsonResp["access_token"]

   fmt.Printf("ACCESS_TOKEN: '%s'\n", access_token)

	// Make an authenticated ESI request with the access token

	request, _ = http.NewRequest("GET", fmt.Sprintf("https://esi.evetech.net/latest/markets/structures/%d/?datasource=tranquility&page=1", structure_id), nil)

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))

	response, err = client.Do(request)
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil{
		fmt.Println("Failed reading response data")
		os.Exit(1)
	}
	defer response.Body.Close()

   fmt.Println("HTTP Response Status:", response.StatusCode, http.StatusText(response.StatusCode))

	var allOrders []Order

	currOrders := make([]Order, 0)
	json.Unmarshal(responseData, &currOrders)
	allOrders = append(allOrders, currOrders...)

	pages, _ := strconv.Atoi(response.Header.Get("X-Pages"))
	fmt.Printf("Fetching a total of %d pages\n", pages)

	for i := 0; i < concurrency; i++{
		
		wg.Add(1)

		go func(){
			defer wg.Done()

			for page := range ch{

				req, _ := http.NewRequest("GET", fmt.Sprintf("https://esi.evetech.net/latest/markets/structures/%d/?datasource=tranquility&page=%d", structure_id, page), nil)
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
				
				resp, err := client.Do(req)
				respData, err := ioutil.ReadAll(resp.Body)
				if err != nil{
					fmt.Println("Failed reading response data")
					os.Exit(1)
				}
				defer resp.Body.Close()

				cOrders := make([]Order, 0)
				json.Unmarshal(respData, &cOrders)

				mutex.Lock()

				allOrders = append(allOrders, cOrders...)
				fmt.Printf("Fetched %d orders from page %d\n", len(cOrders), page)

				mutex.Unlock()

			}
		}()
	}

	for i := 2; i <= pages; i++{
		ch <- i
	}
	close(ch)
	wg.Wait()

	j, _ := json.MarshalIndent(allOrders, "", "   ")
    ioutil.WriteFile(outputfile, j, 0644)

	fmt.Printf("Fetched a total of %d orders", len(allOrders))
}
