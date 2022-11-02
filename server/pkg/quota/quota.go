package quota

const BandwidthQuotaInMiB int32 = 1024 * 1000 // 1000 Gib

func BytesAsMib(numBytes int64) int32 {
	return int32(numBytes / 1048576)
}
