package aggregator

import (
	"backend/internal/models"
	"backend/internal/pools"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"log/slog"
	"time"
)

func handleEvents(ctx context.Context, c *app.RequestContext) {

	req := pools.EventFromGamePool.Get().(*models.UpdateFromGame)
	defer pools.EventFromGamePool.Put(req)

	if err := c.BindAndValidate(req); err != nil {
		slog.Error("Unable to unpack request", "err", err)
		c.JSON(consts.StatusBadRequest, utils.H{"status": "error"})
		return
	}

	// Defaults event time to current time
	if req.Time.IsZero() {
		req.Time = time.Now()
	}

	// TODO API key implementation (project tagging and API key creation)
	//if err := req.ApiKey.Validate(); err != nil {
	//	slog.Error("Improper API key")
	//}

	slog.Debug("Received event update from game server", "event", req)

	// Acquire agg instance
	agg := GetAggregator()

	// Send events to kafka
	err := agg.EmitBatchKillEvent(ctx, req.ExtractKillRecords())
	if err != nil {
		slog.Error("Unable to send kill events to kafka")
		c.JSON(consts.StatusInternalServerError, utils.H{"status": "error"})
		return
	}

	c.JSON(consts.StatusOK, utils.H{"status": "ok"})
	return
}
