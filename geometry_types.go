package rasterizer

type GeometryType string

const (
	GeometryTypeGeometry          GeometryType = "Geometry"
	GeometryTypeAxis              GeometryType = "Axis"
	GeometryTypeFilterBeam        GeometryType = "FilterBeam"
	GeometryTypeFilterColor       GeometryType = "FilterColor"
	GeometryTypeFilterGobo        GeometryType = "FilterGobo"
	GeometryTypeFilterShaper      GeometryType = "FilterShaper"
	GeometryTypeBeam              GeometryType = "Beam"
	GeometryTypeMediaServerLayer  GeometryType = "MediaServerLayer"
	GeometryTypeMediaServerCamera GeometryType = "MediaServerCamera"
	GeometryTypeMediaServerMaster GeometryType = "MediaServerMaster"
	GeometryTypeDisplay           GeometryType = "Display"
	GeometryTypeGeometryReference GeometryType = "GeometryReference"
	GeometryTypeLaser             GeometryType = "Laser"
	GeometryTypeWiringObject      GeometryType = "WiringObject"
	GeometryTypeInventory         GeometryType = "Inventory"
	GeometryTypeStructure         GeometryType = "Structure"
	GeometryTypeSupport           GeometryType = "Support"
	GeometryTypeMagnet            GeometryType = "Magnet"
)

var geometryTypes []GeometryType = []GeometryType{
	GeometryTypeGeometry,
	GeometryTypeAxis,
	GeometryTypeFilterBeam,
	GeometryTypeFilterColor,
	GeometryTypeFilterGobo,
	GeometryTypeFilterShaper,
	GeometryTypeBeam,
	GeometryTypeMediaServerLayer,
	GeometryTypeMediaServerCamera,
	GeometryTypeMediaServerMaster,
	GeometryTypeDisplay,
	GeometryTypeGeometryReference,
	GeometryTypeLaser,
	GeometryTypeWiringObject,
	GeometryTypeInventory,
	GeometryTypeStructure,
	GeometryTypeSupport,
	GeometryTypeMagnet,
}
