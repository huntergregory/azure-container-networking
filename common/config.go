// Copyright 2017 Microsoft. All rights reserved.
// MIT License

package common

// Command line options.
const (
	// Operating environment.
	OptEnvironment             = "environment"
	OptEnvironmentAlias        = "e"
	OptEnvironmentAzure        = "azure"
	OptEnvironmentMAS          = "mas"
	OptEnvironmentFileIpam     = "fileIpam"
	OptEnvironmentIPv6NodeIpam = "ipv6NodeIpam"

	// API server URL.
	OptAPIServerURL      = "api-url"
	OptAPIServerURLAlias = "u"
	OptCnsURL            = "cns-url"
	OptCnsURLAlias       = "c"

	// Logging level.
	OptLogLevel      = "log-level"
	OptLogLevelAlias = "l"
	OptLogLevelInfo  = "info"
	OptLogLevelDebug = "debug"

	// Logging target.
	OptLogTarget       = "log-target"
	OptLogTargetAlias  = "t"
	OptLogTargetSyslog = "syslog"
	OptLogTargetStderr = "stderr"
	OptLogTargetFile   = "logfile"
	OptLogStdout       = "stdout"
	OptLogMultiWrite   = "stdoutfile"

	// Logging location
	OptLogLocation      = "log-location"
	OptLogLocationAlias = "o"

	// IPAM query URL.
	OptIpamQueryUrl      = "ipam-query-url"
	OptIpamQueryUrlAlias = "q"

	// IPAM query interval.
	OptIpamQueryInterval      = "ipam-query-interval"
	OptIpamQueryIntervalAlias = "i"

	// Start CNM
	OptStartAzureCNM      = "start-azure-cnm"
	OptStartAzureCNMAlias = "startcnm"

	// Interval to send reports to host
	OptReportToHostInterval      = "report-interval"
	OptReportToHostIntervalAlias = "hostinterval"

	// Periodic Interval Time
	OptIntervalTime        = "interval"
	OptIntervalTimeAlias   = "it"
	OptDefaultIntervalTime = "defaultinterval"

	// Version.
	OptVersion      = "version"
	OptVersionAlias = "v"

	// Help.
	OptHelp      = "help"
	OptHelpAlias = "h"

	// CNI binary location
	OptNetPluginPath      = "net-plugin-path"
	OptNetPluginPathAlias = "np"

	// CNI binary location
	OptNetPluginConfigFile      = "net-plugin-config-file"
	OptNetPluginConfigFileAlias = "npconfig"

	// Telemetry config Location
	OptTelemetryConfigDir      = "telemetry-config-file"
	OptTelemetryConfigDirAlias = "d"

	// Create ext Hns network
	OptCreateDefaultExtNetworkType      = "create-defaultextnetwork-type"
	OptCreateDefaultExtNetworkTypeAlias = "defaultextnetworktype"

	// Disable Telemetry
	OptTelemetry      = "telemetry"
	OptTelemetryAlias = "dt"

	// HTTP connection timeout
	OptHttpConnectionTimeout      = "http-connection-timeout"
	OptHttpConnectionTimeoutAlias = "httpcontimeout"

	// HTTP response header timeout
	OptHttpResponseHeaderTimeout      = "http-response-header-timeout"
	OptHttpResponseHeaderTimeoutAlias = "httprespheadertimeout"

	// Store file location
	OptStoreFileLocation      = "store-file-path"
	OptStoreFileLocationAlias = "storefilepath"
)
