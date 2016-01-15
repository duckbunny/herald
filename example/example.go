package main

import (
	"fmt"

	"github.com/duckbunny/herald/registry"
)

func init() {
	registry.RegisterAll()
}

func main() {
	fmt.Println("started")
}
