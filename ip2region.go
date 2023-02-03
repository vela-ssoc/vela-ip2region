package ip2region

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-ip2region/xdb"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
)

type region struct {
	lua.SuperVelaData
	cfg    *config
	search *xdb.Searcher
}

func (r *region) Name() string {
	return r.cfg.name
}

func (r *region) Type() string {
	return typeof
}

func (r *region) Start() error {
	switch r.cfg.method {
	case "file":
		return r.byFile()
	case "index":
		return r.byIndex()
	case "memory":
		return r.byMemory()
	default:
		return fmt.Errorf("not found invalid method")
	}
}

func (r *region) Close() error {
	if r.search == nil {
		return nil
	}

	old := r.search
	r.search = nil

	old.Close()
	return nil
}

func (r *region) byIndex() error {
	index, err := xdb.LoadVectorIndexFromFile(r.cfg.xdb)
	if err != nil {
		return err
	}

	search, err := xdb.NewWithVectorIndex(r.cfg.xdb, index)
	if err != nil {
		return err
	}

	r.search = search
	return nil
}

func (r *region) byMemory() error {
	buff, err := xdb.LoadContentFromFile(r.cfg.xdb)
	if err != nil {
		return err
	}

	search, err := xdb.NewWithBuffer(buff)
	if err != nil {
		return err
	}

	r.search = search
	return nil
}

func (r *region) byFile() error {
	search, err := xdb.NewWithFileOnly(r.cfg.xdb)
	if err != nil {
		return err
	}
	r.search = search
	return nil
}

func (r *region) Search(ip string) (*vela.IPv4Info, error) {
	if r.search == nil {
		return nil, fmt.Errorf("not found search")
	}

	raw, err := r.search.SearchByStr(ip)
	if err != nil {
		xEnv.Debugf("%s %s ip region got fail %v", r.Name(), ip, err)
		return nil, err
	}

	return vela.NewIPv4Info(0, auxlib.S2B(raw)), nil
}

func newRegion(cfg *config) *region {
	return &region{cfg: cfg}
}
