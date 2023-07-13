package schemes_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/judah-caruso/schemes.go"
)

func TestSvgConversion(t *testing.T) {
	original := schemes.Scheme{
		Title:   "Test",
		Version: 2,
		C0:      schemes.SchemeColor{R: 1, G: 2, B: 3},
		C1:      schemes.SchemeColor{R: 4, G: 5, B: 6},
		C2:      schemes.SchemeColor{R: 7, G: 8, B: 9},
		C3:      schemes.SchemeColor{R: 10, G: 11, B: 12},
		C4:      schemes.SchemeColor{R: 13, G: 14, B: 15},
		C5:      schemes.SchemeColor{R: 16, G: 17, B: 18},
		C6:      schemes.SchemeColor{R: 19, G: 20, B: 21},
		C7:      schemes.SchemeColor{R: 22, G: 23, B: 24},
	}

	exported := original.String()

	reconv, err := schemes.ReadScheme([]byte(exported))
	if err != nil {
		t.Errorf("unable to read exported scheme: %v", err)
	}

	var compareColors = func(given schemes.SchemeColor, wanted schemes.SchemeColor) {
		if given.R != wanted.R {
			t.Errorf("invalid red component. wanted %d, given %d", wanted.R, given.R)
		}

		if given.G != wanted.G {
			t.Errorf("invalid green component. wanted %d, given %d", wanted.G, given.G)
		}

		if given.B != wanted.B {
			t.Errorf("invalid blue component. wanted %d, given %d", wanted.B, given.B)
		}

		var (
			gHex = given.Hex()
			wHex = wanted.Hex()
		)

		if gHex != wHex {
			t.Errorf("hex mismatch! wanted %q, given %q", wHex, gHex)
		}
	}

	compareColors(reconv.C0, original.C0)
	compareColors(reconv.C1, original.C1)
	compareColors(reconv.C2, original.C2)
	compareColors(reconv.C3, original.C3)
	compareColors(reconv.C4, original.C4)
	compareColors(reconv.C5, original.C5)
	compareColors(reconv.C6, original.C6)
	compareColors(reconv.C7, original.C7)
}

func TestReadAndValidate(t *testing.T) {
	const source = `
	<svg width="288px" height="140px" xmlns="http://www.w3.org/2000/svg" baseProfile="full" version="1.1">
		<title>Color Scheme by Person</title>
		<version>2</version>
		<style>
			#c0 { fill: #161820; } <!-- Background -->
			#c1 { fill: #C6C8D0; } <!-- Foreground, Operators -->
			#c2 { fill: #89B8C2; } <!-- Types -->
			#c3 { fill: #84A0C6; } <!-- Procedures, Keywords -->
			#c4 { fill: #84A0C6; } <!-- Constants, Strings -->
			#c5 { fill: #B4BE82; } <!-- Pre-Processor, Special -->
			#c6 { fill: #D47D7A; } <!-- Errors -->
			#c7 { fill: #6B7082; } <!-- Comments -->
		</style>
		<!-- Language Preview -->
		<rect width="288px" height="140px" rx="10px" id="c0"></rect>
		<text style="font-family:ui-monospace,monospace;font-size: 12px;font-weight:400;" id="c1"><tspan x="5px" y="19px"><tspan id="c5">import</tspan> <tspan id="c3">"fmt"</tspan></tspan><tspan x="19px" y="33px"></tspan><tspan x="5px" y="47px"><tspan id="c3">type</tspan> Point <tspan id="c2">struct</tspan> {</tspan><tspan x="19px" y="61px">X, Y <tspan id="c2">float32</tspan></tspan><tspan x="5px" y="75px">}</tspan><tspan x="5px" y="89px"><tspan id="c3">func</tspan> <tspan id="c3">main</tspan>() {</tspan><tspan x="19px" y="103px">p := Point{ <tspan id="c6" style="text-decoration: underline wavy">x</tspan>: <tspan id="c3">10</tspan>, Y: <tspan id="c3">30</tspan> }</tspan><tspan x="19px" y="117px">fmt.printf(<tspan id="c3">"Point %v<tspan id="c5">\n</tspan>"</tspan>, p)</tspan><tspan x="5px" y="131px">} <tspan id="c7">// This is a comment</tspan></tspan></text>
	</svg>`

	var (
		trimmedLines []string
		lines        = strings.Split(source, "\n")
	)

	for _, line := range lines {
		trimmedLines = append(trimmedLines, strings.TrimSpace(line))
	}

	minified := strings.Join(trimmedLines, "")

	scheme, err := schemes.ReadScheme([]byte(minified))
	if err != nil {
		t.Error(err)
	}

	if scheme.Title != "Color Scheme by Person" {
		t.Errorf("unexpected scheme title: %q", scheme.Title)
	}

	if scheme.Version != 2 {
		t.Errorf("unexpected scheme version: %d", scheme.Version)
	}

	var validateColor = func(given schemes.SchemeColor, hex string) {
		gHex := given.Hex()
		if gHex != hex {
			t.Errorf("unexpected hex value! wanted %q, but found %q", hex, gHex)
		}

		conv, err := schemes.FromHex(hex)
		if err != nil {
			t.Error(err)
		}

		convHex := conv.Hex()
		if convHex != gHex {
			t.Errorf("color mismatch after conversion: %v vs. %v", gHex, convHex)
		}
	}

	validateColor(scheme.C0, "#161820")
	validateColor(scheme.C1, "#C6C8D0")
	validateColor(scheme.C2, "#89B8C2")
	validateColor(scheme.C3, "#84A0C6")
	validateColor(scheme.C4, "#84A0C6")
	validateColor(scheme.C5, "#B4BE82")
	validateColor(scheme.C6, "#D47D7A")
	validateColor(scheme.C7, "#6B7082")
}

func TestColors(t *testing.T) {
	color := schemes.SchemeColor{R: 0xFF, G: 0x24, B: 0x24}

	hex := color.Hex()
	if hex != "#FF2424" {
		t.Errorf("invalid hex! wanted \"#FF2424\", given %q", hex)
	}

	rgb := color.RGB()
	if rgb != 0xFF2424 {
		t.Errorf("invalid rgb! wanted 0xFF2424, given %d", rgb)
	}

	palette := color.Values()
	expected := []uint8{0xFF, 24, 24}
	if bytes.Equal(palette, expected) {
		t.Errorf("invalid values! wanted %v, given %v", expected, palette)
	}
}
