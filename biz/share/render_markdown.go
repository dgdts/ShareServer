package share

import (
	"bytes"
	"context"
	"html/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

const (
	ShareNoteTemplatePath = "../../template/share_note.html"
)

func RenderMarkdown(ctx context.Context, markdownRawData []byte) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	html := markdown.ToHTML(markdownRawData, parser, nil)

	return template.HTML(html)
}

func GenerateShareNoteHTML(ctx context.Context, note *MarkdownNoteData, title string, userName string) ([]byte, error) {
	html := RenderMarkdown(ctx, []byte(note.Content))

	tmpl, err := template.ParseFiles(ShareNoteTemplatePath)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, map[string]interface{}{
		"Title":    title,
		"Content":  html,
		"UserName": userName,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
