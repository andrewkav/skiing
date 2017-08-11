package main

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	inf = 9999999999
)

func read(fileName string) (N, M int, A [][]int) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	l := scanner.Text()

	nm := strings.Split(string(l), " ")
	N, err = strconv.Atoi(nm[0])
	if err != nil {
		log.Fatal(err)
	}
	M, err = strconv.Atoi(nm[1])
	if err != nil {
		log.Fatal(err)
	}

	A = make([][]int, N)
	i := 0
	for scanner.Scan() {
		A[i] = make([]int, M)
		l := scanner.Text()
		la := strings.Split(string(l), " ")

		for j := 0; j < M; j++ {
			aij, err := strconv.Atoi(la[j])
			if err != nil {
				log.Fatal(err)
			}
			A[i][j] = aij
		}
		i++
	}
	return
}

type edge struct {
	U int
	V int
}

func sourcesAndEdges(N, M int, A [][]int) (sources []int, E []edge) {
	E = make([]edge, 0, N*M)
	var (
		aij                   int
		left, right, up, down int
	)
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			v := i*M + j

			aij = A[i][j]

			// left
			if j+1 < M && aij > A[i][j+1] {
				E = append(E, edge{v, v + 1})
			}
			// right
			if j-1 > -1 && aij > A[i][j-1] {
				E = append(E, edge{v, v - 1})
			}
			// up
			if i+1 < N && aij > A[i+1][j] {
				E = append(E, edge{U: v, V: v + M})
			}
			// down
			if i-1 > -1 && aij > A[i-1][j] {
				E = append(E, edge{U: v, V: v - M})
			}

			left, right, up, down = -inf, -inf, -inf, -inf
			if j > 0 {
				left = A[i][j-1]
			}
			if j < M-1 {
				right = A[i][j+1]
			}
			if i > 0 {
				up = A[i-1][j]
			}
			if i < N-1 {
				down = A[i+1][j]
			}

			if aij > left && aij > right && aij > up && aij > down {
				sources = append(sources, v)
			}
		}
	}

	return
}

func process(N, M int, E []edge, S []int) (maxPathLen, resV int) {

	var (
		NN        = N * M
		d         = make([]int, NN)
		globalMin = inf
		drop      = -1
	)

	for count, s := range S {
		if count%100 == 0 {
			log.Println("processing ", count)
		}
		for i := 0; i < NN; i++ {
			d[i] = inf
		}

		d[s] = 0

		localMin := inf
		var relaxed = false
		var du, dv, v, lenE int

		lenE = len(E)
		for i := 0; i < NN-1; i++ {
			relaxed = false

			for j := 0; j < lenE; j++ {
				v = E[j].V
				du = d[E[j].U]
				dv = d[v]
				if du != inf && dv > du-1 {
					d[v] = du - 1
					relaxed = true
				}
			}
			if !relaxed {
				break
			}

			for ii := 0; ii < NN; ii++ {
				if d[ii] < localMin {
					localMin = d[ii]
				}
			}
		}

		if localMin < globalMin {
			globalMin = localMin
			drop = s
		}
	}

	log.Printf("drop: %d, Path length: %d", drop, -globalMin+1)

	return -globalMin + 1, drop
}

func main() {
	N, M, A := read("map1.txt")

	S, E := sourcesAndEdges(N, M, A)

	startedAt := time.Now()

	splitNum := runtime.NumCPU()
	splitLen := len(S) / splitNum
	type result struct {
		drop    int
		pathLen int
	}

	globalResult := make([]result, splitNum)

	wg := &sync.WaitGroup{}
	for i := 0; i < splitNum; i++ {
		wg.Add(1)
		go func(idx int) {
			right := (idx + 1) * splitLen
			if right >= len(S) {
				right = len(S)
			}
			pathLen, drop := process(N, M, E, S[idx*splitLen:right])
			globalResult[idx] = result{
				drop:    drop,
				pathLen: pathLen,
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	maxPathLen := 0
	drop := -1
	for _, r := range globalResult {
		if r.pathLen > maxPathLen {
			drop = r.drop
			maxPathLen = r.pathLen
		}
	}

	log.Println("--------------------------------------")
	log.Printf("Elapsed time: %.2f sec", time.Now().Sub(startedAt).Seconds())
	log.Printf("Drop: %d, Path length: %d", A[drop/M][drop%M], maxPathLen)
}
