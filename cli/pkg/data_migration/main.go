package data_migration

import (
	"flag"
	"log"
)

var (
	mongouri string
)

func init() {
	flag.StringVar(&mongouri, "mongouri", "", "old mongouri for db migration")
}

func Start() {
	flag.Parse()

	if mongouri == "" {
		flag.PrintDefaults()

		return
	}

	if err := prepareDB(); err != nil {
		log.Fatalln(err)
	}

	if err := migrate(mongouri); err != nil {
		log.Fatalln(err)
	}
}
