package management

import "time"

type StatManager struct {
	m *Management
}

func NewStatManager(m *Management) *StatManager {
	return &StatManager{m}
}

func (sm *StatManager) ActiveUsers() (int, error) {
	var i int
	err := sm.m.get(sm.m.uri("stats/active-users"), &i)
	return i, err
}

type DailyStat struct {
	Date            time.Time `json:"date"`
	Logins          int       `json:"logins"`
	Signups         int       `json:"signups"`
	LeakedPasswords int       `json:"leaked_passwords"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (sm *StatManager) Daily() ([]*DailyStat, error) {
	var ds []*DailyStat
	err := sm.m.get(sm.m.uri("stats/daily"), &ds)
	return ds, err
}
