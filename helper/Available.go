package helper

import (
	"time"

	Init "github.com/shaikhzidhin/initiializer"
)

func FindAvailableRoomIDs(fromDate, toDate time.Time, roomIDs []uint) ([]uint, error) {
	var availableRoomIDs []uint
	rows, err := Init.DB.Raw(`
        SELECT DISTINCT rooms.id as room_id
        FROM rooms
        LEFT JOIN available_rooms ON rooms.id = available_rooms.room_id
        WHERE rooms.id IN (?) AND rooms.is_available = ? 
        AND (
            (available_rooms.room_id IS NULL) OR
            (
                NOT (? < ANY(available_rooms.check_in) AND ? < ANY(available_rooms.checkout)) AND 
                NOT (? > ANY(available_rooms.check_in) AND ? > ANY(available_rooms.checkout))
            )
        )
    `, roomIDs, true, toDate, fromDate, fromDate, toDate).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var roomID uint
		if err := rows.Scan(&roomID); err != nil {
			return nil, err
		}

		availableRoomIDs = append(availableRoomIDs, roomID)
	}

	return availableRoomIDs, nil
}

func GetRoomCountsByCategory() (map[string]int, error) {
	roomCounts := make(map[string]int)

	// Use a SQL query to count rooms in each category
	rows, err := Init.DB.Raw(`
        SELECT rc.name as category_name, COUNT(*) as room_count
        FROM rooms r
        JOIN room_categories rc ON r.room_category_id = rc.id
        GROUP BY rc.name
    `).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryName string
		var roomCount int
		if err := rows.Scan(&categoryName, &roomCount); err != nil {
			return nil, err
		}

		roomCounts[categoryName] = roomCount
	}

	return roomCounts, nil
}
