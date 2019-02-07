package legacy

import (
	"bytes"
	"context"
	"io"
	"testing"

	cmds "gx/ipfs/QmSXUokcP4TJpFfqozT69AVAYRtzXVMUjzQVkYX41R9Svs/go-ipfs-cmds"
	oldcmds "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/commands"
	cmdkit "gx/ipfs/Qmde5VP1qUkyQXKCfmEUA7bP64V2HAptbJ7phuPp7jXWwg/go-ipfs-cmdkit"
)

type WriteNopCloser struct {
	io.Writer
}

func (wc WriteNopCloser) Close() error {
	return nil
}

func TestNewCommand(t *testing.T) {
	root := &cmds.Command{
		Subcommands: map[string]*cmds.Command{
			"test": NewCommand(&oldcmds.Command{
				Run: func(req oldcmds.Request, res oldcmds.Response) {
					res.SetOutput("Test.")
				},
				Marshalers: map[oldcmds.EncodingType]oldcmds.Marshaler{
					oldcmds.Text: func(res oldcmds.Response) (io.Reader, error) {
						ch, ok := res.Output().(<-chan interface{})
						if !ok {
							t.Fatalf("output is not <-chan interface{} but %T", ch)
						}

						v := <-ch
						str, ok := v.(string)
						if !ok {
							t.Fatalf("read value is not string but %T", v)
						}

						buf := bytes.NewBuffer(nil)
						_, err := io.WriteString(buf, str)
						if err != nil {
							t.Fatal(err)
						}

						return buf, nil
					},
				},
				Subcommands: map[string]*oldcmds.Command{
					"sub": &oldcmds.Command{
						Options: []cmdkit.Option{
							cmdkit.NewOption(cmdkit.String, "test", "t", "some random test flag"),
						},
					},
				},
			}),
		},
	}

	path := []string{"test"}
	req, err := cmds.NewRequest(context.TODO(), path, nil, nil, nil, root)
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer(nil)

	// test calling "test" command
	testCmd := root.Subcommands["test"]

	re, err := cmds.NewWriterResponseEmitter(WriteNopCloser{buf}, req)
	if err != nil {
		t.Fatal(err)
	}

	var env oldcmds.Context

	root.Call(req, re, &env)

	expected := `"Test."
`

	if buf.String() != expected {
		t.Fatalf("expected string %#v but got %#v", expected, buf.String())
	}

	// test getting subcommand
	subCmd := testCmd.Subcommands["sub"]
	if subCmd == nil {
		t.Fatal("got nil subcommand")
	}

	if nOpts := len(subCmd.Options); nOpts != 1 {
		t.Fatalf("subcommand has %v options, expected 1", nOpts)
	}

	opt := subCmd.Options[0]

	if nNames := len(opt.Names()); nNames != 2 {
		t.Fatalf("option has %v names, expected 2", nNames)
	}

	names := opt.Names()
	if names[0] != "test" {
		t.Fatalf("option has name %q, expected %q", names[0], "test")
	}

	if names[1] != "t" {
		t.Fatalf("option has name %q, expected %q", names[1], "t")
	}
}

func TestPipePair(t *testing.T) {
	cmd := NewCommand(&oldcmds.Command{Type: "string"})

	req, err := cmds.NewRequest(context.TODO(), nil, nil, nil, nil, cmd)
	if err != nil {
		t.Fatal(err)
	}

	r, w := io.Pipe()
	re, err := cmds.NewWriterResponseEmitter(w, req)
	if err != nil {
		t.Fatal(err)
	}
	res, err := cmds.NewReaderResponse(r, req)
	if err != nil {
		t.Fatal(err)
	}

	wait := make(chan interface{})

	expect := "abc"
	go func() {
		err := re.Emit(expect)
		if err != nil {
			t.Fatal(err)
		}

		err = re.Close()
		if err != nil {
			t.Fatal(err)
		}

		close(wait)
	}()

	v, err := res.Next()
	if err != nil {
		t.Fatal(err)
	}
	str, ok := v.(*string)
	if !ok {
		t.Fatalf("expected type %T but got %T", expect, v)
	}
	if *str != expect {
		t.Fatalf("expected value %#v but got %#v", expect, v)
	}

	_, err = res.Next()
	if err != io.EOF {
		t.Fatal("expected io.EOF, got:", err)
	}

	<-wait
}

func TestChanPair(t *testing.T) {
	cmd := NewCommand(&oldcmds.Command{Type: "string"})

	req, err := cmds.NewRequest(context.TODO(), nil, nil, nil, nil, cmd)
	if err != nil {
		t.Fatal(err)
	}

	re, res := cmds.NewChanResponsePair(req)

	wait := make(chan interface{})

	expect := "abc"
	go func() {
		err := re.Emit(expect)
		if err != nil {
			t.Fatal(err)
		}

		err = re.Close()
		if err != nil {
			t.Fatal(err)
		}

		close(wait)
	}()

	v, err := res.Next()
	if err != nil {
		t.Fatal(err)
	}
	str, ok := v.(string)
	if !ok {
		t.Fatalf("expected type %T but got %T", expect, v)
	}
	if str != expect {
		t.Fatalf("expected value %#v but got %#v", expect, v)
	}

	_, err = res.Next()
	if err != io.EOF {
		t.Fatal("expected io.EOF, got:", err)
	}

	<-wait
}
