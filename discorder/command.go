package discorder

import (
	"encoding/json"
	"github.com/nsf/termbox-go"
	"strings"
)

type ArgumentDataType int

const (
	ArgumentDataTypeInt ArgumentDataType = iota
	ArgumentDataTypeFloat
	ArgumentDataTypeString
	ArgumentDataTypeBool
)

type Command struct {
	Name        string
	Description string
	Args        []*ArgumentDef
	Category    string
	Run         func(app *App, args []*Argument)
}

type ArgumentDef struct {
	Name     string
	Optional bool
	Datatype ArgumentDataType
}

type Argument struct {
	Name string      `json:"name"`
	Val  interface{} `json:"val"`
}

func (a *Argument) Int() (int, bool) {
	fVal, ok := a.Val.(float64)
	if !ok {
		return 0, false
	}

	return int(fVal), true
}

type KeyBind struct {
	Command string         `json:"command"`
	Args    []*Argument    `json:"args"`
	KeyComb KeyCombination `json:"key"`
	Alt     bool           `json:"alt"`
}

func (k KeyBind) Check(seq []termbox.Event) (partialMatch, fullMatch bool) {
	if len(seq) > len(k.KeyComb.Keys) {
		return
	}

	for i, event := range seq {
		keybindKey := k.KeyComb.Keys[i]
		if (event.Mod&termbox.ModAlt != 0 && !keybindKey.Alt) || (event.Mod&termbox.ModAlt == 0 && keybindKey.Alt) {
			return
		}
		if keybindKey.Char != "" {
			if string(event.Ch) != keybindKey.Char {
				return
			}
		} else {
			if event.Key != keybindKey.Special {
				return
			}
		}
	}

	if len(seq) < len(k.KeyComb.Keys) {
		partialMatch = true
	} else {
		fullMatch = true
	}
	return
}

type KeyCombination struct {
	Keys []*KeybindKey
}

// Alt+CtrlX-A
func (k *KeyCombination) UnmarshalJSON(data []byte) error {
	raw := ""
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	k.Keys = make([]*KeybindKey, 0)

	seqSplit := strings.Split(raw, "-")
	for _, sequence := range seqSplit {
		modSplit := strings.Split(sequence, "+")
		key := &KeybindKey{}
		for _, mod := range modSplit {
			if mod == "Alt" || mod == "alt" {
				key.Alt = true
				continue
			}
			special, ok := SpecialKeys[mod]
			if ok {
				key.Special = special
				continue
			}
			key.Char = mod
		}
	}
	return nil
}

type KeybindKey struct {
	Alt     bool
	Special termbox.Key
	Char    string
}

type CommandHandler interface {
	OnCommand(cmd *Command)
}

var SpecialKeys = map[string]termbox.Key{
	"F1":         termbox.KeyF1,
	"F2":         termbox.KeyF2,
	"F3":         termbox.KeyF3,
	"F4":         termbox.KeyF4,
	"F5":         termbox.KeyF5,
	"F6":         termbox.KeyF6,
	"F7":         termbox.KeyF7,
	"F8":         termbox.KeyF8,
	"F9":         termbox.KeyF9,
	"F10":        termbox.KeyF10,
	"F11":        termbox.KeyF11,
	"F12":        termbox.KeyF12,
	"Insert":     termbox.KeyInsert,
	"Delete":     termbox.KeyDelete,
	"Home":       termbox.KeyHome,
	"End":        termbox.KeyEnd,
	"Pgup":       termbox.KeyPgup,
	"Pgdn":       termbox.KeyPgdn,
	"ArrowUp":    termbox.KeyArrowUp,
	"ArrowDown":  termbox.KeyArrowDown,
	"ArrowLeft":  termbox.KeyArrowLeft,
	"ArrowRight": termbox.KeyArrowRight,

	"MouseLeft":      termbox.MouseLeft,
	"MouseMiddle":    termbox.MouseMiddle,
	"MouseRight":     termbox.MouseRight,
	"MouseRelease":   termbox.MouseRelease,
	"MouseWheelUp":   termbox.MouseWheelUp,
	"MouseWheelDown": termbox.MouseWheelDown,

	"CtrlTilde":      termbox.KeyCtrlTilde,
	"CtrlSpace":      termbox.KeyCtrlSpace,
	"CtrlA":          termbox.KeyCtrlA,
	"CtrlB":          termbox.KeyCtrlB,
	"CtrlC":          termbox.KeyCtrlC,
	"CtrlD":          termbox.KeyCtrlD,
	"CtrlE":          termbox.KeyCtrlE,
	"CtrlF":          termbox.KeyCtrlF,
	"CtrlG":          termbox.KeyCtrlG,
	"Backspace":      termbox.KeyBackspace,
	"CtrlH":          termbox.KeyCtrlH,
	"Tab":            termbox.KeyTab,
	"CtrlI":          termbox.KeyCtrlI,
	"CtrlJ":          termbox.KeyCtrlJ,
	"CtrlK":          termbox.KeyCtrlK,
	"CtrlL":          termbox.KeyCtrlL,
	"Enter":          termbox.KeyEnter,
	"CtrlM":          termbox.KeyCtrlM,
	"CtrlN":          termbox.KeyCtrlN,
	"CtrlO":          termbox.KeyCtrlO,
	"CtrlP":          termbox.KeyCtrlP,
	"CtrlQ":          termbox.KeyCtrlQ,
	"CtrlR":          termbox.KeyCtrlR,
	"CtrlS":          termbox.KeyCtrlS,
	"CtrlT":          termbox.KeyCtrlT,
	"CtrlU":          termbox.KeyCtrlU,
	"CtrlV":          termbox.KeyCtrlV,
	"CtrlW":          termbox.KeyCtrlW,
	"CtrlX":          termbox.KeyCtrlX,
	"CtrlY":          termbox.KeyCtrlY,
	"CtrlZ":          termbox.KeyCtrlZ,
	"Esc":            termbox.KeyEsc,
	"CtrlLsqBracket": termbox.KeyCtrlLsqBracket,
	"CtrlBackslash":  termbox.KeyCtrlBackslash,
	"CtrlRsqBracket": termbox.KeyCtrlRsqBracket,
	"CtrlSlash":      termbox.KeyCtrlSlash,
	"CtrlUnderscore": termbox.KeyCtrlUnderscore,
	"Space":          termbox.KeySpace,
	"Backspace2":     termbox.KeyBackspace2,
	"Ctrl2":          termbox.KeyCtrl2,
	"Ctrl3":          termbox.KeyCtrl3,
	"Ctrl4":          termbox.KeyCtrl4,
	"Ctrl5":          termbox.KeyCtrl5,
	"Ctrl6":          termbox.KeyCtrl6,
	"Ctrl7":          termbox.KeyCtrl7,
	"Ctrl8":          termbox.KeyCtrl8,
}
