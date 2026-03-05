package rasterizer

import (
	"fmt"
	"image/color"

	GDTFTypes "github.com/Patch2PDF/GDTF-Parser/pkg/types"
)

type colorMap map[ModelType]map[GeometryType]color.NRGBA
type OverrideColorMap map[ModelType]map[GeometryType]*color.NRGBA

var default_colors = map[GeometryType]color.NRGBA{
	GeometryTypeGeometry:          {25, 25, 25, 255},
	GeometryTypeAxis:              {25, 25, 25, 255},
	GeometryTypeFilterBeam:        {0, 0, 0, 255},
	GeometryTypeFilterColor:       {0, 0, 0, 255},
	GeometryTypeFilterGobo:        {0, 0, 0, 255},
	GeometryTypeFilterShaper:      {0, 0, 0, 255},
	GeometryTypeBeam:              {200, 200, 200, 255},
	GeometryTypeMediaServerLayer:  {0, 0, 0, 255},
	GeometryTypeMediaServerCamera: {0, 0, 0, 255},
	GeometryTypeMediaServerMaster: {0, 0, 0, 255},
	GeometryTypeDisplay:           {0, 0, 0, 255},
	GeometryTypeGeometryReference: {0, 0, 0, 255},
	GeometryTypeLaser:             {0, 0, 0, 255},
	GeometryTypeWiringObject:      {0, 0, 0, 255},
	GeometryTypeInventory:         {0, 0, 0, 255},
	GeometryTypeStructure:         {0, 0, 0, 255},
	GeometryTypeSupport:           {0, 0, 0, 255},
	GeometryTypeMagnet:            {0, 0, 0, 255},
}

var fromGDTFGeometryMap = map[GDTFTypes.GeometryType]GeometryType{
	GDTFTypes.GeometryTypeAxis:              GeometryTypeAxis,
	GDTFTypes.GeometryTypeBeam:              GeometryTypeBeam,
	GDTFTypes.GeometryTypeDisplay:           GeometryTypeDisplay,
	GDTFTypes.GeometryTypeFilterBeam:        GeometryTypeFilterBeam,
	GDTFTypes.GeometryTypeFilterColor:       GeometryTypeFilterColor,
	GDTFTypes.GeometryTypeFilterGobo:        GeometryTypeFilterGobo,
	GDTFTypes.GeometryTypeFilterShaper:      GeometryTypeFilterShaper,
	GDTFTypes.GeometryTypeGeometry:          GeometryTypeGeometry,
	GDTFTypes.GeometryTypeGeometryReference: GeometryTypeGeometryReference,
	GDTFTypes.GeometryTypeInventory:         GeometryTypeInventory,
	GDTFTypes.GeometryTypeLaser:             GeometryTypeLaser,
	GDTFTypes.GeometryTypeMagnet:            GeometryTypeMagnet,
	GDTFTypes.GeometryTypeMediaServerCamera: GeometryTypeMediaServerCamera,
	GDTFTypes.GeometryTypeMediaServerLayer:  GeometryTypeMediaServerLayer,
	GDTFTypes.GeometryTypeMediaServerMaster: GeometryTypeMediaServerMaster,
	GDTFTypes.GeometryTypeStructure:         GeometryTypeStructure,
	GDTFTypes.GeometryTypeSupport:           GeometryTypeSupport,
	GDTFTypes.GeometryTypeWiringObject:      GeometryTypeWiringObject,
}

// TODO: add "global" override

func getOverrideColors(overrides OverrideColorMap) colorMap {
	result := make(colorMap)
	for i := range modelTypes {
		result[ModelType(i)] = make(map[GeometryType]color.NRGBA, len(default_colors))
		for _, geometryType := range geometryTypes {
			if overrides[ModelType(i)][geometryType] != nil {
				result[ModelType(i)][geometryType] = *overrides[ModelType(i)][geometryType]
			} else {
				result[ModelType(i)][geometryType] = default_colors[geometryType]
			}
		}
	}
	return result
}

func (colors colorMap) getColor(modelType ModelType, geometryType GDTFTypes.GeometryType) (color.NRGBA, error) {
	nonGDTFGeometryType, ok := fromGDTFGeometryMap[geometryType]
	if !ok {
		return color.NRGBA{}, fmt.Errorf("Invalid geometry type / could not find mapping for %d", geometryType)
	}
	modelColors, ok := colors[modelType]
	if !ok {
		return color.NRGBA{}, fmt.Errorf("Invalid model type %d", modelType)
	}
	return modelColors[nonGDTFGeometryType], nil
}
