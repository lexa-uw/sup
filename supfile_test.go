package sup

import (
	"testing"
)

func TestParseSample(t *testing.T) {

	_, err := NewSupfile("./invalid_file")

	if err == nil {
		t.Error("Should have failed on invalid file")
	}

	_, err = NewSupfile("./sample_supfile.yml")

	if err != nil {
		t.Error(err)
	}

}

func TestExpandHosts(t *testing.T) {
	supfile, _ := NewSupfile("./sample_supfile.yml")
	err := supfile.LoadNetwork("prod")

	if err != nil {
		t.Error(err)
	}

	var findInSlice = func(haystack []string, needle string) bool {
		for _, t := range haystack {
			if t == needle {
				return true
			}
		}
		return false
	}

	prod := supfile.Networks["prod"]
	appendedHosts := findInSlice(prod.Hosts, "x1.dev") && findInSlice(prod.Hosts, "x2.dev")

	if !appendedHosts {
		t.Error("Failed to append dynamic hosts")
	}

}
