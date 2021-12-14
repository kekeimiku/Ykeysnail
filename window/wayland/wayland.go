package wayland

import (
	"github.com/Difrex/gosway/ipc"
)

type Wayland struct{}

var swayConn *ipc.SwayConnection

func init() {
	var err error
	swayConn, err = ipc.NewSwayConnection()
	if err != nil {
		panic(err)
	}
	_, err = swayConn.SendCommand(ipc.IPC_SUBSCRIBE, `["window"]`)
	if err != nil {
		panic(err)
	}
}
