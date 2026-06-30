package metricsdomain

import (
	"context"
	"math"

	"github.com/jmoiron/sqlx"
	"github.com/yomiAdenaike01/support-queue/internals/db/queries"
)

type Repository struct {
	db *sqlx.DB
}

type DbDashboardSummaryResponse struct {
	TotalTicketsToday     int `db:"total_tickets_today"`
	TotalTicketsYday      int `db:"total_tickets_yday"`
	PendingTicketsToday   int `db:"pending_tickets_today"`
	PendingTicketsYday    int `db:"pending_tickets_yday"`
	ProcessingTickets     int `db:"processing_tickets_today"`
	ProcessingTicketsYday int `db:"processing_tickets_yday"`
	ProcessedTickets      int `db:"processed_tickets_today"`
	ProcessedTicketsYday  int `db:"processed_tickets_yday"`
	FailedTickets         int `db:"failed_tickets_today"`
	FailedTicketsYday     int `db:"failed_tickets_yday"`
	ResolvedTickets       int `db:"resolved_tickets_today"`
	ResolvedTicketsYday   int `db:"resolved_tickets_yday"`
}

type DailySummary struct {
	Today int `json:"today"`
	Yday  int `json:"yesterday"`
	Diff  int `json:"percentageDifference"`
}

type DashboardSummaryResponse struct {
	TotalTickets      DailySummary `json:"totalTickets"`
	PendingTickets    DailySummary `json:"pendingTickets"`
	ProcessingTickets DailySummary `json:"processingTickets"`
	ProcessedTickets  DailySummary `json:"processedTickets"`
	FailedTickets     DailySummary `json:"failedTickets"`
	ResolvedTickets   DailySummary `json:"resolvedTickets"`
}

func (r *Repository) GetDashboardSummary() (DashboardSummaryResponse, error) {
	var dbSummary DbDashboardSummaryResponse

	if err := r.db.Get(&dbSummary, queries.DASHBOARD_SUMMARY); err != nil {
		return DashboardSummaryResponse{}, err
	}

	return dbSummary.ToDashboardSummary(), nil
}

func (s DbDashboardSummaryResponse) ToDashboardSummary() DashboardSummaryResponse {
	return DashboardSummaryResponse{
		TotalTickets:      newDailySumary(s.TotalTicketsToday, s.TotalTicketsYday),
		PendingTickets:    newDailySumary(s.PendingTicketsToday, s.PendingTicketsYday),
		ProcessingTickets: newDailySumary(s.ProcessingTickets, s.ProcessingTicketsYday),
		ProcessedTickets:  newDailySumary(s.ProcessedTickets, s.ProcessedTicketsYday),
		FailedTickets:     newDailySumary(s.FailedTickets, s.FailedTicketsYday),
		ResolvedTickets:   newDailySumary(s.ResolvedTickets, s.ResolvedTicketsYday),
	}
}

func newDailySumary(today, yday int) DailySummary {
	return DailySummary{
		Today: today,
		Yday:  yday,
		Diff:  percentageDiff(today, yday),
	}
}

func percentageDiff(today, yday int) int {
	if yday == 0 {
		if today == 0 {
			return 0
		}
		return 100
	}

	return int(math.Round(float64(today-yday) / float64(yday) * 100))
}

type TicketsOvertimeInput struct {
	StartDatetime string `form:"start_datetime" binding:"required"`
	EndDatetime   string `form:"end_datetime" binding:"required"`
}

type TicketsOvertime struct {
	Hour    int `json:"hour" db:"hour"`
	Tickets int `json:"tickets" db:"tickets"`
}

func (r *Repository) GetTicketsOvertime(ctx context.Context, input TicketsOvertimeInput) ([]TicketsOvertime, error) {
	var output []TicketsOvertime
	if err := r.db.SelectContext(ctx, &output, queries.TICKETS_OVER_TIME, input.StartDatetime, input.EndDatetime); err != nil {
		return nil, err
	}
	return output, nil
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}
