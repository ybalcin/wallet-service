package utility

import "testing"

func TestStrLength(t *testing.T) {
	t.Run("whitespace str", func(t *testing.T) {
		l := StrLength(" ")

		if l != 0 {
			t.Errorf("expected str length is 0, but got: %d", l)
		}
	})

	t.Run("empty str", func(t *testing.T) {
		l := StrLength("")

		if l != 0 {
			t.Errorf("expected str length is 0, but got: %d", l)
		}
	})
}

func TestIsStrEmpty(t *testing.T) {
	t.Run("whitespace str", func(t *testing.T) {
		empty := IsStrEmpty(" ")

		if !empty {
			t.Errorf("expected str isEmpty value is true, but got false")
		}
	})

	t.Run("empty str", func(t *testing.T) {
		empty := IsStrEmpty("")

		if !empty {
			t.Errorf("expected str isEmpty value is true, but got false")
		}
	})
}
