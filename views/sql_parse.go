package views

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func SqlParsePre(myApp fyne.App) {
	w3 := myApp.NewWindow("sql字串输入")
	we := widget.NewMultiLineEntry()
	we.SetMinRowsVisible(15) // 设置高度为 400
	sqlBox := container.NewVBox(we)
	//sqlBox.Resize(fyne.NewSize(300, 400)) // 设置高度为 400
	button := widget.NewButton("确认", func() {
		SqlStatement = we.Text
		fmt.Println("监听到sql", we.Text)
		w3.Close()
	})

	wec := container.New(layout.NewVBoxLayout(), widget.NewLabel("请输入sql 建表语句"), sqlBox, button)
	w3.SetContent(wec)
	w3.Resize(fyne.NewSize(300, 400))
	w3.Show()
}
