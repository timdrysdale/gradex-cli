# pdfpagedata
read/write gradex pagedata frompdf pages

## Why
Tracking pages in PDF documents when they are split into separate files can be done a couple of ways

 -- postpend something to the filename for page
 -- put info in metadata
 -- use some page text, off page
 -- protocol buf into stream object

Unfortunately, processing thousands of pages, though different hands, needs more safety than fragile filenames can provide. Plus we will have duplicate files with non-duplicate annotations. So what do we do when two people re-upload their different files with the same name, and then either overwrite or make some modification to the filename? How do we recover from an unfortunate name choice here?

We could stash info in the metadata, but that tends to be file-level, so it is not clear how to handle duplicate custom metadata fields when multiple files from different documents are joined, then split, then joined again etc. Bseides, I've seen editors mess with the metadata, and I don't fancy users editing it either.

Off-page page text seems fragile too, but I get some comfort from reading that people [an NOT crop when they want to](https://community.adobe.com/t5/acrobat/discarding-cropped-areas-of-pages/td-p/4304473?page=1). A test has been included for this very purpose -which is passing.

## Wrinkles

text written in the same place gets read back out in some sort of merged way, so pageData is written in a tiny font (like 0.00001) and randomly scattered around a location that is far off the page. Tag destruction is detected (such as for clases), and multiple page datas on a page are supported.

A collision _is_ possible ... we could always consider writing each hidden data twice ...

## Future

Protocol buf into a stream object seems like a more robust way (and it avoids crop and collision worries) but it is probably about a half-day or a day to develop so that makes it a roadmap item for now.

