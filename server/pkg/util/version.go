package util

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type VersionChange string

const (
	MajorVersion VersionChange = "major"
	MinorVersion VersionChange = "minor"
	PatchVersion VersionChange = "patch"
)

func NextVersion(currentVersion string, change VersionChange) (string, error) {
	major, minor, patch, err := ParseSemver(currentVersion)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse input version")
	}
	switch change {
	case MajorVersion:
		major += 1
	case MinorVersion:
		minor += 1
	case PatchVersion:
		patch += 1
	}
	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

func ParseSemver(ver string) (int, int, int, error) {
	parts := strings.Split(ver, ".")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("nor enough parts to semantic version: %s", ver)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, errors.Wrapf(err, "MAJOR version number %s is not an integer", parts[0])
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, errors.Wrapf(err, "MINOR version number %s is not an integer", parts[1])
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, errors.Wrapf(err, "PATCH version number %s is not an integer", parts[2])
	}
	return major, minor, patch, nil
}
