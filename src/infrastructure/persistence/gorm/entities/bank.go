package entities

import "time"

func (BankEntity) TableName() string {
	return "BANKS"
}

// Bank represents a financial institution that holds customers and issues cards.
type BankEntity struct {
	ID        uint   `gorm:"primaryKey"`
	Cuit      string `gorm:"size:255"`
	Address   string `gorm:"size:255"`
	Telephone string `gorm:"size:255"`
	//CustomersIds []int
	CreatedAt time.Time
	UpdatedAt time.Time
}
