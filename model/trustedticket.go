package model

type TrustedTicketRequest struct {
	Username   string `form:"username"`
	Targetsite string `form:"target_site"`
}

type TrustedTicket struct {
	Value string
}
