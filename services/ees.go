/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package services

import (
	"fmt"
	"sync"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/ees"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/servmanager"
	"github.com/cgrates/cgrates/utils"
	"github.com/cgrates/rpcclient"
)

// NewEventExporterService constructs EventExporterService
func NewEventExporterService(cfg *config.CGRConfig, filterSChan chan *engine.FilterS,
	connMgr *engine.ConnManager, server *utils.Server, exitChan chan bool,
	intConnChan chan rpcclient.ClientConnector) servmanager.Service {
	return &EventExporterService{
		cfg:         cfg,
		filterSChan: filterSChan,
		connMgr:     connMgr,
		server:      server,
		exitChan:    exitChan,
		intConnChan: intConnChan,
		rldChan:     make(chan struct{}),
	}
}

// EventExporterService is the service structure for EventExporterS
type EventExporterService struct {
	sync.RWMutex

	cfg         *config.CGRConfig
	filterSChan chan *engine.FilterS
	connMgr     *engine.ConnManager
	server      *utils.Server
	exitChan    chan bool
	intConnChan chan rpcclient.ClientConnector
	rldChan     chan struct{}

	eeS *ees.EventExporterS
}

// GetIntenternalChan is deprecated and it will be removed shortly
func (es *EventExporterService) GetIntenternalChan() (conn chan rpcclient.ClientConnector) {
	panic("deprecated method")
}

// IsRunning returns if the service is running
func (es *EventExporterService) IsRunning() bool {
	es.RLock()
	defer es.RUnlock()
	return es != nil && es.eeS != nil
}

// ServiceName returns the service name
func (es *EventExporterService) ServiceName() string {
	return utils.EventExporterS
}

// ShouldRun returns if the service should be running
func (es *EventExporterService) ShouldRun() (should bool) {
	es.cfg.RLocks(config.EEsJson)
	should = es.cfg.EEsCfg().Enabled
	es.cfg.RUnlocks(config.EEsJson)
	return
}

// Reload handles the change of config
func (es *EventExporterService) Reload() (err error) {
	es.rldChan <- struct{}{}
	return // for the momment nothing to reload
}

// Shutdown stops the service
func (es *EventExporterService) Shutdown() (err error) {
	es.Lock()
	defer es.Unlock()
	if err = es.eeS.Shutdown(); err != nil {
		return
	}
	es.eeS = nil
	<-es.intConnChan
	return
}

// Start should handle the service start
func (es *EventExporterService) Start() (err error) {
	if es.IsRunning() {
		return fmt.Errorf("service aleady running")
	}

	fltrS := <-es.filterSChan
	es.filterSChan <- fltrS

	es.Lock()
	es.eeS = ees.NewEventExporterS(es.cfg, fltrS, es.connMgr)
	es.Unlock()
	if err != nil {

		return
	}
	es.intConnChan <- es.eeS
	return es.eeS.ListenAndServe(es.exitChan, es.rldChan)
}
