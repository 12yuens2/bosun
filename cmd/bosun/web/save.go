package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"bosun.org/cmd/bosun/conf"
	"github.com/MiniProfiler/go/miniprofiler"
)

func SaveConfig(t miniprofiler.Timer, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	data := struct {
		Filename string
		Config   string
		Diff     string
		User     string
		Message  string
		Other    []string
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	if data.User != "" && !userCanOverwriteUsername(r) {
		http.Error(w, "Not Authorized to set User", 400)
		return nil, nil
	} else if data.User == "" {
		data.User = getUsername(r)
	}
	backup, err := ioutil.ReadFile(data.Filename)
	if err != nil {
		return nil, err
	}
	err = schedule.RuleConf.SaveRawText(data.Filename, data.Config, data.Diff, data.User, data.Message, data.Other...)
	if err != nil {
		ioutil.WriteFile(data.Filename, []byte(backup), os.FileMode(int(0640)))
		return nil, err
	}
	fmt.Fprint(w, "save successful")
	return nil, nil
}

func DiffConfig(t miniprofiler.Timer, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	data := struct {
		Config  string
		Message string
		User    string
		Other   []string
	}{}
	if data.User != "" && !userCanOverwriteUsername(r) {
		http.Error(w, "Not Authorized to set User", 400)
		return nil, nil
	} else if data.User == "" {
		data.User = getUsername(r)
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	diff, err := schedule.RuleConf.RawDiff(data.Config)
	if err != nil {
		return nil, err
	}
	fmt.Fprint(w, diff)
	return nil, nil
}

func ConfigRunningHash(t miniprofiler.Timer, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	hash := schedule.RuleConf.GetHash()
	return struct {
		Hash string
	}{
		hash,
	}, nil
}

func BulkEdit(t miniprofiler.Timer, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	bulkEdit := conf.BulkEditRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&bulkEdit); err != nil {
		return nil, err
	}
	err := schedule.RuleConf.BulkEdit(bulkEdit)
	if err != nil {
		return nil, err
	}
	fmt.Fprint(w, "edit successful")
	return nil, nil
}

func SaveEnabled(t miniprofiler.Timer, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return schedule.SystemConf.SaveEnabled(), nil
}
