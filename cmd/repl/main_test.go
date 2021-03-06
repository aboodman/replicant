package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"
	"github.com/stretchr/testify/assert"

	"github.com/aboodman/replicant/db"
	"github.com/aboodman/replicant/util/time"
)

func TestCommands(t *testing.T) {
	assert := assert.New(t)
	defer time.SetFake()()

	tc := []struct {
		label string
		in    string
		args  string
		code  int
		out   string
		err   string
	}{
		{
			"log empty",
			"",
			"log --no-pager",
			0,
			"",
			"",
		},
		{
			"bundle put good",
			"function futz(k, v){ db.put(k, v) }; function echo(v) { return v; };",
			"bundle put",
			0,
			"",
			"Replacing unversioned bundle 2eulo8v8rihcjm0e93brv14dopakkder with a2gj33rbobbebq5r89c2vmr8k2so3mo0\n",
		},
		{
			"bundle get good",
			"",
			"bundle get",
			0,
			"function futz(k, v){ db.put(k, v) }; function echo(v) { return v; };",
			"",
		},
		{
			"log bundle put",
			"",
			"log --no-pager",
			0,
			fmt.Sprintf("commit gedd9dnmlb2doi034ltrce8h18102n4m\nOrigin:      cli\nCreated:     %s\nStatus:      PENDING\nMerged:      %s\nTransaction: .putBundle(blob(a2gj33rbobbebq5r89c2vmr8k2so3mo0))\n\n", time.Now(), time.Now()),
			"",
		},
		{
			"exec unknown-function",
			"",
			"exec monkey",
			1,
			"",
			"Unknown function: monkey\n",
		},
		{
			"exec missing-key",
			"",
			"exec futz",
			1,
			"",
			"Error: Invalid id\n    at bootstrap.js:20:14\n    at bootstrap.js:26:4\n    at futz (bundle.js:1:22)\n    at apply (<native code>)\n    at recv (bootstrap.js:64:12)\n\n",
		},
		{
			"exec missing-val",
			"",
			"exec futz foo",
			1,
			"",
			"Error: Invalid value\n    at bootstrap.js:29:15\n    at futz (bundle.js:1:22)\n    at apply (<native code>)\n    at recv (bootstrap.js:64:12)\n\n",
		},
		{
			"exec good",
			"",
			"exec futz foo bar",
			0,
			"",
			"",
		},
		{
			"log exec good",
			"",
			"log --no-pager",
			0,
			fmt.Sprintf("commit gekh5qruoqaq5nk9atkcqmru26qmj220\nOrigin:      cli\nCreated:     %s\nStatus:      PENDING\nMerged:      %s\nTransaction: futz(\"foo\", \"bar\")\n(root) {\n+   \"foo\": \"bar\"\n  }\n\ncommit gedd9dnmlb2doi034ltrce8h18102n4m\nOrigin:      cli\nCreated:     %s\nStatus:      PENDING\nMerged:      %s\nTransaction: .putBundle(blob(a2gj33rbobbebq5r89c2vmr8k2so3mo0))\n\n", time.Now(), time.Now(), time.Now(), time.Now()),
			"",
		},
		{
			"exec echo",
			"",
			"exec echo monkey",
			0,
			`"monkey"`,
			"",
		},
		{
			"has missing-arg",
			"",
			"has",
			1,
			"",
			"required argument 'id' not provided\n",
		},
		{
			"has good",
			"",
			"has foo",
			0,
			"true\n",
			"",
		},
		{
			"get bad missing-arg",
			"",
			"get",
			1,
			"",
			"required argument 'id' not provided\n",
		},
		{
			"get good",
			"",
			"get foo",
			0,
			"\"bar\"\n",
			"",
		},
		{
			"del bad missing-arg",
			"",
			"del",
			1,
			"",
			"required argument 'id' not provided\n",
		},
		{
			"del good no-op",
			"",
			"del monkey",
			0,
			"No such id.\n",
			"",
		},
		{
			"del good",
			"",
			"del foo",
			0,
			"",
			"",
		},
	}

	td, err := ioutil.TempDir("", "")
	fmt.Println("test database:", td)
	assert.NoError(err)

	for _, c := range tc {
		ob := &strings.Builder{}
		eb := &strings.Builder{}
		code := 0
		args := append([]string{"--db=" + td}, strings.Split(c.args, " ")...)
		impl(args, strings.NewReader(c.in), ob, eb, func(c int) {
			code = c
		})

		assert.Equal(c.code, code, c.label)
		assert.Equal(c.out, ob.String(), c.label)
		assert.Equal(c.err, eb.String(), c.label)
	}
}

func TestDrop(t *testing.T) {
	assert := assert.New(t)
	tc := []struct {
		in      string
		errs    string
		deleted bool
	}{
		{"no\n", "", false},
		{"N\n", "", false},
		{"balls\n", "", false},
		{"n\n", "", false},
		{"windows\r\n", "", false},
		{"y\n", "", true},
		{"y\r\n", "", true},
	}

	for i, t := range tc {
		d, dir := db.LoadTempDB(assert)
		d.Put("foo", types.String("bar"))
		val, err := d.Get("foo")
		assert.NoError(err)
		assert.Equal("bar", string(val.(types.String)))

		desc := fmt.Sprintf("test case %d, input: %s", i, t.in)
		args := append([]string{"--db=" + dir, "drop"})
		out := strings.Builder{}
		errs := strings.Builder{}
		code := 0
		impl(args, strings.NewReader(t.in), &out, &errs, func(c int) { code = c })

		assert.Equal(dropWarning, out.String(), desc)
		assert.Equal(t.errs, errs.String(), desc)
		assert.Equal(0, code, desc)
		sp, err := spec.ForDatabase(dir)
		assert.NoError(err)
		noms := sp.GetDatabase()
		ds := noms.GetDataset(db.LOCAL_DATASET)
		assert.Equal(!t.deleted, ds.HasHead())
	}
}

func TestServe(t *testing.T) {
	assert := assert.New(t)
	_, dir := db.LoadTempDB(assert)
	args := append([]string{"--db=" + dir, "serve", "--port=8674"})
	go impl(args, strings.NewReader(""), os.Stdout, os.Stderr, func(_ int) {})

	sp, err := spec.ForDatabase("http://localhost:8674/serve/sandbox/foo")
	assert.NoError(err)
	d, err := db.New(sp.GetDatabase(), "test")
	assert.NoError(err)

	err = d.PutBundle(types.NewBlob(d.Noms(), strings.NewReader("function setFoo(val) { db.put('foo', val); }")))
	assert.NoError(err)
	_, err = d.Exec("setFoo", types.NewList(d.Noms(), types.String("bar")))
	assert.NoError(err)
	v, err := d.Get("foo")
	assert.NoError(err)
	assert.Equal("bar", string(v.(types.String)))
}

func TestEmptyInput(t *testing.T) {
	assert := assert.New(t)
	db.LoadTempDB(assert)
	var args []string

	// Just testing that they don't crash.
	// See https://github.com/aboodman/replicant/issues/120
	impl(args, strings.NewReader(""), ioutil.Discard, ioutil.Discard, func(_ int) {})
	args = []string{"--db=/tmp/foo"}
	impl(args, strings.NewReader(""), ioutil.Discard, ioutil.Discard, func(_ int) {})
}
