package terrgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Group struct {
	// not fully implemented,currently just limits channel buffer size
	// defaults to 10
	MaxThreads *int
	g          *errgroup.Group
	ctx        context.Context
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	return &Group{
		g:   g,
		ctx: ctx,
	}, ctx
}

// generally, a transformer is an array or a slice.
// that will asynchronously transform its elements to be held by an appender
type Transformer interface {
	Length() int
	Transform(int) (interface{}, error)
}

// in place appender
// be sure that the appender used to hold the transformed values can handle the resulting transformed values
// as there is no error handling in the receipt available except to not add the resulting element
type Appender interface {
	Append(interface{})
}

//
// in place appender that relies on pre-allocated memory
// affords significant performance boost
// be sure that the appender used to hold the transformed values can handle the resulting transformed values
// as there is no error handling in the receipt available except to not add the resulting element
type Injector interface {
	InjectAt(i int, value interface{})
}

// This function will not maintain order between transformer and appender
// and should be used when the cardinality of the non-error case is not known
// if the expected cardinality is known, use GoExactTransform instead
func (s *Group) GoTransform(input Transformer, output Appender) error {
	var g *errgroup.Group
	var ctx context.Context
	if s.g == nil || s.ctx == nil {
		g, ctx = errgroup.WithContext(context.Background())
	} else {
		g, ctx = s.g, s.ctx
	}
	var mt int
	if s.MaxThreads != nil {
		mt = *s.MaxThreads
	} else {
		//defaults to 10
		mt = 10
	}

	outputChan := make(chan interface{}, mt)
	//takes transform function and converts to errGroup friendly func
	fn := func(ctx context.Context, i int, inputTransformer Transformer, out chan<- interface{}) func() error {

		return func() error {
			v, err := inputTransformer.Transform(i)
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				out <- v
				return nil
			}
		}
	}
	for i := 0; i < input.Length(); i++ {
		g.Go(fn(ctx, i, input, outputChan))
	}
	// g.Wait() MUST be called after g.Go(), otherwise we have a race condition.
	errChan := make(chan error)
	go func() {
		errChan <- g.Wait()
	}()
	merge := func(outChan chan interface{}, errChan <-chan error, appender Appender) error {
		for {
			select {
			case err := <-errChan:
				if err != nil {
					return err
				} else {
					//avoid race condition
					close(outChan)
				}
			case result, more := <-outChan:
				if more {
					appender.Append(result)
				} else {
					return nil
				}
			}
		}
	}
	return merge(outputChan, errChan, output)
}

// NOT INTENDED TO MAINTAIN ORDER FROM input
// this function relies on pre-allocated memory receivers to boost performance
// and should be preferred when the cardinality of the expected result is known
func (s *Group) GoExactTransform(input Transformer, injector Injector) error {
	var g *errgroup.Group
	var ctx context.Context
	if s.g == nil {
		g, ctx = errgroup.WithContext(context.Background())
	} else {
		g, ctx = s.g, s.ctx
	}
	var mt int
	if s.MaxThreads != nil {
		mt = *s.MaxThreads
	} else {
		//defaults to 10
		mt = 10
	}

	outputChan := make(chan interface{}, mt)
	//takes transform function and converts to errGroup friendly func
	fn := func(ctx context.Context, i int, inputTransformer Transformer, out chan<- interface{}) func() error {

		return func() error {
			v, err := inputTransformer.Transform(i)
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				out <- v
				return nil
			}
		}
	}
	for i := 0; i < input.Length(); i++ {
		g.Go(fn(ctx, i, input, outputChan))
	}
	errChan := make(chan error)
	go func() {
		errChan <- g.Wait()
	}()
	merge := func(outChan chan interface{}, errChan <-chan error, injector Injector) error {
		var i int
		for {
			select {
			case err := <-errChan:
				if err != nil {
					return err
				} else {
					//avoid race condition
					close(outChan)
				}
			case result, more := <-outChan:
				if more {
					injector.InjectAt(i, result)
					i++
				} else {
					return nil
				}
			}
		}
	}
	return merge(outputChan, errChan, injector)
}
