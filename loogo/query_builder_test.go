package loogo

import "testing"

func TestInq(t *testing.T) {
	p1 := QueryParam{
		QueryType: "Inq",
		Field:     "status",
		Values:    []string{"done", "new"},
	}

	want := "filter[where][status][inq]=done&filter[where][status][inq]=new"
	if got := Inq(p1, false); got != want {
		t.Errorf("Inq() = %v, want %v", got, want)
	}
}

func TestEq(t *testing.T) {
	p1 := QueryParam{
		QueryType: "Eq",
		Field:     "status",
		Values:    []string{"new"},
	}

	want := "filter[where][status]=new"
	if got := Eq(p1, false); got != want {
		t.Errorf("Eq() = %v, want %v", got, want)
	}
}

func TestBetween(t *testing.T) {
	p1 := QueryParam{
		QueryType: "Between",
		Field:     "status",
		Values:    []string{"1", "2"},
	}

	want := "filter[where][status][between][0]=1&filter[where][status][between][1]=2"
	if got := Between(p1, false); got != want {
		t.Errorf("Between() = %v, want %v", got, want)
	}
}

func TestBuildQuery_not_for_count(t *testing.T) {
	p1 := QueryParam{
		QueryType: "Eq",
		Field:     "status",
		Values:    []string{"done"},
	}

	p2 := QueryParam{
		QueryType: "Between",
		Field:     "timestamp_ms",
		Values:    []string{"1", "2"},
	}

	params := QueryParams{
		p1,
		p2,
	}

	want := "?filter[where][status]=done&filter[where][timestamp_ms][between][0]=1&filter[where][timestamp_ms][between][1]=2"
	if got := BuildQuery(params, false); got != want {
		t.Errorf("BuildQuery() = %v, want %v", got, want)
	}
}

func TestBuildQuery_for_count(t *testing.T) {
	p1 := QueryParam{
		QueryType: "Eq",
		Field:     "status",
		Values:    []string{"new", "done"},
	}

	params := QueryParams{
		p1,
	}

	want := "?[where][status]=new"
	if got := BuildQuery(params, true); got != want {
		t.Errorf("BuildQuery() = %v, want %v", got, want)
	}
}
