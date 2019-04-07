package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gocarina/gocsv"

	trader "github.com/mridul-sahu/zerodha-trading"
	kt "github.com/zerodhatech/gokiteconnect"
	ktick "github.com/zerodhatech/gokiteconnect/ticker"
)

const (
	apiKey    string = "my_api_key"
	apiSecret string = "my_api_secret"
)

func writeInstruments(instruments kt.Instruments, filename string) {
	instumentsFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("Cannot open/create file (%s): %v", filename, err)
	}
	defer instumentsFile.Close()

	if err := gocsv.MarshalFile(instruments, instumentsFile); err != nil {
		log.Fatalf("Cannot write instruments to file: %v", err)
	}
}

func main() {
	kc := kt.New(apiKey)
	fmt.Println(kc.GetLoginURL())
	var requestToken string

	fmt.Println("Please Enter Request Token")
	fmt.Scan(&requestToken)
	data, err := kc.GenerateSession(requestToken, apiSecret)
	if err != nil {
		log.Fatalf("Cannot generate session: %v", err)
	}
	kc.SetAccessToken(data.AccessToken)
	ticker := ktick.New(apiKey, data.AccessToken)

	instruments, err := kc.GetInstruments()
	if err != nil {
		log.Fatalf("Cannot get Instruemnts: %v", err)
	}

	writeInstruments(instruments, "instruments.csv")

	var ids []uint32

	for _, inst := range instruments {
		ids = append(ids, uint32(inst.InstrumentToken))
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

	trader := trader.NewPaperTrader(instruments, broker, feed)
	trader.StartTrading()

	go func() {
		ticker.Serve()
	}()

	var command string
	for {
		fmt.Println("Want to quit? (Yes/No)")
		fmt.Scan(&command)
		if command == "Yes" {
			ticker.Unsubscribe(ids)
			trader.End()
			if orders, err := json.Marshal(broker.GetOrders()); err == nil {
				if err := ioutil.WriteFile("Orders.josn", orders, os.ModePerm); err != nil {
					fmt.Println("Error writing orders: ", err)
				}
			}
			return
		}
	}
}
