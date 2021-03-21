package cmdline

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var testResults []string

type testCfg1 struct {
	ID string `description:"Identifier" barevalue:"True"`
}

func (cfg testCfg1) Init() error {
	testResults = append(testResults, fmt.Sprintf("testCfg1 Init %s", cfg.ID))
	return nil
}

func (cfg testCfg1) Prepare() error {
	testResults = append(testResults, fmt.Sprintf("testCfg1 Prepare %s", cfg.ID))
	return nil
}

func (cfg testCfg1) Run() error {
	testResults = append(testResults, fmt.Sprintf("testCfg1 Run %s", cfg.ID))
	return nil
}

type testCfg2 struct {
	Value string `required:"True"`
}

type testCfg3 struct {
	Value string `default:"98765"`
}

func (cfg testCfg3) Run() error {
	testResults = append(testResults, fmt.Sprintf("testCfg3 Run %s", cfg.Value))
	return nil
}

var tcf4 testCfg4

type testCfg4 struct {
	I1  int
	I2  int8
	I3  int16
	I4  int32
	I5  int64
	F1  float32
	F2  float64
	S   string
	Ls  []string
	Li  []int
	Mss map[string]string
}

func (cfg testCfg4) Run() error {
	tcf4 = cfg
	return nil
}

func TestPhases(t *testing.T) {
	cl := NewCmdline()
	cl.AddConfigType("PhaseTest", "Phase Test", testCfg1{})
	testResults = make([]string, 0)
	err := cl.ParseAndRun([]string{"--PhaseTest", "ID=test1", "--PhaseTest", "ID=test2"}, []string{"Init", "Prepare", "Run"})
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(testResults, []string{
		"testCfg1 Init test1",
		"testCfg1 Init test2",
		"testCfg1 Prepare test1",
		"testCfg1 Prepare test2",
		"testCfg1 Run test1",
		"testCfg1 Run test2",
	}) {
		t.Error("Actual results did not match expected")
	}
}

func TestRequired(t *testing.T) {
	cl := NewCmdline()
	cl.AddConfigType("test", "", testCfg2{}, Required)
	err := cl.ParseAndRun([]string{}, []string{})
	if err == nil {
		t.Error("Missing required config type did not produce error")
	}
	err = cl.ParseAndRun([]string{"--test"}, []string{})
	if err == nil {
		t.Error("Missing required field did not produce error")
	}
}

func TestDefault(t *testing.T) {
	cl := NewCmdline()
	cl.AddConfigType("test", "", testCfg3{})
	testResults = make([]string, 0)
	err := cl.ParseAndRun([]string{"--test"}, []string{"Run"})
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(testResults, []string{
		"testCfg3 Run 98765",
	}) {
		t.Error("Actual results did not match expected")
	}
}

func TestSingleton(t *testing.T) {
	cl := NewCmdline()
	cl.AddConfigType("test", "", testCfg2{}, Singleton)
	err := cl.ParseAndRun([]string{"--test", "value=abc"}, []string{})
	if err != nil {
		t.Error(err)
	}
	err = cl.ParseAndRun([]string{"--test", "value=abc", "--test", "value=def"}, []string{})
	if err == nil {
		t.Error("Two values of same singleton did not produce error")
	}
}

func TestExclusive(t *testing.T) {
	cl := NewCmdline()
	cl.AddConfigType("test1", "", testCfg2{}, Exclusive)
	cl.AddConfigType("test2", "", testCfg3{})
	err := cl.ParseAndRun([]string{"--test1", "value=abc"}, []string{})
	if err != nil {
		t.Error(err)
	}
	if cl.WhatRan() != "test1" {
		t.Error("WhatRan was not set correctly")
	}
	err = cl.ParseAndRun([]string{"--test1", "value=abc", "--test2", "value=def"}, []string{})
	if err == nil {
		t.Error("Exclusive item with other item did not produce error")
	}
}

func TestHidden(t *testing.T) {
	cl := NewCmdline()
	buf := new(bytes.Buffer)
	cl.SetOutput(buf)
	cl.AddConfigType("test1", "", testCfg2{})
	cl.AddConfigType("xyzzy", "", testCfg3{}, Hidden)
	err := cl.ShowHelp()
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(buf.String(), "xyzzy") {
		t.Error("Hidden config type appeared in help")
	}
}

func TestSection(t *testing.T) {
	sec1 := &ConfigSection{
		Description: "Section 1",
		Order:       1,
	}
	sec2 := &ConfigSection{
		Description: "Section 2",
		Order:       2,
	}
	cl := NewCmdline()
	buf := new(bytes.Buffer)
	cl.SetOutput(buf)
	cl.AddConfigType("abcdef", "", testCfg2{}, Section(sec2))
	cl.AddConfigType("zyxwvu", "", testCfg3{}, Section(sec1))
	err := cl.ShowHelp()
	if err != nil {
		t.Error(err)
	}
	re := regexp.MustCompile("(?s)zyxwvu.*abcdef")
	if !re.MatchString(buf.String()) {
		t.Error("Help output not ordered as expected by sections")
	}
}

func TestBarevalue(t *testing.T) {
	cl := NewCmdline()
	cl.AddConfigType("test", "", testCfg1{})
	testResults = make([]string, 0)
	err := cl.ParseAndRun([]string{"--test", "plugh"}, []string{"Run"})
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(testResults, []string{
		"testCfg1 Run plugh",
	}) {
		t.Error("Actual results did not match expected")
	}
}

func TestSimpleTypes(t *testing.T) {
	cl := NewCmdline()
	cl.AddConfigType("test", "", testCfg4{})
	tcf4 = testCfg4{}
	err := cl.ParseAndRun([]string{"--test",
		"i1=1",
		"i2=2",
		"i3=3",
		"i4=4",
		"i5=5",
		"f1=1.0",
		"f2=2.0",
		"s=hello",
	}, []string{"Run"})
	if err != nil {
		t.Error(err)
	}
	if tcf4.I1 != 1 ||
		tcf4.I2 != 2 ||
		tcf4.I3 != 3 ||
		tcf4.I4 != 4 ||
		tcf4.I5 != 5 ||
		tcf4.F1 != 1.0 ||
		tcf4.F2 != 2.0 ||
		tcf4.S != "hello" {
		t.Error("Parameters did not receive expected values")
	}
}

func TestMultiString(t *testing.T) {
	cl := NewCmdline()
	cl.AddConfigType("test", "", testCfg4{})
	tcf4 = testCfg4{}
	err := cl.ParseAndRun([]string{"--test",
		"ls=hello",
		"ls=goodbye",
	}, []string{"Run"})
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(tcf4.Ls, []string{
		"hello",
		"goodbye",
	}) {
		t.Error("Actual results did not match expected")
	}
}

func TestJSON(t *testing.T) {
	cl := NewCmdline()
	cl.AddConfigType("test", "", testCfg4{})
	tcf4 = testCfg4{}
	err := cl.ParseAndRun([]string{"--test",
		"li=[1, 2, 3]",
		`mss={"a": "b", "c": "d"}`,
	}, []string{"Run"})
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(tcf4.Li, []int{1, 2, 3}) ||
		!reflect.DeepEqual(tcf4.Mss, map[string]string{"a": "b", "c": "d"}) {
		t.Error("Actual results did not match expected")
	}
}

var yamlData = `
---
- test:
    i1: 37
    f1: 2.3
    s: "hello"
    ls:
      - hello
      - goodbye
    li: [1, 2, 3]
    mss:
      a: b
      c: d
`

func TestYAML(t *testing.T) {
	yamlFile, err := ioutil.TempFile("", "test*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(yamlFile.Name())
	}()
	_, err = fmt.Fprintf(yamlFile, yamlData)
	cl := NewCmdline()
	cl.AddConfigType("test", "", testCfg4{})
	tcf4 = testCfg4{}
	err = cl.ParseAndRun([]string{"--test", "--config", yamlFile.Name()}, []string{"Run"})
	if err != nil {
		t.Error(err)
	}
	if tcf4.I1 != 37 ||
		tcf4.F1 != 2.3 ||
		tcf4.S != "hello" ||
		!reflect.DeepEqual(tcf4.Ls, []string{"hello", "goodbye"}) ||
		!reflect.DeepEqual(tcf4.Li, []int{1, 2, 3}) ||
		!reflect.DeepEqual(tcf4.Mss, map[string]string{"a": "b", "c": "d"}) {
		t.Error("Actual results did not match expected")
	}
}
