package decoder

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type DecodedInput map[string]interface{}

var fileContents map[string]string

func init() {
	// Directory containing the files
	dir := "./with_parameter_names"

	// Read the contents of the files and store them in a map
	var err error
	fileContents, err = readFiles(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func readFiles(dir string) (map[string]string, error) {
	// Initialize an empty map to store the file contents
	files := make(map[string]string)

	// Get a list of all files in the directory
	fileList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Read the contents of each file and store them in the map
	for _, fileInfo := range fileList {
		filename := fileInfo.Name()
		filePath := filepath.Join(dir, filename)
		fileContents, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		files[filename] = string(fileContents)
	}

	return files, nil
}

// TODO: should we be able to accept bytes ever?
func decode(input string) (string, []byte, error) {
	// Basic function
	// 1. Check if input valid
	// 2. Use 4byte to scan a cache / map for if we know it
	// 3. Grab the method arguments string
	// 4. Make an abi class out of it
	// 5. Retrieve the method name
	// 6. Retrieve the inputs
	// 7. Return

	// -- Everything below is a test

	fmt.Println("Testing with input: ", input)

	input = remove0xPrefix(input)

	inputBytes, err := hex.DecodeString(input[8:])
	if err != nil {
		return "", nil, err
	}

	signature, err := getSignature(input[:8])
	if err != nil {
		return "", nil, err
	}

	functionName, args, err := parseSignature(signature)
	if err != nil {
		return "", nil, err
	}

	fmt.Println("Parsed method name: ", functionName)

	decodedInput := make(map[string]interface{})

	err = args.UnpackIntoMap(decodedInput, inputBytes)
	if err != nil {
		return "", nil, err
	}

	// Convert the map to JSON
	jsonData, err := json.MarshalIndent(decodedInput, "", "  ")
	if err != nil {
		log.Fatal("Error converting decoded input to JSON:", err)
	}

	fmt.Printf("Parsed input arguments: %s\n", jsonData)

	return functionName, jsonData, nil
}

func remove0xPrefix(str string) string {
	if strings.HasPrefix(str, "0x") {
		return str[2:]
	}
	return str
}

func getSignature1(sigHash string) (string, error) {
	// TODO: some sort of map / data structure retireval here!
	// Use this: https://github.com/ethereum-lists/4bytes

	// We curl to the github currently for tests
	// URL with the path parameter
	fmt.Println("sigHash: ", sigHash)
	url := "https://raw.githubusercontent.com/ethereum-lists/4bytes/master/with_parameter_names/" + sigHash

	// Make a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	// Read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	fmt.Println("Response: ", string(body))

	parts := strings.Split(string(body), ";")
	if len(parts) > 1 {
		result := parts[1]
		fmt.Println("Just grabbing the last collision for now!", result)
		return result, nil
	}

	// TODO: Remove this hardcode
	return "", nil

}

func getSignature(sigHash string) (string, error) {
	contents := fileContents[sigHash]

	parts := strings.Split(contents, ";")
	if len(parts) > 1 {
		result := parts[1]
		fmt.Println("Just grabbing the last collision for now!", result)
		return result, nil
	}

	return contents, nil
}

func parseSignature(signature string) (string, abi.Arguments, error) {
	re := regexp.MustCompile(`(?P<name>\w+)\((?P<params>.*)\)`)
	matches := re.FindStringSubmatch(signature)

	if len(matches) != 3 {
		return "", nil, fmt.Errorf("invalid signature")
	}

	functionName := matches[1]
	rawParams := strings.Split(matches[2], ",")

	args := abi.Arguments{}
	for _, rawParam := range rawParams {
		parts := strings.Split(rawParam, " ")
		if len(parts) == 1 {
			argType, err := abi.NewType(normalizeType(parts[0]), "", nil)
			if err != nil {
				return "", nil, err
			}
			args = append(args, abi.Argument{Type: argType})
		} else if len(parts) == 2 {
			argType, err := abi.NewType(normalizeType(parts[0]), "", nil)
			if err != nil {
				return "", nil, err
			}
			args = append(args, abi.Argument{Name: parts[1], Type: argType})
		} else {
			return "", nil, fmt.Errorf("invalid parameter format")
		}
	}

	return functionName, args, nil
}

func normalizeType(paramType string) string {
	if strings.HasPrefix(paramType, "int") && !strings.Contains(paramType, "[") {
		if paramType == "int" {
			return "int256"
		}
		return paramType
	}

	if strings.HasPrefix(paramType, "uint") && !strings.Contains(paramType, "[") {
		if paramType == "uint" {
			return "uint256"
		}
		return paramType
	}

	return paramType
}
