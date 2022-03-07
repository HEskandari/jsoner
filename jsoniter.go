// package jsoner implements encoding and decoding of JSON as defined in
// RFC 4627 and provides interfaces with identical syntax of standard lib encoding/json.
// Converting from encoding/json to jsoner is no more than replacing the package with jsoner
// and variable type declarations (if any).
// jsoner interfaces gives 100% compatibility with code using standard lib.
//
// "JSON and Go"
// (https://golang.org/doc/articles/json_and_go.html)
// gives a description of how Marshal/Unmarshal operate
// between arbitrary or predefined json objects and bytes,
// and it applies to jsoner.Marshal/Unmarshal as well.
//
// Besides, jsoner.Iterator provides a different set of interfaces
// iterating given bytes/string/reader
// and yielding parsed elements one by one.
// This set of interfaces reads input as required and gives
// better performance.
package jsoner