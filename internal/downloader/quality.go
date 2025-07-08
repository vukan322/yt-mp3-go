package downloader

type AudioQuality string

const (
	QualityLow    AudioQuality = "low"
	QualityMedium AudioQuality = "medium"
	QualityHigh   AudioQuality = "high"
)

func (q AudioQuality) ToYtDlp() string {
	switch q {
	case QualityLow:
		return "9"
	case QualityMedium:
		return "5"
	case QualityHigh:
		return "0"
	default:
		return "0"
	}
}
