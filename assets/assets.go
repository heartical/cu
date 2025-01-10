package assets

import (
	"embed"
)

//go:embed *.html
//go:embed images/*.xml
//go:embed images/*.png
//go:embed fonts/*.ttf
var resources embed.FS

var (
	IndexHTML         []byte
	UIPackRPGSheetPNG []byte
	UIPackRPGSheetXML []byte

	UIPackSpaceSheetPNG []byte
	UIPackSpaceSheetXML []byte

	Font []byte
)

// Init инициализирует ресурсы.
func Init() error {
	if err := initResources(); err != nil {
		return err
	}
	return nil
}

func initResources() error {
	resourceNames := []struct {
		name string
		file *[]byte
	}{
		{"_index.html", &IndexHTML},
		{"images/uipack_rpg_sheet.xml", &UIPackRPGSheetXML},
		{"images/uipack_rpg_sheet.png", &UIPackRPGSheetPNG},

		{"images/uipackSpace_sheet.xml", &UIPackSpaceSheetXML},
		{"images/uipackSpace_sheet.png", &UIPackSpaceSheetPNG},

		{"fonts/x14y20pxScoreDozer.ttf", &Font},
	}

	for _, v := range resourceNames {
		buf, err := resources.ReadFile(v.name)
		if err != nil {
			return err
		}
		*v.file = buf
	}

	return nil
}
