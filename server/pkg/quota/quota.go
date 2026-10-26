package quota

const BandwidthQuotaInMiB float64 = 2000 * 1024 // 2000 Gib

func BytesAsMib(numBytes int64) float64 {
	return float64(numBytes) / 1048576
}
