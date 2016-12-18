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
		Users: []string{
			"gerben",
			"robpike",
			"davecheney",
		},
	}
	s.wrangler = &core.Wrangler{}
	s.boot()
	return s
}

func TestSheriffNext(t *testing.T) {
	s := createSheriff()

	// validate current sheriff
	want := 0 // "davecheney"
	got := s.current
	if got != want {
		t.Errorf("wrong sheriff. Got: %d Want: %d", got, want)
	}

	// trigger next and validate
	s.rotateSheriff(true)
	want = 1 // "gerben"
	got = s.current
	if got != want {
		t.Errorf("wrong next sheriff. Got: %d Want: %d", got, want)
	}
}

func TestSheriffPreviousWrapAround(t *testing.T) {
	s := createSheriff()

	// validate current sheriff
	want := 0 // "davecheney"
	got := s.current
	if got != want {
		t.Errorf("wrong sheriff. Got: %d Want: %d", got, want)
	}

	// trigger next and validate
	s.rotateSheriff(false)
	want = 2 // "robpike"
	got = s.current
	if got != want {
		t.Errorf("wrong next sheriff. Got: %d Want: %d", got, want)
	}
}
