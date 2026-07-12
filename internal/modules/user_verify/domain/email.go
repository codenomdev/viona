package domain

import (
	"fmt"
	"time"

	"github.com/codenomdev/viona/pkg/util"
	"github.com/gookit/goutil/strutil"
	"gorm.io/gorm"
)

const TableName string = "user_verify_email"

type EmailVerification struct {
	ID         int64     `gorm:"primaryKey;autoIncrement:false;"`
	UserID     int64     `gorm:"not null;index"`
	TokenHash  string    `gorm:"size:255;not null;uniqueIndex;"`
	ExpiresAt  time.Time `gorm:"not null"`
	VerifiedAt *time.Time
	IPAddress  string `gorm:"size:45"`
	UserAgent  string `gorm:"size:500"`
	CreatedAt  time.Time
}

func (EmailVerification) TableName() string {
	return TableName
}

// before create
func (e *EmailVerification) BeforeCreate(tx *gorm.DB) error {
	if e.TokenHash == "" {
		e.GenerateRandomToken()
	}
	if e.ExpiresAt.IsZero() {
		e.GenerateExpiredAt()
	}
	return nil
}

func (e *EmailVerification) GenerateRandomToken() error {
	tokenTpl := fmt.Sprintf("%s%s", strutil.AlphaBet1, util.ParseInt64String(e.UserID))
	e.TokenHash = util.GenerateStrTpl(32, tokenTpl)
	return nil
}

// Generate token expired
// expired default 2 hours.
func (e *EmailVerification) GenerateExpiredAt() error {
	e.ExpiresAt = util.GenerateTimeHours()
	return nil
}
