package main

import (
	"log"
	"strconv"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/reptile/pkg/handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/spf13/cobra"
)

func init() {
	if err := database.Setup(); err != nil {
		log.Fatalf("database.Setup err: %v", err)
	}
}

func RunPullInfo(cmd *cobra.Command, args []string) error {
	defaultEndPos := 6000

	if len(args) > 1 {
		arguEndPos, err := strconv.Atoi(args[1])
		if err != nil {
			logger.Warnf("invalid end position: %s", args[1])
		} else {
			defaultEndPos = arguEndPos
		}
	}

	handler.PullInformation(defaultEndPos)

	return nil
}

func RunSetDBDataLower(cmd *cobra.Command, args []string) error {
	handler.SetDBDataDressToLower()

	return nil
}

func RunGetResultByStage(cmd *cobra.Command, args []string) error {
	handler.GetResultByStage()

	return nil
}

var rootCmd = &cobra.Command{Use: "reptile"}

func main() {
	rootCmd.AddCommand(&cobra.Command{
		Use:  "pullinfo",
		RunE: RunPullInfo,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:  "setdbdata",
		RunE: RunSetDBDataLower,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:  "getresultbystage",
		RunE: RunGetResultByStage,
	})

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
