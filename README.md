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

The bit that handles merging questions coming back from markers is not integrated yet, so for this demo, we just manually move our files to where they will be after marking is analysed:

```
cp -r $GRADEX_CLI_ROOT/usr/exam/Demo/22-marker-back/* $GRADEX_CLI_ROOT/usr/exam/Demo/26-marked-ready/
```

Now we can prepare for moderating. The bit of the system that puts papers back into by-script files is not currently integrated. At this stage of the workflow, both by-script and by-questions processes return to the same path. With many scripts in this folder, the system automatically splits the set of scripts into a set to be actively moderated, with a green sidebar. The rest get a smaller grey "inactive" sidebar.


### TODO -- FINISH THIS SECTION




## Guidance to markers

Markers need only use Adobe Acrobat Reader (Free). The onedrive PDF app works on ipad, and Master PDF works on Linux. Most other viewers don't implement enough support for acroforms.

Markers:

- can use a keyboard, or stylus
- do not need to rename their file

## Custom templates

For detailed information on how to customise the templates using Inkscape, [see here](https://github.com/timdrysdale/gradex-cli/blob/master/parsesvg/README.md).


## TODO

- handle incoming marked/moderated/checked work
    - merge pages
	- report bad pages detected by markers
	- report results into csv, similar to [this](https://github.com/timdrysdale/gradex-extract)
- integrate [optical check box](https://github.com/timdrysdale/opticalcheckbox)
- integrate reporting
- integrate [optical handwriting recognition](https://github.com/sausheong/gonn)
- live marking tool to show staff running averages/totals/percentage completion
- integrate tree view [from here](https://github.com/timdrysdale/dt)


## Test coverage

```
ok  	github.com/timdrysdale/gradex-cli/comment	0.031s	coverage: 93.8% of statements
ok  	github.com/timdrysdale/gradex-cli/extract	0.041s	coverage: 49.3% of statements
ok  	github.com/timdrysdale/gradex-cli/ingester	14.762s	coverage: 59.1% of statements
ok  	github.com/timdrysdale/gradex-cli/pagedata	0.122s	coverage: 74.2% of statements
ok  	github.com/timdrysdale/gradex-cli/parselearn	0.010s	coverage: 87.6% of statements
ok  	github.com/timdrysdale/gradex-cli/parsesvg	15.017s	coverage: 81.6% of statements
```

[identity]: ./img/identity.png "csv file with anoynmous identity"
[flattened]: ./img/flattened.png "scanned page with header added"
[logo]: ./parsesvg/img/gradexTMlogo2-50pc.png "gradex logo"
[marking]: ./img/marking.png "red bars added for marking"
[status]: https://img.shields.io/badge/build-passing-green "build passing"
[teaching-matters]: ./img/teaching-matters.png "scroll of parchment"
