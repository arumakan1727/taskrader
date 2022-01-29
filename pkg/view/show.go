package view

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/arumakan1727/taskrader/pkg/assignment"
	"github.com/fatih/color"
	"golang.org/x/text/width"
)

const (
	borderTopL = "╭"
	borderTopR = "╮"
	borderBotL = "╰"
	borderBotR = "╯"
	borderVert = "│"
	borderHori = "─"
	borderVH   = "├"
	borderHV   = "┤"
)

type lineInfo struct {
	ltext  string // 左寄せにする文字列
	rtext  string // 右寄せにする文字列
	lwidth int    // 左寄せ文字列の表示幅
	rwidth int    // 右寄せ文字列の表示幅
}

type cardContent struct {
	lines []lineInfo
}

func (l *lineInfo) totalWidth() int {
	return l.lwidth + l.rwidth
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// 締切の昇順にソートする。
// ただし、UnknownDeadline はnow+(48時間1ナノ秒)として扱う (色付きで警告表示される行の直後に来るようにする)
func SortAssignments(ass []*assignment.Assignment) {
	if len(ass) == 0 {
		return
	}

	afterTomorrow := time.Now().Add(48*time.Hour + 1)

	sort.SliceStable(ass, func(i, j int) bool {
		ti := ass[i].Deadline
		tj := ass[j].Deadline
		if ti.Equal(assignment.UnknownDeadline()) {
			ti = afterTomorrow
		}
		if tj.Equal(assignment.UnknownDeadline()) {
			tj = afterTomorrow
		}
		return ti.Before(tj)
	})
}

func Show(ass []*assignment.Assignment, writer io.Writer) {
	if len(ass) == 0 {
		color.Green("未提出の課題はありません")
		return
	}
	cOrd := color.New(color.FgWhite).SprintfFunc()
	cTitle := color.New(color.FgCyan, color.Bold).SprintfFunc()
	cCourse := color.New(color.FgCyan, color.Bold).SprintfFunc()
	cOrig := color.New(color.FgMagenta).SprintfFunc()

	cards := make([]cardContent, len(ass))

	maxWidth := 0
	for i, a := range ass {
		l1 := lineInfo{
			ltext:  cOrd("%02d  ", i+1) + cTitle("%s", a.Title) + cCourse(" 【%s】", a.Course),
			rtext:  "",
			lwidth: CalcTextWidth(fmt.Sprintf("%02d  %s 【%s】", i+1, a.Title, a.Course), 1),
			rwidth: 0,
		}
		l2 := lineInfo{
			ltext:  "",
			rtext:  "     " + fmtDeadline(a.Deadline, true) + " " + cOrig("%s", fmtOrig(a.Origin)),
			lwidth: 0,
			rwidth: CalcTextWidth(fmt.Sprintf("     %s %s", fmtDeadline(a.Deadline, false), fmtOrig(a.Origin)), 1),
		}
		cards[i] = cardContent{
			lines: []lineInfo{l1, l2},
		}
		maxWidth = max(maxWidth, l1.totalWidth())
		maxWidth = max(maxWidth, l2.totalWidth())
	}

	// 最低でも40以上確保
	maxWidth = max(40, maxWidth)

	writeBorderTop(writer, maxWidth+2)

	bar := " " + borderVH + strings.Repeat(borderHori, maxWidth+2) + borderHV + "\n"

	for i, card := range cards {
		for _, l := range card.lines {
			io.WriteString(writer, " "+borderVert+" ")
			io.WriteString(writer, l.ltext)
			io.WriteString(writer, strings.Repeat(" ", maxWidth-l.totalWidth()))
			io.WriteString(writer, l.rtext)
			fmt.Fprint(writer, " "+borderVert+"\n")
		}

		if i < len(cards)-1 {
			io.WriteString(writer, bar)
		}
	}
	writeBorderBottom(writer, maxWidth+2)
}

func writeBorderTop(w io.Writer, width int) {
	io.WriteString(w, " "+borderTopL)
	io.WriteString(w, strings.Repeat(borderHori, width))
	io.WriteString(w, borderTopR+"\n")
}

func writeBorderBottom(w io.Writer, width int) {
	io.WriteString(w, " "+borderBotL)
	io.WriteString(w, strings.Repeat(borderHori, width))
	io.WriteString(w, borderBotR+"\n")
}

// 全角なら2, 半角なら1として各文字の表示幅を求め、その総和を返す。
// '①' や '▼' のような曖昧幅の文字は引数 `ambigousWidth` の値とする。
func CalcTextWidth(text string, ambigousWidth int) int {
	wsum := 0

	for _, chr := range text {
		kind := width.LookupRune(chr).Kind()

		switch kind {
		case width.Neutral, width.EastAsianNarrow, width.EastAsianHalfwidth:
			wsum += 1
		case width.EastAsianWide, width.EastAsianFullwidth:
			wsum += 2
		case width.EastAsianAmbiguous:
			wsum += ambigousWidth
		default:
			panic("Unknown Kind of rune width: " + kind.String())
		}
	}
	return wsum
}

func fmtOrig(orig assignment.Origin) string {
	switch orig {
	case assignment.OrigGakujo:
		return "  @学情"
	case assignment.OrigEdStem:
		return "@EdStem"
	case assignment.OrigTeams:
		return " @Teams"
	default:
		return ""
	}
}

var (
	decoratorOverDue  = color.New(color.BgHiRed, color.FgHiWhite, color.Bold)
	decolatorUntil24h = color.New(color.BgRed, color.FgWhite, color.Bold)
	decolatorUntil48h = color.New(color.BgYellow, color.FgBlack, color.Bold)
)

func fmtDeadline(t time.Time, enableColor bool) string {
	if t.Equal(assignment.UnknownDeadline()) {
		return "締切不明"
	}

	noColorMemo := color.NoColor
	color.NoColor = !enableColor
	defer func() {
		color.NoColor = noColorMemo
	}()

	layoutDay := "1月2日"
	if t.Year() != time.Now().Year() {
		layoutDay = "2006年1月2日"
	}
	const layoutTime = "15:04"
	wdays := []string{"日", "月", "火", "水", "木", "金", "土"}

	var res string

	// x 日の 00:XX は (x-1)日の 24:XX として表示
	if t.Hour() == 0 {
		xt := t.AddDate(0, 0, -1)
		res = xt.Format(layoutDay) + "(" + wdays[xt.Weekday()] + ") " + fmt.Sprintf("24:%02d", xt.Minute())
	} else {
		res = t.Format(layoutDay) + "(" + wdays[t.Weekday()] + ") " + t.Format(layoutTime)
	}

	bold := color.New(color.Bold)
	if diff := time.Until(t); diff < 0 {
		label := decoratorOverDue.Sprint("締切超過")
		return label + " " + bold.Add(color.FgRed).Sprint(res)
	} else if diff <= 24*time.Hour {
		label := decolatorUntil24h.Sprintf("残り%d時間", int(diff.Hours()))
		return label + " " + bold.Add(color.FgRed).Sprint(res)
	} else if diff <= 48*time.Hour {
		label := decolatorUntil48h.Sprintf("残り%d時間", int(diff.Hours()))
		return label + " " + bold.Add(color.FgYellow).Sprint(res)
	} else {
		return res
	}
}
