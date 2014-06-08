package instances

import (
	"fmt"
)

// TODO(vmarmol): In real life, we'd probably discover this.
// Map of instance name to number of cores.
var instanceMapping = map[string]int{
	"vertigo-0": 1,
	"vertigo-1": 2,
	"vertigo-2": 4,
}

// Get an instance one size larger than the specified instance.
func GetLargerInstance(instanceName string) (string, error) {
	val, ok := instanceMapping[instanceName]
	if !ok {
		return "", fmt.Errorf("unknown instance %q", instanceName)
	}

	// Found one twice as big.
	val = val * 2
	for instance, size := range instanceMapping {
		if val == size {
			return instance, nil
		}
	}

	return "", fmt.Errorf("failed to find a larger instance than %q", instanceName)
}

// Get an instance one size smaller then the specified instance.
func GetSmallerInstance(instanceName string) (string, error) {
	val, ok := instanceMapping[instanceName]
	if !ok {
		return "", fmt.Errorf("unknown instance %q", instanceName)
	}

	// Found one half as big.
	val = val / 2
	for instance, size := range instanceMapping {
		if val == size {
			return instance, nil
		}
	}

	return "", fmt.Errorf("failed to find a smaller instance than %q", instanceName)
}
