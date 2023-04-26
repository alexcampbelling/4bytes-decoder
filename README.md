# 4bytes-decoder

- Taken ~100mb of known function names from the https://www.4byte.directory/ project
- This is read into a map in < 1 second on start up
- We can now O(1) get the argument names and types for any transaction input
- The mapping is then applied to the input string and decoded JUST like we currently do, but with no backend!

Example:

Testing with input:  0xa9059cbb000000000000000000000000f3fffd51075a4f60735ed8505d30ac46951303a90000000000000000000000000000000000000000000000000000000005f5e100
Parsed method name:  transfer
Parsed input arguments: {
  "to": "0xf3fffd51075a4f60735ed8505d30ac46951303a9",
  "val": 100000000
}


No calls to the OWL microservice, no full ABI being parsed or held in a backend database, no fussssss