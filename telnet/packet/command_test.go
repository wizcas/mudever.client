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
	Convey("Given a control command", t, func() {
		p := NewControlCommandPacket(protocol.GA)
		Convey("Print its content", func() {
			So(p.String(), ShouldEqual, "[CMD] GA")
		})
	})
	Convey("Given an option command packet", t, func() {
		p := NewOptionCommandPacket(protocol.WILL, protocol.Echo)
		Convey("Print its content", func() {
			So(p.String(), ShouldEqual, "[CMD] WILL > ECHO")
		})
	})
	Convey("Given an unknown control command", t, func() {
		p := NewControlCommandPacket(22)
		Convey("Print its command as decimal value", func() {
			So(p.String(), ShouldEqual, "[CMD] 22")
		})
	})
	Convey("Given an unknown option command", t, func() {
		p := NewOptionCommandPacket(22, 254)
		Convey("Print its command and option as decimal value", func() {
			So(p.String(), ShouldEqual, "[CMD] 22 > 254")
		})
	})
}

func TestControlCommandSerialize(t *testing.T) {
	Convey("Given a control command", t, func() {
		p := NewControlCommandPacket(protocol.GA)
		Convey("Serialize into 2 bytes", func() {
			result, err := p.Serialize()
			_assertCmd(result, []byte{byte(protocol.IAC), byte(protocol.GA)}, err, 2)
		})
	})
}

func TestOptionCommandSerialize(t *testing.T) {
	Convey("Given an option command", t, func() {
		p := NewOptionCommandPacket(protocol.WILL, protocol.Echo)
		Convey("Serialize into 3 bytes", func() {
			result, err := p.Serialize()
			_assertCmd(result, []byte{byte(protocol.IAC), byte(protocol.WILL), byte(protocol.Echo)}, err, 3)
		})
	})
}
