package geofence

// Geofence is a struct for efficient search whether a point is in polygon
type Geofence struct {
	vertices    []*Point
	tiles       map[float64]byte
	granularity int64
	minX        float64
	maxX        float64
	minY        float64
	maxY        float64
	tileWidth   float64
	tileHeight  float64
	minTileX    float64
	maxTileX    float64
	minTileY    float64
	maxTileY    float64
}

const (
	TILE_IN  = 0x01
	TILE_OUT = 0x02
	TILE_EITHER   = 0x03
)

const defaultGranularity = 20

// NewGeofence is the construct for Geofence, vertices: {{(1,2),(2,3)}, {(1,0)}}.
// 1st array contains polygon vertices. 2nd array contains holes.
func NewGeofence(points []*Point, args ...interface{}) *Geofence {
	geofence := &Geofence{}
	if len(args) > 0 {
		geofence.granularity = args[0].(int64)
	} else {
		geofence.granularity = defaultGranularity
	}
	geofence.vertices = points
	geofence.tiles = make(map[float64]byte)

	geofence.setInclusionTiles()
	return geofence
}

// Inside checks whether a given point is inside the geofence
func (geofence *Geofence) Inside(point *Point) bool {
	// Bbox check first
	if point.Lat() < geofence.minX || point.Lat() > geofence.maxX || point.Lng() < geofence.minY || point.Lng() > geofence.maxY {
		return false
	}

	tileHash := (project(point.Lng(), geofence.tileHeight)-geofence.minTileY)*float64(geofence.granularity) + (project(point.Lat(), geofence.tileWidth) - geofence.minTileX)
	intersects := geofence.tiles[tileHash]

	if intersects == TILE_IN {
		return true
	} else if intersects == TILE_EITHER {
		polygon := NewPolygon(geofence.vertices)
		inside := polygon.Contains(point)
		return inside
	} else {
		return false
	}
}

func (geofence *Geofence) setInclusionTiles() {
	xVertices := geofence.getXVertices()
	yVertices := geofence.getYVertices()

	geofence.minX = getMin(xVertices)
	geofence.minY = getMin(yVertices)
	geofence.maxX = getMax(xVertices)
	geofence.maxY = getMax(yVertices)

	xRange := geofence.maxX - geofence.minX
	yRange := geofence.maxY - geofence.minY
	geofence.tileWidth = xRange / float64(geofence.granularity)
	geofence.tileHeight = yRange / float64(geofence.granularity)

	geofence.minTileX = project(geofence.minX, geofence.tileWidth)
	geofence.minTileY = project(geofence.minY, geofence.tileHeight)
	geofence.maxTileX = project(geofence.maxX, geofence.tileWidth)
	geofence.maxTileY = project(geofence.maxY, geofence.tileHeight)

	geofence.setExclusionTiles(geofence.vertices, true)
}

func (geofence *Geofence) setExclusionTiles(vertices []*Point, inclusive bool) {
	var tileHash float64
	var bBoxPoly []*Point
	for tileX := geofence.minTileX; tileX <= geofence.maxTileX; tileX++ {
		for tileY := geofence.minTileY; tileY <= geofence.maxTileY; tileY++ {
			tileHash = (tileY-geofence.minTileY)*float64(geofence.granularity) + (tileX - geofence.minTileX)
			bBoxPoly = []*Point{NewPoint(tileX*geofence.tileWidth, tileY*geofence.tileHeight), NewPoint((tileX+1)*geofence.tileWidth, tileY*geofence.tileHeight), NewPoint((tileX+1)*geofence.tileWidth, (tileY+1)*geofence.tileHeight), NewPoint(tileX*geofence.tileWidth, (tileY+1)*geofence.tileHeight), NewPoint(tileX*geofence.tileWidth, tileY*geofence.tileHeight)}

			if haveIntersectingEdges(bBoxPoly, vertices) || hasPointInPolygon(vertices, bBoxPoly) {
				geofence.tiles[tileHash] = TILE_EITHER
			} else if hasPointInPolygon(bBoxPoly, vertices) {
				if inclusive {
					geofence.tiles[tileHash] = TILE_IN
				} else {
					geofence.tiles[tileHash] = TILE_OUT
				}
			} // else all points are outside the poly
		}
	}
}

func (geofence *Geofence) getXVertices() []float64 {
	xVertices := make([]float64, len(geofence.vertices))
	for i := 0; i < len(geofence.vertices); i++ {
		xVertices[i] = geofence.vertices[i].Lat()
	}
	return xVertices
}

func (geofence *Geofence) getYVertices() []float64 {
	yVertices := make([]float64, len(geofence.vertices))
	for i := 0; i < len(geofence.vertices); i++ {
		yVertices[i] = geofence.vertices[i].Lng()
	}
	return yVertices
}

func getMin(slice []float64) float64 {
	var min float64
	if len(slice) > 0 {
		min = slice[0]
	}
	for i := 1; i < len(slice); i++ {
		if slice[i] < min {
			min = slice[i]
		}
	}
	return min
}

func getMax(slice []float64) float64 {
	var max float64
	if len(slice) > 0 {
		max = slice[0]
	}
	for i := 1; i < len(slice); i++ {
		if slice[i] > max {
			max = slice[i]
		}
	}
	return max
}
