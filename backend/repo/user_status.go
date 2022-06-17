package repo

import (
	"github.com/filedrive-team/filfind/backend/models"
	"gorm.io/gorm/clause"
)

func (m *Manager) UpsertUserOnline(uo *models.UserStatus) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}},
		DoUpdates: clause.AssignmentColumns([]string{"online"}),
	}).Create(uo).Error
}

type PartnerStatus struct {
	Uid    string `json:"uid"`
	Online bool   `json:"online"`
}

func (m *Manager) PartnersStatusByUid(uid string) (list []*PartnerStatus, err error) {
	stmt := `
select 
t.partner as uid,
COALESCE(uo.online),0) as online
from
(
select ur.partner 
from user_relation ur 
where ur.uid=?
) t
left join user_online uo on t.partner=uo.uid
`
	err = m.db.Raw(stmt, uid).Scan(&list).Error
	return list, err
}
