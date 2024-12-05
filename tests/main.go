package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fastly/waf-simulator-automation/ngwaf"
	"gopkg.in/yaml.v2"
)

type Rule struct {
	Tests []Test `yaml:"tests"`
}

type Test struct {
	Name     string `yaml:"name"`
	Site     string `yaml:"site"`
	RuleID   string `yaml:"rule_id"`
	Type     string `yaml:"type"`
	Request  string `yaml:"request"`
	Response string `yaml:"response"`
	Expect   struct {
		WafResponse int            `yaml:"waf_response"`
		Signals     []ngwaf.Signal `yaml:"signals"`
	} `yaml:"expect"`
}

// converts yaml config to struct
func (r *Rule) getConf(file string) (*Rule, error) {

	fileType := filepath.Ext(file)
	if fileType != ".yml" && fileType != ".yaml" {
		return nil, fmt.Errorf("validation err: unsupported file type %s, file type must be .yml or .yaml", fileType)
	}
	ymlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("ymlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(ymlFile, r)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return r, nil
}

// tests rules
func testSimulations(file string, sc *ngwaf.Client) (errCnt int, err error) {
	var rule Rule
	_, err = rule.getConf(file)
	if err != nil {
		return errCnt, err
	}

	for _, test := range rule.Tests {

		testBody := ngwaf.Simulation{
			SampleRequest:  test.Request,
			SampleResponse: test.Response,
		}

		response, err := sc.SimulationTest(testBody, test.Site)
		if err != nil {
			fmt.Println(err)
		}

		if !validateTest(response, test) {
			errCnt += 1
		}

	}

	return errCnt, err
}

// validates the simulation output against the expected output
func validateTest(response ngwaf.SimulationOutput, test Test) bool {

	pass := true

	// Check waf response code
	if response.Data.WafResponse != test.Expect.WafResponse {
		fmt.Printf("%s failed: WafResponse %d != %d\n", test.Name, response.Data.WafResponse, test.Expect.WafResponse)
		pass = false
	}

	// Create a map for quick lookups
	responseSignals := make(map[string]ngwaf.Signal)
	for _, signal := range response.Data.Signals {
		responseSignals[signal.Name] = signal
	}

	for _, expectedSignal := range test.Expect.Signals {
		found := false
		var failureReason string

		for _, respSignal := range response.Data.Signals {
			if expectedSignal.Type != respSignal.Type {
				failureReason = fmt.Sprintf("Signal mismatch: expected %q, got %q", expectedSignal.Type, respSignal.Type)
				continue
			}
			if expectedSignal.Value != "" && expectedSignal.Value != respSignal.Value {
				failureReason = fmt.Sprintf("Value mismatch: expected %q, got %q", expectedSignal.Value, respSignal.Value)
				break
			}
			if expectedSignal.Location != "" && expectedSignal.Location != respSignal.Location {
				failureReason = fmt.Sprintf("Location mismatch: expected %q, got %q", expectedSignal.Location, respSignal.Location)
				break
			}
			if expectedSignal.Name != "" && expectedSignal.Name != respSignal.Name {
				failureReason = fmt.Sprintf("Name mismatch: expected %q, got %q", expectedSignal.Name, respSignal.Name)
				break
			}
			if expectedSignal.Detector != "" && expectedSignal.Detector != respSignal.Detector {
				failureReason = fmt.Sprintf("Detector mismatch: expected %q, got %q", expectedSignal.Detector, respSignal.Detector)
				break
			}
			if expectedSignal.Redaction != 0 && expectedSignal.Redaction != respSignal.Redaction {
				failureReason = fmt.Sprintf("Redaction mismatch: expected %v, got %v", expectedSignal.Redaction, respSignal.Redaction)
				break
			}

			found = true
			break
		}

		if !found {
			fmt.Printf("%s failed: Reason: %s\n", test.Name, failureReason)
			pass = false
		}
	}

	return pass

}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func getTestFiles() (files []string, err error) {

	err = filepath.Walk("tests/rules",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if dir, err := isDirectory(path); !dir {
				if err != nil {
					return err
				}
				files = append(files, path)
			}
			return nil
		})

	return files, err

}

func main() {

	fc := ngwaf.NewTokenClient(os.Getenv("SIGSCI_EMAIL"), os.Getenv("SIGSCI_TOKEN"), os.Getenv("SIGSCI_CORP"))

	files, err := getTestFiles()
	if err != nil {
		fmt.Printf("Error getting test files: %v\n", err)
		return
	}

	// Check if the files array is empty
	if len(files) == 0 {
		panic(fmt.Errorf("No tests found."))
	}

	errCnt := 0
	for _, f := range files {
		cnt, err := testSimulations(f, &fc)
		if err != nil {
			fmt.Printf("Error in test simulations for file %s: %v\n", f, err)
			continue
		}
		errCnt += cnt
	}
	if errCnt > 0 {
		panic(fmt.Errorf("%v test(s) failed", errCnt))
	}

}
