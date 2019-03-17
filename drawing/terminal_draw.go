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
	return descDraw.showLesson(lesson, lessonType, page)
	//}else if lessonType == lessons.Quiz {

	//}

}

//func (drawer *lessonDrawer) showLesson(lesson lessons.Lesson, lessonType int, page int) error {

//lesson navigation
//	err := descDraw.showDescription(lesson, )
//
//	//switch {
//	//case err == errBack:
//	//	fmt.Println("go back")
//	//	descDraw.page = i - 1
//	//	descDraw.showDescription()
//	//case err == errNext:
//	//	fmt.Println("go next")
//	//	descDraw.page = i + 1
//	//	descDraw.showDescription()
//	//default:
//	//	fmt.Println(err.Error())
//	//	return err
//	//}
//	return nil
//}

func (drawer *lessonDrawer) showLesson(lesson lessons.Lesson, dataType int, page int) error {

	if err := drawer.initDrawer(lesson, dataType, page); err != nil {
		return err
	}

	defer drawer.termbox.close()

	for {
		drawer.draw(10 * time.Millisecond)
		//TODO define errors
		err := drawer.readKey()
		switch {
		case err == ErrAbort:
		case err == errBack:
			goBack(lesson, dataType, page)
		case err == errNext:
			goNext(lesson, dataType, page)
		case err != nil:
			fmt.Println(err.Error())
			return fmt.Errorf(err.Error(), "failed to read a key")
		default:
			return nil
		}
	}
}

func (drawer *lessonDrawer) initDrawer(lesson lessons.Lesson, dataType int, page int) error {

	if drawer.termbox == nil {
		drawer.termbox = &termImpl{}
	}

	if err := drawer.termbox.init(); err != nil {
		return errors.Wrap(err, "failed to initialize termbox")
	}
	switch dataType {

	case lessons.Desc:
		drawer.drawTimer = time.AfterFunc(0, func() {
			drawer.drawDescription(lesson.Descriptions, dataType, page)
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

func (drawer *lessonDrawer) drawDescription(cards []lessons.DescriptionCard, dataType int, page int) {

	const pipeline = '│'
	const backArrow = '<'
	const nextArrow = '>'
	const arrowBody = '='

	width, height := drawer.termbox.size()
	drawer.termbox.clear(termbox.ColorDefault, termbox.ColorDefault)
	card := cards[page]
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

	if isBackArrowExist(dataType, page) {
		backValue := "  " + getBackValue(dataType, page) + " "
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
	if isNextArrowExist(len(cards), dataType, page) {
		nextValue := "  " + getNextValue(len(cards), dataType, page) + " "
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
			return errBack
			return fmt.Errorf("error no back arrow exist")
		case termbox.KeyArrowUp, termbox.KeyArrowRight, termbox.KeyCtrlF, termbox.KeyEnter:
			descDraw.rwMutex.RLock()
			defer descDraw.rwMutex.RUnlock()
			return errNext
		case termbox.KeyCtrlA:

		case termbox.KeyCtrlE:

		case termbox.KeyCtrlW:
		case termbox.KeyCtrlU:
		case termbox.KeyCtrlK, termbox.KeyCtrlP:
		case termbox.KeyCtrlJ, termbox.KeyCtrlN:
		case termbox.KeyTab:
		default:
			fmt.Println(e.Ch)
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

func goBack(lesson lessons.Lesson, dataType int, page int) {
	if dataType == lessons.Desc {
		descDraw.showLesson(lesson, dataType, page-1)
	} else if dataType == lessons.Interactive {
		if page == 0 {
			descDraw.showLesson(lesson, lessons.Desc, len(lesson.Descriptions)-1)
		} else {
			descDraw.showLesson(lesson, dataType, page-1)
		}
	} else if dataType == lessons.Quiz {
		if page == 0 {
			descDraw.showLesson(lesson, lessons.Interactive, len(lesson.InteractiveActions)-1)
		} else {
			descDraw.showLesson(lesson, dataType, page-1)
		}
	}
}

func goNext(lesson lessons.Lesson, dataType int, page int) {
	if dataType == lessons.Desc {
		if len(lesson.Descriptions) == page+1 {
			descDraw.showLesson(lesson, lessons.Interactive, 0)
		} else {
			descDraw.showLesson(lesson, dataType, page+1)
		}
	} else if dataType == lessons.Interactive {
		if len(lesson.InteractiveActions) == page+1 {
			descDraw.showLesson(lesson, lessons.Quiz, 0)
		} else {
			descDraw.showLesson(lesson, dataType, page+1)
		}
	} else if dataType == lessons.Quiz {
		if len(lesson.Quiz) == page+1 {
			descDraw.showEndPage()
		} else {
			descDraw.showLesson(lesson, dataType, page+1)
		}
	}
}

//TODO
func isNextArrowExist(length int, dataType int, page int) bool {
	if dataType == lessons.Desc {
		return !(length == page+1 && descDraw.lesson.InteractiveActions != nil && descDraw.lesson.Quiz != nil)
	} else if dataType == lessons.Interactive {
		return !(length == page+1 && descDraw.lesson.Quiz != nil)
	} else if dataType == lessons.Quiz {
		return !(length == page+1)
	}
	return false
}

func isBackArrowExist(dataType int, page int) bool {
	if dataType == lessons.Desc {
		return !(page == 0)
	} else if dataType == lessons.Interactive {
		return !(page == 0 && descDraw.lesson.Descriptions != nil)
	} else if dataType == lessons.Quiz {
		return !(page == 0 && descDraw.lesson.InteractiveActions != nil && descDraw.lesson.Descriptions != nil)
	}
	return false
}

func getBackValue(dataType int, page int) string {
	if dataType == lessons.Desc {
		return descDraw.lesson.Descriptions[page-1].Header
	} else if dataType == lessons.Interactive {
		if page == 0 && descDraw.lesson.Descriptions != nil {
			return descDraw.lesson.Descriptions[len(descDraw.lesson.Descriptions)-1].Header
		} else {
			return string(page-1) + ". Lab"
		}

	} else if dataType == lessons.Quiz {
		if page == 0 && descDraw.lesson.InteractiveActions != nil {
			return string(len(descDraw.lesson.InteractiveActions)) + ". Lab"
		} else if page == 0 && descDraw.lesson.Descriptions != nil {
			return descDraw.lesson.Descriptions[len(descDraw.lesson.Descriptions)-1].Header
		} else {
			return string(page-1) + ". Question"
		}
	} else {
		return "ERROR -- getBackValue :)"
	}
}
func getNextValue(length int, dataType int, page int) string {
	if dataType == lessons.Desc {
		if length == page+1 && descDraw.lesson.InteractiveActions != nil {
			return "1. Lab"
		} else if length == page+1 && descDraw.lesson.Quiz != nil {
			return "1. Question"
		} else {
			return descDraw.lesson.Descriptions[page+1].Header
		}

	} else if dataType == lessons.Interactive {
		if length == page+1 && descDraw.lesson.Quiz != nil {
			return "1. Question"
		} else {
			return string(page+1) + ". Lab"
		}
	} else if dataType == lessons.Quiz {
		return string(page+1) + ". Question"
	} else {
		return "getNextValue"
	}
}
