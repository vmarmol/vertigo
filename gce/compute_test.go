package gce

import "testing"

func TestNewMachine(t *testing.T) {
	mngr, err := NewGceManager()
	if err != nil {
		t.Error(err)
	}
	spec := &VirtualMachineSpec{}
	_, err = mngr.NewMachine(spec)
	if err != nil {
		t.Error(err)
	}
}
