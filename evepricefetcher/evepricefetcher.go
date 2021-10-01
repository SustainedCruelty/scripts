package main

import (
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os"
        "encoding/json"
        "bufio"
        "time"
        "sync"
        "math"
        "strconv"
        "flag"
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
        SystemID int `json:"system_id"`
        TypeID int `json:"type_id"`
        VolumeRemain int `json:"volume_remain"`
        VolumeTotal int `json:"volume_total"`

}

var prices = make(map[string]float64)
var mutex = &sync.Mutex{}

func main() {

        var concurrency int
        flag.IntVar(&concurrency, "c", 20, "set the concurrency level")

        var inputfile string
        flag.StringVar(&inputfile, "i", "","filename of the input file")

        var outputfile string
        flag.StringVar(&outputfile, "o","prices.json","filename of the output file")
		
		var verbose bool
		flag.BoolVar(&verbose, "v", false, "whether to enable verbose output")
		
		var otype string
		flag.StringVar(&otype, "t", "sell", "whether to fetch buy or sell price")

        flag.Parse()

        start := time.Now()
		
		var scanner *bufio.Scanner
		
		//Check if stdin has data
		stdin := os.Stdin
		stat, _ := stdin.Stat()
		size := stat.Size()
		
		//Read the inputfile from -o if stdin is empty
		if size > 0{
				scanner = bufio.NewScanner(os.Stdin)
		}else if inputfile != ""{
				f,_ := os.Open(inputfile)
				defer f.Close()
				scanner = bufio.NewScanner(f)
		}else{ //No data at all
				fmt.Println("Couldnt read any data from stdin and inputfile (-o) isnt set. Exiting.")
				os.Exit(0)
		}
		

        ch := make(chan string)
        wg := sync.WaitGroup{}
        concurrency = int(math.Min(float64(concurrency),float64(linecount(inputfile))))
		if verbose{
				log.Printf("Using a concurrency of %s",strconv.Itoa(concurrency))
		}
		if otype == "sell"{
			for t := 0; t < concurrency; t++{
                wg.Add(1)
                go min_sell_threading(ch, &wg, verbose)
			}
		}
		if otype == "buy"{
			for t := 0; t < concurrency; t++{
                wg.Add(1)
                go max_buy_threading(ch, &wg, verbose)
        }
		}
        
        for scanner.Scan(){
                ch <- scanner.Text()
        }
        if err := scanner.Err(); err != nil {
                log.Fatal(err)
        }
        close(ch)
        wg.Wait()
        j, _ := json.MarshalIndent(prices,"","    ")
        ioutil.WriteFile(outputfile,j,0644)
        elapsed := time.Since(start)
		if verbose{
				log.Printf("Fetched %s price(s) in %s",strconv.Itoa(linecount(inputfile)),elapsed)
		}
}

func min_sell_threading(ch chan string, wg *sync.WaitGroup, verbose bool){
        for line := range ch{
                min_sell(line, verbose)
        }
        wg.Done()
}

func max_buy_threading(ch chan string, wg *sync.WaitGroup, verbose bool){
        for line := range ch{
                max_buy(line, verbose)
        }
        wg.Done()
}


func min_sell(typeid string, verbose bool){
        response, err := http.Get("https://esi.evetech.net/latest/markets/10000002/orders/?datasource=tranquility&order_type=sell&page=1&type_id="+typeid)
        if err != nil{
                fmt.Print(err.Error())
                os.Exit(1)
		}
        responseData, err := ioutil.ReadAll(response.Body)
		if err != nil{
				fmt.Print(err.Error())
				os.Exit(1)
		}
        orders := make([]Order,0)
        json.Unmarshal(responseData, &orders)
        min := float64(0)
        if len(orders)>0{
                min = orders[0].Price
        }
        for i := 0; i < len(orders); i++ {
                if (orders[i].Price < min){
                        min = orders[i].Price
                }
        }
		if verbose{
				fmt.Println(typeid)
		}
        mutex.Lock()
        prices[typeid] = min
        mutex.Unlock()
}

func max_buy(typeid string, verbose bool){
		response, err := http.Get("https://esi.evetech.net/latest/markets/10000002/orders/?datasource=tranquility&order_type=buy&page=1&type_id="+typeid)
		if err != nil{
				fmt.Print(err.Error())
				os.Exit(1)
		}
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil{
				fmt.Print(err.Error())
				os.Exit(1)
		}
		orders := make([]Order, 0)
		json.Unmarshal(responseData, &orders)
		max := float64(0)
		for i := 0; i<len(orders);i++{
				if (orders[i].Price > max){
						max = orders[i].Price
				}
		}
		if verbose{
				fmt.Println(typeid)
		}
		mutex.Lock()
		prices[typeid] = max
		mutex.Unlock()
}

func linecount(filename string) int {
        f, _:= os.Open(filename)
        scanner := bufio.NewScanner(f)
        count := 0
        for scanner.Scan(){
                count++
        }
        return count

}                