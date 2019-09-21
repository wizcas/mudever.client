package packet

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wizcas/mudever.svc/telnet/telbyte"
)

func _assertSub(result, expect []byte, err error) {
	So(result, ShouldResemble, expect)
	So(err, ShouldBeNil)
}

func TestSubStringify(t *testing.T) {
	Convey("Given a complex subnegotiation packet", t, func() {
		p := NewSubPacket(telbyte.TTYPE, []byte{0}, []byte("MUD"))
		Convey("Print its content", func() {
			So(p.String(), ShouldEqual, "[SUB] TERMINAL-TYPE > [0 77 85 68]")
		})
	})
	Convey("Given a no-parameter subnegotiation packet", t, func() {
		p := NewSubPacket(telbyte.TTYPE)
		Convey("Print its content", func() {
			So(p.String(), ShouldEqual, "[SUB] TERMINAL-TYPE > []")
		})
	})
	Convey("Given a no-parameter subnegotiation packet of unknown option", t, func() {
		p := NewSubPacket(254)
		Convey("Print its content with decimal option value", func() {
			So(p.String(), ShouldEqual, "[SUB] 254 > []")
		})
	})
}

func TestSimpleSub(t *testing.T) {
	Convey("Given a simple subnegotiation (1 parameter)", t, func() {
		p := NewSubPacket(telbyte.TTYPE, []byte{1})
		Convey("Serialize correctly", func() {
			expect := []byte{
				byte(telbyte.IAC), byte(telbyte.SB),
				byte(telbyte.TTYPE),
				1,
				byte(telbyte.IAC), byte(telbyte.SE),
			}
			result, err := p.Serialize()
			_assertSub(result, expect, err)
		})
	})
}

func TestComplexSub(t *testing.T) {
	Convey("Given a subnegotiation with 2 parameters", t, func() {
		p := NewSubPacket(telbyte.TTYPE, []byte{0}, []byte("MUDEVER"))
		Convey("Serialize correctly", func() {
			expect := []byte{byte(telbyte.IAC), byte(telbyte.SB), byte(telbyte.TTYPE), 0,
				'M', 'U', 'D', 'E', 'V', 'E', 'R',
				byte(telbyte.IAC), byte(telbyte.SE)}
			result, err := p.Serialize()
			_assertSub(result, expect, err)
		})
	})
	Convey("Given a subnegotiation with 0xFF in parameter", t, func() {
		p := NewSubPacket(telbyte.TTYPE, []byte{0, 255}, []byte("MUDEVER"))
		Convey("Serialize correctly with 0xFF escaped", func() {
			expect := []byte{byte(telbyte.IAC), byte(telbyte.SB), byte(telbyte.TTYPE), 0,
				255, 255,
				'M', 'U', 'D', 'E', 'V', 'E', 'R',
				byte(telbyte.IAC), byte(telbyte.SE)}
			result, err := p.Serialize()
			_assertSub(result, expect, err)
		})
	})
}
