# A knowledge database in command line
`kdb` is a knowledgebase command line tool. You can use it as is or integrate it with 3rd party tools to make your own Retrieval Augmented Generation (RAG) system.

# How to ask questions?
Simply use the `--query` parameter to ask questions. You can use the `--top` parameter to indicate how many of most relevant entries you want to retrieve.
```
./kdb --query "What is life meaning?" --top 1
```
as you can see it doesn't return anything, because it doesn't have any knowledge yet.


# How to add knowledge?
You can use the `--embed` paramter to pass a piece of knowledge to the database.
```
./kdb --embed "Seeking purpose helps us navigate life's challenges and find motivation and fulfillment"

./kdb --embed "A life worth living is often marked by meaningful connections, personal growth, and contributing to something greater than oneself."

./kdb --embed "The idea of a universal purpose varies, but many philosophies and religions suggest common themes like love, compassion, and personal growth"

./kdb --embed "Life's meaning often comes from forming deep and meaningful relationships."

./kdb --embed "Life's meaning can be found in learning, understanding, and gaining wisdom."

./kdb --embed "Leaving a lasting legacy and impact on the world provides a sense of purpose."
```

# How to list all entries?
You can do this by using `--list` parameter and entering the max number of entries to be listed.
```
./kdb --list 100
```

# How to make it a RAG solution?
You can easily integrate it with Ollama to make it a RAG solution. For example:
```
echo "Answer to the question by following the given context. IF you don't know the answer simply say you don't know. Do not do yapping. The Context: $(./kdb --query "What is the meaning of life?" --top 3)" | ollama run phi3
```