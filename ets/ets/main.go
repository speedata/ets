package main

import (
	"fmt"
	"os"
	"time"

	"github.com/speedata/ets/core"
	"github.com/speedata/optionparser"
)

func dothings() error {
	op := optionparser.NewOptionParser()
	op.Command("version", "Show version information")

	err := op.Parse()
	if err != nil {
		return err
	}

	if len(op.Extra) == 0 {
		return fmt.Errorf("please specify a file to run")
	}
	return core.Dothings(op.Extra[0])
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
