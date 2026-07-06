package domain

import (
	"database/sql/driver"
	"time"

	"github.com/codenomdev/viona/pkg/util"
	"github.com/gookit/goutil/strutil"
	"gorm.io/gorm"
)

const TableName string = "users"

// Status user type
type StatusUser int

const (
	InActive StatusUser = iota
	Active
	Suspend
	Blocked
)

type StatusUserString string

const (
	StatusInActive StatusUserString = "inactive"
	StatusActive   StatusUserString = "active"
	StatusSuspend  StatusUserString = "suspend"
	StatusBlocked  StatusUserString = "blocked"
	StatusUnknown  StatusUserString = "unknown"
)

type User struct {
	ID              int64      `gorm:"primaryKey;autoIncrement:false;"`
	FullName        string     `gorm:"column:fullname;size:255;not null;check:LENGTH(fullname) >=4;"`
	Username        string     `gorm:"size:50;uniqueIndex;check:LENGTH(username) >=4;not null;"`
	Email           string     `gorm:"size:255;not null;uniqueIndex;"`
	Password        string     `gorm:"column:password;size:255;not null;check:LENGTH(password) >= 8;"`
	Avatar          string     `gorm:"column:avatar;size:255;default:null;"`
	IsActive        StatusUser `gorm:"column:is_active;type:int;size:1;default:1;not null;index:,type:btree;"`
	EmailVerifiedAt *time.Time

	LastLoginDate time.Time `gorm:"column:last_login_date;"`
	DeletedAt     gorm.DeletedAt
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (User) TableName() string {
	return TableName
}

// Hook gorm global user before create.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Lower case email
	if err := u.GenerateEmailLower(); err != nil {
		return err
	}

	// Hash password before create
	if err := u.GeneratePasswordHash(); err != nil {
		return err
	}

	return nil
}

func (st *StatusUser) Scan(value any) error {
	b, ok := value.(int64)
	if !ok {
		*st = StatusUser(b)
	}
	*st = StatusUser(int(b))
	return nil
}

func (st *StatusUser) Value() (driver.Value, error) {
	return int(*st), nil
}

// util generate email to lower case
func (u *User) GenerateEmailLower() error {
	u.Email = strutil.Lowercase(u.Email)
	return nil
}

// util generate password hash.
func (u *User) GeneratePasswordHash() error {
	crypto := util.NewCryptoHash()

	hash, err := crypto.CreatePasswordHash(u.Password)
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

// Status User to string.
func (u *User) GetStatusUser() string {
	switch u.IsActive {
	case Active:
		return string(StatusActive)
	case Suspend:
		return string(StatusSuspend)
	case Blocked:
		return string(StatusBlocked)
	case InActive:
		return string(StatusInActive)
	default:
		return string(StatusUnknown)
	}
}
