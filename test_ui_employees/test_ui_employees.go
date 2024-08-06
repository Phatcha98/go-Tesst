package main

import (
    "github.com/lxn/walk"
    "github.com/lxn/walk/declarative"
)

func main() {
    mw := new(MyMainWindow)
    mw.Run()
}

type MyMainWindow struct {
    *walk.MainWindow
    nameLineEdit    *walk.LineEdit
    ageSpinBox      *walk.NumberEdit
    countryComboBox *walk.ComboBox
    positionLineEdit *walk.LineEdit
    wageSpinBox     *walk.NumberEdit
}

func (mw *MyMainWindow) Run() {
    if _, err := (declarative.MainWindow{
        AssignTo: &mw.MainWindow,
        Title:    "Employee Information",
        Layout:   declarative.VBox{},
        Children: []declarative.Widget{
            declarative.Label{Text: "Name:"},
            declarative.LineEdit{AssignTo: &mw.nameLineEdit},
            declarative.Label{Text: "Age:"},
            declarative.NumberEdit{AssignTo: &mw.ageSpinBox},
            declarative.Label{Text: "Country:"},
            declarative.ComboBox{
                AssignTo: &mw.countryComboBox,
                Model:    []string{"USA", "UK", "Canada", "Australia"},
            },
            declarative.Label{Text: "Position:"},
            declarative.LineEdit{AssignTo: &mw.positionLineEdit},
            declarative.Label{Text: "Wage:"},
            declarative.NumberEdit{AssignTo: &mw.wageSpinBox},
            declarative.PushButton{
                Text: "Submit",
                OnClicked: func() {
                    name := mw.nameLineEdit.Text()
                    age := mw.ageSpinBox.Value()
                    position := mw.positionLineEdit.Text()
                    wage := mw.wageSpinBox.Value()

                    println("Name:", name)
                    println("Age:", age)
                    println("Position:", position)
                    println("Wage:", wage)
                },
            },
        },
    }).Run(); err != nil {
        panic(err)
    }
}
