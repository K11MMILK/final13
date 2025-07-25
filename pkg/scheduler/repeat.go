package scheduler

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string, strictlyAfterNow bool) (string, error) {
	start, err := parseDateString(dstart)
	if err != nil {
		return "", err
	}
	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}
	if start.Year() < 1900 || start.Year() > 2100 {
		start = time.Date(now.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
		strictlyAfterNow = false
	}

	parts := strings.Fields(repeat)
	if len(parts) < 1 {
		return "", errors.New("repeat rule format error")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid daily format")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid daily interval")
		}
		t := start

		if strictlyAfterNow {
			t = t.AddDate(0, 0, days)
		}

		for i := 0; i <= 400; i++ {
			if afterNow(t, now, strictlyAfterNow) {
				return t.Format("20060102"), nil
			}
			t = t.AddDate(0, 0, days)
		}
		return "", errors.New("no valid daily date found")

	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid yearly format")
		}

		t := start

		if strictlyAfterNow {
			t = t.AddDate(1, 0, 0)
		}

		for i := 0; i <= 1000; i++ {
			if afterNow(t, now, strictlyAfterNow) {
				return t.Format("20060102"), nil
			}
			t = t.AddDate(1, 0, 0)
		}
		return "", errors.New("no valid yearly date found")

	case "w":
		if len(parts) < 2 {
			return "", errors.New("missing weekdays for weekly repeat")
		}

		var days []int
		for _, s := range strings.Split(parts[1], ",") {
			d, err := strconv.Atoi(s)
			if err != nil {
				return "", fmt.Errorf("invalid weekday: %w", err)
			}
			if d < 1 || d > 7 {
				return "", errors.New("weekday must be 1 (Mon) to 7 (Sun)")
			}
			days = append(days, d)
		}

		t := start

		if strictlyAfterNow {
			t = t.AddDate(0, 0, 1)
		}

		for i := 0; i < 1000; i++ {
			date := t.AddDate(0, 0, i)
			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7
			}
			for _, d := range days {
				if wd == d && afterNow(date, now, strictlyAfterNow) {
					return date.Format("20060102"), nil
				}
			}
		}

		return "", errors.New("could not find valid weekday")

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid monthly format")
		}

		days, err := parseIntList(parts[1], -31, 31)
		if err != nil {
			return "", err
		}

		months := []int{}
		if len(parts) == 3 {
			months, err = parseIntList(parts[2], 1, 12)
			if err != nil {
				return "", err
			}
		}

		if len(months) == 0 {
			months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		}
		sort.Ints(months)

		t := start

		if strictlyAfterNow {
			t = t.AddDate(0, 0, 1)
		}

		for i := 0; i <= 1000; i++ {
			checkDate := t.AddDate(0, i, 0)
			y, m := checkDate.Year(), checkDate.Month()

			for j := 0; j < len(months); j++ {
				candidateMonth := months[(sort.SearchInts(months, int(m))+j)%len(months)]
				candidateYear := y
				if candidateMonth < int(m) || (candidateMonth == int(m) && j > 0) {
					candidateYear++
				}

				firstOfNext := time.Date(candidateYear, time.Month(candidateMonth), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
				lastOfMonth := firstOfNext.AddDate(0, 0, -1)

				sort.Slice(days, func(i, j int) bool {
					di, dj := days[i], days[j]
					if di >= 0 && dj >= 0 {
						return di < dj
					}
					if di < 0 && dj < 0 {
						return di < dj
					}
					return di >= 0
				})

				for _, d := range days {
					if d < -2 {
						return "", errors.New("no valid monthly date found")
					}
				}

				for _, d := range days {
					var day int

					if d > 0 {
						day = d
					} else {
						day = lastOfMonth.Day() + d + 1
					}

					if day < 1 || day > lastOfMonth.Day() {
						continue
					}

					t := time.Date(candidateYear, time.Month(candidateMonth), day, 0, 0, 0, 0, time.UTC)
					if afterNow(t, now, strictlyAfterNow) {
						return t.Format("20060102"), nil
					}
				}
			}
		}

		return "", errors.New("no valid monthly date found")

	default:
		return "", errors.New("unknown repeat rule type")
	}
}

func afterNow(t, now time.Time, strictly bool) bool {
	t = truncateToDate(t)
	now = truncateToDate(now)
	if strictly {
		return t.After(now)
	}
	return !t.Before(now)
}

func parseDateString(s string) (time.Time, error) {
	return time.Parse("20060102", s)
}

func parseIntList(s string, min, max int) ([]int, error) {
	var result []int
	for _, part := range strings.Split(s, ",") {
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		if n < min || n > max {
			return nil, fmt.Errorf("value %d out of bounds (%d..%d)", n, min, max)
		}
		result = append(result, n)
	}
	return result, nil
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
