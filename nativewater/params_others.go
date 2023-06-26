//go:build !linux && !darwin && !windows
// +build !linux,!darwin,!windows

package nativewater

// PlatformSpeficParams
type PlatformSpecificParams struct {
}

func defaultPlatformSpecificParams() PlatformSpecificParams {
	return PlatformSpecificParams{}
}
