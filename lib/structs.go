package mp4

const BOX_HDR_SZ = 8;
const BOX_HDR_SZ_EXT = 16;

type Box struct {
  Type string `json:"type"`
  Size uint64 `json:"size"`
  headerSize uint64
}

type FullBox struct {
  Box Box `json:"box"`
  Version uint8 `json:"version"`
  Flags [3]byte `json:"flags"`
}

type FileTypeBox struct {
  Box Box `json:"box"`
  MajorBrand string `json:"majorBrand"`
  MinorVersion uint32 `json:"minorVersion"`
  CompatibleBrands []string `json:"compatibleBrands"`
}

type FreeSpaceBox struct {
  Box Box `json:"box"`
  data []byte
}

type MediaDataBox struct {
  Box Box `json:"box"`
  data []byte
}

type MovieHeaderBox struct {
  Box FullBox `json:"fullBox"`
  Ctime uint64 `json:"creationTime"`
  Mtime uint64 `json:"modificationTime"`
  Timescale uint32 `json:"timescale"`
  Duration uint64 `json:"duration"`
  Rate float32 `json:"rate"`
  Volume float32 `json:"volume"`
  Matrix [3][3]float32 `json:"matrix"`
  NextTrackID uint32 `json:"nextTrackID"`
}

type TrackHeaderBox struct {
  Box FullBox `json:"fullBox"`
  Ctime uint64 `json:"creationTime"`
  Mtime uint64 `json:"modificationTime"`
  TrackID uint32 `json:"trackID"`
  Duration uint64 `json:"duration"`
  Width float32 `json:"width"`
  Height float32 `json:"height"`
}

type MediaHeaderBox struct {
  Box FullBox `json:"fullBox"`
  Ctime uint64 `json:"creationTime"`
  Mtime uint64 `json:"modificationTime"`
  Timescale uint32 `json:"timescale"`
  Duration uint64 `json:"duration"`
  Language string `json:"language"`
}

type HandlerBox struct {
  Box FullBox `json:"fullBox"`
  HandlerType string `json:"handlerType"`
  Name string `json:"name"`
}

type AVCcBox struct {
  Box Box `json:"box"`
  Version uint8 `json:"version"`
  Profile uint8 `json:"profile"`
  Level uint8 `json:"level"`
  SizeLen uint8 `json:"sizeLen"`
  SPS [][]byte `json:"sps"`
  PPS [][]byte `json:"pps"`
}

type PixelAspectRatioBox struct {
  Box Box `json:"box"`
  HSpacing uint32 `json:"hSpacing"`
  VSpacing uint32 `json:"vSpacing"`
}

type VideoSampleDescription struct {
  Width uint16 `json:"width"`
  Height uint16 `json:"height"`
  Depth uint16 `json:"depth"`
}

type ESDescriptor struct {
  Tag uint8 `json:"tag"`
  Length uint8 `json:"length"`
  Id uint16 `json:"id"`
  Config []byte `json:"config"`
}

type ElementaryStreamDescBox struct {
  Box Box `json:"box"`
  Esd ESDescriptor `json:"esd"`
}

type SoundSampleDescription struct {
  Channels uint16 `json:"channels"`
  SampleSize uint16 `json:"sampleSize"`
  SampleRate float32 `json:"sampleRate"`
}

type SampleEntry struct {
  Box Box `json:"box"`
  DataRefIndex uint16 `json:"dataRefIndex"`
  SampleDesc interface{} `json:"sampleDesc"`
  Extensions []interface{} `json:"extensions"`
}

type SampleDescriptionBox struct {
  Box FullBox `json:"box"`
  EntryCount uint32 `json:"entryCount"`
  Entries []SampleEntry `json:"entries"`
}

type TimeToSampleBox struct {
  Box FullBox `json:"fullBox"`
  EntryCount uint32 `json:"entryCount"`
  SampleCount []uint32 `json:"sampleCount"`
  SampleDelta []uint32 `json:"sampleDelta"`
}

type SyncSampleBox struct {
  Box FullBox `json:"fullBox"`
  EntryCount uint32 `json:"entryCount"`
  SampleNumber []uint32 `json:"sampleNumber"`
}

type CompTimeToSampleBox struct {
  Box FullBox `json:"fullBox"`
  EntryCount uint32 `json:"entryCount"`
  SampleCount []int32 `json:"sampleCount"`
  SampleOffset []int32 `json:"sampleOffset"`
}

type SampleToChunkBox struct {
  Box FullBox `json:"fullBox"`
  EntryCount uint32 `json:"entryCount"`
  FirstChunk []int32 `json:"firstChunk"`
  SamplesPerChunk []int32 `json:"samplesPerChunk"`
  SampleDescIndex []int32 `json:"sampleDescindex"`
}

type SampleSizeBox struct {
  Box FullBox `json:"fullBox"`
  SampleSize uint32 `json:"sampleSize"`
  SampleCount uint32 `json:"sampleCount"`
  EntrySize []uint32 `json:"entrySize"`
}

type ChunkOffsetBox struct {
  Box FullBox `json:"fullBox"`
  EntryCount uint32 `json:"entryCount"`
  ChunkOffset []uint32 `json:"chunkOffset"`
}

type SampleTableBox struct {
  Box Box `json:"box"`
  Stsd SampleDescriptionBox `json:"stsd"`
  Stts TimeToSampleBox `json:"stts"`
  Stss SyncSampleBox `json:"stss"`
  Ctss CompTimeToSampleBox `json:"ctss"`
  Stsc SampleToChunkBox `json:"stsc"`
  Stsz SampleSizeBox `json:"stsz"`
  Stco ChunkOffsetBox `json:"stco"`
}

type MediaInfoBox struct {
  Box Box `json:"box"`
  Stbl SampleTableBox `json:"stbl"`
}

type MediaBox struct {
  Box Box `json:"box"`
  Mdhd MediaHeaderBox `json:"mdhd"`
  Hdlr HandlerBox `json:"hdlr"`
  Minf MediaInfoBox `json:"minf"`
}

type TrackBox struct {
  Box Box `json:"box"`
  Tkhd TrackHeaderBox `json:"tkhd"`
  Mdia MediaBox `json:"mdia"`
}

type MovieBox struct {
  Box Box `json:"box"`
  Mvhd MovieHeaderBox `json:"mvhd"`
  Tracks []TrackBox `json:"tracks"`
}
