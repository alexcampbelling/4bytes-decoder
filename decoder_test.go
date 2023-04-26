package decoder

import (
	"regexp"
	"testing"
)

// A simple transfer input for a erc20 contract
// This specific example was taken from https://etherscan.io/tx/0x84fda0fd363078988979dd1eb652d6954c7cea1b169e6d56f4cddd84789e97cb
func TestDecodeErc20Transfer(t *testing.T) {
	input := "0xa9059cbb000000000000000000000000f3fffd51075a4f60735ed8505d30ac46951303a90000000000000000000000000000000000000000000000000000000005f5e100"

	wantFunctionName := regexp.MustCompile("transfer")
	wantArgs := regexp.MustCompile("{\n  \"to\": \"0xf3fffd51075a4f60735ed8505d30ac46951303a9\",\n  \"val\": 100000000\n}")

	functionName, args, err := decode(input)
	argsString := string(args)

	if !wantFunctionName.MatchString(functionName) || err != nil {
		t.Fatalf(`The decode function name returned as: %q, %v. We want this: %#q`, functionName, err, wantFunctionName)
	}
	if !wantArgs.MatchString(argsString) || err != nil {
		t.Fatalf(`The decode function arguments returned with: %q, %v. We want this: %#q`, argsString, err, wantArgs)
	}
}
