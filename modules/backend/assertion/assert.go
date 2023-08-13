package assertion

type AssertionFunc func() error

func Assert(funcList ...AssertionFunc) error {
	for _, f := range funcList {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
