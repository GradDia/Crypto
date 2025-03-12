package cases

type Options struct {
	Titles []string
}

type Option func(*Options)

func WithTitles(titles []string) Option {
	return func(o *Options) {
		o.Titles = titles
	}

}
