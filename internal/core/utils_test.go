package core

import (
	"path/filepath"
	"testing"
)

func TestHashBotStateFile(t *testing.T) {
	for i, each := range []struct {
		input string
		want  string
	}{
		{"Sheriff", filepath.Join(StateFileFolder, "sheriff.json")},
		{"Versionista!", filepath.Join(StateFileFolder, "versionista.json")},
		{"Jenkins Integration Robot", filepath.Join(StateFileFolder, "jenkins-integration-robot.json")},
		{"робат", filepath.Join(StateFileFolder, "робат.json")},
		{"世界", filepath.Join(StateFileFolder, "世界.json")},
		{"!Trim me!", filepath.Join(StateFileFolder, "trim-me.json")},
		{"Sheriff :D", filepath.Join(StateFileFolder, "sheriff-d.json")},
	} {
		if got := hashBotStateFile(each.input); got != each.want {
			t.Errorf("[%v] got %#v want %#v", i, got, each.want)
		}
	}
}
