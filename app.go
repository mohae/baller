// app.go contains all of your application specific settings. Most changes
// needed for a new application should be here, including environment variable
// names, default settings, etc.
package main

import (
	"fmt"
	"os"

	log "github.com/cihub/seelog"
	"github.com/mohae/cli"
	"github.com/mohae/contour"
	"github.com/mohae/quine/cmd")

// Name is the name of the application
var Name string = "quine"

// The git commit that was compiled. This will be filled in by the compiler
var GitCommit string

// The main version number that is being run at the moment.
const Version = "0.0.10"

// A pre-release marker for the version. If this is "" (empty string)
// then it means that it is a final release. Otherwise, this is a pre-release
// such as "dev" (in development), "beta", "rc1", etc.
const VersionPrerelease = "dev"

// AppCode is the code for the application. This is used to prefix the
// environment variable. It can be left empty.
var AppCode string

// ConfigFilename is the configuration file for the application.
var ConfigFilename string = "config.json"

// LogConfigFilename is the name for the logging configuration file.
var LogConfigFilename string = "seelog.xml"

// Logging: whether or not application logging is enabled by default.
// Initialize to true if it should automatically be enabled.
var Logging bool

// Environment variables
var (
	EnvConfigFilename string = "configfilename"
	EnvLogConfigFilename string = "logconfigfilename"
	EnvLogging string = "logging"
)

// Config a pointer to the AppConfig. The AppConfig can either be updated by
// calling the contour function or the Config's method, both of which will be
// the same other than being a function or method. 
//
// If you want a different Config object to use for your configuration, call 
// contour.NewConfig() instead. This will return a new Config object. You will
// need to use its methods to work with it, calling contour's function won't 
// apply to it.
var Config *contour.Config = contour.GetAppConfig()

// set-up the application defaults and let contour know about the app.
// Setting also saves them to their relative environment variable.
func init() {
	// Calling Register* saves the configuration setting information
	// without writing it to its respective environment variable. This
	// allows any already set environment variables to override non-core
	// vars.
	//
	// Only settings that have been initialized are recognized by contour.

	// The config filename is handled differently, calling this function
	// also sets the ConfigFile format automatically,based on the
	// extension, if it can be determined. If it cannot, the extension is
	// left blank and must be set. 
	contour.RegisterConfigFilename(EnvConfigFilename, ConfigFilename)

	//// Alternative way, manually setting the values
	//contour.RegisterString(EnvConfigFilename, ConfigFilename) 
	//contour.RegisterString(EnvConfigExt, "json") 

	// Core settings are only settable by the application, and once set are
	// immutable
	contour.RegisterCoreString("appname", Name) 

	// Immutable settings are only settable once. If its value isn't set
	// during registration, it can be set at a later time. Once it is set,
	// immutable values cannot be changed. Because of this, and the fact
	// that initialization causes bools to be set, bools cannot be made
	// immutable.

	// This is set in the config file.
	contour.RegisterImmutableString(EnvLogConfigFilename, "")
	
	// Set*Flag allows you to add settings that are also exposed as
	// command-line flags. Default implicit values to settings:
	//	IsFlag = true
	//	IsIdempotent = false
	//	IsCore = false
	// The shortcode, 2nd parameter, can be left as an empty string, ""
	// if this flag doesn't support a shortcode.
	contour.RegisterBoolFlag(EnvLogging, Logging, "") 

	// AddSettingAlias sets an alias for the setting.
	contour.AddSettingAlias(EnvLogging, "logenabled")

	InitApp()

	// Now that the configuration in
}

// InitApp is the best place to add custom defaults for your application,
func InitApp() {
	contour.RegisterBoolFlag("lower", false, "")
}

// InitConfig initialized the application's configuration. When the config is
// has been initialized, the preset-enivronment variables, application 
// defaults, and your application's configuration file, should it have one,
// will all be merged according to the setting properties.
//
// After this, only overrides can occur via command flags.
func InitConfig() error {
	// Set config:
	//    Checks environment variables for settings, follows update rules.
	//    Retrieves config file and applies those settings, if and where
	//      applicable.
	//    Writes the resulting configuration settings to their env vars.
	//  After set config, only command flags can override the settings.
	//  If this is an interactive application, preference changes would
	//    also override certain settings. It may necessitate an additional
	//    flag or two.
	return contour.SetConfig()
}


// SetAppLogging sets the logger for package loggers and allow for custom-
// ization of the applications logging. This is where app specific code for
// setting up the application's logging should be.
// 
// SetAppLogging assumes that logging is enabled if it has been called as its
// caller should be SetLogging(). If you are going to call this from elsewhere,
// first make sure that logging is enabled.
//
// This uses seelog.
func SetAppLogging() error {
	contour.UseLogger(logger)
	cmd.UseLogger(logger)
	return nil
}

// appMain, is the actual main for the application. This keeps all changes
// needed for a new application to one file in the main application directory.
// In addition to this, only commands/ needs to be modified, adding the app's
// commands and any handler codes for those commands, like the 'cmd' package.
//
// No logging is done until the flags are processed, since the flags could
// enable/disable output, alter it, or alter its output locations. Everything
// must go to stdout until then.
func appMain() int {
	defer cmd.FlushLog()
	defer contour.FlushLog()
	defer log.Flush()

	// Initialize the applications's defaults
	err := InitConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error configuring %s: %s\n", Name, err.Error())
		return 1
	}
	
	// Set the Logging
	err = SetLogging()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up logging for %s: %s\n", Name, err.Error())
		return 1
	}

	// Get the command line args.
	args := os.Args[1:]

	// Setup the args, Commands, and Help info.
	cli := &cli.CLI{
		Name: Name,
		Version: Version,
		Commands: Commands,
		Args: args,
		HelpFunc: cli.BasicHelpFunc(),
	}

	// Run the passed command, recieve back a message and error object.
	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	// Return the exitcode.
	return exitCode
}
