
// In the range [0, 2^24-1] there are 1,077,871 primes

package main

import (
	"fmt"
	"math"
	"sort"
	"runtime"
	"reflect"
	"time"
)

func checkIfPrime(val uint64) bool {
	if val < 2 {
		return false
	}
	var ISPRIME bool = true
    for i := uint64(2); i <= uint64(math.Sqrt(float64(val))); i++ {
    	if val%i == 0 {
   			ISPRIME = false
   			break
    	}
    }
    return ISPRIME
}

func findprimes(cid int, c chan uint64, min_prime, max_prime uint64) {
    fmt.Printf("{%d} --> checking (0x%08X, 0x%08X]\n", cid, min_prime, max_prime)
    for j := min_prime; j <= max_prime; j++ {
	    if checkIfPrime(j) {
	   		c <- j
	    }
    }
    close(c)
}

func main() {
	timeStart := time.Now()

	fmt.Printf("NumCPU = %d, GOMAXPROCS = %d %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(8), runtime.GOMAXPROCS(8))
	numWorkers := 8
    var workers []chan uint64 = make([](chan uint64), numWorkers)
	for i := range workers {
	   workers[i] = make(chan uint64)
	}
	
	var mymin uint64 = 0
	var mymax uint64 = 0x00FFFFFF

	var myrange uint64 = (mymax-mymin+1)/uint64(numWorkers)
	for i := 0; i < numWorkers; i++ {
		min := myrange*uint64(i)
		max := myrange*(uint64(i)+1)-1
		fmt.Printf("Dispatching %d  0x%08X 0x%08X\n", uint64(i), min, max)
		go findprimes(i, workers[i], min, max)
	}
	
	var ok bool
	var msg reflect.Value
	var chosen int
	var primes = make(map[uint64]bool)

	cases := make([]reflect.SelectCase, len(workers))
	for i, ch := range workers {
	    cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}

	remaining := len(cases)
	for remaining > 0 {
		chosen, msg, ok = reflect.Select(cases)
		if ok {
			primes[msg.Uint()] = true
		} else {
			cases[chosen].Chan = reflect.ValueOf(nil) // Channel closed
			fmt.Printf("{%d} --> DONE!\n", chosen)
			remaining -= 1
			continue
		}
		//fmt.Printf("Read from channel %#v and received %d\n", workers[chosen], msg.Uint())
	}

	fmt.Printf("Sorting output\n")
	var keys []int
	for k := range primes {
	    keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, k := range keys {
	    //fmt.Printf("Key: %d  Value: %d", k, primes[k])
	    //fmt.Printf("%d\n", k)
	    _=k
	}

	timeEnd := time.Now()
	timeDelta := timeEnd.Sub(timeStart)
	fmt.Printf("Found %d primes in %s\n", len(keys), timeDelta)
}

