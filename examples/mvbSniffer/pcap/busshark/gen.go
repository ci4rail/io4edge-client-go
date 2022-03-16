package busshark

func i64tobBe(val uint64) []byte {
	r := make([]byte, 8)
	for i := uint32(0); i < 8; i++ {
		r[i] = byte((val >> (8 * (7 - i))) & 0xff)
	}
	return r
}

func crc(b []byte) uint8 {
	return 0xCC // TODO
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func mvbDataWithCRCs(mvbData []byte) []byte {
	r := make([]byte, 0)
	for idx := 0; idx < len(mvbData); idx += 8 {
		sliceLen := min(len(mvbData)-idx, 8)
		curSlice := mvbData[idx : idx+sliceLen]
		r = append(r, curSlice...)
		r = append(r, crc(curSlice))
	}
	return r
}

// line: 0=undef 1=A 2=B
// frameType: 0=undef 1=master 2=slave
func statusField(line int, frameType int, len int) []byte {
	len += 2 // add 2 for start and end delimiter
	r := make([]byte, len)
	for i := 0; i < len; i++ {
		v := byte(0)
		if line == 2 {
			v |= 0x20
		}
		switch frameType {
		case 0:
			v |= 0x03 // Unknown
		case 1:
			v |= 0x01 // master
		case 2:
			v |= 0x02 // slave
		}
		if i == len-1 {
			v |= 0x7 // frame end
		}
		r[i] = v
	}
	return r
}

// Pkt generates the bytes for a MVB busshark packet
// line: 0=undef 1=A 2=B
// frameType: 0=undef 1=master 2=slave
// mvbData: Frame data w/o CRC
func Pkt(frameNumber uint64, receiveTime20nsUnits uint64, line int, frameType int, mvbData []byte) []byte {
	r := make([]byte, 0)
	// ethernet header
	r = append(r, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff)
	r = append(r, 0xc0, 0x01, 0xc0, 0xca, 0xc0, 0x1a)
	r = append(r, 0xc0, 0x01)

	// fill
	r = append(r, "MVBSNIFFER-Ci4Rail-Format 1.00"...)

	// insert dummy CRCs, they have been removed by MVB sniffer
	mvbDataInclCRCs := mvbDataWithCRCs(mvbData)
	r = append(r, i64tobBe(frameNumber)...)
	r = append(r, i64tobBe(receiveTime20nsUnits)...)
	r = append(r, statusField(line, frameType, len(mvbDataInclCRCs))...)
	r = append(r, 0x00) // start delimiter
	r = append(r, mvbDataInclCRCs...)
	r = append(r, 0x00) // end delimiter
	return r
}
