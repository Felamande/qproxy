package vmproxy

import (
	"fmt"

	"github.com/KnicKnic/go-powershell/pkg/powershell"
)

func GetVMAdapters() (m []VmInfo) {
	m = make([]VmInfo, 0)
	cmd := "$(get-vm|select name).name"
	rs := powershell.CreateRunspaceSimple()
	ret := rs.ExecScript(cmd, false, nil, "OS")
	ret.Close()

	for _, obj := range ret.Objects {
		name := obj.ToString()
		vmAdapterCmd := fmt.Sprintf(`$(get-vmnetworkadapter -vmname %s|select ipaddresses).ipaddresses[0]`, name)
		ret := rs.ExecScript(vmAdapterCmd, false, nil, "OS")
		ret.Close()

		if len(ret.Objects) >= 1 {
			ip := ret.Objects[0].ToString()
			ret.Objects[0].IsNull()
			if ret.Objects[0].IsNull() {
				ip = ""
			}
			m = append(m, VmInfo{IP: ip, Name: name})
		}
	}
	return
}
