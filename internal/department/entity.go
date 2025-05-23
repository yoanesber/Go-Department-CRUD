package department

import (
	"time"

	validate "github.com/yoanesber/Go-Department-CRUD/pkg/validator"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

var v *validator.Validate

// Department represents the department entity in the database.
type Department struct {
	ID        string          `gorm:"column:id;type:varchar(4);primaryKey;not null" json:"id" validate:"required,len=4"`
	DeptName  string          `gorm:"column:dept_name;type:varchar(40);unique;not null" json:"deptName" validate:"required,max=40"`
	Active    bool            `gorm:"column:active;type:bool;not null" json:"active"`
	CreatedBy *int64          `gorm:"column:created_by" json:"createdBy,omitempty"`
	CreatedAt *time.Time      `gorm:"column:created_at;type:timestamptz;autoCreateTime;default:now()" json:"createdAt,omitempty"`
	UpdatedBy *int64          `gorm:"column:updated_by" json:"updatedBy,omitempty"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;type:timestamptz;autoUpdateTime;default:now()" json:"updatedAt,omitempty"`
	DeletedBy *int64          `gorm:"column:deleted_by" json:"deletedBy,omitempty"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz;index" json:"deletedAt,omitempty"`
}

// Override the TableName method to specify the table name
// in the database. This is optional if you want to use the default naming convention.
func (Department) TableName() string {
	return "department"
}

// Equals compares two Department objects for equality.
func (d *Department) Equals(other *Department) bool {
	if d == nil && other == nil {
		return true
	}

	if d == nil || other == nil {
		return false
	}

	if (d.ID != other.ID) ||
		(d.DeptName != other.DeptName) ||
		(d.Active != other.Active) {
		return false
	}

	return true
}

// EqualsIgnoreID compares two Department objects for equality,
// ignoring the ID field. This is useful for comparing
// Department objects where the ID is not relevant.
func (d *Department) EqualsIgnoreID(other *Department) bool {
	if d == nil && other == nil {
		return true
	}

	if d == nil || other == nil {
		return false
	}

	if (d.DeptName != other.DeptName) ||
		(d.Active != other.Active) {
		return false
	}

	return true
}

// Validate validates the Department struct using the validator package.
// It checks if the struct fields meet the validation rules defined in the struct tags.
func (d *Department) Validate() error {
	v = validate.GetValidator()

	if err := v.Struct(d); err != nil {
		return err
	}

	return nil
}
