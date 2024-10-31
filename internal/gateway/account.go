package gateway

import "github.com/wallison/fc-ms-wallet/internal/entity"

type AccountGateway interface {
	Save(account *entity.Account) error
	FindById(id string) (*entity.Account, error)
}
