package unit
import "time"

type Unit byte

const (
	Seconds Unit = iota
	Minutes
	Hours
	Days
	Weeks
	Months
	Years)

var unitLabels = []string{
	"Seconds",
	"Minutes",
	"Hours",
	"Days",
	"Weeks",
	"Months",
	"Years"}

var unitNanos = []int64{
	int64(time.Second),
	int64(time.Minute),
	int64(time.Hour),
	24 * int64(time.Hour),
	7 * 24 * int64(time.Hour),
	30 * 24 * int64(time.Hour),
	265 * 24 * int64(time.Hour)}

var unitIndex = []Unit {Seconds, Minutes, Hours, Days, Weeks, Months, Years}
func (u Unit) ConvertSeconds(seconds int64) int64 {
	nanos := seconds * int64(time.Second)
	return nanos / unitNanos[u]
}

func (u Unit) String() string {
	return unitLabels[u]
}

func FindUnit(t time.Duration) Unit {
	for i, nanos := range unitNanos {
		if int64(t) == nanos {
			return unitIndex[i]
		}
		if int64(t) < nanos {
			return unitIndex[i - 1]
		}
	}

	return unitIndex[len(unitNanos) - 1]
}