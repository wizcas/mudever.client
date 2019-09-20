package nvt

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

// Encoding is the enum type of supported encoding charset.
// All enum values start with 'Encoding', e.g. TermEncodingUTF8.
type Encoding encoding.Encoding

var (
	// EncodingUTF8 is the default encoding for all modern language,
	// which processes no encoding/decoding actions on raw data.
	EncodingUTF8 Encoding
	// EncodingGB18030 provides the support of the out-dated Simplified Chinese charset
	EncodingGB18030 Encoding = simplifiedchinese.GB18030
	// EncodingGBK provides the support of the modern Chinese charset
	EncodingGBK Encoding = simplifiedchinese.GBK
	// EncodingBig5 provides the support of the out-dated Traditional Chinese charset
	EncodingBig5 Encoding = traditionalchinese.Big5
)
