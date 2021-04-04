package n26api

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeviceID_NotEmpty(t *testing.T) {
	fixedUUID := uuid.New()

	// UUID is not empty
	assert.Equal(t, fixedUUID, deviceID(fixedUUID))
}

func TestDeviceID_FromEnv(t *testing.T) {
	currentDeviceID := os.Getenv(envDeviceID)

	t.Cleanup(func() {
		_ = os.Setenv(envDeviceID, currentDeviceID)
	})

	emptyUUID := uuid.UUID{}

	t.Run("valid device id", func(t *testing.T) {
		newUUID := uuid.New()
		_ = os.Setenv(envDeviceID, newUUID.String())

		assert.Equal(t, newUUID, deviceID(emptyUUID))
	})

	t.Run("invalid device id", func(t *testing.T) {
		_ = os.Setenv(envDeviceID, "hello world")

		assert.Panics(t, func() {
			deviceID(emptyUUID)
		})
	})
}

func TestDeviceID_New(t *testing.T) {
	emptyUUID := uuid.UUID{}

	assert.NotEqual(t, emptyUUID, deviceID(emptyUUID))
}
