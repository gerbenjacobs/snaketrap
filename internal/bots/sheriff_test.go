package bots

import (
	"testing"

	"github.com/gerbenjacobs/snaketrap/internal/core"
)

func createSheriff() *Sheriff {
	s := &Sheriff{}
	s.sheriffCfg = SheriffConfig{
		Time:  "12:00",
		Topic: "Current sheriff: %s - Some MOTD here..",
		Users: []string{
			"gerben",
			"robpike",
			"davecheney",
		},
	}
	s.wrangler = &core.Wrangler{
		DefaultRoom: "1234",
	}
	s.boot()
	return s
}

func TestSheriffNext(t *testing.T) {
	s := createSheriff()

	// validate current sheriff
	want := 0 // "davecheney"
	got := s.currentSheriff
	if got != want {
		t.Errorf("wrong sheriff. Got: %d Want: %d", got, want)
	}

	// trigger next and validate
	s.next()
	want = 1 // "gerben"
	got = s.currentSheriff
	if got != want {
		t.Errorf("wrong next sheriff. Got: %d Want: %d", got, want)
	}
}

func TestSheriffPreviousWrapAround(t *testing.T) {
	s := createSheriff()

	// validate current sheriff
	want := 0 // "davecheney"
	got := s.currentSheriff
	if got != want {
		t.Errorf("wrong sheriff. Got: %d Want: %d", got, want)
	}

	// trigger next and validate
	s.previous()
	want = 2 // "robpike"
	got = s.currentSheriff
	if got != want {
		t.Errorf("wrong next sheriff. Got: %d Want: %d", got, want)
	}
}

func TestSheriffExtractAction(t *testing.T) {
	s := createSheriff()

	testCases := []struct {
		in   string
		want string
	}{
		{"/bot sheriff next", "next"},
		{"/bot sheriff previous", "previous"},
		{"/bot sheriff away", "away"},
		{"/bot sheriff back", "back"},

		{"/bot sheriff remove", ""},
		{"/bot sheriff status back", ""},
		{"/bot sheriff ", ""},
		{"/bot sheriff", ""},
		{"/bot sheriff --help", ""},
	}

	for _, each := range testCases {
		got := s.extractAction(each.in)
		if got != each.want {
			t.Errorf("failed to extract action, got: %s want: %s", got, each.want)
		}
	}
}
