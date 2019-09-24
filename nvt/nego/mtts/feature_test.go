package mtts

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewFeature(t *testing.T) {
	Convey("Create an emtpy featureSet", t, func() {
		f := newFeatureSet()
		Convey("It should have value 0", func() {
			So(f.Value(), ShouldEqual, 0)
		})
	})
	Convey("Create a featureSet with given flags", t, func() {
		f := newFeatureSet(FeatANSI, FeatVT100)
		Convey("It should have value of the initiating flags", func() {
			So(f.Value(), ShouldEqual, 3)
		})
	})
}

func TestAddFeature(t *testing.T) {
	Convey("Given an empty featureSet", t, func() {
		f := newFeatureSet()
		Convey("Add a single flag", func() {
			f.add(FeatANSI)
			Convey("The flag should be set", func() {
				So(f.Value(), ShouldEqual, 1)
				So(f.has(FeatANSI), ShouldBeTrue)
			})
		})
		Convey("Add a single flag for multiple times", func() {
			f.add(FeatVT100)
			f.add(FeatVT100)
			f.add(FeatVT100)
			Convey("The flag should be set only once", func() {
				So(f.Value(), ShouldEqual, 2)
				So(f.has(FeatVT100), ShouldBeTrue)
				So(f.has(FeatANSI), ShouldBeFalse)
			})
		})
		Convey("Add multiple distinct values as a batch", func() {
			f.addAll(Feat256Colors, FeatUTF8, FeatProxy)
			Convey("All flags should be set", func() {
				So(f.Value(), ShouldEqual, 1<<2+1<<3+1<<7)
				So(f.has(FeatUTF8), ShouldBeTrue)
				So(f.has(Feat256Colors), ShouldBeTrue)
				So(f.has(FeatProxy), ShouldBeTrue)
			})
		})
		Convey("Add multiple duplicated values as a batch", func() {
			f.addAll(FeatTrueColor, FeatTrueColor, FeatTrueColor)
			Convey("All flags should be set", func() {
				So(f.Value(), ShouldEqual, 1<<8)
				So(f.has(FeatTrueColor), ShouldBeTrue)
			})
		})
	})
}

func TestRemoveFeature(t *testing.T) {
	Convey("Given a featureSet set with some flags", t, func() {
		f := newFeatureSet(FeatANSI, FeatVT100, Feat256Colors)
		Convey("Remove one of the set flags", func() {
			f = f.remove(FeatVT100)
			Convey("Should contain only the remaing flags", func() {
				So(f.Value(), ShouldEqual, 1<<0+1<<3)
				So(f.has(FeatVT100), ShouldBeFalse)
			})
		})
		Convey("Remove a flag that is not set", func() {
			f = f.remove(FeatUTF8)
			Convey("Should remain unchanged", func() {
				So(f.Value(), ShouldEqual, 1<<0+1<<1+1<<3)
			})
		})
		Convey("Clear the featureSet", func() {
			f.clear()
			Convey("Should have value 0", func() {
				So(f.Value(), ShouldEqual, 0)
			})
		})
	})
}

func TestFeatureString(t *testing.T) {
	Convey("Given a featureSet set with some valid flags", t, func() {
		f := newFeatureSet(FeatANSI, FeatVT100, Feat256Colors)
		Convey("Should be able to print the flags' names", func() {
			s := f.String()
			So(s, ShouldEqual, "ANSI, VT100, 256 Colors")
		})
	})
	Convey("Given a featureSet set with flags including invalid ones", t, func() {
		f := newFeatureSet(FeatANSI, FeatVT100, 1<<28)
		Convey("Should be able to print the names of valid flags as well as numeric value of invalid flags", func() {
			s := f.String()
			So(s, ShouldResemble, "ANSI, VT100, 268435456")
		})
	})
}
