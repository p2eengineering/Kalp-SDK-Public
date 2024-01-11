// Sample code for creating a SmartContract with kalpsdk library
package smartcontract

import (
	"log"

	// Import Kalpsdk library
	kalpsdk "github.com/p2eengineering/kalp-sdk-public/kalpsdk"
)

func main() {

	// Creating a sample payable contract object
	contract := kalpsdk.Contract{IsPayableContract: true}

	// Creating a KalpSDK Logger object
	contract.Logger = kalpsdk.NewLogger()

	// Create a new instance of your KalpContractChaincode with your smart contract
	chaincode, err := kalpsdk.NewChaincode(&SmartContract{contract})
	contract.Logger.Info("My KAPL SDK sm4")

	if err != nil {
		log.Panicf("Error creating KalpContractChaincode: %v", err)
	}

	// Start the chaincode
	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
