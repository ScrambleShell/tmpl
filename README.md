# Terminal Command Template Engine

This is a template engine for configuring and running complex terminal commands. 

Take ffmpeg for example. To transcode a HD video you might use something along the lines of:

```
ffmpeg -i "in.mkv" -c:v libx264 -preset slow -crf 18 \
-filter:v "crop=1904:800:10:140" -c:a libfdk_aac \
-b:a 448k -c:s copy -threads 4 -y "out.mkv"
```

Wouldn't it be much nicer to have a yaml configuration file for defining the various parameters? For example:

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

...

### Creating configuration files

...

### Applying and "Running" templates

...
