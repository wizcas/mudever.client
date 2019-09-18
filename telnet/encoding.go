package telnet

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

// TermEncoding is the enum type of supported encoding charset.
// All enum values start with 'TermEncoding', e.g. TermEncodingUTF8.
type TermEncoding encoding.Encoding

var (
	// TermEncodingUTF8 is the default encoding for all modern language,
	// which processes no encoding/decoding actions on raw data.
	TermEncodingUTF8 TermEncoding
	// TermEncodingGB18030 provides the support of the out-dated Simplified Chinese charset
	TermEncodingGB18030 TermEncoding = simplifiedchinese.GB18030
	// TermEncodingGBK provides the support of the modern Chinese charset
	TermEncodingGBK TermEncoding = simplifiedchinese.GBK
	// TermEncodingBig5 provides the support of the out-dated Traditional Chinese charset
	TermEncodingBig5 TermEncoding = traditionalchinese.Big5
)
