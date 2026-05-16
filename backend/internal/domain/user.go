package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Role         Role      `json:"role" gorm:"column:role;type:user_role"`
	FullName     string    `json:"fullName" gorm:"column:full_name"`
	Phone        string    `json:"phone" gorm:"column:phone"`
	Email        *string   `json:"email" gorm:"column:email"`
	PasswordHash string    `json:"-" gorm:"column:password_hash"`
	IsActive     bool      `json:"isActive" gorm:"column:is_active"`
	City         *string   `json:"city" gorm:"column:city"`
	Timestamps
}

func (User) TableName() string { return "users" }

type Student struct {
	UserID        uuid.UUID  `json:"userId" gorm:"column:user_id;type:uuid;primaryKey"`
	School        *string    `json:"school" gorm:"column:school"`
	Grade         *string    `json:"grade" gorm:"column:grade"`
	ParentID      *uuid.UUID `json:"parentId" gorm:"column:parent_id;type:uuid"`
	CoordinatorID *uuid.UUID `json:"coordinatorId" gorm:"column:coordinator_id;type:uuid"`
	Notes         *string    `json:"notes" gorm:"column:notes"`
	Timestamps

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (Student) TableName() string { return "students" }

type Coach struct {
	UserID     uuid.UUID   `json:"userId" gorm:"column:user_id;type:uuid;primaryKey"`
	Status     CoachStatus `json:"status" gorm:"column:status;type:coach_status"`
	Bio        *string     `json:"bio" gorm:"column:bio"`
	Experience *string     `json:"experience" gorm:"column:experience"`
	Timestamps

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (Coach) TableName() string { return "coaches" }

type Parent struct {
	UserID uuid.UUID `json:"userId" gorm:"column:user_id;type:uuid;primaryKey"`
	Notes  *string   `json:"notes" gorm:"column:notes"`
	Timestamps

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (Parent) TableName() string { return "parents" }

type Coordinator struct {
	UserID         uuid.UUID `json:"userId" gorm:"column:user_id;type:uuid;primaryKey"`
	FoundationName *string   `json:"foundationName" gorm:"column:foundation_name"`
	Notes          *string   `json:"notes" gorm:"column:notes"`
	Timestamps

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (Coordinator) TableName() string { return "coordinators" }
