package blockchain

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/core/logger"
	mockresponses "github.com/smartcontractkit/external-initiator/integration/mock-client/blockchain/mock-responses"
)

type cannedResponses map[string][]JsonrpcMessage

// GetCannedResponses returns the static responses from a file, if such exists. JSON-RPC ID not set!
func GetCannedResponses(platform string) (cannedResponses, bool) {
	bz, err := mockresponses.Get(platform)
	if err != nil {
		logger.Debug(err)
		return nil, false
	}

	var responses cannedResponses
	err = json.Unmarshal(bz, &responses)
	if err != nil {
		logger.Error(err)
		return nil, false
	}

	return responses, true
}

// GetCannedResponse returns the static response from a file, if such exists for the JSON-RPC method.
func GetCannedResponse(platform string, msg JsonrpcMessage) ([]JsonrpcMessage, bool) {
	responses, ok := GetCannedResponses(platform)
	if !ok {
		return nil, false
	}

	responseList, ok := responses[msg.Method]
	if !ok {
		return nil, false
	}

	return setJsonRpcId(msg.ID, responseList), true
}

//nolint:stylecheck,golint
func setJsonRpcId(id json.RawMessage, msgs []JsonrpcMessage) []JsonrpcMessage {
	for i := 0; i < len(msgs); i++ {
		msgs[i].ID = id
	}
	return msgs
}
