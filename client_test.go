package n26api_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api"
	"github.com/nhatthm/n26api/pkg/testkit"
)

var (
	n26Username = "john.doe"
	n26Password = "123456"
)

func mockServer(deviceID uuid.UUID, mocks ...testkit.ServerOption) testkit.ServerMocker {
	return testkit.MockServer(n26Username, n26Password, deviceID, mocks...)
}

func TestClient_DeviceID(t *testing.T) {
	t.Parallel()

	t.Run("no device id", func(t *testing.T) {
		t.Parallel()

		c := n26api.NewClient()
		assert.NotEqual(t, uuid.UUID{}, c.DeviceID())
	})

	t.Run("with device id", func(t *testing.T) {
		deviceID := uuid.New()

		c := n26api.NewClient(n26api.WithDeviceID(deviceID))
		assert.Equal(t, deviceID, c.DeviceID())
	})
}
