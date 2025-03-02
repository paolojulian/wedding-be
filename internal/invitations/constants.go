package invitations

import "errors"

var (
	ErrNameIsRequired        = errors.New("name is required")
	ErrVoucherCodeIsRequired = errors.New("voucher code is required")

	ErrInvalidIDFormat           = errors.New("invalid id format")
	ErrNoFieldsToUpdate          = errors.New("no fields to update")
	ErrInvitationNotFound        = errors.New("invitation is not found")
	ErrInvitationMessageNotFound = errors.New("invitation message is not found")

	ErrCannotUpdateInDB = errors.New("cannot update database")

	ErrIsAlreadyLocked = errors.New("the app is already locked")
)
