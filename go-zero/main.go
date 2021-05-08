package main

/*
Released under the The MIT License (MIT)
see ./LICENSE
*/

import (
	"log"
)

import (
	"github.com/eyedeekay/zerobundle"
)

func main() {
	if err := zerobundle.ZeroMain(); err != nil {
		log.Println(err)
	}
}
