package verify

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/verify/ethers"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/verify/json_util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/verify/nacl"
)

type agent struct {
	Pubkey        string `json:"pubkey"`
	Signature     string `json:"signature"`
	Authorization string `json:"authorization"`
	App           string `json:"app"`
	DateExpired   string `json:"date_expired"`
}

// Verifies if the current json file has a valid signature.
func Signature(jsonBytes []byte, address, instanceUri string) (bool, error) {
	jsonBytes, err := json_util.SortJsonByKeys(jsonBytes, &json_util.SortOptions{NoSignProperties: true})
	if err != nil {
		return false, err
	}

	var ji map[string]interface{}

	err = json.Unmarshal(jsonBytes, &ji)
	if err != nil {
		return false, err
	}

	if ji == nil {
		return false, fmt.Errorf("json is nil")
	}

	// check if signature is present
	if ji["signature"] == nil {
		return false, fmt.Errorf("json has no signature field")
	}

	// check if signature is valid
	signature, ok := ji["signature"].(string)
	if !ok {
		return false, fmt.Errorf("signature field is not a string")
	}

	// check if agents is present
	if ji["agents"] == nil {
		// check stringified json signature
		ethersOk, err := ethers.VerifyMessage(jsonBytes, signature, address)
		if ethersOk {
			return true, err
		}
	}

	// check if agents is valid
	agents, ok := ji["agents"].([]agent)
	if !ok {
		return false, fmt.Errorf("'agents' field is not valid")
	}

	// check if any of the agents has a valid signature
	for _, agent := range agents {
		// verify if user has authorization to sign
		ethersOk, _ := ethers.VerifyMessage(getAgentSignatureMessage(agent.App, agent.Pubkey, instanceUri), agent.Authorization, address)

		// verify if file signature is valid
		naclOk, _ := nacl.Verify(jsonBytes, []byte(agent.Signature), []byte(agent.Pubkey))

		if ethersOk && naclOk {
			return true, nil
		}
	}

	return false, fmt.Errorf("no valid signature found")
}

// `[RSS3] I am well aware that this APP (name: ${app}) can use
// the following agent instead of me (${InstanceURI}) to
// modify my files and I would like to authorize this agent (${pubkey})`
func getAgentSignatureMessage(appname, pubkey, instanceUri string) []byte {
	var buf bytes.Buffer

	buf.WriteString("[RSS3] I am well aware that this APP (name: ")
	buf.WriteString(appname)
	buf.WriteString(") can use the following agent instead of me (")
	buf.WriteString(instanceUri)
	buf.WriteString(") to modify my files and I would like to authorize this agent (")
	buf.WriteString(pubkey)
	buf.WriteString(")")

	return buf.Bytes()
}
