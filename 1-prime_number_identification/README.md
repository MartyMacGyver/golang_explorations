# Golang Explorations
## Prime number identification

Here we use Go to identify prime numbers in a given range using language features such as channels and goroutines to divide the workload.

While there is a progression of functionality (e.g., command-line argument handling) and some optimization, this is intended as a demonstration of the language and its features, as well as the pros and cons involved as we progress through multiple iterations of the program.

The optimization between primer5.go and primer6.go is of particular interest here - ~~while the number of cycles used is much better (an approximately 6-fold reduction) the timings when spread over 8 cores were not noticeably improved. The timings over two cores on a different machine were, however, also improved 6-fold. Why is this?~~ and does in fact work as expected (either there was a problem with the version of Go originaly used or some system issue manifested that caused it to take unexpectedly long in initial tests). The optimization depends on the presence of a list of contiguous primes in the set [2,3,5,...,sqrt(n)] for testing the primality of a value *n*. Note that this would help little if at all if the values of *n* were not in an ascending sequence (e.g., random or desending)... nor will primer6.go work correctly if the min value is greater than 2, yet.

