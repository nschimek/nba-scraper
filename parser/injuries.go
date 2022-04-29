package parser

import "time"

type Injuries struct {
	PlayerId, TeamId string
	Season           int
	UpdateDate       time.Time
	Description      string
}
