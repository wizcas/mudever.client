package packet

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func _makePacket(data []byte) *DataPacket {
	return &DataPacket{data}
}

func TestNoEscapeSerialize(t *testing.T) {
	Convey("Given bytes without 0xFF", t, func() {
		data := []byte{1, 2, 3, 4, 5}
		p := _makePacket(data)
		Convey("Serialization returns the same value", func() {
			result := p.Serialize()
			So(result, ShouldResemble, data)
			So(result, ShouldNotEqual, data)
		})
	})
}

func TestEscapeSerialize1(t *testing.T) {
	Convey("Given bytes with 0xFF in the beginning", t, func() {
		data := []byte{255, 2, 3, 4, 5}
		expect := []byte{255, 255, 2, 3, 4, 5}
		p := _makePacket(data)
		Convey("Serialization make it escaped", func() {
			result := p.Serialize()
			So(result, ShouldResemble, expect)
			So(result, ShouldNotEqual, data)
		})
	})
}

func TestEscapeSerialize2(t *testing.T) {
	Convey("Given bytes with 0xFF in middle", t, func() {
		data := []byte{1, 2, 255, 4, 5}
		expect := []byte{1, 2, 255, 255, 4, 5}
		p := _makePacket(data)
		Convey("Serialization make it escaped", func() {
			result := p.Serialize()
			So(result, ShouldResemble, expect)
			So(result, ShouldNotEqual, data)
		})
	})
}

func TestEscapeSerialize3(t *testing.T) {
	Convey("Given bytes with 0xFF in the end", t, func() {
		data := []byte{1, 2, 3, 4, 255}
		expect := []byte{1, 2, 3, 4, 255, 255}
		p := _makePacket(data)
		Convey("Serialization make it escaped", func() {
			result := p.Serialize()
			So(result, ShouldResemble, expect)
			So(result, ShouldNotEqual, data)
		})
	})
}

func TestEscapeSerialize4(t *testing.T) {
	Convey("Given bytes with continuous 0xFF", t, func() {
		data := []byte{1, 2, 3, 255, 255}
		expect := []byte{1, 2, 3, 255, 255, 255, 255}
		p := _makePacket(data)
		Convey("Serialization make them escaped", func() {
			result := p.Serialize()
			So(result, ShouldResemble, expect)
			So(result, ShouldNotEqual, data)
		})
	})
}
