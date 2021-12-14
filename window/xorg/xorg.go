package xorg

import (
	"log"
	"ykeysnail/window"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type Xorg struct{}

var (
	activeWindow string = "_NET_ACTIVE_WINDOW"
	wmName       string = "_NET_WM_NAME"
	wmClass      string = "WM_CLASS"
)

var (
	conn                *xgb.Conn
	active, name, class *xproto.InternAtomReply
)

func init() {
	var err error
	conn, err = xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	active, err = xproto.InternAtom(conn, true, uint16(len(activeWindow)),
		activeWindow).Reply()
	if err != nil {
		log.Fatal(err)
	}

	name, err = xproto.InternAtom(conn, true, uint16(len(wmName)),
		wmName).Reply()
	if err != nil {
		log.Fatal(err)
	}

	class, err = xproto.InternAtom(conn, true, uint16(len(wmClass)),
		wmClass).Reply()
	if err != nil {
		log.Fatal(err)
	}

}

func (w *Xorg) Window() *window.WindowInfo {

	reply, err := xproto.GetProperty(conn, false, xproto.Setup(conn).DefaultScreen(conn).Root, active.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		log.Print(err)
	}

	windowId := xproto.Window(xgb.Get32(reply.Value))

	replyName, err := xproto.GetProperty(conn, false, windowId, name.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		log.Print(err)
		return &window.WindowInfo{}
	}

	replyClass, err := xproto.GetProperty(conn, false, windowId, class.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		log.Print(err)
		return &window.WindowInfo{Title: string(replyName.Value)}
	}

	return &window.WindowInfo{Title: string(replyName.Value), Class: string(replyClass.Value)}
}
