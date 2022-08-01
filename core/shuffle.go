package core

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cube/lib/context"
)

type ShuffleArgs struct {
	Candidates []string `arg:"positional"`
}

type Shuffle struct{}

var r, _ = regexp.Compile(`(-?\d+)-(-?\d+)`)

func extractArgs(args []string) (result []string, err error) {
	result = []string{}

	for _, s := range args {
		if !r.MatchString(s) {
			result = append(result, s)
		} else {
			m := r.FindStringSubmatch(s)

			i, err := strconv.Atoi(m[1])
			if err != nil {
				return []string{}, fmt.Errorf("invalid begin value: %v", m[1])
			}
			end, err := strconv.Atoi(m[2])
			if err != nil {
				return []string{}, fmt.Errorf("invalid end value: %v", m[2])
			}

			for ; i <= end; i++ {
				result = append(result, strconv.Itoa(i))
			}
		}
	}
	return
}

func (h *Shuffle) Handle(req *context.ChatContext, args *ShuffleArgs) context.IResponse {
	candidates, err := extractArgs(args.Candidates)
	if err != nil {
		return context.NewTextResponse(err.Error())
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })

	return context.NewTextResponse(strings.Join(candidates, "\n"))
}
