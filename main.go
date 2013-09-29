package main

import (
	"fmt"
	"github.com/robmerrell/comandante"
	"github.com/robmerrell/folkpocket/cmds"
	"os"
)

func main() {
	bin := comandante.New("folkpocket", "Read folklore.org in pocket")
	bin.IncludeHelp()

	// scrape command
	scrapeCmd := comandante.NewCommand("scrape", "Scrape folklore.org story URLs and save them in the database.", cmds.ScrapeAction)
	scrapeCmd.Documentation = cmds.ScrapeDoc
	bin.RegisterCommand(scrapeCmd)

	// serve command
	if err := bin.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
