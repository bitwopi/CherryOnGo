package client

import "github.com/google/uuid"

func parseUuid(v string) uuid.UUID {
	value, _ := uuid.Parse(v)
	return value
}

var Plans = map[string]*RemnaPlan{
	"test": {
		DayLimit:    1,
		DeviceLimit: 3,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"3:30": {
		DayLimit:    30,
		DeviceLimit: 3,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"4:30": {
		DayLimit:    30,
		DeviceLimit: 4,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"5:30": {
		DayLimit:    30,
		DeviceLimit: 5,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"6:30": {
		DayLimit:    30,
		DeviceLimit: 6,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"3:90": {
		DayLimit:    90,
		DeviceLimit: 3,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"4:90": {
		DayLimit:    90,
		DeviceLimit: 4,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"5:90": {
		DayLimit:    90,
		DeviceLimit: 5,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"6:90": {
		DayLimit:    90,
		DeviceLimit: 6,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"3:180": {
		DayLimit:    180,
		DeviceLimit: 3,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"4:180": {
		DayLimit:    180,
		DeviceLimit: 4,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"5:180": {
		DayLimit:    180,
		DeviceLimit: 5,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"6:180": {
		DayLimit:    180,
		DeviceLimit: 6,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"3:365": {
		DayLimit:    365,
		DeviceLimit: 3,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"4:365": {
		DayLimit:    365,
		DeviceLimit: 4,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"5:365": {
		DayLimit:    365,
		DeviceLimit: 5,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
	"6:365": {
		DayLimit:    365,
		DeviceLimit: 6,
		Squad:       parseUuid("f0bb8401-22ee-4b67-b256-d24cd64ee102"),
	},
}
