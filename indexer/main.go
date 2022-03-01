package main

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

func main() {
	if err := config.Setup(); err != nil {
		logger.Fatalf("config.Setup err: %v", err)
	}
}
