package models

type InvitationStatus string

type Database struct {
	Users       []User       `json:"users"`
	Invitations []Invitation `json:"invitations"`
}

type Invitation struct {
	ID string `json:"id"`
	// This is the index to determine the position of the invitation in the list
	Index         int    `json:"index"`
	VoucherCode   string `json:"voucher_code"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	GuestsAllowed int    `json:"guests_allowed"`
	GuestsToBring int    `json:"guests_to_bring"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
