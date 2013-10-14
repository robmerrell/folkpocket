package main

import (
	"fmt"
	"github.com/robmerrell/comandante"
	"github.com/robmerrell/folkpocket/cmds"
	"github.com/robmerrell/folkpocket/config"
	"os"
)

func main() {
	config.LoadConfigFile("config.toml")

	bin := comandante.New("folkpocket", "Read folklore.org in pocket")
	bin.IncludeHelp()

	// scrape command
	scrapeCmd := comandante.NewCommand("scrape", "Scrape folklore.org story URLs and save them in the database.", cmds.ScrapeAction)
	scrapeCmd.Documentation = cmds.ScrapeDoc
	bin.RegisterCommand(scrapeCmd)

	// serve command
	serveCmd := comandante.NewCommand("serve", "Serve the folkpocket webiste", cmds.ServeAction)
	serveCmd.Documentation = cmds.ServeDoc
	bin.RegisterCommand(serveCmd)

	// run the commands
	if err := bin.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
