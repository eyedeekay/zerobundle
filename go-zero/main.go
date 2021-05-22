package main

/*
Released under the The MIT License (MIT)
see ./LICENSE
*/

import (
	"flag"
	"log"
)

import (
	"github.com/eyedeekay/zerobundle"
)

func main() {
	bemysam := flag.Bool("sam", false, "run as a SAM bridge on another router's I2CP port.")
	flag.Parse()
	switch *bemysam {
	case true:
		if err := zerobundle.ZeroAsFreestandingSAM(); err != nil {
			log.Println(err)
		}
	default:
		if err := zerobundle.ZeroMain(); err != nil {
			log.Println(err)
		}
	}
}
