package instances

import (
	"fmt"

	"code.google.com/p/google-api-go-client/compute/v1"
)

// Tag for vertigo instances.
var VertigoTag = "vertigo"

// State for vertigo instance.
var StateTag = "vertigo-state"

// Possible states.
var StateOk = "ok"
var StateMigrating = "migrating"
var StateWarmingUp = "warming-up"

func ClearVertigoState(instance string, serv *compute.Service) error {
	inst, err := serv.Instances.Get(*gceProject, *gceZone, instance).Do()
	if err != nil {
		return err
	}

	// Remove the Vertigo tag.
	for i, data := range inst.Metadata.Items {
		if data.Key == VertigoTag {
			inst.Metadata.Items = append(inst.Metadata.Items[:i], inst.Metadata.Items[i+1:]...)
			break
		}
	}

	// Update the state.
	op, err := serv.Instances.SetMetadata(*gceProject, *gceZone, instance, inst.Metadata).Do()
	if err != nil {
		return nil
	}
	if op.Error != nil {
		return fmt.Errorf("failed to clear Vertigo state for %q: %s", instance, op.Error)
	}

	return nil
}

// Sets the state of the instance.
func SetInstanceState(state, instance string, serv *compute.Service) error {
	inst, err := serv.Instances.Get(*gceProject, *gceZone, instance).Do()
	if err != nil {
		return err
	}

	// Set the state.
	stateSet := false
	for _, data := range inst.Metadata.Items {
		if data.Key == StateTag {
			data.Value = state
			stateSet = true
			break
		}
	}

	// Set the state if not there.
	if !stateSet {
		inst.Metadata.Items = append(inst.Metadata.Items, &compute.MetadataItems{
			Key:   StateTag,
			Value: state,
		})
	}

	// Add Vertigo tag if not there.
	hasVertigo := getTag(VertigoTag, inst.Metadata) != ""
	if !hasVertigo {
		inst.Metadata.Items = append(inst.Metadata.Items, &compute.MetadataItems{
			Key:   VertigoTag,
			Value: VertigoTag,
		})
	}

	// Update the state.
	op, err := serv.Instances.SetMetadata(*gceProject, *gceZone, instance, inst.Metadata).Do()
	if err != nil {
		return nil
	}
	if op.Error != nil {
		return fmt.Errorf("failed to set metadata for %q: %s", instance, op.Error)
	}

	return nil
}

// Get the specified tag
func getTag(tag string, metadata *compute.Metadata) string {
	for _, data := range metadata.Items {
		if data.Key == tag {
			return data.Value
		}
	}

	return ""
}
