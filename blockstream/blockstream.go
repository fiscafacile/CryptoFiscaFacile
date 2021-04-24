package blockstream

type Blockstream struct {
	apiTXs []apiTX
	done   chan error
}

func New() *Blockstream {
	blkst := &Blockstream{}
	blkst.done = make(chan error)
	return blkst
}

func (blkst *Blockstream) WaitFinish() error {
	return <-blkst.done
}
