package mapper

import (
	"testing"

	"github.com/ABA-Developer/go-backend-template/internal/presentation/role/presenter"

	"github.com/stretchr/testify/assert"
)

func TestToUpdateRoleAccessPayload(t *testing.T) {
	allowed := true
	denied := false
	request := presenter.UpdateRoleAccessRequest{
		AccessItem: []presenter.UpdateRoleAccessItem{
			{AccessID: 1, HasAccess: &allowed},
			{AccessID: 2, HasAccess: &denied},
		},
	}

	payload := ToUpdateRoleAccessPayload(request)

	assert.Len(t, payload.AccessItem, 2)
	assert.Equal(t, 1, payload.AccessItem[0].AccessID)
	assert.Equal(t, &allowed, payload.AccessItem[0].HasAccess)
	assert.Equal(t, 2, payload.AccessItem[1].AccessID)
	assert.Equal(t, &denied, payload.AccessItem[1].HasAccess)
}
