package helper

import (
	"time"

	Init "github.com/shaikhzidhin/initiializer"
)

func FindAvailableRoomIDs(fromDate, toDate time.Time, roomIDs []uint) ([]uint, error) {
	var availableRoomIDs []uint
	rows, err := Init.DB.Raw(`
        SELECT MIN(id) as room_id, room_category_id, COUNT(*) as room_count
        FROM rooms
        WHERE id IN (?) AND is_available = ? AND ? NOT IN (SELECT unnest(checkout) FROM available_rooms) AND ? NOT IN (SELECT unnest(check_in) FROM available_rooms)
        GROUP BY room_category_id
    `, roomIDs, true, fromDate, toDate).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var roomID uint
		var roomCategoryID uint
		var roomCount int
		if err := rows.Scan(&roomID, &roomCategoryID, &roomCount); err != nil {
			return nil, err
		}

		availableRoomIDs = append(availableRoomIDs, roomID)
	}

	return availableRoomIDs, nil
}

func GetRoomCountsByCategory() (map[uint]int, error) {
	roomCounts := make(map[uint]int)

	// Use a SQL query to count rooms in each category
	rows, err := Init.DB.Raw(`
        SELECT room_category_id, COUNT(*) as room_count
        FROM rooms
        GROUP BY room_category_id
    `).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var roomCategoryID uint
		var roomCount int
		if err := rows.Scan(&roomCategoryID, &roomCount); err != nil {
			return nil, err
		}

		roomCounts[roomCategoryID] = roomCount
	}

	return roomCounts, nil
}