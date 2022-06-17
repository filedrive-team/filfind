package utils

import (
	"testing"
	"time"
)

func TestGetEpochByTime(t *testing.T) {
	cases := []struct {
		Name     string
		Time     time.Time
		Expected int64
	}{
		{"case 0", time.Unix(1641523050, 0), 1440555},
		{"case 1", time.Unix(GenesisUnixTime, 0), 0},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := GetEpochByTime(c.Time); ans != c.Expected {
				t.Fatalf("input epoch[%v] expected [%v], but ans[%v] got", c.Time, c.Expected, ans)
			}
		})
	}
}

func TestGetBlockTimeByEpoch(t *testing.T) {
	cases := []struct {
		Name     string
		Epoch    int64
		Expected int64
	}{
		{"case 0", 1, GenesisUnixTime + 30},
		{"case 1", 1440555, 1641523050},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := GetBlockTimeByEpoch(c.Epoch); ans.Unix() != c.Expected {
				t.Fatalf("input epoch[%v] expected [%v], but ans[%v] got", c.Epoch, c.Expected, ans)
			}
		})
	}
}

func TestGetDurationByEpoch(t *testing.T) {
	cases := []struct {
		Name     string
		Epoch    int64
		Expected time.Duration
	}{
		{"case 0", 1, 30 * time.Second},
		{"case 1", 100, 3000 * time.Second},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := GetDurationByEpoch(c.Epoch); ans != c.Expected {
				t.Fatalf("input epoch[%v] expected [%v], but ans[%v] got", c.Epoch, c.Expected, ans)
			}
		})
	}
}

func TestMonthBegin(t *testing.T) {
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", "2022-05-17 11:14:00", time.Now().Location())
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", "2022-05-01 00:00:00", time.Now().Location())
	cases := []struct {
		Name     string
		Time     time.Time
		Expected time.Time
	}{
		{"case 0", t1, t2},
		{"case 1", t2, t2},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := MonthBegin(c.Time); ans != c.Expected {
				t.Fatalf("input epoch[%v] expected [%v], but ans[%v] got", c.Time, c.Expected, ans)
			}
		})
	}
}
