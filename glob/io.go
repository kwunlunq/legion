package glob

import (
	"context"
	"errors"
	"io"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"gitlab.paradise-soft.com.tw/glob/tracer"
	"golang.org/x/time/rate"
)

func IoBind(dst io.ReadWriter, srcs []io.ReadWriter, fn func(isSrcErr bool, err error), cfn func(count int,
	isPositive bool), bytesPreSec float64) {
	var one = &sync.Once{}
	// var srcReader []io.Reader
	// for _, src := range srcs {
	// 	srcReader = append(srcReader, src)
	// }
	go func() {
		defer func() {
			if e := recover(); e != nil {
				tracer.Errorf("testrp", "IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
			}
		}()
		var err error
		var isSrcErr bool
		_, isSrcErr, err = ioCopy("1", dst, srcs, func(src io.ReadWriter) {
			go func() {
				defer func() {
					if e := recover(); e != nil {
						tracer.Errorf("testrp", "IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
					}
				}()
				var err error
				var isSrcErr bool
				// fmt.Println("9")

				_, isSrcErr, err = ioCopy2("2", src, dst)
				// fmt.Println("8")
				// if bytesPreSec > 0 {
				// 	newReader := NewReader(dst)
				// 	newReader.SetRateLimit(bytesPreSec)
				// 	_, isSrcErr, err = ioCopy("2", src, []io.Reader{newReader}, func(c int) {
				// 		cfn(c, true)
				// 	})
				// } else {
				//
				//
				// }
				if err != nil {
					one.Do(func() {
						fn(isSrcErr, err)
					})
				}
			}()
		}, func(c int) {
			cfn(c, false)
		})
		// fmt.Println("7")
		// if bytesPreSec > 0 {
		// 	newreader := NewReader(src)
		// 	newreader.SetRateLimit(bytesPreSec)
		// 	_, isSrcErr, err = ioCopy("1", dst, []io.Reader{newreader}, func(c int) {
		// 		cfn(c, false)
		// 	})
		//
		// } else {
		//
		// }
		if err != nil {
			one.Do(func() {
				fn(isSrcErr, err)
			})
		}
	}()
	// go func() {
	// 	defer func() {
	// 		if e := recover(); e != nil {
	// 			tracer.Errorf("testrp", "IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
	// 		}
	// 	}()
	// 	var err error
	// 	var isSrcErr bool
	// 	_, isSrcErr, err = ioCopy("2", src, []io.Reader{dst}, func(c int) {
	// 		cfn(c, true)
	// 	})
	// 	fmt.Println("8")
	// 	// if bytesPreSec > 0 {
	// 	// 	newReader := NewReader(dst)
	// 	// 	newReader.SetRateLimit(bytesPreSec)
	// 	// 	_, isSrcErr, err = ioCopy("2", src, []io.Reader{newReader}, func(c int) {
	// 	// 		cfn(c, true)
	// 	// 	})
	// 	// } else {
	// 	//
	// 	//
	// 	// }
	// 	if err != nil {
	// 		one.Do(func() {
	// 			fn(isSrcErr, err)
	// 		})
	// 	}
	// }()
}

func ioCopy(step string, dst io.ReadWriter, srcs []io.ReadWriter, sfn func(io.ReadWriter), fn ...func(count int)) (written int64,
	isSrcErr bool,
	err error) {
	// defer func() {
	// 	fmt.Println(step, " end")
	// }()

	isSuccess := false
	if len(srcs) > 0 {
		for i, src := range srcs {
			buf := make([]byte, 32*1024)
			isFirst := true
			for {
				nr, er := src.Read(buf)
				// fmt.Println(step)
				if isFirst && isSuccess {
					break
				}

				if isFirst && (strings.HasPrefix(string(buf), "HTTP/1.1 200") || strings.HasPrefix(string(buf), "HTTP/1.0 200")) {
					isFirst = false
					// fmt.Println(isFirst)
					sfn(src)
				} else if isFirst {
					// fmt.Println(string(buf))
					if i == len(srcs)-1 {
						isSrcErr = true
						err = errors.New("no proxy success")
					}
					break
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

func ioCopy_bak(step string, dst io.ReadWriter, srcs []io.ReadWriter, sfn func(io.ReadWriter),
	fn ...func(count int)) (written int64,
	isSrcErr bool,
	err error) {
	// defer func() {
	// 	fmt.Println(step, " end")
	// }()

	isSuccess := false
	_ = isSuccess
	var lock sync.Mutex
	var wait sync.WaitGroup
	if len(srcs) > 0 {
		for _, src := range srcs {
			wait.Add(1)
			go func() {
				defer func() {
					wait.Done()
				}()
				buf := make([]byte, 32*1024)
				isFirst := true
				for {
					nr, er := src.Read(buf)
					// fmt.Println(step)
					// fmt.Println(string(buf))
					lock.Lock()

					if isFirst && isSuccess {
						lock.Unlock()
						break
					}

					if isFirst && (strings.HasPrefix(string(buf), "HTTP/1.1 200") || strings.HasPrefix(string(buf), "HTTP/1.0 200")) {
						isFirst = false
						// fmt.Println(isFirst)
						sfn(src)
					}

					isSuccess = true

					lock.Unlock()

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
			}()

			// if isSuccess{
			// 	break
			// }
		}

	}

	wait.Wait()

	return written, isSrcErr, err
}

func ioCopy2(step string, dst io.Writer, src io.Reader) (written int64, isSrcErr bool,
	err error) {
	// defer func() {
	// 	fmt.Println(step, " end")
	// }()

	buf := make([]byte, 32*1024)

	for {
		nr, er := src.Read(buf)
		// fmt.Println(step)

		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])

			if nw > 0 {
				written += int64(nw)

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
