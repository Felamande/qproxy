package vmproxy

import (
	"fmt"
	"sync"

	"github.com/Felamande/qproxy/portfwd"
)

const globalConstRdpPort string = "3389"

var defaultVm *VmProxyManager
var initOk = false

func InitDefault(logChan chan error) *VmProxyManager {
	if initOk {
		return defaultVm
	}
	defaultVm = NewVmProxyManager(logChan)
	initOk = true
	return defaultVm
}

func Default() *VmProxyManager {
	if initOk {
		return defaultVm
	} else {
		return nil
	}
}

type VmInfo struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
}

type VmFwdInfo struct {
	vi        *VmInfo
	fwd       *portfwd.PortForwarder `json:"-"`
	FwdStatus bool
}

type VmProxyManager struct {
	vmFwdMap map[string]*VmFwdInfo

	updateLock sync.Mutex

	logChan chan error
}

func NewVmProxyManager(logChan chan error) *VmProxyManager {
	return &VmProxyManager{
		vmFwdMap: make(map[string]*VmFwdInfo),
		logChan:  logChan,
	}
}

func (v *VmProxyManager) UpdateVms(vms []VmInfo) error {
	v.updateLock.Lock()
	defer v.updateLock.Unlock()

	v.diff(vms,
		func(vi VmInfo) error { //delete
			err := v.vmFwdMap[vi.Name].fwd.Stop()
			delete(v.vmFwdMap, vi.Name)
			if err != nil {
				return err
			}
			return nil
		},
		func(vi VmInfo) error { //update
			old := v.vmFwdMap[vi.Name]
			if old.vi.IP == vi.IP {
				return nil
			}

			old.vi.IP = vi.IP
			if old.fwd != nil {
				err := old.fwd.Stop()
				old.fwd = portfwd.NewPortForwarder(v.logChan)
				old.FwdStatus = false

				return err
			}

			return nil
		},
		func(vi VmInfo) error { //add
			v.vmFwdMap[vi.Name] = &VmFwdInfo{
				vi: &VmInfo{
					IP:   vi.IP,
					Name: vi.Name,
				},
				fwd:       portfwd.NewPortForwarder(v.logChan),
				FwdStatus: false,
			}
			return nil
		})

	return nil

}

func (v *VmProxyManager) StartForward(name string) error {
	v.updateLock.Lock()
	defer v.updateLock.Unlock()

	fwdInfo, exist := v.vmFwdMap[name]
	if !exist {
		return fmt.Errorf("cannot find vm name:%v", name)
	}

	if fwdInfo.fwd == nil {
		fwdInfo.fwd = portfwd.NewPortForwarder(v.logChan)
	}

	if fwdInfo.fwd.IsRunning() {
		go func() {
			v.logChan <- fmt.Errorf("vmProxyManager use exist forward to %s:%s", fwdInfo.vi.IP, globalConstRdpPort)
		}()
		return nil
	}

	go func() {
		v.logChan <- fmt.Errorf("vmProxyManager open new forward to %s:%s", fwdInfo.vi.IP, globalConstRdpPort)
	}()

	if fwdInfo.vi.IP == "" {
		invalidIpErr := fmt.Errorf("the vm has no valid ip, maybe not started:name=%s", name)
		go func() {
			v.logChan <- invalidIpErr
		}()
		return invalidIpErr
	}

	return fwdInfo.fwd.Start(fwdInfo.vi.IP, globalConstRdpPort)

}

func (v *VmProxyManager) diff(newInfos []VmInfo, delete, update, add func(v VmInfo) error) (errs []error) {
	tmpNewInfoMap := make(map[string]VmInfo)

	for _, newInfo := range newInfos {
		tmpNewInfoMap[newInfo.Name] = newInfo
	}

	for name, info := range tmpNewInfoMap {
		if _, exist := v.vmFwdMap[name]; exist {
			err := update(info)
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			err := add(info)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	for name, _ := range v.vmFwdMap {
		if _, exist := tmpNewInfoMap[name]; !exist {
			delete(VmInfo{Name: name})
		}
	}

	return
}

func (v *VmProxyManager) GetPort(name string) int {
	fwdInfo, exist := v.vmFwdMap[name]
	if !exist || !fwdInfo.fwd.IsRunning() {
		return 0
	}
	port, _, _ := fwdInfo.fwd.Info()
	return port
}
