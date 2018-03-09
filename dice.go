// Copyright 2018 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
)

type Wordlist map[uint32]string

func LoadWordlist(name string) (Wordlist, error) {
	if fp, err := os.Open(name); err != nil {
		return nil, err
	} else {
		defer fp.Close()
		return readWordlist(fp)
	}
}

func readWordlist(r io.Reader) (Wordlist, error) {
	scanner := bufio.NewScanner(r)
	results := make(Wordlist)
	for scanner.Scan() {
		pieces := strings.Split(scanner.Text(), "\t")
		if s, err := strconv.ParseUint(pieces[0], 10, 32); err != nil {
			return nil, err
		} else {
			results[uint32(s)] = pieces[1]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// Remap a value in the range [0, 7775] to a base-6 1-indexed integer
func IntToIndex(val, numOptions uint32) uint32 {
	pos := float64(0)
	result := uint32(0)
	posMax := math.Log(float64(numOptions)) / math.Log(6)
	for pos < posMax {
		result += ((val % 6) + 1) * uint32(math.Pow(10, pos))
		val /= 6
		pos++
	}
	return result
}

// Get a random integer under max.  Panics if not possible.
func GetRandUInt(max uint32) uint32 {
	maxInt := big.NewInt(int64(max))
	if n, err := rand.Int(rand.Reader, maxInt); err == nil {
		return uint32(n.Int64())
	} else {
		panic(fmt.Sprintf("Error in rand! %s", err.Error()))
	}
}

func main() {
	numWords := uint32(6)
	wordlist, err := readWordlist(strings.NewReader(eff_wordlist))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading wordlist: %s\n", err.Error())
		return
	}
	if len(os.Args) > 1 {
		if argNumWords, err := strconv.ParseUint(os.Args[1], 10, 32); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing argument: %s\n", err.Error())
			return
		} else {
			numWords = uint32(argNumWords)
		}
	}
	outwords := make([]string, 0, numWords)
	numOptions := uint32(len(wordlist))
	for ; numWords > 0; numWords-- {
		idx := IntToIndex(GetRandUInt(numOptions), numOptions)
		word, ok := wordlist[idx]
		if !ok {
			fmt.Fprintf(os.Stderr, "WTF, no word %d\n", idx)
		}
		// fmt.Printf("%05d: %s\n", idx, word)
		outwords = append(outwords, word)
	}
	fmt.Printf("%s\n", strings.Join(outwords, " "))
}
