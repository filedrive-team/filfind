package repo

import (
	"fmt"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/types"
)

func (m *Manager) CreateMessage(msg *models.Message) error {
	return m.db.Create(msg).Error
}

func (m *Manager) MessagesByUser(oneUid, anotherUid string, before *types.UnixTime, limit int) (list []*models.Message, err error) {
	stmtParams := map[string]interface{}{
		"one":     oneUid,
		"another": anotherUid,
		"limit":   limit,
	}

	subStmt := `
(select * from message m where m.sender=@one and m.recipient=@another %s order by m.created_at desc limit @limit)
union
(select * from message m where m.sender=@another and m.recipient=@one %s order by m.created_at desc limit @limit)
`
	beforeCond := ""
	if before != nil {
		stmtParams["before"] = before
		beforeCond = "and m.created_at<@before"
	}
	subStmt = fmt.Sprintf(subStmt, beforeCond, beforeCond)
	stmt := fmt.Sprintf(`
select * from 
(%s) t
order by t.created_at desc
limit @limit
`, subStmt)

	err = m.db.Raw(stmt, stmtParams).Scan(&list).Error
	return list, err
}

func (m *Manager) UpdateMessagesChecked(sender, recipient string) error {
	return m.db.Model(&models.Message{}).
		Where("sender=? and recipient=? and not checked", sender, recipient).
		Updates(&models.Message{
			Checked: true,
		}).Error
}

func (m *Manager) UncheckedMessagesByRecipient(recipient string) (total int64, err error) {
	err = m.db.Model(&models.Message{}).Where("recipient=? and not checked", recipient).Count(&total).Error
	return
}

type UncheckedMessageItem struct {
	Partner string `json:"partner"`
	Number  int64  `json:"number"`
}

func (m *Manager) UncheckedMessageGroupListByRecipient(recipient string) (list []*UncheckedMessageItem, err error) {
	err = m.db.Model(&models.Message{}).
		Where("recipient=? and not checked", recipient).
		Group("sender").
		Select("sender as partner,count(*) as number").
		Scan(&list).Error
	return
}

func (m *Manager) UncheckedMessagesByPair(recipient, sender string) (list []*UncheckedMessageItem, err error) {
	err = m.db.Model(&models.Message{}).
		Where("sender=? and recipient=? and not checked", sender, recipient).
		Group("sender").
		Select("sender as partner,count(*) as number").
		Scan(&list).Error
	return
}

func (m *Manager) ContactsClientSpLastMonth() (count int64, err error) {
	err = m.db.Model(&models.Message{}).
		Where("created_at>=DATE_FORMAT(CURDATE(), '%Y-%m-01 00:00:00')").
		Count(&count).Error
	return
}

type ContactsMonthly struct {
	Month string `json:"month"`
	Value int64  `json:"value"`
}

func (m *Manager) ContactsClientSpMonthly() (list []*AccessesDaily, err error) {
	err = m.db.Model(&models.Message{}).
		Select("DATE_FORMAT(created_at, '%Y%m') as month,COUNT(id) as value").
		Group("month").
		Scan(&list).Error
	return
}
