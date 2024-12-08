package excel

import (
	"fmt"
	"strconv"

	devicepart "github.com/poriamsz55/BoosterPump-webapp/internal/models/device_part"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/project"
	"github.com/xuri/excelize/v2"
)

type ExcelWriter struct {
	// Define colors as constants
	greenColor      string
	redColor        string
	blueColor       string
	orangeColor     string
	backgroundColor string
	yellowColor     string

	// Define styles as constants
	headerStyle     int
	titleStyle      int
	bodyStyle       int
	deviceStyle     int
	partStyle       int
	extraPriceStyle int
	priceStyle      int
}

func NewExcelWriter() *ExcelWriter {
	return &ExcelWriter{
		greenColor:      "92D050",
		redColor:        "FF0000",
		blueColor:       "0000FF",
		orangeColor:     "FFA500",
		backgroundColor: "D9D9D9",
		yellowColor:     "FFFF99",
	}
}

func GenerateProjectReport(project *project.Project, fileName string) error {

	ew := NewExcelWriter()

	// Create a new Excel file
	f := excelize.NewFile()

	// Create Project Report Sheet
	sheetName := "گزارش پروژه"
	f.NewSheet(sheetName)

	err := ew.SetStyles(f)
	if err != nil {
		return err
	}

	// Write header
	rowi := 1
	coli := 1
	excelTitle := fmt.Sprintf("پروژه %s", project.Name)
	f.SetCellValue(sheetName, getCellName(coli, rowi), excelTitle)
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+13, rowi+3), ew.headerStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+13, rowi+3))

	// Write device header
	rowi += 4
	f.SetCellValue(sheetName, getCellName(coli, rowi), "دستگاه ها")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+13, rowi+2), ew.deviceStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+13, rowi+2))

	// Write title row
	rowi += 3
	titles := []string{"نام دستگاه", "تبدیل", "صافی", "تعداد(مقدار)", "بهای واحد(ریال)", "بهای کل(ریال)"}
	colSpans := []int{3, 2, 1, 2, 3, 3}

	currentCol := coli
	for i, title := range titles {
		f.SetCellValue(sheetName, getCellName(currentCol, rowi), title)
		f.SetCellStyle(sheetName, getCellName(currentCol, rowi), getCellName(currentCol+colSpans[i]-1, rowi), ew.titleStyle)
		if colSpans[i] > 1 {
			f.MergeCell(sheetName, getCellName(currentCol, rowi), getCellName(currentCol+colSpans[i]-1, rowi))
		}
		currentCol += colSpans[i]
	}

	// -----------------------
	// ---- Write devices ----
	// -----------------------
	rowi++
	var totalDevicePrice uint64

	for _, projectDevice := range project.ProjectDeviceList {
		coli := 1
		device := projectDevice.Device

		// Device name
		f.SetCellValue(sheetName, getCellName(coli, rowi), device.Name)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

		// Converter
		coli += 3
		f.SetCellValue(sheetName, getCellName(coli, rowi), device.Converter)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+1, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+1, rowi))

		// Filter
		coli += 2
		filterText := "ندارد"
		if device.Filter {
			filterText = "دارد"
		}
		f.SetCellValue(sheetName, getCellName(coli, rowi), filterText)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)

		// Count
		coli += 1
		f.SetCellValue(sheetName, getCellName(coli, rowi), projectDevice.Count)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+1, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+1, rowi))

		// Price per unit
		coli += 2
		f.SetCellValue(sheetName, getCellName(coli, rowi), device.Price)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

		// Total price
		coli += 3
		deviceTotalPrice := projectDevice.Price
		f.SetCellValue(sheetName, getCellName(coli, rowi), deviceTotalPrice)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

		totalDevicePrice += deviceTotalPrice
		rowi++
	}

	// Write total devices price row
	coli = 1
	f.SetCellValue(sheetName, getCellName(coli, rowi), "جمع کل بهای دستگاه ها(ریال)")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+10, rowi), ew.titleStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+10, rowi))

	coli += 11
	f.SetCellValue(sheetName, getCellName(coli, rowi), totalDevicePrice)
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.priceStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

	// -----------------------
	// ---- Write extras ----
	// -----------------------
	rowi += 2
	coli = 1

	// Create Extra Price Header
	// Write "Extra Prices" header
	f.SetCellValue(sheetName, getCellName(coli, rowi), "هزینه های اضافی")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+5, rowi+2), ew.extraPriceStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+5, rowi+2))

	// Write column headers
	rowi += 3
	// Title column
	f.SetCellValue(sheetName, getCellName(coli, rowi), "عنوان")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.titleStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

	// Price column
	coli += 3
	f.SetCellValue(sheetName, getCellName(coli, rowi), "هزینه(ریال)")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.titleStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

	// Write extra prices data
	rowi++
	var totalExtraPrice uint64

	for _, extraPrice := range project.ExtraPriceList {
		coli = 1

		// Write extra price name
		f.SetCellValue(sheetName, getCellName(coli, rowi), extraPrice.Name)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

		// Write extra price amount
		coli += 3
		f.SetCellValue(sheetName, getCellName(coli, rowi), extraPrice.Price)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

		totalExtraPrice += extraPrice.Price
		rowi++
	}

	// Write total extra prices row
	coli = 1
	f.SetCellValue(sheetName, getCellName(coli, rowi), "جمع کل(ریال)")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.titleStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

	coli += 3
	f.SetCellValue(sheetName, getCellName(coli, rowi), totalExtraPrice)
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.priceStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

	// Write project total price
	coli = 1
	rowi += 2
	f.SetCellValue(sheetName, getCellName(coli, rowi), "بهای کل پروژه(ریال)")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+10, rowi), ew.titleStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+10, rowi))

	coli += 11
	f.SetCellValue(sheetName, getCellName(coli, rowi), project.Price)
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.priceStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

	// ---------------------
	// ---- Write parts ----
	// ---------------------
	// Create Part Report Sheet
	partSheetName := "گزارش قطعات"
	f.NewSheet(partSheetName)

	// Write header
	rowi = 1
	coli = 1
	f.SetCellValue(partSheetName, getCellName(coli, rowi), "قطعات")
	f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli+12, rowi+3), ew.deviceStyle)
	f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli+12, rowi+3))

	// Write column headers
	rowi += 4
	titlesStruct := []struct {
		title   string
		colSpan int
	}{
		{"نام قطعه", 3},
		{"سایز", 1},
		{"جنس", 1},
		{"برند", 2},
		{"تعداد(مقدار)", 2},
		{"بهای واحد(ریال)", 2},
		{"بهای کل(ریال)", 2},
	}

	currentCol = coli
	for _, title := range titlesStruct {
		f.SetCellValue(partSheetName, getCellName(currentCol, rowi), title.title)
		f.SetCellStyle(partSheetName, getCellName(currentCol, rowi), getCellName(currentCol+title.colSpan-1, rowi), ew.titleStyle)
		if title.colSpan > 1 {
			f.MergeCell(partSheetName, getCellName(currentCol, rowi), getCellName(currentCol+title.colSpan-1, rowi))
		}
		currentCol += title.colSpan
	}

	// Create map to merge parts
	partSummaryMap := make(map[int]*devicepart.DevicePartMerged)

	// Collect and merge parts from all devices
	for _, projectDevice := range project.ProjectDeviceList {
		for _, devicePart := range projectDevice.Device.DevicePartList {
			if merged, exists := partSummaryMap[devicePart.Part.Id]; !exists {
				// Create new merged part
				partSummaryMap[devicePart.Part.Id] = devicepart.NewDevicePartMerged(
					devicePart,
					devicePart.Count*projectDevice.Count,
				)
			} else {
				// Update existing merged part
				merged.MergeDevicePart(devicePart, projectDevice.Count)
			}
		}
	}

	// Write merged parts data
	rowi++
	var totalPartsPrice uint64

	for _, mergedPart := range partSummaryMap {
		coli = 1
		part := mergedPart.DevicePart.Part

		// Part name
		f.SetCellValue(partSheetName, getCellName(coli, rowi), part.Name)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.bodyStyle)
		f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

		// Size
		coli += 3
		f.SetCellValue(partSheetName, getCellName(coli, rowi), part.Size)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)

		// Material
		coli += 1
		f.SetCellValue(partSheetName, getCellName(coli, rowi), part.Material)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)

		// Brand
		coli += 1
		f.SetCellValue(partSheetName, getCellName(coli, rowi), part.Brand)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli+1, rowi), ew.bodyStyle)
		f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli+1, rowi))

		// Count
		coli += 2
		f.SetCellValue(partSheetName, getCellName(coli, rowi), mergedPart.Count)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli+1, rowi), ew.bodyStyle)
		f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli+1, rowi))

		// Price per unit
		coli += 2
		f.SetCellValue(partSheetName, getCellName(coli, rowi), part.Price)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli+1, rowi), ew.bodyStyle)
		f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli+1, rowi))

		// Total price
		coli += 2
		f.SetCellValue(partSheetName, getCellName(coli, rowi), mergedPart.Price)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli+1, rowi), ew.bodyStyle)
		f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli+1, rowi))

		totalPartsPrice += mergedPart.Price
		rowi++
	}

	// Write total parts price
	coli = 1
	f.SetCellValue(partSheetName, getCellName(coli, rowi), "جمع کل(ریال)")
	f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli+9, rowi), ew.titleStyle)
	f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli+9, rowi))

	coli += 10
	f.SetCellValue(partSheetName, getCellName(coli, rowi), totalPartsPrice)
	f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.priceStyle)
	f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli+2, rowi))

	// Adjust column widths for better readability
	// You may need to adjust these values based on your needs
	columns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M"}
	for _, col := range columns {
		f.SetColWidth(partSheetName, col, col, 15)
	}

	// -----------------------
	// ---- Device report ----
	// -----------------------
	// Create individual device report sheets
	for _, projectDevice := range project.ProjectDeviceList {
		device := projectDevice.Device

		// Create sheet name
		deviceSheetText := fmt.Sprintf("%s %s ", device.Name, device.Converter)
		if device.Filter {
			deviceSheetText += "با صافی"
		} else {
			deviceSheetText += "بدون صافی"
		}

		// Create new sheet
		f.NewSheet(deviceSheetText)

		// Write header
		rowi = 1
		coli = 1
		headerText := fmt.Sprintf("%s برای %s - %d عدد", deviceSheetText, project.Name, projectDevice.Count)
		f.SetCellValue(deviceSheetText, getCellName(coli, rowi), headerText)
		f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli+12, rowi+3), ew.deviceStyle)
		f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli+12, rowi+3))

		// Write column headers
		rowi += 4
		titles := []struct {
			title   string
			colSpan int
		}{
			{"اقلام", 3},
			{"سایز", 1},
			{"جنس", 1},
			{"برند", 2},
			{"تعداد(مقدار)", 2},
			{"بهای واحد(ریال)", 2},
			{"بهای کل(ریال)", 2},
		}

		currentCol := coli
		for _, title := range titles {
			f.SetCellValue(deviceSheetText, getCellName(currentCol, rowi), title.title)
			f.SetCellStyle(deviceSheetText, getCellName(currentCol, rowi), getCellName(currentCol+title.colSpan-1, rowi), ew.titleStyle)
			if title.colSpan > 1 {
				f.MergeCell(deviceSheetText, getCellName(currentCol, rowi), getCellName(currentCol+title.colSpan-1, rowi))
			}
			currentCol += title.colSpan
		}

		// Write device parts data
		rowi++
		var deviceTotalPrice uint64

		for _, devicePart := range device.DevicePartList {
			coli = 1
			part := devicePart.Part

			// Part name
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), part.Name)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.bodyStyle)
			f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli+2, rowi))

			// Size
			coli += 3
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), part.Size)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)

			// Material
			coli += 1
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), part.Material)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)

			// Brand
			coli += 1
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), part.Brand)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli+1, rowi), ew.bodyStyle)
			f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli+1, rowi))

			// Count
			coli += 2
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), devicePart.Count)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli+1, rowi), ew.bodyStyle)
			f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli+1, rowi))

			// Price per unit
			coli += 2
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), part.Price)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli+1, rowi), ew.bodyStyle)
			f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli+1, rowi))

			// Total price for this part
			coli += 2
			partTotalPrice := devicePart.Price
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), partTotalPrice)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli+1, rowi), ew.bodyStyle)
			f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli+1, rowi))

			deviceTotalPrice += partTotalPrice
			rowi++
		}

		// Write device total price
		coli = 1
		f.SetCellValue(deviceSheetText, getCellName(coli, rowi), "جمع کل(ریال)")
		f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli+9, rowi), ew.titleStyle)
		f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli+9, rowi))

		coli += 10
		f.SetCellValue(deviceSheetText, getCellName(coli, rowi), deviceTotalPrice)
		f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli+2, rowi), ew.priceStyle)
		f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli+2, rowi))

		// Adjust column widths for better readability
		columns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M"}
		for _, col := range columns {
			f.SetColWidth(deviceSheetText, col, col, 15)
		}

		type PointerBool *bool
		var pointerBool PointerBool = PointerBool(new(bool))
		*pointerBool = true
		// Set RTL direction for the sheet
		f.SetSheetView(deviceSheetText, 0, &excelize.ViewOptions{
			RightToLeft: pointerBool,
		})
	}

	type PointerBool *bool
	var pointerBool PointerBool = PointerBool(new(bool))
	*pointerBool = true
	// Set RTL direction for all sheets
	f.SetSheetView("گزارش پروژه", 0, &excelize.ViewOptions{
		RightToLeft: pointerBool,
	})
	f.SetSheetView("گزارش قطعات", 0, &excelize.ViewOptions{
		RightToLeft: pointerBool,
	})

	// Set the first sheet as active
	defaultSheet, err := f.GetSheetIndex("گزارش پروژه")
	if err != nil {
		return err
	}
	f.SetActiveSheet(defaultSheet)

	// Save the file to Downloads and get the path
	filePath, err := ew.SaveToDownloads(f, fileName)
	if err != nil {
		return fmt.Errorf("failed to save report: %v", err)
	}

	// Print success message with file path
	fmt.Printf("Excel file saved successfully at: %s\n", filePath)

	// Close the workbook
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close Excel file: %v", err)
	}

	return nil
}

func getCellName(col, row int) string {
	// Convert column number to Excel column name (A, B, C, ... AA, AB, etc.)
	colName := ""
	for col > 0 {
		col--
		colName = string('A'+col%26) + colName
		col /= 26
	}
	return colName + strconv.Itoa(row)
}
