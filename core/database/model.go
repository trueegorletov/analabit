package database

import (
	"analabit/core"
	"gorm.io/gorm"
	"time"
)

type Varsity struct {
	gorm.Model
	ID       int
	Code     string
	Name     string
	Headings []Heading
}

type Heading struct {
	gorm.Model
	ID        int
	VarsityID int
	Varsity   Varsity
	Capacity  int
	Name      string
}

type Application struct {
	gorm.Model
	HeadingID       int
	Heading         Heading
	StudentID       string
	RatingPlace     int
	Score           int
	competitionType core.Competition
	Priority        int
	Iteration       int
	UpdatedAt       time.Time
}

type Admission struct {
	gorm.Model
	HeadingID     int
	Heading       Heading
	StudentID     string
	AdmittedPlace int
	Iteration     int
	UpdatedAt     time.Time
}
