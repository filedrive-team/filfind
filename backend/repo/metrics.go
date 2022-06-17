package repo

import (
	"fmt"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/filedrive-team/filfind/backend/types"
	"github.com/filedrive-team/filfind/backend/utils/utils"
	"github.com/shopspring/decimal"
	"time"
)

type MetricsOverview struct {
	RegisteredProviders int64           `json:"registered_providers"`
	AutoFilledProviders int64           `json:"auto_filled_providers"`
	RegisteredSpRatio   decimal.Decimal `json:"registered_sp_ratio"`

	InternalContacts     int64           `json:"internal_contacts"`
	AverageAccessesDaily decimal.Decimal `json:"average_accesses_daily"`
	TotalAccesses        int64           `json:"total_accesses"`
}

func (m *Manager) MetricsOverview() (*MetricsOverview, error) {
	var res MetricsOverview
	var manager utils.AsyncManager
	errors := make([]error, 0)
	manager.AddTask(func() {
		sql := `
select COALESCE(sum(t.num),0) as registered_providers, 
count(*)-COALESCE(sum(t.num),0) as auto_filled_providers,
COALESCE(sum(t.num),0)/(count(*)-COALESCE(sum(t.num),0)) as registered_sp_ratio
from provider p 
left join (
select address_id, 1 as num 
from user u 
where u.type ='sp_owner' 
and LENGTH(address_id)>0
) t on p.owner = t.address_id 
`
		err := m.db.Raw(sql).Scan(&res).Error
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.AddTask(func() {
		err := m.db.Model(&models.Message{}).Count(&res.InternalContacts).Error
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.AddTask(func() {
		sql := `
select avg(num) as average_accesses_daily,
sum(num) as total_accesses
from (
select cast(created_at as date) as day,COUNT(id) as num from access_log al group by day
) t1
`
		err := m.db.Raw(sql).Scan(&res).Error
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.Wait()
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return &res, nil
}

func (m *Manager) IncrMetricsOverviewLast24h() (*MetricsOverview, error) {
	var res MetricsOverview
	var manager utils.AsyncManager
	errors := make([]error, 0)
	manager.AddTask(func() {
		sql := `
select t1.registered_providers-t2.registered_providers as registered_providers,
t1.auto_filled_providers-t2.auto_filled_providers as auto_filled_providers,
t1.registered_sp_ratio-t2.registered_sp_ratio as registered_sp_ratio
from (
select COALESCE(sum(t.num),0) as registered_providers, 
count(*)-COALESCE(sum(t.num),0) as auto_filled_providers,
COALESCE(sum(t.num),0)/(count(*)-COALESCE(sum(t.num),0)) as registered_sp_ratio
from provider p 
left join (
select address_id, 1 as num 
from user u 
where u.type ='sp_owner' 
and LENGTH(address_id)>0
) t on p.owner = t.address_id 
) t1, (
select COALESCE(sum(t.num),0) as registered_providers, 
count(*)-COALESCE(sum(t.num),0) as auto_filled_providers,
COALESCE(sum(t.num),0)/(count(*)-COALESCE(sum(t.num),0)) as registered_sp_ratio
from provider p 
left join (
select address_id, 1 as num 
from user u 
where u.type ='sp_owner' 
and LENGTH(address_id)>0
and u.created_at < (NOW()-INTERVAL 1 day)
) t on p.owner = t.address_id 
where p.created_at < (NOW()-INTERVAL 1 day)
) t2
`
		err := m.db.Raw(sql).Scan(&res).Error
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.AddTask(func() {
		err := m.db.Model(&models.Message{}).Where("created_at>(NOW()-INTERVAL 1 day)").
			Count(&res.InternalContacts).Error
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.AddTask(func() {
		sql := `
select avg(t1.num)-avg(t2.num) as average_accesses_daily, 
t3.num as total_accesses
from (
select cast(created_at as date) as day,COUNT(id) as num from access_log al group by day
) t1, (
select cast(created_at as date) as day,COUNT(id) as num from access_log al 
where al.created_at < (NOW()-INTERVAL 1 day) group by day
) t2, (
select COUNT(id) as num from access_log al where al.created_at > (NOW()-INTERVAL 1 day)
) t3
`
		err := m.db.Raw(sql).Scan(&res).Error
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.Wait()
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return &res, nil
}

type MetricsSpOverview struct {
	AverageRetrievalSuccessRatio decimal.Decimal `json:"average_retrieval_success_ratio"`
	SpsToDealNewClient           int64           `json:"sps_to_deal_new_client"`

	IncrAverageRetrievalSuccessRatio decimal.Decimal `json:"incr_average_retrieval_success_ratio"`
	IncrSpsToDealNewClientMonth      int64           `json:"incr_sps_to_deal_new_client_month"`
}

func (m *Manager) MetricsSpOverview() (*MetricsSpOverview, error) {
	var res MetricsSpOverview
	var manager utils.AsyncManager
	errors := make([]error, 0)
	manager.AddTask(func() {
		sql := `
select t1.avg_rsr as average_retrieval_success_ratio,
t1.avg_rsr-t2.avg_rsr as incr_average_retrieval_success_ratio 
from (
select avg(p.retrieval_success_rate) as avg_rsr from provider p where p.raw_power >0
) t1, (
select avg(p.retrieval_success_rate) as avg_rsr from provider p where p.raw_power >0 and p.created_at < (NOW()-INTERVAL 1 day)
) t2
`
		err := m.db.Raw(sql).Scan(&res).Error
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.AddTask(func() {
		startMonth := utils.MonthBegin(m.publishDate)
		endMonth := utils.MonthBegin(time.Now().AddDate(0, 1, 0))
		var err error
		res.SpsToDealNewClient, err = m.SpsToDealNewClientByRange(startMonth, endMonth)
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.AddTask(func() {
		var err error
		res.IncrSpsToDealNewClientMonth, err = m.SpsToDealNewClientLastMonth()
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.Wait()
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return &res, nil
}

type MetricsClientOverview struct {
	Clients         int64 `json:"clients"`
	TargetClients   int64 `json:"target_clients"`
	ClientsToDealSp int64 `json:"clients_to_deal_sp"`

	IncrClients              int64 `json:"incr_clients"`
	IncrTargetClients        int64 `json:"incr_target_clients"`
	IncrClientsToDealSpMonth int64 `json:"incr_clients_to_deal_sp_month"`
}

func (m *Manager) MetricsClientOverview() (*MetricsClientOverview, error) {
	var res MetricsClientOverview
	var manager utils.AsyncManager
	errors := make([]error, 0)
	manager.AddTask(func() {
		sql := `
select t1.num as clients,t2.num as target_clients 
from (
select count(*) as num from user u 
where u.type ='data_client' 
) t1, (
select count(*) as num from user u 
where u.type ='data_client' and u.address_id in (
select ci.address_id from client_info ci 
where ci.data_cap+ci.used_data_cap>=100*pow(1024,4))
)t2
`
		err := m.db.Raw(sql).Scan(&res).Error
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.AddTask(func() {
		var err error
		res.ClientsToDealSp, err = m.ClientsToDealSpByRange(utils.GetEpochByTime(m.publishDate), utils.GetEpochByTime(time.Now()))
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.AddTask(func() {
		sql := `
select t1.num as clients,t2.num as target_clients 
from (
select count(*) as num from user u 
where u.type ='data_client' and u.created_at > (NOW()-INTERVAL 1 day)
) t1, (
select count(*) as num from user u 
where u.type ='data_client' and u.address_id in (
select ci.address_id from client_info ci 
where ci.data_cap+ci.used_data_cap>=100*pow(1024,4))
and u.created_at > (NOW()-INTERVAL 1 day)
)t2
`
		err := m.db.Raw(sql).Scan(&res).Error
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.AddTask(func() {
		var err error
		res.IncrClientsToDealSpMonth, err = m.ClientsToDealSpLastMonth()
		if err != nil {
			errors = append(errors, err)
		}
	})

	manager.Wait()
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return &res, nil
}

type MetricsSpDetailItem struct {
	Owner        string          `json:"owner"`
	Providers    string          `json:"providers"`
	NewClients   string          `json:"new_clients"`
	NewClientNum int64           `json:"new_client_num"`
	RegisterTime *types.UnixTime `json:"register_time"`
}

func (m *Manager) MetricsSpToDealNewClientDetail(pagination types.PaginationParams) (int, []*MetricsSpDetailItem, error) {
	list := make([]*MetricsSpDetailItem, 0)
	var total int
	offset, size := utils.PaginationHelper(pagination.Page, pagination.PageSize, settings.DefaultPageSize)

	startMonth := utils.MonthBegin(m.publishDate)
	endMonth := utils.MonthBegin(time.Now().AddDate(0, 1, 0))
	stmtParams := map[string]interface{}{
		"startMonth": startMonth,
		"endMonth":   endMonth,
		"limit":      size,
		"offset":     offset,
	}

	subStmt := `
-- get relation with sp to deal new client
	select provider,client from (
		select t1.provider,t1.client,t2.repeats from (
			select psm.provider,psm.client
			from provider_stats_monthly psm
			where psm.month >= @startMonth and psm.month < @endMonth
			group by psm.provider,psm.client
		) t1 left join (
			select psm.provider,psm.client,1 as repeats
			from provider_stats_monthly psm
			where psm.month < @startMonth
			group by psm.provider,psm.client
		) t2 on t1.provider=t2.provider and t1.client=t2.client
	) t3 where repeats is null
`
	sqlCount := `
select count(*) from (
	select p.owner
	from (
	%s
	) tt1 left join provider p on tt1.provider=p.address 
	where p.owner is not null
	group by 1
) tt2
`
	sqlCount = fmt.Sprintf(sqlCount, subStmt)
	err := m.db.Raw(sqlCount, stmtParams).Scan(&total).Error
	if err != nil {
		return 0, nil, err
	}

	sql := `
select p.owner,
GROUP_CONCAT(distinct tt1.provider) as providers,
GROUP_CONCAT(distinct tt1.client) as clients,
count(distinct tt1.client) as client_num,
u.created_at as register_time
from (
%s
) tt1 left join provider p on tt1.provider=p.address 
left join user u on p.owner=u.address_id and u.type='sp_owner'
where p.owner is not null
group by 1
order by u.created_at, p.owner
limit @limit offset @offset
`
	sql = fmt.Sprintf(sql, subStmt)
	err = m.db.Raw(sql, stmtParams).Scan(&list).Error
	return total, list, err
}

type MetricsClientDetailItem struct {
	Client       string          `json:"client"`
	Owners       string          `json:"owners"`
	Providers    string          `json:"providers"`
	ProviderNum  int64           `json:"provider_num"`
	RegisterTime *types.UnixTime `json:"register_time"`
}

func (m *Manager) MetricsClientToDealSpDetail(pagination types.PaginationParams) (int, []*MetricsClientDetailItem, error) {
	list := make([]*MetricsClientDetailItem, 0)
	var total int
	offset, size := utils.PaginationHelper(pagination.Page, pagination.PageSize, settings.DefaultPageSize)

	startMonth := utils.MonthBegin(m.publishDate)
	endMonth := utils.MonthBegin(time.Now().AddDate(0, 1, 0))
	stmtParams := map[string]interface{}{
		"startMonth": startMonth,
		"endMonth":   endMonth,
		"limit":      size,
		"offset":     offset,
	}

	subStmt := `
-- get relation with client to deal sp
	select psm.client,psm.provider
	from provider_stats_monthly psm
	where psm.month >= @startMonth and psm.month < @endMonth
	and client in (select ci.address_id from client_info ci where ci.data_cap+ci.used_data_cap>=100*pow(1024,4))
	group by 1,2
`
	sqlCount := `
select count(*) from (
	select tt1.client
	from (
	%s
	) tt1
	group by 1
) tt2
`
	sqlCount = fmt.Sprintf(sqlCount, subStmt)
	err := m.db.Raw(sqlCount, stmtParams).Scan(&total).Error
	if err != nil {
		return 0, nil, err
	}

	sql := `
select tt1.client,
GROUP_CONCAT(distinct tt1.provider) as providers,
GROUP_CONCAT(distinct p.owner) as owners,
count(distinct tt1.provider) as provider_num,
u.created_at as register_time
from (
%s
) tt1 left join provider p on tt1.provider=p.address 
left join user u on tt1.client=u.address_id and u.type='data_client'
group by 1
order by u.created_at, p.owner
limit @limit offset @offset
`
	sql = fmt.Sprintf(sql, subStmt)
	err = m.db.Raw(sql, stmtParams).Scan(&list).Error
	return total, list, err
}
