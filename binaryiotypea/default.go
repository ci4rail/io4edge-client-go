package binaryiotypea

const (
	// DefaultChannelFilterMask defines the channels which are considered in the stream by default
	DefaultChannelFilterMask uint32 = 0xFFFFFFFF
	//DefaultBucketSamples defines the max. number of samples per message
	DefaultBucketSamples uint32 = 20
	// DefaultKeepaliveInterval defines the default keep alive interval in ms
	DefaultKeepaliveInterval uint32 = 100000
	// DefaultBufferedSamples defines the default buffer size in samples
	DefaultBufferedSamples uint32 = 50
)
