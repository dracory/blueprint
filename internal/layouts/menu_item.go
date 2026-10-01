package layouts

// MenuItem represents a navigation item rendered by the page navbars.
// Icon holds a raw HTML snippet (e.g. a bootstrap-icons <i> tag).
type MenuItem struct {
	Icon   string
	Title  string
	URL    string
	Target string
}
