package text

import "github.com/tinne26/etxt"

var (
	R    *etxt.Renderer
	Font = "x14y20pxScoreDozer"
)

func LoadFonts(font []byte) {
	fontLib := etxt.NewFontLibrary()
	_, err := fontLib.ParseFontBytes(font)
	if err != nil {
		panic(err)
	}

	// create a new text renderer and configure it
	R = etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	R.SetCacheHandler(glyphsCache.NewHandler())
	R.SetFont(fontLib.GetFont(Font))
	R.SetAlign(etxt.YCenter, etxt.XCenter)
	R.SetSizePx(14)
}
