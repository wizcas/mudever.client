package packet

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wizcas/mudever.svc/telnet/protocol"
)

func _assertCmd(result, expect []byte, err error, size int) {
	So(err, ShouldBeNil)
	So(result, ShouldResemble, expect)
	So(len(result), ShouldEqual, size)
}

func TestCmdStringify(t *testing.T) {
	Convey("Given a mono command packet", t, func() {
		p := NewControlCommandPacket(protocol.GA)
		Convey("Print its content", func() {
			So(p.String(), ShouldEqual, "[CMD] 249")
		})
	})
	Convey("Given an option command packet", t, func() {
		p := NewOptionCommandPacket(protocol.WILL, protocol.Echo)
		Convey("Print its content", func() {
			So(p.String(), ShouldEqual, "[CMD] 251 | 1")
		})
	})
}

func TestControlCommandSerialize(t *testing.T) {
	Convey("Given a mono command packet", t, func() {
		p := NewControlCommandPacket(protocol.GA)
		Convey("Serialize into 2 bytes", func() {
			result, err := p.Serialize()
			_assertCmd(result, []byte{protocol.IAC, protocol.GA}, err, 2)
		})
	})
}

func TestOptionCommandSerialize(t *testing.T) {
	Convey("Given an option command packet", t, func() {
		p := NewOptionCommandPacket(protocol.WILL, protocol.Echo)
		Convey("Serialize into 3 bytes", func() {
			result, err := p.Serialize()
			_assertCmd(result, []byte{protocol.IAC, protocol.WILL, protocol.Echo}, err, 3)
		})
	})
}
