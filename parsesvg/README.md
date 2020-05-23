# parsesvg
parse svg to produce the acroforms needed by **pdf.gradex**™ (an AGPL open source exam marking tool)

![alt text][logo]


![alt text][status]


If you are coming here from pdf.gradex.io, then welcome to the developer side of things where all the details are. We're working on a user workflow that is a piece of cake, that will automate the fiddly bits. Meanwhile, if you like knowing how things work, feel free to take a look around.

## Installation

you'll need to have image magick available on the path, so that visual comparisons can be done with compare.

These are a bit flakey, so sometimes need to be run a couple of times for all the tests to pass.


## Context

The work of the support tool is combine script images and marking acroforms, at each stage of the process. The render tool take cares of adding the acroforms to the images. We let another tool put this module to use - our concern in this module is the layout of the pages we're producing (and understanding the SVG design that specifies them). If you want to make your own designs, then this is the module README that you need.... mostly these are notes to myself but we'll back fill in the gaps as we identify them.

![alt text][support-workflow]


## Why parse SVG, I just want to design custom PDF acroforms workflows for my Uni?

It's a lot easier to do geometry entry using a GUI, particularly when you want it to look nice. So why not leverage the fabulous [inkscape](https://inkscape.org/) for a double win? Pretty forms AND output in a format that we can parse to the find the exact coordinates of the acroforms we want to embed. Goodbye editing golang structs by hand in one window, with inkscape's ruler hard at work in another.

## How?
SVG is XML, so parsing for the position and size of an object is straightforward. There are some domain-specific constraints:

- inkscape has no idea we are doing this, so can't guide you on the drawing/annotating process 
- it's still _way_ easier than editing structs by hand, but it doesn't avoid having to think
- the output boxes always have pointy corners, so beware of your design-frustration sky-rocketing when you do nice rounded borders and find the light blue sharp corners of the TextField ruining your design vision. You can always hand craft some structs to relax again.
- conventions are a moving target ... you'll be naming a bunch of objects in inkscape, then re-doing it again later, just saying.
- there are transformations, such as translate, that we need to account for in calculating the position (but it seems that transforms are flattened on each application, so we should only need to do this once per object). There seems to be a global translate in all the svg I have looked at so far ... (hence the use of the reference ```anchors```)


```svg
  <g
     inkscape:label="Layer 1"
     inkscape:groupmode="layer"
     id="layer1"
     transform="translate(0,-247)">
    <rect
       style="opacity:1;fill:#ffffff;fill-opacity:0.75117373;stroke:#000000;stroke-width:0.26499999;stroke-miterlimit:4;stroke-dasharray:none;stroke-dashoffset:9.44881821;stroke-opacity:1"
       id="rect47789"
       width="24.735001"
       height="24.735001"
       x="14.363094"
       y="259.41382">
      <title
         id="title47791">markbox-00-00</title>
    </rect>
  </g>

```

## Procedure

This is bit finicky, but we'll make it make sense. I was on the receiving end of a procedure a bit like this (in terms of finickyness) when [creating custom content](https://github.com/timdrysdale/rf3a) for ```$BRAND_NAME``` model flying simulator - it was set up that you could copy exactly what they did, and worked if you had the same versions of the CAD software,  but nothing else. It was painful but I was grateful they left the door open for customisation. So here's the crack in the door .... 

### Set up a page to exactly the size your marks ladder, and name the ladder

Assuming you already have a design you are modifying, or some mock designs, you will know the page size you want. I find it helpful to set the page size and export the page to ```png``` when making images, and will usually end up resizing the page at least once during any drawing. So this step is not essential, now. ```Ctrl-Shift-D``` when you are ready. While you are there, put the ladder name in the ```Title``` field of the metadata. Resist the urge to hit either of the strangely tempting buttons in the metadata dialogue as this title data is unique to the ladder - it seems you needn't explicitly save in this dialogue. 

### Set up layers

We want at least three layers - your pretty design (the ```chrome```), the reference and position ```anchors``` at least one acroforms layer (one layer per type of form element).

- [```textfields```]
- [```dropdowns```]  --not implemented
- [```checkboxes```] --not implemented
- [```comboboxes```] --not implemented
- ```anchors```
- ```chrome```

For exporting your pretty chrome that goes around the textfields, just make textfields and anchors invisible. They're on top of the chrome so you can check alignment. Their names are expected exactly in this lower case format, and pluralised for ```textfields``` and ```anchors```. 

### Set up the reference anchor

So we can make sense of things at import, we're going to put an anchor in the top left corner of the ladder. We'll use this as an aid to position the ladder and textFields in a future step. Meanwhile, how do we make an anchor? We need an unambiguous reference that is optically distinguishable from a rectangular box, hence a circle. 
```
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="50"/>
</svg>
```
If we used another box, we'd be forever unsure if we had the wrong corner, and if we made it small enough to avoid that, we'd never see it. Whereas, the ```x,y``` coordinates of an ```svg``` circle are its centre. You can make the radius, colour of your anchor(s) anything you like. You just have to get the centre on the top left page corner as you perceive it. Snapping a circle onto a corner of a page avoids errors creeping into the rest of your relevant measuresments - the exact point of the anchor circle is used. Turn on the snap to page boundaries and snap to object centres, as well as the default snap settings. These are marked with a star:

![alt text][anchor]


### Anchor types

We're supporting an iterative work flow (multiple people grading things), so we need to anticipate that ladders are added at each step, and that we may wish to make one of them a different size at some point, without having to re-process all the other ladder configurations. So we are going to draw each ladder in its own file, so once we've got all the cells named, we can just leave it alone...

When combining multiple ladders into a workflow, we'll want to use the individual ladders as leaf cells, and arrange their respective positions on the final page by placing ```position anchors```. The ```reference anchor``` for a particular ladder is just mapped onto the ```position anchor```. We'll do some fu with naming schemes to sort this out. To add position anchors to a layout, name the anchor avoiding the reserved ```ref-anchor``` and in the metadata description file, place the base of the filename that of the ```svg```  and ```jpg``` versions of that sidebar or header. Note that BOTH files are needed (we get the forms elements from the ```svg``` and the image of the chrome we get Inkscape to render. (TODO: Yes, that means the sidebars and headers are not vector graphics, but they get rendered to images during any later steps in the process anyway. Although in future if ```svg`` rendering or ```pdf``` insertion becomes an option, it would help differentiate the active area from the previously-edited, because when active it would presumably be in vector format.) Inkscape does not produce ```jpg``` and the ```pdf``` library doesn't speak ```png``` so we simply use ImageMagick to convert
``` convert sidebar-312pt-mark-flow.png sidebar-31pt-mark-flow.jpg```



## Acroforms

Acroforms supports several types of field. I'm ignoring signature boxes for now because we can do [opticalcheckboxes](https://github.com/timdrysdale/opticalcheckbox) which play better with the idea of freely annotating anywhere. (TODO: So far support is only provided for textfields, but dropboxes are needed for the checking workflow)

- [```textfields```]
- [```dropdowns```] --not implemented
- [```checkboxes```] --not implemented
- [```comboboxes```] --not implemented

### Labelling and annotating

#### Ladders

In the document properties tab, ```Ctrl-Shift-D``` set the name of the layout element in the Title field of the metadata. Make sure it matches the filename, and the exported image.

![alt text][element-name]
![alt text][element-filename]

```svg
inkscape:version="0.92.4 (5da689c313, 2019-01-14)"
   sodipodi:docname="sidebar-312pt-check-flow.svg"
   inkscape:export-filename="/home/<snip>/parsesvg/test/sidebar-312pt-check-flow.png"
   inkscape:export-xdpi="299.86111"
   inkscape:export-ydpi="299.86111">
  <title
     id="title11203">sidebar-312pt-check-flow</title>
  <defs
     id="defs2">
	 ```
#### Anchors

- Give the reference anchor the title ```ref-anchor```. Behaviour is undefined if you add more than one reference anchor.
- Give any position anchors the title ```pos-anchor-<element_name>```, where <element_name> is meaningful to you. The ```svg``` for the element will be found in the file mentioned in the document metadata title, so the name of the pos-anchor does not have to match, i.e. you don't have to name your element's svg pos-anchor-whatever, you can just call it ```whatever.svg``` and put ```whatever``` in the document metadata as the document title.

Note that the previous images have a specific format that must be respected ....

The anchor must be name ```img-previous-<your-spread-name>```, while the box on the image layer must be called ```image-previous-<your-spread-name>```

```golang
	// Obtain the special "previous-image" which is flattened/rendered to image version of this page at the last step
	anchorName:= fmt.Sprintf("img-previous-%s",spread.Name)
	dimName := 	fmt.Sprintf("previous-%s",spread.Name)
	offset := DiffPosition(layout.Anchors[anchorName], layout.Anchor) 

	previousImage := ImageInsert{
		Filename:              previousImagePath,
		Corner:                offset,                            
		Dim:                   layout.ImageDims[dimName], 
	}
```

#### Textfields

Textfields don't necessarily need to be prefilled (but they can), whereas constrained-choice selections must be pre-populated. Let's do that in ```inkscape``` for an easy life. You can label and describe SVG elements in ```inkscape``` by ```Ctrl-Shift-O``` (remember to hit the 'Set' button - I kept forgetting first time out, so do check when you go back to an object that the data has persisted.) We'll use these to pass extra information into the parser, e.g. ```choiceBox``` options, or format strings that might help with hydrating the ```id``` to include page numbers etc. This bit is going to move rapidly ... so consider any implied API to be experimental and subject to change from minute to minute.

```
  <g
     inkscape:label="Layer 1"
     inkscape:groupmode="layer"
     id="layer1"
     transform="translate(0,-247)">
    <rect
       style="opacity:1;fill:#ffffff;fill-opacity:0.75117373;stroke:#000000;stroke-width:0.26499999;stroke-miterlimit:4;stroke-dasharray:none;stroke-dashoffset:9.44881821;stroke-opacity:1"
       id="rect47789"
       width="24.735001"
       height="24.735001"
       x="14.363094"
       y="259.41382">
      <desc
         id="desc50373">This is a description field I wonder if it makes it into the svg file .....</desc>
      <title
         id="title47791">markbox-00-00</title>
    </rect>
  </g>
 ```

Textfield trouble shooting - if your chrome is present on the page, then the anchor is good - because they use the same anchor. Check that
-- they are on the correct layer
-- they are NOT grouped 

### Tab order of acroforms elements

The order in which elements are written into the ```pdf``` determines the tab order as experienced by the user (which box you go to next when you hit tab). This strongly affects the ease of use of the workflow so it needs to be set logically (e.g. running from top to bottom) to avoid causing extra work to markers and checkers using keyboards. Inkscape does not offer a way to manipulate the order of elements in the ```xml```, e.g.  modifying the ID does not cause a reordering (for obvious efficiency reasons). Therefore, a sorting provision is included in the parser, that re-orders based on the tab number appended to the id as follows ...


![alt text][taborder]


### Exporting your sidebar/header element chrome for usage in a layout

The chrome layer is the only layer you should export.
- insert a white layer (or whatever suits your documents) in the background, or else the unfilled areas become black when converted to jpg
- turn off the visibility of the anchor and and forms layers, e.g. text fields.

What you want is this:

![alt text][layer-visibility]

A trouble-shooting graphic shows you the main things you might get wrong:

![alt text][export-troubleshooting]


## Example

The example we'll work with has three ```textfields```, that are floating in fairly random places, and at random sizes. A great example of where this parsing comes in handy - else lining it up by hand will be fiddly and tedious.


 All layers visible: ![alt text][example]

Only the chrome: ![alt text][example-chrome]

```anchor``` & ```textfields```: ![alt text][example-nonchrome]

We don't need to examine the chrome, because we'll import that in a separate step as a background image, using the position ```anchor``` from the overall layout (TODO: create section on that). The ```textfields``` and ```anchor``` are grouped into their own layers:

```svg
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
         id="desc51388">Enter your intials here</desc>
      <title
         id="title50383">initials</title>
    </rect>
  </g>
```

Now that we have an example, we can fire up the parser, and punt out a test pdf.

## Actual usage

Our actual marking sidebars will have many textfields. Easy! We're using a GUI, so you can copy and paste, making sure to

- append the tab order code to the ID (e.g. ```<original-id>-tab-012```
- add the name of the element in the Title field, so you can identify it in extracted data later
 
This example is part of the tests (exported here with textfields layer turned on, showing in a light blue):

![alt text][sidebar-example]

## Layouts

We've got one sidebar working. Great! What next? The exams are going to make visits to moderators, and checkers too. They have their own sidebars. We've got this far using the GUI - let's continue to use it to organise where the individual sidebars and headers go. Here we make a separate ```layout.svg``` which contains the layers

- textfields
- anchors
- [chrome] (optional but useful) 
- pages

The chrome is optional, but you'll want to include it for helping get your anchors in the right place. Delete the white backgrounds in your chrome, so you can see the pages. These will be used to set the paper size

For example:

![]alt text][layout-example]

This example represents a three stage process where all the incoming scans have been scaled to A4-portrait. Handling landscape is straightforward, by setting a flag that triggers the use of an alternate layout set tuned to landscape - that flagging process is not part of ```timdrysdale/parsesvg``` - search the ```pdf.gradex``` ecosystem for more details. In this case we are assuming 100% of papers are moderated, and that therefore we will always include the moderation sidebar. For many, standard moderation means only a subset of papers are moderated, and it would be helpful to give a visual indication whether a paper was selected for moderation or not, and not include the moderation sidebar in its 'active' form if a paper is not to be moderated. We leave aside for now issues of auto-grouping files to send to those who will work on them, as those matters rest outside this library (but we need to support bifurcating workflows with the most straightforward specification of the desired behaviour that we can arrange).

### Dynamic sidebar selection

If our moderation process did not take place on a paper, we'd rather slot in a different moderation bar, that clearly indicates this. For example, we could have an active, coloured, version, and a thinner, grey, inactive version, of which we choose one:

![alt text][compare-active-inactive]

Our static layout isn't going to be able to take advantage of the space saving in the inactive case. So we introduce the concept of static and dynamic pages, and allow the engine to work out what to do in each case. We don't make decisions in the parse on what to trigger - we just need to collect enough information from the user and the svg to support those decisions in terms of layout. That avoids the parser having to interact with the page layout engine - it just throws a complete set of specifications over the wall to the engine, and lets it handle the dynamic elements.

For a dynamic element, it is likely that we have opinion about the size of one of the dimensions. So we can represent the dynamic page as  narrowest possible rectangle. Long in the known direction, narrow in the dynamic direction. It seems you can't have a zero width rect in Inkscape, so we have set a single-digit threshold (TODO - report threshold chosen here). This has the added benefit of making a dynamic layout file look different to the static layout. As before we place the new sidebar next to the right hand edge of the page - except now it is zero-width, so the check sidebar is on the far left of the screen, even though we know it ends up on the right hand side of the screen. If there is only dynamic element in the layout, then this is relatively easy to understand. If there are multiple dynamic elements, then the chrome will get quite busy in that area, and it may be worth considering having two or more layout docs during editing, then hold your breath and paste them together into one superlayout that you load into the tool (remember don't duplicate the ```ref-anchor``` or any other meaningful element). Make any edits in one of the two docs, rather than trying to edit the combined doc. I've not tried this - just making a note for later in case things get busy with any of our mark flows.

![dynamic-layout]

Note that we have two anchors for the red mark sidebars - the ```ladder``` for subtotals and the ```flow``` for totals have been split to demonstrate multiple elements being added to the same page.

### Paper size

We still need to finish off specifying the paper size - the parser is not a layout engine, and the layout engine might become hard to test if it works out a page size from the included elements and automagically figures out aesthetically pleasing padding etc. So we require the user to specify this by adding a ```page-static-<yourpagename>``` adjacent to the near-zero-width ```page-dynamic-[width/height]-<yourpagename>```. So you can have either

- ```page-dynamic-width-<yourpagename>```, or
-```page-dynamic-height-<yourpagename>```,

but not both, and the key-word ```page-dynamic-height-width-<*>``` is NOt implemented.

In this scheme,the previous two pages are labelled as `page-static=<somepage>` and `page-static=<someotherpage>` because they are static. 

### Previous-image size

We also need to let the page layout engine know about how large to make the image of the previous stage of the process, using the "previous-image-<yourpagename>" ID. For the case of the first two processing stages (red, green), the image is a fixed size. We auto-scale to make the red image, then the green image is the right size as a knock on effect (if we draw it around the red page correctly).
For the dynamic pages, the input image is the thing that varies in size, so this takes a near-zero wide rectangle in the dynamic direction (judt duplicate and rename the dynamic page rect, and move to the ```images``` layers)

## Spreads

A ```spread``` is the subsection of the overall layout that we pass to the layout engine for the construction of the page. Making the spread object is a separate job to the parser ... but we put in a partial implementation to test the idea, and it worked, so here it stays (for now).


## A note on coordinates

PDF coordinates are (0,0) in the lower left corner [says this info](https://www.pdfscripting.com/public/PDF-Page-Coordinates.cfm) - same as inkscape
BUT! The media box coordinates are [with respect to the top left](https://www.leadtools.com/help/sdk/v20/dh/to/pdf-coordinate-system.html)

Thus Y_position_pdf = page_height - Y_position_inkscape

### anchor styling - must be a circular path. No stroke.


## What next?

There's a [commmand line tool](https://github.com/timdrysdale/gradex-overlay) under development now, and a GUI to follow.

### Troubleshooting notes

Missing boxes in Adobe Reader/Pro, but present in other Readers:
check for a conflict in the box title inkscape - this becomes the form ID. Duplicate fields are not editable in Adobe Reader (this may well be what the spec says - PDF is notorious for having ambiguities / variance in implementation). Worth having a release approval process that includes ALL the editors you are supporting ... (standard, right!)

No prefill text? You MUST include the textSize field in the description of the textprefill box in Inkscape. There is no "good" default text size - and someone who knows the default behaviour of undefined variables in golang might find it weird that an opinion is given here given that text prefilling is also useful for hidden information. Open to opinions on this as it had me chasing my tail for a day as it is!
```
{"text":"THIS IS A PREFILL TEXT AREA","textSize":20}
```

![alt text][textprefill-example]

[anchor]: ./img/inkscape-anchor-alignment.png "circle on corner of page and snap settings bar"
[compare-active-inactive]: ./img/compare-active-inactive-sidebar.png "green coloured active moderate bar and grey thin inactive moderate sidebar"
[dynamic-layout]: ./img/dynamic-layout-60pc.png "screen showing three side bars and the layers dialong"
[element-name]: ./img/element-name.png "name of the layout element entered into metadata"
[element-filename]: ./img/element-filename.png "name of the layout element used for saving image of the chrome"
[example]: ./img/example.png "example of three textfields with pretty surrounds in red"
[example-chrome]: ./img/example-chrome.png "just showing the pretty surrounds, not the anchor or textfields"
[example-nonchrome]: ./img/example-nonchrome.png "just showing the anchor or textfields"
[export-troubleshooting]: ./img/export-troubleshooting.png "examples with black background, anchor and textfields showing"
[layout-example]: ./img/layout-example-with-layer-dialogue.png "drawing with a header and three side bars surrounding an original scan of work to be marked"
[metadata-title]: ./img/metadata-title.png "inkscape metadata tab in document properties, showing title is ladder3-rect"
[layer-visibility]: ./img/layer-visibility-at-export.png "Layers dialog with the chrome layer set to visible, and all others invisible"
[logo]: ./img/gradexTMlogo2-50pc.png "gradex logo"
[sidebar-example]: ./img/sidebar-mark.png "marking sidebar with nearly 30 textfields"
[status]: https://img.shields.io/badge/build-passing-green "Build passing"
[support-workflow]: ./img/supportWorkflow.png "three stages in the workflow, from marking to moderaing, to checking"
[taborder]: ./img/taborder.png "object properties dialogue showing tab-02 appended to ID, setting tab order"
[textprefill-example]: ./img/text-prefills-test-output.png "example of text prefills being positioned on page"