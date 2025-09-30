package errors

import (
	"fmt"
	"strings"
)

type AppModelNotFoundError struct {
	ID       string
	Entity   string
	Comments string
}

var AppModelNotFound *AppModelNotFoundError

func NewAppModelNotFoundError(id string, entity string, comments string) *AppModelNotFoundError {
	return &AppModelNotFoundError{
		ID:       id,
		Entity:   entity,
		Comments: comments,
	}
}

func (n *AppModelNotFoundError) Error() string {
	if strings.TrimSpace(n.Comments) == "" {
		return fmt.Sprintf("model [%s] with id [%s] not found", n.Entity, n.ID)
	}
	return fmt.Sprintf("model [%s] with id [%s] not found with comments [%s]", n.Entity, n.ID, n.Comments)
}
