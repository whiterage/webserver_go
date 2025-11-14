package clock

import "time"

var moscowLocation = loadLocation()

func loadLocation() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.FixedZone("MSK", 3*60*60)
	}
	return loc
}

func Now() time.Time {
	return time.Now().In(moscowLocation)
}

func Location() *time.Location {
	return moscowLocation
}
