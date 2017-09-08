package mp4

import (
  //"fmt"
  "bufio"
  "encoding/binary"
)

const BOX_HDR_SZ = 8;
const BOX_HDR_SZ_EXT = 16;

type Box struct {
  Size uint64
  Type string
}

type FullBox struct {
  box Box
  version uint8
  flags [3]byte
}

type FileTypeBox struct {
  box Box
  major_brand string
  minor_version uint32
  compatible_brands []string
}

type FreeSpaceBox struct {
  box Box
  data []uint8
}

type MP4 struct {
  Boxes []interface{}
}

func parseBox(r *bufio.Reader) (*Box, error) {
  var b Box;

  data,err := r.Peek(BOX_HDR_SZ);

  if (err != nil) {
    return nil,err;
  }

  b.Size = uint64(binary.BigEndian.Uint32(data[0:4]));

  data = data[4:8];

  b.Type = string(data);

  r.Discard(BOX_HDR_SZ);

  if (b.Size == 1) {
    data,err = r.Peek(BOX_HDR_SZ_EXT - BOX_HDR_SZ);
    if (err != nil) {
      return nil,err;
    }
    b.Size = binary.BigEndian.Uint64(data);
    r.Discard(BOX_HDR_SZ_EXT - BOX_HDR_SZ);
  }

  return &b,nil;
}

func parseFileTypeBox(r *bufio.Reader, b *Box) (*FileTypeBox, error) {
  ftb := FileTypeBox{box: *b};

  data,err := r.Peek(8);

  if (err != nil) {
    return nil, err;
  }

  ftb.major_brand = string(data[0:4]);
  ftb.minor_version = binary.BigEndian.Uint32(data[4:8]);

  r.Discard(8);

  n := int((ftb.box.Size - 16) / 4);

  ftb.compatible_brands = make([]string, n);

  for i := 0; i < n; i++ {

    data,err = r.Peek(4);

    if (err != nil) {
      return nil, err;
    }

    ftb.compatible_brands[i] = string(data);

    r.Discard(4);
  }

  return &ftb, nil;
}

func parseFreeSpaceBox(r *bufio.Reader, b *Box) (*FreeSpaceBox, error) {
  fsb := FreeSpaceBox{box: *b};
  return &fsb, nil;
}

func Parse(r *bufio.Reader) (*MP4, error) {

  res := MP4{make([]interface{}, 0)};

  for i := 0; i < 3; i++ {

    b,err := parseBox(r);

    if (err != nil) {
      return nil, err;
    }

    switch b.Type {
    case "ftyp":
      ftb,_ := parseFileTypeBox(r, b);
      res.Boxes = append(res.Boxes, *ftb);
    case "free":
      fsb,_ := parseFreeSpaceBox(r, b);
      res.Boxes = append(res.Boxes, *fsb);
    }
  }

  // if (err != nil) {
  //   return nil, err;
  // }

  return &res, nil;
}
