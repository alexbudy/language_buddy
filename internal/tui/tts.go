package tui

import (
	"os"
	"os/exec"
	"strings"

	"github.com/alexbudy/go_spanish_rewrite/internal/store"
)

var spanishVoices = []string{
	"Microsoft Sabina Desktop",
	"Microsoft Helena Desktop",
	"Microsoft Laura Desktop",
}

var englishVoices = []string{
	"Microsoft Zira Desktop",
	"Microsoft David Desktop",
}

// these voices do not come through standard Windows install - see README
var ukrainianVoices = []string{
	"Natalia",
	"Anatol",
}

var langCodeToVoices = map[string][]string{
	"en": englishVoices,
	"es": spanishVoices,
	"uk": ukrainianVoices,
}

func pronounceWord(word store.Word, lang string) {
	word1 := ""

	switch lang { // TODO potentially generalize
	case "en":
		word1 = word.English
	case "es":
		word1 = word.Spanish
	case "uk":
		word1 = word.Ukrainian
	}

	go func() {
		pronounceText(word1, langCodeToVoices[lang])
	}()
}

// Given word, pronounce it in lang1 then lang 2
func pronounceWordPair(word store.Word, lang1 string, lang2 string) {
	word1 := ""
	word2 := ""

	switch lang1 { // TODO potentially generalize
	case "en":
		word1 = word.English
	case "es":
		word1 = word.Spanish
	case "uk":
		word1 = word.Ukrainian
	}
	switch lang2 {
	case "en":
		word2 = word.English
	case "es":
		word2 = word.Spanish
	case "uk":
		word2 = word.Ukrainian
	}	

	pronounceWordPairStringsOnly(word1, word2, lang1, lang2)
}

// this function operates on strings, not the word struct
func pronounceWordPairStringsOnly(word1 string, word2 string, lang1 string, lang2 string) {
	go func() {
		pronounceText(word1, langCodeToVoices[lang1])
		go func() {
			pronounceText(word2, langCodeToVoices[lang2])
		}()
	}()
}

// pronounce the text in the preferred voices, in order of preference. If none of the preferred voices are available, it will use the default voice.
func pronounceText(text string, preferredVoices []string) {
	script := `
		Add-Type -AssemblyName System.Speech

		$s = New-Object System.Speech.Synthesis.SpeechSynthesizer

		foreach ($name in $env:TTS_VOICES -split '\|') {
			$voice = $s.GetInstalledVoices() |
				Where-Object { $_.VoiceInfo.Name -eq $name } |
				Select-Object -First 1

			if ($voice) {
				$s.SelectVoice($name)
				break
			}
		}

		$s.Speak($env:TTS_TEXT)
		`

	cmd := exec.Command(
		"powershell.exe",
		"-NoProfile",
		"-Command",
		script,
	)

	cmd.Env = append(
		os.Environ(),
		"TTS_TEXT="+text,
		"TTS_VOICES="+strings.Join(preferredVoices, "|"),
	)

	cmd.Run()
}
