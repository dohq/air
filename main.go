package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cosmtrek/air/runner"
)

var (
	cfgPath       string
	debugMode     bool
	showVersion   bool
	outputExample bool
)

//go:embed air_example.toml
var exampleConfigFile string

func init() {
	flag.StringVar(&cfgPath, "c", "", "config path")
	flag.BoolVar(&debugMode, "d", false, "debug mode")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&outputExample, "e", false, "output example config (e.g. air -e > .air.toml)")
	flag.Parse()
}

func main() {
	if outputExample {
		fmt.Print(exampleConfigFile)
		return
	}

	fmt.Printf(`
  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ %s // live reload for Go apps, with Go %s

`, airVersion, goVersion)

	if showVersion {
		return
	}

	if debugMode {
		fmt.Println("[debug] mode")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var err error
	r, err := runner.NewEngine(cfgPath, debugMode)
	if err != nil {
		log.Fatal(err)
		return
	}
	go func() {
		<-sigs
		r.Stop()
	}()

	defer func() {
		if e := recover(); e != nil {
			log.Fatalf("PANIC: %+v", e)
		}
	}()

	r.Run()
}
