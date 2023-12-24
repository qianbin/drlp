package drlp

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

type test struct {
	val  interface{}
	want string
}

var tests = []test{
	// integers
	{0, "80"},
	{127, "7F"},
	{128, "8180"},
	{256, "820100"},
	{1024, "820400"},
	{0xFFFFFF, "83FFFFFF"},
	{0xFFFFFFFF, "84FFFFFFFF"},
	{0xFFFFFFFF, "84FFFFFFFF"},
	{0xFFFFFFFFFF, "85FFFFFFFFFF"},
	{0xFFFFFFFFFFFF, "86FFFFFFFFFFFF"},
	{0xFFFFFFFFFFFFFF, "87FFFFFFFFFFFFFF"},
	{uint64(0xFFFFFFFFFFFFFFFF), "88FFFFFFFFFFFFFFFF"},

	// strings
	{[]byte{}, "80"},
	{[]byte{0x7E}, "7E"},
	{[]byte{0x7F}, "7F"},
	{[]byte{0x80}, "8180"},
	{[]byte{1, 2, 3}, "83010203"},
	{
		"Lorem ipsum dolor sit amet, consectetur adipisicing eli",
		"B74C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C69",
	},
	{
		"Lorem ipsum dolor sit amet, consectetur adipisicing elit",
		"B8384C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C6974",
	},
	{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris magna, suscipit sed vehicula non, iaculis faucibus tortor. Proin suscipit ultricies malesuada. Duis tortor elit, dictum quis tristique eu, ultrices at risus. Morbi a est imperdiet mi ullamcorper aliquet suscipit nec lorem. Aenean quis leo mollis, vulputate elit varius, consequat enim. Nulla ultrices turpis justo, et posuere urna consectetur nec. Proin non convallis metus. Donec tempor ipsum in mauris congue sollicitudin. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Suspendisse convallis sem vel massa faucibus, eget lacinia lacus tempor. Nulla quis ultricies purus. Proin auctor rhoncus nibh condimentum mollis. Aliquam consequat enim at metus luctus, a eleifend purus egestas. Curabitur at nibh metus. Nam bibendum, neque at auctor tristique, lorem libero aliquet arcu, non interdum tellus lectus sit amet eros. Cras rhoncus, metus ac ornare cursus, dolor justo ultrices metus, at ullamcorper volutpat",
		"B904004C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E73656374657475722061646970697363696E6720656C69742E20437572616269747572206D6175726973206D61676E612C20737573636970697420736564207665686963756C61206E6F6E2C20696163756C697320666175636962757320746F72746F722E2050726F696E20737573636970697420756C74726963696573206D616C6573756164612E204475697320746F72746F7220656C69742C2064696374756D2071756973207472697374697175652065752C20756C7472696365732061742072697375732E204D6F72626920612065737420696D70657264696574206D6920756C6C616D636F7270657220616C6971756574207375736369706974206E6563206C6F72656D2E2041656E65616E2071756973206C656F206D6F6C6C69732C2076756C70757461746520656C6974207661726975732C20636F6E73657175617420656E696D2E204E756C6C6120756C74726963657320747572706973206A7573746F2C20657420706F73756572652075726E6120636F6E7365637465747572206E65632E2050726F696E206E6F6E20636F6E76616C6C6973206D657475732E20446F6E65632074656D706F7220697073756D20696E206D617572697320636F6E67756520736F6C6C696369747564696E2E20566573746962756C756D20616E746520697073756D207072696D697320696E206661756369627573206F726369206C756374757320657420756C74726963657320706F737565726520637562696C69612043757261653B2053757370656E646973736520636F6E76616C6C69732073656D2076656C206D617373612066617563696275732C2065676574206C6163696E6961206C616375732074656D706F722E204E756C6C61207175697320756C747269636965732070757275732E2050726F696E20617563746F722072686F6E637573206E69626820636F6E64696D656E74756D206D6F6C6C69732E20416C697175616D20636F6E73657175617420656E696D206174206D65747573206C75637475732C206120656C656966656E6420707572757320656765737461732E20437572616269747572206174206E696268206D657475732E204E616D20626962656E64756D2C206E6571756520617420617563746F72207472697374697175652C206C6F72656D206C696265726F20616C697175657420617263752C206E6F6E20696E74657264756D2074656C6C7573206C65637475732073697420616D65742065726F732E20437261732072686F6E6375732C206D65747573206163206F726E617265206375727375732C20646F6C6F72206A7573746F20756C747269636573206D657475732C20617420756C6C616D636F7270657220766F6C7574706174",
	},
	// lists
	{[]interface{}{}, "C0"},
	{[]interface{}{1, 2, 3}, "C3010203"},
	{
		// [ [], [[]], [ [], [[]] ] ]
		[]interface{}{[]interface{}{}, []interface{}{[]interface{}{}}, []interface{}{[]interface{}{}, []interface{}{[]interface{}{}}}},
		"C7C0C1C0C3C0C1C0",
	},
	{
		[]interface{}{"aaa", "bbb", "ccc", "ddd", "eee", "fff", "ggg", "hhh", "iii", "jjj", "kkk", "lll", "mmm", "nnn", "ooo"},
		"F83C836161618362626283636363836464648365656583666666836767678368686883696969836A6A6A836B6B6B836C6C6C836D6D6D836E6E6E836F6F6F",
	},
	{
		[]interface{}{1, 0xFFFFFF, []interface{}{[]interface{}{4, 5, 5}}, "abc"},
		"CE0183FFFFFFC4C304050583616263",
	},
	{[]interface{}{
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"},
		[]interface{}{"asdf", "qwer", "zxcv"}},
		"F90200CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376"},
}

func encode(v interface{}, b *Buffer) {
	switch v := v.(type) {
	case []interface{}:
		l := b.List()
		for _, c := range v {
			encode(c, b)
		}
		l.End()
	case int:
		b.PutUint(uint64(v))
	case uint64:
		b.PutUint(v)
	case []byte:
		b.PutString(v)
	case string:
		b.PutString([]byte(v))
	default:
		panic(errors.New("unexpected value type"))
	}
}

func TestBuffer(t *testing.T) {
	var buf Buffer
	for i, tt := range tests {
		buf.Reset()
		encode(tt.val, &buf)
		if got := strings.ToUpper(hex.EncodeToString(buf)); got != tt.want {
			t.Errorf("#%v: got %v, want %v", i, got, tt.want)
		}
	}
}
