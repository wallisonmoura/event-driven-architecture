package createtransaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wallison/fc-ms-wallet/internal/entity"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) FindById(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("Client 1", "j@j.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("Client 2", "j@j2.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	mockAccount := &AccountGatewayMock{}
	mockAccount.On("FindById", account1.ID).Return(account1, nil)
	mockAccount.On("FindById", account2.ID).Return(account2, nil)

	mockTransaction := &TransactionGatewayMock{}
	mockTransaction.On("Create", mock.Anything).Return(nil)

	inputDTO := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}

	uc := NewCreateTransactionUseCase(mockTransaction, mockAccount)
	output, err := uc.Execute(inputDTO)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockAccount.AssertExpectations(t)
	mockTransaction.AssertExpectations(t)
	mockAccount.AssertNumberOfCalls(t, "FindById", 2)
	mockTransaction.AssertNumberOfCalls(t, "Create", 1)
}
