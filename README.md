# godelbrot and webdelbrot

## Summary

A web-based mandelbrot viewer in Go.

godelbrot has no dependencies outside of the go compiler, and its standard 
library.

The godelbrot developers are using go 1.4, mostly because godebug does not
support 1.5 at the date of publication.

## Usage

This package has been hacked up quickly, and we haven't made it go-gettable yet.

Wide distribution isn't a priority till we reach out development goals (more on 
them to come!)

To get started with the command line tool, clone this repository and then

    $ export GOPATH=/path/to/godelbrot
    $ cd /path/to/godelbrot
    $ go install functorama.com/demo/godelbrot
    $ bin/godelbrot

To run the web app

    $ export GOPATH=/path/to/godelbrot
    $ cd /path/to/godelbrot
    $ go install functorama.com/demo/webdelbrot
    $ bin/webdelbrot

Now point your browser to localhost:8080 and you should see the fractal.  Note
the --addr argument allows you to specify the interface.

Both applications have a set of command line options.  Try --help.

## Web app controls

Left click to begin highlighting zoom region.  Left click again to zoom.

Middle quick or "q" to cancel zoom selection.

## Credits

**John Morrice**

http://functorama.com

https://github.com/johnny-morrice

**Gavin Leech**

https://github.com/technicalities
