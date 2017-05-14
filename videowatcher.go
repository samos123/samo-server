package samo

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"net/http"
	"time"
)

func SetUpVideoWatcher(dirs ...string) *watcher.Watcher {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Create, watcher.Rename)
	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	for _, dir := range dirs {
		// Watch test_folder recursively for changes.
		if err := w.AddRecursive(dir); err != nil {
			log.Fatalln(err)
		}
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	go func() {
		// Start the watching process - it'll check for changes every 100ms.
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}()

	return w
}

func DirAddHTTPHandler(w *watcher.Watcher) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		path := r.FormValue("path")
		if err := w.AddRecursive(path); err != nil {
			log.Fatalln(err)
		}

		//out, err := exec.Command("find", path, "-name \"*.mp4\" -or -name \"*.mkv\"").Output()
		fmt.Fprintln(rw, "Added path:", path)
	})
}
