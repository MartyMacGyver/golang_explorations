
package main

// In the range [0, 2^24-1] there are 1,077,871 primes

import (
	"fmt"
	"runtime"
	"reflect"
	"time"
	"math"
	"sort"
)


type WorkerOutput struct {
    value uint64
    isPrime bool
    statsCycles uint64
    statsPasses uint64
}

func checkIfPrime(val uint64) (IsPrime bool, cycles uint64) {
	IsPrime = true
	cycles++
	if val < 2 {
		IsPrime = false
		return
	}
    for i := uint64(2); i <= uint64(math.Sqrt(float64(val))); i++ {
    	cycles++
    	if val%i == 0 {
   			IsPrime = false
   			break
    	}
    }
    return
}

func myWorker(cid int, cOut chan WorkerOutput, cIn chan uint64) {
	fmt.Printf("{%d} ++ myWorker initialized\n", cid)
    for {
        select {
            case value, ok := <- cIn:
                if ok {
				    //fmt.Printf("{%d} --> checking 0x%08X\n", cid, value)
				    statsPasses := uint64(0)
				    IsPrime, statsCycles := checkIfPrime(value)
				    statsPasses++
				    cOut <- WorkerOutput{value, IsPrime, statsCycles, statsPasses}
                } else {
               		fmt.Printf("{%d} --> cIn not OK!\n", cid)
                }
        }
    }
    close(cOut)
	fmt.Printf("{%d} xx myWorker exiting\n", cid)
}

func main() {
 	var valMin uint64 = 0
	var valMax uint64 = 0x00FFFFFF

	timeStart := time.Now()

	numWorkers := runtime.NumCPU()
	runtime.GOMAXPROCS(numWorkers)

	var gCycles []uint64
	var gPasses []uint64

	gCycles = make([]uint64, numWorkers)
	gPasses = make([]uint64, numWorkers)

    var downChans []chan WorkerOutput = make([](chan WorkerOutput), numWorkers)
	for i := range downChans {
		downChans[i] = make(chan WorkerOutput)
	}
	upChan := make(chan uint64)

	for i := 0; i < numWorkers; i++ {
		fmt.Printf("Creating worker %d\n", i)
		go myWorker(i, downChans[i], upChan)
	}

	cases := make([]reflect.SelectCase, numWorkers+1) // Add room for default
	for i, ch := range downChans {
	    cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	defaultCid := numWorkers
    cases[defaultCid] = reflect.SelectCase{Dir: reflect.SelectDefault}

	var primes = make(map[uint64]bool)

	poolValue := uint64(valMin)
	
	remaining := 0
	for i := 0; i < numWorkers && poolValue <= valMax; i++ {
		upChan <- poolValue // Primes the queue
		remaining++
		poolValue++
	}

	for remaining > 0 {
		var ok bool
		var msg reflect.Value
		var cid int
		//time.Sleep(time.Millisecond*00)
		cid, msg, ok = reflect.Select(cases)
		if cid == defaultCid {
			//fmt.Printf("{%d} <-- DEFAULT\n", cid)
			continue
		} else if ok {
			value       := msg.FieldByName("value").Uint()
			isPrime     := msg.FieldByName("isPrime").Bool()
			statsCycles := msg.FieldByName("statsCycles").Uint()
			statsPasses := msg.FieldByName("statsPasses").Uint()
		    gCycles[cid]+=statsCycles
		    gPasses[cid]+=statsPasses
			if isPrime {
				primes[value] = true
			}
			if (poolValue <= valMax) {
				upChan <- poolValue
				poolValue++
			} else {
				remaining--
				fmt.Printf("Remaining channels: %d\n", remaining)
			}
		} else {
			cases[cid].Chan = reflect.ValueOf(nil) // Channel closed
			fmt.Printf("{%d} --> DONE!\n", cid)
			remaining -= 1
			continue
		}
		//fmt.Printf("Read from channel %#v and received %d\n", downChans[cid], msg.Uint())
	}

	time.Sleep(time.Millisecond*1)

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
	totalCycles := uint64(0)
	totalPasses := uint64(0)
	for i := 0; i < numWorkers; i++ {
		totalCycles+=gCycles[i]
		totalPasses+=gPasses[i]
	}
	fmt.Printf("%d workers ran %12d cycles          for %11d values (%5.2f%% ideal)\n", numWorkers,
		totalCycles, totalPasses, 100.0/float64(numWorkers))
	for i := 0; i < numWorkers; i++ {
		fmt.Printf("Worker #%d ran %12d cycles (%5.2f%%) for %11d values (%5.2f%%)\n", i,
			gCycles[i], float64(gCycles[i])/float64(totalCycles)*100.0,
			gPasses[i], float64(gPasses[i])/float64(totalPasses)*100.0,
		)
	}
}
