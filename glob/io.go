package glob

import (
	"context"
	"io"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"gitlab.paradise-soft.com.tw/glob/tracer"
	"golang.org/x/time/rate"
)

func IoBind(dst io.ReadWriter, src io.ReadWriter, fn func(isSrcErr bool, err error), cfn func(count int, isPositive bool), bytesPreSec float64) {
	var one = &sync.Once{}
	go func() {
		defer func() {
			if e := recover(); e != nil {
				tracer.Errorf("testrp", "IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
			}
		}()
		var err error
		var isSrcErr bool
		if bytesPreSec > 0 {
			newreader := NewReader(src)
			newreader.SetRateLimit(bytesPreSec)
			_, isSrcErr, err = ioCopy("1", dst, []io.Reader{newreader}, func(c int) {
				cfn(c, false)
			})

		} else {
			//
			_, isSrcErr, err = ioCopy("1", dst, []io.Reader{src}, func(c int) {
				cfn(c, false)
			})
		}
		if err != nil {
			one.Do(func() {
				fn(isSrcErr, err)
			})
		}
	}()
	go func() {
		defer func() {
			if e := recover(); e != nil {
				tracer.Errorf("testrp", "IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
			}
		}()
		var err error
		var isSrcErr bool
		if bytesPreSec > 0 {
			newReader := NewReader(dst)
			newReader.SetRateLimit(bytesPreSec)
			_, isSrcErr, err = ioCopy("2", src, []io.Reader{newReader}, func(c int) {
				cfn(c, true)
			})
		} else {
			_, isSrcErr, err = ioCopy("2", src, []io.Reader{dst}, func(c int) {
				cfn(c, true)
			})
			// fmt.Println("8")

		}
		if err != nil {
			one.Do(func() {
				fn(isSrcErr, err)
			})
		}
	}()
}

func ioCopy(step string, dst io.Writer, srcs []io.Reader, fn ...func(count int)) (written int64, isSrcErr bool,
	err error) {
	// defer func() {
	// 	fmt.Println("6")
	// }()
	buf := make([]byte, 32*1024)
	isFirst := true
	_ = isFirst
	isSuccess := false
	_ = isSuccess

	if len(srcs) > 0 {
		for _, src := range srcs {
			for {
				nr, er := src.Read(buf)
				if isFirst && strings.HasPrefix(string(buf), "HTTP/1.1 200 Connection established") {
					isFirst = false
				}

				isSuccess = true
				if nr > 0 {
					nw, ew := dst.Write(buf[0:nr])

					if nw > 0 {
						written += int64(nw)
						if len(fn) == 1 {
							fn[0](nw)
						}
					}
					if ew != nil {
						err = ew
						break
					}
					if nr != nw {
						err = io.ErrShortWrite
						break
					}
				}
				if er != nil {
					err = er
					isSrcErr = true
					break
				}
			}
			// if isSuccess{
			// 	break
			// }
		}

	}

	return written, isSrcErr, err
}

type Reader struct {
	r       io.Reader
	limiter *rate.Limiter
	ctx     context.Context
}

// NewReader returns a reader that implements io.Reader with rate limiting.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:   r,
		ctx: context.Background(),
	}
}

const burstLimit = 1000 * 1000 * 1000

// SetRateLimit sets rate limit (bytes/sec) to the reader.
func (s *Reader) SetRateLimit(bytesPerSec float64) {
	s.limiter = rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	s.limiter.AllowN(time.Now(), burstLimit) // spend initial burst
}

// Read reads bytes into p.
func (s *Reader) Read(p []byte) (int, error) {
	if s.limiter == nil {
		return s.r.Read(p)
	}
	n, err := s.r.Read(p)
	if err != nil {
		return n, err
	}
	if err := s.limiter.WaitN(s.ctx, n); err != nil {
		return n, err
	}
	return n, nil
}
