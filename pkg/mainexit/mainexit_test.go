package mainexit

import (
	"fmt"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func Test_Analyzer(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), Analyzer, "./...")
	fmt.Println(results)
}
