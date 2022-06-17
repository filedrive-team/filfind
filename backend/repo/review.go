package repo

import (
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/filedrive-team/filfind/backend/types"
	"github.com/filedrive-team/filfind/backend/utils/utils"
)

func (m *Manager) InsertReview(review *models.Review) error {
	return m.db.Create(review).Error
}

type Review struct {
	CreatedAt types.UnixTime `json:"created_at"`
	Client    string         `json:"client"`
	Provider  string         `json:"provider"`
	Score     int            `json:"score"`
	Content   string         `json:"content"`
	Title     string         `json:"title"`
}

func (m *Manager) ReviewsByClient(pagination types.PaginationParams, clientId string) (total int, list []*Review, err error) {
	offset, size := utils.PaginationHelper(pagination.Page, pagination.PageSize, settings.DefaultPageSize)

	stmtCount := `
select 
count(*) 
from review 
where client=?
	`
	stmt := `
select 
* 
from review 
where client=? 
order by id desc
limit ? offset ?
	`

	err = m.db.Raw(stmtCount, clientId).Scan(&total).Error
	if err != nil {
		return total, nil, err
	}

	err = m.db.Raw(stmt, clientId, size, offset).Scan(&list).Error
	return total, list, err
}

func (m *Manager) ReviewsBySpOwner(pagination types.PaginationParams, ownerId string) (total int, list []*Review, err error) {
	offset, size := utils.PaginationHelper(pagination.Page, pagination.PageSize, settings.DefaultPageSize)

	stmtCount := `
select 
count(*) 
from review 
where provider in (select address from provider where owner=?)
	`
	stmt := `
select 
* 
from review 
where provider in (select address from provider where owner=?)
order by id desc
limit ? offset ?
	`

	err = m.db.Raw(stmtCount, ownerId).Scan(&total).Error
	if err != nil {
		return total, nil, err
	}

	err = m.db.Raw(stmt, ownerId, size, offset).Scan(&list).Error
	return total, list, err
}

type ReviewStatsItem struct {
	Provider string  `json:"provider"`
	AvgScore float64 `json:"avg_score"`
	Reviews  uint64  `json:"reviews"`
}

func (m *Manager) StatsReviews() (list []*ReviewStatsItem, err error) {
	err = m.db.Model(models.Review{}).
		Group("provider").
		Select("provider,avg(score) as avg_score,count(*) as reviews").
		Scan(&list).Error
	return
}
