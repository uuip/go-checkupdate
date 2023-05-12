package rule

import "github.com/go-resty/resty/v2"

var cssRules = map[string]string{
	"PDF-XChange":       "#bh-history>li:first-of-type>a",
	"SecureCRT":         "#download-tabs>h4",
	"Registry Workshop": "p",
	"Firefox":           ".c-release-version",
	"Navicat[Mac]":      `.release-notes-table[platform="M"] td>.note-title`,
	"Navicat":           `.release-notes-table[platform="W"] td>.note-title`,
	"Everything":        "h2",
	"Python":            "p.download-buttons>a",
	"Contexts [Mac]":    ".section--history__item__header>h1",
	"WGestures 2":       "a#download:nth-of-type(1)",
	"WGestures 2 [Mac]": "a#download:nth-of-type(2)",
	"Git":               ".version",
	"AIDA64":            "td.version",
	"Beyond Compare":    ".hasicon",
}

var fnRules = map[string]func(resp *resty.Response) (string, error){
	"DevManView": parseDevManView,
	"FS Capture": parseFaststone,
	"FS Viewer":  parseFaststone,
	"Python":     parsePython,
	"VMware":     parseVmware,
	"WinRAR":     parseWinrar,
}
