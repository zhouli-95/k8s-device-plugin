package cdi

type Option func(*CDIHandler)

func WithVendor(vendor string) Option {
	return func(c *CDIHandler) {
		c.Vendor = vendor
	}
}

func WithClass(class string) Option {
	return func(c *CDIHandler) {
		c.Class = class
	}
}
