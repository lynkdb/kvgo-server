package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/lynkdb/kvgo"
	kv2 "github.com/lynkdb/kvspec/v2/go/kvspec"
)

type baseCommand = cobra.Command

var (
	version = "0.9.0"
	release = ""
)

var rootCmd = &baseCommand{
	Use:   "kvgo-convert",
	Short: "Storage Engine Convert Tool",
}

func main() {

	rootCmd.AddCommand(newPebbleToGoleveldbCommand())

	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}

	log.Println("DONE")
}

type goleveldbToPebble struct {
	cmd          *baseCommand
	argInputDir  string
	argOutputDir string
}

func newPebbleToGoleveldbCommand() *baseCommand {

	c := &goleveldbToPebble{
		cmd: &baseCommand{
			Use:   "goleveldb-to-pebble",
			Short: "Convert Goleveldb to Pebble",
		},
	}

	c.cmd.Flags().StringVarP(&c.argInputDir, "input", "i",
		"", "source dir")

	c.cmd.Flags().StringVarP(&c.argOutputDir, "output", "o",
		"", "destination dir")

	c.cmd.RunE = c.run

	return c.cmd
}

func (it *goleveldbToPebble) run(cmd *baseCommand, args []string) error {

	if it.argInputDir == "" {
		return fmt.Errorf("--input-dir not setup")
	}

	if it.argOutputDir == "" {
		return fmt.Errorf("--output-dir not setup")
	}

	if it.argInputDir == it.argOutputDir {
		return fmt.Errorf("invalid dir path")
	}

	src, err := kvgo.StorageLevelDBOpen(it.argInputDir, &kv2.StorageOptions{})
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := kvgo.StoragePebbleOpen(it.argOutputDir, nil)
	if err != nil {
		return err
	}
	defer dst.Close()

	type keyValue struct {
		key   []byte
		value []byte
	}

	var (
		queue       = make(chan *keyValue, 1000)
		quit        = false
		num   int64 = 0
		t1          = time.Now()
	)

	go func() {

		var (
			rg = &kv2.StorageIteratorRange{
				Start: []byte{},
				Limit: []byte{0xff, 0xff},
			}
			iter = src.NewIterator(rg)
		)
		defer iter.Release()

		for ok := iter.First(); ok && !quit; ok = iter.Next() {
			queue <- &keyValue{
				key:   kvgo.BytesClone(iter.Key()),
				value: kvgo.BytesClone(iter.Value()),
			}
		}

		queue <- nil
	}()

	go func() {
		tr := time.NewTicker(1e9)
		defer tr.Stop()

		for !quit {
			_ = <-tr.C
			log.Printf("db put %d keys, time %v", num, time.Since(t1))
		}
	}()

	for {
		item := <-queue
		if item == nil {
			break
		}

		if ss := dst.Put(item.key, item.value, nil); ss.OK() {
			num += 1
		} else {
			log.Printf("put fail %v", ss.Error())
			break
		}
	}

	quit = true
	time.Sleep(2e9)

	return nil
}
