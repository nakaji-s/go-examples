package main

import (
	"github.com/marcusolsson/tui-go"
)

func main() {
	buffer := tui.NewTextEdit()
	buffer.SetSizePolicy(tui.Expanding, tui.Expanding)
	buffer.SetText(body)
	buffer.SetFocused(true)
	buffer.SetWordWrap(true)

	status := tui.NewStatusBar("lorem.txt")

	root := tui.NewVBox(buffer, status)

	ui := tui.New(root)
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

const body = `hello world!`
