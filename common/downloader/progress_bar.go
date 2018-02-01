package downloader

import (
	"io"

	pb "gopkg.in/cheggaaa/pb.v1"
)

// ProgressBar is an implementation of ProxyReader that displays a process bar
// during file downloading.
type ProgressBar struct {
	bar *pb.ProgressBar
	out io.Writer
}

func NewProgressBar(out io.Writer) *ProgressBar {
	return &ProgressBar{out: out}
}

func (p *ProgressBar) Proxy(totalSize int64, reader io.Reader) io.Reader {
	p.bar = pb.New(int(totalSize)).SetUnits(pb.U_BYTES)
	p.bar.Output = p.out
	p.bar.Start()

	return p.bar.NewProxyReader(reader)
}

func (p *ProgressBar) Finish() {
	p.bar.Finish()
}
