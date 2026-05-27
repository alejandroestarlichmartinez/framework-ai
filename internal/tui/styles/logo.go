package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// logoLines contains the ASCII art for the Framework-AI wolf logo.
var logoLines = []string{
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣦⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣼⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣷⡄⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⣆⠀⠀⠀⠀⠀⢀⣴⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⣴⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣧⣀⣾⣿⣿⣿⣿⣿⣿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⢹⣿⣶⣦⣤⣀⡀⠀⠀⠀⠀⠀⣼⣿⣿⣿⡿⠿⠟⠛⠛⠿⠿⣿⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⢿⣿⣿⣿⣿⣿⣿⣶⣶⣤⣤⡿⠟⠉⢴⣶⣿⣿⣿⣿⣿⣷⣦⣍⠻⣿⣿⣿⡇⠀⠀⠀⠀⠀⣀⣀⣠⣤⣶⡶",
	"⠀⠀⠀⠀⠀⠈⣿⣿⣿⣿⣿⣿⣿⡿⠟⣋⣀⣙⡻⢶⣝⢿⣿⣿⣿⣿⣿⣿⣿⣿⣌⠻⣿⣷⣶⣶⣿⣿⣿⣿⣿⣿⣿⠏⠀",
	"⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⣿⠏⣴⣿⡿⠿⢿⣿⣦⡙⢦⣽⣿⣿⣿⣿⣿⣿⣿⣿⡧⠹⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⣿⣿⡆⢉⣥⣶⣾⣶⣌⠻⣿⣎⠻⣿⣿⣿⡿⠟⣋⣭⣴⣶⡄⢹⣿⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⣿⢃⣿⣿⡿⠿⠿⠿⣧⡙⢿⣷⣶⣶⣶⣶⣿⠿⠟⠋⣩⣴⡌⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⢀⣸⣿⣿⡟⢸⠟⣡⣶⣾⣿⣿⣶⣌⠲⣬⣉⠉⣉⣥⣴⣾⣿⣷⣦⡙⣧⢹⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⢀⣠⣴⣾⣿⣿⣿⣿⡇⡎⣼⣿⣿⣿⣿⣿⣿⠉⢢⢹⡿⢰⣿⣿⣿⣿⣿⣿⠉⣳⠈⢸⣿⣿⡋⠀⠀⠀⠀⠀⠀⠀",
	"⠠⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⢁⡇⣿⣿⣿⣿⣿⣿⣿⣿⣿⢸⡇⣾⣿⣿⣿⣿⣿⣿⣿⣿⠀⢸⣿⣿⣿⣷⣶⣤⣄⣀⣀⠀",
	"⠀⠀⠉⠻⢿⣿⣿⣿⣿⣿⣿⢸⡇⢿⣿⣿⣿⣿⣿⣿⣿⠇⣼⣧⠸⣿⣿⣿⣿⣿⣿⣿⡿⢠⢸⣿⣿⣿⣿⣿⣿⣿⣿⠟⠁",
	"⠀⠀⠀⠀⠀⠈⠛⢿⣿⣿⣿⢸⣿⣌⠻⢿⣿⣿⣿⡿⢋⣼⣿⣿⣧⡙⠿⣿⣿⣿⡿⠟⣡⣿⢸⣿⣿⣿⣿⣿⡿⠋⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⣠⣾⣿⣿⣿⣾⣿⣿⣿⣶⣤⣤⣤⣶⣿⠋⣿⣿⢻⣿⣷⣶⣤⣴⣶⣿⣿⣿⢸⣿⣿⣿⡿⠋⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⣠⣾⣿⣿⣿⣿⣿⢹⣿⣷⣬⣛⣛⠛⣛⣩⣽⠀⣿⣿⢀⣷⣬⣙⡛⠛⣛⣫⣴⣿⢸⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠠⣾⣿⣿⣿⣿⣿⣿⠟⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⣿⣿⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⠘⢿⣿⣷⡀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠈⠙⠻⢿⣿⣿⢃⣾⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣌⣡⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⣦⡙⣿⣿⣷⣤⣀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠈⣿⢸⣿⡇⣿⣿⣿⣿⡿⢿⣿⡿⢻⣟⢹⡟⢻⣟⠻⣿⣿⣿⣿⣿⣿⣿⢸⣿⡇⣿⣿⣿⠿⠟⠁⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⢰⣿⣦⡙⠇⢸⣿⣿⡟⡰⠁⠈⠁⠀⠁⠀⠀⠀⠁⠀⠉⠀⠙⣌⢻⣿⣿⠘⣋⣴⠉⠁⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⢀⣿⣿⣿⣿⣷⡌⣿⣿⢰⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⡌⣿⡇⣼⣿⣿⡆⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠘⠛⠛⠛⠛⠻⣷⠹⣿⠸⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⡇⣿⢡⣿⣿⣿⣷⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⣇⢻⣧⡙⠿⠶⠴⢦⡀⠶⣶⣶⡶⠆⢠⣤⠴⢏⣴⢃⡎⠀⠈⠉⠉⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠟⠛⠀⠻⣿⣷⣶⣾⣿⣿⡇⢹⠏⣴⣶⣶⣶⣶⡿⠃⠚⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠻⣿⣿⣿⣿⣿⣤⣾⣿⣿⣿⣿⠟⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠛⠛⠻⠿⠿⠛⠛⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
}

// BrandYellow is the primary Framework-AI branding color (#f7df1e — JS yellow).
const BrandYellow = lipgloss.Color("#f7df1e")

// gradientColors defines the top-to-bottom gradient for the logo.
// Distributed across rows: light yellow → gold → brand yellow → amber → deep gold.
var gradientColors = []lipgloss.Color{
	lipgloss.Color("#ffeb3b"), // band 1 — light yellow
	lipgloss.Color("#fdd835"), // band 2 — golden yellow
	BrandYellow,               // band 3 — brand yellow
	lipgloss.Color("#f9a825"), // band 4 — amber
	lipgloss.Color("#f57f17"), // band 5 — deep gold
}

// RenderLogo returns the ASCII logo with a top-to-bottom gradient.
func RenderLogo() string {
	total := len(logoLines)
	if total == 0 {
		return ""
	}

	bands := len(gradientColors)
	var b strings.Builder

	for i, line := range logoLines {
		bandIdx := (i * bands) / total
		if bandIdx >= bands {
			bandIdx = bands - 1
		}
		style := lipgloss.NewStyle().Foreground(gradientColors[bandIdx])
		b.WriteString(style.Render(line))
		if i < total-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}
