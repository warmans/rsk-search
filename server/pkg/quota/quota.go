package quota

// BandwidthQuotaInMiB doesn't seem correct. The sum actually gives us like 5,000,000Gib - but the usage is much lower in reality so it's not
// clear where the quota should be set.
const BandwidthQuotaInMiB float64 = 5 * 1024 * 100000 // ???

func BytesAsMib(numBytes int64) float64 {
	return float64(numBytes) / 1048576
}
