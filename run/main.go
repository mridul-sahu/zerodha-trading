package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	trader "github.com/mridul-sahu/zerodha-trading"
	kt "github.com/zerodhatech/gokiteconnect"
	ktick "github.com/zerodhatech/gokiteconnect/ticker"
)

//var ids = []uint32{7712001, 738561, 7670273, 2748929, 969473, 2912513, 470529, 408065, 3693569, 177665,
//	424961, 2730497, 3771393, 340481, 1510401, 2977281, 112129, 4451329, 1213441, 3491073}

var ids = []uint32{7712001}

func mockRun(root string) {
	instumentsFile, err := os.OpenFile(filepath.Join(root, "instruments.csv"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Cannot open file: %v", err)
	}
	var instruments kt.Instruments

	if err := gocsv.UnmarshalFile(instumentsFile, &instruments); err != nil {
		log.Fatalf("Cannot Unmarshall File: %v", err)
	}

	var instToProcess kt.Instruments

	for i := range instruments {
		id := uint32(instruments[i].InstrumentToken)
		for _, b := range ids {
			if b == id {
				instToProcess = append(instToProcess, instruments[i])
				break
			}
		}
	}

	var ticks []ktick.Tick
	for _, id := range ids {
		var ts []ktick.Tick
		data, err := ioutil.ReadFile(filepath.Join(root, strconv.FormatUint(uint64(id), 10)+"_ticks.json"))
		if err != nil {
			log.Fatalf("Cannot read file: %v", err)
		}
		if err := json.Unmarshal(data, &ts); err != nil {
			log.Fatalf("Cannot read ticks: %v", err)
		}
		ticks = append(ticks, ts...)
	}

	sort.Slice(ticks, func(i, j int) bool {
		return ticks[i].Timestamp.Time.Before(ticks[i].Timestamp.Time)
	})

	feed := trader.NewFeed(ids)
	broker := trader.NewPaperBroker(10000)

	trader := trader.NewPaperTrader(instToProcess, broker, feed, trader.NewSuperTrendStrategy)
	trader.StartTrading()

	for _, t := range ticks {
		<-time.After(time.Millisecond * 10)
		feed.OnTick(t)
	}
	trader.End()
	broker.SaveOrdersToFile("Orders.json")
}

func main() {
	apiKey := flag.String("key", "", "API KEY")
	apiSecret := flag.String("secret", "", "API SECRET")
	mockPath := flag.String("mock-path", "", "Mock Data Folder")
	flag.Parse()

	if *mockPath != "" {
		mockRun(*mockPath)
		return
	}

	if *apiKey == "" || *apiSecret == "" {
		log.Fatalln("Could not find a vaid api key or secret")
	}

	kc := kt.New(*apiKey)
	fmt.Println(kc.GetLoginURL())
	var requestToken string

	fmt.Println("Please Enter Request Token")
	fmt.Scan(&requestToken)
	data, err := kc.GenerateSession(requestToken, *apiSecret)
	if err != nil {
		log.Fatalf("Cannot generate session: %v", err)
	}
	kc.SetAccessToken(data.AccessToken)
	ticker := ktick.New(*apiKey, data.AccessToken)

	instruments, err := kc.GetInstruments()
	if err != nil {
		log.Fatalf("Cannot get Instruemnts: %v", err)
	}

	var instToProcess kt.Instruments

	for i := range instruments {
		id := uint32(instruments[i].InstrumentToken)
		for _, b := range ids {
			if b == id {
				instToProcess = append(instToProcess, instruments[i])
				break
			}
		}
	}

	feed := trader.NewFeed(ids)
	broker := trader.NewPaperBroker(10000)

	ticker.OnError(func(err error) {
		log.Println("Ticker Error: ", err)
	})

	ticker.OnConnect(func() {
		fmt.Println("Connected")
		if err := ticker.Subscribe(ids); err != nil {
			fmt.Println("Suscribe Error: ", err)
		}
		ticker.SetMode(ktick.ModeFull, ids)
	})

	ticker.OnReconnect(func(attempt int, delay time.Duration) {
		fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
	})

	ticker.OnNoReconnect(func(attempt int) {
		fmt.Println("Maximum no of reconnect attempt reached: ", attempt)
	})

	ticker.OnTick(feed.OnTick)
	//ticker.OnOrderUpdate()

	ticker.OnClose(func(code int, reason string) {
		fmt.Println("Close: ", code, reason)
	})

	trader := trader.NewPaperTrader(instToProcess, broker, feed, trader.NewSuperTrendStrategy)
	trader.StartTrading()

	go func() {
		ticker.Serve()
	}()

	var command string
	for {
		fmt.Scan(&command)
		if command == "Stop" {
			ticker.Unsubscribe(ids)
			trader.End()
			broker.SaveOrdersToFile("Orders.json")
			return
		} else if command == "Save" {
			broker.SaveOrdersToFile("Orders.json")
		} else {
			log.Println("Unknown Command: ", command)
		}
	}
}
