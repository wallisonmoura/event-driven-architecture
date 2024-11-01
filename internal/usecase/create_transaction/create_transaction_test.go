package create_transaction

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wallison/fc-ms-wallet/internal/entity"
	"github.com/wallison/fc-ms-wallet/internal/event"
	"github.com/wallison/fc-ms-wallet/internal/usecase/mocks"
	"github.com/wallison/fc-ms-wallet/pkg/events"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("Client 1", "j@j.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("Client 2", "j@j2.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	inputDTO := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}

	dispather := events.NewEventDispatcher()
	event := event.NewTransactionCreated()
	ctx := context.Background()

	uc := NewCreateTransactionUseCase(mockUow, dispather, event)
	output, err := uc.Execute(ctx, inputDTO)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
