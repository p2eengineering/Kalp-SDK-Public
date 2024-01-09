package kalpsdk

import (
	//Standard Libs
	"os"
	"testing"
	"time"

	//Custom Build Libs
	"github.com/p2eengineering/kalp-sdk/mocks"

	//Third party Libs
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/stretchr/testify/require"
)

type MockChaincodeStub struct {
	shim.ChaincodeStubInterface
}

func (m *MockChaincodeStub) GetFunctionAndParameters() (string, []string) {
	return "", []string{}
}

func (m *MockChaincodeStub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) peer.Response {
	return peer.Response{}
}

func TestInit(t *testing.T) {
	chaincode := &ContractChaincode{}
	stub := &MockChaincodeStub{}

	// Call the Init function
	response := chaincode.Init(stub)

	// require the response status code is as expected
	require.Equal(t, shim.OK, int(response.Status))
}

func TestInvoke(t *testing.T) {
	stub := &mocks.ChaincodeStubInterface{}

	stub.On("GetFunctionAndParameters").Return("test", []string{})

	chaincode := &ContractChaincode{}

	// Call the Invoke function
	response := chaincode.Invoke(stub)

	// require the response status code is as expected
	// require.Equal(t, shim.OK, int(response.Status))
	require.NotNil(t, response)
}

func TestStart(t *testing.T) {
	// Set the CORE_CHAINCODE_ID_NAME environment variable
	err := os.Setenv("CORE_CHAINCODE_ID_NAME", "myChaincode")
	require.NoError(t, err, "Failed to set CORE_CHAINCODE_ID_NAME environment variable")

	err = os.Setenv("CHAINCODE_SERVER_ADDRESS", ":6050")
	require.NoError(t, err, "Failed to set CORE_CHAINCODE_ID_NAME environment variable")

	// Create a new instance of ContractChaincode
	chaincode := &ContractChaincode{}

	go func() {
		err = chaincode.Start()
		// require that the Start method did not return an error
		require.NoError(t, err, "Start method returned an error")
	}()

	// Delay to check the chaincode server starts without any error
	time.Sleep(time.Second)

	// Remove the CORE_CHAINCODE_ID_NAME environment variable
	err = os.Unsetenv("CORE_CHAINCODE_ID_NAME")
	require.NoError(t, err, "Failed to unset CORE_CHAINCODE_ID_NAME environment variable")

	// Remove the CHAINCODE_SERVER_ADDRESS environment variable
	err = os.Unsetenv("CHAINCODE_SERVER_ADDRESS")
	require.NoError(t, err, "Failed to unset CHAINCODE_SERVER_ADDRESS environment variable")
}
