package kalpsdk

import (
	//Standard Libs
	"fmt"
	"testing"

	//Custom Build Libs
	"github.com/hyperledger/fabric-chaincode-go/shim"

	//Third party Libs
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/p2eengineering/kalp-sdk-public/mocks"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	mockClientIdentity := new(mocks.ClientIdentity)
	mockStub := new(mocks.ChaincodeStubInterface)
	ctx := &TransactionContext{
		clientIdentity: mockClientIdentity,
		stub:           mockStub,
	}

	expectedId := "eDUwOTo6Q049VGVzdE93bmVyLDEyMw=="
	t.Run("Check for success response", func(t *testing.T) {
		mockStub.On("GetState", "smartContractOwner").Return([]byte(""), nil).Once()
		// Set up the expected behavior of the mock stub
		mockClientIdentity.On("GetID").Return(expectedId, nil).Once()

		mockStub.On("PutState", "smartContractOwner", []byte("TestOwner")).Return(nil).Once()

		err := ctx.Initialize()
		require.NoError(t, err)
	})
}

func TestTransferOwner(t *testing.T) {
	mockClientIdentity := new(mocks.ClientIdentity)
	mockStub := new(mocks.ChaincodeStubInterface)
	ctx := &TransactionContext{
		clientIdentity: mockClientIdentity,
		stub:           mockStub,
	}

	expectedId := "eDUwOTo6Q049VGVzdE93bmVyLDEyMw=="
	t.Run("Check for success response", func(t *testing.T) {
		mockClientIdentity.On("GetID").Return(expectedId, nil).Once()
		
		expectedResponse := peer.Response{Status: shim.OK, Payload: []byte("true")}
		mockStub.On("InvokeChaincode", "kyc", [][]byte{[]byte("KycExists"), []byte("TestOwner")}, "universalkyc").Return(expectedResponse).Once()

		mockStub.On("PutState", "smartContractOwner", []byte("TestOwner")).Return(nil).Once()

		err := ctx.TransferOwner()
		require.NoError(t, err)
	})
}

func TestPutStateWithoutKYC(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	tx := &TransactionContext{
		stub: mockStub,
	}

	// Check for success response
	t.Run("check for success response", func(t *testing.T) {
		mockStub.On("PutState", "key", []byte("value")).Return(nil).Once()

		err := tx.PutStateWithoutKYC("key", []byte("value"))
		// require that no error occurred
		require.NoError(t, err)
	})

	// Check for failure response
	t.Run("check for failure response", func(t *testing.T) {
		mockStub.On("PutState", "key", []byte("value")).Return(fmt.Errorf("failed to put key,value in state"))

		err := tx.PutStateWithoutKYC("key", []byte("value"))
		// require that the expected error occurred
		expectedErr := fmt.Errorf("failed to put key,value in state")
		require.EqualError(t, err, expectedErr.Error())
	})
}

func TestDelStateWithoutKYC(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	tx := &TransactionContext{
		stub: mockStub,
	}

	// Check for success response
	t.Run("Check for success response", func(t *testing.T) {
		mockStub.On("DelState", "key").Return(nil).Once()

		err := tx.DelStateWithoutKYC("key")
		// require that no error occurred
		require.NoError(t, err)
	})

	// Check for failure response
	t.Run("Check for failure response", func(t *testing.T) {
		mockStub.On("DelState", "key").Return(fmt.Errorf("failed to delete the value on state"))

		err := tx.DelStateWithoutKYC("key")
		// require that the expected error occurred
		expectedErr := fmt.Errorf("failed to delete the value on state")
		require.EqualError(t, err, expectedErr.Error())
	})
}

func TestInvokeChaincode(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	tx := &TransactionContext{
		stub: mockStub,
	}

	channelName := "universalkyc"
	chaincodeName := "kyc"
	args := [][]byte{[]byte("arg1"), []byte("arg2")}

	// Check for success response
	t.Run("Check for success response", func(t *testing.T) {
		expectedResponse := peer.Response{Status: shim.OK, Payload: []byte("true")}
		mockStub.On("InvokeChaincode", chaincodeName, args, channelName).Return(expectedResponse)

		response := tx.InvokeChaincode(chaincodeName, args, channelName)
		// require that no error occurred
		require.Equal(t, expectedResponse, response.Response)
	})

	// Check for failure response
	t.Run("Check for failure response", func(t *testing.T) {
		expectedResponse := peer.Response{Status: shim.ERROR, Payload: []byte("false")}
		mockStub.On("InvokeChaincode", "", args, channelName).Return(expectedResponse)

		response := tx.InvokeChaincode("", args, channelName)
		// require that the expected error occurred
		require.Equal(t, expectedResponse, response.Response)
	})
}

func TestPutKYC(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	tx := &TransactionContext{
		stub: mockStub,
	}

	params := []string{"CreateKyc", "sampleId", "kycId", "kycHash"}
	invokeArgs := make([][]byte, len(params))
	for i, arg := range params {
		invokeArgs[i] = []byte(arg)
	}

	// Check for success response
	t.Run("Check for success response", func(t *testing.T) {
		expectedResponse := peer.Response{Status: shim.OK, Payload: []byte("true")}
		mockStub.On("InvokeChaincode", "kyc", invokeArgs, "universalkyc").Return(expectedResponse).Once()

		err := tx.PutKYC("sampleId", "kycId", "kycHash")
		// require that no error occurred
		require.NoError(t, err)
	})

	// Check for failure response
	t.Run("check for failure response", func(t *testing.T) {
		expectedResponse := peer.Response{Status: shim.ERROR, Payload: []byte("failed to query kyc chaincode")}
		mockStub.On("InvokeChaincode", "kyc", invokeArgs, "universalkyc").Return(expectedResponse).Once()

		err := tx.PutKYC("sampleId", "kycId", "kycHash")
		// require that the expected error occurred
		require.Error(t, err)
	})
}

func TestPutStateWithKYC(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	mockClientIdentity := new(mocks.ClientIdentity)
	tx := &TransactionContext{
		stub:           mockStub,
		clientIdentity: mockClientIdentity,
	}

	// Check for success response
	t.Run("Check for success response", func(t *testing.T) {
		expectedResponse := peer.Response{Status: shim.OK, Payload: []byte("true")}
		mockClientIdentity.On("GetID").Return("eDUwOTo6Q049VGVzdE93bmVyLDEyMw==", nil)
		mockStub.On("InvokeChaincode", "kyc", [][]byte{[]byte("KycExists"), []byte("TestOwner")}, "universalkyc").Return(expectedResponse)
		mockStub.On("GetKYC", "eDUwOTo6Q049VGVzdE93bmVyLDEyMw==").Return(false, nil)
		mockStub.On("PutState", "key", []byte("value")).Return(nil).Once()

		err := tx.PutStateWithKYC("key", []byte("value"))
		// require that no error occurred
		require.NoError(t, err, "Expected no error for user with completed kyc")
	})

	// Check for failure response
	t.Run("Check for failure response", func(t *testing.T) {
		expectedResponse := peer.Response{Status: shim.ERROR, Payload: []byte("failed to query kyc chaincode")}
		mockClientIdentity.On("GetID").Return("eDUwOTo6Q049VGVzdE93bmVyLDEyMw==", nil)
		mockStub.On("InvokeChaincode", "kyc", [][]byte{[]byte("KycExists"), []byte("TestOwner")}, "universalkyc").Return(expectedResponse)
		mockStub.On("GetKYC", "eDUwOTo6Q049VGVzdE93bmVyLDEyMw==").Return(false, nil)
		mockStub.On("PutState", "key", []byte("value")).Return(fmt.Errorf("failed to put key and value"))

		//  failure response
		err := tx.PutStateWithKYC("key", []byte("value"))
		// require that the expected error occurred
		require.Error(t, err)
	})
}

func TestDelStateWithKYC(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	mockClientIdentity := new(mocks.ClientIdentity)
	tx := &TransactionContext{
		stub:           mockStub,
		clientIdentity: mockClientIdentity,
	}

	// Check for success response
	t.Run("Check for success response", func(t *testing.T) {
		expectedResponse := peer.Response{Status: shim.OK, Payload: []byte("true")}
		mockClientIdentity.On("GetID").Return("eDUwOTo6Q049VGVzdE93bmVyLDEyMw==", nil)
		mockStub.On("InvokeChaincode", "kyc", [][]byte{[]byte("KycExists"), []byte("TestOwner")}, "universalkyc").Return(expectedResponse)
		mockStub.On("GetKYC", "eDUwOTo6Q049VGVzdE93bmVyLDEyMw==").Return(false, nil)
		mockStub.On("DelState", "key").Return(nil).Once()

		err := tx.DelStateWithKYC("key")
		// require that no error occurred
		require.NoError(t, err, "Expected no error for DelStateWithKYC")
	})

	// Check for failure response
	t.Run("Check for failure response", func(t *testing.T) {
		expectedResponse := peer.Response{Status: shim.ERROR, Payload: []byte("failed to query kyc chaincode")}
		mockClientIdentity.On("GetID").Return("eDUwOTo6Q049VGVzdE93bmVyLDEyMw==", nil)
		mockStub.On("InvokeChaincode", "kyc", [][]byte{[]byte("KycExists"), []byte("TestOwner")}, "universalkyc").Return(expectedResponse)
		mockStub.On("GetKYC", "eDUwOTo6Q049VGVzdE93bmVyLDEyMw==").Return(false, nil)
		mockStub.On("DelState", "key").Return(fmt.Errorf("failed to Delete the state from the world state. "))

		err := tx.DelStateWithKYC("key")
		// require that the expected error occurred
		require.Error(t, err)
	})
}
