package helper

import (
	"time"

	"github.com/shaikhzidhin/initializer"
	Init "github.com/shaikhzidhin/initializer"
)

// FindAvailableRoomIDs retrieves room IDs that are available for a specific date range.
func FindAvailableRoomIDs(fromDate, toDate time.Time, roomIDs []uint) ([]uint, error) {
	var availableRoomIDs []uint
	rows, err := initializer.DB.Raw(`
        SELECT rooms.id as room_id
        FROM rooms
        LEFT JOIN available_rooms ON rooms.id = available_rooms.room_id
        WHERE rooms.id IN (?) AND rooms.is_available = ?
        `, roomIDs, true).Rows()
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

// GetRoomCountsByCategory helps to count the room catagory numbers
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
