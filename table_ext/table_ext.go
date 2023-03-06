package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/osquery/osquery-go"
	"github.com/osquery/osquery-go/plugin/table"
)

func main() {
	var (
		socket  = flag.String("socket", "", "")
		timeout = flag.Int("timeout", 0, "")
		_       = flag.Int("interval", 0, "")
		_       = flag.Bool("verbose", true, "")
	)
	flag.Parse()
	if *socket == "" {
		log.Fatalf(`Usage: %s --socket SOCKET_PATH`, os.Args[0])
	}

	server, err := osquery.NewExtensionManagerServer("foobar", *socket, osquery.ServerTimeout(time.Duration(*timeout)*time.Second))
	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	// Create and register a new table plugin with the server.
	// table.NewPlugin requires the table plugin name,
	// a slice of Columns and a Generate function.
	server.RegisterPlugin(table.NewPlugin("foobar", FoobarColumns(), FoobarGenerate))
	server.RegisterPlugin(table.NewPlugin("foobar2", Foobar2Columns(), FoobarGenerate))
	server.RegisterPlugin(table.NewPlugin("foobar3", Foobar3Columns(), FoobarGenerate))
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}

// FoobarColumns returns the columns that our table will return.
func FoobarColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("foo"),
		table.TextColumn("baz"),
		table.TextColumn("timestamp"),
	}
}

func Foobar2Columns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("foo2"),
		table.TextColumn("baz"),
		table.TextColumn("timestamp"),
	}
}

func Foobar3Columns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("foo3"),
		table.TextColumn("baz"),
		table.TextColumn("timestamp"),
	}
}

// FoobarGenerate will be called whenever the table is queried. It should return
// a full table scan.
func FoobarGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	ret := []map[string]string{}
	size := rand.Intn(10)
	for i := 0; i < size; i++ {
		ret = append(ret, map[string]string{
			"foo":       "bar",
			"baz":       "baz",
			"timestamp": time.Now().UTC().GoString(),
		})
	}

	return ret, nil
}

func Foobar2Generate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	ret := []map[string]string{}
	size := rand.Intn(10)
	for i := 0; i < size; i++ {
		ret = append(ret, map[string]string{
			"foo2":      "bar",
			"baz":       "baz",
			"timestamp": time.Now().UTC().GoString(),
		})
	}

	return ret, nil
}

func Foobar3Generate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	ret := []map[string]string{}
	size := rand.Intn(100)
	for i := 0; i < size; i++ {
		ret = append(ret, map[string]string{
			"foo3":      "bar",
			"baz":       "baz",
			"timestamp": time.Now().UTC().GoString(),
		})
	}

	return ret, nil
}
