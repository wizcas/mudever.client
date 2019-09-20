package receiver

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const base = byte((1 << 5) + (1 << 1)) // 00100010

func TestAddState(t *testing.T) {
	Convey("Given a state of not nil value", t, func() {
		state := &processor{value: base}
		Convey("With a bit flag which is not set", func() {
			flag := byte(1 << 2)
			Convey("The bit should be set correctly with addState()", func() {
				expect := byte((1 << 5) + (1 << 2) + (1 << 1))
				state.addState(flag)
				So(state.value, ShouldEqual, expect)
			})
		})
		Convey("With a bit flag which is already set", func() {
			flag := byte(1 << 1)
			Convey("The state's value should keep unchanged", func() {
				state.addState(flag)
				So(state.value, ShouldEqual, base)
			})
		})
	})
}

func TestDelState(t *testing.T) {
	Convey("Given a bit flag and a state with base value", t, func() {
		state := &processor{value: base}
		Convey("With a bit flag which is set", func() {
			flag := byte(1 << 5)
			Convey("The bit should be unset correctly with delState()", func() {
				expect := byte(1 << 1)
				state.delState(flag)
				So(state.value, ShouldEqual, expect)
			})
		})
		Convey("With a bit flag which is not set yet", func() {
			flag := byte(1 << 6)
			Convey("The state's value should keep unchanged", func() {
				state.delState(flag)
				So(state.value, ShouldEqual, base)
			})
		})
	})
}

func TestInState(t *testing.T) {
	Convey("Given a state", t, func() {
		state := &processor{value: base}
		Convey("With a bit flag that is set", func() {
			flag := byte(1 << 1)
			Convey("It should pass the inState() test", func() {
				So(state.inState(flag), ShouldBeTrue)
			})
		})
		Convey("With a bit flag that is not set", func() {
			flag := byte(1 << 2)
			Convey("It should fail the inState() test", func() {
				So(state.inState(flag), ShouldBeFalse)
			})
		})
	})
}
