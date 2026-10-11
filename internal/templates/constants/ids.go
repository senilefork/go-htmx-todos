package constants

// These DOM ids are the contract between templates and Go code.
// Templates put them in hx-target attributes; Go code (templatetargets,
// handlers) uses them to decide what to render. Renaming one side without
// the other breaks navigation silently — grep for the constant instead.
const (
	// The <main> element that page content swaps into.
	MainContentID = "main-content"
	// Wraps every page below the nav; a second, finer swap target.
	PageLayoutID = "page-layout"
	// The search box targets just the table, so the input keeps focus.
	TodosTableID = "todos-table"
	// Fixed container where error toasts are injected.
	GlobalToastID = "global-toast"
)
