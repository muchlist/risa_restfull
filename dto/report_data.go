package dto

type ReportResponse struct {
	TargetTime     int64
	CctvDaily      *VendorCheck   `json:"cctv_daily"`
	CctvMonthly    *VenPhyCheck   `json:"cctv_monthly"`
	CctvQuarterly  *VenPhyCheck   `json:"cctv_quarterly"`
	AltaiDaily     *AltaiCheck    `json:"altai_daily"`
	AltaiMonthly   *AltaiPhyCheck `json:"altai_monthly"`
	AltaiQuarterly *AltaiPhyCheck `json:"altai_quarterly"`
}
