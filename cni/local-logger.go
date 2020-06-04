// Copyright 2017 Microsoft. All rights reserved.
// MIT License

package cni

import "github.com/Azure/azure-container-networking/log"

var logger = log.NewLogger("azure-container-networking", log.LevelInfo, log.TargetStderr, "")
logger.SetComponentName(runtime.Caller(0))
