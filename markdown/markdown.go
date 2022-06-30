package markdown

import (
	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

const (
	boldOn           = "\x1b[1m"
	boldOff          = "\x1b[21m"
	italicOn         = "\x1b[3m"
	italicOff        = "\x1b[23m"
	crossedOutOn     = "\x1b[9m"
	crossedOutOff    = "\x1b[29m"
	greenOn          = "\x1b[32m"
	resetAll         = "\x1b[0m"
	colorOff         = "\x1b[39m"
	colorBlockQuote  = "\x1b[38;5;10m"
	colorBlockQuote2 = "\x1b[38;5;13m"
	colorListItem    = "\x1b[38;5;13m"
	colorHeading     = "\x1b[38;5;11m"
)

// Extensions returns the bitmask of extensions supported by this renderer.
// The output of this function can be used to instantiate a new markdown
// parser using the `NewWithExtensions` function.
func Extensions() parser.Extensions {
	extensions := parser.NoIntraEmphasis        // Ignore emphasis markers inside words
	extensions |= parser.Tables                 // Parse tables
	extensions |= parser.FencedCode             // Parse fenced code blocks
	extensions |= parser.Autolink               // Detect embedded URLs that are not explicitly marked
	extensions |= parser.Strikethrough          // Strikethrough text using ~~test~~
	extensions |= parser.SpaceHeadings          // Be strict about prefix heading rules
	extensions |= parser.HeadingIDs             // specify heading IDs  with {#id}
	extensions |= parser.BackslashLineBreak     // Translate trailing backslashes into line breaks
	extensions |= parser.DefinitionLists        // Parse definition lists
	extensions |= parser.LaxHTMLBlocks          // more in HTMLBlock, less in HTMLSpan
	extensions |= parser.NoEmptyLineBeforeBlock // no need for new line before a list

	return extensions
}

func Render(source string, lineWidth int, leftPad int) []byte {
	p := parser.NewWithExtensions(Extensions())
	nodes := md.Parse([]byte(source), p)
	renderer := NewRenderer(lineWidth, leftPad)

	return md.Render(nodes, renderer)
}
