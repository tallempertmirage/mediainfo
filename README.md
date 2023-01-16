# mediainfo

## Installation

1. Download mediainfo for the machine os https://mediaarea.net/en/MediaInfo/Download
2. Add to PATH

Ubuntu - https://snapcraft.io/install/mediainfo/ubuntu
Mac - brew install mediainfo

## Usage
## Example
```go
package main

import (
	"encoding/json"
	"fmt"
	"gitlab.com/mirage-dynamics/indexing/mediainfo"
)

func main() {
	mediainfo, _ := mediainfo.GetMediaInfo("/Users/Fang/Movies/11.flv")
	info, _ := json.Marshal(mediainfo)
	fmt.Println(string(info), err)
}

