// Package parser provides functionality for reading and representing .env files.
//
// An .env file is a newline-delimited list of KEY=VALUE pairs with optional
// inline comments (# ...) and support for single- and double-quoted values.
//
// Example usage:
//
//	f, err := os.Open(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//
//	em, err := parser.Parse(f)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if entry, ok := em.Get("DATABASE_URL"); ok {
//		fmt.Println(entry.Value)
//	}
//
// Supported syntax:
//
//	# Full-line comments are ignored.
//	KEY=value
//	KEY="quoted value"
//	KEY='single quoted'
//	KEY=value # inline comment
//
// Duplicate keys and lines missing an '=' sign are treated as parse errors.
package parser
