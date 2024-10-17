package models

type InvitationStatus string

type Database struct {
	Users       []User       `json:"users"`
	Invitations []Invitation `json:"invitations"`
}

type Invitation struct {
	ID string `json:"id" bson:"_id,omitempty"` // ObjectID stored as a string
	// This is the index to determine the position of the invitation in the list
	Index         int    `json:"index" bson:"index"`
	VoucherCode   string `json:"voucher_code" bson:"voucher_code"`
	Name          string `json:"name" bson:"name"`
	Status        string `json:"status" bson:"status"`
	GuestsAllowed int    `json:"guests_allowed" bson:"guests_allowed"`
	GuestsToBring int    `json:"guests_to_bring" bson:"guests_to_bring"`
}

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"` // ObjectID stored as a string
	Username string `json:"username"`
	Password string `json:"-" bson:"password"` // Exclude password from JSON responses
}
