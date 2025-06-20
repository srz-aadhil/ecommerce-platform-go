package domain

import "time"

type BankAccount struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	UserId      int       `json:"userid"`
	BankAccount uint      `json:"bankaccount" gorm:"index;unique;not null"`
	SwiftCode   string    `json:"swiftcode"`
	PaymentType string    `json:"paymenttype"`
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp"`
}
