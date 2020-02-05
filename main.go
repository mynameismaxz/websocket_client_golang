package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

// var addr = flag.String("addr", "api.debug.rtt.in.th:11055", "http service address")

func main() {
	// load environment data
	loadEnvironment()

	// get parameter name
	wsURL := os.Getenv("HOSTNAME")
	protocolWs := os.Getenv("PROTOCOL")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: protocolWs, Host: wsURL}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	for {
		select {
		case <-interrupt:
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Fatal("Error to close message")
				return
			}
			defer c.Close()
			return
		}
	}
}

func loadEnvironment() {
	log.Println("Checking env...")
	requiredVariable := []struct {
		name         string
		defaultValue string
	}{
		{"HOSTNAME", ""},
		{"PROTOCOL", ""},
	}
	for _, env := range requiredVariable {
		value, set := os.LookupEnv(env.name)
		if !set {
			log.Fatal("Require parameter: \"%s\"", env.name)
		}
		log.Printf("%s: %s", env.name, value)
	}
}
