package runstats

import "os"
import "fmt"
import "flag"
import "time"
import "runtime"

import "github.com/bmhatfield/g2s"

var pause *int = flag.Int("pause", 10, "Collection pause interval")
var CPU *bool = flag.Bool("cpu", true, "Collect CPU Statistics")
var MEM *bool = flag.Bool("mem", true, "Collect Memory Statistics")
var GC *bool = flag.Bool("gc", true, "Collect GC Statistics (requires Memory be enabled)")

var prefix string
var s g2s.Statter

var err error

func init() {
	s, err = g2s.Dial("udp", "localhost:8125")

	if err != nil {
		panic("Unable to connect to Statsd")
	}

	flag.Parse()

	prefix, err = os.Hostname()

	if err != nil {
		prefix = "unknown_host.go"
	} else {
		prefix = fmt.Sprintf("%s.go", prefix)
	}

	go collector()
}

func collector() {
	if *CPU {
		// Goroutines
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "cpu.goroutines"), string(runtime.NumGoroutine()))

		// CGo calls
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "cpu.cgo_calls"), string(runtime.NumCgoCall()))
	}

	if *MEM {
		m := new(runtime.MemStats)
		runtime.ReadMemStats(m)

		// General
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.alloc"), string(m.Alloc))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.total"), string(m.TotalAlloc))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.sys"), string(m.Sys))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.lookups"), string(m.Lookups))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.malloc"), string(m.Mallocs))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.frees"), string(m.Frees))

		// Heap
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.heap.alloc"), string(m.HeapAlloc))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.heap.sys"), string(m.HeapSys))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.heap.idle"), string(m.HeapIdle))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.heap.inuse"), string(m.HeapInuse))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.heap.released"), string(m.HeapReleased))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.heap.objects"), string(m.HeapObjects))

		// Stack
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.stack.inuse"), string(m.StackInuse))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.stack.sys"), string(m.StackSys))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.stack.mspan_inuse"), string(m.MSpanInuse))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.stack.mspan_sys"), string(m.MSpanSys))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.stack.mcache_inuse"), string(m.MCacheInuse))
		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.stack.mcache_sys"), string(m.MCacheSys))

		s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.othersys"), string(m.OtherSys))

		if *GC {
			s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.gc.sys"), string(m.GCSys))
			s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.gc.next"), string(m.NextGC))
			s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.gc.last"), string(m.LastGC))
			s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.gc.pause_total"), string(m.PauseTotalNs))
			s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.gc.pause"), string(m.PauseNs[(m.NumGC+255)%256]))
			s.Gauge(1.0, fmt.Sprintf("%s.%s", prefix, "mem.gc.count"), string(m.NumGC))
		}
	}

	// Don't collect for a time interval
	time.Sleep(time.Duration(*pause) * time.Second)
}
