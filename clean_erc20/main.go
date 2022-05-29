package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/clean_erc20/internal"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	mapset "github.com/deckarep/golang-set"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
)

func init() {
	if err := database.Setup(); err != nil {
		log.Fatalf("database.Setup err: %v", err)
	}
}

// const GetNotesLimit = 20000

const GetNotesLimit = 1
const loopTime = 500 * time.Millisecond

var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary

func getApiKey() string {
	apiKey, err := jsoni.MarshalToString(config.Config.Indexer.Moralis.ApiKey)
	if err != nil {
		return ""
	}

	return strings.Trim(apiKey, "\"")
}

func RunReplaceWrongEThEndpoint(cmd *cobra.Command, args []string) error {
	for {
		notes, err := internal.GetDataFromDB(GetNotesLimit)
		if err != nil {
			logger.Infof("get data from db err:%v", err)

			return err
		}

		if len(notes) == 0 {
			logger.Infof("mission completed")
		}

		time.Sleep(time.Second * 5)

		// change db
		internal.ReplaceEndpoint(notes)

		logger.Infof("notes[0]: %v", notes[0])

		if _, err := database.CreateNotes(database.DB, notes, true); err != nil {
			logger.Errorf("err:%v", err)

			continue
		}

		time.Sleep(loopTime)
	}
}

// nolint:funlen // TODO
func RunFixEmptyTokenSymbol(cmd *cobra.Command, args []string) error {
	logger.Debugf("start")
	var chainType = moralis.ChainType(moralis.ETH)

	for {
		// get one err account
		// identifier, err := internal.GetOneTokenSymbolEmptyIdentifier(chainType)
		// if err != nil {
		// 	logger.Errorf("get one token symbol empty account err[%v], account[%s]", err, identifier)

		// 	continue
		// }

		// if identifier == "" {
		// 	break
		// }

		identifier := "rss3://account:0x639090cdd215010fe54c36a49bbd1604f034e1d4@ethereum"

		logger.Debugf("identifier:%s", identifier)

		account, err := internal.GetAccountByIdentifier(identifier)
		if err != nil {
			logger.Warnf("get account by identifier err[%v], identifier[%s]", err, identifier)

			continue
		}

		logger.Debugf("account:%v", account)

		// get this one all err notes
		notesMap, err := internal.GetAllNotesAboutErc20ByIdentifier(identifier, chainType)
		if err != nil {
			logger.Warnf("get all notes about erc20 by account err[%v], account[%s]", err, account)

			continue
		}

		logger.Debugf("len(notesMap):%d", len(notesMap))

		// logger.Debugf("notesMap:%v", notesMap)

		// get this one msg from api
		// result, err := moralis.GetErc20Transfers(context.Background(), account, chainType, "", getApiKey())
		// if err != nil {
		// 	logger.Warnf("get erc20 transfers err[%v]", err)

		// 	continue
		// }

		// if len(result) == 0 {
		// 	continue
		// }

		tokenAddressSet := mapset.NewSet()
		tokenAddresses := []string{}

		for _, note := range notesMap {
			tokenAddressSet.Add(note.TokenAddress)
		}

		for _, tokenAddress := range tokenAddressSet.ToSlice() {
			addressStr, ok := tokenAddress.(string)
			if !ok {
				logger.Warnf("token address[%v] is not string", addressStr)

				continue
			}

			tokenAddresses = append(tokenAddresses, addressStr)
		}

		// get the token metadata
		erc20Tokens, err := moralis.GetErc20TokenMetaData(context.Background(), chainType, tokenAddresses, getApiKey())
		if err != nil {
			logger.Errorf("chain type[%s], get erc20 token metadata [%v]",
				chainType.GetNetworkSymbol().String(), err)

			continue
		}

		logger.Debugf("0xa58a4f5c4bb043d2cc1e170613b74e767c94189b symbol:", erc20Tokens["0xa58a4f5c4bb043d2cc1e170613b74e767c94189b"])

		// logger.Debugf("erc20Tokens:%v", erc20Tokens)

		// get resp ,update notes
		notes, err := internal.ChangeNotesTokenSymbolMsg(notesMap, erc20Tokens)
		if err != nil {
			logger.Errorf("set erc20 token symbol msg in notes err[%v]", err)

			continue
		}

		// set in db
		if _, err := database.CreateNotes(database.DB, notes, true); err != nil {
			logger.Errorf("err:%v", err)

			continue
		}
		/**/
		break
	}

	return nil
}

var rootCmd = &cobra.Command{Use: "clean_erc20"}

func main() {
	rootCmd.AddCommand(&cobra.Command{
		Use:  "replace_wrong_eth_endpoint",
		RunE: RunReplaceWrongEThEndpoint,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:  "run_fix_empty_token_symbol",
		RunE: RunFixEmptyTokenSymbol,
	})

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
