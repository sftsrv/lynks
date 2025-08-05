package theme

import lg "github.com/charmbracelet/lipgloss"

// color reference: https://codehs.com/uploads/7c2481e9158534231fcb3c9b6003d6b3

var black = lg.Color("0")
var white = lg.Color("15")
var highlight = lg.Color("147")
var grey = lg.Color("7")

var Heading = lg.NewStyle().Bold(true).Foreground(white).Background(highlight)

var fgFaded = lg.Color("12")
var Faded = lg.NewStyle().Foreground(grey)
