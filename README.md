# PNG Embed

Embed key-value data into png images.

## Sample application

There is a sample application in the `cmd` directory which takes a source image and attempts to encode data into it. You can run this using `go run cmd/main.go` with the following arguments:

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

## Usage

The pngencode package only supports one method `Embed()` which accepts a png file as an input, a key and value to store in a `tEXt` section right after the header section in the image. `Embed()` returns a slice of bytes and an appropriate error.

