package suite

import (
	"reflect"
	"strings"
	"testing"
)

func RunSuite(t *testing.T, s SuiteStore, isParallel bool) {
	testSuiteType := reflect.TypeOf(s)
	testNames := []string{}
	for i := 0; i < testSuiteType.NumMethod(); i++ {
		m := testSuiteType.Method(i)
		if strings.HasPrefix(m.Name, "Test") && len(m.Name) > 4 {
			testNames = append(testNames, m.Name)
		}
	}
	if len(testNames) == 0 {
		t.Log("no tests to run")
	}
	setupTest, setupTestFound := testSuiteType.MethodByName("SetupTest")
	tearDownTest, tearDownTestFound := testSuiteType.MethodByName("TearDownTest")
	for _, name := range testNames {
		t.Run(name, func(t *testing.T) {
			name := name
			var testData *TestData
			if isParallel {
				t.Parallel()
			}
			// Setup
			if setupTestFound {
				setupTestArgs := []reflect.Value{
					reflect.ValueOf(s),
					reflect.ValueOf(t),
				}
				testData = setupTest.Func.Call(setupTestArgs)[0].Interface().(*TestData)
				s.StoreDB(t.Name(), testData.DB)
			}
			// Execution
			testMethod, testMethodFound := testSuiteType.MethodByName(name)
			if testMethodFound {
				testMethodArgs := []reflect.Value{
					reflect.ValueOf(s),
					reflect.ValueOf(t),
				}
				testMethod.Func.Call(testMethodArgs)
			}
			// TearDown
			if tearDownTestFound {
				tearDownTestArgs := []reflect.Value{
					reflect.ValueOf(s),
					reflect.ValueOf(t),
					reflect.ValueOf(testData),
				}
				tearDownTest.Func.Call(tearDownTestArgs)
			}
		})
	}
}
