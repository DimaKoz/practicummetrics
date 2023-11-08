// staticlint implements set of static checks.
//
// Following checks are included:
//
// 1. All checks from https://golang.org/x/tools/go/analysis/passes
//
// 2. All SA checks from https://staticcheck.io/docs/checks/
//
// 3. ST1012 check from https://staticcheck.io/docs/checks/#ST1012
//
// 4. S1004 check from https://staticcheck.io/docs/checks/#S1004
//
// 5. The 1st external linter  https://github.com/sashamelentyev/usestdlibvars/
//
// 6. The 2nd external linter  https://github.com/fatih/errwrap
//
// 7. A linter which looks for 'os.Exit' calls in 'main' function of 'main' package
//
// Example:
//
//	multichecker -SA1028 <project path>
//
// Perform SA1028 analysis for given project.
// For more details run:
//
//	multichecker -help
//
// mainexit looks for os.Exit call in 'main' function of 'main' package. Run this check with following command:
//
//	multichecker -mainexit
package main
