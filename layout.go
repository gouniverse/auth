package auth

import "github.com/gouniverse/hb"

func (a Auth) layout(content string) string {
	h := hb.NewSection().
		Attr("style", "padding:80px 0px").
		AddChild(hb.NewHTML(content))
	return h.ToHTML()
}
