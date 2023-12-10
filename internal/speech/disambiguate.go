package speech

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Word []Letter

func (w Word) String() string {
	ret := ""
	for _, l := range w {
		ret += l.String()
	}
	return ret
}

type Letter struct {
	Letter rune
	Vowel  rune
	Shadda bool
}

func (l Letter) String() string {
	var shadda string
	if l.Shadda {
		shadda = Shadda
	} else {
		shadda = ""
	}
	return fmt.Sprintf("%v%v%v", string(l.Letter), string(l.Vowel), shadda)
}

func Disambiguate(text string) ([]Word, error) {
	data := fmt.Sprintf(`{"dialect": "msa", "sentence": "%v"}`,
		text)
	body := strings.NewReader(data)
	resp, err := http.Post("https://camelira.abudhabi.nyu.edu/api/disambig",
		"application/json", body)
	if err != nil {
		return nil, errors.Join(ErrRequest, err)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(ErrBadResponse, err)
	}

	inter := camelResponse{}
	err = json.Unmarshal(res, &inter)
	if err != nil {
		return nil, errors.Join(BadFormatError{Text: string(res),
			ExpectedFormat: "json"}, err)
	}

	words := make([]Word, len(inter.Output.Disambig))
	for i, cWord := range inter.Output.Disambig {
		words[i] = []Letter{}
		lastLetter := -1
		for _, cc := range cWord.Analyses[0].Analysis.Diac {
			if IsArabicLetter(cc) {
				words[i] = append(words[i], Letter{Letter: cc})
				lastLetter += 1
			} else if IsVowel(cc) {
				words[i][lastLetter].Vowel = cc
			} else if IsShadda(cc) {
				words[i][lastLetter].Shadda = true
			} else if string(cc) != SuperscriptAlef {
				return nil, UnrecognizedCharacterError{Character: cc}
			}
		}
	}
	return words, nil
}

type camelResponse struct {
	Output camelOutput `json:"output"`
}

type camelOutput struct {
	Disambig []camelWord `json:"disambig"`
}

type camelWord struct {
	Analyses []camelAnalysisMeta `json:"analyses"`
}

type camelAnalysisMeta struct {
	Analysis camelAnalysis `json:"analysis"`
}

type camelAnalysis struct {
	Diac string `json:"diac"`
}