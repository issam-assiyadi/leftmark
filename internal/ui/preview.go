package ui

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/awesome-gocui/gocui"

	"github.com/issam-assiyadi/leftmark/domain"
)

const previewStyleName = "monokai"

const gutterStyle = "\x1b[38;5;244m"

func (a *App) previewItem(item domain.Item) error {
	data, err := a.Service.ReadFile(item.File)
	if err != nil {
		return err
	}

	a.previewLines = highlightLines(item.File, data)
	focusLine := item.Line - 1

	a.previewFocusText = ""
	if plainLines := strings.Split(string(data), "\n"); focusLine >= 0 && focusLine < len(plainLines) {
		a.previewFocusText = plainLines[focusLine]
	}

	a.focused = a.Preview.Open(a.focused, fmt.Sprintf(" %s ", item.File), focusLine)
	return nil
}

func previewGutterWidth(lineCount int) int {
	return len(strconv.Itoa(lineCount))
}

func (a *App) openSelectedInPreview(g *gocui.Gui, v *gocui.View) error {
	item, ok := a.selectedItem()
	if !ok {
		return nil
	}

	if err := a.previewItem(item); err != nil {
		log.Println("preview: read file:", err)
		return nil
	}

	if g == nil {
		return nil
	}

	if err := a.layout(g); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(a.focused); err != nil {
		log.Println("preview: unable to focus view:", err)
	}
	return nil
}

func (a *App) closePreview(g *gocui.Gui, v *gocui.View) error {
	if !a.Preview.IsOpen() {
		return nil
	}

	a.previewLines = nil
	a.focused = a.Preview.Close(g)

	if g == nil {
		return nil
	}
	if _, err := g.SetCurrentView(a.focused); err != nil {
		log.Println("preview: unable to restore focus:", err)
	}
	return nil
}

func highlightLines(path string, data []byte) []string {
	lexer := lexers.Match(path)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	var lines []string
	iterator, err := lexer.Tokenise(nil, string(data))
	if err != nil {
		lines = strings.Split(string(data), "\n")
	} else {
		var buf bytes.Buffer
		if err := formatters.TTY256.Format(&buf, styles.Get(previewStyleName), iterator); err != nil {
			lines = strings.Split(string(data), "\n")
		} else {
			lines = strings.Split(buf.String(), "\n")
		}
	}

	return addLineNumbers(lines)
}

// addLineNumbers prefixes each line with a right-aligned, dimmed line
// number and gutter separator, and resets color state at the start of
// every line. The reset matters even for lines with no gutter-related
// color of their own: chroma can leave a color open across an embedded
// newline inside a multi-line token (e.g. a block comment), and gocui
// parses escapes per rendered line, so without it a color can bleed into
// the next visual line.
func addLineNumbers(lines []string) []string {
	width := previewGutterWidth(len(lines))
	numbered := make([]string, len(lines))
	for i, line := range lines {
		numbered[i] = fmt.Sprintf("%s%*d │ %s%s", gutterStyle, width, i+1, rowStyleReset, line)
	}
	return numbered
}
