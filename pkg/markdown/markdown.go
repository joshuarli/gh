package markdown

// Option is a no-op rendering option retained for API compatibility with callers.
type Option struct{}

func WithoutIndentation() Option { return Option{} }

func WithWrap(_ int) Option { return Option{} }

func WithTheme(_ string) Option { return Option{} }

func WithBaseURL(_ string) Option { return Option{} }

// Render returns text as-is. Glamour markdown rendering has been removed.
func Render(text string, _ ...Option) (string, error) {
	return text, nil
}
