package user

import (
	fmt "fmt"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("failed to execute before create hook on user: %s", err)
	}

	return scope.SetColumn("Id", uuid.String())
}
