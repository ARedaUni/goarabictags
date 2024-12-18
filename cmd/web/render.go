package main

import (
	"strconv"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/amrojjeh/arabic-tags/ui/pages"
	"github.com/amrojjeh/arabic-tags/ui/partials"
	"github.com/amrojjeh/kalam"
	g "github.com/maragudk/gomponents"
)

func renderEdit(u url,
	e models.Excerpt, user models.User, ws []models.Word,
	selectedId int,
	error, warning string) g.Node {
	props := pages.EditProps{
		ExcerptTitle: e.Title,
		Inspector:    nil,
		Text:         renderText(u, e, ws, selectedId),
		ReadOnly:     user.Email != e.AuthorEmail,
		Nav:          renderNav(u, user),
		Error:        error,
		Warning:      warning,
		TitleUrl:     u.excerptTitle(e.Id),
		ExportUrl:    u.excerptExport(e.Id),
	}

	for _, w := range ws {
		if w.Id == selectedId {
			props.Inspector = renderInspector(u, user, e, w)
		}
	}

	return pages.EditPage(props)
}

func renderText(u url,
	e models.Excerpt, words []models.Word,
	selectedId int) g.Node {
	wps := []partials.WordProps{}
	for _, word := range words {
		wps = append(wps, partials.WordProps{
			Id:          strconv.Itoa(word.Id),
			Word:        word.Word,
			Punctuation: word.Punctuation,
			Connected:   word.Connected,
			Selected:    word.Id == selectedId,
			GetUrl:      u.wordSelect(e.Id, word.Id),
		})
	}

	return partials.Text(wps)
}

func renderNav(u url,
	user models.User) []g.Node {
	return partials.SimpleNav(user.Username, u.login(), u.register(), u.logout())
}

func renderInspector(u url,
	user models.User, e models.Excerpt, w models.Word) g.Node {
	if w.Punctuation {
		return renderInspectorPunctuation(u, user, e, w)
	}
	return renderInspectorWord(u, user, e, w)
}

func renderInspectorWord(u url,
	user models.User, e models.Excerpt, w models.Word) g.Node {
	props := partials.InspectorProps{
		Id:            strconv.Itoa(w.Id),
		Word:          w.Word,
		Letters:       []partials.LetterProps{},
		ReadOnly:      user.Email != e.AuthorEmail,
		Connected:     w.Connected,
		Ignore:        w.Ignore,
		SentenceStart: w.SentenceStart,
		CaseOptions: []partials.DropdownOption{
			{
				Value:    "",
				Selected: false,
			},
		},
		StateOptions:     []partials.DropdownOption{},
		EditWordUrl:      u.wordEditArgs(e.Id, w.Id),
		CaseUrl:          u.wordCase(e.Id, w.Id),
		StateUrl:         u.wordState(e.Id, w.Id),
		MoveRightUrl:     u.wordRight(e.Id, w.Id),
		MoveLeftUrl:      u.wordLeft(e.Id, w.Id),
		AddWordUrl:       u.wordAdd(e.Id, w.Id),
		RemoveWordUrl:    u.wordRemove(e.Id, w.Id),
		ConnectedUrl:     u.wordConnect(e.Id, w.Id),
		IgnoreUrl:        u.wordIgnore(e.Id, w.Id),
		SentenceStartUrl: u.wordSentenceStart(e.Id, w.Id),
	}

	for _, c := range kalam.Cases {
		props.CaseOptions = append(props.CaseOptions,
			partials.DropdownOption{
				Value:    c,
				Selected: w.Case == c,
			})
	}

	if w.Case != "" {
		for _, s := range kalam.States[w.Case] {
			props.StateOptions = append(props.StateOptions,
				partials.DropdownOption{
					Value:    s,
					Selected: w.State == s,
				})
		}
	}

	ls := kalam.LetterPacks(w.Word)
	for i, l := range ls {
		props.Letters = append(props.Letters,
			partials.LetterProps{
				Letter:          l.Unpointed(false),
				ShortVowel:      l.Vowel,
				Shadda:          l.Shadda,
				SuperscriptAlef: l.SuperscriptAlef,
				Index:           i,
				PostUrl:         u.letterEdit(e.Id, w.Id, i),
			})
	}

	return partials.Inspector(props)
}

func renderInspectorPunctuation(u url,
	user models.User, e models.Excerpt, w models.Word) g.Node {
	props := partials.InspectorProps{
		Id:               strconv.Itoa(w.Id),
		Word:             w.Word,
		ReadOnly:         user.Email != e.AuthorEmail,
		Connected:        w.Connected,
		SentenceStart:    w.SentenceStart,
		EditWordUrl:      u.wordEditArgs(e.Id, w.Id),
		MoveRightUrl:     u.wordRight(e.Id, w.Id),
		MoveLeftUrl:      u.wordLeft(e.Id, w.Id),
		AddWordUrl:       u.wordAdd(e.Id, w.Id),
		RemoveWordUrl:    u.wordRemove(e.Id, w.Id),
		ConnectedUrl:     u.wordConnect(e.Id, w.Id),
		SentenceStartUrl: u.wordSentenceStart(e.Id, w.Id),
	}

	return partials.Inspector(props)
}

func renderManuscript(u url,
	e models.Excerpt, ms models.Manuscript, user models.User,
	error, warning string) g.Node {
	props := pages.ManuscriptProps{
		ExcerptTitle:        e.Title,
		ReadOnly:            user.Email != e.AuthorEmail,
		AcceptedPunctuation: kalam.PunctuationRegex().String(),
		Content:             ms.Content,
		Warning:             warning,
		Error:               error,
		Username:            user.Username,
		SaveUrl:             u.excerpt(e.Id),
		NextUrl:             u.excerptLock(e.Id),
		TitleUrl:            u.excerptTitle(e.Id),
		LoginUrl:            u.login(),
		RegisterUrl:         u.register(),
		LogoutUrl:           u.logout(),
	}

	return pages.ManuscriptPage(props)
}

func renderEditLetter(u url,
	e models.Excerpt, w models.Word) g.Node {
	return partials.EditLetter(strconv.Itoa(w.Id), u.wordEditArgs(e.Id, w.Id),
		u.wordSelect(e.Id, w.Id), w.Word, w.Connected)
}
