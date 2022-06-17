package repo

import (
	"github.com/filedrive-team/filfind/backend/models"
)

func (m *Manager) CreateAccessLog(logs []*models.AccessLog) error {
	return m.db.Create(logs).Error
}

func (m *Manager) AccessesLastDay() (count int64, err error) {
	err = m.db.Model(&models.AccessLog{}).
		Where("created_at>=DATE_FORMAT(NOW(),'%Y-%m-%d 00:00:00')").
		Count(&count).Error
	return
}

type AccessesDaily struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

func (m *Manager) AccessesDaily() (list []*AccessesDaily, err error) {
	err = m.db.Model(&models.AccessLog{}).
		Where("al.created_at >=CONVERT(DATE_SUB(NOW(),interval 1 month), date)").
		Select("CONVERT(created_at, date) as date,COUNT(id) as value").
		Group("date").
		Scan(&list).Error
	return
}
