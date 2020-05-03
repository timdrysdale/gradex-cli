package parsesvg

import "encoding/xml"

type Csvg__svg struct {
	XMLName                               xml.Name              `xml:"svg,omitempty" json:"svg,omitempty"`
	AttrXmlnscc                           string                `xml:"xmlns cc,attr"  json:",omitempty"`
	AttrXmlnsdc                           string                `xml:"xmlns dc,attr"  json:",omitempty"`
	AttrSodipodiSpacedocname              string                `xml:"http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd docname,attr"  json:",omitempty"`
	AttrInkscapeSpaceexport_dash_filename string                `xml:"http://www.inkscape.org/namespaces/inkscape export-filename,attr"  json:",omitempty"`
	AttrInkscapeSpaceexport_dash_xdpi     string                `xml:"http://www.inkscape.org/namespaces/inkscape export-xdpi,attr"  json:",omitempty"`
	AttrInkscapeSpaceexport_dash_ydpi     string                `xml:"http://www.inkscape.org/namespaces/inkscape export-ydpi,attr"  json:",omitempty"`
	Height                                string                `xml:"height,attr"  json:",omitempty"`
	Attrid                                string                `xml:"id,attr"  json:",omitempty"`
	AttrXmlnsinkscape                     string                `xml:"xmlns inkscape,attr"  json:",omitempty"`
	AttrXmlnsrdf                          string                `xml:"xmlns rdf,attr"  json:",omitempty"`
	AttrXmlnssodipodi                     string                `xml:"xmlns sodipodi,attr"  json:",omitempty"`
	AttrXmlnssvg                          string                `xml:"xmlns svg,attr"  json:",omitempty"`
	Attrversion                           string                `xml:"version,attr"  json:",omitempty"`
	AttrviewBox                           string                `xml:"viewBox,attr"  json:",omitempty"`
	Width                                 string                `xml:"width,attr"  json:",omitempty"`
	Attrxmlns                             string                `xml:"xmlns,attr"  json:",omitempty"`
	Cdefs__svg                            *Cdefs__svg           `xml:"http://www.w3.org/2000/svg defs,omitempty" json:"defs,omitempty"`
	Cg__svg                               []*Cg__svg            `xml:"http://www.w3.org/2000/svg g,omitempty" json:"g,omitempty"`
	Cmetadata__svg                        *Cmetadata__svg       `xml:"http://www.w3.org/2000/svg metadata,omitempty" json:"metadata,omitempty"`
	Cnamedview__sodipodi                  *Cnamedview__sodipodi `xml:"http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd namedview,omitempty" json:"namedview,omitempty"`
	Ctitle__svg                           *Ctitle__svg          `xml:"http://www.w3.org/2000/svg title,omitempty" json:"title,omitempty"`
}

type Ctitle__svg struct {
	XMLName xml.Name `xml:"title,omitempty" json:"title,omitempty"`
	Attrid  string   `xml:"id,attr"  json:",omitempty"`
	String  string   `xml:",chardata" json:",omitempty"`
}

type Cdefs__svg struct {
	XMLName xml.Name `xml:"defs,omitempty" json:"defs,omitempty"`
	Attrid  string   `xml:"id,attr"  json:",omitempty"`
}

type Cnamedview__sodipodi struct {
	XMLName                                xml.Name `xml:"namedview,omitempty" json:"namedview,omitempty"`
	Attrbordercolor                        string   `xml:"bordercolor,attr"  json:",omitempty"`
	Attrborderopacity                      string   `xml:"borderopacity,attr"  json:",omitempty"`
	AttrInkscapeSpacecurrent_dash_layer    string   `xml:"http://www.inkscape.org/namespaces/inkscape current-layer,attr"  json:",omitempty"`
	AttrInkscapeSpacecx                    string   `xml:"http://www.inkscape.org/namespaces/inkscape cx,attr"  json:",omitempty"`
	AttrInkscapeSpacecy                    string   `xml:"http://www.inkscape.org/namespaces/inkscape cy,attr"  json:",omitempty"`
	AttrInkscapeSpacedocument_dash_units   string   `xml:"http://www.inkscape.org/namespaces/inkscape document-units,attr"  json:",omitempty"`
	Attrid                                 string   `xml:"id,attr"  json:",omitempty"`
	Attrpagecolor                          string   `xml:"pagecolor,attr"  json:",omitempty"`
	AttrInkscapeSpacepageopacity           string   `xml:"http://www.inkscape.org/namespaces/inkscape pageopacity,attr"  json:",omitempty"`
	AttrInkscapeSpacepageshadow            string   `xml:"http://www.inkscape.org/namespaces/inkscape pageshadow,attr"  json:",omitempty"`
	Attrshowgrid                           string   `xml:"showgrid,attr"  json:",omitempty"`
	AttrInkscapeSpacesnap_dash_center      string   `xml:"http://www.inkscape.org/namespaces/inkscape snap-center,attr"  json:",omitempty"`
	AttrInkscapeSpacesnap_dash_global      string   `xml:"http://www.inkscape.org/namespaces/inkscape snap-global,attr"  json:",omitempty"`
	AttrInkscapeSpacesnap_dash_page        string   `xml:"http://www.inkscape.org/namespaces/inkscape snap-page,attr"  json:",omitempty"`
	AttrInkscapeSpacewindow_dash_height    string   `xml:"http://www.inkscape.org/namespaces/inkscape window-height,attr"  json:",omitempty"`
	AttrInkscapeSpacewindow_dash_maximized string   `xml:"http://www.inkscape.org/namespaces/inkscape window-maximized,attr"  json:",omitempty"`
	AttrInkscapeSpacewindow_dash_width     string   `xml:"http://www.inkscape.org/namespaces/inkscape window-width,attr"  json:",omitempty"`
	AttrInkscapeSpacewindow_dash_x         string   `xml:"http://www.inkscape.org/namespaces/inkscape window-x,attr"  json:",omitempty"`
	AttrInkscapeSpacewindow_dash_y         string   `xml:"http://www.inkscape.org/namespaces/inkscape window-y,attr"  json:",omitempty"`
	AttrInkscapeSpacezoom                  string   `xml:"http://www.inkscape.org/namespaces/inkscape zoom,attr"  json:",omitempty"`
}

type Cmetadata__svg struct {
	XMLName   xml.Name   `xml:"metadata,omitempty" json:"metadata,omitempty"`
	Attrid    string     `xml:"id,attr"  json:",omitempty"`
	CRDF__rdf *CRDF__rdf `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# RDF,omitempty" json:"RDF,omitempty"`
}

type CRDF__rdf struct {
	XMLName   xml.Name   `xml:"RDF,omitempty" json:"RDF,omitempty"`
	CWork__cc *CWork__cc `xml:"http://creativecommons.org/ns# Work,omitempty" json:"Work,omitempty"`
}

type CWork__cc struct {
	XMLName           xml.Name     `xml:"Work,omitempty" json:"Work,omitempty"`
	AttrRdfSpaceabout string       `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# about,attr"  json:",omitempty"`
	Cformat__dc       *Cformat__dc `xml:"http://purl.org/dc/elements/1.1/ format,omitempty" json:"format,omitempty"`
	Ctitle__dc        *Ctitle__dc  `xml:"http://purl.org/dc/elements/1.1/ title,omitempty" json:"title,omitempty"`
	Ctype__dc         *Ctype__dc   `xml:"http://purl.org/dc/elements/1.1/ type,omitempty" json:"type,omitempty"`
}

type Cformat__dc struct {
	XMLName xml.Name `xml:"format,omitempty" json:"format,omitempty"`
	String  string   `xml:",chardata" json:",omitempty"`
}

type Ctype__dc struct {
	XMLName              xml.Name `xml:"type,omitempty" json:"type,omitempty"`
	AttrRdfSpaceresource string   `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# resource,attr"  json:",omitempty"`
}

type Ctitle__dc struct {
	XMLName xml.Name `xml:"title,omitempty" json:"title,omitempty"`
	String  string   `xml:",chardata" json:",omitempty"`
}

type Cg__svg struct {
	XMLName                    xml.Name      `xml:"g,omitempty" json:"g,omitempty"`
	AttrInkscapeSpacegroupmode string        `xml:"http://www.inkscape.org/namespaces/inkscape groupmode,attr"  json:",omitempty"`
	Attrid                     string        `xml:"id,attr"  json:",omitempty"`
	AttrInkscapeSpacelabel     string        `xml:"http://www.inkscape.org/namespaces/inkscape label,attr"  json:",omitempty"`
	Attrstyle                  string        `xml:"style,attr"  json:",omitempty"`
	Cpath__svg                 []*Cpath__svg `xml:"http://www.w3.org/2000/svg path,omitempty" json:"path,omitempty"`
	Crect__svg                 []*Crect__svg `xml:"http://www.w3.org/2000/svg rect,omitempty" json:"rect,omitempty"`
	Transform                  string        `xml:"transform,attr"  json:",omitempty"`
}

type Cpath__svg struct {
	XMLName                xml.Name     `xml:"path,omitempty" json:"path,omitempty"`
	Cx                     string       `xml:"http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd cx,attr"  json:",omitempty"`
	Cy                     string       `xml:"http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd cy,attr"  json:",omitempty"`
	Attrd                  string       `xml:"d,attr"  json:",omitempty"`
	AttrSodipodiSpaceend   string       `xml:"http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd end,attr"  json:",omitempty"`
	ID                     string       `xml:"id,attr"  json:",omitempty"`
	AttrSodipodiSpaceopen  string       `xml:"http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd open,attr"  json:",omitempty"`
	AttrSodipodiSpacerx    string       `xml:"http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd rx,attr"  json:",omitempty"`
	AttrSodipodiSpacery    string       `xml:"http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd ry,attr"  json:",omitempty"`
	AttrSodipodiSpacestart string       `xml:"http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd start,attr"  json:",omitempty"`
	Attrstyle              string       `xml:"style,attr"  json:",omitempty"`
	Transform              string       `xml:"transform,attr"  json:",omitempty"`
	AttrSodipodiSpacetype  string       `xml:"http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd type,attr"  json:",omitempty"`
	Desc                   *Cdesc__svg  `xml:"http://www.w3.org/2000/svg desc,omitempty" json:"desc,omitempty"`
	Title                  *Ctitle__svg `xml:"http://www.w3.org/2000/svg title,omitempty" json:"title,omitempty"`
}

type Crect__svg struct {
	XMLName   xml.Name     `xml:"rect,omitempty" json:"rect,omitempty"`
	Height    string       `xml:"height,attr"  json:",omitempty"`
	Id        string       `xml:"id,attr"  json:",omitempty"`
	Attrstyle string       `xml:"style,attr"  json:",omitempty"`
	Width     string       `xml:"width,attr"  json:",omitempty"`
	Rx        string       `xml:"x,attr"  json:",omitempty"`
	Ry        string       `xml:"y,attr"  json:",omitempty"`
	Desc      *Cdesc__svg  `xml:"http://www.w3.org/2000/svg desc,omitempty" json:"desc,omitempty"`
	Title     *Ctitle__svg `xml:"http://www.w3.org/2000/svg title,omitempty" json:"title,omitempty"`
	Transform string       `xml:"transform,attr"  json:",omitempty"`
}

type Cdesc__svg struct {
	XMLName xml.Name `xml:"desc,omitempty" json:"desc,omitempty"`
	Attrid  string   `xml:"id,attr"  json:",omitempty"`
	String  string   `xml:",chardata" json:",omitempty"`
}
