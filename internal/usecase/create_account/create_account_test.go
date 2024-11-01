package create_account

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wallison/fc-ms-wallet/internal/entity"
	"github.com/wallison/fc-ms-wallet/internal/usecase/mocks"
)

func TestAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	clientMock := &mocks.ClientGatewayMock{}
	clientMock.On("Get", client.ID).Return(client, nil)

	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("Save", mock.Anything).Return(nil)

	uc := NewCreateAccountUseCase(accountMock, clientMock)
	inputDTO := CreateAccountInputDTO{
		ClientId: client.ID,
	}

	output, err := uc.Execute(inputDTO)
	assert.Nil(t, err)
	assert.NotNil(t, output.ID)
	clientMock.AssertExpectations(t)
	accountMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t, "Get", 1)
	accountMock.AssertNumberOfCalls(t, "Save", 1)

}
