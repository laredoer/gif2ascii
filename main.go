package main

import (
	"log"
	"os"

	"github.com/wule61/gif2ascii/cli"
)

func main() {
	app := cli.New()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
