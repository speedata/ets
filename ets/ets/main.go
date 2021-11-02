package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/speedata/ets/core"
	"github.com/speedata/optionparser"
)

var (
	version string
)

const (
	cmdRun     = "run"
	cmdHelp    = "help"
	cmdVersion = "version"
)

func dothings() error {
	pathToExefile, err := os.Executable()
	exename := filepath.Base(pathToExefile)
	if err != nil {
		return err
	}

	op := optionparser.NewOptionParser()
	op.Banner = "experimental typesetting system\nrun: ets somefile.lua"
	op.Command(cmdVersion, "Show version information")
	op.Command(cmdHelp, "Show usage help")

	err = op.Parse()
	if err != nil {
		return err
	}

	if len(op.Extra) == 0 {
		return fmt.Errorf("Please specify a command or file to run. See %s --help", exename)
	}
	switch op.Extra[0] {
	case cmdVersion:
		fmt.Printf("%s verson %s\n", exename, version)
		os.Exit(0)
	case cmdHelp:
		op.Help()
		os.Exit(0)
	}
	return core.Dothings(op.Extra[0], exename)
}

func main() {
	start := time.Now()
	err := dothings()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(time.Now().Sub(start))
}
