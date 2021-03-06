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
package utils

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCGREventHasField(t *testing.T) {
	//empty check
	cgrEvent := new(CGREvent)
	rcv := cgrEvent.HasField("")
	if rcv {
		t.Error("Expecting: false, received: ", rcv)
	}
	//normal check
	cgrEvent = &CGREvent{
		Event: map[string]interface{}{
			Usage: time.Duration(20 * time.Second),
		},
	}
	rcv = cgrEvent.HasField("Usage")
	if !rcv {
		t.Error("Expecting: true, received: ", rcv)
	}
}

func TestCGREventCheckMandatoryFields(t *testing.T) {
	//empty check
	cgrEvent := new(CGREvent)
	fldNames := []string{}
	err := cgrEvent.CheckMandatoryFields(fldNames)
	if err != nil {
		t.Error(err)
	}
	cgrEvent = &CGREvent{
		Event: map[string]interface{}{
			Usage:   time.Duration(20 * time.Second),
			"test1": 1,
			"test2": 2,
			"test3": 3,
		},
	}
	//normal check
	fldNames = []string{"test1", "test2"}
	err = cgrEvent.CheckMandatoryFields(fldNames)
	if err != nil {
		t.Error(err)
	}
	//MANDATORY_IE_MISSING
	fldNames = []string{"test4", "test5"}
	err = cgrEvent.CheckMandatoryFields(fldNames)
	if err == nil || err.Error() != NewErrMandatoryIeMissing("test4").Error() {
		t.Errorf("Expected %s, received %s", NewErrMandatoryIeMissing("test4"), err)
	}
}

func TestCGREventFielAsString(t *testing.T) {
	//empty check
	cgrEvent := new(CGREvent)
	fldname := "test"
	_, err := cgrEvent.FieldAsString(fldname)
	if err == nil || err.Error() != ErrNotFound.Error() {
		t.Errorf("Expected %s, received %s", ErrNotFound, err)
	}
	//normal check
	cgrEvent = &CGREvent{
		Event: map[string]interface{}{
			Usage:   time.Duration(20 * time.Second),
			"test1": 1,
			"test2": 2,
			"test3": 3,
		},
	}
	fldname = "test1"
	rcv, err := cgrEvent.FieldAsString(fldname)
	if err != nil {
		t.Error(err)
	}
	if rcv != "1" {
		t.Errorf("Expected: 1, received %+q", rcv)
	}

}

func TestLibRoutesUsage(t *testing.T) {
	se := &CGREvent{
		Tenant: "cgrates.org",
		ID:     "supplierEvent1",
		Event: map[string]interface{}{
			Usage: time.Duration(20 * time.Second),
		},
	}
	seErr := &CGREvent{
		Tenant: "cgrates.org",
		ID:     "supplierEvent1",
		Event:  make(map[string]interface{}),
	}
	answ, err := se.FieldAsDuration(Usage)
	if err != nil {
		t.Error(err)
	}
	if answ != se.Event[Usage] {
		t.Errorf("Expecting: %+v, received: %+v", se.Event[Usage], answ)
	}
	answ, err = seErr.FieldAsDuration(Usage)
	if err != ErrNotFound {
		t.Error(err)
	}
}

func TestCGREventFieldAsTime(t *testing.T) {
	se := &CGREvent{
		Tenant: "cgrates.org",
		ID:     "supplierEvent1",
		Event: map[string]interface{}{
			AnswerTime: time.Now(),
		},
	}
	seErr := &CGREvent{
		Tenant: "cgrates.org",
		ID:     "supplierEvent1",
		Event:  make(map[string]interface{}),
	}
	answ, err := se.FieldAsTime(AnswerTime, "UTC")
	if err != nil {
		t.Error(err)
	}
	if answ != se.Event[AnswerTime] {
		t.Errorf("Expecting: %+v, received: %+v", se.Event[AnswerTime], answ)
	}
	answ, err = seErr.FieldAsTime(AnswerTime, "CET")
	if err != ErrNotFound {
		t.Error(err)
	}
}

func TestCGREventFieldAsString(t *testing.T) {
	se := &CGREvent{
		Tenant: "cgrates.org",
		ID:     "supplierEvent1",
		Event: map[string]interface{}{
			"supplierprofile1": "Supplier",
			"UsageInterval":    time.Duration(1 * time.Second),
			"PddInterval":      "1s",
			"Weight":           20.0,
		},
	}
	answ, err := se.FieldAsString("UsageInterval")
	if err != nil {
		t.Error(err)
	}
	if answ != "1s" {
		t.Errorf("Expecting: %+v, received: %+v", se.Event["UsageInterval"], answ)
	}
	answ, err = se.FieldAsString("PddInterval")
	if err != nil {
		t.Error(err)
	}
	if answ != se.Event["PddInterval"] {
		t.Errorf("Expecting: %+v, received: %+v", se.Event["PddInterval"], answ)
	}
	answ, err = se.FieldAsString("supplierprofile1")
	if err != nil {
		t.Error(err)
	}
	if answ != se.Event["supplierprofile1"] {
		t.Errorf("Expecting: %+v, received: %+v", se.Event["supplierprofile1"], answ)
	}
	answ, err = se.FieldAsString("Weight")
	if err != nil {
		t.Error(err)
	}
	if answ != "20" {
		t.Errorf("Expecting: %+v, received: %+v", se.Event["Weight"], answ)
	}
}

func TestCGREventFieldAsFloat64(t *testing.T) {
	se := &CGREvent{
		Tenant: "cgrates.org",
		ID:     "supplierEvent1",
		Event: map[string]interface{}{
			AnswerTime:         time.Now(),
			"supplierprofile1": "Supplier",
			"UsageInterval":    "54.2",
			"PddInterval":      "1s",
			"Weight":           20.0,
		},
	}
	answ, err := se.FieldAsFloat64("UsageInterval")
	if err != nil {
		t.Error(err)
	}
	if answ != float64(54.2) {
		t.Errorf("Expecting: %+v, received: %+v", se.Event["UsageInterval"], answ)
	}
	answ, err = se.FieldAsFloat64("Weight")
	if err != nil {
		t.Error(err)
	}
	if answ != float64(20.0) {
		t.Errorf("Expecting: %+v, received: %+v", se.Event["Weight"], answ)
	}
	answ, err = se.FieldAsFloat64("PddInterval")
	if err == nil || err.Error() != `strconv.ParseFloat: parsing "1s": invalid syntax` {
		t.Errorf("Expected %s, received %s", `strconv.ParseFloat: parsing "1s": invalid syntax`, err)
	}
	if answ != 0 {
		t.Errorf("Expecting: %+v, received: %+v", 0, answ)
	}

	if _, err := se.FieldAsFloat64(AnswerTime); err == nil || !strings.HasPrefix(err.Error(), "cannot convert field") {
		t.Errorf("Unexpected error : %+v", err)
	}
	if _, err := se.FieldAsFloat64(Account); err == nil || err.Error() != ErrNotFound.Error() {
		t.Errorf("Expected %s, received %s", ErrNotFound, err)
	}
	// }
}

// }
func TestCGREventTenantID(t *testing.T) {
	//empty check
	cgrEvent := new(CGREvent)
	rcv := cgrEvent.TenantID()
	eOut := ":"
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
	//normal check
	cgrEvent = &CGREvent{
		Tenant: "cgrates.org",
		ID:     "supplierEvent1",
	}
	rcv = cgrEvent.TenantID()
	eOut = "cgrates.org:supplierEvent1"
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}

}

func TestCGREventClone(t *testing.T) {
	now := time.Now()
	ev := &CGREvent{
		Tenant: "cgrates.org",
		ID:     "supplierEvent1",
		Time:   &now,
		Event: map[string]interface{}{
			AnswerTime:         time.Now(),
			"supplierprofile1": "Supplier",
			"UsageInterval":    "54.2",
			"PddInterval":      "1s",
			"Weight":           20.0,
		},
	}
	cloned := ev.Clone()
	if !reflect.DeepEqual(ev, cloned) {
		t.Errorf("Expecting: %+v, received: %+v", ev, cloned)
	}
	if cloned.Time == ev.Time {
		t.Errorf("Expecting: different pointer but received: %+v", cloned.Time)
	}
}

func TestCGREventconsumeArgDispatcher(t *testing.T) {
	//empty check
	var opts map[string]interface{}
	rcv := getArgDispatcherFromOpts(opts)
	if rcv != nil {
		t.Errorf("Expecting: nil, received: %+v", rcv)
	}
	//nil check
	opts = nil
	rcv = getArgDispatcherFromOpts(opts)
	if rcv != nil {
		t.Errorf("Expecting: nil, received: %+v", rcv)
	}
	//normal check without APIkey
	routeID := "route"
	opts = map[string]interface{}{
		RouteID: routeID,
	}
	eOut := &ArgDispatcher{
		RouteID: &routeID,
	}
	rcv = getArgDispatcherFromOpts(opts)
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting:  %+v, received: %+v", eOut, rcv)
	}
	//check if *route_id was deleted
	if _, has := opts[RouteID]; has {
		t.Errorf("*route_id wasn't deleted")
	}
	//normal check with routeID and APIKey
	apiKey := "key"
	opts = map[string]interface{}{RouteID: routeID, APIKey: apiKey}
	eOut.APIKey = &apiKey

	rcv = getArgDispatcherFromOpts(opts)
	//check if *api_key and *route_id was deleted
	if _, has := opts[APIKey]; has {
		t.Errorf("*api_key wasn't deleted")
	} else if _, has := opts[RouteID]; has {
		t.Errorf("*route_id wasn't deleted")
	}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting:  %+v, received: %+v", eOut, rcv)
	}
}

func TestCGREventconsumeRoutePaginator(t *testing.T) {
	//empty check
	var opts map[string]interface{}
	rcv, err := getRoutePaginatorFromOpts(opts)
	if err != nil {
		t.Error(err)
	}
	eOut := new(Paginator)
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting:  %+v, received: %+v", eOut, rcv)
	}
	opts = nil
	rcv, err = getRoutePaginatorFromOpts(opts)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting:  %+v, received: %+v", eOut, rcv)
	}
	//normal check
	opts = map[string]interface{}{
		RoutesLimit:  18,
		RoutesOffset: 20,
	}

	eOut = &Paginator{
		Limit:  IntPointer(18),
		Offset: IntPointer(20),
	}
	rcv, err = getRoutePaginatorFromOpts(opts)
	if err != nil {
		t.Error(err)
	}
	//check if *suppliers_limit and *suppliers_offset was deleted
	if _, has := opts[RoutesLimit]; has {
		t.Errorf("*suppliers_limit wasn't deleted")
	} else if _, has := opts[RoutesOffset]; has {
		t.Errorf("*suppliers_offset wasn't deleted")
	}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting:  %+v, received: %+v", eOut, rcv)
	}
	//check without *suppliers_limit, but with *suppliers_offset
	opts = map[string]interface{}{
		RoutesOffset: 20,
	}

	eOut = &Paginator{
		Offset: IntPointer(20),
	}
	rcv, err = getRoutePaginatorFromOpts(opts)
	if err != nil {
		t.Error(err)
	}
	//check if *suppliers_limit and *suppliers_offset was deleted
	if _, has := opts[RoutesLimit]; has {
		t.Errorf("*suppliers_limit wasn't deleted")
	} else if _, has := opts[RoutesOffset]; has {
		t.Errorf("*suppliers_offset wasn't deleted")
	}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting:  %+v, received: %+v", eOut, rcv)
	}
	//check with notAnInt at *suppliers_limit
	opts = map[string]interface{}{
		RoutesLimit: "Not an int",
	}
	eOut = new(Paginator)
	rcv, err = getRoutePaginatorFromOpts(opts)
	if err == nil {
		t.Error("Expected error")
	}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting:  %+v, received: %+v", eOut, rcv)
	}
	//check with notAnInt at and *suppliers_offset
	opts = map[string]interface{}{
		RoutesOffset: "Not an int",
	}
	eOut = new(Paginator)
	rcv, err = getRoutePaginatorFromOpts(opts)
	if err == nil {
		t.Error("Expected error")
	}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting:  %+v, received: %+v", eOut, rcv)
	}
}

func TestCGREventConsumeArgs(t *testing.T) {
	//empty check
	opts := map[string]interface{}{}
	eOut := ExtractedArgs{}
	// false false
	rcv, err := ExtractArgsFromOpts(opts, false, false)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
	// false true
	rcv, err = ExtractArgsFromOpts(opts, false, true)
	if err != nil {
		t.Error(err)
	}
	eOut.RoutePaginator = new(Paginator)
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
	//true false
	eOut = ExtractedArgs{
		ArgDispatcher:  new(ArgDispatcher),
		RoutePaginator: nil,
	}
	rcv, err = ExtractArgsFromOpts(opts, true, false)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
	//true true
	rcv, err = ExtractArgsFromOpts(opts, true, true)
	if err != nil {
		t.Error(err)
	}
	eOut = ExtractedArgs{
		RoutePaginator: new(Paginator),
		ArgDispatcher:  new(ArgDispatcher),
	}
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}

}

func TestNewCGREventWithArgDispatcher(t *testing.T) {
	eOut := new(CGREventWithArgDispatcher)
	rcv := NewCGREventWithArgDispatcher()

	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
}

func TestCGREventWithArgDispatcherClone(t *testing.T) {
	//empty check
	cgrEventWithArgDispatcher := new(CGREventWithArgDispatcher)
	rcv := cgrEventWithArgDispatcher.Clone()
	if !reflect.DeepEqual(cgrEventWithArgDispatcher, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", cgrEventWithArgDispatcher, rcv)
	}
	//nil check
	cgrEventWithArgDispatcher = nil
	rcv = cgrEventWithArgDispatcher.Clone()
	if !reflect.DeepEqual(cgrEventWithArgDispatcher, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", cgrEventWithArgDispatcher, rcv)
	}
	//normal check
	now := time.Now()
	cgrEventWithArgDispatcher = &CGREventWithArgDispatcher{
		CGREvent: &CGREvent{
			Tenant: "cgrates.org",
			ID:     "IDtest",
			Time:   &now,
			Event: map[string]interface{}{
				"test1": 1,
				"test2": 2,
				"test3": 3,
			},
		},
		ArgDispatcher: new(ArgDispatcher),
	}
	rcv = cgrEventWithArgDispatcher.Clone()
	if !reflect.DeepEqual(cgrEventWithArgDispatcher, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", cgrEventWithArgDispatcher, rcv)
	}
	//check vars
	apiKey := "apikey"
	routeID := "routeid"

	rcv.ArgDispatcher = &ArgDispatcher{
		APIKey:  &apiKey,
		RouteID: &routeID,
	}
	if reflect.DeepEqual(cgrEventWithArgDispatcher.ArgDispatcher, rcv.ArgDispatcher) {
		t.Errorf("Expected to be different")
	}

}
