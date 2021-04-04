package n26api_test

import (
	"github.com/google/uuid"

	"github.com/nhatthm/n26api/pkg/testkit"
)

var (
	n26Username = "john.doe"
	n26Password = "123456"
)

func mockServer(deviceID uuid.UUID, mocks ...testkit.ServerOption) testkit.ServerMocker {
	return testkit.MockServer(n26Username, n26Password, deviceID, mocks...)
}
