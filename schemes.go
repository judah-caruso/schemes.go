package schemes

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"math"
	"strings"
)

type Scheme struct {
	Title string
	C0    SchemeColor
	C1    SchemeColor
	C2    SchemeColor
	C3    SchemeColor
	C4    SchemeColor
	C5    SchemeColor
	C6    SchemeColor
	C7    SchemeColor
}

func (s Scheme) Palette() []SchemeColor {
	return []SchemeColor{
		s.C0,
		s.C1,
		s.C2,
		s.C3,
		s.C4,
		s.C5,
		s.C6,
		s.C7,
	}
}

func (s Scheme) String() string {
	return fmt.Sprintf(
		"%s [ %s, %s, %s, %s, %s, %s, %s, %s ]",
		s.Title,
		s.C0.Hex(),
		s.C1.Hex(),
		s.C2.Hex(),
		s.C3.Hex(),
		s.C4.Hex(),
		s.C5.Hex(),
		s.C6.Hex(),
		s.C7.Hex(),
	)
}

type SchemeColor struct {
	R, G, B uint8
}

func (c SchemeColor) Values() []uint8 {
	return []uint8{
		c.R,
		c.G,
		c.B,
	}
}

func (c SchemeColor) RGB() uint32 {
	return (uint32(c.R) << 16) | (uint32(c.G) << 8) | uint32(c.B)
}

func (c SchemeColor) Hex() string {
	return fmt.Sprintf("#%06X", c.RGB())
}

// Parses SVG scheme source into a 'Scheme'
func ReadScheme(source []byte) (*Scheme, error) {
	scheme := &Scheme{}
	scheme.Title = defaultTitle

	var svg struct {
		Title string `xml:"title"`
		Style string `xml:"style"`
	}

	if err := xml.Unmarshal(source, &svg); err != nil {
		return nil, err
	}

	if len(svg.Title) > 0 {
		scheme.Title = svg.Title
	}

	colors := strings.Split(svg.Style, "#c")
	for _, line := range colors {
		line = strings.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}

		var color *SchemeColor

		id := line[0:1]
		switch id {
		case "0":
			color = &scheme.C0
		case "1":
			color = &scheme.C1
		case "2":
			color = &scheme.C2
		case "3":
			color = &scheme.C3
		case "4":
			color = &scheme.C4
		case "5":
			color = &scheme.C5
		case "6":
			color = &scheme.C6
		case "7":
			color = &scheme.C7
		default: // Ignore IDs we don't care about
			continue
		}

		_, rawColor, found := strings.Cut(line, "fill")
		if !found {
			return nil, fmt.Errorf("Unable to find 'fill' property on #%s", id)
		}

		_, rawColor, found = strings.Cut(rawColor, ":")
		if !found {
			return nil, fmt.Errorf("Malformed property for %s", line)
		}

		end := -1
		for i, chr := range rawColor {
			if chr == ';' {
				end = i
				break
			}
		}

		if end == -1 {
			return nil, fmt.Errorf("Malformed property for %s", line)
		}

		var err error

		switch str := strings.TrimSpace(rawColor[:end]); {
		case strings.Contains(str, "rgb"):
			*color, err = FromRGB(str)
			if err != nil {
				return nil, err
			}

		case strings.Contains(str, "hsv"):
			*color, err = FromHSV(str)
			if err != nil {
				return nil, err
			}

		case str[0] == '#':
			*color, err = FromHex(str)
			if err != nil {
				return nil, err
			}
		}
	}

	return scheme, nil
}

// Takes a Scheme struct and turns it into an SVG scheme
func ExportScheme(s *Scheme) string {
	return fmt.Sprintf(
		schemeTemplate,
		s.Title,
		s.C0.Hex(),
		s.C1.Hex(),
		s.C2.Hex(),
		s.C3.Hex(),
		s.C4.Hex(),
		s.C5.Hex(),
		s.C6.Hex(),
		s.C7.Hex(),
	)
}

// Converts an RGB string (in the format of 'rgb(RR, GG, BB)') to aScheme_Color
func FromRGB(str string) (SchemeColor, error) {
	var color SchemeColor

	_, err := fmt.Sscanf(str, "rgb(%d,%d,%d)", &color.R, &color.G, &color.B)
	if err != nil {
		return color, err
	}

	return color, nil
}

// Converts an HSV string (in the format of 'hsv(HH, SS, VV)') to aScheme_Color
func FromHSV(str string) (SchemeColor, error) {
	var hue float64
	var sat float64
	var val float64

	_, err := fmt.Sscanf(str, "hsv(%f,%f,%f)", &hue, &sat, &val)
	if err != nil {
		return SchemeColor{}, err
	}

	hue = clamp(hue, 0, 360) / 360
	sat = clamp(sat, 0, 100) / 100
	val = clamp(val, 0, 100) / 100

	r, g, b := func(hue float64, sat float64, val float64) (float64, float64, float64) {
		px := math.Abs(fract(hue+1)*6-3) - 1
		py := math.Abs(fract(hue+2/3.0)*6-3) - 1
		pz := math.Abs(fract(hue+1/3.0)*6-3) - 1

		px = clamp(px, 0, 1)
		py = clamp(px, 0, 1)
		pz = clamp(px, 0, 1)

		px = lerp(1, px, sat)
		py = lerp(1, py, sat)
		pz = lerp(1, pz, sat)

		return val * px, val * py, val * pz
	}(hue, sat, val)

	return SchemeColor{
		uint8(math.Ceil(r * 255)),
		uint8(math.Ceil(g * 255)),
		uint8(math.Ceil(b * 255)),
	}, nil
}

// Converts a hex string (in the format of '#RRGGBB') to aScheme_Color
func FromHex(str string) (SchemeColor, error) {
	hex, err := hex.DecodeString(str[1:])
	if err != nil {
		return SchemeColor{}, err
	}

	return SchemeColor{
		hex[0],
		hex[1],
		hex[2],
	}, nil
}

func clamp(v float64, min float64, max float64) float64 {
	return math.Min(math.Max(v, min), max)
}

func fract(x float64) float64 {
	return x - math.Floor(x)
}

func lerp(a float64, b float64, t float64) float64 {
	return a + (b-a)*t
}

const (
	defaultTitle = "Color Scheme by Person"

	schemeTemplate = `<svg width="288px" height="140px" xmlns="http://www.w3.org/2000/svg" baseProfile="full" version="1.1">
   <title>%s</title>
   <style>
      #c0 { fill: %s; } <!-- Background -->
      #c1 { fill: %s; } <!-- Foreground, Operators --> 
      #c2 { fill: %s; } <!-- Types --> 
      #c3 { fill: %s; } <!-- Procedures, Keywords --> 
      #c4 { fill: %s; } <!-- Constants, Strings -->
      #c5 { fill: %s; } <!-- Pre-Processor, Special -->
      #c6 { fill: %s; } <!-- Errors -->
      #c7 { fill: %s; } <!-- Comments --> 
   </style>

   <rect width="288" height="140" id="c0"></rect>   
   <text x="10" y="20" style="font-family: monospace" id="c1"><tspan x="10"><tspan id="c5">import</tspan> <tspan id="c4">"fmt"</tspan> <tspan id="c7">// package main</tspan></tspan><tspan x="10" dy="2em"><tspan id="c3">type</tspan> Point <tspan id="c2">struct</tspan> {</tspan><tspan x="25" dy="1em">x, y <tspan id="c2">float32</tspan></tspan><tspan x="10" dy="1em">}</tspan><tspan x="10" dy="2em"><tspan id="c2">func</tspan> <tspan id="c3">main</tspan>() {</tspan><tspan x="25" dy="1em">fmt.println(<tspan id="c4">"Hello, World"</tspan>)</tspan><tspan x="10" dy="1em">}</tspan>
</text>
</svg>`
)
