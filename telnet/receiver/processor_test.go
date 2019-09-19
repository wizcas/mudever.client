package receiver

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wizcas/mudever.svc/telnet/packet"
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

func TestProcPlainData(t *testing.T) {
	Convey("Given plain data", t, func() {
		data := []byte("I'm plain data")
		Convey("Processor should generate a data packet", func() {
			tester := _newProcTester()
			go tester.proc(data, tester.chOutput, tester.chErr)
			p := <-tester.chOutput
			So(p, ShouldHaveSameTypeAs, &packet.DataPacket{})
			So(p.(*packet.DataPacket).Data, ShouldResemble, data)
		})
	})
	Convey("Given plain data with 0xFF", t, func() {
		data := append(append([]byte("I'm plain"), []byte{255, 255}...), []byte("data")...)
		expect := append(append([]byte("I'm plain"), []byte{255}...), []byte("data")...)
		Convey("Processor should unescape 0xFF", func() {
			tester := _newProcTester()
			go tester.proc(data, tester.chOutput, tester.chErr)
			p := <-tester.chOutput
			So(p, ShouldHaveSameTypeAs, &packet.DataPacket{})
			So(p.(*packet.DataPacket).Data, ShouldResemble, expect)
		})
	})
}
