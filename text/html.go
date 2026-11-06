package text

import (
	"fmt"
	"strings"
)

const (
	sectionStart = "<section>"
	sectionEnd   = "</section>"
	articleStart = "<article>"
	articleEnd   = "</article>"
)

type Transformer struct {
	source []byte
}

var (
	marginNoteCount int
	sideNoteCount   int
)

func (t Transformer) Article(tokens []Token) (string, error) {
	var builder strings.Builder
	marginNoteCount = 0
	sideNoteCount = 0

	builder.WriteString(articleStart)
	for _, token := range tokens {
		err := t.toHTML(&builder, token)
		if err != nil {
			return "", err
		}
	}

	builder.WriteString(sectionEnd)
	builder.WriteString(articleEnd)
	return builder.String(), nil
}

// I need a way to decide when to write section start and end
func (t Transformer) toHTML(builder *strings.Builder, token Token) error {
	switch v := token.(type) {
	case *Header:
		if v.Level == header2 {
			builder.WriteString(sectionStart)
		}
		fmt.Fprintf(builder, "<h%d>", int(v.Level))
		builder.Write(t.source[v.Start():v.End()])
		fmt.Fprintf(builder, "</h%d>", int(v.Level))
	case *Subheader:
		builder.WriteString(`<p class="subtitle">`)
		builder.Write(t.source[v.Start():v.End()])
		builder.WriteString("</p>")
	case *MarginNote:
		fmt.Fprintf(builder, `<label for="mn-%d" class="margin-toggle">⊕</label>
				<input type="checkbox" id="mn-%d" class="margin-toggle"/><span class="marginnote">`, marginNoteCount, marginNoteCount)
		children := v.GetChildren()
		for _, c := range children {
			err := t.toHTML(builder, c)
			if err != nil {
				return err
			}
		}
		builder.WriteString("</span>")
		marginNoteCount++
	case *SideNote:
		fmt.Fprintf(builder, `<label for="sn-%d" class="margin-toggle sidenote-number">
				<input type="checkbox" id="mn-%d" class="margin-toggle"/><span class="sidenote">`, sideNoteCount, sideNoteCount)
		children := v.GetChildren()
		for _, c := range children {
			err := t.toHTML(builder, c)
			if err != nil {
				return err
			}
		}
		builder.WriteString("</span>")
		sideNoteCount++
	case *Link:
		builder.WriteString(`<a href="`)
		builder.Write(t.source[v.Start():v.End()])
		fmt.Fprintf(builder, `">%s</a>`, v.Name)
	case *Image:
		builder.WriteString(`<figure><img src="`)
		builder.Write(t.source[v.Start():v.End()])
		builder.WriteString(`"/>`)
		if len(v.Caption) > 0 {
			builder.Write(v.Caption)
		}
		builder.WriteString("</figure>")
	case *InlineCode:
		builder.WriteString("<code>")
		builder.Write(t.source[v.Start():v.End()])
		builder.WriteString("</code>")
	case *BlockCode:
		builder.WriteString(`<pre><code class="language-`)
		builder.Write(v.Language)
		builder.WriteString(`">'`)
		builder.Write(t.source[v.Start():v.End()])
		builder.WriteString("</code></pre>")
	case *Bold:
		builder.WriteString("<b>")
		builder.Write(t.source[v.Start():v.End()])
		builder.WriteString("</b>")
	case *Italic:
		builder.WriteString("<em>")
		builder.Write(t.source[v.Start():v.End()])
		builder.WriteString("</em>")
	case *Text:
		builder.Write(t.source[v.Start():v.End()])
	case *LineFeed:
		builder.Write([]byte{lineFeed})
	default:
		return fmt.Errorf("not valid token %s", v.String())
	}
	return nil
}
