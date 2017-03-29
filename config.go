package goutil

import (
	"encoding/json"
	"bytes"
	"path/filepath"
	"sync"
	"io/ioutil"
	"golang.org/x/exp/inotify"
	"strings"
	"path"
	"reflect"
)

// watch json file and decode it into pointer of struct, t must not be a value of struct not a pointer
// init value will be returned as first return value
// chan[0] is use to send operator, chan[1] is use to return value
// send nil will load the latest value, send non-nil will return a chan to receive changes

func AutoReloader(path string, t interface{}) (interface{}, func() interface{}, func() chan interface{}) {
	d, f := filepath.Split(path)
	if d == "" {
		d = "."
	}
	var (
		initiated bool
		w = new(sync.WaitGroup)
		latest interface{}
		in, out = make(chan int), make(chan interface{})
	)
	w.Add(1)
	go func() {
		fileChan := ReadAndWatchFile(d, f)
		var watchersChan  []chan interface{}
		for {
			select {
			case nb := <-fileChan:
				t1 := reflect.New(reflect.TypeOf(t)).Interface()
				if err := json.NewDecoder(bytes.NewReader(nb.Bz)).Decode(&t1); err != nil {
					if !initiated {
						panic(err)
					} else {
						break
					}
				}
				latest = t1
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

			case op := <-in:
				if op == 0 {
					//load
					out <- latest
				} else {
					//add watch
					c := make(chan interface{}, 1)
					watchersChan = append(watchersChan, c)
					out <- c
				}
			}
		}
	}()
	w.Wait()
	load := func() interface{} {
		in <- 0
		return <-out
	}
	watch := func() chan interface{} {
		in <- 1
		return (<-out).(chan interface{})
	}
	return latest, load, watch
}

type NameBytes struct{ Name string; Bz []byte }

func ReadAndWatchFile(dir string, fileList ...string) chan NameBytes {
	bzChan := make(chan NameBytes, len(fileList))
	watcher, err := inotify.NewWatcher()
	if err != nil {
		return nil
	}
	readSendFn := func(fullPath string) {
		_, fileName := filepath.Split(fullPath)
		if bz, err := ioutil.ReadFile(fullPath); err != nil {
		} else {
			bzChan <- NameBytes{fileName, bz}
		}
	}
	nameMap := make(map[string]bool)
	for _, f := range fileList {
		nameMap[f] = true
	}
	go func() {
		for event := range watcher.Event {
			if (event.Mask & inotify.IN_CLOSE_WRITE == inotify.IN_CLOSE_WRITE ||
				event.Mask & inotify.IN_MOVED_TO == inotify.IN_MOVED_TO ) &&
				nameMap[event.Name[strings.LastIndex(event.Name, "/") + 1:]] {
				readSendFn(event.Name)
			}
		}
	}()
	err = watcher.Watch(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range fileList {
		readSendFn(path.Join(dir, f))
	}
	return bzChan
}
