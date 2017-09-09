package mp4

import (
  "os"
  "fmt"
  "errors"
)

type MP4 struct {
  Boxes []interface{}
}

func Parse(f *os.File) (*MP4, error) {

  res := MP4{make([]interface{}, 0)};

  hdr := make([]byte, BOX_HDR_SZ_EXT);

  for {

    n,err := f.Read(hdr);

    if (err != nil) {
      if (err.Error() == "EOF") {
        break;
      }
      return nil, err;
    }

    if (n != BOX_HDR_SZ_EXT) {
      return nil, errors.New("not enough data");
    }

    b,_ := parseBox(hdr);

    f.Seek(int64(b.headerSize - BOX_HDR_SZ_EXT), 1);

    fmt.Println("-", b.Type);

    var data []byte;

    if (b.Type != "mdat") {
      data = make([]byte, b.Size - b.headerSize);
      n,err := f.Read(data);

      if (err != nil) {
        return nil, err;
      }

      if (n != int(b.Size - b.headerSize)) {
        return nil, errors.New("not enough data");
      }

    } else {
      f.Seek(int64(b.Size - b.headerSize), 1);
    }

    switch b.Type {
    case "ftyp":
      ftb,_ := parseFileTypeBox(data, b);
      res.Boxes = append(res.Boxes, *ftb);
    case "free":
      fallthrough
    case "skip":
      fsb,_ := parseFreeSpaceBox(data, b);
      res.Boxes = append(res.Boxes, *fsb);
    case "mdat":
      mdb,_ := parseMediaDataBox(data, b);
      res.Boxes = append(res.Boxes, *mdb);
    case "moov":
      mb,_ := parseMovieBox(data, b);
      res.Boxes = append(res.Boxes, *mb);
    }
  }

  fmt.Println("");

  return &res, nil;
}
