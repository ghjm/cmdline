package cmdline

import (
	"fmt"
	"reflect"
	"testing"
)

var testResults []string

type testPhasesCfg struct {
	ID string `description:"Identifier"`
}

func (cfg testPhasesCfg) Init() error {
	testResults = append(testResults, fmt.Sprintf("Init %s", cfg.ID))
	return nil
}

func (cfg testPhasesCfg) Prepare() error {
	testResults = append(testResults, fmt.Sprintf("Prepare %s", cfg.ID))
	return nil
}
func (cfg testPhasesCfg) Run() error {
	testResults = append(testResults, fmt.Sprintf("Run %s", cfg.ID))
	return nil
}

func TestCmdlinePhases(t *testing.T) {
	cl := NewCmdline()
	cl.AddConfigType("PhaseTest", "Phase Test", testPhasesCfg{})
	testResults = make([]string, 0)
	err := cl.ParseAndRun([]string{"--PhaseTest", "ID=test1", "--PhaseTest", "ID=test2"}, []string{"Init", "Prepare", "Run"})
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(testResults, []string{
		"Init test1",
		"Init test2",
		"Prepare test1",
		"Prepare test2",
		"Run test1",
		"Run test2",
	}) {
		t.Error("Actual results did not match expected")
	}
}
