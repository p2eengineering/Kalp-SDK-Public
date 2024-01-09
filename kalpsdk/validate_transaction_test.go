package kalpsdk

import (
	//Standard Libs
	"fmt"
	"testing"

	//Custom Build Libs
	"github.com/p2eengineering/kalp-sdk/mocks"

	//Third party Libs
	"github.com/stretchr/testify/require"
)

func TestValidateCreateTokenTransaction(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	mockState := new(mocks.StateQueryIteratorInterface)
	mockClientIdentity := new(mocks.ClientIdentity)
	ctx := &TransactionContext{
		stub:           mockStub,
		clientIdentity: mockClientIdentity,
	}

	id := "sampleId"
	docType := "ASSET-R2CI"
	queryString := fmt.Sprintf(`{"selector": {"id": "%s", "docType": "%s"}}`, id, docType)
	expectedId := "eDUwOTo6Q049VGVzdE93bmVyLDEyMw=="

	mockStub.On("GetQueryResult", queryString).Return(mockState, nil)
	mockState.On("HasNext").Return(false).Once()
	mockClientIdentity.On("GetID").Return(expectedId, nil).Once()

	// Check for success response
	err := ctx.ValidateCreateTokenTransaction(id, docType, []string{"TestOwner"})
	if err != nil {
		t.Errorf("Expected no error but Got:%v", err)
	}
	require.NoError(t, err)

}

func TestIsMinted(t *testing.T) {
	mockStub := new(mocks.ChaincodeStubInterface)
	mockState := new(mocks.StateQueryIteratorInterface)

	ctx := &TransactionContext{
		stub: mockStub,
	}

	id := "sampleId"
	docType := "ASSET-R2CI"
	queryString := fmt.Sprintf(`{"selector": {"id": "%s", "docType": "%s"}}`, id, docType)

	// Check for success response
	t.Run("Ckeck for success response", func(t *testing.T) {
		mockStub.On("GetQueryResult", queryString).Return(mockState, nil)
		mockState.On("HasNext").Return(true).Once()
		expectedbool := true
		actualbool, err := IsMinted(ctx, id, docType)
		require.NoError(t, err)
		require.Equal(t, expectedbool, actualbool)
	})

	// Check for failure response
	t.Run("Ckeck for success response", func(t *testing.T) {
		mockStub.On("GetQueryResult", queryString).Return(nil, fmt.Errorf("failed to get query result from the world state"))
		mockState.On("HasNext").Return(false).Once()
		expectedbool := false

		actualbool, err := IsMinted(ctx, id, docType)
		require.NoError(t, err)
		require.Equal(t, expectedbool, actualbool)
	})
}
