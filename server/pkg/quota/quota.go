package quota

const BandwidthQuotaInMiB float64 = 1024 * 100000 // 100Gib

func BytesAsMib(numBytes int64) float64 {
	return float64(numBytes) / 1048576
}
