package dto

import "fmt"

type ErrorNotFound struct {
	EntityName string
	EntityID   int
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("cannot find %s with id %d", e.EntityName, e.EntityID)
}
