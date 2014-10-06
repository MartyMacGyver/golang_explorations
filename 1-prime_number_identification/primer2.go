
package main

import (
	"fmt"
	"math"
	"sort"
	"runtime"
)

func findprimes(cid, min_prime, max_prime int, c []chan int) {
    println(fmt.Sprintf("{%d} --> checking (%d, %d]", cid, min_prime, max_prime))
    for j := min_prime; j < max_prime; j++ {
    	if j < 2 {
    		continue
    	}
	    ISPRIME := true
	    for i := 2; i <= int(math.Sqrt(float64(j))); i++ {
	    	if j%i == 0 {
	   			ISPRIME = false
	   			break
	    	}
	    }
	    if ISPRIME {
	   		c[cid] <- j
	    }
    }
    close(c[cid])
}

func main() {
	num_gos := 8
	println(runtime.NumCPU(), runtime.GOMAXPROCS(8))
    var c []chan int = make([](chan int), num_gos)
	for i := range c {
	   c[i] = make(chan int)
	}
	
	go findprimes(0, 0x00000000, 0x001FFFFF, c[0:])
	go findprimes(1, 0x00200000, 0x003FFFFF, c[0:])
	go findprimes(2, 0x00400000, 0x005FFFFF, c[0:])
	go findprimes(3, 0x00600000, 0x007FFFFF, c[0:])
	go findprimes(4, 0x00800000, 0x009FFFFF, c[0:])
	go findprimes(5, 0x00A00000, 0x00BFFFFF, c[0:])
	go findprimes(6, 0x00C00000, 0x00DFFFFF, c[0:])
	go findprimes(7, 0x00E00000, 0x00FFFFFF, c[0:])
	
	var ok bool
	var msg int
	var primes = make(map[int]bool)
    isRunning := true
    for isRunning {
        select {
            case msg, ok = <-c[0]:
                if ok {
			        primes[msg] = true
                } else {
                	println("{0} --> DONE!")
                	c[0] = nil
                }
            case msg, ok = <-c[1]:
                if ok {
			        primes[msg] = true
                } else {
                	println("{1} --> DONE!")
                	c[1] = nil
                }
            case msg, ok = <-c[2]:
                if ok {
			        primes[msg] = true
                } else {
                	println("{2} --> DONE!")
                	c[2] = nil
                }
            case msg, ok = <-c[3]:
                if ok {
			        primes[msg] = true
                } else {
                	println("{3} --> DONE!")
                	c[3] = nil
                }
            case msg, ok = <-c[4]:
                if ok {
			        primes[msg] = true
                } else {
                	println("{4} --> DONE!")
                	c[4] = nil
                }
            case msg, ok = <-c[5]:
                if ok {
			        primes[msg] = true
                } else {
                	println("{5} --> DONE!")
                	c[5] = nil
                }
            case msg, ok = <-c[6]:
                if ok {
			        primes[msg] = true
                } else {
                	println("{6} --> DONE!")
                	c[6] = nil
                }
            case msg, ok = <-c[7]:
                if ok {
			        primes[msg] = true
                } else {
                	println("{7} --> DONE!")
                	c[7] = nil
                }
        }
        if c[0] == nil && c[1] == nil && c[2] == nil && c[3] == nil &&
           c[4] == nil && c[5] == nil && c[6] == nil && c[7] == nil {
        	isRunning = false
        }
	}

	println("Sorting output")
	var keys []int
	for k := range primes {
	    keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
	    //fmt.Println("Key:", k, "Value:", primes[k])
	    //fmt.Printf("%d ", k)
	    _=k
	}
	fmt.Printf("Found %d primes\n", len(keys))
}

