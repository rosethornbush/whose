package output

import "testing"

func TestFormatDateUTC(t *testing.T) {
	got := formatDate("2007-10-09T18:20:50Z")
	want := "2007-10-09 18:20:50 UTC"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFormatDateConvertsOffsetToUTC(t *testing.T) {
	got := formatDate("2017-02-17T18:04:32-05:00")
	want := "2017-02-17 23:04:32 UTC"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFormatDatePreservesInvalidValue(t *testing.T) {
	got := formatDate("not-a-date")
	want := "not-a-date"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
