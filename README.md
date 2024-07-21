# A knowledge database in your command line
`kdb` is a knowledgebase command line tool that relies on [Ollama](https://github.com/ollama/ollama). You can use it as is or integrate it with 3rd party tools to make your own Retrieval Augmented Generation (RAG) system.

## Dependencies
* [Ollama](https://github.com/ollama/ollama)

## Building the project
1. Clone the repository
2. `make build`
3. Your binaries are at `bin/` directory

## Run
Here is the list of available commands
```
kdb --help
```

## How embedding works?
```
kdb embed "here is some text"
```
![Image](docs/kdb-embed-diagram.png)
[source](https://www.plantuml.com/plantuml/dumla/RO_BJiCm44Nt_eghY0MLqj3oO16XgekoO4CHBIj0gfEPWjN4hk8n3o7-7LkbFAIiHXgVETTtwaLt4jf6wI3EXhRYJkbQN8S0xP7TUUuojOsGgSrMMy66AaLQTAKzWoeaQXRgpQpwpsNpfj6jpXBWb0eOJj9jYTL1ck2OHSYUkMCO3-zQXl2RP2kLdeTxt5WZkyq4hiJDVB74qUxu0vMZzx9FWa_bVXqld2gLk1yLu-EJ7AFYzEmyHr4KZtjrmgwk5vUtopYMyzttoDWd_s0FQsU5hUJVea5SzMJcVVw1-bjcZCwzUkZrEegOVfg662xSGXCnbRW8mT14JTbIQ9il)

## How data is retrieved?
```
kdb query "find something"
```
![Image](docs/kdb-query-diagram.png)
[source](https://www.plantuml.com/plantuml/dumla/RO-zJlCm68LtNyLHz0tVa8Ro0qC6L47B10fAou0eedRy0bOJkzYlaIh4lLC7HO1WikJpdCS-ay3IS-nQ8ICx6pj5NiY6dKU43CXk0lbCR7QGQSn6MiPQAw4bIoK3GQkUciPcFLK_kQngMzd9B05EY8ZHQUgMMFrB9ruY-IsoHsCrWkk8durzobOYPQE1DTPmkabbL-AwcY-mHvqYZJSefxVbBmFUXViIl58QK-9kNEmV7EOO5qV79pcAehWVl0cRpovdvmiqpZ9PpY6zbzi7RG9Bsz3_Jqk-j7zeS55NL-Z_f7VlXwMjwkcinjajXdxYnFTPDONOCW9nIeR5SJVG6ylmVHp4XAoENVi1)