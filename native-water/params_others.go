//go:build !linux && !darwin && !windows
// +build !linux,!darwin,!windows

package native-water

// PlatformSpeficParams
type PlatformSpecificParams struct {
}

func defaultPlatformSpecificParams() PlatformSpecificParams {
	return PlatformSpecificParams{}
}
