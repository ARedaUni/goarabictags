package main

import (
	"net/http"

	"github.com/amrojjeh/arabic-tags/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// TODO(Amr Ojjeh): Setup secure headers
func (app *application) routes() http.Handler {
	app.u = url{}
	router := httprouter.New()
	router.NotFound = app.notFound()

	router.Handler(http.MethodGet, "/static/*file",
		http.FileServer(http.FS(ui.Files)))

	router.Handler(http.MethodGet, app.u.register(), app.registerGet())
	router.Handler(http.MethodPost, app.u.register(), app.registerPost())
	router.Handler(http.MethodGet, app.u.login(), app.loginGet())
	router.Handler(http.MethodPost, app.u.login(), app.loginPost())

	authRequired := alice.New(app.authRequired)
	// TODO(Amr Ojjeh): Write an index page
	router.Handler(http.MethodGet, app.u.index(), authRequired.Then(app.homeGet()))
	router.Handler(http.MethodPost, app.u.logout(), authRequired.Then(app.logoutPost()))
	router.Handler(http.MethodGet, app.u.home(), authRequired.Then(app.homeGet()))
	router.Handler(http.MethodGet, app.u.createExcerpt(), authRequired.Then(app.createExcerptGet()))
	router.Handler(http.MethodPost, app.u.createExcerpt(), authRequired.Then(app.createExcerptPost()))

	excerptRequired := alice.New(app.excerptRequired)
	router.Handler(http.MethodGet, app.u.excerpt(":id"), excerptRequired.Then(app.excerptGet()))
	router.Handler(http.MethodPost, app.u.excerpt(":id"), excerptRequired.Then(app.excerptPost()))
	router.Handler(http.MethodGet, app.u.excerptTitle(":id"), excerptRequired.Then(app.excerptTitleGet()))
	router.Handler(http.MethodPost, app.u.excerptTitle(":id"), excerptRequired.Then(app.excerptTitlePost()))

	// TODO(Amr Ojjeh): ownerRequired doesn't check if auth is owner
	ownerRequired := excerptRequired.Extend(authRequired)
	router.Handler(http.MethodPost, app.u.excerptLock(":id"), ownerRequired.Then(app.excerptNextPost()))
	router.Handler(http.MethodPost, app.u.excerptEditLetter(":id", ":wid", ":lid"), ownerRequired.Then(app.excerptEditLetterPost()))
	router.Handler(http.MethodGet, app.u.excerptEditWord(":id"), ownerRequired.Then(app.excerptEditWordGet()))
	router.Handler(http.MethodPost, app.u.excerptEditWord(":id"), ownerRequired.Then(app.excerptEditWordPost()))
	router.Handler(http.MethodPost, app.u.wordRight(":id", ":wid"), ownerRequired.Then(app.wordRightPost()))
	router.Handler(http.MethodPost, app.u.wordLeft(":id", ":wid"), ownerRequired.Then(app.wordLeftPost()))
	router.Handler(http.MethodPost, app.u.wordAdd(":id", ":wid"), ownerRequired.Then(app.wordAddPost()))
	router.Handler(http.MethodPost, app.u.wordRemove(":id", ":wid"), ownerRequired.Then(app.wordRemovePost()))
	router.Handler(http.MethodPost, app.u.wordConnect(":id", ":wid"), ownerRequired.Then(app.wordConnectPost()))

	base := alice.New(
		app.recoverPanic,
		app.logRequest,
		app.session.LoadAndSave,
		app.getUser,
	)
	return base.Then(router)
}
