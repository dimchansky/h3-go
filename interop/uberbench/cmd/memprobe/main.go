// Command memprobe measures process-level memory for one workload in one
// implementation and prints a single TSV line. It exists because Go
// benchmark metrics (B/op, allocs/op) only see the Go heap: memory that C
// code allocates through malloc — as the cgo binding does for polygon
// inputs and linked geo structures — is invisible to them. Peak RSS is the
// only comparable, implementation-agnostic ceiling.
//
// One process measures exactly one (implementation, workload) pair so that
// peak RSS is attributable; run the full matrix with the bench-uber-mem
// make target (see interop/uberbench/run.sh).
//
// Usage:
//
//	go run ./cmd/memprobe -impl pure|uber -workload <name> [-iters N]
//	go run ./cmd/memprobe -list
//	go run ./cmd/memprobe -header
//
// Output columns (TSV):
//
//	impl workload iters wall_ms peak_rss_kb heap_alloc_kb total_alloc_kb
//	go_sys_kb mallocs num_gc checksum
//
// peak_rss_kb is the high-water resident set size of the whole process
// (getrusage RU_MAXRSS), which includes Go heap, Go runtime, and any C
// heap. heap_alloc_kb/total_alloc_kb/go_sys_kb/mallocs/num_gc come from
// runtime.ReadMemStats after a final runtime.GC() with the workload result
// still reachable, so they describe retained Go heap — and, by their
// difference from peak RSS, how much the Go accounting misses.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/dimchansky/h3-go/interop/uberbench"
)

func main() {
	impl := flag.String("impl", "", "implementation to measure: pure or uber")
	workload := flag.String("workload", "", "workload name (see -list)")
	iters := flag.Int("iters", 3, "workload iterations")
	list := flag.Bool("list", false, "list workloads and exit")
	header := flag.Bool("header", false, "print the TSV header and exit")
	flag.Parse()

	if *list {
		for _, w := range uberbench.MemWorkloads {
			fmt.Printf("%-22s %s\n", w.Name, w.Description)
		}
		return
	}
	if *header {
		fmt.Println("impl\tworkload\titers\twall_ms\tpeak_rss_kb\theap_alloc_kb\ttotal_alloc_kb\tgo_sys_kb\tmallocs\tnum_gc\tchecksum")
		return
	}

	w, err := uberbench.MemWorkloadByName(*workload)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	var fn func(int) uberbench.MemResult
	switch *impl {
	case "pure":
		fn = w.Pure
	case "uber":
		fn = w.Uber
	default:
		fmt.Fprintf(os.Stderr, "unknown -impl %q (want pure or uber)\n", *impl)
		os.Exit(2)
	}

	start := time.Now()
	res := fn(*iters)
	wall := time.Since(start)

	// Retained-state measurement: collect garbage with the workload result
	// still reachable, then read both the Go view and the OS view.
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	peak, err := peakRSSKb()
	if err != nil {
		fmt.Fprintln(os.Stderr, "getrusage:", err)
		os.Exit(1)
	}

	fmt.Printf("%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%#x\n",
		*impl, *workload, *iters, wall.Milliseconds(), peak,
		ms.HeapAlloc/1024, ms.TotalAlloc/1024, ms.Sys/1024,
		ms.Mallocs, ms.NumGC, res.Checksum)

	runtime.KeepAlive(res.Retained)
}

// peakRSSKb returns the process high-water RSS in KiB. getrusage reports
// bytes on darwin and KiB on linux.
func peakRSSKb() (int64, error) {
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		return 0, err
	}
	if runtime.GOOS == "darwin" {
		return ru.Maxrss / 1024, nil
	}
	return ru.Maxrss, nil
}
