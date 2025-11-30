# Wayland Color Picker (GTK4)

A simple and elegant color picker application for Wayland desktop environments, built with GTK4 and Go.

![Wayland Color Picker](icon.png)

## Features

- Pick colors from anywhere on your screen
- Automatic clipboard copy
- Display color names (e.g., "Sky Blue", "Crimson")
- Color history (stores last 20 colors)
- Click to copy colors from history
- Persistent storage between sessions

## Supported Wayland Compositors

- **GNOME** (via XDG Desktop Portal)
- **KDE Plasma** (via XDG Desktop Portal)
- **Hyprland** (requires `hyprpicker`)
- **Sway** (via `grim` + `slurp`)
- **Other wlroots-based compositors** (via `grim` + `slurp`)

## Dependencies

### Build Dependencies
- Go 1.21 or later
- GTK4 development files
- Cairo development files

### Runtime Dependencies
- GTK4
- One of the following, depending on your compositor:
  - `xdg-desktop-portal` (for GNOME/KDE)
  - `hyprpicker` (for Hyprland)
  - `grim` + `slurp` (for Sway and other wlroots compositors)
  - `imagemagick` (for pixel color extraction)

### Installation of Dependencies

**Arch Linux:**
```bash
sudo pacman -S go gtk4 cairo
# For Hyprland:
sudo pacman -S hyprpicker
# For Sway/wlroots:
sudo pacman -S grim slurp imagemagick
# For GNOME/KDE:
sudo pacman -S xdg-desktop-portal-gtk
```

**Ubuntu/Debian:**
```bash
sudo apt install golang gtk4-dev libcairo2-dev
# For Sway/wlroots:
sudo apt install grim slurp imagemagick
# For GNOME/KDE:
sudo apt install xdg-desktop-portal-gtk
```

**Fedora:**
```bash
sudo dnf install golang gtk4-devel cairo-devel
# For Sway/wlroots:
sudo dnf install grim slurp ImageMagick
# For GNOME/KDE:
sudo dnf install xdg-desktop-portal-gtk
```

## Building

1. Clone the repository:
```bash
git clone https://github.com/yourusername/wayland-color-picker-gtk4.git
cd wayland-color-picker-gtk4
```

2. Download Go dependencies:
```bash
make deps
```

3. Build the application:
```bash
make build
```

## Installation

### System-wide Installation (Recommended)

Build and install in one step:
```bash
make build-install
```

Or separately:
```bash
make build
sudo make install
```

This will install:
- Binary to `/usr/bin/wayland-color-picker-gtk4`
- Desktop file to `/usr/share/applications/`
- Icons to `/usr/share/icons/hicolor/` and `/usr/share/pixmaps/`

After installation, you can launch the app from your application menu or by running:
```bash
wayland-color-picker-gtk4
```

### Uninstall

```bash
sudo make uninstall
```

## Configuration

Color history is stored in:
```
~/.config/color-picker/history.json
```

## Troubleshooting

### Color picking doesn't work

Make sure you have the appropriate tools installed for your compositor:
- **Hyprland:** Install `hyprpicker`
- **Sway/wlroots:** Install `grim`, `slurp`, and `imagemagick`
- **GNOME/KDE:** Install `xdg-desktop-portal-gtk` or `xdg-desktop-portal-gnome`

### Icon not showing in dock/taskbar

This is usually resolved by logging out and back in after installation. The desktop file needs to be registered by your desktop environment.

### Permission errors during installation

Make sure to build as a regular user first, then install with sudo:
```bash
make build        # as regular user
sudo make install # with sudo
```

## Development

### Run without installing
```bash
make run
```

### Build with debug symbols
```bash
make debug
```

### Code checks
```bash
make check
```

### Clean build artifacts
```bash
make clean
```

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- Built with [gotk4](https://github.com/diamondburned/gotk4)
- Color naming provided by [colornameconv](https://github.com/dachinat/colornameconv)
- Inspired by various color picker tools in the Linux ecosystem
