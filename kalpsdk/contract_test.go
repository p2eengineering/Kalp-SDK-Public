package kalpsdk

import (
	//Standard Libs
	"testing"

	//Third party Libs
	"github.com/hyperledger/fabric-contract-api-go/metadata"
	"github.com/stretchr/testify/require"
)

// ReturnsString is a method of myContract that returns a string.
func ReturnsString() string {
	return "Some string"
}

type customContext struct {
	TransactionContext
}

func TestGetInfo(t *testing.T) {
	c := Contract{}
	c.Info = metadata.InfoMetadata{}
	c.Info.Version = "Some String"

	// Checked the expected result
	require.Equal(t, c.Info, c.GetInfo(), "should set the version")
}

func TestGetName(t *testing.T) {
	contract := Contract{}

	// Checked the expected result
	require.Equal(t, "", contract.GetName(), "should have returned blank ns when not set")

	contract.Name = "myname"
	require.Equal(t, "myname", contract.GetName(), "should have returned custom ns when set")
}

func TestGetBeforeTransaction(t *testing.T) {
	var contract Contract
	var beforeFn interface{}

	contract = Contract{}
	beforeFn = contract.GetBeforeTransaction()

	// Checked the beforeFn must be nil
	require.Nil(t, beforeFn, "should not return contractFunction when before transaction not set")

	contract = Contract{}
	contract.BeforeTransaction = ReturnsString
	beforeFn = contract.GetBeforeTransaction()

	// Checked the expected before function loaded
	require.Equal(t, ReturnsString(), beforeFn.(func() string)(), "function returned should be same value as set for before transaction")
}

func TestGetUnknownTransaction(t *testing.T) {
	var contract Contract
	var unknownFn interface{}

	contract = Contract{}
	unknownFn = contract.GetUnknownTransaction()

	// Checked the unknown transaction must be nil
	require.Nil(t, unknownFn, "should not return contractFunction when unknown transaction not set")

	contract = Contract{}
	contract.UnknownTransaction = ReturnsString
	unknownFn = contract.GetUnknownTransaction()

	// Checked the expected unknown transaction function loaded
	require.Equal(t, ReturnsString(), unknownFn.(func() string)(), "function returned should be same value as set for unknown transaction")
}

func TestGetTransactionContextHandler(t *testing.T) {
	// Checked the expected transaction context
	contract := Contract{}
	require.Equal(t, new(TransactionContext), contract.GetTransactionContextHandler(), "should return default transaction context type when unset")

	contract.TransactionContextHandler = new(customContext)
	require.Equal(t, new(customContext), contract.GetTransactionContextHandler(), "should return custom context when set")
}
