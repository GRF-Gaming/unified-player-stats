package aggregator

import (
	"backend/internal/models"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"log/slog"
)

func handleEvents(ctx context.Context, c *app.RequestContext) {

	var req models.Event
	if err := c.BindAndValidate(&req); err != nil {
		slog.Error("Unable to unpack request", "err", err)
		c.JSON(consts.StatusBadRequest, utils.H{})
		return
	}

	if err := req.ApiKey.Validate(); err != nil {
		slog.Error("Improper API key")
	}

	slog.Info("Received", "data", req)

	c.JSON(consts.StatusOK, utils.H{"message": "pong"})
	return
}
