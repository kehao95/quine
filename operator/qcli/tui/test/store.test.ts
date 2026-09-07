// The reducer suite moved to the shared client-core with the data layer
// (WEB.md §3/§6). Importing it registers the same tests in this package's
// run, so the TUI suite still exercises the exact reducer it ships with.
import "../../client-core/test/reducer.test"
