package utils

import (
	"github.com/google/uuid"
)

func GenerateUuid() int64 {
	id := uuid.New()
	return int64(id.ID())
}
