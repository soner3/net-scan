package hosts

import (
	"errors"
	"slices"
)

var (
	ErrExists    = errors.New("host already in the list")
	ErrNotExists = errors.New("host does not exist in the list")
)

type HostList struct {
	Hosts []string
}

func NewHostList() *HostList {
	return &HostList{}
}

func (hl *HostList) Add(host string) error {
	return nil
}

func (hl *HostList) Remove(host string) error {
	return nil
}

func (hl *HostList) Load(host string) error {
	return nil
}

func (hl *HostList) Save(host string) error {
	return nil
}

func (hl *HostList) search(host string) (bool, int) {
	slices.Sort(hl.Hosts)
	i := slices.Index(hl.Hosts, host)
	if i != -1 && hl.Hosts[i] == host {
		return true, i
	}
	return false, -1
}
