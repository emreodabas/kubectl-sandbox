package drawing

import (
	"../lessons"
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	"strings"
	"sync"
	"time"
)

type lessonDrawer struct {
	termbox    iTermbox
	rwMutex    sync.RWMutex
	drawTimer  *time.Timer
	lesson     lessons.Lesson
	page       int
	lessonType int
}

var (
	descDraw        = &lessonDrawer{}
	ErrAbort        = errors.New("abort")
	errNext         = errors.New("entered")
	errBack         = errors.New("go back")
	descScreenRatio = 6
)

func ShowLesson(lesson lessons.Lesson, lessonType int, page int) error {
	descDraw.lesson = lesson
	descDraw.page = page
	descDraw.lessonType = lessonType
	//if lessonType == lessons.Desc {
	return descDraw.showDescriptions(lesson.Descriptions, page)
	//}else if lessonType == lessons.Quiz {

	//}

}

func (drawer *lessonDrawer) showDescriptions(cards []lessons.DescriptionCard, i int) error {

	err := descDraw.showDescription(cards[i])

	switch {
	case err == errBack:
		fmt.Println("go back")
		descDraw.showDescription(cards[i-1])
	case err == errNext:
		fmt.Println("go next")
		descDraw.showDescription(cards[i+1])
	default:
		return err
	}
	return nil
}

func (drawer *lessonDrawer) showDescription(card lessons.DescriptionCard) error {

	if err := drawer.initDrawer("description"); err != nil {
		return err
	}

	defer drawer.termbox.close()

	for {
		drawer.draw(10 * time.Millisecond)
		//TODO define errors
		err := drawer.readKey()
		switch {
		case err == ErrAbort:
			return err
		case err != nil:
			return nil
		}
	}
}

func (drawer *lessonDrawer) initDrawer(dataType string) error {

	if drawer.termbox == nil {
		drawer.termbox = &termImpl{}
	}

	if err := drawer.termbox.init(); err != nil {
		return errors.Wrap(err, "failed to initialize termbox")
	}
	switch dataType {

	case "description":
		drawer.drawTimer = time.AfterFunc(0, func() {
			drawer.drawDescription()
			drawer.termbox.flush()
		})

	case "quiz":
		drawer.drawTimer = time.AfterFunc(0, func() {
			drawer.drawQuiz()
			drawer.termbox.flush()
		})
	case "interactive":
		drawer.drawTimer = time.AfterFunc(0, func() {
			drawer.drawInteractive()
			drawer.termbox.flush()
		})
	}
	drawer.drawTimer.Stop()
	return nil
}
func (drawer *lessonDrawer) drawDescription() {

	const pipeline = '│'
	const backArrow = '<'
	const nextArrow = '>'
	const arrowBody = '='

	width, height := drawer.termbox.size()
	drawer.termbox.clear(termbox.ColorDefault, termbox.ColorDefault)
	card := drawer.getCard()
	//fmt.Println("\n\t\t\t"+card.LessonHeader+"\n\n\n\t"+card.Header+"\t\n\n"+card.Data, "\n")
	sp := strings.Split("\n\t\t\t"+card.LessonHeader+"\n\n\n\t"+card.Header+"\t\n\n"+card.Data, "\n")
	prevLines := make([][]rune, 0, len(sp))
	for _, s := range sp {
		prevLines = append(prevLines, []rune(s))
	}

	descSquareRight := width - (width / descScreenRatio)
	descSquareLeft := (width / descScreenRatio)
	descSquareTop := height // - (height / 5)
	descSquareBottom := 1   //(height / 5)
	backArrowLeft := 5
	backArrowRight := descSquareLeft - 5
	nextArrowLeft := descSquareRight + 5
	nextArrowRight := width - 5
	arrowHeight := height - (height / descScreenRatio)

	if isBackArrowExist() {
		backValue := "  " + getBackValue() + " "
		for i := backArrowLeft; i < backArrowRight; i++ {
			switch {
			case i == backArrowLeft:
				drawer.termbox.setCell(i, arrowHeight, backArrow, termbox.ColorBlack, termbox.ColorDefault)
			case i == backArrowLeft+1, i == backArrowLeft+2:
				drawer.termbox.setCell(i, arrowHeight, arrowBody, termbox.ColorBlack, termbox.ColorDefault)
			default:
				if i-backArrowLeft-3 > 0 {
					if i-backArrowLeft-3 < len(backValue) {
						drawer.termbox.setCell(i, arrowHeight, rune(backValue[i-backArrowLeft-3]), termbox.ColorBlack, termbox.ColorDefault)
					} else {
						drawer.termbox.setCell(i, arrowHeight, arrowBody, termbox.ColorBlack, termbox.ColorDefault)
					}
				}

			}
		}

	}
	if isNextArrowExist() {
		nextValue := "  " + getNextValue() + " "
		for i := nextArrowLeft; i < nextArrowRight; i++ {
			switch {
			case i == nextArrowRight-1:
				drawer.termbox.setCell(i, arrowHeight, nextArrow, termbox.ColorBlack, termbox.ColorDefault)
			case i == nextArrowRight-2, i == nextArrowRight-3:
				drawer.termbox.setCell(i, arrowHeight, arrowBody, termbox.ColorBlack, termbox.ColorDefault)
			default:
				if nextArrowRight-i-3 > 0 {
					if nextArrowRight-i-3 < len(nextValue) {
						drawer.termbox.setCell(i, arrowHeight, rune(nextValue[len(nextValue)-(nextArrowRight-i-3)]), termbox.ColorBlack, termbox.ColorDefault)
					} else {
						drawer.termbox.setCell(i, arrowHeight, arrowBody, termbox.ColorBlack, termbox.ColorDefault)
					}
				}

			}
		}

	}

	//case i == nextArrowRight:
	//drawer.termbox.setCell(i, h, nextArrow, termbox.ColorBlack, termbox.ColorDefault)

	// top line
	for i := descSquareLeft; i < descSquareRight; i++ {
		var r rune
		if i == descSquareLeft {
			r = '┌'
		} else if i == descSquareRight-1 {
			r = '┐'
		} else {
			r = '─'
		}
		drawer.termbox.setCell(i, 0, r, termbox.ColorBlack, termbox.ColorDefault)
	}
	// bottom line
	for i := descSquareLeft; i < descSquareRight; i++ {
		var r rune
		if i == descSquareLeft {
			r = '└'
		} else if i == descSquareRight-1 {
			r = '┘'
		} else {
			r = '─'
		}
		drawer.termbox.setCell(i, descSquareTop-1, r, termbox.ColorBlack, termbox.ColorDefault)
	}

	var wvline = runewidth.RuneWidth(pipeline)
	for h := descSquareBottom; h < descSquareTop-1; h++ {
		w := descSquareLeft
		for i := descSquareLeft; i < descSquareRight; i++ {
			switch {

			// Box Left line
			case i == descSquareLeft:
				drawer.termbox.setCell(i, h, pipeline, termbox.ColorBlack, termbox.ColorDefault)
				w += wvline
				// Box Right line
			case i == descSquareRight-1:
				drawer.termbox.setCell(i, h, pipeline, termbox.ColorBlack, termbox.ColorDefault)
				w += wvline
				// Box left right indentation
			case w == descSquareLeft+wvline, w == descSquareRight-1-wvline:
				drawer.termbox.setCell(w, h, ' ', termbox.ColorDefault, termbox.ColorDefault)
				w++

			default:
				if h-1 >= len(prevLines) {
					w++
					continue
				}
				j := i - descSquareLeft - 2 // Two spaces.
				l := prevLines[h-1]
				if j >= len(l) {
					w++
					continue
				}
				rw := runewidth.RuneWidth(l[j])
				if w+rw > descSquareRight-1-2 {
					drawer.termbox.setCell(w, h, '.', termbox.ColorDefault, termbox.ColorDefault)
					drawer.termbox.setCell(w+1, h, '.', termbox.ColorDefault, termbox.ColorDefault)
					w += 2
					continue
				}

				drawer.termbox.setCell(w, h, l[j], termbox.ColorDefault, termbox.ColorDefault)
				w += rw
			}
		}
	}

}
func getBackValue() string {
	if descDraw.lessonType == lessons.Desc {
		return descDraw.lesson.Descriptions[descDraw.page-1].Header
	} else if descDraw.lessonType == lessons.Quiz {
		return string(descDraw.page-1) + ". Question"
	} else {
		return string(descDraw.page-1) + ". Lab"
	}
}
func getNextValue() string {
	if descDraw.lessonType == lessons.Desc {
		return descDraw.lesson.Descriptions[descDraw.page+1].Header
	} else if descDraw.lessonType == lessons.Quiz {
		return string(descDraw.page+1) + ". Question"
	} else {
		return string(descDraw.page+1) + ". Lab"
	}
}

//TODO
func isNextArrowExist() bool {
	if descDraw.lessonType == lessons.Desc {
		return descDraw.page+1 < len(descDraw.lesson.Descriptions)
	} else if descDraw.lessonType == lessons.Quiz {
		return descDraw.page+1 < len(descDraw.lesson.Quiz)
	} else {
		return descDraw.page+1 < len(descDraw.lesson.InteractiveActions)
	}
}

func isBackArrowExist() bool {
	return descDraw.page > 0
}

func (drawer *lessonDrawer) drawQuiz() {

}
func (drawer *lessonDrawer) drawInteractive() {

}

func (drawer *lessonDrawer) draw(duration time.Duration) {
	drawer.rwMutex.RLock()
	defer drawer.rwMutex.RUnlock()
	drawer.drawTimer.Reset(duration)
}

func (drawer *lessonDrawer) readKey() error {
	switch e := drawer.termbox.pollEvent(); e.Type {
	case termbox.EventKey:
		switch e.Key {
		case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlD:
			return ErrAbort
		case termbox.KeyBackspace, termbox.KeyBackspace2:

		case termbox.KeyDelete:

		case termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyCtrlB:
			fmt.Println("keypressed for back")
			isBackArrowExist()
			{
				fmt.Println("keypressed for back and is back arrow exist")
				return errBack
			}
		case termbox.KeyArrowUp, termbox.KeyArrowRight, termbox.KeyCtrlF, termbox.KeyEnter:
			isNextArrowExist()
			{
				return errNext
			}
		case termbox.KeyCtrlA:

		case termbox.KeyCtrlE:

		case termbox.KeyCtrlW:
		case termbox.KeyCtrlU:
		case termbox.KeyCtrlK, termbox.KeyCtrlP:
		case termbox.KeyCtrlJ, termbox.KeyCtrlN:
		case termbox.KeyTab:
		default:
			//fmt.Println(e.Ch)
		}
	case termbox.EventResize:
		// To get actual window size, clear all buffers.
		// See termbox.Clear's documentation for more details.
		drawer.termbox.clear(termbox.ColorDefault, termbox.ColorDefault)
		drawer.draw(200 * time.Millisecond)
	}
	return nil
}
func (drawer *lessonDrawer) getCard() lessons.DescriptionCard {
	return descDraw.lesson.Descriptions[descDraw.page]
}

func goBack() {

}

func goForward() {

}
