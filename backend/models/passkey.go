package models

import (
	"time"
)

// Passkey 对应数据库中的 Passkeys 表
type Passkey struct {
	Passkey     string    `gorm:"column:Passkey;type:VARCHAR(36);primaryKey;not null;default:uuid()" json:"passkey"`
	Description string    `gorm:"column:Description;type:VARCHAR(511);not null;default:''" json:"description"`
	Operator    string    `gorm:"column:Operator;type:TINYTEXT;not null;default:''" json:"operator"`
	CreatedAt   time.Time `gorm:"column:Created_At;not null;default:NOW()" json:"created_at"`
	IsActive    bool      `gorm:"column:Is_Active;type:TINYINT(1);not null;default:1" json:"is_active"`
	Extends     string    `gorm:"column:Extends;type:TINYTEXT;not null;default:''" json:"extends"`
}

// TableName 指定表名
func (Passkey) TableName() string {
	return "Passkeys"
}
