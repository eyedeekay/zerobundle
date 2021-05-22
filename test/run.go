//+build run

package main

/*
Released under the The MIT License (MIT)
see ./LICENSE
*/

import (
	. "github.com/eyedeekay/zerobundle"
	"log"
)

func main() {
	/*if err := UnpackZeroJavaHome(); err != nil {
		log.Println(err)
	}*/
	if err := UnpackZero(); err != nil {
		log.Println(err)
	}
	latest := LatestZero()
	log.Println("latest zero version is:", latest)
	if err := RunZero(); err != nil {
		log.Fatal(err)
	}
}
