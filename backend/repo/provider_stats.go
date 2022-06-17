package repo

import (
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/utils/utils"
	"gorm.io/gorm/clause"
	"time"
)

func (m *Manager) UpsertProviderStatsMonthly(list []*models.ProviderStatsMonthly) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "month"}, {Name: "provider"}, {Name: "client"}},
		UpdateAll: true,
	}).Create(list).Error
}

func (m *Manager) LastProviderStatsMonthly() (stats *models.ProviderStatsMonthly, err error) {
	stats = &models.ProviderStatsMonthly{}
	err = m.db.Model(&models.ProviderStatsMonthly{}).Order("id desc").First(stats).Error
	return
}

func (m *Manager) SpsToDealNewClientLastMonth() (count int64, err error) {
	begin := utils.MonthBegin(time.Now())
	end := begin.AddDate(0, 1, 0)
	return m.SpsToDealNewClientByRange(begin, end)
}

func (m *Manager) SpsToDealNewClientByRange(startMonth, endMonth time.Time) (count int64, err error) {
	stmt := `
select count(distinct provider) from (
select t1.provider,t1.client,t2.repeats from (
select psm.provider,psm.client
from provider_stats_monthly psm
where psm.month >= ? and psm.month < ?
group by psm.provider,psm.client
) t1 left join (
select psm.provider,psm.client,1 as repeats
from provider_stats_monthly psm
where psm.month < ?
group by psm.provider,psm.client
) t2 on t1.provider=t2.provider and t1.client=t2.client
) t3 where repeats is null
`
	err = m.db.Raw(stmt, startMonth, endMonth, startMonth).Scan(&count).Error
	return
}
