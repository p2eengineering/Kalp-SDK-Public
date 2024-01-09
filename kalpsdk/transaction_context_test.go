package kalpsdk

import (
	//Standard Libs
	"testing"

	//Custom Build Libs
	"github.com/p2eengineering/kalp-sdk/mocks"

	//Third party Libs
	"github.com/stretchr/testify/require"
)

func TestSetStub(t *testing.T) {
	ctx := TransactionContext{}
	stub := &mocks.ChaincodeStubInterface{}
	ctx.SetStub(stub)

	// Check if the retrieved stub is the same as the one previously set
	require.Equal(t, stub, ctx.stub, "should have set the same stub as passed")
}

func TestSetClientIdentity(t *testing.T) {
	client := &mocks.ClientIdentity{}
	ctx := TransactionContext{}
	ctx.SetClientIdentity(client)

	// Check if the retrieved client identity is the same as the one previously set
	require.Equal(t, client, ctx.clientIdentity, "should have set the same client identity as passed")
}

func TestGetStub(t *testing.T) {
	stub := &mocks.ChaincodeStubInterface{}
	ctx := TransactionContext{}
	ctx.stub = stub

	// Check if the retrieved stub is the same as the one previously set
	require.Equal(t, stub, ctx.GetStub(), "should have returned same stub as set")
}

func TestGetClientIdentity(t *testing.T) {
	client := &mocks.ClientIdentity{}
	ctx := TransactionContext{}
	ctx.clientIdentity = client

	// Check if the retrieved client identity is the same as the one previously set
	require.Equal(t, client, ctx.GetClientIdentity(), "should have returned same client identity as set")
}
