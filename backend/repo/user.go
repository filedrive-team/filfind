package repo

import (
	"github.com/filedrive-team/filfind/backend/models"
	uuid "github.com/satori/go.uuid"
)

func (m *Manager) CreateUser(u *models.User) error {
	u.Uid = uuid.NewV4()
	return m.db.Create(u).Error
}

func (m *Manager) ExistUser(userType string, address string) (bool, error) {
	var count int64
	err := m.db.Model(models.User{}).Where("address_robust=? and type=?", address, userType).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *Manager) ExistUserByEmail(eamil string) (bool, error) {
	var count int64
	err := m.db.Model(models.User{}).Where("email=?", eamil).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *Manager) ExistUserByUid(uid string) (bool, error) {
	var count int64
	err := m.db.Model(models.User{}).Where("uid=?", uid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *Manager) QueryUserByEmail(email string) (u *models.User, err error) {
	data := &models.User{}
	err = m.db.Model(data).Where("email=?", email).First(data).Error
	if err != nil {
		return
	}
	return data, err
}

func (m *Manager) QueryUserByUid(uid string) (u *models.User, err error) {
	data := &models.User{}
	err = m.db.Model(data).Where("uid=?", uid).First(data).Error
	if err != nil {
		return
	}
	return data, err
}

func (m *Manager) UpdateUserPassword(u *models.User) (err error) {
	return m.db.Model(u).Update("password", u.HashedPassword).Where("uid=?", u.Uid).Error
}

type Profile struct {
	Uid           string `json:"uid"`
	Type          string `json:"type"`
	AddressRobust string `json:"address_robust"`
	AddressId     string `json:"address_id"`
	Name          string `json:"name"`
	Avatar        string `json:"avatar"`
	Logo          string `json:"logo"`
	Location      string `json:"location"`
	ContactEmail  string `json:"contact_email"`
	Slack         string `json:"slack"`
	Github        string `json:"github"`
	Twitter       string `json:"twitter"`
	Description   string `json:"description"`
}

type ClientProfile struct {
	Profile
}

type SPOwnerProfile struct {
	Profile
	ReputationScore float64 `json:"reputation_score"`
	ReviewScore     float64 `json:"review_score"`
	Reviews         uint64  `json:"reviews"`
}

func (m *Manager) ClientProfile(addrId string) (profile *ClientProfile, err error) {
	profile = &ClientProfile{}
	err = m.db.Model(&models.User{}).
		Where("address_id=? and type=?", addrId, models.ClientRole).
		First(profile).Error
	return profile, err
}

func (m *Manager) SPOwnerProfile(addrId string) (profile *SPOwnerProfile, err error) {
	profile = &SPOwnerProfile{}
	stmt := `
select
*
from (
select
p.owner,
avg(p.score) as reputation_score,
avg(p.review_score) as review_score,
sum(p.reviews) as reviews
from provider p 
where p.owner = ?
group by 1
) t
left join user u 
on t.owner = u.address_id and u.type=?
`
	err = m.db.Raw(stmt, addrId, models.SPOwnerRole).Scan(profile).Error
	return profile, err
}

func (m *Manager) UpdateProfile(u *models.User) error {
	return m.db.Model(u).Where("uid=?", u.Uid).Updates(u).Error
}

type ClientAddress struct {
	Uid       uuid.UUID `json:"uid"`
	AddressId string    `json:"address_id"`
}

func (m *Manager) ClientUserList() (list []*ClientAddress, err error) {
	err = m.db.Model(&models.User{}).
		Where("type=?", models.ClientRole).
		Select("uid,address_id").
		Scan(&list).Error
	return
}
