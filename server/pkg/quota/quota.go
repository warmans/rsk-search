package quota

const BandwidthQuotaInMiB int64 = 1024 * 2500 // 2.5 Gib

func BytesAsMib(numBytes int64) int64 {
	return numBytes / 1048576
}
