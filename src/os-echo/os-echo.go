package main

import (
	/* Std lib */
	"fmt"
	"log"
	"os"
)

/* Homegrown package(s) */
import "envmap"

func main() {
	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Host: %s\n", host)

	fmt.Println("Environment:")
	for key, val := range envmap.Map() {
		fmt.Printf("  %s=%s\n", key, val)
	}
}
