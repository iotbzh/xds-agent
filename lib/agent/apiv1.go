/*
 * Copyright (C) 2017 "IoT.bzh"
 * Author Sebastien Douheret <sebastien@iot.bzh>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package agent

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iotbzh/xds-agent/lib/xdsconfig"
	"github.com/iotbzh/xds-server/lib/xsapiv1"
)

const apiBaseURL = "/api/v1"

// APIService .
type APIService struct {
	*Context
	apiRouter   *gin.RouterGroup
	serverIndex int
}

// NewAPIV1 creates a new instance of API service
func NewAPIV1(ctx *Context) *APIService {
	s := &APIService{
		Context:     ctx,
		apiRouter:   ctx.webServer.router.Group(apiBaseURL),
		serverIndex: 0,
	}

	s.apiRouter.GET("/version", s.getVersion)

	s.apiRouter.GET("/config", s.getConfig)
	s.apiRouter.POST("/config", s.setConfig)

	s.apiRouter.GET("/browse", s.browseFS)

	s.apiRouter.GET("/projects", s.getProjects)
	s.apiRouter.GET("/projects/:id", s.getProject)
	s.apiRouter.PUT("/projects/:id", s.updateProject)
	s.apiRouter.POST("/projects", s.addProject)
	s.apiRouter.POST("/projects/sync/:id", s.syncProject)
	s.apiRouter.DELETE("/projects/:id", s.delProject)

	s.apiRouter.POST("/exec", s.execCmd)
	s.apiRouter.POST("/exec/:id", s.execCmd)
	s.apiRouter.POST("/signal", s.execSignalCmd)

	s.apiRouter.GET("/events", s.eventsList)
	s.apiRouter.POST("/events/register", s.eventsRegister)
	s.apiRouter.POST("/events/unregister", s.eventsUnRegister)

	return s
}

// Stop Used to stop/close created services
func (s *APIService) Stop() {
	for _, svr := range s.xdsServers {
		svr.Close()
	}
}

// AddXdsServer Add a new XDS Server to the list of a server
func (s *APIService) AddXdsServer(cfg xdsconfig.XDSServerConf) (*XdsServer, error) {
	var svr *XdsServer
	var exist, tempoID bool
	tempoID = false

	// First check if not already exist and update it
	if svr, exist = s.xdsServers[cfg.ID]; exist {

		// Update: Found, so just update some settings
		svr.ConnRetry = cfg.ConnRetry

		tempoID = svr.IsTempoID()
		if svr.Connected && !svr.Disabled && svr.BaseURL == cfg.URL && tempoID {
			return svr, nil
		}

		// URL differ or not connected, so need to reconnect
		svr.BaseURL = cfg.URL

	} else {

		// Create a new server object
		if cfg.APIBaseURL == "" {
			cfg.APIBaseURL = apiBaseURL
		}
		if cfg.APIPartialURL == "" {
			cfg.APIPartialURL = "/servers/" + strconv.Itoa(s.serverIndex)
			s.serverIndex = s.serverIndex + 1
		}

		// Create a new XDS Server
		svr = NewXdsServer(s.Context, cfg)

		svr.SetLoggerOutput(s.Config.LogVerboseOut)

		// Passthrough routes (handle by XDS Server)
		grp := s.apiRouter.Group(svr.PartialURL)
		svr.SetAPIRouterGroup(grp)

		// Declare passthrough routes
		s.sdksPassthroughInit(svr)

		// Register callback on Connection
		svr.ConnectOn(func(server *XdsServer) error {

			// Add server to list
			s.xdsServers[server.ID] = svr

			// Register event forwarder
			if err := s.sdksEventsForwardInit(server); err != nil {
				s.Log.Errorf("XDS Server %v - sdk event forwarding error: %v", server.ID, err)
			}

			// Load projects
			if err := s.projects.Init(server); err != nil {
				s.Log.Errorf("XDS Server %v - project init error: %v", server.ID, err)
			}

			// Registered to all events
			if err := server.EventRegister(xsapiv1.EVTAll, ""); err != nil {
				s.Log.Errorf("XDS Server %v - register all events error: %v", server.ID, err)
			}

			return nil
		})
	}

	// Established connection
	err := svr.Connect()

	// Delete temporary ID with it has been replaced by right Server ID
	if tempoID && !svr.IsTempoID() {
		delete(s.xdsServers, cfg.ID)
	}

	return svr, err
}

// DelXdsServer Delete an XDS Server from the list of a server
func (s *APIService) DelXdsServer(id string) error {
	if _, exist := s.xdsServers[id]; !exist {
		return fmt.Errorf("Unknown Server ID %s", id)
	}
	// Don't really delete, just disable it
	s.xdsServers[id].Close()
	return nil
}
