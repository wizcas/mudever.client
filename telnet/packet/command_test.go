package packet

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wizcas/mudever.svc/telnet/telbyte"
)

func _assertCmd(result, expect []byte, err error, size int) {
	So(err, ShouldBeNil)
	So(result, ShouldResemble, expect)
	So(len(result), ShouldEqual, size)
}

func TestCmdStringify(t *testing.T) {
	Convey("Given a control command", t, func() {
		p := NewControlCommandPacket(telbyte.GA)
		Convey("Print its content", func() {
			So(p.String(), ShouldEqual, "[CMD] GA")
		})
	})
	Convey("Given an option command packet", t, func() {
		p := NewOptionCommandPacket(telbyte.WILL, telbyte.ECHO)
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
		p := NewControlCommandPacket(telbyte.GA)
		Convey("Serialize into 2 bytes", func() {
			result, err := p.Serialize()
			_assertCmd(result, []byte{byte(telbyte.IAC), byte(telbyte.GA)}, err, 2)
		})
	})
}

func TestOptionCommandSerialize(t *testing.T) {
	Convey("Given an option command", t, func() {
		p := NewOptionCommandPacket(telbyte.WILL, telbyte.ECHO)
		Convey("Serialize into 3 bytes", func() {
			result, err := p.Serialize()
			_assertCmd(result, []byte{byte(telbyte.IAC), byte(telbyte.WILL), byte(telbyte.ECHO)}, err, 3)
		})
	})
}
