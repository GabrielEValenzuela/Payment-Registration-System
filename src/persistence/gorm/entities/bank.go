package entities

import "time"

func (BankEntity) TableName() string {
	return "BANKS"
}

// Bank represents a financial institution that holds customers and issues cards.
type BankEntity struct {
	ID        uint             `gorm:"primaryKey;autoIncrement"`
	Name      string           `gorm:"size:255"`
	Cuit      string           `gorm:"size:255"`
	Address   string           `gorm:"size:255"`
	Telephone string           `gorm:"size:255"`
	Customers []CustomerEntity `gorm:"many2many:CUSTOMERS_BANKS;"`
	CreatedAt time.Time        `gorm:"autoCreateTime"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime"`
}
