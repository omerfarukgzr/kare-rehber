package service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

type SummaryReport struct {
	ActiveStudents       int64 `json:"activeStudents"`
	ActiveCoaches        int64 `json:"activeCoaches"`
	ActiveCoordinators   int64 `json:"activeCoordinators"`
	ActiveAssignments    int64 `json:"activeAssignments"`
	TotalEvaluations     int64 `json:"totalEvaluations"`
	PendingEvaluations   int64 `json:"pendingEvaluations"`
	ApprovedEvaluations  int64 `json:"approvedEvaluations"`
	OpenThreads          int64 `json:"openThreads"`
	PendingRegistrations int64 `json:"pendingRegistrations"`
}

func (s *ReportService) Summary() (*SummaryReport, error) {
	r := &SummaryReport{}
	s.db.Raw(`SELECT COUNT(*) FROM users WHERE role='student' AND is_active=true`).Scan(&r.ActiveStudents)
	s.db.Raw(`SELECT COUNT(*) FROM users WHERE role='coach' AND is_active=true`).Scan(&r.ActiveCoaches)
	s.db.Raw(`SELECT COUNT(*) FROM users WHERE role='coordinator' AND is_active=true`).Scan(&r.ActiveCoordinators)
	s.db.Raw(`SELECT COUNT(*) FROM coach_student_assignments WHERE is_active=true`).Scan(&r.ActiveAssignments)
	s.db.Raw(`SELECT COUNT(*) FROM evaluations`).Scan(&r.TotalEvaluations)
	s.db.Raw(`SELECT COUNT(*) FROM evaluations WHERE status='pending'`).Scan(&r.PendingEvaluations)
	s.db.Raw(`SELECT COUNT(*) FROM evaluations WHERE status='approved'`).Scan(&r.ApprovedEvaluations)
	s.db.Raw(`SELECT COUNT(*) FROM message_threads WHERE status='open'`).Scan(&r.OpenThreads)
	s.db.Raw(`SELECT COUNT(*) FROM registrations WHERE status='pending'`).Scan(&r.PendingRegistrations)
	return r, nil
}

type CoachPerformanceRow struct {
	CoachID         string  `json:"coachId" gorm:"column:coach_id"`
	CoachName       string  `json:"coachName" gorm:"column:coach_name"`
	StudentCount    int     `json:"studentCount" gorm:"column:student_count"`
	EvaluationCount int     `json:"evaluationCount" gorm:"column:evaluation_count"`
	ApprovedCount   int     `json:"approvedCount" gorm:"column:approved_count"`
	AvgMotivation   float64 `json:"avgMotivation" gorm:"column:avg_motivation"`
	AvgBehavior     float64 `json:"avgBehavior" gorm:"column:avg_behavior"`
}

func (s *ReportService) CoachPerformance() ([]CoachPerformanceRow, error) {
	var rows []CoachPerformanceRow
	q := `
		SELECT
		  u.id AS coach_id,
		  u.full_name AS coach_name,
		  COUNT(DISTINCT a.student_id) FILTER (WHERE a.is_active = true) AS student_count,
		  COUNT(e.id) AS evaluation_count,
		  COUNT(e.id) FILTER (WHERE e.status = 'approved') AS approved_count,
		  COALESCE(AVG(e.motivation), 0) AS avg_motivation,
		  COALESCE(AVG(e.behavior), 0) AS avg_behavior
		FROM users u
		LEFT JOIN coach_student_assignments a ON a.coach_id = u.id
		LEFT JOIN evaluations e ON e.assignment_id = a.id
		WHERE u.role = 'coach'
		GROUP BY u.id, u.full_name
		ORDER BY evaluation_count DESC, u.full_name`
	err := s.db.Raw(q).Scan(&rows).Error
	return rows, err
}

type CityDistributionRow struct {
	City    *string `json:"city" gorm:"column:city"`
	Count   int     `json:"count" gorm:"column:count"`
}

func (s *ReportService) CityDistribution(role string) ([]CityDistributionRow, error) {
	var rows []CityDistributionRow
	q := `SELECT city, COUNT(*) AS count FROM users WHERE role = ? AND is_active = true GROUP BY city ORDER BY count DESC`
	err := s.db.Raw(q, role).Scan(&rows).Error
	return rows, err
}

type WeekStatsRow struct {
	WeekID            string `json:"weekId" gorm:"column:week_id"`
	WeekNo            int    `json:"weekNo" gorm:"column:week_no"`
	Label             string `json:"label" gorm:"column:label"`
	StartDate         string `json:"startDate" gorm:"column:start_date"`
	TotalEvaluations  int    `json:"totalEvaluations" gorm:"column:total_evaluations"`
	ApprovedCount     int    `json:"approvedCount" gorm:"column:approved_count"`
	ActiveAssignments int    `json:"activeAssignments" gorm:"column:active_assignments"`
}

func (s *ReportService) WeekStats() ([]WeekStatsRow, error) {
	var rows []WeekStatsRow
	q := `
		SELECT
		  w.id::text AS week_id, w.week_no, w.label, w.start_date::text,
		  COUNT(e.id) AS total_evaluations,
		  COUNT(e.id) FILTER (WHERE e.status = 'approved') AS approved_count,
		  (SELECT COUNT(*) FROM coach_student_assignments WHERE is_active = true) AS active_assignments
		FROM evaluation_weeks w
		LEFT JOIN evaluations e ON e.evaluation_week_id = w.id
		GROUP BY w.id
		ORDER BY w.week_no`
	err := s.db.Raw(q).Scan(&rows).Error
	return rows, err
}

type StudentTrendRow struct {
	WeekNo     int     `json:"weekNo" gorm:"column:week_no"`
	WeekLabel  string  `json:"weekLabel" gorm:"column:week_label"`
	Motivation *int16  `json:"motivation" gorm:"column:motivation"`
	Behavior   *int16  `json:"behavior" gorm:"column:behavior"`
	Status     string  `json:"status" gorm:"column:status"`
}

func (s *ReportService) StudentTrend(studentID uuid.UUID, onlyApproved bool) ([]StudentTrendRow, error) {
	var rows []StudentTrendRow
	approvedClause := ""
	if onlyApproved {
		approvedClause = " AND e.status = 'approved'"
	}
	q := `
		SELECT w.week_no, w.label AS week_label, e.motivation, e.behavior, e.status::text AS status
		FROM evaluations e
		JOIN evaluation_weeks w ON w.id = e.evaluation_week_id
		JOIN coach_student_assignments a ON a.id = e.assignment_id
		WHERE a.student_id = ?` + approvedClause + `
		ORDER BY w.week_no`
	err := s.db.Raw(q, studentID).Scan(&rows).Error
	return rows, err
}
