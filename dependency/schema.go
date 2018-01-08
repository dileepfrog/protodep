package dependency

import (
	"strings"
)

type ProtoDep struct {
	ProtoOutdir  string               `toml:"proto_outdir"`
	Dependencies []ProtoDepDependency `toml:"dependencies"`
}

type ProtoDepDependency struct {
	Target   string `toml:"target"`
	Revision string `toml:"revision"`
	Branch   string `toml:"branch"`
	Tag      string `toml:"tag"`
}

func (d *ProtoDepDependency) Repository() string {
	tokens := strings.Split(d.Target, "/")
	if len(tokens) > 3 {
		return strings.Join(tokens[0:3], "/")
	} else {
		return d.Target
	}
}

func (d *ProtoDepDependency) Directory() string {
	r := d.Repository()

	if d.Target == r {
		return "."
	} else {
		return "." + strings.Replace(d.Target, r, "", 1)
	}
}
