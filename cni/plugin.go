// Copyright 2017 Microsoft. All rights reserved.
// MIT License

package cni

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/Azure/azure-container-networking/common"
	"github.com/Azure/azure-container-networking/log"
	"github.com/Azure/azure-container-networking/platform"
	"github.com/Azure/azure-container-networking/store"

	cniInvoke "github.com/containernetworking/cni/pkg/invoke"
	cniSkel "github.com/containernetworking/cni/pkg/skel"
	cniTypes "github.com/containernetworking/cni/pkg/types"
	cniTypesCurr "github.com/containernetworking/cni/pkg/types/current"
	cniVers "github.com/containernetworking/cni/pkg/version"
)

// Plugin is the parent class for CNI plugins.
type Plugin struct {
	*common.Plugin
	version string
}

// NewPlugin creates a new CNI plugin.
func NewPlugin(name, version string) (*Plugin, error) {
	// Setup base plugin.
	plugin, err := common.NewPlugin(name, version)
	if err != nil {
		return nil, err
	}

	return &Plugin{
		Plugin:  plugin,
		version: version,
	}, nil
}

// Initialize initializes the plugin.
func (plugin *Plugin) Initialize(config *common.PluginConfig) error {
	// Initialize the base plugin.
	plugin.Plugin.Initialize(config)

	return nil
}

// Uninitialize uninitializes the plugin.
func (plugin *Plugin) Uninitialize() {
	plugin.Plugin.Uninitialize()
}

// Execute executes the CNI command.
func (plugin *Plugin) Execute(api PluginApi) (err error) {
	// Recover from panics and convert them to CNI errors.
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1<<12)
			len := runtime.Stack(buf, false)

			cniErr := &cniTypes.Error{
				Code:    ErrRuntime,
				Msg:     fmt.Sprintf("%v", r),
				Details: string(buf[:len]),
			}
			cniErr.Print()
			err = cniErr

			log.Printf("[cni] Recovered panic: %v %v\n", cniErr.Msg, cniErr.Details)
		}
	}()

	// Set supported CNI versions.
	pluginInfo := cniVers.PluginSupports(supportedVersions...)

	// Parse args and call the appropriate cmd handler.
	cniErr := cniSkel.PluginMainWithError(api.Add, api.Get, api.Delete, pluginInfo, plugin.version)
	if cniErr != nil {
		cniErr.Print()
		return cniErr
	}

	return nil
}

// DelegateAdd calls the given plugin's ADD command and returns the result.
func (plugin *Plugin) DelegateAdd(pluginName string, nwCfg *NetworkConfig) (*cniTypesCurr.Result, error) {
	var result *cniTypesCurr.Result
	var err error

	log.Printf("hellooooo from Printf")
	log.WriteToLog(log.LevelInfo, "hellooooo from WriteToLog")

	testLogUpdates()

	log.Printf("[cni] Calling plugin %v ADD nwCfg:%+v.", pluginName, nwCfg)
	defer func() { log.Printf("[cni] Plugin %v returned result:%+v, err:%v.", pluginName, result, err) }()

	os.Setenv(Cmd, CmdAdd)

	res, err := cniInvoke.DelegateAdd(context.TODO(), pluginName, nwCfg.Serialize(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to delegate: %v", err)
	}

	result, err = cniTypesCurr.NewResultFromResult(res)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert result: %v", err)
	}

	return result, nil
}

// DelegateDel calls the given plugin's DEL command and returns the result.
func (plugin *Plugin) DelegateDel(pluginName string, nwCfg *NetworkConfig) error {
	var err error

	log.Printf("[cni] Calling plugin %v DEL nwCfg:%+v.", pluginName, nwCfg)
	defer func() { log.Printf("[cni] Plugin %v returned err:%v.", pluginName, err) }()

	os.Setenv(Cmd, CmdDel)

	err = cniInvoke.DelegateDel(context.TODO(), pluginName, nwCfg.Serialize(), nil)
	if err != nil {
		return fmt.Errorf("Failed to delegate: %v", err)
	}

	return nil
}

// Error creates and logs a structured CNI error.
func (plugin *Plugin) Error(err error) *cniTypes.Error {
	var cniErr *cniTypes.Error
	var ok bool

	// Wrap error if necessary.
	if cniErr, ok = err.(*cniTypes.Error); !ok {
		cniErr = &cniTypes.Error{Code: 100, Msg: err.Error()}
	}

	log.Printf("[%v] %+v.", plugin.Name, cniErr.Error())

	return cniErr
}

// Errorf creates and logs a custom CNI error according to a format specifier.
func (plugin *Plugin) Errorf(format string, args ...interface{}) *cniTypes.Error {
	return plugin.Error(fmt.Errorf(format, args...))
}

// Initialize key-value store
func (plugin *Plugin) InitializeKeyValueStore(config *common.PluginConfig) error {
	// Create the key value store.
	if plugin.Store == nil {
		var err error
		plugin.Store, err = store.NewJsonFileStore(platform.CNIRuntimePath + plugin.Name + ".json")
		if err != nil {
			log.Printf("[cni] Failed to create store: %v.", err)
			return err
		}

		// Force unlock the json store if the lock file is left on the node after reboot
		if lockFileModTime, err := plugin.Store.GetLockFileModificationTime(); err == nil {
			rebootTime, err := platform.GetLastRebootTime()
			log.Printf("[cni] reboot time %v storeLockFile mod time %v", rebootTime, lockFileModTime)
			if err == nil && rebootTime.After(lockFileModTime) {
				log.Printf("[cni] Detected Reboot")

				if err := plugin.Store.Unlock(true); err != nil {
					log.Printf("[cni] Failed to force unlock store due to error %v", err)
				} else {
					log.Printf("[cni] Force unlocked the store successfully")
				}
			}
		}
	}

	// Acquire store lock.
	if err := plugin.Store.Lock(true); err != nil {
		log.Printf("[cni] Failed to lock store: %v.", err)
		return err
	}

	config.Store = plugin.Store

	return nil
}

// Uninitialize key-value store
func (plugin *Plugin) UninitializeKeyValueStore(force bool) error {
	if plugin.Store != nil {
		err := plugin.Store.Unlock(force)
		if err != nil {
			log.Printf("[cni] Failed to unlock store: %v.", err)
			return err
		}
	}
	plugin.Store = nil

	return nil
}

// check if safe to remove lockfile
func (plugin *Plugin) IsSafeToRemoveLock(processName string) (bool, error) {
	if plugin != nil && plugin.Store != nil {
		// check if get process command supported
		if cmdErr := platform.GetProcessSupport(); cmdErr != nil {
			log.Errorf("Get process cmd not supported. Error %v", cmdErr)
			return false, cmdErr
		}

		// Read pid from lockfile
		lockFileName := plugin.Store.GetLockFileName()
		content, err := ioutil.ReadFile(lockFileName)
		if err != nil {
			log.Errorf("Failed to read lock file :%v, ", err)
			return false, err
		}

		if len(content) <= 0 {
			log.Errorf("Num bytes read from lock file is 0")
			return false, fmt.Errorf("Num bytes read from lock file is 0")
		}

		log.Printf("Read from Lock file:%s", content)
		// Get the process name if running and
		// check if that matches with our expected process
		pName, err := platform.GetProcessNameByID(string(content))
		if err != nil {
			return true, nil
		}

		log.Printf("[CNI] Process name is %s", pName)

		if pName != processName {
			return true, nil
		}
	}

	log.Errorf("Plugin store is nil")
	return false, fmt.Errorf("plugin store nil")
}

// var alreadyTested = false
func testLogUpdates() {
	// if alreadyTested {
	// 	return
	// }
	// alreadyTested = true
	log.SetLevel(log.LevelInfo)
	log.WriteToLog(log.LevelInfo, "starting logs")

	// normal
	writeItAll()

	// test error-tolerance
	logWithFolder("thebest\\azure-container-networking/\\folder\\structure\\ever.go")
	logWithFolder("thebest/folder/structure/ever.go")
	logWithFolder("onefolder/withfile.go")
	logWithFolder("nogo/file/infolders")
	logWithFolder("onefolder")
	logWithFolder("onefolderwithslash/")
	logWithFolder("")
	logWithFolder("azure-container-networking")
	logWithFolder("superfolder/azure-container-networking")
	logWithFolder("superfolder/azure-container-networking/")

	logWithFolder("azure-container-networking/subfolder")
	logWithFolder("azure-container-networking/subfolder/")
	logWithFolder("azure-container-networking/subfolder/file.go")
	logWithFolder("azure-container-networking/subfolder/lowerfolder/file.go")
	logWithFolder("azure-container-networking/subfolder/lowerfolder")
	logWithFolder("azure-container-networking/subfolder/lowerfolder/")

	// failure from Caller()
	x, y, z, ok := runtime.Caller(0)
	ok = false
	log.SetComponentName(x, y, z, ok)
	writeItAll()
}

func logWithFolder(folderName string) {
	x, _, z, ok := runtime.Caller(0)
	log.SetComponentName(x, folderName, z, ok)
	log.WriteToLog(log.LevelInfo, "set folder name %s", folderName)
	writeItAll()
}

func writeItAll() {
	log.WriteToLog(log.LevelDebug, "round 1 shouldn't show up")
	log.WriteToLog(log.LevelWarning, "round 2: is %s", "good")
	log.WriteToLog(log.LevelError, "testing %d %d %s", 1, 2, "3")
	log.WriteToLog(log.LevelAlert, "round 4: testing alert")
}

// UNIT TESTING
// "thebest\\azure-container-networking/\\folder\\structure\\ever.go"
// "thebest/folder/structure/ever.go"
// "onefolder/withfile.go"
// "nogo/file/infolders"
// ""
// "onefolder"
// "onefolderwithslash/"

// ok := false
// check that all messages are in the log
