package pools

import (
	"backend/internal/models"
	"sync"
)

var EventKillRecordPool = sync.Pool{
	New: func() interface{} {
		return new(models.EventKillRecord)
	},
}
