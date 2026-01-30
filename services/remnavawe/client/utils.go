package client

import "github.com/google/uuid"

func parseUuid(v string) uuid.UUID {
	value, _ := uuid.Parse(v)
	return value
}

var Plans = map[string]*RemnaPlan{
	"3:30": {
		DayLimit:    30,
		DeviceLimit: 3,
		Squad:       parseUuid("81360544-6e21-45be-aa48-c308074f0b0b"),
	},
	"5:30": {
		DayLimit:    30,
		DeviceLimit: 5,
		Squad:       parseUuid(""),
	},
	"3:90": {
		DayLimit:    90,
		DeviceLimit: 3,
		Squad:       parseUuid(""),
	},
	"5:90": {
		DayLimit:    90,
		DeviceLimit: 5,
		Squad:       parseUuid(""),
	},
	"3:180": {
		DayLimit:    180,
		DeviceLimit: 3,
		Squad:       parseUuid(""),
	},
	"5:180": {
		DayLimit:    180,
		DeviceLimit: 5,
		Squad:       parseUuid(""),
	},
	"3:365": {
		DayLimit:    365,
		DeviceLimit: 3,
		Squad:       parseUuid(""),
	},
	"5:365": {
		DayLimit:    365,
		DeviceLimit: 5,
		Squad:       parseUuid(""),
	},
}
