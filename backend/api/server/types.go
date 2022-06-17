package server

import (
	"github.com/filedrive-team/filfind/backend/repo"
	"github.com/filedrive-team/filfind/backend/types"
)

type SignUpParam struct {
	Type     string `json:"type" binding:"required,oneof=sp_owner data_client" example:"sp_owner"`
	Name     string `json:"name" binding:"required,max=128" example:"hello"`
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,password" example:"Hello123"`
	// example privateKey hex: 7b2254797065223a22736563703235366b31222c22507269766174654b6579223a22347179396744327177724f68764148305267667738624b346b52546e575535595337622f4338374435686f3d227d
	Address string `json:"address" binding:"required,address,min=3" example:"f1gxcq2s72oepgufqrkbblbwgnxosrkwn3jib3bmy"`
	// raw data:
	//Signature for filfind
	//f1gxcq2s72oepgufqrkbblbwgnxosrkwn3jib3bmy
	//2022-04-14T12:03:45.169Z
	// raw data hex: 5369676e617475726520666f722066696c66696e640d0a663167786371327337326f657067756671726b62626c6277676e786f73726b776e336a696233626d790d0a323032322d30342d31345431323a30333a34352e3136395a
	Message   string `json:"message" binding:"required,hexadecimal" example:"5369676e617475726520666f722066696c66696e640d0a663167786371327337326f657067756671726b62626c6277676e786f73726b776e336a696233626d790d0a323032322d30342d31345431323a30333a34352e3136395a"`
	Signature string `json:"signature" binding:"required,hexadecimal" example:"01032aa043dea7ef185a3f5345d8bc9ee672b91584b090fca782e04f1bf36792211addd36fe6cb1e8829045d95175a48a3e431061362044bf830275829d9a7fe0d00"`
}

type UserParam struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,password" example:"Hello123"`
}

type ResetPwdParams struct {
	Email       string `json:"email" binding:"required,email" example:"test@example.com"`
	VCode       string `json:"vcode" binding:"required,number,len=6" example:"123456"`
	NewPassword string `json:"new_password" binding:"required,password" example:"Hello456"`
}

type EmailVcodeParam struct {
	Email string `json:"email" binding:"required,email" example:"test@example.com"`
}

type ModifyPasswordParam struct {
	Password    string `json:"password" binding:"required" example:"Hello123"`
	NewPassword string `json:"new_password" binding:"required,password" example:"Hello456"`
}

type ProviderListParam struct {
	types.PaginationParams
	repo.OrderParam
	repo.FilterParam
	Search string `json:"search" form:"search" binding:"omitempty,lte=255"` // miner id/name/location
}

type AddressIdParam struct {
	AddressId string `json:"address_id" form:"address_id" binding:"required,addressid" example:"f01624861"`
}

type ProfileParam struct {
	Name         string `json:"name" binding:"omitempty,max=128" example:"example name"`
	Avatar       string `json:"avatar" binding:"omitempty,max=1024" example:""`
	Logo         string `json:"logo" binding:"omitempty,max=1024" example:""`
	Location     string `json:"location" binding:"omitempty,max=128" example:"Shanghai,China"`
	ContactEmail string `json:"contact_email" binding:"omitempty,email,max=256" example:"public@example.com"`
	Slack        string `json:"slack" binding:"omitempty,max=128" example:""`
	Github       string `json:"github" binding:"omitempty,max=128" example:""`
	Twitter      string `json:"twitter" binding:"omitempty,max=128" example:""`
	Description  string `json:"description" binding:"omitempty,max=2048" example:"More information about us."`
}

type ClientListParam struct {
	types.PaginationParams
	repo.ClientOrderParam
	Search string `json:"search" form:"search" binding:"omitempty,lte=255"` // client id/name/location
}

type ClientDetailParams struct {
	Bandwidth          string `json:"bandwidth" binding:"omitempty,max=128" example:"300M"`
	MonthlyStorage     string `json:"monthly_storage" binding:"omitempty,max=128" example:"10TiB"`
	UseCase            string `json:"use_case" binding:"omitempty,max=128" example:"Entertainment/Media/Science"`
	ServiceRequirement string `json:"service_requirement" binding:"omitempty,max=1024" example:"More information about us."`
}

type ProviderDetailParams struct {
	Address         string `json:"address" binding:"omitempty,addressid" example:"f01234"`
	AvailableDeals  string `json:"available_deals" binding:"omitempty,max=128" example:"10TiB/D"`
	Bandwidth       string `json:"bandwidth" binding:"omitempty,max=128" example:"300M"`
	SealingSpeed    string `json:"sealing_speed" binding:"omitempty,max=128" example:"10TiB/D"`
	ParallelDeals   string `json:"parallel_deals" binding:"omitempty,max=128" example:"10"`
	RenewableEnergy string `json:"renewable_energy" binding:"omitempty,max=128" example:"1MWh"`
	Certification   string `json:"certification" binding:"omitempty,max=128" example:"PCI Compliance"`
	IsMember        string `json:"is_member" binding:"omitempty,oneof=Yes No" example:"No"`
	Experience      string `json:"experience" binding:"omitempty,max=128" example:"Textile/Estuary"`
}

type ClientHistoryDealStatsParams struct {
	types.PaginationParams
	AddressIdParam
}

type ClientReviewsParam struct {
	types.PaginationParams
	AddressIdParam
}

type spOwnerReviewsParam struct {
	types.PaginationParams
	AddressIdParam
}

type ReviewParams struct {
	Provider string `json:"provider" binding:"required,addressid" example:"f01662887"`
	Score    int    `json:"score" binding:"required,min=0,max=5" example:"5"`
	Content  string `json:"content" binding:"required,max=1024" example:"Great. It's the ultimate experience."`
	Title    string `json:"title" binding:"required,max=128" example:"Ultimate Experience"`
}

type AdminUserParam struct {
	Name     string `json:"name" binding:"required" example:"filfind"`
	Password string `json:"password" binding:"required,password" example:"filFind123"`
}
