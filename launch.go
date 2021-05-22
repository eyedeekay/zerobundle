package zerobundle

/*
Released under the The MIT License (MIT)
see ./LICENSE
*/

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

import (
	"github.com/eyedeekay/checki2cp"
)

var i2cpConf = `i2cp.tcp4.host=127.0.0.1
i2cp.tcp4.port=7654
`

func zd() (dir string) {
	dir, _ = UnpackZeroDir()
	os.MkdirAll(dir, 0755)
	return
}

func baseargs() (args string) {
	args = "--i2p.dir.base=" + filepath.Join(zd(), "base")
	os.MkdirAll(filepath.Join(zd(), "base"), 0755)
	return
}

func configargs() (args string) {
	args = "--i2p.dir.config=" + filepath.Join(zd(), "config")
	os.MkdirAll(filepath.Join(zd(), "config"), 0755)
	return
}

func WriteI2CPConf() error {
	dir, err := UnpackZeroDir()
	if err != nil {
		return err
	}
	os.MkdirAll(dir, 0755)
	os.Setenv("I2CP_HOME", dir)
	os.Setenv("GO_I2CP_CONF", "/.i2cp.conf")
	home := os.Getenv("I2CP_HOME")
	conf := os.Getenv("GO_I2CP_CONF")
	if err := ioutil.WriteFile(filepath.Join(home, conf), []byte(i2cpConf), 0644); err != nil {
		return err
	}
	return nil
}

func ZeroMain() error {
	if err := WriteI2CPConf(); err != nil {
		return err
	}
	if ok, err := checki2p.ConditionallyLaunchI2P(); ok {
		if err != nil {
			return err
		}
	} else {
		if err := UnpackZero(); err != nil {
			log.Println(err)
		}
		latest := LatestZero()
		log.Println("latest zero version is:", latest)
		if !CheckZeroIsRunning() {
			log.Println("Zero doesn't appear to be running.", latest)
			if err := StartZero(); err != nil {
				return err
			}
		}
		if ok, conn := Available(); ok {
			log.Println("Starting SAM")
			time.Sleep(3 * time.Second)
			if err := SAM(conn); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("I2P router availability failure")
		}
	}
	time.Sleep(1 * time.Second)
	return nil
}

// ZeroAsFreestandingSAM need a SAM API? Don't have one? Launch a zero instance
// and tell it to start SAM because sometimes you want things to be easy.
func ZeroAsFreestandingSAM() error {
	if err := UnpackZero(); err != nil {
		log.Println(err)
	}
	latest := LatestZero()
	log.Println("latest zero version is:", latest)
	if !CheckZeroIsRunning() {
		log.Println("Zero doesn't appear to be running.", latest)
		if err := StartZero(); err != nil {
			return err
		}
	}
	if ok, conn := Available(); ok {
		log.Println("Starting SAM")
		time.Sleep(3 * time.Second)
		if err := SAM(conn); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("I2P router availability failure")
	}
	time.Sleep(1 * time.Second)
	return nil
}

var cmd *exec.Cmd

func loopbackInterface() string {
	if runtime.GOOS != "windows" {
		return "127.0.0.1"
	}
	ift, err := net.Interfaces()
	if err != nil {
		return "localhost"
	}
	log.Println("Searching for appropriate loopback interface")
	for _, ifi := range ift {
		if ifi.Flags&net.FlagLoopback != 0 && ifi.Flags&net.FlagUp != 0 {
			log.Println("Searching", ifi.Name)
			a, err := ifi.Addrs()
			if err != nil {
				return "localhost"
			}
			if !strings.Contains(a[0].String(), "::") {
				return strings.Split(a[0].String(), "/")[0]
			}
		}
	}
	return "localhost"
}

func CheckZeroIsRunning() bool {
	conn, err := net.Dial("tcp4", net.JoinHostPort(loopbackInterface(), "8051"))
	if err != nil {
		log.Println("Connecting error:", err)
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}

func GetZeroCMD() *exec.Cmd {
	return cmd
}

func GetZeroPID() int {
	return cmd.Process.Pid
}

func GetZeroProcess() *os.Process {
	return cmd.Process
}

func LatestZero() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(LatestZeroBinDir(), "i2p-zero.exe")
	} else {
		return filepath.Join(LatestZeroBinDir(), "i2p-zero")
	}
}

func LatestZeroJavaHome() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(LatestZeroBinDirJavaHome(), "i2p-zero.exe")
	} else {
		return filepath.Join(LatestZeroBinDirJavaHome(), "i2p-zero")
	}
}

func LatestZeroBinDir() string {
	var dir string
	var err error
	if dir, err = UnpackZeroDir(); err == nil {
		ks, er := ioutil.ReadDir(dir)
		fs := []os.FileInfo{}
		for _, k := range ks {
			if k.IsDir() {
				if strings.Contains(k.Name(), "i2p-zero-") {
					fs = append(fs, k)
				}
			}
		}
		if er != nil {
			log.Fatal(er)
		}
		if runtime.GOOS == "windows" {
			return filepath.Join(dir, fs[len(fs)-1].Name(), "router")
		} else {
			return filepath.Join(dir, fs[len(fs)-1].Name(), "router", "bin")
		}
	} else {
		log.Fatal(err)
	}
	return ""
}

func LatestZeroBinDirJavaHome() string {
	ks, er := ioutil.ReadDir(JAVA_I2P_OPT_DIR)
	fs := []os.FileInfo{}
	for _, k := range ks {
		if k.IsDir() {
			if strings.Contains(k.Name(), "i2p-zero-") {
				fs = append(fs, k)
			}
		}
	}
	if er != nil {
		log.Fatal(er)
	}
	if runtime.GOOS == "windows" {
		return filepath.Join(JAVA_I2P_OPT_DIR, fs[len(fs)-1].Name(), "router")
	} else {
		return filepath.Join(JAVA_I2P_OPT_DIR, fs[len(fs)-1].Name(), "router", "bin")
	}
}

func StopZero() {
	if runtime.GOOS == "windows" {
		GetZeroProcess().Signal(os.Kill)
	} else {
		GetZeroProcess().Signal(os.Interrupt)
	}
}

func CommandZero() (*exec.Cmd, error) {
	if err := UnpackZero(); err != nil {
		log.Println(err)
	}
	latest := LatestZero()
	cmd = exec.Command(latest, baseargs(), configargs())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

func CommandZeroContext(ctx context.Context) (*exec.Cmd, error) {
	if err := UnpackZero(); err != nil {
		log.Println(err)
	}
	latest := LatestZero()
	cmd = exec.CommandContext(ctx, latest, baseargs(), configargs())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

func RunZero() error {
	var err error
	cmd, err = CommandZero()
	if err != nil {
		return err
	}
	return cmd.Run()
}

func StartZero() error {
	var err error
	cmd, err = CommandZero()
	if err != nil {
		return err
	}
	return cmd.Start()
}

func CommandZeroJavaHome() (*exec.Cmd, error) {
	if err := UnpackZeroJavaHome(); err != nil {
		log.Println(err)
	}
	latest := LatestZeroJavaHome()
	cmd = exec.Command(latest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

func CommandZeroJavaHomeContext(ctx context.Context) (*exec.Cmd, error) {
	if err := UnpackZeroJavaHome(); err != nil {
		log.Println(err)
	}
	latest := LatestZeroJavaHome()
	cmd = exec.CommandContext(ctx, latest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

func RunZeroJavaHome() error {
	var err error
	cmd, err = CommandZeroJavaHome()
	if err != nil {
		return err
	}
	return cmd.Run()
}

func StartZeroJavaHome() error {
	var err error
	cmd, err = CommandZeroJavaHome()
	if err != nil {
		return err
	}
	return cmd.Start()
}

func Available() (bool, net.Conn) {
	i := 0
	for {
		conn, err := net.Dial("tcp4", net.JoinHostPort(loopbackInterface(), "8051"))
		if err != nil {
			log.Println("Connecting error:", err)
		}
		if conn != nil {
			log.Println("Zero is started.", err)
			return true, conn
		}
		i++
		time.Sleep(time.Duration(5) * time.Second)
	}
	return false, nil
}

func SAM(conn net.Conn) error {
	defer conn.Close()
	if runtime.GOOS == "windows" {
		conn.Write([]byte("sam.create\r\n"))
	} else {
		conn.Write([]byte("sam.create\n"))
	}
	i := 0
	for {
		samconn, err := net.Dial("tcp4", net.JoinHostPort(loopbackInterface(), "7656"))
		if err != nil {
			log.Println("Connecting error:", err)
		}
		if samconn != nil {
			log.Println("Started SAM.")
			conn.Close()
			return nil
		}
		i++
		time.Sleep(time.Duration(5) * time.Second)
	}
	return fmt.Errorf("Error connecting to %s", "SAM port")
}
