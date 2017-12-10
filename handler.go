package function

import (
	"encoding/json"
	"fmt"
	"github.com/ChrisCGH/mono"
	"strconv"
	"strings"
	"sync"
)

type Solution struct {
	Solution   string
	Score      int
	Key        string
	Iterations int
	Elapsed    string
}

func get_one_solution(solver mono.Mono_Solver, fixed *mono.Fixed_Key) Solution {
	solver.Set_fixed(*fixed)
	solver.Solve()
	s := Solution{solver.Solution(), solver.Score(), solver.Key(), solver.Iterations(), solver.Elapsed().String()}
	return s
}

// Handle a serverless request
func Handle(req []byte) string {
	// {"Ciphertext" : "lasdkjalskjd", "Limit" : 2000000, "Data_file" : "english.tri", "Crib" : "onceupon" }
	type Request struct {
		Ciphertext string
		Limit      string
		Data_file  string
		Crib       string
		Timeout    string
		Verbose    bool
	}
	var request Request
	err := json.Unmarshal(req, &request)
	if err != nil {
		panic(err)
	}
	ciphertext := strings.ToUpper(request.Ciphertext)
	limit, err := strconv.Atoi(request.Limit)
	timeout, err := strconv.Atoi(request.Timeout)
	if err != nil {
		timeout = 300
	}
	verbose := request.Verbose
	data_file := request.Data_file
	the_crib := strings.ToLower(request.Crib)
	crib_positions := mono.Possible_positions(ciphertext, the_crib)
	err = mono.Init()
	if err != nil {
		panic(err)
	}
	config := mono.NewConfig()
	_, file_type, path, err := config.FindDataFile(data_file)
	if err != nil {
		return fmt.Sprintf("Data file %s not found in configuration\n", data_file)
	}
	solver := mono.NewMono_Solver()
	switch file_type {
	case "trigram":
		solver.Set_trigraph_scoring(path)
	case "tetragram":
		solver.Set_tetragraph_scoring(path)
	case "ngraph":
		solver.Set_ngraph_scoring(path)
	}
	solver.Set_cipher_text(ciphertext)
	if limit > 0 {
		solver.Set_max_iterations(limit)
	} else {
		solver.Set_max_iterations(2000000)
	}
	solver.Set_timeout(timeout)
	if verbose {
		solver.Set_verbose()
	}
	f := mono.NewFixed_Key()
	fixed := &f
	solutions := make([]Solution, 0)
	if the_crib != "" {
		solutions_channel := make(chan Solution)
		wg := &sync.WaitGroup{}
		wg.Add(len(crib_positions))

		for the_crib_position := range crib_positions {
			cc := mono.NewCrib(ciphertext, the_crib, the_crib_position)
			if !cc.Is_possible() {
				(&cc).Next_right()
			}
			if cc.Is_possible() {
				fixed = cc.Get_fixed_key()
			}
			go func(solver mono.Mono_Solver, fixed *mono.Fixed_Key, wg *sync.WaitGroup) {
				defer wg.Done()
				s := get_one_solution(solver, fixed)
				solutions_channel <- s
			}(solver, fixed, wg)
		}
		go func() {
			wg.Wait()
			close(solutions_channel)
		}()
		for s := range solutions_channel {
			solutions = append(solutions, s)
		}
	} else {
		s := get_one_solution(solver, fixed)
		solutions = append(solutions, s)
	}
	solutionsJson, err := json.Marshal(solutions)
	if err != nil {
		return fmt.Sprintf("{\"error\" : \"Failed to Marshal solutions : %s\"}", err)
	}
	return fmt.Sprintf("%s", string(solutionsJson))
}
