package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dachinat/colornameconv"
	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/godbus/dbus/v5"
)

// ColorHistory represents a saved color
type ColorHistory struct {
	Color     string   `json:"color"`
	ColorName string   `json:"color_name"`
	RGB       [3]uint8 `json:"rgb"`
	Timestamp int64    `json:"timestamp"`
}

// Global history storage
var colorHistory []ColorHistory

func main() {
	app := gtk.NewApplication("com.github.wayland-color-picker-gtk4", 0)

	app.ConnectActivate(func() {
		// Load color historyf
		loadColorHistory()

		// Main window
		win := gtk.NewApplicationWindow(app)
		win.SetTitle("Wayland Color Picker GTK4")
		win.SetDefaultSize(400, 500)

		win.ConnectCloseRequest(func() bool {
			saveColorHistory()
			return false
		})

		// Main box layout
		box := gtk.NewBox(gtk.OrientationVertical, 10)
		box.SetMarginTop(20)
		box.SetMarginBottom(20)
		box.SetMarginStart(20)
		box.SetMarginEnd(20)
		win.SetChild(box)

		// Color info container for main display
		colorInfoMainBox := gtk.NewBox(gtk.OrientationVertical, 5)
		box.Append(colorInfoMainBox)

		// Main color label
		colorLabel := gtk.NewLabel("No color selected")
		colorInfoMainBox.Append(colorLabel)

		// Color name label for main display
		colorNameLabel := gtk.NewLabel("")
		colorNameLabel.SetMarkup("<i>Select a color to see its name</i>")
		colorInfoMainBox.Append(colorNameLabel)

		// Color preview box
		colorPreview := gtk.NewDrawingArea()
		colorPreview.SetSizeRequest(100, 100)
		colorPreview.SetHExpand(false)
		colorPreview.SetVExpand(false)
		box.Append(colorPreview)

		// CSS provider for styling
		provider := gtk.NewCSSProvider()
		display := gdk.DisplayGetDefault()
		gtk.StyleContextAddProviderForDisplay(display, provider, gtk.STYLE_PROVIDER_PRIORITY_USER)

		// Current color state
		var currentColor [3]uint8

		// Set up drawing function for color preview
		colorPreview.SetDrawFunc(func(da *gtk.DrawingArea, cr *cairo.Context, width, height int) {
			// Draw the color rectangle
			r := float64(currentColor[0]) / 255.0
			g := float64(currentColor[1]) / 255.0
			b := float64(currentColor[2]) / 255.0
			cr.SetSourceRGB(r, g, b)
			cr.Rectangle(0, 0, float64(width), float64(height))
			cr.Fill()
		})

		// Button
		button := gtk.NewButtonWithLabel("Pick Color")
		box.Append(button)

		// History section
		historyLabel := gtk.NewLabel("Color History:")
		historyLabel.SetHAlign(gtk.AlignStart)
		box.Append(historyLabel)

		// Scrolled window for history
		scrolled := gtk.NewScrolledWindow()
		scrolled.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
		scrolled.SetSizeRequest(-1, 200)
		scrolled.SetVExpand(true)
		box.Append(scrolled)

		// History container
		historyBox := gtk.NewBox(gtk.OrientationVertical, 5)
		scrolled.SetChild(historyBox)

		// Function to update history display
		updateHistoryDisplay := func() {
			// Clear existing children
			// Collect all children first, then remove them
			var children []gtk.Widgetter
			for child := historyBox.FirstChild(); child != nil; {
				children = append(children, child)
				// Get the actual widget to access NextSibling
				if w, ok := child.(interface{ NextSibling() gtk.Widgetter }); ok {
					child = w.NextSibling()
				} else {
					break
				}
			}
			for _, child := range children {
				historyBox.Remove(child)
			}

			// Add history items (most recent first)
			for i := len(colorHistory) - 1; i >= 0; i-- {
				color := colorHistory[i]

				// Create row container
				row := gtk.NewBox(gtk.OrientationHorizontal, 10)
				row.SetMarginTop(2)
				row.SetMarginBottom(2)

				// Color preview square
				colorSquare := gtk.NewDrawingArea()
				colorSquare.SetSizeRequest(30, 30)

				// Capture color for closure
				rgb := color.RGB
				colorSquare.SetDrawFunc(func(da *gtk.DrawingArea, cr *cairo.Context, width, height int) {
					r := float64(rgb[0]) / 255.0
					g := float64(rgb[1]) / 255.0
					b := float64(rgb[2]) / 255.0
					cr.SetSourceRGB(r, g, b)
					cr.Rectangle(0, 0, float64(width), float64(height))
					cr.Fill()

					// Draw border
					cr.SetSourceRGB(0.5, 0.5, 0.5)
					cr.SetLineWidth(1)
					cr.Rectangle(0, 0, float64(width), float64(height))
					cr.Stroke()
				})

				// Frame for the color square
				colorFrame := gtk.NewFrame("")
				colorFrame.SetChild(colorSquare)
				colorFrame.SetSizeRequest(32, 32)

				// Color info container
				colorInfoBox := gtk.NewBox(gtk.OrientationVertical, 2)
				colorInfoBox.SetHExpand(true)

				// HEX color label
				hexLabel := gtk.NewLabel(color.Color)
				hexLabel.SetHAlign(gtk.AlignStart)
				hexLabel.SetMarkup(fmt.Sprintf("<b>%s</b>", color.Color))
				colorInfoBox.Append(hexLabel)

				// Color name label (if available)
				if color.ColorName != "" {
					nameLabel := gtk.NewLabel(color.ColorName)
					nameLabel.SetHAlign(gtk.AlignStart)
					nameLabel.SetMarkup(fmt.Sprintf("<small><i>%s</i></small>", color.ColorName))
					colorInfoBox.Append(nameLabel)
				}

				// Copy button
				copyBtn := gtk.NewButtonWithLabel("Copy")
				copyBtn.SetSizeRequest(60, -1)

				// Make variables available to closure
				colorStr := color.Color
				copyBtn.ConnectClicked(func() {
					clipboard := display.Clipboard()
					clipboard.SetText(colorStr)
					colorLabel.SetText("Copied: " + colorStr)
				})

				// Click gesture for color frame
				clickGesture := gtk.NewGestureClick()
				clickGesture.ConnectPressed(func(n int, x, y float64) {
					clipboard := display.Clipboard()
					clipboard.SetText(colorStr)
					colorLabel.SetText("Copied: " + colorStr)
				})
				colorFrame.AddController(clickGesture)

				row.Append(colorFrame)
				row.Append(colorInfoBox)
				row.Append(copyBtn)

				historyBox.Append(row)
			}
		}

		// Initial history display
		updateHistoryDisplay()

		button.ConnectClicked(func() {
			go func() {
				color, ok := pickColor()
				if ok {
					glib.IdleAdd(func() {
						colorStr := fmt.Sprintf("#%02X%02X%02X", color[0], color[1], color[2])

						// Get color name
						colorName, err := colornameconv.New(colorStr)
						if err != nil {
							colorName = ""
						}

						colorLabel.SetText("Selected: " + colorStr)

						// Update color name label
						if colorName != "" {
							colorNameLabel.SetMarkup(fmt.Sprintf("<i>%s</i>", colorName))
						} else {
							colorNameLabel.SetMarkup("<i>Color name not available</i>")
						}

						// Update current color and redraw preview
						currentColor = color
						colorPreview.QueueDraw()

						// Add to history
						addToHistory(ColorHistory{
							Color:     colorStr,
							ColorName: colorName,
							RGB:       color,
							Timestamp: time.Now().Unix(),
						})

						// Update history display
						updateHistoryDisplay()

						// Copy to clipboard
						clipboard := display.Clipboard()
						clipboard.SetText(colorStr)
					})
				} else {
					glib.IdleAdd(func() {
						colorLabel.SetText("Color picking cancelled or failed")
					})
				}
			}()
		})

		win.Present()
	})

	app.Run(os.Args)
}

// pickColor attempts multiple methods to pick a color
func pickColor() ([3]uint8, bool) {
	// Try XDG Portal first (works on GNOME and KDE)
	if color, ok := pickColorViaPortal(); ok {
		return color, true
	}

	// Detect window manager and try specific methods
	wm := detectWindowManager()
	log.Printf("Detected window manager: %s", wm)

	switch wm {
	case "hyprland":
		if color, ok := pickColorHyprland(); ok {
			return color, true
		}
	case "sway":
		if color, ok := pickColorSway(); ok {
			return color, true
		}
	}

	// Fall back to grim + slurp (works on most wlroots compositors)
	if color, ok := pickColorGrimSlurp(); ok {
		return color, true
	}

	return [3]uint8{}, false
}

// detectWindowManager detects the current window manager
func detectWindowManager() string {
	// Check for Hyprland
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" {
		return "hyprland"
	}

	// Check for Sway
	if os.Getenv("SWAYSOCK") != "" {
		return "sway"
	}

	// Check XDG_CURRENT_DESKTOP
	desktop := strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))
	if strings.Contains(desktop, "hyprland") {
		return "hyprland"
	}
	if strings.Contains(desktop, "sway") {
		return "sway"
	}

	return "unknown"
}

// pickColorHyprland uses Hyprland's IPC to pick a color
func pickColorHyprland() ([3]uint8, bool) {
	var color [3]uint8

	// Use hyprpicker if available
	cmd := exec.Command("hyprpicker", "-a")
	output, err := cmd.Output()
	if err != nil {
		log.Println("hyprpicker failed:", err)
		return color, false
	}

	hexColor := strings.TrimSpace(string(output))
	if !strings.HasPrefix(hexColor, "#") {
		log.Println("Invalid color format from hyprpicker:", hexColor)
		return color, false
	}

	return parseHexColor(hexColor)
}

// pickColorSway uses Sway's IPC with grim and slurp
func pickColorSway() ([3]uint8, bool) {
	return pickColorGrimSlurp()
}

// pickColorGrimSlurp uses grim and slurp to pick a color
func pickColorGrimSlurp() ([3]uint8, bool) {
	var color [3]uint8

	// Check if grim and slurp are available
	if _, err := exec.LookPath("grim"); err != nil {
		log.Println("grim not found")
		return color, false
	}
	if _, err := exec.LookPath("slurp"); err != nil {
		log.Println("slurp not found")
		return color, false
	}

	// Create a temporary file for the screenshot
	tmpfile, err := os.CreateTemp("", "colorpick-*.png")
	if err != nil {
		log.Println("Failed to create temp file:", err)
		return color, false
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Use slurp to select a pixel and pipe to grim
	slurpCmd := exec.Command("slurp", "-p")
	slurpOutput, err := slurpCmd.Output()
	if err != nil {
		log.Println("slurp failed or cancelled:", err)
		return color, false
	}

	// Parse slurp output
	coords := strings.TrimSpace(string(slurpOutput))
	if coords == "" {
		log.Println("slurp returned empty coordinates")
		return color, false
	}

	log.Printf("slurp output: %s", coords)

	spaceParts := strings.Split(coords, " ")
	coordPart := spaceParts[0]

	parts := strings.Split(coordPart, ",")
	if len(parts) != 2 {
		log.Println("Invalid slurp coordinate format:", coords)
		return color, false
	}

	// Parse coordinates
	x := strings.TrimSpace(parts[0])
	y := strings.TrimSpace(parts[1])

	// Format for grim
	grimGeometry := fmt.Sprintf("%s,%s 1x1", x, y)

	// Take screenshot of 1x1 pixel at the selected coordinates
	grimCmd := exec.Command("grim", "-g", grimGeometry, tmpfile.Name())
	grimOutput, err := grimCmd.CombinedOutput()
	if err != nil {
		log.Printf("grim with 1x1 geometry failed: %v\nOutput: %s\nGeometry: %s", err, string(grimOutput), grimGeometry)
		log.Println("Trying alternative method: full screen capture")

		color, ok := pickColorGrimAlternative(x, y)
		if ok {
			return color, true
		}

		log.Println("Alternative method also failed")
		return color, false
	}

	log.Printf("grim succeeded, reading pixel from: %s", tmpfile.Name())

	return readPixelColor(tmpfile.Name())
}

// pickColorGrimAlternative captures full screen then reads specific pixel
func pickColorGrimAlternative(x, y string) ([3]uint8, bool) {
	var color [3]uint8

	tmpfile, err := os.CreateTemp("", "colorpick-full-*.png")
	if err != nil {
		log.Println("Failed to create temp file:", err)
		return color, false
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	grimCmd := exec.Command("grim", tmpfile.Name())
	if err := grimCmd.Run(); err != nil {
		log.Println("Full screen grim failed:", err)
		return color, false
	}

	convertCmd := exec.Command("convert", tmpfile.Name(), "-format", fmt.Sprintf("%%[pixel:p{%s,%s}]", x, y), "info:")
	output, err := convertCmd.Output()
	if err != nil {
		log.Println("convert failed:", err)
		return color, false
	}

	colorStr := strings.TrimSpace(string(output))
	return parseRGBColor(colorStr)
}

// readPixelColor reads color from a 1x1 pixel image
func readPixelColor(filename string) ([3]uint8, bool) {
	var color [3]uint8

	convertCmd := exec.Command("convert", filename, "-format", "%[pixel:p{0,0}]", "info:")
	output, err := convertCmd.Output()
	if err != nil {
		log.Println("ImageMagick convert failed, trying magick:", err)

		magickCmd := exec.Command("magick", filename, "-format", "%[pixel:p{0,0}]", "info:")
		output, err = magickCmd.Output()
		if err != nil {
			log.Println("magick command also failed:", err)
			return color, false
		}
	}

	colorStr := strings.TrimSpace(string(output))
	return parseRGBColor(colorStr)
}

// parseHexColor parses a hex color string like "#RRGGBB"
func parseHexColor(hexColor string) ([3]uint8, bool) {
	var color [3]uint8

	hexColor = strings.TrimPrefix(hexColor, "#")
	if len(hexColor) != 6 {
		return color, false
	}

	r, err1 := strconv.ParseUint(hexColor[0:2], 16, 8)
	g, err2 := strconv.ParseUint(hexColor[2:4], 16, 8)
	b, err3 := strconv.ParseUint(hexColor[4:6], 16, 8)

	if err1 != nil || err2 != nil || err3 != nil {
		return color, false
	}

	color[0] = uint8(r)
	color[1] = uint8(g)
	color[2] = uint8(b)

	return color, true
}

// parseRGBColor parses RGB color strings like "srgb(255,128,64)" or "rgb(255,128,64)"
func parseRGBColor(rgbStr string) ([3]uint8, bool) {
	var color [3]uint8

	rgbStr = strings.TrimPrefix(rgbStr, "srgb(")
	rgbStr = strings.TrimPrefix(rgbStr, "rgb(")
	rgbStr = strings.TrimSuffix(rgbStr, ")")

	parts := strings.Split(rgbStr, ",")
	if len(parts) != 3 {
		return color, false
	}

	r, err1 := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 8)
	g, err2 := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 8)
	b, err3 := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 8)

	if err1 != nil || err2 != nil || err3 != nil {
		return color, false
	}

	color[0] = uint8(r)
	color[1] = uint8(g)
	color[2] = uint8(b)

	return color, true
}

// addToHistory adds a color to the history, avoiding duplicates
func addToHistory(newColor ColorHistory) {
	for _, existing := range colorHistory {
		if existing.Color == newColor.Color {
			return
		}
	}

	colorHistory = append(colorHistory, newColor)

	if len(colorHistory) > 20 {
		colorHistory = colorHistory[1:]
	}

	saveColorHistory()
}

// getConfigDir returns the config directory for the application
func getConfigDir() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".config")
	}
	appConfigDir := filepath.Join(configDir, "color-picker")
	os.MkdirAll(appConfigDir, 0755)
	return appConfigDir
}

// saveColorHistory saves the color history to a JSON file
func saveColorHistory() {
	configDir := getConfigDir()
	historyFile := filepath.Join(configDir, "history.json")

	data, err := json.MarshalIndent(colorHistory, "", "  ")
	if err != nil {
		log.Printf("Error marshaling history: %v", err)
		return
	}

	err = os.WriteFile(historyFile, data, 0644)
	if err != nil {
		log.Printf("Error saving history: %v", err)
	}
}

// loadColorHistory loads the color history from a JSON file
func loadColorHistory() {
	configDir := getConfigDir()
	historyFile := filepath.Join(configDir, "history.json")

	data, err := os.ReadFile(historyFile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Error reading history file: %v", err)
		}
		return
	}

	err = json.Unmarshal(data, &colorHistory)
	if err != nil {
		log.Printf("Error unmarshaling history: %v", err)
		colorHistory = []ColorHistory{}
	}
}

// pickColorViaPortal uses xdg-desktop-portal Screenshot.PickColor
func pickColorViaPortal() ([3]uint8, bool) {
	var color [3]uint8

	conn, err := dbus.SessionBus()
	if err != nil {
		log.Println("Failed to connect to D-Bus:", err)
		return color, false
	}

	obj := conn.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")

	handleToken := "go-color-picker"
	var handle dbus.ObjectPath
	call := obj.Call("org.freedesktop.portal.Screenshot.PickColor", 0, handleToken, map[string]dbus.Variant{})
	if call.Err != nil {
		log.Println("Portal PickColor call failed:", call.Err)
		return color, false
	}
	handle = call.Body[0].(dbus.ObjectPath)

	signalChan := make(chan *dbus.Signal, 5)
	conn.Signal(signalChan)

	matchRule := fmt.Sprintf("type='signal',sender='org.freedesktop.portal.Desktop',path='%s'", handle)
	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule)

	for sig := range signalChan {
		if sig.Path == handle && sig.Name == "org.freedesktop.portal.Request.Response" {
			response := sig.Body[0].(uint32)
			if response == 0 {
				results := sig.Body[1].(map[string]dbus.Variant)
				if c, ok := results["color"]; ok {
					switch c.Signature().String() {
					case "u":
						val := c.Value().(uint32)
						color[0] = uint8((val >> 16) & 0xFF)
						color[1] = uint8((val >> 8) & 0xFF)
						color[2] = uint8(val & 0xFF)
						return color, true
					case "(ddd)":
						if tuple, ok := c.Value().([]interface{}); ok && len(tuple) == 3 {
							r := tuple[0].(float64)
							g := tuple[1].(float64)
							b := tuple[2].(float64)
							color[0] = uint8(r * 255)
							color[1] = uint8(g * 255)
							color[2] = uint8(b * 255)
							return color, true
						}
					}
				}
			}
			break
		}
	}
	return color, false
}
