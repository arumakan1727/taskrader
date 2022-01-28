package view

import (
	"fmt"
	"io"
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
	Text  string // 行の文字列
	Width int    // 表示幅 (e.g. "abc" なら半角3文字なので3, "あ"なら2)
}

type cardContent struct {
	lines []lineInfo
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func Show(ass []*assignment.Assignment, writer io.Writer) {
	if len(ass) == 0 {
		color.Green("未提出の課題はありません")
		return
	}

	cards := make([]cardContent, len(ass))

	maxWidth := 0
	for i, a := range ass {
		l1 := fmt.Sprintf("[%02d] %s 【%s】", i+1, a.Title, a.Course)
		l2 := fmt.Sprintf("     %s     @%s", FmtDeadline(&a.Deadline), a.Origin)
		w1 := CalcTextWidth(l1, 1)
		w2 := CalcTextWidth(l2, 1)
		cards[i] = cardContent{
			lines: []lineInfo{{l1, w1}, {l2, w2}},
		}
		maxWidth = max(maxWidth, w1)
		maxWidth = max(maxWidth, w2)
	}

	// 最低でも40以上確保 || 末尾の余白は5欲しい
	maxWidth = max(40, maxWidth+5)

	writeBorderTop(writer, maxWidth+2)

	bar := " " + borderVH + strings.Repeat(borderHori, maxWidth+2) + borderHV + "\n"

	for i, card := range cards {
		for _, l := range card.lines {
			io.WriteString(writer, " "+borderVert+" ")
			io.WriteString(writer, l.Text)
			io.WriteString(writer, strings.Repeat(" ", maxWidth-l.Width+1))
			fmt.Fprint(writer, borderVert+"\n")
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

func FmtDeadline(t *time.Time) string {
	if t.Equal(assignment.UnknownDeadline()) {
		return "<締切不明>"
	}
	const layoutDay = "2006/01/02"
	const layoutTime = "15:04"
	wdays := []string{"日", "月", "火", "水", "木", "金", "土"}
	return t.Format(layoutDay) + " (" + wdays[t.Weekday()] + ") " + t.Format(layoutTime)
}

func printWithColor(output string, deadline time.Time) {
	if diff := time.Until(deadline).Hours(); 0 <= diff && diff <= 24 {
		fmt.Println(color.RedString(output))
	} else if 24 < diff && diff <= 48 {
		fmt.Println(color.YellowString(output))
	} else {
		fmt.Println(output)
	}
}
