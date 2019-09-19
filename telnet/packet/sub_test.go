package packet

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wizcas/mudever.svc/telnet/protocol"
)

func _assertSub(result, expect []byte, err error) {
	So(result, ShouldResemble, expect)
	So(err, ShouldBeNil)
}

func TestSubStringify(t *testing.T) {
	Convey("Given a complex subnegotiation packet", t, func() {
		p := NewSubPacket(protocol.TerminalType, []byte{0}, []byte("MUD"))
		Convey("Print its content", func() {
			So(p.String(), ShouldEqual, "[SUB] 24 | [0 77 85 68]")
		})
	})
	Convey("Given a no-parameter subnegotiation packet", t, func() {
		p := NewSubPacket(protocol.TerminalType)
		Convey("Print its content", func() {
			So(p.String(), ShouldEqual, "[SUB] 24 | []")
		})
	})
}

func TestSimpleSub(t *testing.T) {
	Convey("Given a simple subnegotiation (1 parameter)", t, func() {
		p := NewSubPacket(protocol.TerminalType, []byte{1})
		Convey("Serialize correctly", func() {
			expect := []byte{protocol.IAC, protocol.SB, protocol.TerminalType, 1, protocol.IAC, protocol.SE}
			result, err := p.Serialize()
			_assertSub(result, expect, err)
		})
	})
}

func TestComplexSub(t *testing.T) {
	Convey("Given a subnegotiation with 2 parameters", t, func() {
		p := NewSubPacket(protocol.TerminalType, []byte{0}, []byte("MUDEVER"))
		Convey("Serialize correctly", func() {
			expect := []byte{protocol.IAC, protocol.SB, protocol.TerminalType, 0,
				'M', 'U', 'D', 'E', 'V', 'E', 'R',
				protocol.IAC, protocol.SE}
			result, err := p.Serialize()
			_assertSub(result, expect, err)
		})
	})
}
