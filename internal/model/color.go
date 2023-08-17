package model

import "github.com/jomei/notionapi"

var CalendarColorMap = map[string]string{
	"1":  "#ac725e",
	"2":  "#d06b64",
	"3":  "#f83a22",
	"4":  "#fa573c",
	"5":  "#ff7537",
	"6":  "#ffad46",
	"7":  "#42d692",
	"8":  "#16a765",
	"9":  "#7bd148",
	"10": "#b3dc6c",
	"11": "#fbe983",
	"12": "#fad165",
	"13": "#92e1c0",
	"14": "#9fe1e7",
	"15": "#9fc6e7",
	"16": "#4986e7",
	"17": "#9a9cff",
	"18": "#b99aff",
	"19": "#c2c2c2",
	"20": "#cabdbf",
	"21": "#cca6ac",
	"22": "#f691b2",
	"23": "#cd74e6",
	"24": "#a47ae2",
}

var EventColorMap = map[string]string{
	"1":  "#a4bdfc",
	"2":  "#7ae7bf",
	"3":  "#dbadff",
	"4":  "#ff887c",
	"5":  "#fbd75b",
	"6":  "#ffb878",
	"7":  "#46d6db",
	"8":  "#e1e1e1",
	"9":  "#5484ed",
	"10": "#51b749",
	"11": "#dc2127",
}

// default, gray, brown, orange, yellow, green, blue, purple, pink, red
var GCalendaToNotionColorMap = map[string]notionapi.Color{
	"1":  notionapi.ColorBrown,
	"2":  notionapi.ColorBrown,
	"3":  notionapi.ColorRed,
	"4":  notionapi.ColorOrange,
	"5":  notionapi.ColorOrange,
	"6":  notionapi.ColorYellow,
	"7":  notionapi.ColorGreen,
	"8":  notionapi.ColorGreen,
	"9":  notionapi.ColorGreen,
	"10": notionapi.ColorGreen,
	"11": notionapi.ColorYellow,
	"12": notionapi.ColorYellow,
	"13": notionapi.ColorGreen,
	"14": notionapi.ColorBlue,
	"15": notionapi.ColorBlue,
	"16": notionapi.ColorBlue,
	"17": notionapi.ColorPurple,
	"18": notionapi.ColorPurple,
	"19": notionapi.ColorDefault,
	"20": notionapi.ColorGray,
	"21": notionapi.ColorBrown,
	"22": notionapi.ColorRed,
	"23": notionapi.ColorPurple,
	"24": notionapi.ColorPurple,
}

// eventは例外のカレンダーの色が来ることがある
var GEventToNotionColorMap = map[string]notionapi.Color{
	"1":  notionapi.ColorBlue,
	"2":  notionapi.ColorGreen,
	"3":  notionapi.ColorPurple,
	"4":  notionapi.ColorOrange,
	"5":  notionapi.ColorYellow,
	"6":  notionapi.ColorOrange,
	"7":  notionapi.ColorBlue,
	"8":  notionapi.ColorDefault,
	"9":  notionapi.ColorBlue,
	"10": notionapi.ColorGreen,
	"11": notionapi.ColorRed,
}
