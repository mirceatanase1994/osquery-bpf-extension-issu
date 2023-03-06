package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/osquery/osquery-go"
	"github.com/osquery/osquery-go/plugin/logger"
)

func main() {
	var (
		socket = flag.String("socket", "", "")
		_      = flag.Int("timeout", 0, "")
		_      = flag.Int("interval", 0, "")
		_      = flag.Bool("verbose", true, "")
	)
	flag.Parse()
	if *socket == "" {
		log.Fatalf(`Usage: %s --socket SOCKET_PATH`, os.Args[0])
	}

	server, err := osquery.NewExtensionManagerServer("example_logger", *socket)
	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}
	// create and register the plugin
	server.RegisterPlugin(logger.NewPlugin("example_logger", LogString))
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}

func LogString(ctx context.Context, typ logger.LogType, logText string) error {
	log.Printf("%s: %s\n", typ, logText)
	data := url.Values{
		"log": {logText},
	}

	resp, err := http.PostForm("http://localhost:80", data)

	if err != nil {
		fmt.Printf("Got error for log %s: %s\n", logText, err)
		return nil
	}

	fmt.Printf("Got resp for log %s: %s\n", logText, resp)
	return nil
}
