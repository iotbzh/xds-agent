// TODO add Doc
//
package main

import (
	"log"
	"os"
	"time"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/iotbzh/xds-agent/lib/agent"
	"github.com/iotbzh/xds-agent/lib/syncthing"
	"github.com/iotbzh/xds-agent/lib/xdsconfig"
	"github.com/iotbzh/xds-agent/lib/xdsserver"
)

const (
	appName        = "xds-agent"
	appDescription = "X(cross) Development System Agent is a web server that allows to remotely cross build applications."
	appCopyright   = "Apache-2.0"
	appUsage       = "X(cross) Development System Agent"
)

var appAuthors = []cli.Author{
	cli.Author{Name: "Sebastien Douheret", Email: "sebastien@iot.bzh"},
}

// AppVersion is the version of this application
var AppVersion = "?.?.?"

// AppSubVersion is the git tag id added to version string
// Should be set by compilation -ldflags "-X main.AppSubVersion=xxx"
var AppSubVersion = "unknown-dev"

// xdsAgent main routine
func xdsAgent(cliCtx *cli.Context) error {

	// Create Agent context
	ctx := agent.NewAgent(cliCtx)

	// Load config
	cfg, err := xdsconfig.Init(cliCtx, ctx.Log)
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	ctx.Config = &cfg

	// Start local instance of Syncthing and Syncthing-notify
	ctx.SThg = st.NewSyncThing(ctx.Config.FileConf.SThgConf, ctx.Log)

	ctx.Log.Infof("Starting Syncthing...")
	ctx.SThgCmd, err = ctx.SThg.Start()
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	ctx.Log.Infof("Syncthing started (PID %d)", ctx.SThgCmd.Process.Pid)

	ctx.Log.Infof("Starting Syncthing-inotify...")
	ctx.SThgInotCmd, err = ctx.SThg.StartInotify()
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	ctx.Log.Infof("Syncthing-inotify started (PID %d)", ctx.SThgInotCmd.Process.Pid)

	// Establish connection with local Syncthing (retry if connection fail)
	time.Sleep(3 * time.Second)
	retry := 10
	for retry > 0 {
		if err := ctx.SThg.Connect(); err == nil {
			break
		}
		ctx.Log.Infof("Establishing connection to Syncthing (retry %d/5)", retry)
		time.Sleep(time.Second)
		retry--
	}
	if ctx.SThg == nil {
		err = fmt.Errorf("ERROR: cannot connect to Syncthing (url: %s)", ctx.SThg.BaseURL)
		return cli.NewExitError(err, 2)
	}

	// Retrieve Syncthing config
	id, err := ctx.SThg.IDGet()
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	ctx.Log.Infof("Local Syncthing ID: %s", id)

	// Create and start Web Server
	ctx.WWWServer = xdsserver.NewServer(ctx.Config, ctx.Log)
	if err = ctx.WWWServer.Serve(); err != nil {
		log.Println(err)
		return cli.NewExitError(err, 3)
	}

	return cli.NewExitError("Program exited ", 4)
}

// main
func main() {

	// Create a new instance of the logger
	log := logrus.New()

	// Create a new App instance
	app := cli.NewApp()
	app.Name = appName
	app.Description = appDescription
	app.Usage = appUsage
	app.Version = AppVersion + " (" + AppSubVersion + ")"
	app.Authors = appAuthors
	app.Copyright = appCopyright
	app.Metadata = make(map[string]interface{})
	app.Metadata["version"] = AppVersion
	app.Metadata["git-tag"] = AppSubVersion
	app.Metadata["logger"] = log

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "JSON config file to use\n\t",
			EnvVar: "XDS_CONFIGFILE",
		},
		cli.StringFlag{
			Name:   "log, l",
			Value:  "error",
			Usage:  "logging level (supported levels: panic, fatal, error, warn, info, debug)\n\t",
			EnvVar: "LOG_LEVEL",
		},
	}

	// only one action
	app.Action = xdsAgent

	app.Run(os.Args)
}
