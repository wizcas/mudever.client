package packet

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func _assertData(result, expect, raw []byte, err error) {
	So(result, ShouldResemble, expect)
	So(result, ShouldNotEqual, raw)
	So(err, ShouldBeNil)
}

func TestDataStringify(t *testing.T) {
	Convey("Given a data packet", t, func() {
		p := NewDataPacket(make([]byte, 22))
		Convey("Print its length", func() {
			So(p.String(), ShouldEqual, "[TXT] (22 bytes)")
		})
	})
}

func TestDataKind(t *testing.T) {
	Convey("Data packet is of KindData", t, func() {
		So(NewDataPacket([]byte{}).GetKind(), ShouldEqual, KindData)
	})
}

func TestNoEscapeSerialize(t *testing.T) {
	Convey("Given bytes without 0xFF", t, func() {
		data := []byte{1, 2, 3, 4, 5}
		p := NewDataPacket(data)
		Convey("Serialization returns the same value", func() {
			result, err := p.Serialize()
			_assertData(result, data, data, err)
		})
	})
}

func TestEscapeSerialize1(t *testing.T) {
	Convey("Given bytes with 0xFF in the beginning", t, func() {
		data := []byte{255, 2, 3, 4, 5}
		expect := []byte{255, 255, 2, 3, 4, 5}
		p := NewDataPacket(data)
		Convey("Serialization make it escaped", func() {
			result, err := p.Serialize()
			_assertData(result, expect, data, err)

		})
	})
}

func TestEscapeSerialize2(t *testing.T) {
	Convey("Given bytes with 0xFF in middle", t, func() {
		data := []byte{1, 2, 255, 4, 5}
		expect := []byte{1, 2, 255, 255, 4, 5}
		p := NewDataPacket(data)
		Convey("Serialization make it escaped", func() {
			result, err := p.Serialize()
			_assertData(result, expect, data, err)
		})
	})
}

func TestEscapeSerialize3(t *testing.T) {
	Convey("Given bytes with 0xFF in the end", t, func() {
		data := []byte{1, 2, 3, 4, 255}
		expect := []byte{1, 2, 3, 4, 255, 255}
		p := NewDataPacket(data)
		Convey("Serialization make it escaped", func() {
			result, err := p.Serialize()
			_assertData(result, expect, data, err)
		})
	})
}

func TestEscapeSerialize4(t *testing.T) {
	Convey("Given bytes with continuous 0xFF", t, func() {
		data := []byte{1, 2, 3, 255, 255}
		expect := []byte{1, 2, 3, 255, 255, 255, 255}
		p := NewDataPacket(data)
		Convey("Serialization make them escaped", func() {
			result, err := p.Serialize()
			_assertData(result, expect, data, err)
		})
	})
}
