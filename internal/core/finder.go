package core

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/martenvd/gofind/internal/utils"
	"github.com/rivo/tview"
)

func Find(dirs []string) {

	app := tview.NewApplication()

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	inputField := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldTextColor(tcell.ColorGreen)
	inputField.SetBorder(true)

	resultsList := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)
	resultsList.SetBorder(true)

	vimInfo := tview.NewTextArea()
	vimInfo.SetBorder(true)
	vimInfo.SetText("--INSERT--", true)

	filteredResults := utils.GetFilteredResults(currentDir, "", dirs)
	for _, result := range filteredResults {
		resultsList.AddItem(result, "", 0, nil)
	}

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 3, 1, true).
		AddItem(resultsList, 0, 10, false).
		AddItem(vimInfo, 3, 1, false)

	inputField.SetChangedFunc(func(text string) {
		currentCursor := resultsList.GetCurrentItem()

		resultsList.Clear()
		filteredResults := utils.GetFilteredResults(currentDir, text, dirs)
		for _, result := range filteredResults {
			resultsList.AddItem(result, "", 0, nil)
		}

		if currentCursor > resultsList.GetItemCount()-1 {
			currentCursor = resultsList.GetItemCount() - 1
		}

		resultsList.SetCurrentItem(currentCursor)
	})

	vimKeys := false
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp, tcell.KeyDown:
			app.SetFocus(resultsList)
			resultsList.InputHandler()(event, nil)
			return nil
		case tcell.KeyEnter:
			currentPath, _ := resultsList.GetItemText(resultsList.GetCurrentItem())
			app.Stop()
			utils.OpenInVSCodeFromFinder(currentPath, resultsList.GetItemCount())
			return nil
		case tcell.KeyEscape:
			app.SetFocus(resultsList)
			vimKeys = true
			vimInfo.SetText("--NORMAL--", true)
			return nil
		}
		return event
	})

	colonPressed := false
	input := ""
	resultsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Rune() {
		case 'k':
			if colonPressed {
				input += "k"
				vimInfo.SetText(input, true)
			} else {
				currentIndex := resultsList.GetCurrentItem()
				if currentIndex > 0 {
					resultsList.SetCurrentItem(currentIndex - 1)
				}
			}
		case 'j':
			if colonPressed {
				input += "j"
				vimInfo.SetText(input, true)
			} else {
				currentIndex := resultsList.GetCurrentItem()
				if currentIndex < resultsList.GetItemCount()-1 {
					resultsList.SetCurrentItem(currentIndex + 1)
				}
			}
		case ':':
			input = ":"
			vimInfo.SetText(input, true)
			colonPressed = true
		case 'i':
			if !colonPressed {
				vimKeys = false
				vimInfo.SetText("--INSERT--", true)
				app.SetFocus(inputField)
				return nil
			} else {
				input += "i"
				vimInfo.SetText(input, true)
			}
		default:
			if colonPressed && event.Key() != tcell.KeyBackspace2 && event.Key() != tcell.KeyEnter {
				input += string(event.Rune())
				vimInfo.SetText(input, true)
			}
		}
		switch event.Key() {
		case tcell.KeyUp:
			currentIndex := resultsList.GetCurrentItem()
			if currentIndex > 0 {
				resultsList.SetCurrentItem(currentIndex - 1)
			}
		case tcell.KeyDown:
			currentIndex := resultsList.GetCurrentItem()
			if currentIndex < resultsList.GetItemCount()-1 {
				resultsList.SetCurrentItem(currentIndex + 1)
			}
		case tcell.KeyEnter:
			if colonPressed && input == ":q" {
				app.Stop()
			} else if colonPressed {
				input = ""
				vimInfo.SetText("--NORMAL--", true)
				colonPressed = false
			} else if !colonPressed {
				currentPath, _ := resultsList.GetItemText(resultsList.GetCurrentItem())
				app.Stop()
				utils.OpenInVSCodeFromFinder(currentPath, resultsList.GetItemCount())
			}
		case tcell.KeyBackspace2:

			if input == "" {
				colonPressed = false
				vimInfo.SetText("--NORMAL--", true)
			}
			if colonPressed {
				input = input[:len(input)-1]
				vimInfo.SetText(input, true)
			}
			return nil
		case tcell.KeyEscape:
			colonPressed = false
			input = ""
			vimInfo.SetText("--NORMAL--", true)
			return nil
		default:
			if !vimKeys {
				inputField.InputHandler()(event, nil)
				return nil
			}
		}
		return nil
	})

	if err := app.SetRoot(flex, true).SetFocus(inputField).Run(); err != nil {
		panic(err)
	}
}
