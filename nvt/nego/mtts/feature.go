package mtts

import (
	"fmt"
	"strings"
	"unsafe"
)

type featureFlag uint

// MUD Terminal Feature Flags
const (
	FeatANSI = featureFlag(1 << iota)
	FeatVT100
	FeatUTF8
	Feat256Colors
	FeatMouseTracking
	FeatOscColorPalette
	FeatScreenReader
	FeatProxy
	FeatTrueColor
)

var featureNames = map[featureFlag]string{
	FeatANSI:            "ANSI",
	FeatVT100:           "VT100",
	FeatUTF8:            "UTF8",
	Feat256Colors:       "256 Colors",
	FeatMouseTracking:   "Mouse Tracking",
	FeatOscColorPalette: "OSC Color Palette",
	FeatScreenReader:    "Screen Reader",
	FeatProxy:           "Proxy",
	FeatTrueColor:       "True Color",
}

type featureSet uint

func newFeatureSet(flags ...featureFlag) *featureSet {
	f := featureSet(0)
	return (&f).addAll(flags...)
}

func (f *featureSet) add(flag featureFlag) *featureSet {
	*f |= featureSet(flag)
	return f
}

func (f *featureSet) addAll(flags ...featureFlag) *featureSet {
	if flags != nil {
		for _, flag := range flags {
			*f |= featureSet(flag)
		}
	}
	return f
}

func (f *featureSet) remove(flag featureFlag) *featureSet {
	*f &= ^featureSet(flag)
	return f
}

func (f *featureSet) clear() *featureSet {
	*f = 0
	return f
}

func (f *featureSet) has(flag featureFlag) bool {
	return (uint(*f) & uint(flag)) == uint(flag)
}

func (f *featureSet) String() string {
	fval := *f
	bitSize := uint(unsafe.Sizeof(fval) * 8)
	names := []string{}
	for bitIndex := uint(0); bitIndex < bitSize; bitIndex++ {
		bitValue := fval & (1 << bitIndex)
		if bitValue>>bitIndex == 0 {
			continue
		}
		bitName, ok := featureNames[featureFlag(bitValue)]
		if !ok {
			bitName = fmt.Sprintf("%d", bitValue)
		}
		names = append(names, bitName)
	}
	return strings.Join(names, ", ")
}

func (f *featureSet) Value() uint {
	return uint(*f)
}
