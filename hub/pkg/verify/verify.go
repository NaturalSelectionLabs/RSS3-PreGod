package verify

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/verify/ethers"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/verify/json_util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/verify/nacl"
)

// Verifies if the current json file has a valid signature.
func Signature(jsonBytes []byte, address, instanceUri string) (bool, error) { //nolint:funlen,gocognit // TODO:
	jsonBytes, err := json_util.SortJsonByKeys(jsonBytes, nil)
	if err != nil {
		return false, err
	}

	jsonBytesWithoutSign, err := json_util.SortJsonByKeys(jsonBytes, &json_util.SortOptions{NoSignProperties: true})
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

	// check file type
	var jIndex *protocol.Index

	var jOther *protocol.SignedBase

	//nolint:nestif // TODO:
	if ji["profile"] != nil {
		// index file
		jIndex = new(protocol.Index)
		if err = json.Unmarshal(jsonBytes, jIndex); err != nil {
			return false, err
		} else if jIndex.Signature == "" {
			return false, fmt.Errorf("json has no signature field")
		}
	} else {
		// other file
		jOther = new(protocol.SignedBase)
		if err = json.Unmarshal(jsonBytes, jOther); err != nil {
			return false, err
		} else if jOther.Signature == "" {
			return false, fmt.Errorf("json has no signature field")
		}
	}

	// assign the signature
	var signature string
	if jIndex.Signature != "" {
		signature = jIndex.Signature
	} else {
		signature = jOther.Signature
	}

	// check if agents is present
	var retErr error

	if jIndex.Agents == nil && jOther.Agents == nil {
		// check stringified json signature
		var ethersOk bool
		ethersOk, retErr = ethers.VerifyMessage(getFileSignatureMessage(string(jsonBytesWithoutSign), instanceUri), signature, address)

		if ethersOk {
			return true, nil
		}
	}

	// assgin the agents
	var agents []protocol.Agent
	if jIndex.Agents != nil {
		agents = jIndex.Agents
	} else {
		agents = jOther.Agents
	}

	// check if any of the agents has a valid signature
	for _, agent := range agents {
		// verify if user has authorization to sign
		ethersOk, _ := ethers.VerifyMessage(getAgentSignatureMessage(agent.App, agent.Pubkey, instanceUri), agent.Authorization, address)

		// verify if file signature is valid
		naclOk, _ := nacl.Verify(jsonBytesWithoutSign, []byte(agent.Signature), []byte(agent.Pubkey))

		if ethersOk && naclOk {
			return true, nil
		}
	}

	// if jOther, no further checking is needed
	if jIndex == nil {
		return true, nil
	}

	// check if platform signature is valid
	for _, account := range jIndex.Profile.Accounts {
		signature := account.Signature
		if signature == "" {
			continue
		}

		ethersOk, _ := ethers.VerifyMessage(getPlatformSignatureMessage(signature, instanceUri), signature, address)
		if !ethersOk {
			return false, fmt.Errorf("'platform' signature is not valid (id: %s, sign: %s)", account.Identifier, account.Signature)
		}
	}

	if retErr == nil {
		return true, nil
	} else {
		return false, retErr
	}
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

func getFileSignatureMessage(fileJSON, instanceUri string) []byte {
	var buf bytes.Buffer

	buf.WriteString("[RSS3] I am confirming the results of changes to my file ")
	buf.WriteString(instanceUri)
	buf.WriteString(": ")
	buf.WriteString(fileJSON)

	return buf.Bytes()
}

func getPlatformSignatureMessage(platformID, instanceUri string) []byte {
	var buf bytes.Buffer

	buf.WriteString("[RSS3] I am adding account ")
	buf.WriteString(platformID)
	buf.WriteString(" to my RSS3 Instance ")
	buf.WriteString(instanceUri)

	return buf.Bytes()
}
