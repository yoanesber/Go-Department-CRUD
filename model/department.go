package model

import "time"

type Department struct {
	ID          string    `db:"id" json:"id"`
	DeptName    string    `db:"dept_name" json:"deptName"`
	Active      bool      `db:"active" json:"active"`
	CreatedBy   int64     `db:"created_by" json:"createdBy"`
	CreatedDate time.Time `db:"created_date" json:"createdDate"`
	UpdatedBy   int64     `db:"updated_by" json:"updatedBy"`
	UpdatedDate time.Time `db:"updated_date" json:"updatedDate"`
}
