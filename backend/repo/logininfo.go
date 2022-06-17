package repo

import "github.com/filedrive-team/filfind/backend/models"

func (m *Manager) CreateLoginInfo(info *models.LoginInfo) error {
	return m.db.Create(info).Error
}
