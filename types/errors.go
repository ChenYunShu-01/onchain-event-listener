package types

import "errors"

var (
	ErrAddressMissing  = errors.New("contract address is missing")
	ErrStarkKeyMissing = errors.New("stark key is missing")
	
	ErrAmountInvalid   = errors.New("amount is invalid")
	ErrTokenIdMissing  = errors.New("token id is missing")
)
