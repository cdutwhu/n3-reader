package n3reader

import (
	"context"
	"fmt"
	"log"
	"os"

	jt "github.com/digisan/json-tool"
)

func (n4r *Nats4Reader) PubAsJSON(fileName, meta string) error {
	fmt.Println("Publishing:", fileName)

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cOut, jarr := jt.ScanObject(ctx, f, false, true, jt.OUT_ORI)
	if !jarr {
		log.Println("not json array")
	}
	for result := range cOut {
		if result.Err != nil {
			log.Println(result.Err)
			return result.Err
		}

		// merge file watcher meta & nats reader meta
		meta = jt.MergeSgl(meta, n4r.meta())

		// make publish data
		data := jt.Minimize(fmt.Sprintf(`{"meta":%s, "data":%s}`, meta, result.Obj), true)

		// publish action
		ack, err := n4r.js.Publish(n4r.subject, []byte(data))
		if err != nil {
			return err
		}
		fmt.Println("ACK:", ack.Stream)
	}

	return nil
}
