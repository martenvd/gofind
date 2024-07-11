package core

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Find(dirs []string) {
	app := tview.NewApplication()

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

	filteredResults := getFilteredResults("", dirs)
	for _, result := range filteredResults {
		resultsList.AddItem(result, "", 0, nil)
	}

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 3, 1, true).
		AddItem(resultsList, 0, 10, false).
		AddItem(vimInfo, 0, 1, false)

	inputField.SetChangedFunc(func(text string) {
		resultsList.Clear()
		filteredResults := getFilteredResults(text, dirs)
		for _, result := range filteredResults {
			resultsList.AddItem(result, "", 0, nil)
		}
	})

	vimKeys := false
	colonPressed := false
	quitPressed := false
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp, tcell.KeyDown:
			app.SetFocus(resultsList)
			resultsList.InputHandler()(event, nil)
			return nil
		case tcell.KeyEnter:
			openInVSCodeFromFinder(resultsList.GetItemCount(), resultsList, app)
			return nil
		case tcell.KeyEscape:
			app.SetFocus(resultsList)
			vimKeys = true
			vimInfo.SetText("--NORMAL--", true)
			return nil
		}
		return event
	})

	resultsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'k':
			currentIndex := resultsList.GetCurrentItem()
			if currentIndex > 0 {
				resultsList.SetCurrentItem(currentIndex - 1)
			}
		case 'j':
			currentIndex := resultsList.GetCurrentItem()
			if currentIndex < resultsList.GetItemCount()-1 {
				resultsList.SetCurrentItem(currentIndex + 1)
			}
		case ':':
			vimInfo.SetText(":", true)
			colonPressed = true
		case 'q':
			if colonPressed {
				vimInfo.SetText(":q", true)
				quitPressed = true
			}
		case 'i':
			vimKeys = false
			vimInfo.SetText("--INSERT--", true)
			app.SetFocus(inputField)
			return nil
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
			if colonPressed && quitPressed {
				app.Stop()
			} else {
				openInVSCodeFromFinder(resultsList.GetItemCount(), resultsList, app)
			}
		case tcell.KeyBackspace2:
			if quitPressed {
				quitPressed = false
				vimInfo.SetText(":", true)
			} else if colonPressed {
				colonPressed = false
				vimInfo.SetText("--NORMAL--", true)
			}
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

func getFilteredResults(input string, dirs []string) []string {
	filteredResults := []string{}
	for _, dir := range dirs {
		if strings.Contains(strings.ToLower(dir), strings.ToLower(input)) {
			filteredResults = append(filteredResults, dir)
		}
	}
	return filteredResults
}

func openInVSCodeFromFinder(resultlistCount int, list *tview.List, app *tview.Application) {
	if resultlistCount != 0 {
		currentPath, _ := list.GetItemText(list.GetCurrentItem())
		app.Stop()
		currentItemName := strings.Split(currentPath, "/")[len(strings.Split(currentPath, "/"))-1]
		fmt.Println("To open the directory type:")
		fmt.Println()
		fmt.Print("cd ", currentPath, "\n")
		fmt.Println()
		fmt.Println("Opening:", currentItemName)
		cmd := exec.Command("code", currentPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		panic("No results found")
	}
}
