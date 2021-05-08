//+build generate

package main

/*
Released under the The MIT License (MIT)
see ./LICENSE
*/

import (
	"github.com/zserge/lorca"

	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var chunkNum = 128

var mod = `module github.com/eyedeekay/zerobundle/parts/REPLACEME

go 1.14`

var zeroversion = "v1.20"

var unpacker = `package REPLACEME

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func userFind() string {
	if os.Geteuid() == 0 {
		log.Fatal("Do not run this application as root!")
	}
	if un, err := os.UserHomeDir(); err == nil {
		os.MkdirAll(filepath.Join(un, "i2p"), 0755)
		return un
	}
	return ""
}

var userdir = filepath.Join(userFind(), "/i2p/opt/i2p-zero")

func writeFile(val os.FileInfo, system *fs) ([]byte, error) {
	if !val.IsDir() {
		file, err := system.Open(val.Name())
		if err != nil {
			return nil, err
		}
		sys := bytes.NewBuffer(nil)
		if _, err := io.Copy(sys, file); err != nil {
			return nil, err
		} else {
			return sys.Bytes(), nil
		}
	} else {
		log.Println(filepath.Join(userdir, val.Name()), "ignored", "contents", val.Sys())
	}
	return nil, fmt.Errorf("undefined unpacker error")
}

func WriteBrowser(FS *fs) ([]byte, error) {
	if embedded, err := FS.Readdir(-1); err != nil {
		log.Fatal("Extension error, embedded extension not read.", err)
	} else {
		for _, val := range embedded {
			if val.IsDir() {
				os.MkdirAll(filepath.Join(userdir, val.Name()), val.Mode())
			} else {
				return writeFile(val, FS)
			}
		}
	}
	return nil, nil
}
`

func main() {
	// You can also run "npm build" or webpack here, or compress assets, or
	// generate manifests, or do other preparations for your assets.
	if err := Download(); err != nil {
		log.Fatal(err)
	}
	if err := generateGoGenerator("linux"); err != nil {
		log.Fatal(err)
	}
	if err := generateGoGenerator("windows"); err != nil {
		log.Fatal(err)
	}
	if err := generateGoGenerator("darwin"); err != nil {
		log.Fatal(err)
	}
	if err := deleteDirectories(); err != nil {
		log.Fatal(err)
	}
	if err := createDirectories(); err != nil {
		log.Fatal(err)
	}
	if err := generateGoUnpacker(); err != nil {
		log.Fatal(err)
	}
	if err := generateGoMod(); err != nil {
		log.Fatal(err)
	}
	if err := splitBinaries("i2p-zero-linux." + zeroversion + ".zip"); err != nil {
		log.Fatal(err)
	}
	if err := updateAllChunks("linux", "i2p-zero-linux."+zeroversion+".zip"); err != nil {
		log.Fatal(err)
	}
	if err := splitBinaries("i2p-zero-win." + zeroversion + ".zip"); err != nil {
		log.Fatal(err)
	}
	if err := updateAllChunks("windows", "i2p-zero-win."+zeroversion+".zip"); err != nil {
		log.Fatal(err)
	}
	if err := splitBinaries("i2p-zero-darwin." + zeroversion + ".zip"); err != nil {
		log.Fatal(err)
	}
	if err := updateAllChunks("darwin", "i2p-zero-darwin."+zeroversion+".zip"); err != nil {
		log.Fatal(err)
	}

}

var libs = calculateChunks()

func updateChunk(chunk, platform, file string) error {
	err := lorca.Embed("iz"+chunk, "parts/"+chunk+"/chunk_"+platform+".go", file+"."+chunk)
	if err != nil {
		return err
	}
	log.Println("embedded iz" + chunk)
	return nil
}

func updateAllChunks(platform, file string) error {
	for _, lib := range libs {
		updateChunk(lib, platform, file)
	}
	return nil
}

func calculateChunks() []string {
	/*fileToBeChunked := GS_VERSION
	bytes, err := ioutil.ReadFile(fileToBeChunked)
	if err != nil {
		log.Fatal(err)
	}*/
	var libs []string
	for i := 0; i < chunkNum; i++ {
		libs = append(libs, strconv.Itoa(i))
	}
	return libs
}

func chunkSize(fileToBeChunked string) int {
	//	fileToBeChunked := GS_VERSION
	bytes, err := ioutil.ReadFile(fileToBeChunked)
	if err != nil {
		log.Fatal(err)
	}
	chunkSize := len(bytes) / chunkNum
	return chunkSize
}

func splitBinaries(fileToBeChunked string) error {
	bytes, err := ioutil.ReadFile(fileToBeChunked)
	if err != nil {
		return err
	}
	chunkSize := chunkSize(fileToBeChunked)
	for index, lib := range libs {
		start := index * chunkSize
		finish := ((index + 1) * chunkSize)
		if index == chunkNum-1 {
			finish = len(bytes)
		}
		outBytes := bytes[start:finish]
		err := ioutil.WriteFile(fileToBeChunked+"."+lib, outBytes, 0644)
		if err != nil {
			return err
		}
		log.Printf("Started at: %d,  Ended at: %d", start, finish)
	}
	return nil
}

func deleteDirectories() error {
	for _, dir := range libs {
		err := os.RemoveAll(filepath.Join("parts", dir))
		if err != nil {
			return err
		}
	}
	return nil
}

func createDirectories() error {
	for _, dir := range libs {
		err := os.MkdirAll(filepath.Join("parts", dir), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateGoUnpacker() error {
	for index, dir := range libs {
		contents := strings.Replace(unpacker, "REPLACEME", "iz"+libs[index], -1)
		if err := ioutil.WriteFile(filepath.Join("parts", dir, "unpacker.go"), []byte(contents), 0644); err != nil {
			return err
		}
	}
	return nil
}

func generateGoMod() error {
	for index, dir := range libs {
		contents := strings.Replace(mod, "REPLACEME", libs[index], -1)
		if err := ioutil.WriteFile(
			filepath.Join("parts", dir, "go.mod"), []byte(contents), 0644); err != nil {
			return err
		}
	}
	return nil
}

func Download() error {
	if err := download("i2p-zero-linux."+zeroversion+".zip", "https://github.com/i2p-zero/i2p-zero/releases/download/"+zeroversion+"/i2p-zero-linux."+zeroversion+".zip"); err != nil {
		return err
	}
	if err := download("i2p-zero-win."+zeroversion+".zip", "https://github.com/i2p-zero/i2p-zero/releases/download/"+zeroversion+"/i2p-zero-win-gui."+zeroversion+".zip"); err != nil {
		return err
	}
	if err := download("i2p-zero-darwin."+zeroversion+".zip", "https://github.com/i2p-zero/i2p-zero/releases/download/"+zeroversion+"/i2p-zero-mac."+zeroversion+".zip"); err != nil {
		return err
	}
	return nil
}

func download(path string, url string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Println("fetching", path, "from", url)
		// Get the data
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		// Create the file
		out, err := os.Create(path)
		if err != nil {
			return err
		}
		defer out.Close()
		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		return err
	}
	return nil
}

func generateGoGenerator(platform string) error {
	newfile := `package zerobundle

import (
`
	for index, dir := range libs {
		newfile += "	\"github.com/eyedeekay/zerobundle/parts/" + dir + "\"\n"
		log.Println(dir, libs[index])
	}

	newfile += ")\n\n"

	newfile += `func TBZBytes() ([]byte, error) {
	var bytes []byte
	`

	for index, dir := range libs {
		newfile += `	b` + dir + `, err := iz` + dir + `.WriteBrowser(iz` + dir + `.FS)
	if err != nil {
		return nil, err
	}
	bytes = append(bytes, b` + dir + `...)
`
		log.Println(dir, libs[index])
	}

	newfile += `	return bytes, nil
}`

	err := ioutil.WriteFile("import/embed_"+platform+".go", []byte(newfile), 0644)
	if err != nil {
		return err
	}
	return err
}
