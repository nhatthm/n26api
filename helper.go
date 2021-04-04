package n26api

import (
	"os"

	"github.com/google/uuid"
)

func deviceID(id uuid.UUID) uuid.UUID {
	if id != emptyUUID {
		return id
	}

	if envUUID := os.Getenv(envDeviceID); envUUID != "" {
		return uuid.MustParse(envUUID)
	}

	return uuid.New()
}
