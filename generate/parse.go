package generate

import (
	"regexp"

	"github.com/aoldershaw/regen"
)

var headerPrefix = regen.Sequence(
	regen.LineStart,
	regen.Sequence(
		regen.String("#").Repeat().Min(1).Max(2),
		regen.Whitespace.Repeat().Min(1),
	),
)

var releaseNoteHeader = regen.Sequence(
	headerPrefix,
	regen.String("Release Note"),
	regen.String("s").Optional(),
	regen.Whitespace.Repeat(),
	regen.CharSet('\n', '\r').Repeat().Min(1),
).Group().NoCapture().SetFlags(regen.FlagCaseInsensitive)

var releaseNoteContent = regen.Any.Repeat().Ungreedy().Group().SetFlags(regen.FlagMatchNewLine)

var releaseNoteRegexp = regexp.MustCompile(regen.Sequence(
	releaseNoteHeader,
	releaseNoteContent.Capture(),
	regen.Whitespace.Repeat(),
	regen.OneOf(
		regen.TextEnd,
		headerPrefix,
	).Group().NoCapture(),
).Group().NoCapture().SetFlags(regen.FlagMultiLine).Regexp())

func parseReleaseNote(body string) string {
	groups := releaseNoteRegexp.FindStringSubmatch(body)
	if len(groups) < 2 {
		return ""
	}
	return groups[1]
}
