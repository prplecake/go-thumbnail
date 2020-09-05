# Thumbnail Generation Package for Go

This package provides method to create thumbnails from provided images.

## Installation

Use to `go` command:

```
$ go get git.sr.ht/~mjorgensen/go-thumbnail
```

## Example

```go
package main

import (
	"fmt"
	"io/ioutil"

	"git.sr.ht/~mjorgensen/go-thumbnail"
)

func main() {
	thumbCfg := thumbnail.Configuration{
		Path: "./image.jpeg",
		ContentType: "image/jpeg",
		DestinationPrefix: "thumb_",
	}
	testImage, err := ioutil.ReadFile(thumbCfg.Path)
	if err != nil {
		panic(err)
	}

	err = Create(testImage, thumbCfg)
	if err != nil {
		panic(err)
	}
}
```

## Resources

Comprehensive documentation still needs to be written but will
eventually [be found here][man].

Discussion and patches are welcome and should be directed to my public
inbox for now: [~mjorgensen/public-inbox@lists.sr.ht][lists]. Please use
`--subject-prefix PATCH go-thumbnail` for clarity when sending patches.

Bugs, issues, planning, and tasks can all be found at the tracker:
[~mjorgensen/go-thumbnail][todo].

[man]:https://man.sr.ht/~mjorgensen/go-thumbnail
[lists]:https://lists.sr.ht/~mjorgensen/public-inbox
[todo]:https://todo.sr.ht/~mjorgensen/go-thumbnail
