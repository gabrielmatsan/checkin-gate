package lib

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nrednav/cuid2"
)

type IDType int

const (
	UUID IDType = iota
	CUID2
)

func GenerateID(idType IDType) (string, error) {
	switch idType {
	case UUID:
		return uuid.NewString(), nil
	case CUID2:
		return cuid2.Generate(), nil
	default:
		return "", fmt.Errorf("unsupported id type: %d", idType)
	}
}
