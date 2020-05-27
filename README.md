# gradex-cli

![alt text][logo]

Command line interface to gradex tool

![alt text][status]

## Why?

A certain pandemic pushed all our examinations online at short notice, which meant handling a lot of PDF. In a big higher education department, it's preferable to have straightforward, robust procedures. Some of my thinking behind why PDF on its own was not enough for our needs, is [here](https://www.teaching-matters-blog.ed.ac.uk/spotlight-on-alternative-assessment-methods-remote-exam-marking-holding-on-to-the-philosophy-of-paper/).

[![teaching-matters]](https://www.teaching-matters-blog.ed.ac.uk/spotlight-on-alternative-assessment-methods-remote-exam-marking-holding-on-to-the-philosophy-of-paper/)


Essentially, paper keeps a permanent record of any student work, and any marking, allowing it to be passed from student to marker to moderator to checker without anyone being able to inadvertently delete or lose work from the previous people in the chain. This does the same, except using PDF.

## Features

- Integrated workflow for preparing PDF for marking, moderating and checking.
- Headers for page listing exam, student anonymous ID, and page number.
- colourful sidebars for recording marking, moderating and checking.
- Flattens (i.e. keeps) all student and staff annotations
- Flattens (i.e. keeps) all sticky-note comments
- Uniquely tracks each page, and its processing history, with CRC32-protected data-tags
- Automatically sorts labelled pages by script or by question
- Knows where to store files that it receives at each stage
- Keeps exam exam separately
- Customisable templates
- Dynamically reconfigurable text
- Textfields and Comboboxes using acroforms
- Parallel processing for increased speed

## Installation

### Prerequisites

#### Required

Ghostscript downloads can be found [here](https://www.ghostscript.com/download.html).

For Windows, choose the 64bit version.

#### For testing

ImageMagick must be [installed](https://imagemagick.org/script/download.php), and on the path, so as to allow visual comparisons of rendered PDFs. 


#### Optional

Logging messages can read directly from logging file in ```$GRADEX_CLI_ROOT/var/log/gradex-cli.log```. They're in JSON format, one message per line, with the latest message appearing at the bottom of the file.

For debugging and development, I prefer to be able to search the messages, and analyse the log data using Kibana, which can be installed following this guide [here](https://www.elastic.co/guide/en/elastic-stack/current/installing-elastic-stack.html). I use Elastic, Logstash, Kibana and Filebeat. The installation is standard, in that Filebeat reads the logfile, passes it to logstash, which is configured with a JSON filter:

```
input {
  beats {
    port => 5044
  }
}
filter{
    json{
        source => "message"
    }
}
output {
  elasticsearch {
    hosts => ["http://localhost:9200"]
    index => "%{[@metadata][beat]}-%{[@metadata][version]}" 
  }
}
```
Logstash passes the logging message to elastic, and then you can discover it in Kibana, using the ```filebeat-*``` index.

## Usage

There are two main work flows. One for marking by script, the other for marking by question.

If it is your first time using ```gradex-cli```, then initialise the file structure by issuing the following command

```gradex-cli ingest```

### Ingest

In order to mark an exam, the first step is ingesting the raw pdf files with scanned student work. This is best done on a per-exam basis while you get used to the system. The system currently assumes all work is submitted through Blackboard Learn's Box, which provides a receipt in the following format ```.txt``` file

```
Name: Demonstrator Alpha (s0000000)
Assignment: Demo
Date Submitted: Tuesday, 21 April 2020 06:34:13 o'clock BST
Current Mark: Needs Marking

Submission Field:
There is no student submission text data for this assignment.

Comments:
There are no student comments for this assignment.

Files:
	Original filename: something from my computer.pdf
	Filename: demo-a.pdf
```

In the [Windows demo release](https://github.com/timdrysdale/gradex-cli/releases) there are three demo PDFs and associated "fake" Learn receipts in ```./demo-input```.

Place these demo files in ```$GRADEX_CLI_ROOT/ingest``` and invoke

```gradex-cli ingest```

Your files should disappear - although if there is something wrong, they will appear back in the ingest directory, as a way of letting you know what has been rejected.

You can see the files have been ingested into ```$GRADEX_CLI_ROOT\usr\exam\Demo\02-accepted-receipts``` and ```$GRADEX_CLI_ROOT\usr\exam\Demo\02-accepted-papers```.

We can prepare them for further processing by ```flattening``` them, and embedding information that we will use to anonymously track each page.

We require a system for swapping the known identity for an anonymous one. The ```identity.csv``` in ```$GRADEX_CLI_ROOT\etc\identity\``` contains two columns, ```identity``` and ```anonymous```. For the demo system, it comes with a simple version containing fake data:

![alt text][identity]

You can prepare the pages by

```
gradex-cli flatten Demo
```

You can manually inspect the flattened paper directory to see the newly flattened files

```
ls $GRADEX_CLI_ROOT/usr/exam/Demo/05-anonymous-papers
```

![alt text][flattened]



Now you can choose whether to mark by script or by question. Let's mark by question. First we need to add labelling side bars so our talented team of labellers (who we shall call `X`, somewhat mysteriously) can whizz through and tell us which page has what question on it:

```
gradex-cli label X Demo
gradex-cli export marking X Demo
```

Then we can send the files we find in ```$GRADEX_CLI_ROOT/export/Demo-question-ready-X/``` to our labllers. Once they send them back, we put them in ```$GRADEX_CLI_ROOT/ingest``` and ingest them with

```
gradex-cli ingest
```

You can manually inspect them to see that they end up in

```
ls $GRADEX_CLI_ROOT/usr/exam/Demo/10-question-back/X
```

### Marking

We want to prepare a set of pages for a marker with the initials ABC, so we issue

```
gradex-cli sort ABC Demo
gradex-cli export marking ABC Demo
```

We can get the files from export, and send to our marker.
```
$GRADEX_CLI_ROOT/export/Demo-marker-ready-ABC/
```
![alt text][marking]


You can try marking these files yourself, and save direct back to ingest (no need to change the filename, it will see from the hidden data what file it is). With the files back in the ingest directory after marking, we ingest again (same command as before)


### Processing marked files

We flatten the files to preserve the comments, read the textfields and optical boxes and store the data in the file, then we assemble documents that merge together the relevant pages for each file, by script.

Each page is categorised into exactly one of four categories, in order lowest to highest priority

    -- ```skipped``` - no indication from marker that they saw it
	-- ```seen``` - page-ok has had a character entered in a textfield or a stylus mark has been made in more than 2% of the area of the box (or a smaller amount if the box is rectangular)
	-- ```marked``` - something has been entered in one of the other textfields, by keyboard or stylus
	-- ```bad``` - the page-bad box has been ticked

Note that the priority is used to resolve what status to use when more than one applies. For example, a ```marked``` page that is also ```bad```, is given the status ```bad```. A page that is both ```marked``` and ```seen``` is given status ```marked```.


#### Page merge rules

Every page in the script is included at least once, for context. There is a merge summary bar on the side of each page, so you can tell at a glance if you should expect to see a duplicate copy of the page. If there are no ```marked``` pages, then one of the other pages is chosen. If there are more than one pages that are ```marked```, then all ```marked``` pages are included (e.g. if two markers share a question, and there are one or more pages that have material they both ended up marking).

#### Processing adjustments

Textfields are not easily edited by stylus, so for these markers, we expect them to annotate by hand. Then we'll get someone to key in the mark later. So as to retain the benefits of automation, we can use "optical" methods to check whether hand annotations have been made in the textfields, and if so, trigger the same actions as would have happened by typing into the ```page-ok``` and ```page-bad``` boxes.

##### Background colour for optical boxes

We assume a vanilla background (#ffffff) for the boxes, unless the flag ```--background-vanilla=false``` is given, e.g.

```
gradex-cli flatten marked 'Some Exam' --background-vanilla=false
```

in which case, the background is assumed to be chocolate (#000000).

###### Optical Box boundaries
There are some occasions when you get false positives from the optical-boxes, which is attributed without 100% certainty to artefacts from the boundary edges. It's even been the case in testing (before default shrinkage was increased to 6 pixels) where one marker's scripts threw 100% false positives on the ```page-bad``` box, but the other Marker on that script threw far fewer false positives. If you get a bunch of false positives (no marks in box visually, but pagedata contains "markDetected") then try setting the box shrinkage to a higher number. The number is the number of pixels in each direction. A 10mm by 10mm box at 175dpi has 69x69 pixels. The default shrink reduces that to (69-6-6)x(69-6-6) = 57x57 pixels. If you wanted to shrink some more, you could try for (69-10-10)x(69-10-10) = 49x49 pixels with

```
gradex-cli flatten marked 'Some exam' --box-shrink=10
```
Either or both flags can be issued in the same command. Note that flags must come AFTER the exam.

Also note the change from an imperative "mark" from the mark command, to the adjective "marked". Just to keep you on your toes, like. The imperative (command) here is "flatten."

#### Limitations

This page flattening and merging process _should_ work on the by-question batches (but has not been tested yet for that). Note that the flatten and merge phases of this step are implemented separately behind the scenes (for now), but are always performed at the same time, so the single command "flatten" is used to trigger one after the other.


### Moderating

Once our marked work is flattened, we are ready to put on the moderating bars. Since we might be doing this for more than one moderator, we don't link it to the previous step. At this stage of the workflow, both by-script and by-questions processes have return to the same path (```26-marked-ready```). With many scripts in this folder, the system automatically splits the set of scripts into a set to be actively moderated, with a green sidebar. The rest get a smaller grey "inactive" sidebar. Let's say we have moderate FFF who will moderate 10% or 10 scripts (whichever is greater) for 'Some Exam':


```
gradex-cli moderate FFF 'Some Exam'

```

Note: we don't currently support any other split ratios other than 10% or 10, whichever is bigger, but it is straightforward to add flags to do this if needed.


## Further procesing steps

There are further processing steps which are currently partly supported (check bars etc). These will be updated in a future release.

## Guidance to markers

Markers need only use Adobe Acrobat Reader (Free). The onedrive PDF app works on ipad, and Master PDF works on Linux. Most other viewers don't implement enough support for acroforms.

### Markers:

- can use a keyboard, or stylus
- do not need to rename their file


### Tech to avoid:

   -- Apple Preview (Quartz PDF) trashes the page catalog and prevents unipdf from reading the file
   -- Chrome lets you edit, but doesn't save
   -- Edge doesn't autosize the text in the boxes so it is not nice to use
   -- Almost everything on linux
   

### What to use on Linux:

Master PDF which can fill forms without registration being required


## Custom templates

Some exams will have different marking requirements. These can be accommodated by offering different layout templates that offer the same stages as the default process flow (mark, moderate-active, moderate-inactive, check - note that these are intended to be reused for remark remoderate recheck, but these are not fully supported yet. This modification offers an alternative 5-questions-per-page mark template via usage of the ```layout-q5.svg``` layout at mark stage. You can use a custom template with the mark command by issuing the layout flag at the command line. The template path is relative to ```$GRADEX_CLI_ROOT/etc/overlay/template```. For example, for the five-question markbar, issue:

```
gradex-cli mark <marker> <exam> --layout "layout-q5.svg"
```

Note that flags need to come after the exam (one of the positional arguments)

For detailed information on how to customise the templates using Inkscape, [see here](https://github.com/timdrysdale/gradex-cli/blob/master/parsesvg/README.md).
 

## TODO

- handle incoming marked/moderated/checked work
    - merge pages
	- report bad pages detected by markers
	- report results into csv, similar to [this](https://github.com/timdrysdale/gradex-extract)

- report results into csv, similar to [this](https://github.com/timdrysdale/gradex-extract)

### Done

- integrate [optical check box](https://github.com/timdrysdale/opticalcheckbox)
- integrate tree view [from here](https://github.com/timdrysdale/dt)

### Deferred

- integrate [optical handwriting recognition](https://github.com/sausheong/gonn)
- live marking tool to show staff running averages/totals/percentage completion

## Test coverage

```
comment	    coverage: 93.8% of statements
ingester	coverage: 58.1% of statements
optical	    coverage: 81.5% of statements
pagedata    coverage: 74.2% of statements
parselearn  coverage: 87.6% of statements
parsesvg    coverage: 81.8% of statements
tree        coverage: 63.4% of statements
```

## Codebase

Very close to 10 KLOC ....

```
--------------------------------------------------------------------------------
 Language             Files        Lines        Blank      Comment         Code
--------------------------------------------------------------------------------
 Go                      74        15873         3427         1273        11173
 Markdown                 8          788          254            0          534
 Plain Text              18          291           72            0          219
 Bourne Shell             1            5            2            2            1
 JSON                     1            1            0            0            1
--------------------------------------------------------------------------------
 Total                  102        16958         3755         1275        11928
--------------------------------------------------------------------------------
```

Most of the libraries are sub 1K, but the largest are:
```
ingester 4155
parsesvg  2753
```


[identity]: ./img/identity.png "csv file with anoynmous identity"
[flattened]: ./img/flattened.png "scanned page with header added"
[logo]: ./parsesvg/img/gradexTMlogo2-50pc.png "gradex logo"
[marking]: ./img/marking.png "red bars added for marking"
[status]: https://img.shields.io/badge/build-passing-green "build passing"
[teaching-matters]: ./img/teaching-matters.png "scroll of parchment"
