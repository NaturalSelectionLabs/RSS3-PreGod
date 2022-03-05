package main

import (
	"flag"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/pkg/data_migration"
)

var (
	isDataMigrationCmd bool
)

func init() {
	flag.BoolVar(&isDataMigrationCmd, "data-mig", false, "Use Data Migration Command")
}

func main() {
	flag.Parse()

	if isDataMigrationCmd {
		data_migration.Start()
	} else {
		flag.PrintDefaults()
	}
}
