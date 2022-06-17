package repo

import (
	"github.com/filedrive-team/filfind/backend/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm/clause"
)

func (m *Manager) CreateRelationship(oneUid, otherUid string) error {
	if oneUid == otherUid {
		return nil
	}
	one, err := uuid.FromString(oneUid)
	if err != nil {
		return err
	}
	other, err := uuid.FromString(otherUid)
	if err != nil {
		return err
	}
	urs := []*models.UserRelation{
		{
			Uid:     one,
			Partner: other,
		},
		{
			Uid:     other,
			Partner: one,
		},
	}
	return m.UpsertUserRelation(urs)
}

func (m *Manager) UpsertUserRelation(urs []*models.UserRelation) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}, {Name: "partner"}},
		UpdateAll: false,
	}).Create(urs).Error
}

func (m *Manager) HasRelationByUser(userUid, partnerUid string) (has bool, err error) {
	var count int64
	err = m.db.Model(models.UserRelation{}).
		Where("uid=? and partner=?", userUid, partnerUid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

type Partner struct {
	Uid          string `json:"uid"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	AddressId    string `json:"address_id"`
	Avatar       string `json:"avatar"`
	Location     string `json:"location"`
	ContactEmail string `json:"contact_email"`
	Slack        string `json:"slack"`
	Github       string `json:"github"`
	Twitter      string `json:"twitter"`
	Description  string `json:"description"`
}

func (m *Manager) PartnersByUid(uid string) (list []*Partner, err error) {
	stmt := `
select 
u.* 
from
(
select ur.partner 
from user_relation ur 
where ur.uid=?
) t
left join user u on t.partner=u.uid
order by id desc
`
	err = m.db.Raw(stmt, uid).Scan(&list).Error
	return list, err
}
