package repo

import (
	"messaging-service/types/records"

	goerrors "github.com/go-errors/errors"
)

func (r *Repo) GetAuthProfileByEmail(email string) (*records.AuthProfile, error) {
	// result := &records.AuthProfile{}
	result := &records.AuthProfile{}

	res := r.DB.Where("email = ? ", email).Find(result)
	// if errors.Is(res.Error, gorm.ErrRecordNotFound) {
	if res.RowsAffected == 0 {
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
		return goerrors.Wrap(err, 0)
	}
	return nil
}

func (r *Repo) UpdatePassword(email string, hashedPassword string) error {
	err := r.DB.Model(&records.AuthProfile{}).Where("email = ?", email).Updates(&records.AuthProfile{HashedPassword: hashedPassword}).Error
	if err != nil {
		return goerrors.Wrap(err, 0)
	}
	return nil
}
