# gradex-opticalcheckbox

optical check box for use in hand annotated electronic documents

## Why?

Check boxes in acroforms might sometimes not be activated correctly by pen display users, so being able to detect freehand electronic annotations is a useful alternative. The intended procedure is place an acroform in the region of the checkbox so that a keyboard user can type an X there, or a pen user can annotate a cross. There is no need to handle non-ideal backgrounds. Nearly all-white squares are true if Vanilla, or nearly all-black are true if !Vanilla. This approach can also be used to identify rectangles with hand annotations for selectively exporting annotations to a compact summary page

## Demo

The test image is eight squares, of which the second (all white) is ```false``` if Vanilla, and the third to last (all black) is false if !Vanilla. The second to last box, mostly white with a small dot, is also ```false``` if Vanilla. This small threshold allows a small touch with a pen to pass un-noticed, to avoid false positives when the mark is so small the person making it may not notice it.

![alt text][test]

[test]: ./img/test.png "test image comprising eight squares filled variously"