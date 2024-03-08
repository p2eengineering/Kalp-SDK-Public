package kalpsdk

import (
	//Standard Libs
	"fmt"
	"testing"
	"time"

	//Custom Build Libs
	"github.com/p2eengineering/kalp-sdk-public/mocks"

	//Third party Libs
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockHistoryIterator struct {
	mock.Mock
}

func (m *mockHistoryIterator) HasNext() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockHistoryIterator) Next() (*ChaincodeResponse, error) {
	args := m.Called()
	return args.Get(0).(*ChaincodeResponse), args.Error(1)
}

type ChaincodeResponse struct {
	TxId      string
	Value     []byte
	Timestamp *timestampProto
	IsDelete  bool
}

type timestampProto struct {
	Seconds int64
	Nanos   int32
}

func TestIsSmartContractOwner(t *testing.T) {
	mockClientIdentity := new(mocks.ClientIdentity)
	mockStub := new(mocks.ChaincodeStubInterface)
	ctx := &TransactionContext{
		clientIdentity: mockClientIdentity,
		stub:           mockStub,
	}

	expectedId := "eDUwOTo6Q049VGVzdE93bmVyLDEyMw=="
	t.Run("Check for success response", func(t *testing.T) {
		mockStub.On("GetState", "smartContractOwner").Return([]byte("TestOwner"), nil).Once()
		// Set up the expected behavior of the mock stub
		mockClientIdentity.On("GetID").Return(expectedId, nil).Once()

		isOwner, err := ctx.IsSmartContractOwner()
		if !isOwner || err != nil {
			t.Errorf("unexpected result: isOwner=%t, err=%v", isOwner, err)
		}
		require.NoError(t, err)
	})

	t.Run("check for false response", func(t *testing.T) {
		mockStub.On("GetState", "smartContractOwner").Return([]byte("TestOwnerFake"), nil).Once()
		// Set up the expected behavior of the mock stub
		mockClientIdentity.On("GetID").Return(expectedId, nil).Once()

		isOwner, err := ctx.IsSmartContractOwner()
		if err != nil {
			t.Errorf("unexpected result: isOwner=%t, err=%v", isOwner, err)
		}
		require.NoError(t, err)
	})
}

func TestFetchOwnerHistory(t *testing.T) {
	mockClientIdentity := new(mocks.ClientIdentity)
	mockHisQuery := new(mocks.HistoryQueryIteratorInterface)
	mockStub := new(mocks.ChaincodeStubInterface)
	ctx := &TransactionContext{
		clientIdentity: mockClientIdentity,
		stub:           mockStub,
	}

	expectedResult := &queryresult.KeyModification{
		TxId:  "txid",
		Value: []byte("historyforkey"),
	}

	t.Run("Check for success response", func(t *testing.T) {
		mockStub.On("GetHistoryForKey", "key").Return(mockHisQuery, nil).Once()
		mockHisQuery.On("Next").Return(expectedResult, nil)
		mockHisQuery.On("HasNext").Return(false)

		hisQuery, err := ctx.GetHistoryForKey("key")
		require.NoError(t, err)

		km, err := hisQuery.Next()
		require.NoError(t, err)

		// require that no error occurred
		require.Equal(t, expectedResult, km)
		require.False(t, hisQuery.HasNext())
	})
}

func TestSetEvent(t *testing.T) {
	tx := &TransactionContext{
		stub: &shim.ChaincodeStub{},
	}

	err := tx.SetEvent("name", []byte("payload"))
	// require that no error occurred
	require.NoError(t, err)

	// require that the expected error occurred
	err = tx.SetEvent("", []byte("event payload"))
	require.EqualError(t, err, "event name can not be empty string")
}

func TestGetTxID(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)

	// Set up the expected behavior of the mock stub
	mockStub.On("GetTxID").Return("sampleTxID")

	ctx := &TransactionContext{
		stub: mockStub,
	}

	// Call the GetTxID function and require the expected result
	txID := ctx.GetTxID()
	require.Equal(t, "sampleTxID", txID)
}

func TestGetState(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	ctx := &TransactionContext{
		stub: mockStub,
	}

	// Check for success response
	t.Run("Check for success response", func(t *testing.T) {
		// Set up the expected behavior of the mock stub
		mockStub.On("GetState", "key").Return([]byte("myvalue"), nil).Once()
		payload := []byte("myvalue")

		getState, err := ctx.GetState("key")
		require.NoError(t, err)

		// require that no error occurred
		require.Equal(t, payload, getState)
	})

	// Check for failure response
	t.Run("Check the Failure response", func(t *testing.T) {
		// Set up the expected behavior of the mock stub
		mockStub.On("GetState", "key").Return(nil, fmt.Errorf("Failed to get Value: Failure to retrieving the data")).Once()

		_, err := ctx.GetState("key")
		// require that the expected error occurred
		expectedErr := fmt.Errorf("Failed to get Value: Failure to retrieving the data")
		require.EqualError(t, err, expectedErr.Error())
	})
}

func TestGetUserID(t *testing.T) {
	mockClientIdentity := new(mocks.ClientIdentity)
	ctx := &TransactionContext{
		clientIdentity: mockClientIdentity,
	}

	expectedId := "eDUwOTo6Q049VGVzdE93bmVyLDEyMw=="
	decodeId := "TestOwner"

	// Check for success response
	t.Run("Check for success response", func(t *testing.T) {
		// Set up the expected behavior of the mock stub
		mockClientIdentity.On("GetID").Return(expectedId, nil).Once()

		id, err := ctx.GetUserID()
		require.NoError(t, err)

		// require that no error occurred
		require.Equal(t, decodeId, id)
	})

	// Check for failure response
	t.Run("Check for GetId Error", func(t *testing.T) {
		// Set up the expected behavior of the mock stub
		mockClientIdentity.On("GetID").Return("", fmt.Errorf("failure to read Id")).Once()

		_, err := ctx.GetUserID()
		// require that the expected error occurred
		expectedErr := fmt.Errorf("failed to read clientID: failure to read Id")
		require.EqualError(t, err, expectedErr.Error())
	})
}

func TestGetKYC(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	ctx := &TransactionContext{
		stub: mockStub,
	}

	// Check for success response
	t.Run("Check for Success response", func(t *testing.T) {
		expectedResponse := peer.Response{Status: shim.OK, Payload: []byte("true")}
		mockStub.On("InvokeChaincode", "kyc", [][]byte{[]byte("KycExists"), []byte("TestUser")}, "universalkyc").Return(expectedResponse).Once()

		result, err := ctx.GetKYC("TestUser")
		require.True(t, result)

		// require that no error occurred
		require.NoError(t, err)
	})

	// Check for failure response
	t.Run("Check for Failure response", func(t *testing.T) {
		userID := "TestUser"
		expectedResponse := peer.Response{Status: shim.ERROR, Payload: []byte("failed to query kyc chaincode")}
		mockStub.On("InvokeChaincode", "kyc", [][]byte{[]byte("KycExists"), []byte(userID)}, "universalkyc").Return(expectedResponse)

		_, err := ctx.GetKYC(userID)
		// require that the expected error occurred
		require.Error(t, err)
	})
}

func TestGetChannelID(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	ctx := &TransactionContext{
		stub: mockStub,
	}

	mockStub.On("GetChannelID").Return("sampleChannelId")

	channelID := ctx.GetChannelID()
	// require that no error occurred
	require.Equal(t, "sampleChannelId", channelID)
}

func TestGetStateByPartialCompositeKey(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	mockState := new(mocks.StateQueryIteratorInterface)
	ctx := &TransactionContext{
		stub: mockStub,
	}

	expectedResult := &queryresult.KV{
		Key:   "querykey",
		Value: []byte("queryvalue"),
	}

	// Check for success response
	t.Run("Check for success response", func(t *testing.T) {
		// Set up the expected behavior of the mock stub
		mockStub.On("GetStateByPartialCompositeKey", "object", []string{"attr1", "attr2"}).Return(mockState, nil).Once()
		mockState.On("Next").Return(expectedResult, nil)

		stateQuery, err := ctx.GetStateByPartialCompositeKey("object", []string{"attr1", "attr2"})
		require.NoError(t, err)

		kv, err := stateQuery.Next()
		require.NoError(t, err)

		// require that no error occurred
		require.Equal(t, expectedResult, kv)
	})

	// Check for failure response
	t.Run("Check for Failure response", func(t *testing.T) {
		mockStub.On("GetStateByPartialCompositeKey", "", []string{"attr1", "attr2"}).Return(nil, fmt.Errorf("Failed to retrieve the composite keys by partial composite key"))
		mockState.On("Close").Return(nil)

		iterator, err := ctx.GetStateByPartialCompositeKey("", []string{"attr1", "attr2"})
		// require that the expected error occurred
		require.Error(t, err)
		require.Nil(t, iterator)
		require.EqualError(t, err, "Failed to retrieve the composite keys by partial composite key")
	})
}

func TestGetStateByRange(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	mockState := new(mocks.StateQueryIteratorInterface)
	ctx := &TransactionContext{
		stub: mockStub,
	}

	// Check for success response
	t.Run("Check for Success Response", func(t *testing.T) {
		mockStub.On("GetStateByRange", "key1", "key5").Return(mockState, nil)

		iterator, err := ctx.GetStateByRange("key1", "key5")
		// require that no error occurred
		require.Equal(t, mockState, iterator)
		require.NoError(t, err)
	})

	// Check for failure response
	t.Run("Check for Failure Response", func(t *testing.T) {
		mockStub.On("GetStateByRange", "", "key5").Return(nil, fmt.Errorf("Failure in retrieving the keys by range"))

		iterator, err := ctx.GetStateByRange("", "key5")
		// require that the expected error occurred
		require.Error(t, err)
		require.Nil(t, iterator)
		require.EqualError(t, err, "Failure in retrieving the keys by range")
	})
}

func TestGetQueryResult(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	mockState := new(mocks.StateQueryIteratorInterface)
	ctx := &TransactionContext{
		stub: mockStub,
	}

	expectedResult := &queryresult.KV{
		Key:   "querykey",
		Value: []byte("queryvalue"),
	}

	// Check for success response
	t.Run("Check for Success response", func(t *testing.T) {
		mockStub.On("GetQueryResult", "object").Return(mockState, nil).Once()
		mockState.On("Next").Return(expectedResult, nil)

		stateQuery, err := ctx.GetQueryResult("object")
		require.NoError(t, err)

		kv, err := stateQuery.Next()
		require.NoError(t, err)

		// require that no error occurred
		require.Equal(t, expectedResult, kv)
	})

	// Check for failure response
	t.Run("Check for Failure response", func(t *testing.T) {
		mockStub.On("GetQueryResult", "key").Return(nil, fmt.Errorf("Failure to retrieving the query"))

		_, err := ctx.GetQueryResult("key")
		// require that the expected error occurred
		expectedResult := fmt.Errorf("Failure to retrieving the query")
		require.EqualError(t, err, expectedResult.Error())
	})
}

func TestGetHistoryForKey(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	mockHisQuery := new(mocks.HistoryQueryIteratorInterface)
	ctx := &TransactionContext{
		stub: mockStub,
	}

	expectedResult := &queryresult.KeyModification{
		TxId:  "txid",
		Value: []byte("historyforkey"),
	}

	// Check for success response
	t.Run("Check for success Response", func(t *testing.T) {
		mockStub.On("GetHistoryForKey", "key").Return(mockHisQuery, nil).Once()
		mockHisQuery.On("Next").Return(expectedResult, nil)
		mockHisQuery.On("HasNext").Return(false)

		hisQuery, err := ctx.GetHistoryForKey("key")
		require.NoError(t, err)

		km, err := hisQuery.Next()
		require.NoError(t, err)

		// require that no error occurred
		require.Equal(t, expectedResult, km)
		require.False(t, hisQuery.HasNext())
	})

	// Check for failure response
	t.Run("check for failure response", func(t *testing.T) {
		mockStub.On("GetHistoryForKey", "key").Return(nil, fmt.Errorf("failure to retreiving the history")).Once()

		_, err := ctx.GetHistoryForKey("key")
		// require that the expected error occurred
		expectedResult := fmt.Errorf("failure to retreiving the history")
		require.EqualError(t, err, expectedResult.Error())
	})
}

func TestCreateCompositeKey(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	ctx := &TransactionContext{
		stub: mockStub,
	}

	// Check for success response
	t.Run("Check for success response", func(t *testing.T) {
		mockStub.On("CreateCompositeKey", "Owner", []string{"id"}).Return("sampleID", nil)
		key := "sampleID"

		result, err := ctx.CreateCompositeKey("Owner", []string{"id"})
		require.NoError(t, err)

		// require that no error occurred
		require.Equal(t, key, result)
	})

	// Check for failure response
	t.Run("check for failure response", func(t *testing.T) {
		mockStub.On("CreateCompositeKey", "", []string{"id"}).Return("", fmt.Errorf("Failure to create Composite key: Expected error for empty object type"))

		_, err := ctx.CreateCompositeKey("", []string{"id"})
		require.Error(t, err)

		// require that the expected error occurred
		expectedErr := fmt.Errorf("Failure to create Composite key: Expected error for empty object type")
		require.EqualError(t, err, expectedErr.Error())
	})
}

func TestGetTxTimestamp(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	ctx := &TransactionContext{
		stub: mockStub,
	}

	expectedTimestamp := &timestamp.Timestamp{
		Seconds: time.Now().Unix(),
		Nanos:   int32(time.Now().Nanosecond()),
	}

	// Check for success response
	t.Run("Check for success response", func(t *testing.T) {
		mockStub.On("GetTxTimestamp").Return(expectedTimestamp, nil)

		result, err := ctx.GetTxTimestamp()
		require.NoError(t, err)

		actualTime := result.AsTime()
		expectedTime := expectedTimestamp.AsTime()
		require.Equal(t, expectedTime, actualTime)
	})
}

func TestGetFunctionAndParameters(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	ctx := &TransactionContext{
		stub: mockStub,
	}

	expectedFunction := "sampleFunction"
	expectedParams := []string{"param1", "param2"}

	// Check for success response
	t.Run("Check for success response", func(t *testing.T) {
		mockStub.On("GetFunctionAndParameters").Return(expectedFunction, expectedParams)

		actualFunction, actualParams := ctx.GetFunctionAndParameters()
		// require that no error occurred
		require.Equal(t, expectedFunction, actualFunction)
		require.Equal(t, expectedParams, actualParams)
	})
}
