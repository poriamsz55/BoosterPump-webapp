package excel

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"

	devicepart "github.com/poriamsz55/BoosterPump-webapp/internal/models/device_part"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/project"
	projectd "github.com/poriamsz55/BoosterPump-webapp/internal/models/project_device"
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

func cmpProjectDevice(a, b *projectd.ProjectDevice) int {
	return cmp.Compare(a.Device.Name, b.Device.Name)
}

func cmpPartSummary(a, b *devicepart.DevicePartMerged) int {
	return cmp.Compare(a.DevicePart.Part.Name, b.DevicePart.Part.Name)
}

func cmpDevicePart(a, b *devicepart.DevicePart) int {
	return cmp.Compare(a.Part.Name, b.Part.Name)
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
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+5, rowi+3), ew.headerStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+5, rowi+3))

	// Write device header
	rowi += 4
	f.SetCellValue(sheetName, getCellName(coli, rowi), "دستگاه ها")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+5, rowi+2), ew.deviceStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+5, rowi+2))

	// Write title row
	rowi += 3
	titles := []string{"نام دستگاه", "تبدیل", "صافی", "تعداد(مقدار)", "بهای واحد(ریال)", "بهای کل(ریال)"}
	colSpans := []int{1, 1, 1, 1, 1, 1}

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

	// sort project.ProjectDeviceList by name
	slices.SortFunc(project.ProjectDeviceList, cmpProjectDevice)

	for _, projectDevice := range project.ProjectDeviceList {
		coli := 1
		device := projectDevice.Device

		// Device name
		f.SetCellValue(sheetName, getCellName(coli, rowi), device.Name)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		// Converter
		coli += 1
		f.SetCellValue(sheetName, getCellName(coli, rowi), device.Converter)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		// Filter
		coli += 1
		filterText := "ندارد"
		if device.Filter {
			filterText = "دارد"
		}
		f.SetCellValue(sheetName, getCellName(coli, rowi), filterText)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)

		// Count
		coli += 1
		f.SetCellValue(sheetName, getCellName(coli, rowi), projectDevice.Count)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		// Price per unit
		coli += 1
		f.SetCellValue(sheetName, getCellName(coli, rowi), device.Price)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		// Total price
		coli += 1
		deviceTotalPrice := projectDevice.Price
		f.SetCellValue(sheetName, getCellName(coli, rowi), deviceTotalPrice)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		totalDevicePrice += deviceTotalPrice
		rowi++
	}

	// Write total devices price row
	coli = 1
	f.SetCellValue(sheetName, getCellName(coli, rowi), "جمع کل بهای دستگاه ها(ریال)")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+4, rowi), ew.titleStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+4, rowi))

	coli += 5
	f.SetCellValue(sheetName, getCellName(coli, rowi), totalDevicePrice)
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.priceStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

	// -----------------------
	// ---- Write extras ----
	// -----------------------
	rowi += 2
	coli = 1

	// Create Extra Price Header
	// Write "Extra Prices" header
	f.SetCellValue(sheetName, getCellName(coli, rowi), "هزینه های اضافی")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli+1, rowi+2), ew.extraPriceStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli+1, rowi+2))

	// Write column headers
	rowi += 3
	// Title column
	f.SetCellValue(sheetName, getCellName(coli, rowi), "عنوان")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.titleStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

	// Price column
	coli += 1
	f.SetCellValue(sheetName, getCellName(coli, rowi), "هزینه(ریال)")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.titleStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

	// Write extra prices data
	rowi++
	var totalExtraPrice uint64

	for _, extraPrice := range project.ExtraPriceList {
		coli = 1

		// Write extra price name
		f.SetCellValue(sheetName, getCellName(coli, rowi), extraPrice.Name)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		// Write extra price amount
		coli += 1
		f.SetCellValue(sheetName, getCellName(coli, rowi), extraPrice.Price)
		f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		totalExtraPrice += extraPrice.Price
		rowi++
	}

	// Write total extra prices row
	coli = 1
	f.SetCellValue(sheetName, getCellName(coli, rowi), "جمع کل(ریال)")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.titleStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

	coli += 1
	f.SetCellValue(sheetName, getCellName(coli, rowi), totalExtraPrice)
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.priceStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

	// Write project total price
	coli = 1
	rowi += 2
	f.SetCellValue(sheetName, getCellName(coli, rowi), "بهای کل پروژه(ریال)")
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.titleStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

	coli += 1
	f.SetCellValue(sheetName, getCellName(coli, rowi), project.Price)
	f.SetCellStyle(sheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.priceStyle)
	f.MergeCell(sheetName, getCellName(coli, rowi), getCellName(coli, rowi))

	// Adjust column widths for better readability
	// You may need to adjust these values based on your needs
	projectColumns := []string{"C", "D"}
	for _, col := range projectColumns {
		f.SetColWidth(sheetName, col, col, 15)
	}

	projectColumns = []string{"B", "E", "F"}
	for _, col := range projectColumns {
		f.SetColWidth(sheetName, col, col, 25)
	}

	f.SetColWidth(sheetName, "A", "A", 50)

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
	f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli+6, rowi+3), ew.deviceStyle)
	f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli+6, rowi+3))

	// Write column headers
	rowi += 4
	titlesStruct := []struct {
		title   string
		colSpan int
	}{
		{"نام قطعه", 1},
		{"سایز", 1},
		{"جنس", 1},
		{"برند", 1},
		{"تعداد(مقدار)", 1},
		{"بهای واحد(ریال)", 1},
		{"بهای کل(ریال)", 1},
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
	partSummaryList := []*devicepart.DevicePartMerged{}
	for _, mergedPart := range partSummaryMap {
		partSummaryList = append(partSummaryList, mergedPart)
	}

	// sort
	slices.SortFunc(partSummaryList, cmpPartSummary)

	for _, mergedPart := range partSummaryList {
		coli = 1
		part := mergedPart.DevicePart.Part

		// Part name
		f.SetCellValue(partSheetName, getCellName(coli, rowi), part.Name)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		// Size
		coli += 1
		f.SetCellValue(partSheetName, getCellName(coli, rowi), part.Size)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)

		// Material
		coli += 1
		f.SetCellValue(partSheetName, getCellName(coli, rowi), part.Material)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)

		// Brand
		coli += 1
		f.SetCellValue(partSheetName, getCellName(coli, rowi), part.Brand)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		// Count
		coli += 1
		f.SetCellValue(partSheetName, getCellName(coli, rowi), mergedPart.Count)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		// Price per unit
		coli += 1
		f.SetCellValue(partSheetName, getCellName(coli, rowi), part.Price)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		// Total price
		coli += 1
		f.SetCellValue(partSheetName, getCellName(coli, rowi), mergedPart.Price)
		f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
		f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi))

		totalPartsPrice += mergedPart.Price
		rowi++
	}

	// Write total parts price
	coli = 1
	f.SetCellValue(partSheetName, getCellName(coli, rowi), "جمع کل(ریال)")
	f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli+5, rowi), ew.titleStyle)
	f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli+5, rowi))

	coli += 6
	f.SetCellValue(partSheetName, getCellName(coli, rowi), totalPartsPrice)
	f.SetCellStyle(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi), ew.priceStyle)
	f.MergeCell(partSheetName, getCellName(coli, rowi), getCellName(coli, rowi))

	// Adjust column widths for better readability
	// You may need to adjust these values based on your needs
	columns := []string{"B", "C", "D", "E"}
	for _, col := range columns {
		f.SetColWidth(partSheetName, col, col, 15)
	}

	columns = []string{"F", "G"}
	for _, col := range columns {
		f.SetColWidth(partSheetName, col, col, 25)
	}

	f.SetColWidth(partSheetName, "A", "A", 50)

	// -----------------------
	// ---- Device report ----
	// -----------------------
	// Create individual device report sheets
	for pi, projectDevice := range project.ProjectDeviceList {
		device := projectDevice.Device

		deviceSheetText := fmt.Sprintf("دستگاه %v", pi)

		// Create new sheet
		f.NewSheet(deviceSheetText)

		// Write header
		var filterStr string
		if device.Filter {
			filterStr = "دارد"
		} else {
			filterStr = "ندارد"
		}

		headerName := fmt.Sprintf("%s %s %s", device.Name, device.Converter.String(), filterStr)
		rowi = 1
		coli = 1
		headerText := fmt.Sprintf("%s - %v عدد", headerName, int(projectDevice.Count))
		f.SetCellValue(deviceSheetText, getCellName(coli, rowi), headerText)
		f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli+6, rowi+3), ew.deviceStyle)
		f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli+6, rowi+3))

		// Write column headers
		rowi += 4
		titles := []struct {
			title   string
			colSpan int
		}{
			{"اقلام", 1},
			{"سایز", 1},
			{"جنس", 1},
			{"برند", 1},
			{"تعداد(مقدار)", 1},
			{"بهای واحد(ریال)", 1},
			{"بهای کل(ریال)", 1},
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

		// sort
		slices.SortFunc(device.DevicePartList, cmpDevicePart)

		for _, devicePart := range device.DevicePartList {
			coli = 1
			part := devicePart.Part

			// Part name
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), part.Name)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
			f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi))

			// Size
			coli += 1
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), part.Size)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)

			// Material
			coli += 1
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), part.Material)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)

			// Brand
			coli += 1
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), part.Brand)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
			f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi))

			// Count
			coli += 1
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), devicePart.Count)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
			f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi))

			// Price per unit
			coli += 1
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), part.Price)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
			f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi))

			// Total price for this part
			coli += 1
			partTotalPrice := devicePart.Price
			f.SetCellValue(deviceSheetText, getCellName(coli, rowi), partTotalPrice)
			f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi), ew.bodyStyle)
			f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi))

			deviceTotalPrice += partTotalPrice
			rowi++
		}

		// Write device total price
		coli = 1
		f.SetCellValue(deviceSheetText, getCellName(coli, rowi), "جمع کل(ریال)")
		f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli+5, rowi), ew.titleStyle)
		f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli+5, rowi))

		coli += 6
		f.SetCellValue(deviceSheetText, getCellName(coli, rowi), deviceTotalPrice)
		f.SetCellStyle(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi), ew.priceStyle)
		f.MergeCell(deviceSheetText, getCellName(coli, rowi), getCellName(coli, rowi))

		// Adjust column widths for better readability
		columns := []string{"B", "C", "D", "E"}
		for _, col := range columns {
			f.SetColWidth(deviceSheetText, col, col, 15)
		}

		columns = []string{"F", "G"}
		for _, col := range columns {
			f.SetColWidth(deviceSheetText, col, col, 25)
		}

		f.SetColWidth(deviceSheetText, "A", "A", 50)

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

	// Delete Sheet1
	f.DeleteSheet("Sheet1")

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
