package pcap

import (
	"io"
)

// Writer represent the pcap writer
type Writer struct {
	llWriter io.Writer
}

// NewWriter generates an instance of the pcap writer
func NewWriter(llWriter io.Writer) (*Writer, error) {
	w := &Writer{llWriter: llWriter}

	_, err := w.llWriter.Write(globalHeader())
	if err != nil {
		return nil, err
	}

	return w, nil
}

// AddPacket adds a packet to the pcap file
func (w *Writer) AddPacket(ts uint64, data []byte) error {
	// write header
	len := uint32(len(data))
	_, err := w.llWriter.Write(recordHeader(ts, len, len))
	if err != nil {
		return err
	}
	// write payload
	_, err = w.llWriter.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func i32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

func i16tob(val uint16) []byte {
	r := make([]byte, 2)
	for i := uint32(0); i < 2; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

// globalHeader provides a byte slice with the global header data
// https://wiki.wireshark.org/Development/LibpcapFileFormat
//
//	typedef struct pcap_hdr_s {
//		guint32 magic_number;   /* magic number */
//		guint16 version_major;  /* major version number */
//		guint16 version_minor;  /* minor version number */
//		gint32  thiszone;       /* GMT to local correction */
//		guint32 sigfigs;        /* accuracy of timestamps */
//		guint32 snaplen;        /* max length of captured packets, in octets */
//		guint32 network;        /* data link type */
//	} pcap_hdr_t;
func globalHeader() []byte {
	r := make([]byte, 0)
	r = append(r, i32tob(0xa1b2c3d4)...)
	r = append(r, i16tob(0x0002)...)
	r = append(r, i16tob(0x0004)...)
	r = append(r, i32tob(0x00000000)...) // thiszone
	r = append(r, i32tob(0x00000000)...) // sigfigs
	r = append(r, i32tob(0x00040000)...) // max length
	r = append(r, i32tob(0x00000001)...) // data link type = Ethernet
	return r
}

// recordHeader provides a byte slice with the record header data
//
//	typedef struct pcaprec_hdr_s {
//		guint32 ts_sec;         /* timestamp seconds */
//		guint32 ts_usec;        /* timestamp microseconds */
//		guint32 incl_len;       /* number of octets of packet saved in file */
//		guint32 orig_len;       /* actual length of packet */
//	} pcaprec_hdr_t;
func recordHeader(ts uint64, inclLen uint32, origLen uint32) []byte {
	r := make([]byte, 0)
	tsSec := uint32(ts / 1000000)
	tsUsec := uint32(ts % 1000000)
	r = append(r, i32tob(tsSec)...)
	r = append(r, i32tob(tsUsec)...)
	r = append(r, i32tob(inclLen)...)
	r = append(r, i32tob(origLen)...)
	return r
}
