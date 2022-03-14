package main

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
}

func main() {
	command := &cobra.Command{
		Use: "pregod-cli",
	}

	command.AddCommand(migrate.NewMigrateCommand())

	if err := command.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
