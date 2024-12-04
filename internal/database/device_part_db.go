package database

import (
	"database/sql"
	"fmt"
	"log"

	devicepart "github.com/poriamsz55/BoosterPump-webapp/internal/models/device_part"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/part"
)

func AddDevicePartToDB(dvcId int, count float64, prtId int) error {
	query := `INSERT INTO ` + tableDeviceParts + ` (` + columnDevicePartCount + `, ` +
		columnPartIDK + `, ` +
		columnDeviceIDFK + `) 
	          VALUES (?, ?, ?)`

	stmt, err := instance.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(count, prtId, dvcId)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	return nil
}

// Helper function to get device parts for a specific device
func GetDevicePartsByDeviceId(deviceID int) ([]*devicepart.DevicePart, error) {

	query := fmt.Sprintf(`
        SELECT dp.%s, dp.%s, p.%s, p.%s, p.%s, p.%s, p.%s, p.%s
        FROM %s dp
        JOIN %s p ON dp.%s = p.%s
        WHERE dp.%s = ?
    `,
		columnDevicePartID, columnDevicePartCount,
		columnPartID, columnPartName, columnPartSize, columnPartMaterial, columnPartBrand, columnPartPrice,
		tableDeviceParts,
		tableParts, columnPartIDK, columnPartID,
		columnDeviceIDFK)

	rows, err := instance.db.Query(query, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deviceParts []*devicepart.DevicePart
	for rows.Next() {
		var dp devicepart.DevicePart
		var part part.Part

		err := rows.Scan(
			&dp.Id,
			&dp.Count,
			&part.Id,
			&part.Name,
			&part.Size,
			&part.Material,
			&part.Brand,
			&part.Price,
		)
		if err != nil {
			return nil, err
		}

		dp.Part = &part
		dp.UpdatePrice()
		deviceParts = append(deviceParts, &dp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return deviceParts, nil
}

func GetDevicePartByIdFromDB(id int) (*devicepart.DevicePart, error) {

	query := fmt.Sprintf(`
        SELECT dp.%s, dp.%s, dp.%s, dp.%s, p.%s, p.%s, p.%s, p.%s, p.%s, p.%s
        FROM %s dp
        JOIN %s p ON dp.%s = p.%s
        WHERE dp.%s = ?
    `,
		columnDevicePartID, columnDevicePartCount, columnProjectIDFK, columnDeviceIDFK,
		columnPartID, columnPartName, columnPartSize, columnPartMaterial, columnPartBrand, columnPartPrice,
		tableDeviceParts,
		tableParts, columnPartIDK, columnPartID,
		columnDevicePartID)

	var dp devicepart.DevicePart
	var p part.Part
	err := instance.db.QueryRow(query, id).Scan(
		&dp.Id,
		&dp.Count,
		&dp.DeviceId,
		&dp.Part.Id,
		&p.Id,
		&p.Name,
		&p.Size,
		&p.Material,
		&p.Brand,
		&p.Price,
	)
	if err != nil {
		return nil, err
	}

	dp.Part = &p
	dp.UpdatePrice()
	return &dp, nil
}

func DeleteDevicePartFromDB(id int) error {

	query := fmt.Sprintf(`
        DELETE FROM %s 
        WHERE %s = ?
    `, tableDeviceParts, columnDevicePartID)

	result, err := instance.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func UpdateDevicePartInDB(id, projectId int, count float64, deviceId int) error {

	query := fmt.Sprintf(`
        UPDATE %s 
        SET %s = ?, %s = ?, %s = ?
        WHERE %s = ?
    `, tableDeviceParts,
		columnProjectIDFK, columnDevicePartCount, columnDeviceIDFK,
		columnDevicePartID)

	result, err := instance.db.Exec(query, projectId, count, deviceId, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
