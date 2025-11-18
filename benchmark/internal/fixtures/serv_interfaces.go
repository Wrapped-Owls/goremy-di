package fixtures

type (
	ServiceA interface {
		DoWork() string
	}

	ServiceB interface {
		Process() string
	}

	ServiceC interface {
		Execute() string
	}

	Repository interface {
		GetData() string
	}
)
