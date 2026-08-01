package markdown

import "github.com/charmbracelet/glamour"

// Render takes a Markdown string and returns ANSI-styled text
// suitable for printing into a gocui.View.
func Render(md string, width int) (string, error) {
	if width <= 0 {
		width = 80
	}

	r, err := glamour.NewTermRenderer(
		// glamour.WithAutoStyle(),
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}

	return r.Render(md)
}
