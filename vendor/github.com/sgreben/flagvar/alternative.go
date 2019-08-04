package flagvar

import (
	"flag"
	"fmt"
)

// Alternative tries to parse the argument using `Either`, and if that fails, using `Or`.
// `EitherOk` is true if the first attempt succeed.
type Alternative struct {
	Either   flag.Value
	Or       flag.Value
	EitherOk bool
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Alternative) Help() string {
	if fv.Either != nil && fv.Or != nil {
		if eitherHelp, ok := fv.Either.(interface {
			Help() string
		}); ok {
			if orHelp, ok := fv.Or.(interface {
				Help() string
			}); ok {
				return fmt.Sprintf("either %s, or %s", eitherHelp.Help(), orHelp.Help())
			}
		}
	}
	return ""
}

// Set is flag.Value.Set
func (fv *Alternative) Set(v string) error {
	err := fv.Either.Set(v)
	fv.EitherOk = err == nil
	if err != nil {
		return fv.Or.Set(v)
	}
	return nil
}

func (fv *Alternative) String() string {
	if fv.EitherOk {
		return fv.Either.String()
	}
	if fv.Or != nil {
		return fv.Or.String()
	}
	return ""
}
