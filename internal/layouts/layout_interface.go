package layouts

// LayoutInterface is the minimal contract for all page layouts.
// Any layout that can render itself to HTML satisfies this interface.
type LayoutInterface interface {
	ToHTML() string
}
