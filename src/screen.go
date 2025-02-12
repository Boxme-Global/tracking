package omisocial

type screenClass struct {
	minWidth uint16
	class    string
}

// ScreenClasses is a list of typical screen sizes used to group resolutions.
// Everything below is considered "XS" (tiny).
var ScreenClasses = []screenClass{
	{5120, "UHD 5K"},
	{3840, "UHD 4K"},
	{2560, "WQHD"},
	{1920, "Full HD"},
	{1280, "HD"},
	{1024, "XL"},
	{800, "L"},
	{600, "M"},
	{415, "S"},
}

// GetScreenClass returns the screen class for given width in pixels.
func GetScreenClass(width uint16) string {
	if width <= 0 {
		return ""
	}

	for _, class := range ScreenClasses {
		if width >= class.minWidth {
			return class.class
		}
	}

	return "XS"
}
