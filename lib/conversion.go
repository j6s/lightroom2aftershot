package lib

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

type AttributeMapper = func(preset AfterShotPreset, lightroomName string, value string) AfterShotPreset

// Copies one attribute directly into another one
func copyValueDirectly(destination string) AttributeMapper {
	return func(preset AfterShotPreset, lightroomName string, value string) AfterShotPreset {
		preset.Attributes[destination] = value
		return preset
	}
}

// Copies the absolute value
func absInt(destination string) AttributeMapper {
	return func(preset AfterShotPreset, lightroomName string, value string) AfterShotPreset {
		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("[WARN] Could not convert %s to a float: %s", value, err)
		}

		preset.Attributes[destination] = fmt.Sprintf("%d", int(math.Abs(valueFloat)))
		return preset
	}
}

// Multiplies the value in the attribute with the given multiplier
// and writes the result into another attribute.
// Assumes that the input value is float-like
func applyMultiplier(destination string, multiplier float64) AttributeMapper {
	return func(preset AfterShotPreset, lightroomName string, value string) AfterShotPreset {
		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("[WARN] Could not convert %s to a float: %s", value, err)
		}

		preset.Attributes[destination] = fmt.Sprintf("%f", valueFloat*multiplier)
		return preset
	}
}

// Does nothing, print's a warning message
func todo() AttributeMapper {
	return func(preset AfterShotPreset, lightroomName string, value string) AfterShotPreset {
		if value != "" && value != "0" {
			log.Printf(
				"[WARN] Cannot lightroom configuration '%s' with value '%s' to aftershot",
				lightroomName,
				value,
			)
		}
		return preset
	}
}

// Does nothing
func ignore() AttributeMapper {
	return func(preset AfterShotPreset, lightroomName string, value string) AfterShotPreset {
		return preset
	}
}

func NewAftershotPresetFromLightroom(lightroom LightroomPreset) AfterShotPreset {

	emptyAttributeSet := map[string]string{
		"bopt:scont":                       "0",
		"bopt:highlightrecval":             "0",
		"bopt:fillamount":                  "0",
		"bopt:Equalizer_kb.kbs_redhue":     "0",
		"bopt:Equalizer_kb.kbs_orangehue":  "0",
		"bopt:Equalizer_kb.kbs_yellowhue":  "0",
		"bopt:Equalizer_kb.kbs_greenhue":   "0",
		"bopt:Equalizer_kb.kbs_cyanhue":    "0",
		"bopt:Equalizer_kb.kbs_bluehue":    "0",
		"bopt:Equalizer_kb.kbs_magentahue": "0",
		"bopt:Equalizer_kb.kbs_redsat":     "0",
		"bopt:Equalizer_kb.kbs_orangesat":  "0",
		"bopt:Equalizer_kb.kbs_yellowsat":  "0",
		"bopt:Equalizer_kb.kbs_greensat":   "0",
		"bopt:Equalizer_kb.kbs_cyansat":    "0",
		"bopt:Equalizer_kb.kbs_bluesat":    "0",
		"bopt:Equalizer_kb.kbs_magentasat": "0",
		"bopt:Equalizer_kb.kbs_redlum":     "0",
		"bopt:Equalizer_kb.kbs_orangelum":  "0",
		"bopt:Equalizer_kb.kbs_yellowlum":  "0",
		"bopt:Equalizer_kb.kbs_greenlum":   "0",
		"bopt:Equalizer_kb.kbs_cyanlum":    "0",
		"bopt:Equalizer_kb.kbs_bluelum":    "0",
		"bopt:Equalizer_kb.kbs_magentalum": "0",
		"bopt:Equalizer_kb.kbs_enabled":    "true",
	}

	// Attribute mappers are used to map lightroom attributes to aftershot attributes.
	// Keys correspond to attribute names in lightroom configuration, values correspond to
	// attribute mapper functions (see above)
	attributeMappers := map[string]AttributeMapper{
		"Contrast2012":   copyValueDirectly("bopt:scont"),
		"Highlights2012": absInt("bopt:highlightrecval"),

		// Lightroom and Aftershot use a vastly differing scale.
		"Shadows2012": applyMultiplier("bopt:fillamount", 0.01),

		"Whites2012":  todo(),
		"Blacks2012":  todo(),
		"Clarity2021": todo(),
		"Vibrance":    copyValueDirectly("bopt:vibe"),
		"Saturation":  copyValueDirectly("bopt:sat"),

		// Through trial and errror I discovered that 100 in lightroom is roughly equal to 70 in aftershot
		"HueAdjustmentRed":     applyMultiplier("bopt:Equalizer_kb.kbs_redhue", 0.7),
		"HueAdjustmentOrange":  applyMultiplier("bopt:Equalizer_kb.kbs_orangehue", 0.7),
		"HueAdjustmentYellow":  applyMultiplier("bopt:Equalizer_kb.kbs_yellowhue", 0.7),
		"HueAdjustmentGreen":   applyMultiplier("bopt:Equalizer_kb.kbs_greenhue", 0.7),
		"HueAdjustmentAqua":    applyMultiplier("bopt:Equalizer_kb.kbs_cyanhue", 0.7),
		"HueAdjustmentBlue":    applyMultiplier("bopt:Equalizer_kb.kbs_bluehue", 0.7),
		"HueAdjustmentMagenta": applyMultiplier("bopt:Equalizer_kb.kbs_magentahue", 0.7),

		// Saturation values seem to be 1:1
		"SaturationAdjustmentRed":     copyValueDirectly("bopt:Equalizer_kb.kbs_redsat"),
		"SaturationAdjustmentOrange":  copyValueDirectly("bopt:Equalizer_kb.kbs_orangesat"),
		"SaturationAdjustmentYellow":  copyValueDirectly("bopt:Equalizer_kb.kbs_yellowsat"),
		"SaturationAdjustmentGreen":   copyValueDirectly("bopt:Equalizer_kb.kbs_greensat"),
		"SaturationAdjustmentAqua":    copyValueDirectly("bopt:Equalizer_kb.kbs_cyansat"),
		"SaturationAdjustmentBlue":    copyValueDirectly("bopt:Equalizer_kb.kbs_bluesat"),
		"SaturationAdjustmentMagenta": copyValueDirectly("bopt:Equalizer_kb.kbs_magentasat"),

		// Luminance seems to be 1:1
		"LuminanceAdjustmentRed":     copyValueDirectly("bopt:Equalizer_kb.kbs_redlum"),
		"LuminanceAdjustmentOrange":  copyValueDirectly("bopt:Equalizer_kb.kbs_orangelum"),
		"LuminanceAdjustmentYellow":  copyValueDirectly("bopt:Equalizer_kb.kbs_yellowlum"),
		"LuminanceAdjustmentGreen":   copyValueDirectly("bopt:Equalizer_kb.kbs_greenlum"),
		"LuminanceAdjustmentAqua":    copyValueDirectly("bopt:Equalizer_kb.kbs_cyanlum"),
		"LuminanceAdjustmentBlue":    copyValueDirectly("bopt:Equalizer_kb.kbs_bluelum"),
		"LuminanceAdjustmentMagenta": copyValueDirectly("bopt:Equalizer_kb.kbs_magentalum"),

		// Sharpness: Ignore detailed configuration. Main sharpness configuration defaults to 50 in lightroom and 100 in aftershot
		"Sharpness":     applyMultiplier("bopt:newsharpen", 2),
		"SharpenRadius": ignore(),
		"SharpenDetail": ignore(),

		// Lightroom specific metadata
		"Version":                    ignore(),
		"UUID":                       ignore(),
		"PresetType":                 ignore(),
		"crs":                        ignore(),
		"ProcessVersion":             ignore(),
		"SupportsColor":              ignore(),
		"SupportsOutputReferred":     ignore(),
		"SupportsNormalDynamicRange": ignore(),
		"SupportsAmount":             ignore(),
		"SupportsHighDynamicRange":   ignore(),
		"SupportsMonochrome":         ignore(),
		"SupportsSceneReferred":      ignore(),
		"OverrideLookVignette":       ignore(),
		"HasSettings":                ignore(),
		"ToneCurveName2012":          ignore(),
		"CameraProfile":              ignore(),

		// Handled in pass
		"ConvertToGrayscale":            ignore(),
		"Texture":                       ignore(),
		"SplitToningBalance":            ignore(),
		"SplitToningShadowSaturation":   ignore(),
		"SplitToningShadowHue":          ignore(),
		"GrainAmount":                   ignore(),
		"GrainFrequency":                ignore(),
		"GrainSize":                     ignore(),
		"ColorNoiseReduction":           ignore(),
		"ColorNoiseReductionSmoothness": ignore(),
		"ColorNoiseReductionDetail":     ignore(),
		"Dehaze":                        ignore(),
		"ParametricShadows":             ignore(),
		"ParametricDarks":               ignore(),
		"ParametricLights":              ignore(),
		"ParametricHighlights":          ignore(),
		"ParametricShadowSplit":         ignore(),
		"ParametricMidtoneSplit":        ignore(),
		"ParametricHighlightSplit":      ignore(),
		"HueAdjustmentPurple":           ignore(),
		"SaturationAdjustmentPurple":    ignore(),
		"LuminanceAdjustmentPurple":     ignore(),

		// Ignored
		"ColorGradeMidtoneHue": ignore(),
		"ColorGradeBlending":   ignore(),
		"ColorGradeMidtoneSat": ignore(),
	}

	// Custom passes that are applied to the preset after the attribute mapping (see above)
	// has finished. This can be used to add more involved logic.
	postMappingPasses := []func(LightroomPreset, AfterShotPreset) AfterShotPreset{

		// ConvertToGrayscale = 0 saturation
		func(lightroom LightroomPreset, preset AfterShotPreset) AfterShotPreset {
			if lightroom.Attributes["ConvertToGrayscale"] == "True" {
				preset.Attributes["bopt:sat"] = "0"
			}
			return preset
		},

		// Texture to wavelet sharpen USM in clarity mode
		func(lightroom LightroomPreset, preset AfterShotPreset) AfterShotPreset {
			if lightroom.Attributes["Texture"] != "0" {
				log.Printf("[INFO] Texture is translated to usage of the wavelet sharpen plugin. Make sure you have that plugin installed")
				preset.Attributes["bopt:WaveletSharpen2.bSphWaveletUsmon"] = "true"
				preset.Attributes["bopt:WaveletSharpen2.bSphWaveletUsmClarity"] = "true"

				// Approximated value. TODO: Compare with lightroom rendering
				preset.Attributes["bopt:WaveletSharpen2.bSphWaveletUsmRadius"] = "10"

				preset.Attributes["bopt:WaveletSharpen2.bSphWaveletUsmAmount"] = lightroom.Attributes["Texture"]
			}
			return preset
		},

		// Dehaze to local contrast
		func(lightroom LightroomPreset, preset AfterShotPreset) AfterShotPreset {
			if lightroom.Attributes["Dehaze"] != "0" {
				preset.Attributes["bopt:lc_enabled"] = "true"
				preset.Attributes["bopt:lc_strength"] = lightroom.Attributes["Dehaze"]

			}
			return preset
		},

		// Non supported features in aftershot
		// Aftershot does not support split toning
		func(lightroom LightroomPreset, preset AfterShotPreset) AfterShotPreset {
			if lightroom.Attributes["SplitToningBalance"] != "+50" || lightroom.Attributes["SplitToningShadowSaturation"] == "0" {
				log.Printf("[WARN] This preset seems to use split toning. Split toning is not supported by aftershot and will be ignored.")
			}
			if lightroom.Attributes["GrainAmount"] != "+50" {
				log.Printf("[WARN] This preset seems to use grain. Grain is not supported by aftershot and will be ignored.")
			}
			if lightroom.Attributes["ColorNoiseReduction"] != "25" {
				log.Printf("[WARN] This preset seems to use color noise reduction. This is not supported in Aftershot and will be ignored.")
			}
			if lightroom.Attributes["ParametricShadowSplit"] != "25" || lightroom.Attributes["ParametricMidtoneSplit"] != "50" || lightroom.Attributes["ParametricHighlightSplit"] != "75" {
				log.Printf("[WARN] This preset seems to use parametric splits. This is not supported in Aftershot and will be ignored.")
			}

			log.Printf("[INFO] Lightroom has 7 Adjustable colors, aftershot has 6. Purple will be ignored if it has any settings.")

			return preset
		},
	}

	preset := AfterShotPreset{
		Attributes: emptyAttributeSet,
		ToneCurve:  newAfterShotCombinedToneCurveFromLightRoomToneCurve(lightroom.ToneCurve),
	}
	preset.Attributes = map[string]string{
		"bopt:Equalizer_kb.kbs_enabled": "true",
		"bopt:curveson":                 "true",
	}

	for key, value := range lightroom.Attributes {
		if value == "" {
			continue
		}

		mapper, mapperExists := attributeMappers[key]
		if !mapperExists {
			mapper = todo()
		}

		preset = mapper(preset, key, value)
	}

	for _, pass := range postMappingPasses {
		preset = pass(lightroom, preset)
	}

	return preset
}

func newAfterShotToneCurvePointFromLightroomToneCurvePoint(lightroom LightroomToneCurvePoint) AfterShotToneCurvePoint {
	multiplier := AFTERSHOT_CURVE_MAX / LIGHTROOM_CURVE_MAX

	return AfterShotToneCurvePoint{
		In:  lightroom.In * multiplier,
		Out: lightroom.Out * multiplier,
	}
}

func newAfterShotToneCurveChannelFromLightRoomToneCurveChannel(lightroom LightroomToneCurve) AfterShotToneCurveChannel {
	points := make([]AfterShotToneCurvePoint, len(lightroom.Points))
	for index, lightroomPoint := range lightroom.Points {
		points[index] = newAfterShotToneCurvePointFromLightroomToneCurvePoint(lightroomPoint)
	}

	return AfterShotToneCurveChannel{
		Points: points,
	}
}

func newAfterShotCombinedToneCurveFromLightRoomToneCurve(lightroom LightroomCombinedToneCurve) AfterShotCombinedToneCurve {
	return AfterShotCombinedToneCurve{
		Rgb:   newAfterShotToneCurveChannelFromLightRoomToneCurveChannel(lightroom.Rgb),
		Red:   newAfterShotToneCurveChannelFromLightRoomToneCurveChannel(lightroom.Red),
		Green: newAfterShotToneCurveChannelFromLightRoomToneCurveChannel(lightroom.Green),
		Blue:  newAfterShotToneCurveChannelFromLightRoomToneCurveChannel(lightroom.Blue),
	}
}
