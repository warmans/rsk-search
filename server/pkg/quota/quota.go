package quota

const BandwidthQuotaInMiB int64 = 1024 * 1000 // 1000 Gib

func BytesAsMib(numBytes int64) int64 {
	return numBytes / 1048576
}
