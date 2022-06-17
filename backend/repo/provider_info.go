package repo

import (
	"github.com/filedrive-team/filfind/backend/models"
	"gorm.io/gorm/clause"
)

func (m *Manager) UpsertProviderInfo(info *models.ProviderInfo) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		UpdateAll: true,
	}).Create(info).Error
}
