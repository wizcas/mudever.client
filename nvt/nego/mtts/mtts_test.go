package mtts

import (
	"context"
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wizcas/mudever.svc/nvt/nego"
	"github.com/wizcas/mudever.svc/packet"
	"github.com/wizcas/mudever.svc/telbyte"
)

var errMock = errors.New("MOCK")

type mockSender struct {
	bad    bool
	packet packet.Packet
	err    error
}

func (s *mockSender) Send(p packet.Packet) error {
	if s.bad {
		return errMock
	}
	s.packet = p
	return nil
}
func (s *mockSender) OnError(err error) {
	s.err = err
}

func mockMTTSCtx(h nego.OptionHandler, bad bool) (*nego.OptionContext, *mockSender) {
	sender := &mockSender{bad: bad}
	return nego.NewOptionContext(context.Background(), h, sender, sender.OnError), sender
}

func TestNewMTTS(t *testing.T) {
	Convey("By creating a new MTTS with UTF8", t, func() {
		h := New(true)
		Convey("Should have UTF8 flag", func() {
			So(h.Features.has(FeatUTF8), ShouldBeTrue)
		})
	})
	Convey("By creating a new MTTS without UTF8", t, func() {
		h := New(false)
		Convey("Should have no UTF8 flag", func() {
			So(h.Features.has(FeatUTF8), ShouldBeFalse)
		})
	})
}

func TestMTTSOption(t *testing.T) {
	Convey("Given an MTTS handler", t, func() {
		h := New(false)
		Convey("Should have a TTYPE telnet option byte", func() {
			So(h.Option(), ShouldEqual, telbyte.TTYPE)
		})
	})
}

func assertPacketCommand(p packet.Packet, expectCmd telbyte.Command) {
	So(p, ShouldHaveSameTypeAs, &packet.CommandPacket{})
	So(p.(*packet.CommandPacket).Command, ShouldEqual, expectCmd)
}

func TestMTTSHandshake(t *testing.T) {
	Convey("Given an MTTS handler", t, func() {
		h := New(true)
		Convey("With a ctx that can send successfully", func() {
			ctx, snd := mockMTTSCtx(h, false)
			Convey("Should reply to DO with WILL", func() {
				h.Handshake(ctx, telbyte.DO)
				assertPacketCommand(snd.packet, telbyte.WILL)
			})
			Convey("Should reply to DONT with WONT", func() {
				h.Handshake(ctx, telbyte.DONT)
				assertPacketCommand(snd.packet, telbyte.WONT)
			})
			Convey("Should ignore other commands", func() {
				h.Handshake(ctx, telbyte.WILL)
				So(snd.packet, ShouldBeNil)
				So(snd.err, ShouldBeNil)
			})
		})
		Convey("With a ctx that can't send successfully", func() {
			ctx, snd := mockMTTSCtx(h, true)
			Convey("Should get an error on handshaking", func() {
				h.Handshake(ctx, telbyte.DO)
				So(snd.err, ShouldEqual, errMock)
			})
		})
	})
}

func assertSubReply(p packet.Packet, expectBody []byte) {
	So(p, ShouldHaveSameTypeAs, &packet.SubPacket{})
	param := p.(*packet.SubPacket).Parameter
	So(param[0], ShouldEqual, IS)
	So(param[1:], ShouldResemble, expectBody)
}

func TestMTTSSubnego(t *testing.T) {
	Convey("Given an MTTS handler", t, func() {
		h := New(true)
		featureFlag := fmt.Sprintf("MTTS %d", h.Features.Value())
		Convey("With a ctx that can send successfully", func() {
			ctx, snd := mockMTTSCtx(h, false)
			Convey("Should report an error on an empty subnegotiation parameter", func() {
				h.Subnegotiate(ctx, nil)
				So(snd.err, ShouldEqual, nego.ErrLackData)
				h.Subnegotiate(ctx, []byte{})
				So(snd.err, ShouldEqual, nego.ErrLackData)
			})
			Convey("Should ignore an invalid subnegotiation parameter", func() {
				h.Subnegotiate(ctx, []byte{2})
				So(snd.packet, ShouldBeNil)
				So(snd.err, ShouldBeNil)
			})
			Convey("On queryTimes == 0, reply client name", func() {
				h.Subnegotiate(ctx, []byte{SEND})
				assertSubReply(snd.packet, []byte(h.ClientName))
				So(h.queryTimes, ShouldEqual, 1)
			})
			Convey("On queryTimes == 1, reply terminal type", func() {
				h.setQueryTimesForTesting(1)
				h.Subnegotiate(ctx, []byte{SEND})
				assertSubReply(snd.packet, []byte(h.TerminalType))
				So(h.queryTimes, ShouldEqual, 2)
			})
			Convey("On queryTimes == 2, reply feature flag", func() {
				h.setQueryTimesForTesting(2)
				h.Subnegotiate(ctx, []byte{SEND})
				assertSubReply(snd.packet, []byte(featureFlag))
				So(h.queryTimes, ShouldEqual, 3)
			})
			Convey("On queryTimes == 3, reply feature flag", func() {
				h.setQueryTimesForTesting(3)
				h.Subnegotiate(ctx, []byte{SEND})
				assertSubReply(snd.packet, []byte(featureFlag))
				So(h.queryTimes, ShouldEqual, 4)
			})
		})
		Convey("With a ctx that can't send successfully", func() {
			ctx, snd := mockMTTSCtx(h, true)
			Convey("Should get an error on subnegotiating", func() {
				h.Subnegotiate(ctx, []byte{SEND})
				So(snd.err, ShouldEqual, errMock)
			})
		})
	})
}
