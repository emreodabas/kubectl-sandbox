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
	lessonType int
	page       int
	lenQuiz    int
	lenDesc    int
	lenLab     int
}

var (
	drawer          = &lessonDrawer{}
	ErrAbort        = errors.New("abort")
	errNext         = errors.New("entered")
	errBack         = errors.New("go back")
	errDoNothing    = errors.New("No change")
	descScreenRatio = 6
)

func ShowLesson(lesson lessons.Lesson, lessonType int, page int) error {
	drawer.lesson = lesson
	drawer.page = page
	drawer.lessonType = lessonType
	drawer.lenDesc = len(lesson.Descriptions)
	drawer.lenQuiz = len(lesson.Quiz)
	drawer.lenLab = len(lesson.Labs)
	return drawer.showLesson()
	//if lessonType == lessons.Desc {

}

//func (drawer *lessonDrawer) showLesson(lesson lessons.Lesson, lessonType int, page int) error {

//lesson navigation
//	err := drawer.showDescription(lesson, )
//
//	//switch {
//	//case err == errBack:
//	//	fmt.Println("go back")
//	//	drawer.page = i - 1
//	//	drawer.showDescription()
//	//case err == errNext:
//	//	fmt.Println("go next")
//	//	drawer.page = i + 1
//	//	drawer.showDescription()
//	//default:
//	//	fmt.Println(err.Error())
//	//	return err
//	//}
//	return nil
//}

func (drawer *lessonDrawer) showLesson() error {

	if err := drawer.initDrawer(); err != nil {
		return err
	}

	defer drawer.termbox.close()

	for {
		drawer.draw(10 * time.Millisecond)
		//TODO define errors
		err := drawer.readKey()
		switch {
		case err == ErrAbort:
			return ErrAbort
		case err == errBack:
			goBack()
		case err == errNext:
			goNext()
		case err == errDoNothing:
			drawer.showLesson()
		case err != nil:
			return fmt.Errorf(err.Error(), "failed to read a key")
		default:
			return nil
		}
	}
}

func (drawer *lessonDrawer) initDrawer() error {

	if drawer.termbox == nil {
		drawer.termbox = &termImpl{}
	}

	if err := drawer.termbox.init(); err != nil {
		return errors.Wrap(err, "failed to initialize termbox")
	}

	switch drawer.lessonType {

	case lessons.Desc:
		drawer.drawTimer = time.AfterFunc(0, func() {
			drawer.drawDescription()
			drawer.termbox.flush()
		})

	case lessons.Quiz:
		drawer.drawTimer = time.AfterFunc(0, func() {
			drawer.drawQuiz()
			drawer.termbox.flush()
		})
	case lessons.Interactive:
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
	card := drawer.lesson.Descriptions[drawer.page]
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
				drawer.termbox.setCell(i, arrowHeight, backArrow, termbox.ColorDefault, termbox.ColorDefault)
			case i == backArrowLeft+1, i == backArrowLeft+2:
				drawer.termbox.setCell(i, arrowHeight, arrowBody, termbox.ColorDefault, termbox.ColorDefault)
			default:
				if i-backArrowLeft-3 > 0 {
					if i-backArrowLeft-3 < len(backValue) {
						drawer.termbox.setCell(i, arrowHeight, rune(backValue[i-backArrowLeft-3]), termbox.ColorDefault, termbox.ColorDefault)
					} else {
						drawer.termbox.setCell(i, arrowHeight, arrowBody, termbox.ColorDefault, termbox.ColorDefault)
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
				drawer.termbox.setCell(i, arrowHeight, nextArrow, termbox.ColorDefault, termbox.ColorDefault)
			case i == nextArrowRight-2, i == nextArrowRight-3:
				drawer.termbox.setCell(i, arrowHeight, arrowBody, termbox.ColorDefault, termbox.ColorDefault)
			default:
				if nextArrowRight-i-3 > 0 {
					if nextArrowRight-i-3 < len(nextValue) {
						drawer.termbox.setCell(i, arrowHeight, rune(nextValue[len(nextValue)-(nextArrowRight-i-3)]), termbox.ColorDefault, termbox.ColorDefault)
					} else {
						drawer.termbox.setCell(i, arrowHeight, arrowBody, termbox.ColorDefault, termbox.ColorDefault)
					}
				}

			}
		}
	}
	//case i == nextArrowRight:
	//drawer.termbox.setCell(i, h, nextArrow, termbox.ColorDefault, termbox.ColorDefault)

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
		drawer.termbox.setCell(i, 0, r, termbox.ColorDefault, termbox.ColorDefault)
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
		drawer.termbox.setCell(i, descSquareTop-1, r, termbox.ColorDefault, termbox.ColorDefault)
	}

	var wvline = runewidth.RuneWidth(pipeline)
	for h := descSquareBottom; h < descSquareTop-1; h++ {
		w := descSquareLeft
		for i := descSquareLeft; i < descSquareRight; i++ {
			switch {

			// Box Left line
			case i == descSquareLeft:
				drawer.termbox.setCell(i, h, pipeline, termbox.ColorDefault, termbox.ColorDefault)
				w += wvline
				// Box Right line
			case i == descSquareRight-1:
				drawer.termbox.setCell(i, h, pipeline, termbox.ColorDefault, termbox.ColorDefault)
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
			if isBackArrowExist() {
				return errBack
			} else {
				return errDoNothing
			}
		case termbox.KeyArrowUp, termbox.KeyArrowRight, termbox.KeyCtrlF, termbox.KeyEnter:
			if isNextArrowExist() {
				return errNext
			} else {
				return errDoNothing
			}
		case termbox.KeyCtrlA:

		case termbox.KeyCtrlE:

		case termbox.KeyCtrlW:
		case termbox.KeyCtrlU:
		case termbox.KeyCtrlK, termbox.KeyCtrlP:
		case termbox.KeyCtrlJ, termbox.KeyCtrlN:
		case termbox.KeyTab:
		default:
			return fmt.Errorf(string(e.Ch) + "Pressed and exit")
		}
	case termbox.EventResize:
		// To get actual window size, clear all buffers.
		// See termbox.Clear's documentation for more details.
		drawer.termbox.clear(termbox.ColorDefault, termbox.ColorDefault)
		drawer.draw(200 * time.Millisecond)
	}
	return nil
}
func (drawer *lessonDrawer) showEndPage() {
	//TODO
	fmt.Println("congratulations")
}
func (drawer *lessonDrawer) drawInteractive() {

}
func (drawer *lessonDrawer) drawQuiz() {

}

func goBack() {
	if drawer.lessonType == lessons.Desc {
		iteratePage(-1)
		drawer.showLesson()
	} else if drawer.lessonType == lessons.Interactive {
		if drawer.page == 0 {
			setPage(drawer.lenDesc - 1)
			changeLessonType(lessons.Desc)
			drawer.showLesson()
		} else {
			iteratePage(-1)
			drawer.showLesson()
		}
	} else if drawer.lessonType == lessons.Quiz {
		if drawer.page == 0 {
			setPage(drawer.lenLab - 1)
			changeLessonType(lessons.Interactive)
			drawer.showLesson()
		} else {
			iteratePage(-1)
			drawer.showLesson()
		}
	}
}

func goNext() {
	if drawer.lessonType == lessons.Desc {
		if drawer.lenDesc == drawer.page+1 {
			changeLessonType(lessons.Interactive)
			setPage(0)
			drawer.showLesson()
		} else {
			iteratePage(1)
			drawer.showLesson()
		}
	} else if drawer.lessonType == lessons.Interactive {
		if drawer.lenLab == drawer.page+1 {
			changeLessonType(lessons.Quiz)
			setPage(0)
			drawer.showLesson()
		} else {
			iteratePage(1)
			drawer.showLesson()
		}
	} else if drawer.lessonType == lessons.Quiz {
		if drawer.lenQuiz == drawer.page+1 {
			drawer.showEndPage()
		} else {
			iteratePage(1)
			drawer.showLesson()
		}
	}
}

func iteratePage(value int) {
	drawer.page = drawer.page + value
}
func setPage(value int) {
	drawer.page = value
}
func changeLessonType(value int) {
	drawer.lessonType = value
}

//TODO
func isNextArrowExist() bool {
	if drawer.lessonType == lessons.Desc {
		return !(drawer.lenDesc == drawer.page+1 && drawer.lenLab == 0 && drawer.lenQuiz == 0)
	} else if drawer.lessonType == lessons.Interactive {
		return !(drawer.lenLab == drawer.page+1 && drawer.lesson.Quiz != nil)
	} else if drawer.lessonType == lessons.Quiz {
		return !(drawer.lenQuiz == drawer.page+1)
	}
	return false
}

func isBackArrowExist() bool {
	if drawer.lessonType == lessons.Desc {
		return !(drawer.page == 0)
	} else if drawer.lessonType == lessons.Interactive {
		return !(drawer.page == 0 && drawer.lesson.Descriptions != nil)
	} else if drawer.lessonType == lessons.Quiz {
		return !(drawer.page == 0 && drawer.lesson.Labs != nil && drawer.lesson.Descriptions != nil)
	}
	return false
}

func getBackValue() string {
	if drawer.lessonType == lessons.Desc {
		return drawer.lesson.Descriptions[drawer.page-1].Header
	} else if drawer.lessonType == lessons.Interactive {
		if drawer.page == 0 && drawer.lesson.Descriptions != nil {
			return drawer.lesson.Descriptions[drawer.lenDesc-1].Header
		} else {
			return string(drawer.page-1) + ". Lab"
		}

	} else if drawer.lessonType == lessons.Quiz {
		if drawer.page == 0 && drawer.lesson.Labs != nil {
			return string(drawer.lenLab) + ". Lab"
		} else if drawer.page == 0 && drawer.lesson.Descriptions != nil {
			return drawer.lesson.Descriptions[drawer.lenDesc-1].Header
		} else {
			return string(drawer.page-1) + ". Question"
		}
	} else {
		return "ERROR -- getBackValue :)"
	}
}
func getNextValue() string {
	if drawer.lessonType == lessons.Desc {
		if drawer.lenDesc == drawer.page+1 && drawer.lesson.Labs != nil {
			return "1. Lab"
		} else if drawer.lenDesc == drawer.page+1 && drawer.lesson.Quiz != nil {
			return "1. Question"
		} else {
			return drawer.lesson.Descriptions[drawer.page+1].Header
		}

	} else if drawer.lessonType == lessons.Interactive {
		if drawer.lenLab == drawer.page+1 && drawer.lesson.Quiz != nil {
			return "1. Question"
		} else {
			return string(drawer.page+1) + ". Lab"
		}
	} else if drawer.lessonType == lessons.Quiz {
		return string(drawer.page+1) + ". Question"
	} else {
		return "getNextValue"
	}
}
