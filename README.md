# PNG Embed

A very simple (and currently rather dull) data embedder for PNG images.

## Usage

There is a sample application in the `cmd' directory which takes a source image and attempts to encode data into it. You can run this using `go run cmd/main.go` with the following arguments:
```
  -i string
        input file name for the png
  -input string
        input file name for the png
  -k string
        key name for the data to inject
  -key string
        key name for the data to inject
  -o string
        output file name for the png
  -output string
        output file name for the png
  -v string
        value for the data to inject
  -value string
        value for the data to inject
```

## Library

Currently the pngencode library only supports one method `Encode()` which accepts a png file as an input, a key and value to store in a `tEXt` section right after the header section in the image. `Encode()` returns a slice of bytes.