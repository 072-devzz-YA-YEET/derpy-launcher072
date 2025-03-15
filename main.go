package main

import (
	"derpy-launcher072/igdb"
	"derpy-launcher072/library"
	"derpy-launcher072/torrent"
	"derpy-launcher072/utils/settings"
	"embed"
	"fmt"
	"log"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.
var libraryManager *library.Library
var apiManager *igdb.APIManager
var torrentManager *torrent.Manager

var app *application.App

//go:embed all:frontend/dist
var assets embed.FS

type WindowService struct{}

func (w *WindowService) Minimize() {
	app.CurrentWindow().Minimise()
}

func (w *WindowService) Maximize() {
	app.CurrentWindow().Maximise()
}

func (w *WindowService) Close() {
	app.CurrentWindow().Close()
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	// 🐐routine
	settings, err := settings.LoadSettings(filepath.Join("settings.json"))
	if err != nil {
		fmt.Println(err)
		return
	}

	libraryManager = library.GetLibrary()
	apiManager = igdb.NewAPI()
	torrentManager = torrent.StartClient(settings.DownloadPath)

	fmt.Println(settings)

	//go func() {
	//	results := torrent.Scrape_1337x("goat simulator 3")
	//	for _, result := range results {
	//		data := torrent.Get_1337x_data(result)
	//		fmt.Printf("Title: %s\nUploader: %s\nDownloads: %d\nDate: %s\n\n", data.Title, data.Uploader, data.Downloads, data.Date)
	//	}
	//}()

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app = application.New(application.Options{
		Name: "derpyLauncher",
		Services: []application.Service{
			application.NewService(torrentManager),
			application.NewService(apiManager),
			application.NewService(libraryManager),
			application.NewService(&WindowService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:     "derpyLauncher",
		Width:     1200,
		Height:    900,
		MinHeight: 700,
		MinWidth:  1064,
		Frameless: true,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	// Run the application. This blocks until the application has been exited.
	err = app.Run()
	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
