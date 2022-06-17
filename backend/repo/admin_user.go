package repo

import (
	"github.com/filedrive-team/filfind/backend/models"
	uuid "github.com/satori/go.uuid"
)

func (m *Manager) CreateAdminUser(u *models.AdminUser) error {
	u.Uid = uuid.NewV4()
	u.Type = models.AdminRole
	return m.db.Create(u).Error
}

func (m *Manager) ExistAdminUserByName(name string) (bool, error) {
	var count int64
	err := m.db.Model(models.AdminUser{}).Where("name=?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *Manager) ExistAdminUserByUid(uid string) (bool, error) {
	var count int64
	err := m.db.Model(models.AdminUser{}).Where("uid=?", uid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *Manager) QueryAdminUserByName(name string) (u *models.AdminUser, err error) {
	data := &models.AdminUser{}
	err = m.db.Model(data).Where("name=?", name).First(data).Error
	if err != nil {
		return
	}
	return data, err
}

func (m *Manager) QueryAdminUserByUid(uid string) (u *models.AdminUser, err error) {
	data := &models.AdminUser{}
	err = m.db.Model(data).Where("uid=?", uid).First(data).Error
	if err != nil {
		return
	}
	return data, err
}

func (m *Manager) UpdateAdminUserPassword(u *models.AdminUser) (err error) {
	return m.db.Model(u).Update("password", u.HashedPassword).Where("uid=?", u.Uid).Error
}
