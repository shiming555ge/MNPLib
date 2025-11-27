package models

import (
	"time"
)

// Data 对应数据库中的 data 表
type Data struct {
	ID          string     `gorm:"column:ID;type:VARCHAR(10);primaryKey;not null" json:"id"`
	Source      string     `gorm:"column:Source;type:VARCHAR(20);not null" json:"source"`
	ItemName    *string    `gorm:"column:Item_Name;type:VARCHAR(180)" json:"item_name,omitempty"`
	ItemType    *string    `gorm:"column:Item_Type;type:VARCHAR(50)" json:"item_type,omitempty"`
	IUPACName   *string    `gorm:"column:IUPAC_name;type:VARCHAR(150)" json:"iupac_name,omitempty"`
	Description *string    `gorm:"column:Description;type:VARCHAR(100)" json:"description,omitempty"`
	CASNumber   *string    `gorm:"column:CAS_number;type:VARCHAR(30)" json:"cas_number,omitempty"`
	ItemTag     *string    `gorm:"column:Item_Tag;type:VARCHAR(30)" json:"item_tag,omitempty"`
	Formula     *string    `gorm:"column:Formula;type:VARCHAR(80)" json:"formula,omitempty"`
	Structure   *string    `gorm:"column:Structure;type:TEXT" json:"structure,omitempty"`
	MS1         *string    `gorm:"column:MS1;type:VARCHAR(10)" json:"ms1,omitempty"`
	MS2         *string    `gorm:"column:MS2;type:VARCHAR(10)" json:"ms2,omitempty"`
	Bioactivity *string    `gorm:"column:Bioactivity;type:VARCHAR(10)" json:"bioactivity,omitempty"`
	Smiles      string     `gorm:"column:Smiles;type:TEXT;not null" json:"smiles"`
	CreatedAt   *time.Time `gorm:"column:Created_At" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"column:Updated_At" json:"updated_at,omitempty"`
}

// TableName 指定表名
func (Data) TableName() string {
	return "data"
}
