package stream

import (
	"io"
	"log"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRead(t *testing.T) {
	Convey("Given a string reader of a 15-character-string", t, func() {
		rawdata := "Heeeelllllooooo"
		sr := strings.NewReader(rawdata)
		r := NewReader(sr)
		Convey("With a small buffer of size 5", func() {
			l := 5
			buf := make([]byte, l)
			Convey("only 5 characters can be read for once", func() {
				n, err := r.Read(buf)
				So(n, ShouldEqual, l)
				So(err, ShouldBeNil)
				So(string(buf[:n]), ShouldEqual, "Heeee")
			})
			Convey("finish reading to eof with 5 read calls", func() {
				count := 0
				expects := []string{
					"Heeee",
					"lllll",
					"ooooo",
				}
				for count < 5 {
					count++
					log.Printf("read %d", count)
					n, err := r.Read(buf)
					if count < 4 {
						So(n, ShouldEqual, 5)
						s := string(buf[:n])
						So(s, ShouldEqual, expects[count-1])
					} else if count == 4 {
						So(n, ShouldEqual, 0)
						So(err, ShouldBeNil)
					} else {
						So(n, ShouldEqual, 0)
						So(err, ShouldBeError, io.EOF)
						break
					}
				}
				So(count, ShouldEqual, 5)
			})
		})
		Convey("With a buffer of exact 15 bytes long", func() {
			l := 15
			buf := make([]byte, l)
			Convey("Reads the entire string on 1st read, EOP on 2nd read, and EOF on 3rd read", func() {
				n, err := r.Read(buf)
				So(n, ShouldEqual, l)
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, rawdata)
				n, err = r.Read(buf)
				So(n, ShouldEqual, 0)
				So(err, ShouldEqual, nil)
				n, err = r.Read(buf)
				So(n, ShouldEqual, 0)
				So(err, ShouldEqual, io.EOF)
			})
		})
		Convey("With a buffer larger than 15", func() {
			l := 20
			buf := make([]byte, l)
			Convey("Reads the entire string with ErrEOP, then returns empty buffer with EOF on next read", func() {
				n, err := r.Read(buf)
				So(n, ShouldEqual, 15)
				So(err, ShouldEqual, nil)
				So(string(buf[:n]), ShouldEqual, rawdata)
				n, err = r.Read(buf)
				So(n, ShouldEqual, 0)
				So(err, ShouldEqual, nil)
				n, err = r.Read(buf)
				So(n, ShouldEqual, 0)
				So(err, ShouldEqual, io.EOF)
			})
		})
	})
}
