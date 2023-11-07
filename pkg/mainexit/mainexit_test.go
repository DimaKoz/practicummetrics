package mainexit

import (
	"fmt"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Analyzer(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), Analyzer, "./...")
	fmt.Println(results)
}
