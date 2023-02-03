package ip2region

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/lua"
	"reflect"
)

var xEnv vela.Environment

var typeof = reflect.TypeOf((*region)(nil)).String()

func newLuaIP2Region(L *lua.LState) int {
	cfg := newConfig(L)

	proc := L.NewVelaData(cfg.name, typeof)
	if proc.IsNil() {
		proc.Set(newRegion(cfg))
	} else {
		old := proc.Data.(*region)
		old.cfg = cfg
	}

	L.Push(proc)
	return 1
}

func searchL(L *lua.LState) int {
	ip := L.CheckString(1)
	v, err := xEnv.Region(ip)
	if err != nil {
		L.Push(lua.S2L("0|0|0|未知IP|未知IP"))
	} else {
		L.Push(lua.B2L(v.Byte()))
	}
	return 1
}

func newIP2RegionByLoad(L *lua.LState) int {
	dbname := L.CheckString(1)
	info, err := xEnv.Third(dbname)
	if err != nil {
		L.RaiseError("%s ip database load fail %v", dbname, err)
		return 0
	}

	cfg := &config{
		name:   "ip2region." + dbname,
		method: "index",
		xdb:    info.File(),
	}

	proc := L.NewVelaData(cfg.name, typeof)
	if proc.IsNil() {
		proc.Set(newRegion(cfg))
	} else {
		old := proc.Data.(*region)
		old.cfg = cfg
	}

	L.Push(proc)
	return 1
}

func WithEnv(env vela.Environment) {
	xEnv = env

	kv := lua.NewUserKV()
	kv.Set("new", lua.NewFunction(newLuaIP2Region))
	kv.Set("load", lua.NewFunction(newIP2RegionByLoad))
	xEnv.Set("region", lua.NewExport("vela.region.export", lua.WithTable(kv), lua.WithFunc(searchL)))
}
