package telnet

import (
	"io"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRead(t *testing.T) {
	Convey("Given a string reader of a 15-character-string", t, func() {
		rawdata := "Heeeelllllooooo"
		sr := strings.NewReader(rawdata)
		r := newReader(sr)
		Convey("With a small buffer of size 5", func() {
			l := 5
			buf := make([]byte, l)
			Convey("only 5 characters can be read for once", func() {
				n, err := r.read(buf)
				So(n, ShouldEqual, l)
				So(err, ShouldBeNil)
				So(string(buf[:n]), ShouldEqual, "Heeee")
			})
		})
		Convey("With a buffer of exact 15 bytes long", func() {
			l := 15
			buf := make([]byte, l)
			Convey("All 15 characters are read", func() {
				n, err := r.read(buf)
				So(n, ShouldEqual, l)
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, rawdata)
			})
		})
		Convey("With a buffer larger than 15", func() {
			l := 20
			buf := make([]byte, l)
			Convey("15 characters are read after reaching eof", func() {
				n, err := r.read(buf)
				So(n, ShouldEqual, 15)
				So(err, ShouldEqual, io.EOF)
				So(string(buf[:n]), ShouldEqual, rawdata)
			})
		})
	})
}
