package wallutils

import (
	"fmt"
	"os"
)

// Gnome2 windowmanager detector
type Gnome2 struct {
	mode    string
	verbose bool
}

// Name returns the name of this window manager or desktop environment
func (g2 *Gnome2) Name() string {
	return "Gnome2"
}

func (g2 *Gnome2) ExecutablesExists() bool {
	return which("gconftool-2") != ""
}

func (g2 *Gnome2) Running() bool {
	return (os.Getenv("GDMSESSION") == "gnome") || (os.Getenv("DESKTOP_SESSION") == "gnome")
}

func (g2 *Gnome2) SetMode(mode string) {
	g2.mode = mode
}

func (g2 *Gnome2) SetVerbose(verbose bool) {
	g2.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (g2 *Gnome2) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// initialize the mode setting (stretched/tiled etc)
	mode := defaultMode
	if g2.mode != "" {
		mode = g2.mode
	}

	switch mode {
	case "none", "wallpaper", "centered", "scaled", "stretched", "zoom", "spanned":
		break
	case "stretch":
		mode = "stretched"
	case "center":
		mode = "centered"
	case "fill", "scale":
		mode = "scaled"
	case "tile":
		mode = "wallpaper"
	default:
		// Invalid and unrecognized desktop wallpaper mode
		return fmt.Errorf("invalid desktop wallpaper mode for GNOME2: %s", mode)
	}

	// Set the wallpaper mode
	if err := run("gconftool-2", []string{"--type", "string", "--set", "/desktop/gnome/background/picture_options", mode}, g2.verbose); err != nil {
		return err
	}
	// Set the wallpaper image
	return run("gconftool-2", []string{"--type", "string", "--set", "/desktop/gnome/background/picture_filename", imageFilename}, g2.verbose)
}
