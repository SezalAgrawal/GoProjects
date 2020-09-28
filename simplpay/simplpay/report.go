package simplpay

type reporter interface {
	getAlldues()
}

type defaultReporter struct {
}

func (d *defaultReporter) getAlldues() {

}
