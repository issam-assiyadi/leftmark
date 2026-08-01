package app

import (
	"log"

	"github.com/issam-assiyadi/leftmark/internal/content"
	"github.com/issam-assiyadi/leftmark/internal/ui/components/scrollview"
	"github.com/jroimartin/gocui"
)

type App struct {
	Pages       content.Pages
	CurrentPage int

	SidebarView string
	Content     *scrollview.View
}

func New(pages content.Pages) *App {
	return &App{
		Pages:       pages,
		CurrentPage: 0,

		SidebarView: "sidebar",
		Content: scrollview.New(scrollview.Config{
			BaseName:       "content",
			Title:          "Content",
			ScrollbarWidth: 3,
		}),
	}
}

func (app *App) Run() error {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.SetManagerFunc(func(g *gocui.Gui) error {
		if err := app.layout(g); err != nil {
			return err
		}
		return app.Render(g)
	})

	if err := app.BindKeys(g); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(app.SidebarView); err != nil {
		log.Println("unable to set initial view:", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
