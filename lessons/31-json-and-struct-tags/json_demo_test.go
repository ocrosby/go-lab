package jsondemo_test

import (
	"encoding/json"
	"strings"
	"testing"

	jd "github.com/ocrosby/go-lab/lessons/31-json-and-struct-tags"
)

func TestMarshal_UsesTagRenames(t *testing.T) {
	u := jd.User{ID: "u1", Email: "a@b.com", Age: 36}

	b, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("Marshal err = %v", err)
	}
	got := string(b)
	// Field names must be renamed (lowercase, snake) per the tags.
	if !strings.Contains(got, `"id":"u1"`) || !strings.Contains(got, `"email":"a@b.com"`) {
		t.Errorf("output missing renamed fields: %s", got)
	}
}

func TestMarshal_OmitEmptyDropsZeroAge(t *testing.T) {
	u := jd.User{ID: "u1", Email: "a@b.com", Age: 0}

	b, _ := json.Marshal(u)
	if strings.Contains(string(b), "age") {
		t.Errorf("omitempty didn't drop zero age: %s", b)
	}
}

func TestMarshal_MinusTagOmitsNotes(t *testing.T) {
	u := jd.User{ID: "u1", Email: "a@b.com", Notes: "internal-only"}

	b, _ := json.Marshal(u)
	if strings.Contains(string(b), "internal-only") || strings.Contains(string(b), "notes") {
		t.Errorf(`json:"-" tag failed to omit Notes: %s`, b)
	}
}

func TestUnmarshal_PopulatesFields(t *testing.T) {
	input := []byte(`{"id":"u2","email":"c@d.com","age":42}`)

	var u jd.User
	if err := json.Unmarshal(input, &u); err != nil {
		t.Fatalf("Unmarshal err = %v", err)
	}
	if u.ID != "u2" || u.Email != "c@d.com" || u.Age != 42 {
		t.Errorf("got %+v", u)
	}
}

func TestMoney_CustomMarshalRendersAsString(t *testing.T) {
	m := jd.Money(1234) // 12.34 dollars

	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Marshal err = %v", err)
	}
	if string(b) != `"12.34"` {
		t.Errorf("Marshal(1234 cents) = %s, want \"12.34\"", b)
	}
}

func TestMoney_RoundTripsThroughJSON(t *testing.T) {
	original := jd.Money(999) // 9.99

	b, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal err = %v", err)
	}

	var back jd.Money
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("unmarshal err = %v", err)
	}

	if back != original {
		t.Errorf("round-trip lost data: %d → %d", original, back)
	}
}

func TestDecodeStrict_RejectsUnknownFields(t *testing.T) {
	// "emial" is a typo — DisallowUnknownFields catches it.
	input := strings.NewReader(`{"emial":"typo@x.com"}`)

	var u jd.User
	err := jd.DecodeStrict(input, 1024, &u)

	if err == nil {
		t.Fatal("expected error on unknown field, got nil")
	}
	if !strings.Contains(err.Error(), "emial") {
		t.Errorf("err = %v, want to mention the unknown field name", err)
	}
}

func TestDecodeStrict_RejectsMultipleJSONValues(t *testing.T) {
	input := strings.NewReader(`{"id":"u1","email":"a@b.com"}{"id":"u2","email":"c@d.com"}`)

	var u jd.User
	err := jd.DecodeStrict(input, 1024, &u)

	if err == nil {
		t.Fatal("expected error on trailing JSON value, got nil")
	}
	if !strings.Contains(err.Error(), "exactly one") {
		t.Errorf("err = %v, want 'exactly one' complaint", err)
	}
}
