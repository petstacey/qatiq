package validator

import "testing"

func TestCheckIsTrue(t *testing.T) {
	v := New()

	str := "apples"
	v.Check(str == "apples", "str", "sould match 'apples'")
	valid := v.Valid()
	if !valid {
		t.Errorf("error when check is valid. expected 'true', got '%t", valid)
	}
}

func TestCheckIsFalse(t *testing.T) {
	v := New()

	str := "applies"
	v.Check(str == "oranges", "str", "should match 'apples'")
	valid := v.Valid()
	if valid {
		t.Errorf("error when check is not valid. expected 'false', got '%t'", valid)
	}
}

func TestIsPermitted(t *testing.T) {
	permittedStrings := []string{"apples", "oranges"}

	permitted := PermittedValue("apples", permittedStrings...)
	if !permitted {
		t.Errorf("error when permitted value. expected 'true', got '%t'", permitted)
	}
}

func TestIsNotPermitted(t *testing.T) {
	permittedStrings := []string{"apples", "oranges"}

	permitted := PermittedValue("bananas", permittedStrings...)
	if permitted {
		t.Errorf("error when not permitted value. expected 'false', got '%t'", permitted)
	}
}

func TestIsUnique(t *testing.T) {
	uniqueStrings := []string{"apples", "oranges", "bananas"}

	unique := Unique(uniqueStrings)
	if !unique {
		t.Errorf("error when list is unique. expected 'true', got '%t'", unique)
	}
}

func TestIsNotUnique(t *testing.T) {
	uniqueStrings := []string{"apples", "oranges", "bananas", "apples"}

	unique := Unique(uniqueStrings)
	if unique {
		t.Errorf("error when list is unique. expected 'true', got '%t'", unique)
	}
}
