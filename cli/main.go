package main

import (
	"log"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate"
	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use: "pregod-cli",
	}

	command.AddCommand(migrate.NewMigrateCommand())

	if err := command.Execute(); err != nil {
		log.Fatalln(err)
	}
}
