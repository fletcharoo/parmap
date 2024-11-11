package parmap_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/fletcharoo/parmap"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestMain(m *testing.M) {
	r := m.Run()
	snaps.Clean(m, snaps.CleanOpts{Sort: true})
	os.Exit(r)
}

var (
	_errMap = parmap.ErrMap{
		1: fmt.Errorf("error 1"),
		2: fmt.Errorf("error 2"),
		3: fmt.Errorf("error 3"),
	}
	_inputs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
)

func Test_ErrMap_ErrJoin(t *testing.T) {
	snaps.MatchSnapshot(t, _errMap.ErrJoin())
}

func Test_ErrMap_String(t *testing.T) {
	snaps.MatchSnapshot(t, fmt.Sprint(_errMap.ErrJoin()))
}

func Test_Do(t *testing.T) {
	do := func(i int) (int, error) {
		return i * 5, nil
	}

	result, err := parmap.Do(_inputs, do)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	snaps.MatchJSON(t, result)
}

func Test_Do_Error(t *testing.T) {
	do := func(i int) (int, error) {
		if i%3 == 0 {
			return 0, fmt.Errorf("error %d", i)
		}

		return i, nil
	}

	result, err := parmap.Do(_inputs, do)

	t.Run("Error", func(t *testing.T) {
		snaps.MatchSnapshot(t, err)
	})

	t.Run("Result", func(t *testing.T) {
		snaps.MatchSnapshot(t, result)
	})
}
