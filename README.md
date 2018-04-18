# PNG Embed

Embed key-value data into png images.


## Install

```
go get github.com/sabhiram/png-embed
```


## Usage

This library exposes two APIs:

### `Embed(data []byte, key string, value interface{}) ([]byte, error)`

```go
package main

import (
    "io/ioutil"
    pngembed "github.com/sabhiram/png-embed"
)

func main() {
    bs, _ := ioutil.ReadFile("sample.png")

    // Encode the key "FOO" with the value "BAR" (string).
    data, _ := pngembed.Embed(bs, "FOO", "BAR")
    ioutil.WriteFile("sample.png", data, 755)
}
```

### `EmbedFile(filePath, key string, value interface{}) ([]byte, error)`

`EmbedFile` is exactly like `Embed` but it accepts a path to a file instead of the PNG file data.

```go
package main

import (
    "io/ioutil"
    pngembed "github.com/sabhiram/png-embed"
)

func main() {
    // Encode the key "FOO" with the value "BAR" (string).
    data, _ := pngembed.EmbedFile("sample.png", "FOO", "BAR")
    ioutil.WriteFile("sample.png", data, 755)
}
```

### Embedding JSON

Since the `Embed` method accepts an interface and assumes that the interface adheres to JSON encoding, we can pass arbitrary `struct`s or `map[interface{}]interface{}`s to it.

```go
package main

import (
    "io/ioutil"
    pngembed "github.com/sabhiram/png-embed"
)

func main() {
    s := struct {
        Foo string `json:"Foo"`
        Bar string `json:"Bar"`
    }{
        Foo: "FooValue",
        Bar: "BarValue",
    }

    data, _ := pngembed.EmbedFile("sample.png", "FOO", s)
    ioutil.WriteFile("sample.png", data, 755)
}
```


## Sample application

There is a sample application in the `example` directory which takes a source image and attempts to encode data into it. You can run this using `go run example/main.go` with the following arguments:

```
  -input string
        input file name for the png [required]
  -key string
        key name for the data to inject [default="TEST_KEY"]
  -value string
        value for the data to inject [default="TEST_VALUE"]
  -output string
        output file name for the png [default="out.png"]
```

To inject `in.png` with the key value pair "fruit": "apple" and generate out.png:
```shell
$ go run example/main.go -input in.png -key fruit -value apple -output out.png
```

You can then use something like [pngcheck](http://www.libpng.org/pub/png/apps/pngcheck.html) to verify that we did the right thing:
```shell
$ go run example/main.go -input ~/Desktop/test.png -key fruit -value apple -output out.png
$ pngcheck -t out.png
File: out.png (10785 bytes)
fruit:
    apple
OK: out.png (225x225, 8-bit palette, non-interlaced, 78.7%).
```
