package user

import (
	"time"

	"github.com/yoanesber/Go-Department-CRUD/internal/refreshtoken"
	"github.com/yoanesber/Go-Department-CRUD/internal/role"
	validate "github.com/yoanesber/Go-Department-CRUD/pkg/validator"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

var v *validator.Validate

// User represents the user entity in the database.
type User struct {
	ID                        int64                      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserName                  string                     `gorm:"column:username;type:varchar(20);not null;unique" json:"userName" validate:"required,min=3,max=20"`
	Password                  string                     `gorm:"column:password;type:varchar(150);not null" json:"password" validate:"required,min=8"`
	Email                     string                     `gorm:"column:email;type:varchar(100);not null;unique" json:"email" validate:"required,email,max=100"`
	FirstName                 string                     `gorm:"column:firstname;type:varchar(20);not null" json:"firstName" validate:"required,max=20"`
	LastName                  *string                    `gorm:"column:lastname;type:varchar(20)" json:"lastName,omitempty" validate:"omitempty,max=20"`
	IsEnabled                 *bool                      `gorm:"column:is_enabled;not null;default:false" json:"isEnabled,omitempty"`
	IsAccountNonExpired       *bool                      `gorm:"column:is_account_non_expired;not null;default:false" json:"isAccountNonExpired,omitempty"`
	IsAccountNonLocked        *bool                      `gorm:"column:is_account_non_locked;not null;default:false" json:"isAccountNonLocked,omitempty"`
	IsCredentialsNonExpired   *bool                      `gorm:"column:is_credentials_non_expired;not null;default:false" json:"isCredentialsNonExpired,omitempty"`
	IsDeleted                 *bool                      `gorm:"column:is_deleted;not null;default:false" json:"isDeleted,omitempty"`
	AccountExpirationDate     *time.Time                 `gorm:"column:account_expiration_date;type:timestamptz" json:"accountExpirationDate,omitempty"`
	CredentialsExpirationDate *time.Time                 `gorm:"column:credentials_expiration_date;type:timestamptz" json:"credentialsExpirationDate,omitempty"`
	UserType                  string                     `gorm:"column:user_type;type:varchar(20);not null;check:user_type IN ('SERVICE_ACCOUNT','USER_ACCOUNT')" json:"userType" validate:"required,max=20,oneof=SERVICE_ACCOUNT USER_ACCOUNT"`
	LastLogin                 *time.Time                 `gorm:"column:last_login" json:"lastLogin,omitempty"`
	CreatedBy                 *int64                     `gorm:"column:created_by" json:"createdBy,omitempty"`
	CreatedAt                 *time.Time                 `gorm:"column:created_at;type:timestamptz;autoCreateTime;default:now()" json:"createdAt,omitempty"`
	UpdatedBy                 *int64                     `gorm:"column:updated_by" json:"updatedBy,omitempty"`
	UpdatedAt                 *time.Time                 `gorm:"column:updated_at;type:timestamptz;autoUpdateTime;default:now()" json:"updatedAt,omitempty"`
	DeletedBy                 *int64                     `gorm:"column:deleted_by" json:"deletedBy,omitempty"`
	DeletedAt                 *gorm.DeletedAt            `gorm:"column:deleted_at;type:timestamptz;index" json:"deletedAt,omitempty"`
	Roles                     []role.Role                `gorm:"many2many:user_roles;constraint:OnUpdate:RESTRICT,OnDelete:SET NULL" json:"roles,omitempty"`
	RefreshToken              *refreshtoken.RefreshToken `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"refreshToken,omitempty"`
}

// Override the TableName method to specify the table name
// in the database. This is optional if you want to use the default naming convention.
func (User) TableName() string {
	return "users"
}

// Equals compares two User objects for equality.
func (u *User) Equals(other *User) bool {
	if u == nil && other == nil {
		return true
	}

	if u == nil || other == nil {
		return false
	}

	if (u.ID != other.ID) ||
		(u.UserName != other.UserName) ||
		(u.Password != other.Password) ||
		(u.Email != other.Email) ||
		(u.FirstName != other.FirstName) ||
		(u.LastName != other.LastName) ||
		(u.IsEnabled != other.IsEnabled) ||
		(u.IsAccountNonExpired != other.IsAccountNonExpired) ||
		(u.IsAccountNonLocked != other.IsAccountNonLocked) ||
		(u.IsCredentialsNonExpired != other.IsCredentialsNonExpired) ||
		(u.IsDeleted != other.IsDeleted) ||
		(u.AccountExpirationDate != other.AccountExpirationDate) ||
		(u.CredentialsExpirationDate != other.CredentialsExpirationDate) ||
		(u.UserType != other.UserType) ||
		(u.LastLogin != other.LastLogin) {

		return false
	}

	return true
}

// Validate validates the User struct using the validator package.
// It checks if the struct fields meet the specified validation rules.
func (u *User) Validate() error {
	v = validate.GetValidator()

	if err := v.Struct(u); err != nil {
		return err
	}
	return nil
}
