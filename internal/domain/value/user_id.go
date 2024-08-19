package value

import "fmt"

type UserId uint64

const UnknownUserId UserId = 0

func (id UserId) IsUnknown() bool {
	return id == UnknownUserId
}

func ValidateUserId(userId UserId) error {
	if userId < 0 {
		return fmt.Errorf("userId must be greater than or equal to zero, given %d", userId)
	}
	return nil
}
