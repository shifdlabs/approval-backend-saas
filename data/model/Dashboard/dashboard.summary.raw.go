package dashboard

// DashboardSummaryRaw adalah data mentah dari repository sebelum diproses service.
type DashboardSummaryRaw struct {
	// NeedApproval: giliran user di sequence (status=1, step match)
	NeedApprovalTotal  int
	NeedApprovalUrgent int // priority = 1 (High)
	NeedApprovalNormal int // priority != 1
	OldestPendingDays  int // dihitung dari seluruh data, tidak ikut filter period

	// InProgress: status=1, user sebagai author ATAU approver
	InProgressTotal       int
	LongestProcessingDays int // dihitung dari seluruh data, tidak ikut filter period

	// Rejected: status=99
	RejectedTotal     int
	MineNeedsRevision int // hanya milik user yang login

	// Completed: status=2 atau 3
	CompletedTotal     int
	CompletedTotalYear int // selalu aggregate per tahun berjalan, tidak ikut filter
}
