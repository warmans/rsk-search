package quota

const BandwidthQuotaInMiB float64 = 1000 * 1024 // 1000 Gib

func BytesAsMib(numBytes int64) float64 {
	return float64(numBytes) / 1048576
}
