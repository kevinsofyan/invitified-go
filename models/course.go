package models

import "os"

type Course struct {
	ID           int    `json:"id" gorm:"column:course_id;primaryKey"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DepartmentID int    `json:"department_id"`
	Credits      int    `json:"credits"`
}

func (Course) TableName() string {
	schema := os.Getenv("DB_SCHEMA")
	if schema == "" {
		schema = "public"
	}
	return schema + ".courses"
}
