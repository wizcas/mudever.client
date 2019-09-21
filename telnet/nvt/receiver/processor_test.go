package receiver

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wizcas/mudever.svc/telnet/packet"
	"github.com/wizcas/mudever.svc/telnet/telbyte"
)

type procTester struct {
	*processor
	chOutput chan packet.Packet
	chErr    chan error
}

func _newProcTester() *procTester {
	return &procTester{
		processor: newProcessor(),
		chOutput:  make(chan packet.Packet),
		chErr:     make(chan error),
	}
}

func (pt *procTester) test(data []byte) {
	go pt.proc(data, pt.chOutput, pt.chErr)
}

func _assertData(p packet.Packet, expectData []byte) {
	So(p, ShouldHaveSameTypeAs, &packet.DataPacket{})
	So(p.(*packet.DataPacket).Data, ShouldResemble, expectData)
}

func _assertCmd(p packet.Packet, expectCmd telbyte.Command, expectOption telbyte.Option) {
	So(p, ShouldHaveSameTypeAs, &packet.CommandPacket{})
	cmdp := p.(*packet.CommandPacket)
	So(cmdp.Command, ShouldEqual, expectCmd)
	So(cmdp.Option, ShouldEqual, expectOption)
}
func _assertSub(p packet.Packet, expectOption telbyte.Option, expectParameter []byte) {
	So(p, ShouldHaveSameTypeAs, &packet.SubPacket{})
	subp := p.(*packet.SubPacket)
	So(subp.Option, ShouldEqual, expectOption)
	So(subp.Parameter, ShouldResemble, expectParameter)
}

func TestProcPlainData(t *testing.T) {
	Convey("Given plain data", t, func() {
		data := []byte("I'm plain data")
		Convey("Get: DataPacket", func() {
			tester := _newProcTester()
			tester.test(data)
			p := <-tester.chOutput
			_assertData(p, data)
		})
	})
	Convey("Given plain data with 0xFF", t, func() {
		data := append(append([]byte("I'm plain"), []byte{255, 255}...), []byte("data")...)
		expect := append(append([]byte("I'm plain"), []byte{255}...), []byte("data")...)
		Convey("Get: DataPacket containing unescaped 0xFF", func() {
			tester := _newProcTester()
			tester.test(data)
			p := <-tester.chOutput
			_assertData(p, expect)
		})
	})
}

func TestProcCommand(t *testing.T) {
	Convey("Given a control command", t, func() {
		data := []byte{byte(telbyte.IAC), byte(telbyte.GA)}
		Convey("Get: CommandPacket whose Option = IAC", func() {
			tester := _newProcTester()
			tester.test(data)
			p := <-tester.chOutput
			_assertCmd(p, telbyte.GA, telbyte.NoOption)
		})
	})
	Convey("Given an option command", t, func() {
		data := []byte{byte(telbyte.IAC), byte(telbyte.WILL), byte(telbyte.ECHO)}
		Convey("Get: CommandPacket whose Command & Option are both set", func() {
			tester := _newProcTester()
			tester.test(data)
			p := <-tester.chOutput
			_assertCmd(p, telbyte.WILL, telbyte.ECHO)
		})
	})
}

func TestProcSub(t *testing.T) {
	Convey("Given a subnegotiation", t, func() {
		data := []byte{byte(telbyte.IAC), byte(telbyte.SB), byte(telbyte.TTYPE),
			0,
			'M', 'U', 'D', 'E', 'V', 'E', 'R',
			byte(telbyte.IAC), byte(telbyte.SE),
		}
		expectOption := telbyte.TTYPE
		expectParameter := append([]byte{0}, []byte("MUDEVER")...)
		Convey("Get: SubPacket whose Option & Parameter are both set", func() {
			tester := _newProcTester()
			tester.test(data)
			p := <-tester.chOutput
			_assertSub(p, expectOption, expectParameter)
		})
	})
}

func TestProcInvalidCommand(t *testing.T) {
	Convey("Given data includes invalid telnet command", t, func() {
		goodData1 := []byte{1, 2, 3, 4}
		goodData2 := []byte{5, 6, 7, 8}
		badCmd := []byte{byte(telbyte.IAC), 0}
		data := append(goodData1, append(badCmd, goodData2...)...)
		Convey("Get: error then rest of the data", func() {
			tester := _newProcTester()
			tester.test(data)
			p1 := <-tester.chOutput
			_assertData(p1, goodData1)
			err := <-tester.chErr
			So(err, ShouldNotBeNil)
			p2 := <-tester.chOutput
			_assertData(p2, goodData2)
		})
	})
}

func TestProcFlow(t *testing.T) {
	plain := []byte("Hello World")

	ctrlcmd := telbyte.GA
	optcmd := telbyte.WILL
	optval := telbyte.ECHO
	subopt := telbyte.TTYPE
	subparam := []byte{0, 'M', 'U', 'D', 'E', 'V', 'E', 'R'}

	pieceCtrlCmd := []byte{byte(telbyte.IAC), byte(ctrlcmd)}
	pieceOptCmd := []byte{byte(telbyte.IAC), byte(optcmd), byte(optval)}
	pieceSub := append(
		append([]byte{byte(telbyte.IAC), byte(telbyte.SB), byte(subopt)}, subparam...),
		byte(telbyte.IAC), byte(telbyte.SE))

	Convey("Given [CTRL CMD | DATA]", t, func() {
		data := append(pieceCtrlCmd, plain...)
		Convey("Get: CommandPacket + DataPacket", func() {
			tester := _newProcTester()
			tester.test(data)
			p1 := <-tester.chOutput
			_assertCmd(p1, ctrlcmd, telbyte.NoOption)
			p2 := <-tester.chOutput
			_assertData(p2, plain)
		})
	})
	Convey("Given [OPT CMD | DATA]", t, func() {
		data := append(pieceOptCmd, plain...)
		Convey("Get: CommandPacket + DataPacket", func() {
			tester := _newProcTester()
			tester.test(data)
			p1 := <-tester.chOutput
			_assertCmd(p1, optcmd, optval)
			p2 := <-tester.chOutput
			_assertData(p2, plain)
		})
	})
	Convey("Given [DATA | CTRL CMD]", t, func() {
		data := append(plain, pieceCtrlCmd...)
		Convey("Get: DataPacket + CommandPacket", func() {
			tester := _newProcTester()
			tester.test(data)
			p1 := <-tester.chOutput
			_assertData(p1, plain)
			p2 := <-tester.chOutput
			_assertCmd(p2, ctrlcmd, telbyte.NoOption)
		})
	})
	Convey("Given [DATA | OPT CMD]", t, func() {
		data := append(plain, pieceOptCmd...)
		Convey("Get: DataPacket + CommandPacket", func() {
			tester := _newProcTester()
			tester.test(data)
			p1 := <-tester.chOutput
			_assertData(p1, plain)
			p2 := <-tester.chOutput
			_assertCmd(p2, optcmd, optval)
		})
	})
	Convey("Given [SUB | DATA]", t, func() {
		data := append(pieceSub, plain...)
		Convey("Get: SubPacket + DataPacket", func() {
			tester := _newProcTester()
			tester.test(data)
			p1 := <-tester.chOutput
			_assertSub(p1, subopt, subparam)
			p2 := <-tester.chOutput
			_assertData(p2, plain)
		})
	})
	Convey("Given [DATA | SUB]", t, func() {
		data := append(plain, pieceSub...)
		Convey("Get: DataPacket + SubPacket", func() {
			tester := _newProcTester()
			tester.test(data)
			p1 := <-tester.chOutput
			_assertData(p1, plain)
			p2 := <-tester.chOutput
			_assertSub(p2, subopt, subparam)
		})
	})
	Convey("Given [CMD | DATA | SUB | DATA]", t, func() {
		data := append(pieceOptCmd, append(plain, append(pieceSub, plain...)...)...)
		Convey("Get: CommandPacket + DataPacket + SubPacket + DataPacket", func() {
			tester := _newProcTester()
			tester.test(data)
			p1 := <-tester.chOutput
			_assertCmd(p1, optcmd, optval)
			p2 := <-tester.chOutput
			_assertData(p2, plain)
			p3 := <-tester.chOutput
			_assertSub(p3, subopt, subparam)
			p4 := <-tester.chOutput
			_assertData(p4, plain)
		})
	})
}
