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
	if err := UnpackZeroJavaHome(); err != nil {
		log.Println(err)
	}
	latest := LatestZeroJavaHome()
	log.Println("latest zero version is:", latest)
	if err := RunZeroJavaHome(); err != nil {
		log.Fatal(err)
	}
}
