# lightroom2aftershot

A simple preset converter that converts lightroom presets to aftershot.

There are a lot of good lightroom presets out there. However, for aftershot there are pretty much
none outside of the paid packs offered by corel (and even those are of fairly low quality compared
to some lightroom presets out there).

This application aims to convert lightroom presets to aftershot presets.

## Usage

```
$ lightroom2aftershot lightroom-preset.xml > aftershot-preset.xml
```

## This is not perfect

It is important to note, that the conversion being done here is not perfect. Aftershot interprets
some values differently from lightroom and some concepts don't even exist in Aftershot at all.

The aftershot presets generated by this converter probably have to be tweaked - however, they serve
as good starting points if you have a look you are after that is produced by a Lightroom preset you
like.

If you find presets that produce wildly different results (or that don't work at all), be sure to open
an issue containing the lightroom prefix used and, if possible, a RAW file to reproduce the issue
with. It would also be nice to have comparison screenshots of what the output is supposed to look
like in lightroom.