package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/device"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/project"
)

func Merge(path string) error {
	// Open both databases
	newDB, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer newDB.Close()

	// Merge the data
	if err := mergeAllTables(newDB); err != nil {
		return err
	}

	return nil
}

func mergeAllTables(newDB *sql.DB) error {
	// Iterate over all tables
	tables := []string{tableProjects, tableDevices, tableParts, tableDeviceParts, tableProjectDevices, tableExtraPrice}
	var err error
	var projects []*project.Project
	var devices []*device.Device
	var projectIdPairList []idPair
	var deviceIdPairList []idPair
	var partIdPairList []idPair

	for _, tb := range tables {
		switch tb {
		case tableProjects:
			if projects, projectIdPairList, err = mergeProjects(newDB); err != nil {
				return err
			}
		case tableDevices:
			if devices, deviceIdPairList, err = mergeDevices(newDB); err != nil {
				return err
			}
		case tableParts:
			if partIdPairList, err = mergeParts(newDB); err != nil {
				return err
			}
		case tableDeviceParts:
			if err = mergeDeviceParts(devices, deviceIdPairList, partIdPairList); err != nil {
				return err
			}
		case tableProjectDevices:
			if err = mergeProjectDevices(projects, projectIdPairList, deviceIdPairList); err != nil {
				return err
			}
		case tableExtraPrice:
			if err = mergeExtraPrice(projects, projectIdPairList); err != nil {
				return err
			}
		default:
			return err
		}
	}

	return nil
}

type idPair struct {
	srcID int
	newID int
}

// --------------------------
// ----- Merge projects -----
// --------------------------
func mergeProjects(newDB *sql.DB) ([]*project.Project, []idPair, error) {
	projectIdPairList := make([]idPair, 0)

	// get all projects form newDB
	newProjects, err := GetAllProjects(newDB)
	if err != nil {
		return nil, nil, err
	}

	// add new projects to newDB
	for _, p := range newProjects {
		id, err := AddProjectToDB(p)
		if err != nil {
			projectIdPairList = append(projectIdPairList, idPair{srcID: -1, newID: -1})
		}

		projectIdPairList = append(projectIdPairList, idPair{srcID: id, newID: p.Id})
	}

	return newProjects, projectIdPairList, nil
}

// -------------------------
// ----- Merge devices -----
// -------------------------
func mergeDevices(newDB *sql.DB) ([]*device.Device, []idPair, error) {
	deviceIdPairList := make([]idPair, 0)

	// get all devices form newDB
	newDevices, err := GetAllDevices(newDB)
	if err != nil {
		return nil, nil, err
	}

	// add new devices to newDB
	for _, d := range newDevices {
		id, err := AddDeviceToDB(d)
		if err != nil {
			deviceIdPairList = append(deviceIdPairList, idPair{srcID: -1, newID: -1})
		}

		deviceIdPairList = append(deviceIdPairList, idPair{srcID: id, newID: d.Id})
	}

	return newDevices, deviceIdPairList, nil
}

// -------------------------
// ------ Merge parts ------
// -------------------------
func mergeParts(newDB *sql.DB) ([]idPair, error) {
	partIdPairList := make([]idPair, 0)

	// get all parts form newDB
	newParts, err := GetAllParts(newDB)
	if err != nil {
		return nil, err
	}

	// add new parts to newDB
	for _, p := range newParts {
		id, err := AddPartToDB(p)
		if err != nil {
			partIdPairList = append(partIdPairList, idPair{srcID: -1, newID: -1})
		}

		partIdPairList = append(partIdPairList, idPair{srcID: id, newID: p.Id})
	}

	return partIdPairList, nil
}

// -----------------------------
// ----- Merge device parts ----
// -----------------------------
func mergeDeviceParts(devices []*device.Device, deviceIdPairList []idPair, partIdPairList []idPair) error {

	deviceIdMap := make(map[int]int, 0)
	for _, di := range deviceIdPairList {
		deviceIdMap[di.newID] = di.srcID
	}

	partIdMap := make(map[int]int, 0)
	for _, pi := range partIdPairList {
		partIdMap[pi.newID] = pi.srcID
	}

	for _, d := range devices {
		for _, dp := range d.DevicePartList {

			deviceId := deviceIdMap[dp.DeviceId]
			partId := partIdMap[dp.Part.Id]

			// add new device parts to srcDB
			if err := AddDevicePartToDB(deviceId, dp.Count, partId); err != nil {
				continue
			}
		}
	}

	return nil
}

// --------------------------------
// ----- Merge project devices ----
// --------------------------------
func mergeProjectDevices(projects []*project.Project, projectIdPairList []idPair, deviceIdPairList []idPair) error {

	projectIdMap := make(map[int]int, 0)
	for _, pi := range projectIdPairList {
		projectIdMap[pi.newID] = pi.srcID
	}

	deviceIdMap := make(map[int]int, 0)
	for _, di := range deviceIdPairList {
		deviceIdMap[di.newID] = di.srcID
	}

	for _, p := range projects {
		for _, pd := range p.ProjectDeviceList {

			projectId := projectIdMap[pd.ProjectId]
			deviceId := deviceIdMap[pd.Device.Id]

			// add new project devices to srcDB
			if err := AddProjectDeviceToDB(projectId, pd.Count, deviceId); err != nil {
				continue
			}
		}
	}

	return nil
}

// -----------------------------
// ----- Merge extra price -----
// -----------------------------
func mergeExtraPrice(projects []*project.Project, projectIdPairList []idPair) error {

	projectIdMap := make(map[int]int, 0)
	for _, pi := range projectIdPairList {
		projectIdMap[pi.newID] = pi.srcID
	}

	for _, p := range projects {
		for _, ep := range p.ExtraPriceList {

			projectId := projectIdMap[p.Id]

			// add new extra price to srcDB
			if err := AddExtraPriceToDB(projectId, ep.Name, ep.Price); err != nil {
				continue
			}
		}
	}

	return nil
}
