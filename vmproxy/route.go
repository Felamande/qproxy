package vmproxy

import (
	"gitea.com/lunny/tango"
	"gitea.com/tango/binding"
)

type VmProxyRoute struct {
	tango.JSON
	tango.Ctx
	binding.Binder
}

func NewVmProxyRoute(logChan chan error) *VmProxyRoute {
	return &VmProxyRoute{}
}

func (v *VmProxyRoute) Get() interface{} {
	vms := GetVMAdapters()
	Default().UpdateVms(vms)
	v.Logger.Infof("get adapters:%v", vms)
	return map[string]interface{}{
		"error": nil,
		"vms":   vms,
	}
}

func (v *VmProxyRoute) ReqStart() interface{} {
	type NameStruct struct {
		Name string
	}
	var data NameStruct
	errBind := v.Bind(&data)
	if errBind != nil {
		return map[string]interface{}{
			"name":  data.Name,
			"error": errBind,
		}
	}

	err := Default().StartForward(data.Name)
	if err != nil {
		return map[string]interface{}{
			"name":  data.Name,
			"error": err,
		}
	}

	port := Default().GetPort(data.Name)

	return map[string]interface{}{
		"name":  data.Name,
		"port":  port,
		"error": nil,
	}
}
