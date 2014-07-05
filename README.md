# Crawler

**Crawler** is a concurrent URL web crawler written in **Go**.

## Usage

Following setup crawls all URLs at
[google.com?q=golang](https://google.com?q=golang):

    package main

    import (
        "fmt"

        "github.com/bodokaiser/go-crawler"
    )

    func main() {
        // create a new crawler struct
        c := crawler.New()

        // start crawling on following url
        c.Open("https://google.com?q=golang")

        // listen to the received results
        c.Listen(listener)
    }

    func listener(p crawler.Pipeline) {
        r := <-p

        fmt.Printf("Found %d URLs at %s", len(r.URLS), r.Origin)
    }

## Install

Use `go get` to clone the repository:

    $ go get github.com/bodokaiser/go-crawler

## License

Copyright 2014 Bodo Kaiser <i@bodokaiser.io>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
