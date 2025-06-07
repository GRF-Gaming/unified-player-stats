package pools

import (
	"backend/internal/models"
	"sync"
)

var EventFromGamePool = sync.Pool{
	New: func() interface{} {
		return new(models.UpdateFromGame)
	},
}
