# PNG Embed

Embed key-value data into png images.

## Sample application

There is a sample application in the `cmd` directory which takes a source image and attempts to encode data into it. You can run this using `go run cmd/main.go` with the following arguments:

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
$ go run cmd/main.go -input in.png -key fruit -value apple -output out.png
```

You can then use something like [pngcheck](http://www.libpng.org/pub/png/apps/pngcheck.html) to verify that we did the right thing:
```shell
$ go run cmd/main.go -input ~/Desktop/test.png -key fruit -value apple -output out.png
$ pngcheck -t out.png
File: out.png (10785 bytes)
fruit:
    apple
OK: out.png (225x225, 8-bit palette, non-interlaced, 78.7%).
```

## Usage

Two ways to use this library:

### Embed (key, value) strings:

```go
import (
    "io/ioutil"
    "github.com/sabhiram/png-embed"
)

func MyAwesomeFunction(imgPath, outPath string) error {
    // Assuming that 'imgPath' points to a valid PNG image.
    data, err := pngembed.EmbedKeyValue(imgPath, "FOO", "BAR")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(outPath, data, 777)
}
```

### Embed (key, json) map:

If you have a structure that serializes nicely to JSON, you can also just use the `EmbedMap` method to store the serialized version of the data under `key` in the embedded region.

```go
import (
    "io/ioutil"
    "encoding/json"
    "github.com/sabhiram/png-embed"
)

func MyAwesomeFunction(imgPath, outPath string) error {
    s := struct {
        Foo string `json:"Foo"`
        Bar string `json:"Bar"`
    }{
        Foo: "FooValue",
        Bar: "BarValue",
    }

    data, err := pngembed.EmbedMap(imgPath, "FOO", s)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(outPath, data, 777)
}
```

## TODO:

1. Sister function to extract text data from a PNG?
2. Tests