package main

import (
	"os"
	"github.com/s-gheldd/DatenkommunikationPraktikum4/routing"
	"fmt"
)

func main() {

	net, err := routing.ParseFile()
	if err != nil {
		os.Exit(1)
	}
	count := 0
	printRoutingTables(count, net)
	for routing.DistanceVectorRoutingStep(net) {

		count++

		printRoutingTables(count, net)
	}

	printRoutingTables(count, net)
}

func printRoutingTables(count int, net map[rune]routing.Router) {

	fmt.Printf("Round %d:\n", count)
	for i := 'A'; i <= 'M'; i++ {
		fmt.Println(net[i])
	}
	fmt.Printf("\n\n")
}
