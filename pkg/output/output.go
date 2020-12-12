package output

type format string
type Format format

const (
	FormatJSON         Format = "json"
	FormatJSONL        Format = "jsonl"
	FormatMarkdown     Format = "markdown"
	FormatMarkdownWide Format = "markdown-wide"
	FormatJUnit        Format = "junit"
)

var (
	formats = map[string]Format{
		string(FormatJSON):         FormatJSON,
		string(FormatJSONL):        FormatJSONL,
		string(FormatMarkdown):     FormatMarkdown,
		string(FormatMarkdownWide): FormatMarkdownWide,
		string(FormatJUnit):        FormatJUnit,
	}
	FormatNames = make([]string, 0, len(formats))
)

func init() {
	for k := range formats {
		FormatNames = append(FormatNames, k)
	}
}

func ParseFormatName(s string) (Format, bool) {
	f, ok := formats[s]
	return f, ok
}
