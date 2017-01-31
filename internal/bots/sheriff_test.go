package bots

import (
	"testing"

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
