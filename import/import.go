// +build !debian

package zerobundle

/*
Released under the The MIT License (MIT)
see ./LICENSE
*/

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mholt/archiver/v3"
)

var ZERO_VERSION = "v1.20"

func Write() error {
	var platform = "linux"
	if runtime.GOOS == "windows" {
		platform = "win"
	}
	if runtime.GOOS == "darwin" {
		platform = "mac"
	}
	bytes, err := TBZBytes()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("i2p-zero-"+platform+"."+ZERO_VERSION+".zip", bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func FileNotFound(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return true
	}
	return false
}

func Unpack(destinationDirectory string) error {
	var platform = "linux"
	var platform2 = "linux"
	if runtime.GOOS == "windows" {
		platform = "win-gui"
		platform2 = "win"
	}
	if runtime.GOOS == "darwin" {
		platform = "mac"
		platform2 = "mac"
	}
	if destinationDirectory == "" {
		destinationDirectory = "."
	}
	if FileNotFound(filepath.Join(destinationDirectory, "i2p-zero-"+platform+"."+ZERO_VERSION)) {
		err := Write()
		if err != nil {
			return err
		}
		err = archiver.Unarchive("i2p-zero-"+platform2+"."+ZERO_VERSION+".zip", destinationDirectory)
		if err != nil {
			return err
		}
	}
	return nil
}
