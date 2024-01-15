# Kalp-SDK

Welcome to the documentation for the Kalp Software Development Kit (SDK). This guide will walk you through the process of setting up the SDK, configuring it to connect to a Kalptantra network, and performing various operations on the network.

## Overview

The Kalp SDK is a comprehensive Golang package specifically designed to simplify the development process of chaincode (or smart contracts) on the Kalptantra blockchain network. It empowers developers to write and create Kalptantra-compliant chaincode. The SDK provides a set of powerful functionalities that streamline interaction with the network and enhance the overall development experience.

## Key Features

The Kalp SDK offers a range of key features that simplify application development on the Kalptantra blockchain network:

- **Data Read and Write:** The Kalp SDK provides functions to read and write data to the Kalptantra blockchain network. Developers can store and retrieve key-value pairs in the ledger using functions such as PutStateWithKYC, PutStateWithoutKYC, and GetState. This enables seamless integration of data storage and retrieval within smart contracts.

- **Transaction Management:** With the Kalp SDK, developers can efficiently manage transactions on the blockchain network. It provides functions for submitting transactions, querying transaction information, and retrieving transaction history. This simplifies the process of interacting with the blockchain and ensures the integrity of transactional operations.

- **KYC Checks:** The SDK includes built-in functionality for performing Know Your Customer (KYC) checks. Developers can leverage the GetKYC function to verify if a user has completed the KYC process on the network. This feature enhances the security and compliance of applications built on the Kalptantra network.

- **Payment Tracking for Payable Contracts:** The Kalp SDK supports payment tracking for Payable contracts. Developers can easily track payments made for Payable contracts and retrieve payment information using the provided functions. This simplifies the implementation of payment-related functionality within smart contracts, enabling the development of decentralized applications with payment capabilities.

- **Logger functionality:** The Kalp SDK provides a Logger support as well. You can create a Logger object which will give you better visibility of the kalpsdk operations.

## Installation

To install the Kalp-SDK package, use the following command:

```go
go get -u github.com/p2eengineering/kalp-sdk-public/kalpsdk
```

## Examples

### Creating a Contract

To create a contract using the Kalp-SDK, you need to define a new Go struct that represents your contract and embed the kalpsdk.Contract struct into your struct to inherit the base contract functionalities.

```go
type MyContract struct {
	kalpsdk.Contract
}
```

### Implement the Contract Interface

After defining the contract struct, you need to implement the contract interface by defining the `Init` and `Invoke` methods. These methods will contain the logic for initializing the contract and handling the invocations respectively.

```go
func (c *MyContract) Init(ctx kalpsdk.TransactionContextInterface) error {
	// Initialization logic
	return nil
}

func (c *MyContract) Invoke(ctx kalpsdk.TransactionContextInterface, data string) error {
	// Invoke logic
	return nil
}
```

### Creating and Starting Chaincode

To create a new chaincode using the Kalp-SDK and start it, follow these steps:

#### Create a new Chaincode Instance

Create a new instance of the `kalpsdk.Chaincode` struct. Pass your contract struct as an argument to the `NewChaincode` function and specify whether the contract is payable or not.

```go
// Creating a sample payable contract object
contract := kalpsdk.Contract{IsPayableContract: true}

// Creating a KalpSDK Logger object
contract.Logger = kalpsdk.NewLogger()

// Create a new instance of your KalpContractChaincode with your smart contract
chaincode, err := kalpsdk.NewChaincode(&MyContract{contract})

if err != nil {
  log.Panicf("Error creating KalpContractChaincode: %v", err)
}
```

#### Start the Chaincode

Call the Start function on the chaincode instance to start your chaincode.

```go
if err := chaincode.Start(); err != nil {
  panic(fmt.Sprintf("Error starting chaincode: %v", err))
}
```

## Blockchain Data Management

### Writing to the Blockchain

To write data to the Kalptantra blockchain using the Kalp-SDK, you can use the `PutStateWithKyc` and `PutStateWithoutKyc` functions. These functions allow you to store a key-value pair in the ledger with or without KYC verification.

#### PutStateWithKyc

The `PutStateWithKyc` is a blockchain function provided by the Kalp-SDK. It allows you to write data to the ledger with KYC verification. This function ensures that only users who have completed KYC can make changes to the ledger, providing an additional layer of security and compliance.

Function Parameters:

- `key` (String): The key under which the data will be stored in the ledger.
- `value` (Byte Array): The data to be stored in the ledger as a byte array.

```go
err := ctx.PutStateWithKyc("myKey", []byte("myValue"))
if err != nil {
  // Handle error
} else {
  // Data successfully written to the blockchain with KYC verification
}
```

#### PutStateWithoutKyc

The PutStateWithoutKyc is a blockchain function provided by the Kalp-SDK. It allows you to write data to the ledger without requiring KYC verification. This function does not enforce any restrictions based on KYC completion and allows users to change the ledger.

**Note: Using the PutStateWithoutKyc function bypasses the KYC verification requirement, allowing any user to write data to the ledger. However, it's crucial to be aware that this can have implications for security and compliance, as it does not enforce restrictions on who can modify the ledger.**


Function Parameters:

- `key` (String): The key under which the data will be stored in the ledger.
- `value` (Byte Array): The data to be stored in the ledger as a byte array.

```go
err := ctx.PutStateWithoutKyc("myKey", []byte("myValue"))
if err != nil {
  // Handle error
} else {
  // Data successfully written to the blockchain without KYC verification
}
```

### Reading from the Blockchain

To read data from the Kalptantra blockchain using the Kalp-SDK, you can use the `GetState` blockchain function. This function allows you to retrieve the value associated with a given key from the ledger.

### GetState
The GetState is a blockchain function provided by the Kalp-SDK. It allows you to get data from the ledger.

Function Parameters:

- `key` (String): The key for which to retrieve the value from the ledger.

Return Value:

- `value` (Byte Array): The value associated with the specified key.

```go
value, err := ctx.GetState("myKey")
if err != nil {
  // Handle error
} else {
  // Process the retrieved value
}
```

## Deleting from the Blockchain

To delete data from the Kalptantra blockchain using the Kalp-SDK, you can use the `DelStateWithKyc` and `DelStateWithoutKyc` functions. These functions allow you to remove a key-value pair from the ledger with or without KYC verification.

### DelStateWithKyc

The `DelStateWithKyc` is a blockchain function provided by the Kalp-SDK. It allows you to delete data from the ledger with KYC verification. This function ensures that only users who have completed KYC can remove data from the ledger, providing an additional layer of security and compliance.

Function Parameters:

- `key` (String): The key of the data to be deleted from the ledger.

```go
err := ctx.DelStateWithKyc("myKey")
if err != nil {
  // Handle error
} else {
  // Data successfully deleted from the blockchain with KYC verification
}
```

#### DelStateWithoutKyc

The `DelStateWithoutKyc` is a blockchain function provided by the Kalp-SDK. It allows you to delete data from the ledger without requiring KYC verification. This function does not enforce any restrictions based on KYC completion and allows users to remove data from the ledger.

**Note: Using the DelStateWithoutKyc function bypasses the KYC verification requirement, allowing any user to delete data from the ledger. However, it's crucial to be aware that this can have implications for security and compliance, as it does not enforce restrictions on who can modify the ledger.**

Function Parameters:

- `key` (String): The key of the data to be deleted from the ledger.

```go
err := ctx.DelStateWithoutKyc("myKey")
if err != nil {
  // Handle error
} else {
  // Data successfully deleted from the blockchain without KYC verification
}
```

### Checking KYC Status

To check if a user has completed KYC on the network, you can use the `GetKYC` function provided by the Kalp-SDK. It allows you to check if a user has completed KYC on the network.

Function Parameters:

- `userId` (String): The ID of the user to check for KYC completion.

Returns:

- `bool`: `true` if the user has completed KYC, `false` otherwise.
- `error`: An error object if any error occurs during the KYC status check.

```go
Kyced, err := ctx.GetKYC("userId")
if err != nil {
  // Handle error
}

if Kyced {
  // User has completed KYC, proceed with the desired action
  // ...
} else {
  // User has not completed KYC, handle accordingly
  // ...
}
```
##

**Happy coding with the Kalp-SDK and enjoy building innovative decentralized applications on the Kalptantra blockchain network!**
