package go_utils

import (
	"encoding/json"
	"bytes"
	"github.com/Sirupsen/logrus"
	"path/filepath"
	"sync"
	"io/ioutil"
	"golang.org/x/exp/inotify"
	"strings"
	"path"
)

// watch json file and decode it into pointer of struct, t must be a pointer
// init value will be returned as first return value
// chan[0] is use to send operator, chan[1] is use to return value
// send nil will load the latest value, send non-nil will return a chan to receive changes

func AutoReloader(path string, t interface{}, Logger *logrus.Logger) (interface{}, []chan interface{}) {
	if Logger == nil {
		Logger = logrus.New()
	}
	d, f := filepath.Split(path)
	if d == "" {
		d = "."
	}
	var (
		initiated bool
		w = new(sync.WaitGroup)
		latest interface{}
		ops = []chan interface{}{make(chan interface{}), make(chan interface{})}
	)
	w.Add(1)
	go func() {
		fileChan := ReadAndWatchFile(d, Logger, f)
		watchersChan := []chan interface{}{}
		for {
			select {
			case nb := <-fileChan:
				if err := json.NewDecoder(bytes.NewReader(nb.Bz)).Decode(&t); err != nil {
					if !initiated {
						Logger.Panicf("failed to decode to config: %v.", err)
					} else {
						Logger.Warn(err)
					}
				} else {
					latest = t
					bz := new(bytes.Buffer)
					json.Indent(bz, nb.Bz, "", " ")
					Logger.Info("config reloaded.")
					Logger.Info(bz.String())
					if !initiated {
						initiated = true
						w.Done()
					}
					for _, c := range watchersChan {
						select {
						case c <- latest:
						default:
						}
					}
				}
			case op := <-ops[0]:
				if op == nil {
					//load
					ops[1] <- latest
				} else {
					//add watch
					c := make(chan interface{}, 1)
					watchersChan = append(watchersChan, c)
					ops[1] <- c
				}
			}
		}
	}()
	w.Wait()
	return latest, ops
}

type NameBytes struct{ Name string; Bz []byte }

func ReadAndWatchFile(dir string, Logger *logrus.Logger, fileList ...string) chan NameBytes {
	bzChan := make(chan NameBytes, len(fileList))
	watcher, err := inotify.NewWatcher()
	if err != nil {
		Logger.Errorln("failed to create watcher", err)
		return nil
	}
	readSendFn := func(fullPath string) {
		_, fileName := filepath.Split(fullPath)
		if bz, err := ioutil.ReadFile(fullPath); err != nil {
			Logger.WithFields(logrus.Fields{"err": err, "path":fullPath}).Error("read error")
		} else {
			bzChan <- NameBytes{fileName, bz}
		}
	}
	nameMap := make(map[string]bool)
	for _, f := range fileList {
		nameMap[f] = true
	}
	go func() {
		for {
			select {
			case event := <-watcher.Event:
				if (event.Mask & inotify.IN_CLOSE_WRITE == inotify.IN_CLOSE_WRITE ||
					event.Mask & inotify.IN_MOVED_TO == inotify.IN_MOVED_TO ) &&
					nameMap[event.Name[strings.LastIndex(event.Name, "/") + 1:]] {
					Logger.WithFields(logrus.Fields{"event": event}).Info("file change detected, reload file.")
					readSendFn(event.Name)
				}
			case err := <-watcher.Error:
				Logger.Println("error:", err)
			}
		}
	}()
	err = watcher.Watch(dir)
	if err != nil {
		Logger.Panicln("failed to watch", dir, err)
	} else {
		Logger.Infof("watching %v for %v changes.", dir, fileList)
	}
	for _, f := range fileList {
		go readSendFn(path.Join(dir, f))
	}
	return bzChan
}
