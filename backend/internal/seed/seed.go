package seed

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/config"
	"github.com/koc-luk/backend/internal/domain"
)

// Tüm test kullanıcıları için ortak şifre. Login butonlarında kolay olsun diye sabit.
const TestPassword = "Test123!"

// SeedAdmin idempotent: yoksa default admin kullanıcıyı ekler.
func SeedAdmin(db *gorm.DB, cfg *config.Config) error {
	const phone = "05000000000"
	const password = "Admin123!"

	var existing domain.User
	if err := db.Where("phone = ?", phone).First(&existing).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cfg.BcryptCost)
	if err != nil {
		return err
	}
	email := "admin@kare-rehber.local"
	city := "İstanbul"
	u := domain.User{
		ID:           uuid.New(),
		Role:         domain.RoleAdmin,
		FullName:     "Sistem Yöneticisi",
		Phone:        phone,
		Email:        &email,
		PasswordHash: string(hash),
		IsActive:     true,
		City:         &city,
	}
	return db.Create(&u).Error
}

// SeedWeeks idempotent: tablo boşsa 36 haftalık akademik yılı oluşturur.
func SeedWeeks(db *gorm.DB) error {
	var count int64
	if err := db.Model(&domain.EvaluationWeek{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	year := time.Now().Year()
	if time.Now().Month() < time.September {
		year--
	}
	start := firstMondayOf(year, time.September)

	weeks := make([]domain.EvaluationWeek, 0, 36)
	for i := 1; i <= 36; i++ {
		s := start.AddDate(0, 0, (i-1)*7)
		e := s.AddDate(0, 0, 6)
		label := fmt.Sprintf("%d. Hafta (%s-%s)", i, s.Format("02 Jan"), e.Format("02 Jan"))

		w := domain.EvaluationWeek{
			ID:        uuid.New(),
			WeekNo:    i,
			Label:     label,
			StartDate: s,
			EndDate:   e,
			IsOpen:    isWeekCurrentlyOpen(s, e),
		}
		weeks = append(weeks, w)
	}
	return db.Create(&weeks).Error
}

func firstMondayOf(year int, month time.Month) time.Time {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

func isWeekCurrentlyOpen(start, end time.Time) bool {
	now := time.Now()
	cutoffEnd := end.AddDate(0, 0, 7)
	return !now.Before(start) && !now.After(cutoffEnd)
}

// ---------------------------------------------------------------------------
// test users + atamalar + örnek değerlendirmeler
// ---------------------------------------------------------------------------

type seedSpec struct {
	role     domain.Role
	phone    string
	fullName string
	email    string
	city     string
}

// SeedTestUsers idempotent: 3 koç, 6 öğrenci, 3 veli, 2 koordinatör + atamalar +
// 3 örnek değerlendirme oluşturur. Var olan kayıtları atlar.
func SeedTestUsers(db *gorm.DB, cfg *config.Config) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(TestPassword), cfg.BcryptCost)
	if err != nil {
		return err
	}
	hashStr := string(hash)

	specs := []seedSpec{
		{domain.RoleCoach, "05311111111", "Ahmet Koç", "ahmet@kare-rehber.local", "İstanbul"},
		{domain.RoleCoach, "05311111112", "Zeynep Koç", "zeynep@kare-rehber.local", "Ankara"},
		{domain.RoleCoach, "05311111113", "Mehmet Koç", "mehmet@kare-rehber.local", "İzmir"},

		{domain.RoleStudent, "05322222221", "Ali Öğrenci", "ali@kare-rehber.local", "İstanbul"},
		{domain.RoleStudent, "05322222222", "Ayşe Öğrenci", "ayse@kare-rehber.local", "İstanbul"},
		{domain.RoleStudent, "05322222223", "Burak Öğrenci", "burak@kare-rehber.local", "Ankara"},
		{domain.RoleStudent, "05322222224", "Ceren Öğrenci", "ceren@kare-rehber.local", "Ankara"},
		{domain.RoleStudent, "05322222225", "Derya Öğrenci", "derya@kare-rehber.local", "İzmir"},
		{domain.RoleStudent, "05322222226", "Emir Öğrenci", "emir@kare-rehber.local", "İzmir"},

		{domain.RoleParent, "05333333331", "Hasan Veli", "veli1@kare-rehber.local", "İstanbul"},
		{domain.RoleParent, "05333333332", "Fatma Veli", "veli2@kare-rehber.local", "Ankara"},
		{domain.RoleParent, "05333333333", "Selma Veli", "veli3@kare-rehber.local", "İzmir"},

		{domain.RoleCoordinator, "05344444441", "Aslı Koordinatör", "kor1@kare-rehber.local", "İstanbul"},
		{domain.RoleCoordinator, "05344444442", "Murat Koordinatör", "kor2@kare-rehber.local", "Ankara"},
	}

	users := make(map[string]uuid.UUID, len(specs))
	for _, s := range specs {
		var existing domain.User
		err := db.Where("phone = ?", s.phone).First(&existing).Error
		if err == nil {
			users[s.phone] = existing.ID
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		email := s.email
		city := s.city
		u := domain.User{
			ID:           uuid.New(),
			Role:         s.role,
			FullName:     s.fullName,
			Phone:        s.phone,
			Email:        &email,
			PasswordHash: hashStr,
			IsActive:     true,
			City:         &city,
		}
		if err := db.Create(&u).Error; err != nil {
			return fmt.Errorf("create user %s: %w", s.phone, err)
		}

		switch s.role {
		case domain.RoleStudent:
			school := "Anadolu Lisesi"
			grade := "11"
			if err := db.Create(&domain.Student{UserID: u.ID, School: &school, Grade: &grade}).Error; err != nil {
				return err
			}
		case domain.RoleCoach:
			bio := "Deneyimli koç"
			if err := db.Create(&domain.Coach{UserID: u.ID, Status: domain.CoachApproved, Bio: &bio}).Error; err != nil {
				return err
			}
		case domain.RoleParent:
			if err := db.Create(&domain.Parent{UserID: u.ID}).Error; err != nil {
				return err
			}
		case domain.RoleCoordinator:
			fname := "Eğitim Vakfı"
			if err := db.Create(&domain.Coordinator{UserID: u.ID, FoundationName: &fname}).Error; err != nil {
				return err
			}
		}
		users[s.phone] = u.ID
	}

	pairs := [][2]string{
		{"05322222221", "05333333331"},
		{"05322222222", "05333333331"},
		{"05322222223", "05333333332"},
		{"05322222224", "05333333332"},
		{"05322222225", "05333333333"},
		{"05322222226", "05333333333"},
	}
	for _, p := range pairs {
		stuID, parID := users[p[0]], users[p[1]]
		if err := db.Model(&domain.Student{}).Where("user_id = ?", stuID).Update("parent_id", parID).Error; err != nil {
			return err
		}
	}

	corStuMap := []struct {
		coord    string
		students []string
	}{
		{"05344444441", []string{"05322222221", "05322222222", "05322222223", "05322222224"}},
		{"05344444442", []string{"05322222225", "05322222226"}},
	}
	for _, m := range corStuMap {
		corID := users[m.coord]
		for _, sp := range m.students {
			stuID := users[sp]
			if err := db.Model(&domain.Student{}).Where("user_id = ?", stuID).Update("coordinator_id", corID).Error; err != nil {
				return err
			}
		}
	}

	coachStuMap := []struct {
		coach    string
		students []string
	}{
		{"05311111111", []string{"05322222221", "05322222222"}},
		{"05311111112", []string{"05322222223", "05322222224"}},
		{"05311111113", []string{"05322222225", "05322222226"}},
	}
	assignments := make(map[string]uuid.UUID)
	adminID := getAdminID(db)
	for _, m := range coachStuMap {
		coachID := users[m.coach]
		for _, sp := range m.students {
			stuID := users[sp]
			var existing domain.CoachStudentAssignment
			err := db.Where("coach_id = ? AND student_id = ? AND is_active = true", coachID, stuID).First(&existing).Error
			if err == nil {
				assignments[sp] = existing.ID
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			a := domain.CoachStudentAssignment{
				ID:         uuid.New(),
				CoachID:    coachID,
				StudentID:  stuID,
				StartedAt:  time.Now(),
				IsActive:   true,
				AssignedBy: &adminID,
			}
			if err := db.Create(&a).Error; err != nil {
				return err
			}
			assignments[sp] = a.ID
		}
	}

	return seedSampleEvaluations(db, assignments, users)
}

func seedSampleEvaluations(db *gorm.DB, assignments map[string]uuid.UUID, users map[string]uuid.UUID) error {
	var weeks []domain.EvaluationWeek
	if err := db.Where("is_open = true").Order("week_no asc").Find(&weeks).Error; err != nil {
		return err
	}
	if len(weeks) == 0 {
		return nil
	}
	week := weeks[0]

	type evalSpec struct {
		studentPhone string
		coachPhone   string
		course       string
		homework     bool
		motivation   int16
		behavior     int16
		note         string
		approved     bool
	}
	specs := []evalSpec{
		{"05322222221", "05311111111", "Matematik konularında ilerleme var, türev yeterli.", true, 4, 5, "Hafta boyunca planlı çalıştı, deneme sınavında 380 puan.", true},
		{"05322222222", "05311111111", "Türkçe paragraf zorlanıyor, ek kaynak verildi.", false, 3, 4, "Ödev tamamlanmadı, hafta sonu telafi planlandı.", false},
		{"05322222223", "05311111112", "Fizik soruları çözmeye başladı.", true, 5, 5, "Çok motive, kendi başına soru çözüyor.", true},
	}

	for _, s := range specs {
		assignID, ok := assignments[s.studentPhone]
		if !ok {
			continue
		}
		coachID := users[s.coachPhone]

		var existing domain.Evaluation
		err := db.Where("assignment_id = ? AND evaluation_week_id = ?", assignID, week.ID).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		course := s.course
		hw := s.homework
		mot := s.motivation
		beh := s.behavior
		note := s.note
		ev := domain.Evaluation{
			ID:               uuid.New(),
			AssignmentID:     assignID,
			EvaluationWeekID: week.ID,
			CourseStatus:     &course,
			HomeworkDone:     &hw,
			Motivation:       &mot,
			Behavior:         &beh,
			GeneralNote:      &note,
			Status:           domain.EvalPending,
			SubmittedBy:      coachID,
			SubmittedAt:      time.Now(),
		}
		if s.approved {
			ev.Status = domain.EvalApproved
			adminID := getAdminID(db)
			now := time.Now()
			ev.ApprovedBy = &adminID
			ev.ApprovedAt = &now
		}
		if err := db.Create(&ev).Error; err != nil {
			return err
		}
	}
	return nil
}

func getAdminID(db *gorm.DB) uuid.UUID {
	var u domain.User
	db.Where("role = ?", domain.RoleAdmin).First(&u)
	return u.ID
}

// WipeAll bütün test verisini siler (idempotent değil, dikkatle kullan).
func WipeAll(db *gorm.DB) error {
	tables := []string{
		"evaluation_versions",
		"evaluations",
		"messages",
		"message_threads",
		"sms_logs",
		"audit_logs",
		"coach_student_assignments",
		"registrations",
		"students",
		"coaches",
		"parents",
		"coordinators",
		"users",
		"evaluation_weeks",
	}
	for _, t := range tables {
		if err := db.Exec("DELETE FROM " + t).Error; err != nil {
			return fmt.Errorf("wipe %s: %w", t, err)
		}
	}
	return nil
}
