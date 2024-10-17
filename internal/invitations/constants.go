package invitations

import "errors"

var (
	ErrNameIsRequired        = errors.New("name is required")
	ErrVoucherCodeIsRequired = errors.New("voucher code is required")

	ErrInvalidIDFormat = errors.New("invalid id format")
)
