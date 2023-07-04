package repo

import (
	"errors"
	"fmt"
	"messaging-service/types/records"

	goerrors "github.com/go-errors/errors"
	"gorm.io/gorm"
)

func (r *Repo) GetAuthProfileByEmail(email string) (*records.AuthProfile, error) {
	// result := &records.AuthProfile{}
	result := &records.AuthProfile{}

	res := r.DB.Where("email = ? ", email).First(result)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if res.Error != nil {
		return nil, goerrors.Wrap(res.Error, 0)
	}
	return result, nil
}

func (r *Repo) SaveAuthProfile(authProfile *records.AuthProfile) error {
	err := r.DB.Create(authProfile).Error
	if err != nil {
		fmt.Println("RETURNING SAVE ERROR")
		return goerrors.Wrap(err, 0)
	}
	return nil
}
