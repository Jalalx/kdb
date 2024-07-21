# A knowledge database in your command line
`kdb` is a knowledgebase command line tool that relies on [Ollama](https://github.com/ollama/ollama). You can use it as is or integrate it with 3rd party tools to make your own Retrieval Augmented Generation (RAG) system.

## How to query information?
Simply use the `query` parameter to retrieve relevant informations. You can use the `--limit` parameter to indicate how many of most relevant entries you want to retrieve.
```
./kdb --query "What is life meaning?" --top 1
```
as you can see it doesn't return anything, because it doesn't have any knowledge yet.


## How to add knowledge?
You can use the `embed` paramter to pass a piece of knowledge to the database.
```
./kdb --embed "Here is a sample piece of information that will be remembed forever!"
```
There is a script in the `test_content` directory which you can use to add some sample knowledge.
```
cd test_content/
chmod +x ./train.sh
./train.sh
```

This will add a bunch of sample knowledge. You can now ask questions about it:
```
./kdb --query "how many times Alpha is smaller than Gamma?" --top 3
```
and it returns the revelant piece of informations.

## How to ask questions?
You can connect the outcome of `./kdb query "your question"` to Ollama to retrieve the information in a friendly way. It would be exactly like a simple Retrieval-augmented generation (RAG) solution. There is an `ask` command in this repo that makes life easier. For example:
```bash
$ ask "how many times Alpha is smaller than Gamma?"
```
Should return a reply like:
```
Since Alpha is half of Beta, and Beta is smaller than Gamma (which is three times of Beta), we can conclude that Alpha is
one-sixth of Gamma. Therefore, Alpha is six times smaller than Gamma.
```
Or you can make your own prompt for it:
```
echo "Answer to the question by following the given context. IF you don't know the answer simply say you don't know. Do not do yapping. The Context: $(./kdb --query "What is the meaning of life?" --top 3)" | ollama run llama3
```

## How to list all entries?
You can do this by using `list` parameter and entering the max number of entries to be listed.
```
./kdb --list 100
```

## How to delete an entry?
You can simply use the `delete` option:
```
./kdb --delete "b1827d5a-16d2-4f74-acfa-29864434859a"
```

## Need the list of parameters?
You can run `kdb` with the `--help` parameter to get the list of parameters.
```
./kdb --help
```
