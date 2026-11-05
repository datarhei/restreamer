package hevcparser

// based on https://github.com/deepch/vdk/blob/v0.0.21/codec/h265parser/parser.go

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/utils/bits"
	"github.com/datarhei/joy4/utils/bits/pio"
)

type SPSInfo struct {
	ProfileIdc                       uint
	LevelIdc                         uint
	CropLeft                         uint
	CropRight                        uint
	CropTop                          uint
	CropBottom                       uint
	Width                            uint
	Height                           uint
	PicWidthInLumaSamples            uint
	PicHeightInLumaSamples           uint
	generalProfileSpace              uint
	generalTierFlag                  uint
	generalProfileIDC                uint
	generalProfileCompatibilityFlags uint32
	generalConstraintIndicatorFlags  uint64
	generalLevelIDC                  uint
}

const (
	NAL_UNIT_CODED_SLICE_TRAIL_N    = 0
	NAL_UNIT_CODED_SLICE_TRAIL_R    = 1
	NAL_UNIT_CODED_SLICE_TSA_N      = 2
	NAL_UNIT_CODED_SLICE_TSA_R      = 3
	NAL_UNIT_CODED_SLICE_STSA_N     = 4
	NAL_UNIT_CODED_SLICE_STSA_R     = 5
	NAL_UNIT_CODED_SLICE_RADL_N     = 6
	NAL_UNIT_CODED_SLICE_RADL_R     = 7
	NAL_UNIT_CODED_SLICE_RASL_N     = 8
	NAL_UNIT_CODED_SLICE_RASL_R     = 9
	NAL_UNIT_RESERVED_VCL_N10       = 10
	NAL_UNIT_RESERVED_VCL_R11       = 11
	NAL_UNIT_RESERVED_VCL_N12       = 12
	NAL_UNIT_RESERVED_VCL_R13       = 13
	NAL_UNIT_RESERVED_VCL_N14       = 14
	NAL_UNIT_RESERVED_VCL_R15       = 15
	NAL_UNIT_CODED_SLICE_BLA_W_LP   = 16
	NAL_UNIT_CODED_SLICE_BLA_W_RADL = 17
	NAL_UNIT_CODED_SLICE_BLA_N_LP   = 18
	NAL_UNIT_CODED_SLICE_IDR_W_RADL = 19
	NAL_UNIT_CODED_SLICE_IDR_N_LP   = 20
	NAL_UNIT_CODED_SLICE_CRA        = 21
	NAL_UNIT_RESERVED_IRAP_VCL22    = 22
	NAL_UNIT_RESERVED_IRAP_VCL23    = 23
	NAL_UNIT_RESERVED_VCL24         = 24
	NAL_UNIT_RESERVED_VCL25         = 25
	NAL_UNIT_RESERVED_VCL26         = 26
	NAL_UNIT_RESERVED_VCL27         = 27
	NAL_UNIT_RESERVED_VCL28         = 28
	NAL_UNIT_RESERVED_VCL29         = 29
	NAL_UNIT_RESERVED_VCL30         = 30
	NAL_UNIT_RESERVED_VCL31         = 31
	NAL_UNIT_VPS                    = 32
	NAL_UNIT_SPS                    = 33
	NAL_UNIT_PPS                    = 34
	NAL_UNIT_ACCESS_UNIT_DELIMITER  = 35
	NAL_UNIT_EOS                    = 36
	NAL_UNIT_EOB                    = 37
	NAL_UNIT_FILLER_DATA            = 38
	NAL_UNIT_PREFIX_SEI             = 39
	NAL_UNIT_SUFFIX_SEI             = 40
	NAL_UNIT_RESERVED_NVCL41        = 41
	NAL_UNIT_RESERVED_NVCL42        = 42
	NAL_UNIT_RESERVED_NVCL43        = 43
	NAL_UNIT_RESERVED_NVCL44        = 44
	NAL_UNIT_RESERVED_NVCL45        = 45
	NAL_UNIT_RESERVED_NVCL46        = 46
	NAL_UNIT_RESERVED_NVCL47        = 47
	NAL_UNIT_UNSPECIFIED_48         = 48
	NAL_UNIT_UNSPECIFIED_49         = 49
	NAL_UNIT_UNSPECIFIED_50         = 50
	NAL_UNIT_UNSPECIFIED_51         = 51
	NAL_UNIT_UNSPECIFIED_52         = 52
	NAL_UNIT_UNSPECIFIED_53         = 53
	NAL_UNIT_UNSPECIFIED_54         = 54
	NAL_UNIT_UNSPECIFIED_55         = 55
	NAL_UNIT_UNSPECIFIED_56         = 56
	NAL_UNIT_UNSPECIFIED_57         = 57
	NAL_UNIT_UNSPECIFIED_58         = 58
	NAL_UNIT_UNSPECIFIED_59         = 59
	NAL_UNIT_UNSPECIFIED_60         = 60
	NAL_UNIT_UNSPECIFIED_61         = 61
	NAL_UNIT_UNSPECIFIED_62         = 62
	NAL_UNIT_UNSPECIFIED_63         = 63
	NAL_UNIT_INVALID                = 64
)

const (
	MAX_VPS_COUNT  = 16
	MAX_SUB_LAYERS = 7
	MAX_SPS_COUNT  = 32
)

var (
	ErrorHEVCIncorectUnitSize = errors.New("incorrect unit size")
	ErrorHECVIncorectUnitType = errors.New("incorrect unit type")
)

var StartCodeBytes = []byte{0, 0, 1}
var AUDBytes = []byte{0, 0, 0, 1, 0x9, 0xf0, 0, 0, 0, 1} // AUD

const (
	NALU_RAW = iota
	NALU_AVCC
	NALU_ANNEXB
)

func SplitNALUs(b []byte) (nalus [][]byte, typ int) {
	if len(b) < 4 {
		return [][]byte{b}, NALU_RAW
	}
	val3 := pio.U24BE(b)
	val4 := pio.U32BE(b)
	if val4 <= uint32(len(b)) {
		_val4 := val4
		_b := b[4:]
		nalus := [][]byte{}
		for {
			nalus = append(nalus, _b[:_val4])
			_b = _b[_val4:]
			if len(_b) < 4 {
				break
			}
			_val4 = pio.U32BE(_b)
			_b = _b[4:]
			if _val4 > uint32(len(_b)) {
				break
			}
		}
		if len(_b) == 0 {
			return nalus, NALU_AVCC
		}
	}
	if val3 == 1 || val4 == 1 {
		_val3 := val3
		_val4 := val4
		start := 0
		pos := 0
		for {
			if start != pos {
				nalus = append(nalus, b[start:pos])
			}
			if _val3 == 1 {
				pos += 3
			} else if _val4 == 1 {
				pos += 4
			}
			start = pos
			if start == len(b) {
				break
			}
			_val3 = 0
			_val4 = 0
			for pos < len(b) {
				if pos+2 < len(b) && b[pos] == 0 {
					_val3 = pio.U24BE(b[pos:])
					if _val3 == 0 {
						if pos+3 < len(b) {
							_val4 = uint32(b[pos+3])
							if _val4 == 1 {
								break
							}
						}
					} else if _val3 == 1 {
						break
					}
					pos++
				} else {
					pos++
				}
			}
		}
		typ = NALU_ANNEXB
		return
	}

	return [][]byte{b}, NALU_RAW
}

func ParseSPS(sps []byte) (ctx SPSInfo, err error) {
	if len(sps) < 2 {
		err = ErrorHEVCIncorectUnitSize
		return
	}
	rbsp := nal2rbsp(sps[2:])
	br := &bits.GolombBitReader{R: bytes.NewReader(rbsp)}

	// sps_video_parameter_set_id
	if _, err = br.ReadBits(4); err != nil {
		return
	}

	// sps_max_sub_layers_minus1
	spsMaxSubLayersMinus1, err := br.ReadBits(3)
	if err != nil {
		return
	}

	// sps_temporal_id_nesting_flag
	if _, err = br.ReadBit(); err != nil {
		return
	}

	// profile_tier_level( 1, sps_max_sub_layers_minus1 )
	if err = parsePTL(br, &ctx, spsMaxSubLayersMinus1); err != nil {
		return
	}

	// sps_seq_parameter_set_id
	if _, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}

	// chroma_format_idc
	var chroma_format_idc uint
	if chroma_format_idc, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}
	if chroma_format_idc == 3 {
		// separate_colour_plane_flag
		if _, err = br.ReadBit(); err != nil {
			return
		}
	}

	// Table 6-1, Section 6.2
	var subWidthC uint
	var subHeightC uint

	switch chroma_format_idc {
	case 0:
		subWidthC, subHeightC = 1, 1
	case 1:
		subWidthC, subHeightC = 2, 2
	case 2:
		subWidthC, subHeightC = 2, 1
	case 3:
		subWidthC, subHeightC = 1, 1
	}

	// pic_width_in_luma_samples
	if ctx.PicWidthInLumaSamples, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}

	// pic_height_in_luma_samples
	if ctx.PicHeightInLumaSamples, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}

	// conformance_window_flag
	conformanceWindowFlag, err := br.ReadBit()
	if err != nil {
		return
	}

	var conf_win_left_offset uint
	var conf_win_right_offset uint
	var conf_win_top_offset uint
	var conf_win_bottom_offset uint

	if conformanceWindowFlag != 0 {
		// conf_win_left_offset
		conf_win_left_offset, err = br.ReadExponentialGolombCode()
		if err != nil {
			return
		}
		ctx.CropLeft = subWidthC * conf_win_left_offset

		// conf_win_right_offset
		conf_win_right_offset, err = br.ReadExponentialGolombCode()
		if err != nil {
			return
		}
		ctx.CropRight = subWidthC * conf_win_right_offset

		// conf_win_top_offset
		conf_win_top_offset, err = br.ReadExponentialGolombCode()
		if err != nil {
			return
		}
		ctx.CropTop = subHeightC * conf_win_top_offset

		// conf_win_bottom_offset
		conf_win_bottom_offset, err = br.ReadExponentialGolombCode()
		if err != nil {
			return
		}
		ctx.CropBottom = subHeightC * conf_win_bottom_offset
	}

	ctx.Width = ctx.PicWidthInLumaSamples - ctx.CropLeft - ctx.CropRight
	ctx.Height = ctx.PicHeightInLumaSamples - ctx.CropTop - ctx.CropBottom

	// bit_depth_luma_minus8
	if _, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}

	// bit_depth_chroma_minus8
	if _, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}

	// log2_max_pic_order_cnt_lsb_minus4
	_, err = br.ReadExponentialGolombCode()
	if err != nil {
		return
	}

	// sps_sub_layer_ordering_info_present_flag
	spsSubLayerOrderingInfoPresentFlag, err := br.ReadBit()
	if err != nil {
		return
	}
	var i uint
	if spsSubLayerOrderingInfoPresentFlag != 0 {
		i = 0
	} else {
		i = spsMaxSubLayersMinus1
	}
	for ; i <= spsMaxSubLayersMinus1; i++ {
		// sps_max_dec_pic_buffering_minus1[ i ]
		if _, err = br.ReadExponentialGolombCode(); err != nil {
			return
		}
		// sps_max_num_reorder_pics[ i ]
		if _, err = br.ReadExponentialGolombCode(); err != nil {
			return
		}
		// sps_max_latency_increase_plus1[ i ]
		if _, err = br.ReadExponentialGolombCode(); err != nil {
			return
		}
	}

	// log2_min_luma_coding_block_size_minus3
	if _, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}
	// log2_diff_max_min_luma_coding_block_size
	if _, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}
	// log2_min_luma_transform_block_size_minus2
	if _, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}
	// log2_diff_max_min_luma_transform_block_size
	if _, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}
	// max_transform_hierarchy_depth_inter
	if _, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}
	// max_transform_hierarchy_depth_intra
	if _, err = br.ReadExponentialGolombCode(); err != nil {
		return
	}
	return
}

func parsePTL(br *bits.GolombBitReader, ctx *SPSInfo, maxSubLayersMinus1 uint) error {
	var err error
	var ptl SPSInfo
	if ptl.generalProfileSpace, err = br.ReadBits(2); err != nil {
		return err
	}
	if ptl.generalTierFlag, err = br.ReadBit(); err != nil {
		return err
	}
	if ptl.generalProfileIDC, err = br.ReadBits(5); err != nil {
		return err
	}
	if ptl.generalProfileCompatibilityFlags, err = br.ReadBits32(32); err != nil {
		return err
	}
	if ptl.generalConstraintIndicatorFlags, err = br.ReadBits64(48); err != nil {
		return err
	}
	if ptl.generalLevelIDC, err = br.ReadBits(8); err != nil {
		return err
	}
	updatePTL(ctx, &ptl)
	if maxSubLayersMinus1 == 0 {
		return nil
	}
	subLayerProfilePresentFlag := make([]uint, maxSubLayersMinus1)
	subLayerLevelPresentFlag := make([]uint, maxSubLayersMinus1)
	for i := uint(0); i < maxSubLayersMinus1; i++ {
		if subLayerProfilePresentFlag[i], err = br.ReadBit(); err != nil {
			return err
		}
		if subLayerLevelPresentFlag[i], err = br.ReadBit(); err != nil {
			return err
		}
	}
	if maxSubLayersMinus1 > 0 {
		for i := maxSubLayersMinus1; i < 8; i++ {
			if _, err = br.ReadBits(2); err != nil {
				return err
			}
		}
	}
	for i := uint(0); i < maxSubLayersMinus1; i++ {
		if subLayerProfilePresentFlag[i] != 0 {
			if _, err = br.ReadBits32(32); err != nil {
				return err
			}
			if _, err = br.ReadBits32(32); err != nil {
				return err
			}
			if _, err = br.ReadBits32(24); err != nil {
				return err
			}
		}

		if subLayerLevelPresentFlag[i] != 0 {
			if _, err = br.ReadBits(8); err != nil {
				return err
			}
		}
	}
	return nil
}

func updatePTL(ctx, ptl *SPSInfo) {
	ctx.generalProfileSpace = ptl.generalProfileSpace

	if ptl.generalTierFlag > ctx.generalTierFlag {
		ctx.generalLevelIDC = ptl.generalLevelIDC

		ctx.generalTierFlag = ptl.generalTierFlag
	} else {
		if ptl.generalLevelIDC > ctx.generalLevelIDC {
			ctx.generalLevelIDC = ptl.generalLevelIDC
		}
	}

	if ptl.generalProfileIDC > ctx.generalProfileIDC {
		ctx.generalProfileIDC = ptl.generalProfileIDC
	}

	ctx.generalProfileCompatibilityFlags &= ptl.generalProfileCompatibilityFlags

	ctx.generalConstraintIndicatorFlags &= ptl.generalConstraintIndicatorFlags
}

func nal2rbsp(nal []byte) []byte {
	return bytes.Replace(nal, []byte{0x0, 0x0, 0x3}, []byte{0x0, 0x0}, -1)
}

type CodecData struct {
	Record     []byte
	RecordInfo HEVCDecoderConfRecord
	SPSInfo    SPSInfo
}

func (codec CodecData) Type() av.CodecType {
	return av.HEVC
}

func (codec CodecData) HEVCDecoderConfRecordBytes() []byte {
	return codec.Record
}

func (codec CodecData) SPS() []byte {
	return codec.RecordInfo.SPS[0]
}

func (codec CodecData) PPS() []byte {
	return codec.RecordInfo.PPS[0]
}

func (codec CodecData) VPS() []byte {
	return codec.RecordInfo.VPS[0]
}

func (codec CodecData) Width() int {
	return int(codec.SPSInfo.Width)
}

func (codec CodecData) Height() int {
	return int(codec.SPSInfo.Height)
}

func NewCodecDataFromHEVCDecoderConfRecord(record []byte) (self CodecData, err error) {
	self.Record = record
	if _, err = (&self.RecordInfo).Unmarshal(record); err != nil {
		return
	}
	if len(self.RecordInfo.SPS) == 0 {
		err = fmt.Errorf("hevcparser: no SPS found in HEVCDecoderConfRecord")
		return
	}
	if len(self.RecordInfo.PPS) == 0 {
		err = fmt.Errorf("hevcparser: no PPS found in HEVCDecoderConfRecord")
		return
	}
	if len(self.RecordInfo.VPS) == 0 {
		err = fmt.Errorf("hevcparser: no VPS found in HEVCDecoderConfRecord")
		return
	}
	if self.SPSInfo, err = ParseSPS(self.RecordInfo.SPS[0]); err != nil {
		err = fmt.Errorf("hevcparser: parse SPS failed(%s)", err)
		return
	}
	return
}

func NewCodecDataFromVPSAndSPSAndPPS(vps, sps, pps []byte) (self CodecData, err error) {
	recordinfo := HEVCDecoderConfRecord{}
	recordinfo.HEVCProfileIndication = sps[3]
	recordinfo.ProfileCompatibility = sps[4]
	recordinfo.HEVCLevelIndication = sps[5]
	recordinfo.SPS = [][]byte{sps}
	recordinfo.PPS = [][]byte{pps}
	recordinfo.VPS = [][]byte{vps}
	recordinfo.LengthSizeMinusOne = 3
	if self.SPSInfo, err = ParseSPS(sps); err != nil {
		return
	}
	buf := make([]byte, recordinfo.Len())
	recordinfo.Marshal(buf, self.SPSInfo)
	self.RecordInfo = recordinfo
	self.Record = buf
	return
}

type HEVCDecoderConfRecord struct {
	HEVCProfileIndication uint8
	ProfileCompatibility  uint8
	HEVCLevelIndication   uint8
	LengthSizeMinusOne    uint8
	VPS                   [][]byte
	SPS                   [][]byte
	PPS                   [][]byte
}

var ErrDecconfInvalid = fmt.Errorf("hevcparser: HEVCDecoderConfRecord invalid")

func (record *HEVCDecoderConfRecord) Unmarshal(b []byte) (n int, err error) {
	if len(b) < 30 {
		err = ErrDecconfInvalid
		return
	}
	record.HEVCProfileIndication = b[1]
	record.ProfileCompatibility = b[2]
	record.HEVCLevelIndication = b[3]
	record.LengthSizeMinusOne = b[4] & 0x03

	vpscount := int(b[25] & 0x1f)
	n += 26
	for i := 0; i < vpscount; i++ {
		if len(b) < n+2 {
			err = ErrDecconfInvalid
			return
		}
		vpslen := int(pio.U16BE(b[n:]))
		n += 2

		if len(b) < n+vpslen {
			err = ErrDecconfInvalid
			return
		}
		record.VPS = append(record.VPS, b[n:n+vpslen])
		n += vpslen
	}

	if len(b) < n+1 {
		err = ErrDecconfInvalid
		return
	}

	n++
	n++

	spscount := int(b[n])
	n++

	for i := 0; i < spscount; i++ {
		if len(b) < n+2 {
			err = ErrDecconfInvalid
			return
		}
		spslen := int(pio.U16BE(b[n:]))
		n += 2

		if len(b) < n+spslen {
			err = ErrDecconfInvalid
			return
		}
		record.SPS = append(record.SPS, b[n:n+spslen])
		n += spslen
	}

	n++
	n++

	ppscount := int(b[n])
	n++

	for i := 0; i < ppscount; i++ {
		if len(b) < n+2 {
			err = ErrDecconfInvalid
			return
		}
		ppslen := int(pio.U16BE(b[n:]))
		n += 2

		if len(b) < n+ppslen {
			err = ErrDecconfInvalid
			return
		}
		record.PPS = append(record.PPS, b[n:n+ppslen])
		n += ppslen
	}
	return
}

func (record HEVCDecoderConfRecord) Len() (n int) {
	n = 23
	for _, sps := range record.SPS {
		n += 5 + len(sps)
	}
	for _, pps := range record.PPS {
		n += 5 + len(pps)
	}
	for _, vps := range record.VPS {
		n += 5 + len(vps)
	}
	return
}

func (record HEVCDecoderConfRecord) Marshal(b []byte, si SPSInfo) (n int) {
	b[0] = 1
	b[1] = record.HEVCProfileIndication
	b[2] = record.ProfileCompatibility
	b[3] = record.HEVCLevelIndication
	b[21] = 3
	b[22] = 3
	n += 23
	b[n] = (record.VPS[0][0] >> 1) & 0x3f
	n++
	b[n] = byte(len(record.VPS) >> 8)
	n++
	b[n] = byte(len(record.VPS))
	n++
	for _, vps := range record.VPS {
		pio.PutU16BE(b[n:], uint16(len(vps)))
		n += 2
		copy(b[n:], vps)
		n += len(vps)
	}
	b[n] = (record.SPS[0][0] >> 1) & 0x3f
	n++
	b[n] = byte(len(record.SPS) >> 8)
	n++
	b[n] = byte(len(record.SPS))
	n++
	for _, sps := range record.SPS {
		pio.PutU16BE(b[n:], uint16(len(sps)))
		n += 2
		copy(b[n:], sps)
		n += len(sps)
	}
	b[n] = (record.PPS[0][0] >> 1) & 0x3f
	n++
	b[n] = byte(len(record.PPS) >> 8)
	n++
	b[n] = byte(len(record.PPS))
	n++
	for _, pps := range record.PPS {
		pio.PutU16BE(b[n:], uint16(len(pps)))
		n += 2
		copy(b[n:], pps)
		n += len(pps)
	}
	return
}
