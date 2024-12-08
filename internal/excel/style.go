package excel

import "github.com/xuri/excelize/v2"

// styles
func (ew *ExcelWriter) SetStyles(f *excelize.File) error {
	// Set styles
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  24,
			Color: "000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{ew.blueColor},
			Pattern: 1,
		},
	})
	if err != nil {
		return err
	}
	ew.headerStyle = headerStyle

	titleStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  16,
			Color: "000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{ew.orangeColor},
			Pattern: 1,
		},
	})
	if err != nil {
		return err
	}
	ew.titleStyle = titleStyle

	bodyStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  14,
			Color: "000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{ew.backgroundColor},
			Pattern: 1,
		},
	})
	if err != nil {
		return err
	}
	ew.bodyStyle = bodyStyle

	// --------------------------------
	// --------- Device style ---------
	// --------------------------------
	deviceStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  14,
			Color: "000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{ew.greenColor},
			Pattern: 1,
		},
	})
	if err != nil {
		return err
	}
	ew.deviceStyle = deviceStyle

	// --------------------------------
	// --------- Part style ---------
	// --------------------------------
	partStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  14,
			Color: "000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{ew.redColor},
			Pattern: 1,
		},
	})
	if err != nil {
		return err
	}
	ew.partStyle = partStyle

	// --------------------------------
	// ------ Extra Price style -------
	// --------------------------------
	extraPriceStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  14,
			Color: "000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{ew.greenColor},
			Pattern: 1,
		},
	})
	if err != nil {
		return err
	}
	ew.extraPriceStyle = extraPriceStyle

	// --------------------------------
	// --------- Price style ----------
	// --------------------------------
	priceStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  14,
			Color: "000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{ew.yellowColor},
			Pattern: 1,
		},
	})
	if err != nil {
		return err
	}
	ew.priceStyle = priceStyle

	return nil
}
