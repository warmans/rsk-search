package quota

const BandwidthQuotaInMiB int64 = 1024 * 2000 // 2000 Gib

func BytesAsMib(numBytes int64) int64 {
	return numBytes / 1048576
}
