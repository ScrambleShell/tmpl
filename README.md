# Terminal Command Template Engine

This is a template engine for configuring and running complex terminal commands. 

Take ffmpeg for example. To transcode a HD video you might use something along the lines of:

```
ffmpeg -i "in.mkv" -c:v libx264 -preset slow -crf 18 \
-filter:v "crop=1904:800:10:140" -c:a libfdk_aac \
-b:a 448k -c:s copy -threads 4 -y "out.mkv"
```

Wouldn't it be much nicer to have a YAML configuration file for defining the various parameters? For example:

```
video:
    codec: libx264
    preset: slow
    crf: 18
    filter: "crop=1904:800:10:140"
audio:
    codec: libfdk_aac
    bitrate: 448k
subtitles:
    codec: copy
input: $0
output: $1
threads: 4
```

... and then apply it to run the ffmpeg command?

```
tmpl run ffmpeg config.yaml in.mkv out.mkv
```

This is just one of many possible uses of this engine.

## Installation

Make sure you have a recent version of [Go](https://golang.org/) installed. Then use:

```
go get github.com/sf1/tmpl
```

## Usage

### Defining templates

Currently the engine expects all templates to be located in

```
$HOME/templates
```

Templates are defined using Go's template syntax, which is documented [here](http://golang.org/pkg/html/template/).

### Creating configuration files

Configuration files are plain YAML files used to pass parameters when templates are applied or "run". Their structure depends on each template. Configuration files may contain placeholders for additional command line parameters:

```
input: $0
output: $1
```

$0 and $1 refer to parameters that can be specified when the engine is invoked from the terminal. The main use of this feature is to allow specifying input and output files.

### Applying and "Running" templates

Templates can be applied or "run". When templates are applied, the engine outputs the result but does not attempt to run it. For example

```
tmpl apply ffmpeg config.yaml in.mkv out.mkv
```

outputs the ffmpeg command shown at the top of this readme, but does not run it.

```
tmpl run ffmpeg config.yaml in.mkv out.mkv
```

applies the template and attempts to run the result.

