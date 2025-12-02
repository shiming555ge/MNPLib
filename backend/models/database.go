package models

import (
	"time"
)

// Data 对应数据库中的 data 表
type Data struct {
	ID           string     `gorm:"column:ID;type:VARCHAR(12);primaryKey;not null" json:"id"`
	Source       *string    `gorm:"column:Source;type:VARCHAR(255)" json:"source,omitempty"`
	ItemName     *string    `gorm:"column:ItemName;type:TEXT" json:"item_name,omitempty"`
	ItemType     *string    `gorm:"column:ItemType;type:TEXT" json:"item_type,omitempty"`
	Formula      *string    `gorm:"column:Formula;type:VARCHAR(128)" json:"formula,omitempty"`
	SMILES       *string    `gorm:"column:SMILES;type:TEXT" json:"smiles,omitempty"`
	Description  *string    `gorm:"column:Description;type:ENUM('KNOWN COMPOUND','NEW NATURAL PRODUCT','NEW ANALOGS')" json:"description,omitempty"`
	CASNumber    *string    `gorm:"column:CAS_number;type:VARCHAR(100)" json:"cas_number,omitempty"`
	ItemTag      *string    `gorm:"column:ItemTag;type:VARCHAR(255)" json:"item_tag,omitempty"`
	Structure    *string    `gorm:"column:Structure;type:TEXT" json:"structure,omitempty"`
	MS1          *float64   `gorm:"column:MS1;type:DOUBLE" json:"ms1,omitempty"`
	MS2          *string    `gorm:"column:MS2;type:VARCHAR(512)" json:"ms2,omitempty"`
	Bioactivity  *string    `gorm:"column:Bioactivity;type:VARCHAR(512)" json:"bioactivity,omitempty"`
	NMR_13C_data *string    `gorm:"column:NMR_13C_data;type:TEXT" json:"nmr_13c_data,omitempty"`
	Weight       *float32   `gorm:"column:Weight;type:FLOAT" json:"weight,omitempty"`
	FP           *string    `gorm:"column:FP;type:VARCHAR(127)" json:"fp,omitempty"`
	CreatedAt    *time.Time `gorm:"column:Created_At" json:"created_at,omitempty"`
	UpdatedAt    *time.Time `gorm:"column:Updated_At" json:"updated_at,omitempty"`
}

type PublicData struct {
	ID          string     `gorm:"column:ID;type:VARCHAR(12);primaryKey;not null" json:"id"`
	Source      *string    `gorm:"column:Source;type:VARCHAR(255)" json:"source,omitempty"`
	ItemName    *string    `gorm:"column:ItemName;type:TEXT" json:"item_name,omitempty"`
	ItemType    *string    `gorm:"column:ItemType;type:TEXT" json:"item_type,omitempty"`
	Formula     *string    `gorm:"column:Formula;type:VARCHAR(128)" json:"formula,omitempty"`
	SMILES      *string    `gorm:"column:SMILES;type:TEXT" json:"smiles,omitempty"`
	Description *string    `gorm:"column:Description;type:ENUM('KNOWN COMPOUND','NEW NATURAL PRODUCT','NEW ANALOGS')" json:"description,omitempty"`
	CASNumber   *string    `gorm:"column:CAS_number;type:VARCHAR(100)" json:"cas_number,omitempty"`
	ItemTag     *string    `gorm:"column:ItemTag;type:VARCHAR(255)" json:"item_tag,omitempty"`
	Structure   *string    `gorm:"column:Structure;type:TEXT" json:"structure,omitempty"`
	MS1         *float64   `gorm:"column:MS1;type:DOUBLE" json:"ms1,omitempty"`
	Weight      *float32   `gorm:"column:Weight;type:FLOAT" json:"weight,omitempty"`
	FP          *string    `gorm:"column:FP;type:VARCHAR(127)" json:"fp,omitempty"`
	CreatedAt   *time.Time `gorm:"column:Created_At" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"column:Updated_At" json:"updated_at,omitempty"`
}

// 定义只包含保护字段的结构
type ProtectedData struct {
	MS2          *string `json:"ms2,omitempty"`
	Bioactivity  *string `json:"bioactivity,omitempty"`
	NMR_13C_data *string `json:"nmr_13c_data,omitempty"`
}

// TableName 指定表名
func (Data) TableName() string {
	return "data"
}
