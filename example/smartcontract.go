package smartcontract

import (
	//Standard Libs
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"

	//Custom Build Libs
	kalpsdk "github.com/p2eengineering/kalp-sdk-public/kalpsdk"
	
)

const nameKey = "name"
const symbolKey = "symbol"
const statusInProgress = "INPROGRESS"
const statusCompleted = "COMPLETED"

// Smart Contract Object
type SmartContract struct {
	kalpsdk.Contract
}

// NIU Structure
type NIU struct {
	Id          string      `json:"id"`
	DocType     string      `json:"docType"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Desc        string      `json:"desc"`
	Status      string      `json:"status"`
	Account     []string    `json:"account"`
	MetaData    interface{} `json:"metadata"`
	Amount      uint64      `json:"amount,omitempty" metadata:",optional"`
	Uri         string      `json:"uri"`
	AssetDigest string      `json:"assetDigest"`
}

// Initialize function initializes the smart contract by setting the name and symbol for the token.
// It takes the transaction context interface and the token name and symbol as input parameters.
// If the token name already exists in the state, the function returns an error indicating that the client is not authorized to change them.
// If the PutStateWithKYC function fails to set the token name or symbol, the function returns an error message.
func (s *SmartContract) Initialize(sdk kalpsdk.TransactionContextInterface, data string) (bool, error) {
	inputData := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &inputData)
	if err != nil {
		return false, fmt.Errorf("failed to fetch input %v", err)
	}

	name, ok := inputData["name"].(string)
	if !ok {
		return false, fmt.Errorf("name is required field")
	}

	symbol, ok := inputData["symbol"].(string)
	if !ok {
		return false, fmt.Errorf("name is required field")
	}

	// Get the token name from the state
	bytes, err := sdk.GetState(nameKey)
	if err != nil {
		return false, fmt.Errorf("failed to get token name from state: %v", err)
	}
	if bytes != nil {
		return false, fmt.Errorf("contract options are already set, client is not authorized to change them")
	}

	// Set the token name in the state with KYC
	err = sdk.PutStateWithKYC(nameKey, []byte(name))
	if err != nil {
		return false, fmt.Errorf("failed to set token name: %v", err)
	}

	// Set the token symbol in the state with KYC
	err = sdk.PutStateWithKYC(symbolKey, []byte(symbol))
	if err != nil {
		return false, fmt.Errorf("failed to set token symbol: %v", err)
	}

	return true, nil
}

// CreateNIU is a smart contract function which takes the Asset or Token input as JSON and store it in blockchain.
// It takes the transaction context interface and the input data string as input parameters.
// The function returns an error if it fails to parse the input data, if the input data is invalid, or if it fails to mint and store the token and metadata.
func (s *SmartContract) CreateNIU(sdk kalpsdk.TransactionContextInterface, data string) error {
	// Parse input data into NIU struct.
	var niu NIU
	errs := json.Unmarshal([]byte(data), &niu)
	if errs != nil {
		return fmt.Errorf("failed to parse data: %v", errs)
	}

	// Validate input data.
	if niu.AssetDigest == "" {
		return fmt.Errorf("assetDigest can not be null")
	}

	
	if niu.Status != statusInProgress && niu.Status != statusCompleted {
		return fmt.Errorf("not a valid Status")
	}

	// Common validation for mint new token
	err := sdk.ValidateCreateTokenTransaction(niu.Id, niu.DocType, niu.Account)
	if err != nil {
		return err
	}

	// Generate JSON representation of NIU struct.
	niuJSON, err := json.Marshal(niu)
	if err != nil {
		return fmt.Errorf("unable to Marshal Token struct : %v", err)
	}

	// Marshal metadata and put it in state.
	metadataByte, err := json.Marshal(niu.MetaData)
	if err != nil {
		return fmt.Errorf("failed to marshal the metadata: %v", err,metadataByte )
	}

	// // Mint token and store the JSON representation in the state database.
	// err = kaps.MintWithTokenURIMetadata(sdk, niu.Account, niu.Id, niu.Amount, niu.Uri, metadataByte, niu.DocType)
	// if err != nil {
	// 	return err
	// }

	// Store the NIU struct in the state database
	if err := sdk.PutStateWithKYC(niu.Id, niuJSON); err != nil {
		return fmt.Errorf("unable to put Asset struct in statedb: %v", err)
	}

	return nil
}

// ReadNIU retrieves the NIU asset with the given ID from the world state and returns it as a pointer to a NIU struct.
func (s *SmartContract) ReadNIU(sdk kalpsdk.TransactionContextInterface, id string) (*NIU, error) {
	// Get the asset from the ledger using id & check if asset exists
	niuJSON, err := sdk.GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if niuJSON == nil {
		return nil, fmt.Errorf("the NIU asset with ID %v does not exist", id)
	}

	// Unmarshal asset from JSON to struct
	var niu NIU
	err = json.Unmarshal(niuJSON, &niu)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal NIU struct: %v", err)
	}

	// Get the operator's client ID
	operator, err := sdk.GetUserID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client id: %v", err)
	}

	// Check if operator is a valid owner of the asset
	if !slices.Contains(niu.Account, operator) {
		return nil, fmt.Errorf("not a valid owner %v for the NIU asset with ID %v", niu.Account, id)
	}

	return &niu, nil
}

// TransferNIU function transfers NIU tokens from senders to receivers
// using KAPS contract functionality and updates the state of the asset in the world state.
func (s *SmartContract) TransferNIU(sdk kalpsdk.TransactionContextInterface, senders []string, receivers []string, id string, docType string, amount uint64, timeStamp string) error {
	// Retrieve asset from the world state using its ID
	niuJSON, err := sdk.GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if niuJSON == nil {
		return fmt.Errorf("the Asset %v does not exist", id)
	}

	// Unmarshal the asset JSON data into a struct
	var niu NIU
	err = json.Unmarshal(niuJSON, &niu)
	if err != nil {
		return fmt.Errorf("failed to unmarshal struct: %v", err)
	}

	// // Check if DocType is valid
	// if docType == "" || (docType != kaps.DocTypeNIU && docType != kaps.DocTypeAsset) {
	// 	return fmt.Errorf("invalid DocType, expected NIU or ASSET-NIU")
	// }

	// Check KYC status for each recipient
	for i := 0; i < len(receivers); i++ {
		kycCheck, err := sdk.GetKYC(receivers[i])
		if err != nil {
			return fmt.Errorf("failed to perform KYC check for user:%s, error:%v", receivers[i], err)
		}
		if !kycCheck {
			return fmt.Errorf("user %s is not KYCed", receivers[i])
		}
		if slices.Contains(niu.Account, receivers[i]) {
			return fmt.Errorf("transfer to self is not allowed")
		}
	}

	// Update the asset owner to the new recipient
	//var OrgSenders = niu.Account
	niu.Account = receivers

	// Marshal the updated asset struct into JSON data and verify the asset hash
	newniuJSON, err := json.Marshal(niu)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %v", err)
	}

	// Check if the asset is COMPLETED before transferring tokens
	if niu.Status != statusCompleted {
		return fmt.Errorf("asset is not applicable to Transfer")
	}

	// // Transfer tokens using the KAPS contract functionality
	// err = kaps.TransferFrom(sdk, senders, receivers, id, amount, docType, OrgSenders)
	// if err != nil {
	// 	return fmt.Errorf("failed to transfer tokens: %v", err)
	// }

	// Save the updated asset state in the world state
	if err := sdk.PutStateWithKYC(id, newniuJSON); err != nil {
		return fmt.Errorf("unable to put Asset struct in statedb: %v", err)
	}

	// Emit an event
	if err := sdk.SetEvent("TransferNIU", newniuJSON); err != nil {
		return fmt.Errorf("unable to setEvent TransferNIU: %v", err)
	}
	return nil
}

// BurnTokens burns tokens associated with a given asset and deletes the asset from the world state
func (s *SmartContract) BurnTokens(sdk kalpsdk.TransactionContextInterface, id string, account string, amount uint64) error {
	// Retrieve the asset from the world state using its ID
	niuJSON, err := sdk.GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if niuJSON == nil {
		return fmt.Errorf("the asset %v does not exist", id)
	}

	// Unmarshal the asset JSON data into a struct
	var niu NIU
	if err := json.Unmarshal(niuJSON, &niu); err != nil {
		return fmt.Errorf("failed to unmarshal struct: %v", err)
	}

	// Get the operator's client ID
	operator, err := sdk.GetUserID()
	if err != nil {
		return fmt.Errorf("failed to get client ID: %v", err)
	}

	// Check if the operator is a valid owner of the asset
	if !slices.Contains(niu.Account, operator) {
		return fmt.Errorf("only the asset owner can initiate a burn")
	}

	// Check if the asset is COMPLETED before burning tokens
	if niu.Status != statusCompleted {
		return fmt.Errorf("the asset is not applicable for burning tokens")
	}

	// // Burn the tokens using the KAPS contract
	// if niu.DocType == kaps.DocTypeNIU {
	// 	// Burn tokens with amount
	// 	err = kaps.Burn(sdk, []string{account}, id, amount, niu.DocType)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to burn tokens: %v", err)
	// 	}
	// } else if niu.DocType == kaps.DocTypeAsset {
	// 	// Burn tokens without amount (set to 0)
	// 	err = kaps.Burn(sdk, niu.Account, id, 0, niu.DocType)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to burn tokens: %v", err)
	// 	}
	// } else {
	// 	return fmt.Errorf("unknown document type: %s", niu.DocType)
	// }

	// Delete the asset from the world state
	if err := sdk.DelStateWithoutKYC(id); err != nil {
		return fmt.Errorf("unable to delete asset struct in state database: %v", err)
	}

	// Emit an event indicating the asset has been deleted
	if err := sdk.SetEvent("DeleteNIU", niuJSON); err != nil {
		return fmt.Errorf("unable to set event DeleteNIU: %v", err)
	}

	return nil
}
