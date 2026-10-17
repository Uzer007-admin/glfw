package glfw

import "testing"

func TestXIMCaretDirections(t *testing.T) {
	text := []rune("абв abc\n日本\nlast")
	cases := []struct{ cursor, direction, position, want int }{{0, 0, 0, 1}, {0, 1, 0, 0}, {0, 2, 0, 4}, {6, 3, 0, 4}, {9, 4, 0, 1}, {5, 5, 0, 10}, {5, 6, 0, 8}, {9, 7, 0, 0}, {9, 8, 0, 8}, {9, 9, 0, 10}, {2, 10, 9, 9}, {2, 11, 0, 2}}
	for _, c := range cases {
		if got := ximCaret(text, c.cursor, c.position, c.direction); got != c.want {
			t.Errorf("direction=%d cursor=%d: got %d want %d", c.direction, c.cursor, got, c.want)
		}
	}
}
