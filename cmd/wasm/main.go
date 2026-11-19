//go:build js && wasm

// Command wasm exposes the Enigma simulator to JavaScript so it can run in a web browser. Build it with:
//
//	GOOS=js GOARCH=wasm go build -o ultra.wasm ./cmd/wasm
//
// It defines a single global function, ultraRun(configJSON, text), which returns a JSON string. See the web folder for
// the page that uses it.
package main

import (
	"encoding/json"
	"strings"
	"syscall/js"

	"github.com/atbu/ultra/conversion"
	"github.com/atbu/ultra/machine"
)

type step struct {
	Component machine.Component `json:"component"`
	In        string            `json:"in"`
	Out       string            `json:"out"`
}

type result struct {
	Error string `json:"error,omitempty"`
	// The letters that were encrypted, i.e. the text with everything other than A to Z removed.
	Input  string `json:"input"`
	Output string `json:"output"`
	// The letters showing in the left, middle and right rotor windows after the text has been typed.
	Positions      string `json:"positions"`
	FourthPosition string `json:"fourthPosition,omitempty"`
	// The path the signal took for the last letter typed.
	Trace []step `json:"trace,omitempty"`
}

// The run function sets up a fresh machine from the configuration, types the letters of text into it and reports the
// result. Starting fresh each time means the page can simply call it again whenever the settings or the text change.
func run(_ js.Value, args []js.Value) any {
	var res result

	if len(args) != 2 {
		res.Error = "ultraRun expects two arguments: a configuration and some text"
		return encode(res)
	}

	var config machine.Config
	if err := json.Unmarshal([]byte(args[0].String()), &config); err != nil {
		res.Error = "invalid configuration: " + err.Error()
		return encode(res)
	}

	m, err := machine.NewFromConfig(config)
	if err != nil {
		res.Error = err.Error()
		return encode(res)
	}

	var input, output strings.Builder
	letters := lettersOnly(args[1].String())
	for i, key := range letters {
		if i < len(letters)-1 {
			output.WriteRune(m.PressKey(key))
			continue
		}

		out, trace := m.PressKeyWithTrace(key)
		output.WriteRune(out)
		for _, s := range trace {
			res.Trace = append(res.Trace, step{s.Component, string(s.In), string(s.Out)})
		}
	}
	input.WriteString(letters)

	res.Input = input.String()
	res.Output = output.String()
	res.Positions = string([]rune{
		conversion.ConvertIndexToChar(m.LeftRotor.Position),
		conversion.ConvertIndexToChar(m.MiddleRotor.Position),
		conversion.ConvertIndexToChar(m.RightRotor.Position),
	})
	if m.FourthRotor != nil {
		res.FourthPosition = string(conversion.ConvertIndexToChar(m.FourthRotor.Position))
	}

	return encode(res)
}

// The lettersOnly function upper-cases text and drops anything that isn't A to Z, since an Enigma keyboard has no
// other keys.
func lettersOnly(text string) string {
	var b strings.Builder
	for _, char := range strings.ToUpper(text) {
		if char >= 'A' && char <= 'Z' {
			b.WriteRune(char)
		}
	}
	return b.String()
}

func encode(res result) string {
	encoded, _ := json.Marshal(res)
	return string(encoded)
}

func main() {
	js.Global().Set("ultraRun", js.FuncOf(run))

	// Let the page know the simulator is ready to use.
	if ready := js.Global().Get("onUltraReady"); ready.Type() == js.TypeFunction {
		ready.Invoke()
	}

	// Keep the program running so JavaScript can keep calling ultraRun.
	select {}
}
