package models

type InvitationStatus string

// const (
// 	StatusGoing    InvitationStatus = "going"
// 	StatusNotGoing InvitationStatus = "not-going"
// 	StatusPending  InvitationStatus = "pending"
// )

type Invitation struct {
	ID            string `json:"id"`
	VoucherCode   string `json:"voucher_code"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	GuestsAllowed int    `json:"guests_allowed"`
	GuestsToBring int    `json:"guests_to_bring"`
}
