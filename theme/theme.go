package theme

import lg "github.com/charmbracelet/lipgloss"

// color reference: https://codehs.com/uploads/7c2481e9158534231fcb3c9b6003d6b3

var Black = lg.Color("0")
var White = lg.Color("15")
var Grey = lg.Color("7")

var Primary = lg.Color("147")
var Secondary = lg.Color("32")

var Heading = lg.NewStyle().Bold(true).Foreground(White).PaddingLeft(1).PaddingRight(1).Background(Primary)

var fgFaded = lg.Color("12")
var Faded = lg.NewStyle().Foreground(Grey)
