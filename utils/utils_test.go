package utils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoryHostAndSlug(t *testing.T) {
	testCases := map[string][]string{
		"nbedos/cistern": {
			// SSH git URL
			"git@github.com:nbedos/cistern.git",
			"git@github.com:nbedos/cistern",
			// HTTPS git URL
			"https://github.com/nbedos/cistern.git",
			// Host shouldn't matter
			"git@gitlab.com:nbedos/cistern.git",
			// Web URLs should work too
			"https://gitlab.com/nbedos/cistern",
			// URL scheme shouldn't matter
			"http://gitlab.com/nbedos/cistern",
			"gitlab.com/nbedos/cistern",
			// Trailing slash
			"gitlab.com/nbedos/cistern/",
		},
		"namespace/nbedos/cistern": {
			"https://gitlab.com/namespace/nbedos/cistern",
		},
		"long/namespace/nbedos/cistern": {
			"git@gitlab.com:long/namespace/nbedos/cistern.git",
			"https://gitlab.com/long/namespace/nbedos/cistern",
			"git@gitlab.com:long/namespace/nbedos/cistern.git",
		},
	}

	for expectedSlug, urls := range testCases {
		for _, u := range urls {
			t.Run(fmt.Sprintf("URL: %v", u), func(t *testing.T) {
				_, slug, err := RepositoryHostAndSlug(u)
				if err != nil {
					t.Fatal(err)
				}
				if slug != expectedSlug {
					t.Fatalf("expected %q but got %q", expectedSlug, slug)
				}
			})
		}
	}

	badURLs := []string{
		// Missing 1 path component
		"git@github.com:nbedos.git",
		// Invalid url (colon)
		"https://github.com:nbedos/cistern.git",
	}

	for _, u := range badURLs {
		t.Run(fmt.Sprintf("URL: %v", u), func(t *testing.T) {
			if _, _, err := RepositoryHostAndSlug(u); err == nil {
				t.Fatalf("expected error but got nil for URL %q", u)
			}
		})
	}
}

func TestNullDuration_String(t *testing.T) {
	testCases := []struct {
		name string
		d    time.Duration
		s    string
	}{
		{
			name: "zero duration",
			d:    0,
			s:    "0s",
		},
		{
			name: "less than a second",
			d:    time.Nanosecond,
			s:    "<1s",
		},
		{
			name: "less than a minute",
			d:    time.Second,
			s:    "1s",
		},
		{
			name: "less than an hour",
			d:    time.Minute + 2*time.Second,
			s:    "1m02s",
		},
		{
			name: "more than an hour",
			d:    time.Hour + 2*time.Minute + 3*time.Second,
			s:    "1h02m03s",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			d := NullDuration{
				Valid:    true,
				Duration: testCase.d,
			}
			if diff := cmp.Diff(testCase.s, d.String()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestPollingStrategy_NextInterval(t *testing.T) {
	s := PollingStrategy{
		InitialInterval: time.Millisecond,
		Multiplier:      2,
		Randomizer:      0.25,
		MaxInterval:     time.Minute,
	}

	rand.Seed(0)
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			interval := s.InitialInterval
			for j := 0; j < 10; j++ {
				nextInterval := s.NextInterval(interval)
				if float64(nextInterval) < float64(interval)*(2-0.25) || float64(nextInterval) > float64(interval)*(2+0.25) {
					t.Fail()
				}
				interval = nextInterval
			}
		})
	}
}

func TestNewPollingStrategy(t *testing.T) {
	t.Run("any zero interval duration must be replaced by the corresponding duration of the default strategy", func(t *testing.T) {
		d := PollingStrategy{
			InitialInterval: 1,
			Multiplier:      2,
			Randomizer:      3,
			MaxInterval:     4,
		}
		s, err := NewPollingStrategy(0, 0, false, d)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(s, d); diff != "" {
			t.Fatal(diff)
		}
	})
}
