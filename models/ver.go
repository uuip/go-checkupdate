package models

import (
	"database/sql"
	"encoding/json"
)

type VerModel struct {
	Name       string     `gorm:"primaryKey" json:"name"`
	Ver        string     `json:"ver"`
	Url        string     `json:"url"`
	Newversion NullString `json:"newversion"`
	Json       int8       `json:"json"`
}

type Tabler interface {
	TableName() string
}

func (VerModel) TableName() string {
	return "ver_tab"
}

type NullString struct {
	sql.NullString
}

func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.String = *s
	} else {
		v.Valid = false
	}
	return nil
}
