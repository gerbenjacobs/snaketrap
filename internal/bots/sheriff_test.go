package bots

import (
	"testing"

	"time"

	"github.com/gerbenjacobs/snaketrap/internal/core"
)

func createSheriff() *Sheriff {
	s := &Sheriff{}
	s.config = SheriffConfig{
		Days:  []int{1, 2, 3, 4, 5},
		Time:  "12:00",
		Topic: "Current sheriff: %s - Some MOTD here..",
		// SherrifUsers will be alphabetically sorted.
		SheriffUsers: SheriffUsers{
			{Name: "gerben", Away: true},
			{Name: "robpike", Away: false},
			{Name: "davecheney", Away: false},
			{Name: "jessfraz", Away: false},
		},
	}
	s.wrangler = &core.Wrangler{}
	s.boot()
	return s
}

func TestSheriffNext(t *testing.T) {
	s := createSheriff()
	s.current = 1 // set sheriff to "gerben"

	// trigger next and validate
	got := s.nextSheriff(true, s.current)
	want := 2 // "jessfraz"
	if got != want {
		t.Errorf("wrong next sheriff. Got: %d Want: %d", got, want)
	}
}

func TestSheriffPreviousWrapAround(t *testing.T) {
	s := createSheriff()
	s.current = 0 // set sheriff to "davecheney"

	// trigger previous and validate wrapround
	got := s.nextSheriff(false, s.current)
	want := 3 // "robpike"
	if got != want {
		t.Errorf("wrong next sheriff. Got: %d Want: %d", got, want)
	}
}

func TestSheriffNextSkipsUnavailable(t *testing.T) {
	s := createSheriff()
	s.current = 0 // set sheriff to "davecheney"

	// trigger next and validate that 1:"gerben" is skipped because of away status
	got := s.nextSheriff(true, s.current)
	want := 2 // "jessfraz"
	if got != want {
		t.Errorf("wrong next sheriff. Got: %d Want: %d", got, want)
	}
}

func TestIsActiveDay(t *testing.T) {
	s := createSheriff()

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Errorf("failed to load timezone UTC: %s", err)
	}

	for _, each := range []struct {
		name   string
		day    time.Time
		active bool
	}{
		{"mon", time.Date(2017, 01, 30, 12, 0, 0, 0, loc), true},
		{"tue", time.Date(2017, 01, 31, 12, 0, 0, 0, loc), true},
		{"wed", time.Date(2017, 02, 1, 12, 0, 0, 0, loc), true},
		{"thu", time.Date(2017, 02, 2, 12, 0, 0, 0, loc), true},
		{"fri", time.Date(2017, 02, 3, 12, 0, 0, 0, loc), true},
		{"sat", time.Date(2017, 02, 4, 12, 0, 0, 0, loc), false},
		{"sun", time.Date(2017, 02, 5, 12, 0, 0, 0, loc), false},
	} {
		got := s.isActiveDay(each.day)
		if got != each.active {
			t.Errorf("[%v] isActiveDay returned %t, want %t", each.name, got, each.active)
		}
	}
}

func TestIsActiveDayTomorrow(t *testing.T) {
	s := createSheriff()

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Errorf("failed to load timezone UTC: %s", err)
	}

	for _, each := range []struct {
		name   string
		day    time.Time
		active bool
	}{
		{"fri", time.Date(2017, 02, 3, 12, 0, 0, 0, loc), false},
		{"sat", time.Date(2017, 02, 4, 12, 0, 0, 0, loc), false},
		{"sun", time.Date(2017, 02, 5, 12, 0, 0, 0, loc), true},
		{"mon", time.Date(2017, 01, 30, 12, 0, 0, 0, loc), true},
	} {
		got := s.isActiveDayTomorrow(each.day)
		if got != each.active {
			t.Errorf("[%v] isActiveDay returned %t, want %t", each.name, got, each.active)
		}
	}
}
