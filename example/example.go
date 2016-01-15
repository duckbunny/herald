package main

import (
	"log"

	"github.com/duckbunny/herald"
	"github.com/duckbunny/herald/registry"
)

func init() {
	registry.RegisterAll()
}

func main() {
	h, err := herald.This()
	if err != nil {
		log.Fatal(err)
	}
	err = h.Init()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = h.Declare()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = h.StartPool()
	if err != nil {
		log.Fatal(err)
		return
	}

}
