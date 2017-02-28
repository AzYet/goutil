package go_utils

import (
	"testing"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"fmt"
	"time"
)

type TestConfg struct {
	List []string
	Map  map[string]int
}

func TestAutoReloader(t *testing.T) {
	//test init
	testFile := `{"list":["abc","cd"], "map":{"a":1,"b":2}}`
	if err := ioutil.WriteFile("/tmp/test.json", []byte(testFile), 0644); err != nil {
		t.Error("failed to prepare file" + err.Error())
	}
	i, cs := AutoReloader("/tmp/test.json", &TestConfg{}, logrus.New())
	c, ok := i.(*TestConfg)
	cStr := fmt.Sprintln(c)
	if !ok {
		t.Errorf("return wrong type:%T %v", i, i)
	}
	fmt.Println("init value", cStr)
	if len(c.List) == 0 {
		t.Error("failed to init")
	}


	// test reload
	testFile = `{"list":["abcd","cde"], "map":{"a":12,"b":23}}`
	if err := ioutil.WriteFile("/tmp/test.json", []byte(testFile), 0644); err != nil {
		t.Error("failed to change file" + err.Error())
	}
	time.Sleep(time.Second)
	cs[0] <- nil
	r := <-cs[1]
	c1, ok := r.(*TestConfg)
	if !ok {
		t.Errorf("return wrong type:%T %v", i, i)
	}
	c1Str := fmt.Sprintln(c1)
	fmt.Println("new value", c1)
	if cStr == c1Str {
		t.Error("failed to refresh")
	}

	// test watch
	cs[0] <- 1
	ch, ok := (<-cs[1]).(chan interface{})
	if !ok {
		t.Errorf("return wrong type:%T %v", i, i)
	}
	testFile = `{"list":["abcde","cdef"], "map":{"a":123,"b":234}}`
	if err := ioutil.WriteFile("/tmp/test.json", []byte(testFile), 0644); err != nil {
		t.Error("failed to change file" + err.Error())
	}
	time.Sleep(time.Second)
	c3, ok := (<-ch).(*TestConfg)
	if !ok {
		t.Errorf("return wrong type:%T %v", i, i)
	}
	c3Str := fmt.Sprintln(c3)
	fmt.Println("new value from watch", c3)
	if c3Str == c1Str {
		t.Error("failed to refresh")
	}
}
