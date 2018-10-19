package lab7

import (
	"testing"
)

var timetests = []struct {
	in  string
	out string
}{
	{"2010/12/22, 20:22:32, 10.213.223.232, NRK2, NRK1", "20:22:32"},
	{"2010/12/22, 20:23:47, 10.213.223.232, Viasat 4, NRK2", "20:23:47"},
	{"2010/12/24, 10:48:26, 10.200.8.24, NRK 3, Disney XD", "10:48:26"},
	{"2010/12/24, 00:00:00, 10.200.8.24, NRK 3, Disney XD", "00:00:00"},
	{"2010/12/24, 23:59:59, 10.200.8.24, NRK 3, Disney XD", "23:59:59"},
	{"2010/12/24, 00:00:00, 10.200.8.24, NRK 3, Disney XD", "00:00:00"},
}

func TestSTBTime(t *testing.T) {
	for _, tt := range timetests {
		zap, _, _ := NewSTBEvent(tt.in)
		s := zap.Time.Format(timeOnly)
		if s != tt.out {
			t.Errorf("NewSTBEvent(%q) => %q, want %q", tt.in, s, tt.out)
		}
	}
}

var timeerrtests = []struct {
	in  string
	out string
}{
	{"2010/12/24, 24:00:00, 10.200.8.24, NRK 3, Disney XD", "NewSTBEvent: failed to parse timestamp"},
	{"2010/12/24, 00:60:00, 10.200.8.24, NRK 3, Disney XD", "NewSTBEvent: failed to parse timestamp"},
	{"2010/12/24, 00:00:60, 10.200.8.24, NRK 3, Disney XD", "NewSTBEvent: failed to parse timestamp"},
}

func TestSTBTimeErr(t *testing.T) {
	for _, tt := range timeerrtests {
		zap, schng, err := NewSTBEvent(tt.in)
		if zap != nil || schng != nil || err == nil {
			t.Errorf("NewSTBEvent(%q) => (%q, %q, %q), want (nil, nil, %q)",
				tt.in, zap, schng, err, tt.out)
		}
		if err.Error() != tt.out {
			t.Errorf("NewSTBEvent(%q) => (nil, nil, %q), want (nil, nil, %q)",
				tt.in, err, tt.out)
		}
	}
}

var toofewfieldstests = []struct {
	in  string
	out string
}{
	{"2010/12/24, 00:00:00", "NewSTBEvent: too short event string: 2010/12/24, 00:00:00"},
	{"2010/12/24, 00:00:00, 10.200.8.24 ", "NewSTBEvent: event with too few fields: 2010/12/24, 00:00:00, 10.200.8.24 "},
}

func TestSTBTooFewFields(t *testing.T) {
	for _, tt := range toofewfieldstests {
		zap, schng, err := NewSTBEvent(tt.in)
		if zap != nil || schng != nil || err == nil {
			t.Errorf("NewSTBEvent(%q) => (%q, %q, %q), want (nil, nil, %q)",
				tt.in, zap, schng, err, tt.out)
		}
		if err.Error() != tt.out {
			t.Errorf("NewSTBEvent(%q) => (nil, nil, %q), want (nil, nil, %q)",
				tt.in, err, tt.out)
		}
	}
}

var statuschangetests = []struct {
	in  string
	out string
}{
	{"2013/07/20, 21:57:42, 203.124.29.72, Volume: 50", "Volume: 50"},
	{"2013/07/20, 21:57:42, 203.124.29.72, Mute_Status: 0", "Mute_Status: 0"},
}

func TestSTBStatusChange(t *testing.T) {
	for _, tt := range statuschangetests {
		zap, schng, err := NewSTBEvent(tt.in)
		if zap != nil || schng == nil || err != nil {
			t.Errorf("NewSTBEvent(%q) => (%q, %q, %q), want (nil, %q, nil)",
				tt.in, zap, schng, err, tt.out)
		}
		//TODO activate code later
		// if schng.Status != tt.out {
		// 	t.Errorf("NewSTBEvent(%q) => (nil, %q, nil), want (nil, %q, nil)",
		// 		tt.in, schng.Status, tt.out)
		// }
	}
}
