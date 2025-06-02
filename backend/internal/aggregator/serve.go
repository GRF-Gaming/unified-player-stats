package aggregator

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Serve(port int) {

	h := server.Default(server.WithHostPorts("localhost:8080"))

	h.POST("/events", handleEvents)

	h.Spin()
}
