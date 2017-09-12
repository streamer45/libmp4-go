package mp4

import (
  "fmt"
  "math"
  "errors"
  "encoding/binary"
)

func parseBox(data []byte) (*Box, error) {
  var b Box;

  hdr := data[0:BOX_HDR_SZ];

  b.headerSize = BOX_HDR_SZ;

  b.Size = uint64(binary.BigEndian.Uint32(hdr[0:4]));
  b.Type = string(hdr[4:8]);

  if (b.Size == 1) {
    b.headerSize = BOX_HDR_SZ_EXT;
    hdr = data[BOX_HDR_SZ:BOX_HDR_SZ_EXT]
    b.Size = binary.BigEndian.Uint64(hdr);
  }

  return &b,nil;
}

func parseFullBox(data []byte, b *Box) (*FullBox, error) {
  fb := FullBox{Box: *b};
  fb.Version = data[0];
  copy(fb.Flags[:], data[1:4]);
  return &fb, nil;
}

func parseFileTypeBox(data []byte, b *Box) (*FileTypeBox, error) {
  ftb := FileTypeBox{Box: *b};

  ftb.MajorBrand = string(data[0:4]);
  ftb.MinorVersion = binary.BigEndian.Uint32(data[4:8]);

  data = data[8:];

  n := int((ftb.Box.Size - ftb.Box.headerSize - 8) / 4);

  ftb.CompatibleBrands = make([]string, n);

  for i := 0; i < n; i++ {
    ftb.CompatibleBrands[i] = string(data[i * 4: (i + 1) * 4]);
  }

  return &ftb, nil;
}

func parseFreeSpaceBox(data []byte, b *Box) (*FreeSpaceBox, error) {
  fsb := FreeSpaceBox{Box: *b};
  fsb.data = data;
  return &fsb, nil;
}

func parseMediaDataBox(data []byte, b *Box) (*MediaDataBox, error) {
  mdb := MediaDataBox{Box: *b};
  mdb.data = data;
  return &mdb, nil;
}

func parseMovieHeaderBox(data []byte, b *Box) (*MovieHeaderBox, error) {
  fb,_ := parseFullBox(data, b);
  mhb := MovieHeaderBox{Box: *fb};

  data = data[4:];

  if (fb.Version == 1) {
    mhb.Ctime = binary.BigEndian.Uint64(data[0:8]);
    mhb.Mtime = binary.BigEndian.Uint64(data[8:16]);
    mhb.Timescale = binary.BigEndian.Uint32(data[16:20]);
    mhb.Duration = binary.BigEndian.Uint64(data[20:28]);
    data = data[28:];
  } else {
    mhb.Ctime = uint64(binary.BigEndian.Uint32(data[0:4]));
    mhb.Mtime = uint64(binary.BigEndian.Uint32(data[4:8]));
    mhb.Timescale = binary.BigEndian.Uint32(data[8:12]);
    mhb.Duration = uint64(binary.BigEndian.Uint32(data[12:16]));
    data = data[16:];
  }

  fixed := binary.BigEndian.Uint32(data[0:4]);
  mhb.Rate = float32(fixed) / float32(math.Pow(2, 16));

  fixed2 := binary.BigEndian.Uint16(data[4:6]);
  mhb.Volume = float32(fixed2) / float32(math.Pow(2, 8));

  data = data[16:];

  for i:=0; i < 3; i++ {
    mhb.Matrix[i][0] = float32(binary.BigEndian.Uint32(data[0:4])) / float32(math.Pow(2, 16));
    mhb.Matrix[i][1] = float32(binary.BigEndian.Uint32(data[4:8])) / float32(math.Pow(2, 16));
    mhb.Matrix[i][2] = float32(binary.BigEndian.Uint32(data[8:12])) / float32(math.Pow(2, 30));
    data = data[12:];
  }

  data = data[9 * 4 + 6 * 4:];

  mhb.NextTrackID = binary.BigEndian.Uint32(data[0:4]);

  return &mhb, nil;
}

func parseTrackHeaderBox(data []byte, b *Box) (*TrackHeaderBox, error) {
  fb,_ := parseFullBox(data, b);
  thb := TrackHeaderBox{Box: *fb};

  data = data[4:];

  if (fb.Version == 1) {
    thb.Ctime = binary.BigEndian.Uint64(data[0:8]);
    thb.Mtime = binary.BigEndian.Uint64(data[8:16]);
    thb.TrackID = binary.BigEndian.Uint32(data[16:20]);
    thb.Duration = binary.BigEndian.Uint64(data[24:32]);
    data = data[32:];
  } else {
    thb.Ctime = uint64(binary.BigEndian.Uint32(data[0:4]));
    thb.Mtime = uint64(binary.BigEndian.Uint32(data[4:8]));
    thb.TrackID = binary.BigEndian.Uint32(data[8:12]);
    thb.Duration = uint64(binary.BigEndian.Uint32(data[16:20]));
    data = data[20:];
  }

  data = data[52:];

  fixed := binary.BigEndian.Uint32(data[0:4]);
  thb.Width = float32(fixed) / float32(math.Pow(2, 16));

  fixed2 := binary.BigEndian.Uint32(data[4:8]);
  thb.Height = float32(fixed2) / float32(math.Pow(2, 16));

  return &thb, nil;
}

func parseMediaHeaderBox(data []byte, b *Box) (*MediaHeaderBox, error) {
  fb,_ := parseFullBox(data, b);
  mhb := MediaHeaderBox{Box: *fb};

  data = data[4:];

  if (fb.Version == 1) {
    mhb.Ctime = binary.BigEndian.Uint64(data[0:8]);
    mhb.Mtime = binary.BigEndian.Uint64(data[8:16]);
    mhb.Timescale = binary.BigEndian.Uint32(data[16:20]);
    mhb.Duration = binary.BigEndian.Uint64(data[20:28]);
    data = data[28:];
  } else {
    mhb.Ctime = uint64(binary.BigEndian.Uint32(data[0:4]));
    mhb.Mtime = uint64(binary.BigEndian.Uint32(data[4:8]));
    mhb.Timescale = binary.BigEndian.Uint32(data[8:12]);
    mhb.Duration = uint64(binary.BigEndian.Uint32(data[12:16]));
    data = data[16:];
  }

  code := [3]byte{};

  code[0] = (data[0] >> 2) + 0x60;
  code[1] = (((data[0] & 0x03) << 3) | data[1] >> 5) + 0x60;
  code[2] = (data[1] & 0x1f) + 0x60;

  mhb.Language = string(code[0:3]);

  return &mhb, nil;
}

func parseHandlerBox(data []byte, b *Box) (*HandlerBox, error) {
  fb,_ := parseFullBox(data, b);
  hb := HandlerBox{Box: *fb};
  data = data[4:];
  hb.HandlerType = string(data[4:8]);
  hb.Name = string(data[8 + 12:b.Size - b.headerSize - 4 - 1]);
  return &hb, nil;
}

func parsePixelAspectRatioBox(data []byte, b *Box) (*PixelAspectRatioBox, error) {
  pasp := PixelAspectRatioBox{Box: *b};
  pasp.HSpacing = binary.BigEndian.Uint32(data[0:4]);
  pasp.VSpacing = binary.BigEndian.Uint32(data[4:8]);
  return &pasp, nil;
}

func parseAVCcBox(data []byte, b *Box) (*AVCcBox, error) {
  avcc := AVCcBox{Box: *b};

  avcc.Version = data[0];
  avcc.Profile = data[1];
  avcc.Level = data[3];
  avcc.SizeLen = (data[4] & 0x03) + 1;

  nsps := int(data[5] & 0x1f);
  data = data[6:];

  for i := 0; i < nsps; i++ {
    len := binary.BigEndian.Uint16(data[0:2]);
    sps := make([]byte, len);
    copy(sps, data[2: len + 2]);
    avcc.SPS = append(avcc.SPS, sps);
    data = data[len + 2:];
  }

  npps := int(data[0]);
  data = data[1:];

  for i := 0; i < npps; i++ {
    len := binary.BigEndian.Uint16(data[0:2]);
    pps := make([]byte, len);
    copy(pps, data[2: len + 2]);
    avcc.PPS = append(avcc.PPS, pps);
    data = data[len + 2:];
  }

  return &avcc, nil;
}

func parseVideoSampleDesc(data []byte, entry *SampleEntry) (*VideoSampleDescription, error) {
  vsd := VideoSampleDescription{};
  vsd.Width = binary.BigEndian.Uint16(data[16:18]);
  vsd.Height = binary.BigEndian.Uint16(data[18:20]);
  vsd.Depth = binary.BigEndian.Uint16(data[66:68]);
  // TODO add missing info

  data = data[70:];

  tsize := entry.Box.headerSize + 8 + 70;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("              -", b.Type);

    switch (b.Type) {
    case "avcC":
      avcc,_ := parseAVCcBox(data[b.headerSize:], b);
      entry.Extensions = append(entry.Extensions, *avcc);
    case "pasp":
      pasp,_ := parsePixelAspectRatioBox(data[b.headerSize:], b);
      entry.Extensions = append(entry.Extensions, *pasp);
    }

    tsize += b.Size;

    if (tsize == entry.Box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &vsd, nil;
}

func parseESDescriptor(data []byte) (*ESDescriptor, error) {
  d := ESDescriptor{};

  d.Tag = data[0];
  data = data[1:];

  if ((data[0] & 0x80) != 0x00) {
    data = data[3:];
  }

  d.Length = data[0];
  d.Id = binary.BigEndian.Uint16(data[1:3]);

  data = data[4:];
  tag := data[0];
  data = data[1:];

  if (tag != 0x04) {
    return nil, errors.New("invalid descriptor tag");
  }

  if ((data[0] & 0x80) != 0x00) {
    data = data[3:];
  }

  len := data[0];
  data = data[14:];
  tag = data[0];
  data = data[1:];

  if (tag != 0x05) {
    return nil, errors.New("invalid descriptor tag");
  }

  if ((data[0] & 0x80) != 0x00) {
    data = data[3:];
  }

  len = data[0];
  data = data[1:];
  d.Config = make([]byte, len);
  copy(d.Config, data[:len]);

  return &d, nil;
}

func parseElementaryStreamDescBox(data []byte, b *Box) (*ElementaryStreamDescBox, error) {
  esds := ElementaryStreamDescBox{Box: *b};
  data = data[4:];

  d,err := parseESDescriptor(data);

  if (err != nil) {
    fmt.Println(err)
    return &esds, nil;
  }

  esds.Esd = *d;

  return &esds, nil;
}

func parseSoundSampleDesc(data []byte, entry *SampleEntry) (*SoundSampleDescription, error) {
  ssd := SoundSampleDescription{};
  ssd.Channels = binary.BigEndian.Uint16(data[8:10]);
  ssd.SampleSize = binary.BigEndian.Uint16(data[10:12]);
  fixed := binary.BigEndian.Uint32(data[16:20]);
  ssd.SampleRate = float32(fixed) / float32(math.Pow(2, 16));

  tsize := entry.Box.headerSize + 8 + 20;

  data = data[20:];

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("              -", b.Type);

    switch (b.Type) {
    case "esds":
      esds,_ := parseElementaryStreamDescBox(data[b.headerSize:b.Size], b);
      entry.Extensions = append(entry.Extensions, *esds);
    }

    tsize += b.Size;

    if (tsize == entry.Box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &ssd, nil;
}

func parseSampleDescBox(data []byte, b *Box) (*SampleDescriptionBox, error) {
  fb,_ := parseFullBox(data, b);
  sdb := SampleDescriptionBox{Box: *fb};

  data = data[4:];

  sdb.EntryCount = binary.BigEndian.Uint32(data[0:4]);

  data = data[4:];

  tsize := b.headerSize + 4 + 4;

  for i := 0; i < int(sdb.EntryCount); i++ {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("            -", b.Type);

    entry := SampleEntry{Box: *b};

    data = data[b.headerSize:];

    entry.DataRefIndex = binary.BigEndian.Uint16(data[6:8]);

    data = data[8:];

    switch (b.Type) {
    case "avc1":
      vsd,_ := parseVideoSampleDesc(data, &entry);
      entry.SampleDesc = *vsd;
    case "mp4a":
      ssd,_ := parseSoundSampleDesc(data, &entry);
      entry.SampleDesc = *ssd;
    }

    sdb.Entries = append(sdb.Entries, entry);

    tsize += b.Size;

    if (tsize == sdb.Box.Box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &sdb, nil;
}

func parseTimeToSampleBox(data []byte, b *Box) (*TimeToSampleBox, error) {
  fb,_ := parseFullBox(data, b);
  ttsb := TimeToSampleBox{Box: *fb};
  data = data[4:];

  ttsb.EntryCount = binary.BigEndian.Uint32(data[0:4]);

  data = data[4:];

  ttsb.SampleCount = make([]uint32, ttsb.EntryCount);
  ttsb.SampleDelta = make([]uint32, ttsb.EntryCount);

  for i := 0; i < int(ttsb.EntryCount); i++ {
    ttsb.SampleCount[i] = binary.BigEndian.Uint32(data[0:4]);
    ttsb.SampleDelta[i] = binary.BigEndian.Uint32(data[4:8]);
    data = data[8:];
  }

  return &ttsb, nil;
}

func parseCompTimeToSampleBox(data []byte, b *Box) (*CompTimeToSampleBox, error) {
  fb,_ := parseFullBox(data, b);
  ctsb := CompTimeToSampleBox{Box: *fb};
  data = data[4:];

  ctsb.EntryCount = binary.BigEndian.Uint32(data[0:4]);

  data = data[4:];

  ctsb.SampleCount = make([]int32, ctsb.EntryCount);
  ctsb.SampleOffset = make([]int32, ctsb.EntryCount);

  for i := 0; i < int(ctsb.EntryCount); i++ {
    ctsb.SampleCount[i] = int32(binary.BigEndian.Uint32(data[0:4]));
    ctsb.SampleOffset[i] = int32(binary.BigEndian.Uint32(data[4:8]));
    data = data[8:];
  }

  return &ctsb, nil;
}

func parseSyncSampleBox(data []byte, b *Box) (*SyncSampleBox, error) {
  fb,_ := parseFullBox(data, b);
  ssb := SyncSampleBox{Box: *fb};
  data = data[4:];

  ssb.EntryCount = binary.BigEndian.Uint32(data[0:4]);

  data = data[4:];

  ssb.SampleNumber = make([]uint32, ssb.EntryCount);

  for i := 0; i < int(ssb.EntryCount); i++ {
    ssb.SampleNumber[i] = binary.BigEndian.Uint32(data[0:4]);
    data = data[4:];
  }

  return &ssb, nil;
}

func parseSampleToChunkBox(data []byte, b *Box) (*SampleToChunkBox, error) {
  fb,_ := parseFullBox(data, b);
  stcb := SampleToChunkBox{Box: *fb};
  data = data[4:];

  stcb.EntryCount = binary.BigEndian.Uint32(data[0:4]);

  data = data[4:];

  stcb.FirstChunk = make([]int32, stcb.EntryCount);
  stcb.SamplesPerChunk = make([]int32, stcb.EntryCount);
  stcb.SampleDescIndex = make([]int32, stcb.EntryCount);

  for i := 0; i < int(stcb.EntryCount); i++ {
    stcb.FirstChunk[i] = int32(binary.BigEndian.Uint32(data[0:4]));
    stcb.SamplesPerChunk[i] = int32(binary.BigEndian.Uint32(data[4:8]));
    stcb.SampleDescIndex[i] = int32(binary.BigEndian.Uint32(data[8:12]));
    data = data[12:];
  }

  return &stcb, nil;
}

func parseSampleSizeBox(data []byte, b *Box) (*SampleSizeBox, error) {
  fb,_ := parseFullBox(data, b);
  ssb := SampleSizeBox{Box: *fb};
  data = data[4:];

  ssb.SampleSize = binary.BigEndian.Uint32(data[0:4]);
  ssb.SampleCount = binary.BigEndian.Uint32(data[4:8]);

  data = data[8:];

  if (ssb.SampleSize == 0) {
    ssb.EntrySize = make([]uint32, ssb.SampleCount);
    for i := 0; i < int(ssb.SampleCount); i++ {
      ssb.EntrySize[i] = binary.BigEndian.Uint32(data[0:4]);
      data = data[4:];
    }
  }

  return &ssb, nil;
}

func parseChunkOffsetBox(data []byte, b *Box) (*ChunkOffsetBox, error) {
  fb,_ := parseFullBox(data, b);
  cob := ChunkOffsetBox{Box: *fb};
  data = data[4:];

  cob.EntryCount = binary.BigEndian.Uint32(data[0:4]);

  data = data[4:];

  cob.ChunkOffset = make([]uint32, cob.EntryCount);
  for i := 0; i < int(cob.EntryCount); i++ {
    cob.ChunkOffset[i] = binary.BigEndian.Uint32(data[0:4]);
    data = data[4:];
  }

  return &cob, nil;
}

func parseSampleTableBox(data []byte, b *Box) (*SampleTableBox, error) {
  stb := SampleTableBox{Box: *b};

  tsize := b.headerSize;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("          -", b.Type);

    switch b.Type {
    case "stsd":
      sdb,_ := parseSampleDescBox(data[b.headerSize:b.Size], b);
      stb.Stsd = *sdb;
    case "stts":
      ttsb,_ := parseTimeToSampleBox(data[b.headerSize:b.Size], b);
      stb.Stts = *ttsb;
    case "stss":
      ssb,_ := parseSyncSampleBox(data[b.headerSize:b.Size], b);
      stb.Stss = *ssb;
    case "ctts":
      ctts,_ := parseCompTimeToSampleBox(data[b.headerSize:b.Size], b);
      stb.Ctss = *ctts;
    case "stsc":
      stcb,_ := parseSampleToChunkBox(data[b.headerSize:b.Size], b);
      stb.Stsc = *stcb;
    case "stsz":
      ssb,_ := parseSampleSizeBox(data[b.headerSize:b.Size], b);
      stb.Stsz = *ssb;
    case "stco":
      cob,_ := parseChunkOffsetBox(data[b.headerSize:b.Size], b);
      stb.Stco = *cob;
    }

    tsize += b.Size;

    if (tsize == stb.Box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &stb, nil;

}

func parseMediaInfoBox(data []byte, b *Box) (*MediaInfoBox, error) {
  mib := MediaInfoBox{Box: *b};

  tsize := b.headerSize;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("        -", b.Type);

    switch b.Type {
    case "stbl":
      stb,_ := parseSampleTableBox(data[b.headerSize:], b);
      mib.Stbl = *stb;
    }

    tsize += b.Size;

    if (tsize == mib.Box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &mib, nil;
}

func parseMediaBox(data []byte, b *Box) (*MediaBox, error) {
  mb := MediaBox{Box: *b};

  tsize := b.headerSize;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("      -", b.Type);

    switch b.Type {
    case "mdhd":
      mhb,_ := parseMediaHeaderBox(data[b.headerSize:], b);
      mb.Mdhd = *mhb;
    case "hdlr":
      hb,_ := parseHandlerBox(data[b.headerSize:], b);
      mb.Hdlr = *hb;
    case "minf":
      mib,_ := parseMediaInfoBox(data[b.headerSize:], b);
      mb.Minf = *mib;
    }

    tsize += b.Size;

    if (tsize == mb.Box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &mb, nil;
}

func parseTrackBox(data []byte, b *Box) (*TrackBox, error) {
  tb := TrackBox{Box: *b};

  tsize := b.headerSize;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("    -", b.Type);

    switch b.Type {
    case "tkhd":
      thb,_ := parseTrackHeaderBox(data[b.headerSize:], b);
      tb.Tkhd = *thb;
    case "mdia":
      mb,_ := parseMediaBox(data[b.headerSize:], b);
      tb.Mdia = *mb;
    }

    tsize += b.Size;

    if (tsize == tb.Box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &tb, nil;
}

func parseMovieBox(data []byte, b *Box) (*MovieBox, error) {
  mb := MovieBox{Box: *b};

  tsize := b.headerSize;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("  -", b.Type);

    switch b.Type {
    case "mvhd":
      mhb,_ := parseMovieHeaderBox(data[b.headerSize:], b);
      mb.Mvhd = *mhb;
    case "trak":
      tb,_ := parseTrackBox(data[b.headerSize:], b);
      mb.Tracks = append(mb.Tracks, *tb);
    }

    tsize += b.Size;

    if (tsize == mb.Box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &mb, nil;
}
