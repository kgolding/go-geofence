package geofence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeofenceGroup(t *testing.T) {
	group := NewGeofenceGroup()

	wogf1 := NewGeofence([]*Point{NewPoint(10, 10), NewPoint(100, 10), NewPoint(100, 50), NewPoint(10, 50)})
	bogf1 := NewGeofence([]*Point{NewPoint(20, 20), NewPoint(40, 20), NewPoint(40, 30), NewPoint(20, 30)})
	wogf2 := NewGeofence([]*Point{NewPoint(200, 5), NewPoint(300, 5), NewPoint(300, 100), NewPoint(200, 100)})
	bogf2 := NewGeofence([]*Point{NewPoint(50, 22), NewPoint(60, 22), NewPoint(60, 26), NewPoint(50, 26)})

	group.Add(1, []*Geofence{wogf1, wogf2}, []*Geofence{bogf1, bogf2})
	group.Add(2, []*Geofence{wogf1}, nil)
	group.Add(3, []*Geofence{}, nil)
	group.Add(4, nil, []*Geofence{bogf2})

	// far out should only in 3,4
	expectedResult := make(map[int]bool)
	expectedResult[3] = true
	expectedResult[4] = true
	assert.Equal(t, group.GetValidKeys(NewPoint(1000, 1000)), expectedResult)

	// point in wogf1 & not in any bogf
	expectedResult[1] = true
	expectedResult[2] = true
	assert.Equal(t, group.GetValidKeys(NewPoint(15, 15)), expectedResult)

	// point in wogf1 & in bogf1
	expectedResult = make(map[int]bool)
	expectedResult[2] = true
	expectedResult[3] = true
	expectedResult[4] = true
	assert.Equal(t, group.GetValidKeys(NewPoint(25, 25)), expectedResult)

	// point in wogf2
	expectedResult = make(map[int]bool)
	expectedResult[1] = true
	expectedResult[3] = true
	expectedResult[4] = true
	assert.Equal(t, group.GetValidKeys(NewPoint(250, 10)), expectedResult)
}
