package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	trader "github.com/mridul-sahu/zerodha-trading"
	kt "github.com/zerodhatech/gokiteconnect"
	ktick "github.com/zerodhatech/gokiteconnect/ticker"
)

func main() {
	apiKey := flag.String("key", "", "API KEY")
	apiSecret := flag.String("secret", "", "API SECRET")
	flag.Parse()

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

	ids := []uint32{7712001, 738561, 7670273, 2748929, 969473, 2912513, 470529, 408065, 3693569, 177665,
		424961, 2730497, 3771393, 340481, 1510401, 2977281, 112129, 4451329, 1213441, 3491073}

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

	trader := trader.NewPaperTrader(instToProcess, broker, feed)
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
