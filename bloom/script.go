package bloom

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

func convert2Arg(arg []uint64) []string {
	var offsets = make([]string, len(arg))
	for index, value := range arg {
		offsets[index] = strconv.FormatUint(value, 10)
	}
	return offsets
}

func (bf *BloomFilter) setBit(ctx context.Context, offsets []uint64) error {
	file, err := os.ReadFile("C:\\Users\\王佳宝20031205\\GolandProjects\\grpc_project\\bloom\\setBit.lua")
	if err != nil {
		return err
	}
	fmt.Println(string(file))
	arg := convert2Arg(offsets)
	script := redis.NewScript(string(file))
	cmd := script.Run(ctx, bf.store, []string{bf.key}, arg)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (bf *BloomFilter) getBit(offsets []uint64) (bool, error) {
	file, err := os.ReadFile("C:\\Users\\王佳宝20031205\\GolandProjects\\grpc_project\\bloom\\getbit.lua")
	if err != nil {
		return false, err
	}
	args := convert2Arg(offsets)

	script := redis.NewScript(string(file))
	cmd := script.Run(context.Background(), bf.store,
		[]string{bf.key}, args)

	// BloomFilter has not been created,which represent no a single data in the filter
	if cmd.Err() == redis.Nil {
		return false, nil
	} else if cmd.Err() != nil {
		return false, cmd.Err()
	}
	result, err := cmd.Bool()
	if err != nil {
		return false, err
	}
	return result, nil
}
