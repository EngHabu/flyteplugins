package data

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSyncFuture(t *testing.T) {
	type args struct {
		val interface{}
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{"val", args{val: "val"}},
		{"nil-val", args{}},
		{"error", args{err: fmt.Errorf("err")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSyncFuture(tt.args.val, tt.args.err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.args.val, got.val)
			assert.Equal(t, tt.args.err, got.err)
			assert.True(t, got.Ready())
			v, err := got.Get(nil)
			assert.Equal(t, tt.args.val, v)
			assert.Equal(t, tt.args.err, err)
		})
	}
}

func TestAsyncFuture(t *testing.T) {

	t.Run("immediate-return-val", func(t *testing.T) {
		v := "val"
		err := fmt.Errorf("err")
		af := NewAsyncFuture(context.TODO(), func(ctx context.Context) (interface{}, error) {
			return v, err
		})
		assert.NotNil(t, af)
		rv, rerr := af.Get(context.TODO())
		assert.Equal(t, v, rv)
		assert.Equal(t, err, rerr)
		assert.True(t, af.Ready())
	})

	t.Run("wait-return-val", func(t *testing.T) {
		v := "val"
		err := fmt.Errorf("err")
		af := NewAsyncFuture(context.TODO(), func(ctx context.Context) (interface{}, error) {
			time.Sleep(time.Second * 1)
			return v, err
		})
		runtime.Gosched()
		assert.NotNil(t, af)
		rv, rerr := af.Get(context.TODO())
		assert.Equal(t, v, rv)
		assert.Equal(t, err, rerr)
		assert.True(t, af.Ready())
	})

	t.Run("timeout", func(t *testing.T) {
		v := "val"
		ctx := context.TODO()
		af := NewAsyncFuture(ctx, func(ctx context.Context) (interface{}, error) {
			time.Sleep(time.Second * 5)
			return v, nil
		})
		runtime.Gosched()
		cctx, cancel := context.WithCancel(ctx)
		cancel()
		_, rerr := af.Get(cctx)
		assert.Error(t, rerr)
		assert.Equal(t, AsyncFutureCanceledErr, rerr)
	})
}
