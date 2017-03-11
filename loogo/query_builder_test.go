package loogo

import "testing"

func TestInq(t *testing.T) {
	p1 := QueryParam{
		QueryType: "inq",
		Field:     "status",
		Values:    []string{"done", "new"},
	}

	want := "filter[where][status][inq]=done&filter[where][status][inq]=new"
	if got := inq(p1, false); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestEq(t *testing.T) {
	p1 := QueryParam{
		QueryType: "eq",
		Field:     "status",
		Values:    []string{"new"},
	}

	want := "filter[where][status]=new"
	if got := eq(p1, false); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBetween(t *testing.T) {
	p1 := QueryParam{
		QueryType: "between",
		Field:     "status",
		Values:    []string{"1", "2"},
	}

	want := "filter[where][status][between][0]=1&filter[where][status][between][1]=2"
	if got := between(p1, false); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildQueryNotForCount(t *testing.T) {
	p1 := QueryParam{
		QueryType: "eq",
		Field:     "status",
		Values:    []string{"done"},
	}

	p2 := QueryParam{
		QueryType: "between",
		Field:     "timestamp_ms",
		Values:    []string{"1", "2"},
	}

	params := QueryParams{
		p1,
		p2,
	}

	want := "?filter[where][status]=done&filter[where][timestamp_ms][between][0]=1&filter[where][timestamp_ms][between][1]=2"
	if got := buildQuery(params, false); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildQueryForCount(t *testing.T) {
	p1 := QueryParam{
		QueryType: "eq",
		Field:     "status",
		Values:    []string{"new", "done"},
	}

	params := QueryParams{
		p1,
	}

	want := "?[where][status]=new"
	if got := buildQuery(params, true); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
