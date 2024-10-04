# List of commands

### /help

Brings up this text, regardless of screen.

### /stats

Brings up the Stats screen, regardless of screen.

## Stats Screen (The start screen)

### /search

This will take you to the Search screen where you can perform key word or embeddings searches on the entire code base

### /chat

This will take you to the Chat screen where you can chat with all repositories

## Search Screen

### /repositories <search text>

Performs a key word (’search engine’) search on all repositories in the code based and returns the most relevant packages (ordered by relevance). Note that this is the default command, so just entering search text will also result in a repository search.

### /packages <search text>

Performs a key word (’search engine’) search on all packages in the code base and returns the most relevant packages (ordered by relevance).

### /files <search text>

Performs a key word (’search engine’) search on all files in the code base and returns the most relevant files (ordered by relevance).

### /entities <search text>

Performs a key word (’search engine’) search on all functions and variables in the code base and returns and returns the most relevant variables and functions (ordered by relevance).

### /all <search text>

Performs a key word (’search engine’) search on the entire code base, regardless of whether it is a repository, a package, a file or a variable/function.

### /embeddings <search text>

Performs an embeddings search on all repositories in the code base. 

### /embeddings packages <search text>

Performs an embeddings search on all packages in the code base

### /embeddings files <search text>

Performs an embeddings search on all files in the code base

### /clear

Clears the search results

### /chat

This will take you to the Chat screen where you can chat with all repositories

## Chat screen

### /clear

Clears the chat history

### /chat

Entering /chat while in the chat screen will clear the chat history and change to chatting with the default agent (The ‘Entire code base’ agent, where you can chat with all repositories).

### /search

This will take you to the Search screen where you can perform key word or embeddings searches on the entire code base

## Repository screen

Entering words into the prompt will perform a key word search on all packages in the current repository

### /chat

Will take you to the Chat screen where you can chat with an agent specialised in the current repository

## Package screen

Entering words into the prompt will perform a key word search on all files in the current package

### /chat

Will take you to the Chat screen where you can chat with an agent specialised in the current package

## File screen

Entering words into the prompt will perform a key word search on all functions and variables in the current file

### /file

Will take you to the Chat screen where you can chat with an agent specialised in the current file