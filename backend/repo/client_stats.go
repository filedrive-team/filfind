package repo

import (
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

func (m *Manager) UpsertClientStats(list []*models.ClientStats) error {
	return m.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "client"}},
		DoUpdates: clause.AssignmentColumns([]string{"storage_capacity", "storage_deals",
			"used_data_cap", "verified_deals", "section_deal_id"}),
	}).Create(list).Error
}

func (m *Manager) LastClientStats() (stats *models.ClientStats, err error) {
	stats = &models.ClientStats{}
	err = m.db.Model(&models.ClientStats{}).Order("id desc").First(stats).Error
	return
}

type ClientStats struct {
	Client          string
	UsedDataCap     decimal.Decimal
	VerifiedDeals   int64
	StorageCapacity decimal.Decimal
	StorageDeals    int64
}

func (m *Manager) GetClientStatsSum() (list []*ClientStats, err error) {
	stmt := `
select cs.client,
sum(cs.storage_capacity) as storage_capacity,
sum(cs.storage_deals) as storage_deals,
sum(cs.used_data_cap) as used_data_cap,
sum(cs.verified_deals) as verified_deals
from client_stats cs
where cs.client in (select ci.address_id from client_info ci)
group by cs.client
`
	err = m.db.Raw(stmt).Scan(&list).Error
	return list, err
}
