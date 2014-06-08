package instances

import (
	"flag"

	"code.google.com/p/google-api-go-client/compute/v1"
)

var gceProject = flag.String("gce-project", "lmctfy-prod", "GCE project to use")
var gceZone = flag.String("gce-zone", "us-central1-a", "GCE zone running Vertigo")

type Instance struct {
	Name  string
	State string
}

// List Vertigo instances in this zone.
func GetVertigoInstances(serv *compute.Service) ([]*Instance, error) {
	// Get instances.
	instances, err := serv.Instances.List(*gceProject, *gceZone).Do()
	if err != nil {
		return nil, err
	}
	output := make([]*Instance, 0, len(instances.Items))
	for _, instance := range instances.Items {
		// Skip non-vertigo instances.
		if getTag(VertigoTag, instance.Metadata) == "" {
			continue
		}
		output = append(output, &Instance{
			Name:  instance.Name,
			State: getTag(StateTag, instance.Metadata),
		})
	}
	return output, nil
}
