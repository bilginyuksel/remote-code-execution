package codexec_test

import (
	"testing"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/stretchr/testify/assert"
)

func TestPoolAdd(t *testing.T) {
	pool := codexec.NewPool()

	pool.Add("cid-1")
	pool.Add("cid-2")
	pool.Add("cid-3")

	// head tail setup
	assert.Equal(t, "cid-1", pool.Head.ID)
	assert.Equal(t, "cid-3", pool.Tail.ID)

	// Tail 360 backward
	assert.Equal(t, "cid-2", pool.Tail.Prev.ID)
	assert.Equal(t, "cid-1", pool.Tail.Prev.Prev.ID)
	assert.Equal(t, "cid-3", pool.Tail.Prev.Prev.Prev.ID)

	// Tail 360 forward
	assert.Equal(t, "cid-1", pool.Tail.Next.ID)
	assert.Equal(t, "cid-2", pool.Tail.Next.Next.ID)
	assert.Equal(t, "cid-3", pool.Tail.Next.Next.Next.ID)

	// Head 360 forward
	assert.Equal(t, "cid-2", pool.Head.Next.ID)
	assert.Equal(t, "cid-3", pool.Head.Next.Next.ID)
	assert.Equal(t, "cid-1", pool.Head.Next.Next.Next.ID)

	// Head 360 backward
	assert.Equal(t, "cid-3", pool.Head.Prev.ID)
	assert.Equal(t, "cid-2", pool.Head.Prev.Prev.ID)
	assert.Equal(t, "cid-1", pool.Head.Prev.Prev.Prev.ID)

	// nodes are added to the Nodes in pool
	assert.Contains(t, pool.Nodes, "cid-1")
	assert.Contains(t, pool.Nodes, "cid-2")
	assert.Contains(t, pool.Nodes, "cid-3")
}

func TestPoolNext(t *testing.T) {
	pool := codexec.NewPool()

	// when pool is empty, return nil
	assert.Nil(t, pool.Get())

	pool.Add("cid-1")
	pool.Add("cid-2")
	pool.Add("cid-3")

	first := pool.Get().ID
	second := pool.Get().ID
	third := pool.Get().ID
	fourth := pool.Get().ID

	assert.Equal(t, "cid-1", first)
	assert.Equal(t, "cid-2", second)
	assert.Equal(t, "cid-3", third)
	assert.Equal(t, "cid-1", fourth)
}

func TestPoolRemove(t *testing.T) {

}
