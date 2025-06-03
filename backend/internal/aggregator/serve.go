package aggregator

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Serve(port int) {

	h := server.Default(server.WithHostPorts(fmt.Sprintf("0.0.0.0:%d", port)))

	h.POST("/events", handleEvents)

	h.Spin()
}
