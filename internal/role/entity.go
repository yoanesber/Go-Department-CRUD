package role

import (
	validate "github.com/yoanesber/Go-Department-CRUD/pkg/validator"
	"gopkg.in/go-playground/validator.v9"
)

var v *validator.Validate

// Role represents the role entity in the database.
type Role struct {
	ID   uint   `gorm:"column:id;primaryKey;autoIncrement" json:"roleId"`
	Name string `gorm:"column:name;type:varchar(20);not null;check:name IN ('ROLE_USER','ROLE_MODERATOR','ROLE_ADMIN')" json:"roleName" validate:"required,max=20,oneof=ROLE_USER ROLE_MODERATOR ROLE_ADMIN"`
}

// UserRole represents the many-to-many relationship between users and roles.
type UserRole struct {
	UserID int64 `gorm:"column:user_id;primaryKey;not null"`
	RoleID int   `gorm:"column:role_id;primaryKey;not null"`
}

// Override the TableName method to specify the table name
// in the database. This is optional if you want to use the default naming convention.
func (Role) TableName() string {
	return "roles"
}

// Override the TableName method to specify the table name
// in the database. This is optional if you want to use the default naming convention.
func (UserRole) TableName() string {
	return "user_roles"
}

// Validate validates the Role struct using the validator package.
// It checks if the struct fields meet the specified validation rules.
func (r *Role) Validate() error {
	v = validate.GetValidator()

	if err := v.Struct(r); err != nil {
		return err
	}
	return nil
}

// Equals compares two Role objects for equality.
func (r *Role) Equals(other *Role) bool {
	if r == nil && other == nil {
		return true
	}

	if r == nil || other == nil {
		return false
	}

	if (r.ID != other.ID) ||
		(r.Name != other.Name) {
		return false
	}

	return true
}
