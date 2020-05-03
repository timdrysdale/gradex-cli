package parsesvg

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/mattetti/filebuffer"
	"github.com/timdrysdale/geo"
	"github.com/timdrysdale/unipdf/v3/annotator"
	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/model"
	"github.com/timdrysdale/unipdf/v3/model/optimize"
)

const testInkscapeSvg = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!-- Created with Inkscape (http://www.inkscape.org/) -->

<svg
   xmlns:dc="http://purl.org/dc/elements/1.1/"
   xmlns:cc="http://creativecommons.org/ns#"
   xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
   xmlns:svg="http://www.w3.org/2000/svg"
   xmlns="http://www.w3.org/2000/svg"
   xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
   xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
   width="50mm"
   height="50mm"
   viewBox="0 0 50 50"
   version="1.1"
   id="svg47244"
   inkscape:version="0.92.4 (5da689c313, 2019-01-14)"
   sodipodi:docname="example.svg">
  <title
     id="title51367">example</title>
  <defs
     id="defs47238" />
  <sodipodi:namedview
     id="base"
     pagecolor="#ffffff"
     bordercolor="#666666"
     borderopacity="1.0"
     inkscape:pageopacity="0.0"
     inkscape:pageshadow="2"
     inkscape:zoom="1.979899"
     inkscape:cx="79.390843"
     inkscape:cy="25.175961"
     inkscape:document-units="mm"
     inkscape:current-layer="layer8"
     showgrid="false"
     inkscape:window-width="1850"
     inkscape:window-height="1136"
     inkscape:window-x="70"
     inkscape:window-y="27"
     inkscape:window-maximized="1"
     inkscape:snap-page="true"
     inkscape:snap-center="true"
     inkscape:snap-global="false" />
  <metadata
     id="metadata47241">
    <rdf:RDF>
      <cc:Work
         rdf:about="">
        <dc:format>image/svg+xml</dc:format>
        <dc:type
           rdf:resource="http://purl.org/dc/dcmitype/StillImage" />
        <dc:title>example</dc:title>
      </cc:Work>
    </rdf:RDF>
  </metadata>
  <g
     inkscape:label="chrome"
     inkscape:groupmode="layer"
     id="layer1"
     transform="translate(0,-247)"
     style="display:inline"
     sodipodi:insensitive="true">
    <g
       id="g20789"
       transform="matrix(0.83470821,0,0,0.82232123,-166.30073,-13.565732)"
       style="stroke-width:1.20701277">
      <a
         id="a17055"
         transform="translate(-22.054358,-11.322927)"
         style="stroke-width:1.20701277">
        <g
           id="g17044"
           transform="matrix(1,0,0,-1,12.843917,688.24148)"
           style="stroke-width:1.20701277">
          <rect
             style="opacity:1;fill:#fdcccb;fill-opacity:1;stroke:#ed1a1f;stroke-width:0.90525961;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1"
             id="rect16988"
             width="24.25"
             height="19.25"
             x="240.98444"
             y="327.56677"
             ry="1.2991039" />
          <rect
             ry="0.96345305"
             y="327.74417"
             x="241.16528"
             height="5.3170075"
             width="23.901506"
             id="rect16990"
             style="opacity:1;fill:#ed1a1f;fill-opacity:1;stroke:none;stroke-width:0.90525955;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
          <path
             inkscape:connector-curvature="0"
             id="path16998"
             d="m 241.24277,333.15092 h 24"
             style="fill:none;stroke:#ed1a1f;stroke-width:0.90525961;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
        </g>
      </a>
      <g
         transform="translate(-444.57853,80.364147)"
         id="g17050"
         style="stroke-width:1.20701277">
        <path
           style="opacity:1;fill:#ed1a1f;fill-opacity:1;stroke:none;stroke-width:0.90525961;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1"
           id="path17046"
           sodipodi:type="arc"
           sodipodi:cx="682.55292"
           sodipodi:cy="256.52966"
           sodipodi:rx="4.6202097"
           sodipodi:ry="4.6202097"
           sodipodi:start="0.02702045"
           sodipodi:end="0"
           d="m 687.17144,256.65449 a 4.6202097,4.6202097 0 0 1 -4.71215,4.49443 4.6202097,4.6202097 0 0 1 -4.52616,-4.68168 4.6202097,4.6202097 0 0 1 4.651,-4.55768 4.6202097,4.6202097 0 0 1 4.589,4.6201"
           sodipodi:open="true" />
        <path
           style="fill:none;stroke:#fefefe;stroke-width:1.20701277;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1"
           d="m 679.8289,256.29253 1.5372,2.66238 3.0279,-5.1079"
           id="path17048"
           inkscape:connector-curvature="0" />
      </g>
      <text
         id="text16996"
         y="347.87659"
         x="243.90221"
         style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:4.23333311px;line-height:1.25;font-family:comfortaa;-inkscape-font-specification:comfortaa;text-align:center;letter-spacing:0px;word-spacing:0px;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke:none;stroke-width:0.31935549"
         xml:space="preserve"><tspan
           style="text-align:center;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.31935549"
           y="347.87659"
           x="243.90221"
           id="tspan16992"
           sodipodi:role="line">Mark</tspan><tspan
           id="tspan16994"
           style="text-align:center;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.31935549"
           y="353.16824"
           x="243.90221"
           sodipodi:role="line" /></text>
      <rect
         ry="0.64644408"
         y="332.40796"
         x="244.92342"
         height="8.75"
         width="8.75"
         id="rect17020"
         style="opacity:1;fill:#ffffff;fill-opacity:1;stroke:#ed1a1f;stroke-width:0.90525949;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
    </g>
    <g
       id="g21100"
       transform="translate(-265.30594,214.1246)">
      <rect
         style="opacity:1;fill:#fdcccb;fill-opacity:1;stroke:#ed1a1f;stroke-width:0.74999994;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1"
         id="rect21035"
         width="11.75"
         height="24.03001"
         x="270.11682"
         y="37.904831"
         ry="1.6216873" />
      <g
         style="stroke-width:0.99797541"
         transform="matrix(1.0040615,0,0,1,13.261788,41.096524)"
         id="g21065">
        <path
           style="fill:#ed1a1f;fill-opacity:1;stroke:none;stroke-width:0.26404768px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1"
           d="m 257.489,8.5897632 h 8.34805 l 0.009,-7.66175497 -2.27922,-2.33824533 h -6.07834 z"
           id="path21047"
           inkscape:connector-curvature="0"
           sodipodi:nodetypes="cccccc" />
        <path
           style="fill:#ffffff;stroke:none;stroke-width:0.26404768px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1"
           d="m 265.41413,1.0915791 h -2.07938 v -2.07786084 z"
           id="path21049"
           inkscape:connector-curvature="0" />
        <path
           style="fill:none;stroke:#f7f7f7;stroke-width:0.4989877;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1"
           d="m 259.07487,2.3060688 1.26751,1.2675083"
           id="path21051"
           inkscape:connector-curvature="0" />
        <g
           style="stroke-width:0.56135339;stroke-miterlimit:4;stroke-dasharray:none"
           id="g21055"
           transform="matrix(0.88890122,0,0,0.88890122,41.288596,34.147288)">
          <path
             style="fill:none;stroke:#f7f7f7;stroke-width:0.56135339;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1"
             d="m 246.4321,-35.820876 -1.42593,1.425927"
             id="path21053"
             inkscape:connector-curvature="0" />
        </g>
        <path
           inkscape:connector-curvature="0"
           id="path21057"
           d="m 262.36753,2.3060688 1.26751,1.2675083"
           style="fill:none;stroke:#f7f7f7;stroke-width:0.4989877;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
        <g
           style="stroke-width:0.56135339;stroke-miterlimit:4;stroke-dasharray:none"
           id="g21061"
           transform="matrix(0.88890122,0,0,0.88890122,44.581235,34.147288)">
          <path
             inkscape:connector-curvature="0"
             id="path21059"
             d="m 246.4321,-35.820876 -1.42593,1.425927"
             style="fill:none;stroke:#f7f7f7;stroke-width:0.56135339;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
        </g>
        <path
           style="fill:none;stroke:#ffffff;stroke-width:0.4989877;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1"
           d="m 258.91687,6.9695931 c 3.9614,-2.2871135 4.6757,0 4.6757,0"
           id="path21063"
           inkscape:connector-curvature="0"
           sodipodi:nodetypes="cc" />
      </g>
      <rect
         ry="0.64644408"
         y="51.374069"
         x="271.59906"
         height="8.75"
         width="8.7855377"
         id="rect21067"
         style="opacity:1;fill:#ffffff;fill-opacity:1;stroke:#ed1a1f;stroke-width:0.74999988;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
    </g>
    <rect
       style="opacity:1;fill:#fdcccb;fill-opacity:1;stroke:#ed1a1f;stroke-width:0.74999994;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1"
       id="rect51322"
       width="40.420525"
       height="8.6134062"
       x="7.6522551"
       y="-292.15524"
       ry="0.58128363"
       transform="scale(1,-1)" />
    <text
       id="text51342"
       y="291.77893"
       x="14.460031"
       style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:3.50728083px;line-height:1.25;font-family:comfortaa;-inkscape-font-specification:comfortaa;text-align:center;letter-spacing:0px;word-spacing:0px;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke:none;stroke-width:0.26458332"
       xml:space="preserve"
       transform="scale(1.0075036,0.99255229)"><tspan
         style="text-align:center;text-anchor:middle;fill:#ff0000;fill-opacity:1;stroke-width:0.26458332"
         y="291.77893"
         x="14.460031"
         id="tspan51338"
         sodipodi:role="line">Initials</tspan><tspan
         id="tspan51340"
         style="text-align:center;text-anchor:middle;fill:#ff0000;fill-opacity:1;stroke-width:0.26458332"
         y="296.16302"
         x="14.460031"
         sodipodi:role="line" /></text>
    <rect
       ry="0.46247476"
       y="284.90372"
       x="22.102682"
       height="6.2598677"
       width="24.54258"
       id="rect51344"
       style="opacity:1;fill:#ffffff;fill-opacity:1;stroke:#ed1a1f;stroke-width:0.74999988;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
  </g>
  <g
     inkscape:groupmode="layer"
     id="layer9"
     inkscape:label="anchors"
     style="display:inline"
     sodipodi:insensitive="true">
    <path
       transform="translate(0,-247)"
       style="opacity:1;fill:#00c041;fill-opacity:0.55399062;stroke:none;stroke-width:2;stroke-miterlimit:4;stroke-dasharray:none;stroke-dashoffset:9.44881821;stroke-opacity:1"
       id="path50406"
       sodipodi:type="arc"
       sodipodi:cx="0"
       sodipodi:cy="247"
       sodipodi:rx="4.8191962"
       sodipodi:ry="4.8191962"
       sodipodi:start="1.586911"
       sodipodi:end="1.582216"
       sodipodi:open="true"
       d="m -0.07765641,251.81857 a 4.8191962,4.8191962 0 0 1 -4.74100189,-4.89057 4.8191962,4.8191962 0 0 1 4.88500291,-4.74674 4.8191962,4.8191962 0 0 1 4.75246949,4.87943 4.8191962,4.8191962 0 0 1 -4.87384655,4.75819">
      <title
         id="title50411">ref-anchor</title>
    </path>
  </g>
  <g
     inkscape:groupmode="layer"
     id="layer8"
     inkscape:label="textfields"
     style="display:inline">
    <rect
       style="display:inline;opacity:1;fill:#e2e4ff;fill-opacity:0.75117373;stroke:none;stroke-width:0.26499999;stroke-miterlimit:4;stroke-dasharray:none;stroke-dashoffset:9.44881821;stroke-opacity:1"
       id="rect47789"
       width="7.8166327"
       height="8.0839024"
       x="6.810286"
       y="18.934586">
      <title
         id="title47791">badfile</title>
    </rect>
    <rect
       y="13.149623"
       x="38.650368"
       height="6.3086619"
       width="6.3086619"
       id="rect50379"
       style="display:inline;opacity:1;fill:#e2e4ff;fill-opacity:0.75117373;stroke:none;stroke-width:0.26499999;stroke-miterlimit:4;stroke-dasharray:none;stroke-dashoffset:9.44881821;stroke-opacity:1">
      <title
         id="title50377">markok</title>
    </rect>
    <rect
       style="display:inline;opacity:1;fill:#e2e4ff;fill-opacity:0.75117373;stroke:none;stroke-width:0.26499999;stroke-miterlimit:4;stroke-dasharray:none;stroke-dashoffset:9.44881821;stroke-opacity:1"
       id="rect50385"
       width="23.535746"
       height="5.3614168"
       x="22.629934"
       y="38.408676">
      <desc
         id="desc51388">Enter your initials here</desc>
      <title
         id="title50383">initials</title>
    </rect>
  </g>
</svg>`

func TestParseSvg(t *testing.T) {

	var svg Csvg__svg

	err := xml.Unmarshal([]byte(testInkscapeSvg), &svg)

	if err != nil {
		t.Errorf(err.Error())
	}

	jsonData, _ := json.Marshal(svg)

	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, jsonData, "", "\t")
	if error != nil {
		t.Errorf("JSON parse error: %v", error)
	}
}

var expectedLadder = &Ladder{
	Anchor: geo.Point{X: 0, Y: 0},
	Dim:    geo.Dim{Width: 141.73228346456693, Height: 141.73228346456693},
	ID:     "example",
	TextFields: []TextField{
		TextField{
			Rect: geo.Rect{
				Corner: geo.Point{X: 19.304747716535434, Y: 53.67284220472441},
				Dim:    geo.Dim{Width: 22.157384031496065, Height: 22.91499892913386},
			},
			ID: "badfile",
		},
		TextField{
			Rect: geo.Rect{
				Corner: geo.Point{X: 109.56009826771654, Y: 37.27452188976378},
				Dim:    geo.Dim{Width: 17.88282113385827, Height: 17.88282113385827},
			},
			ID: "markok",
		},
		TextField{
			Rect: geo.Rect{
				Corner: geo.Point{X: 64.14784440944882, Y: 108.87498708661418},
				Dim:    geo.Dim{Width: 66.71550047244095, Height: 15.197716913385827},
			},
			ID:      "initials",
			Prefill: "Enter your initials here",
		},
	},
}

/* Output from 	fmt.Printf("%v", ladder):
&{{0 0} {141.73228346456693 141.73228346456693}  [{{{19.304747716535434 88.05944125984252} {22.157384031496065 22.91499892913386}} badfile } {{{109.56009826771654 104.45776157480316} {17.88282113385827 17.88282113385827}} markok } {{{64.14784440944882 32.85729637795275} {66.71550047244095 15.197716913385827}} initials Enter your intials here}]}
Used to constructed expected result after close visual insection of the output.
*/

const textPrefillSVG = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!-- Created with Inkscape (http://www.inkscape.org/) -->

<svg
   xmlns:dc="http://purl.org/dc/elements/1.1/"
   xmlns:cc="http://creativecommons.org/ns#"
   xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
   xmlns:svg="http://www.w3.org/2000/svg"
   xmlns="http://www.w3.org/2000/svg"
   xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
   xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
   width="50mm"
   height="50mm"
   viewBox="0 0 50 50"
   version="1.1"
   id="svg8"
   inkscape:version="0.92.4 (5da689c313, 2019-01-14)"
   sodipodi:docname="textprefill.svg">
  <defs
     id="defs2" />
  <sodipodi:namedview
     id="base"
     pagecolor="#ffffff"
     bordercolor="#666666"
     borderopacity="1.0"
     inkscape:pageopacity="0.0"
     inkscape:pageshadow="2"
     inkscape:zoom="1.979899"
     inkscape:cx="87.774531"
     inkscape:cy="144.00533"
     inkscape:document-units="mm"
     inkscape:current-layer="layer2"
     showgrid="false"
     inkscape:snap-object-midpoints="true"
     inkscape:snap-center="false"
     inkscape:snap-page="true"
     inkscape:window-width="1850"
     inkscape:window-height="1136"
     inkscape:window-x="70"
     inkscape:window-y="27"
     inkscape:window-maximized="1" />
  <metadata
     id="metadata5">
    <rdf:RDF>
      <cc:Work
         rdf:about="">
        <dc:format>image/svg+xml</dc:format>
        <dc:type
           rdf:resource="http://purl.org/dc/dcmitype/StillImage" />
        <dc:title />
      </cc:Work>
    </rdf:RDF>
  </metadata>
  <g
     inkscape:label="anchors"
     inkscape:groupmode="layer"
     id="layer1"
     transform="translate(0,-247)"
     style="display:inline">
    <path
       style="opacity:0.98999999;fill:#3dff74;fill-opacity:1;stroke:none;stroke-width:0.2;stroke-miterlimit:4;stroke-dasharray:none;stroke-dashoffset:9.44881821;stroke-opacity:1"
       id="path817"
       sodipodi:type="arc"
       sodipodi:cx="0"
       sodipodi:cy="247"
       sodipodi:rx="2.7395127"
       sodipodi:ry="2.7395127"
       sodipodi:start="1.5935991"
       sodipodi:end="1.5918851"
       sodipodi:open="true"
       d="m -0.06246307,249.7388 a 2.7395127,2.7395127 0 0 1 -2.67636393,-2.80009 2.7395127,2.7395127 0 0 1 2.7989429,-2.67756 2.7395127,2.7395127 0 0 1 2.6787626,2.79779 2.7395127,2.7395127 0 0 1 -2.79664718,2.67996" />
  </g>
  <g
     inkscape:groupmode="layer"
     id="layer2"
     inkscape:label="textprefills"
     style="display:inline">
    <rect
       transform="translate(0,-247)"
       style="display:inline;opacity:0.98999999;fill:#0000ff;fill-opacity:1;stroke:none;stroke-width:0.2;stroke-miterlimit:4;stroke-dasharray:none;stroke-dashoffset:9.44881821;stroke-opacity:1"
       id="rect819"
       width="30.067822"
       height="10.022608"
       x="9.0871639"
       y="254.77141"
       ry="0">
      <desc
         id="desc824">{&quot;text&quot;:&quot;someContents&quot;,&quot;textFont&quot;:&quot;Helvetica&quot;,&quot;textSize&quot;:10,&quot;lineHeight&quot;:1,&quot;alignment&quot;:&quot;left&quot;,&quot;enableWrap&quot;:true,&quot;wrapWidth&quot;:50,&quot;angle&quot;:0,&quot;absolutePositioning&quot;:false,&quot;margins&quot;:null,&quot;xpos&quot;:50,&quot;ypos&quot;:50,&quot;colorHex&quot;:&quot;#ffee33&quot;}
</desc>
      <title
         id="title822">topbox</title>
    </rect>
  </g>
  <g
     inkscape:groupmode="layer"
     id="layer3"
     inkscape:label="textfields"
     style="display:inline">
    <rect
       style="opacity:0.98999999;fill:#00ffff;fill-opacity:1;stroke:none;stroke-width:0.2;stroke-miterlimit:4;stroke-dasharray:none;stroke-dashoffset:9.44881821;stroke-opacity:1"
       id="rect846"
       width="17.372519"
       height="14.432554"
       x="16.704346"
       y="28.083899"
       ry="0">
      <desc
         id="desc850">This is the prefill text for the bottom box, a textfield i.e acroforms, and not a uneditable prefill like the other box</desc>
      <title
         id="title848">bottombox</title>
    </rect>
  </g>
</svg>`

var expectedTextPrefill = &Ladder{
	Anchor: geo.Point{X: 0, Y: 0},
	Dim:    geo.Dim{Width: 141.73228346456693, Height: 141.73228346456693},
	ID:     "",
	TextFields: []TextField{
		TextField{
			Rect: geo.Rect{
				Corner: geo.Point{X: 47.350902047244105, Y: 79.60790267716536},
				Dim:    geo.Dim{Width: 49.2449357480315, Height: 40.911176692913386},
			},
			ID:      "bottombox",
			Prefill: "This is the prefill text for the bottom box, a textfield i.e acroforms, and not a uneditable prefill like the other box",
		},
	},
	TextPrefills: []TextPrefill{
		TextPrefill{
			Rect: geo.Rect{
				Corner: geo.Point{X: 25.758889795275593, Y: 22.029193700787413},
				Dim:    geo.Dim{Width: 85.23162141732284, Height: 28.410542362204726},
			},
			ID:         "topbox",
			Properties: "{\"text\":\"someContents\",\"textFont\":\"Helvetica\",\"textSize\":10,\"lineHeight\":1,\"alignment\":\"left\",\"enableWrap\":true,\"wrapWidth\":50,\"angle\":0,\"absolutePositioning\":false,\"margins\":null,\"xpos\":50,\"ypos\":50,\"colorHex\":\"#ffee33\"}\n",
			Text: Paragraph{
				Text:                "someContents",
				TextFont:            "Helvetica",
				TextSize:            10,
				LineHeight:          1,
				Alignment:           "left",
				EnableWrap:          true,
				WrapWidth:           50,
				Angle:               0,
				AbsolutePositioning: false,
				Margins:             nil,
				XPos:                50,
				YPos:                50,
				ColorHex:            "#ffee33",
			},
		},
	},
}

func TestDefineLadderFromSvg(t *testing.T) {

	ladder, err := DefineLadderFromSVG([]byte(testInkscapeSvg))
	if err != nil {
		t.Errorf("Error defining ladder %v", err)
	}

	if !reflect.DeepEqual(ladder, expectedLadder) {
		t.Errorf("Ladder does not match expected")
		fmt.Println("----------EXPECTED------------------")
		PrettyPrintStruct(expectedLadder)
		fmt.Println("----------GOT------------------")
		PrettyPrintStruct(ladder)
	}

}

func TestTextPrefills(t *testing.T) {

	ladder, err := DefineLadderFromSVG([]byte(textPrefillSVG))
	if err != nil {
		t.Errorf("Error defining ladder %v", err)
	}

	if !reflect.DeepEqual(ladder, expectedTextPrefill) {
		t.Errorf("Ladder does not match expected")
		t.Errorf("Ladder does not match expected")
		fmt.Println("----------EXPECTED------------------")
		PrettyPrintStruct(expectedTextPrefill)
		fmt.Println("----------GOT------------------")
		PrettyPrintStruct(ladder)
	}
}

func TestPrintParsedExample(t *testing.T) {

	c := creator.New()

	c.SetPageMargins(0, 0, 0, 0) // we're not printing

	jpegFilename := "./test/example-chrome.jpg"
	pageFilename := "./test/example-with-textfields.pdf"
	img, err := c.NewImageFromFile(jpegFilename)

	if err != nil {
		t.Errorf("Error opening image file: %s", err)
	}
	writeParsedGeometry([]byte(testInkscapeSvg), img, pageFilename, c, t)
}

func TestPrintParsedLargeExample(t *testing.T) {

	c := creator.New()

	c.SetPageMargins(0, 0, 0, 0) // we're not printing

	svgFilename := "./test/ladders-a4-portrait-mark.svg"
	jpegFilename := "./test/ladders-a4-portrait-mark.jpg"
	pageFilename := "./test/ladders-a4-portrait-mark.pdf"

	svgBytes, err := ioutil.ReadFile(svgFilename)
	if err != nil {
		t.Error(err)
	}

	img, err := c.NewImageFromFile(jpegFilename)

	if err != nil {
		t.Errorf("Error opening image file: %s", err)
	}

	writeParsedGeometry(svgBytes, img, pageFilename, c, t)
}

func writeParsedLayout(svg []byte, img *creator.Image, pageFilename string, c *creator.Creator, t *testing.T) {

	ladder, err := DefineLadderFromSVG(svg)
	if err != nil {
		t.Errorf("Error defining ladder %v", err)
	}

	// scale and position image
	img.ScaleToHeight(ladder.Dim.Height)
	img.SetPos(ladder.Anchor.X, ladder.Anchor.Y) //TODO check this has correct sense for non-zero offsets

	// create new page with image
	c.SetPageSize(creator.PageSize{ladder.Dim.Width, ladder.Dim.Height})
	c.NewPage()
	c.Draw(img)

	// write to memory
	var buf bytes.Buffer

	err = c.Write(&buf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := model.NewPdfReader(fbuf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	pdfWriter := model.NewPdfWriter()

	page, err := pdfReader.GetPage(1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	form := model.NewPdfAcroForm()

	for _, tf := range ladder.TextFields {

		tfopt := annotator.TextFieldOptions{Value: tf.Prefill} //TODO - MaxLen?!
		name := fmt.Sprintf("Page-00-%s", tf.ID)
		textf, err := annotator.NewTextField(page, name, formRect(tf, ladder.Dim), tfopt)
		if err != nil {
			panic(err)
		}
		*form.Fields = append(*form.Fields, textf.PdfField)
		page.AddAnnotation(textf.Annotations[0].PdfAnnotation)
	}

	err = pdfWriter.SetForms(form)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = pdfWriter.AddPage(page)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	of, err := os.Create(pageFilename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer of.Close()

	pdfWriter.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    80,
		ImageUpperPPI:                   100,
	}))

	pdfWriter.Write(of)
}

func writeParsedGeometry(svg []byte, img *creator.Image, pageFilename string, c *creator.Creator, t *testing.T) {

	ladder, err := DefineLadderFromSVG(svg)
	if err != nil {
		t.Errorf("Error defining ladder %v", err)
	}

	// scale and position image
	img.ScaleToHeight(ladder.Dim.Height)
	img.SetPos(ladder.Anchor.X, ladder.Anchor.Y) //TODO check this has correct sense for non-zero offsets

	// create new page with image
	c.SetPageSize(creator.PageSize{ladder.Dim.Width, ladder.Dim.Height})
	c.NewPage()
	c.Draw(img)

	// write to memory
	var buf bytes.Buffer

	err = c.Write(&buf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := model.NewPdfReader(fbuf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	pdfWriter := model.NewPdfWriter()

	page, err := pdfReader.GetPage(1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	form := model.NewPdfAcroForm()

	for _, tf := range ladder.TextFields {

		tfopt := annotator.TextFieldOptions{Value: tf.Prefill} //TODO - MaxLen?!
		name := fmt.Sprintf("Page-00-%s", tf.ID)
		textf, err := annotator.NewTextField(page, name, formRect(tf, ladder.Dim), tfopt)
		if err != nil {
			panic(err)
		}
		*form.Fields = append(*form.Fields, textf.PdfField)
		page.AddAnnotation(textf.Annotations[0].PdfAnnotation)
	}

	err = pdfWriter.SetForms(form)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = pdfWriter.AddPage(page)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	of, err := os.Create(pageFilename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer of.Close()

	pdfWriter.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    80,
		ImageUpperPPI:                   100,
	}))

	pdfWriter.Write(of)
}
