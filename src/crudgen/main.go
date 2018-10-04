package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	table = flag.String("table", "", "table name")
	logic = flag.String("logic", "", "logic name")
)

func init() {
	err := os.Setenv("APPPATH", os.Getenv("GOPATH")+"/src/"+os.Getenv("GOAPP")+"/")
	if err != nil {
		log.Fatal(err)
		panic(1)
		return
	}

	flag.Parse()
	check()
}

func check() {
	if os.Getenv("GOPATH") == "" {
		fmt.Println("SET GOPATH")
		os.Exit(1)
		return
	}
	if os.Getenv("GOAPP") == "" {
		fmt.Println("SET GOAPP")
		os.Exit(1)
		return
	}
	if *logic == "" {
		fmt.Println("Define -logic={logic_name}")
		os.Exit(1)
		return
	}
	if *table == "" {
		fmt.Println("Define -table={table_name}")
		os.Exit(1)
		return
	}
}

func main() {

}
