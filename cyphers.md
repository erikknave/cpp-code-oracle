# Main cypher with total model
MATCH (r1:Repository)-[:HAS_MODULE]->(m1:Module)<-[:PART_OF_MODULE]-(p1:Package)-[:CONTAINS]->(f1:File)-[:DEFINES]->(e1:Entity)-[:USES]->(e2:Entity)<-[:DEFINES]-(f2:File)<-[:CONTAINS]-(p2:Package)-[:PART_OF_MODULE]->(m2:Module)<-[:HAS_MODULE]-(r2:Repository)
MATCH (fc1:FileCommit)-[:AFFECTS]->(f1:File)
MATCH (fc2:FileCommit)-[:AFFECTS]->(f2:File)
RETURN r1, m1, p1, f1, e1, e2, f2, p2, m2, r2, fc1, fc2

The following are the node labels and corresponding node properties:
Repository:
 - dbid  (used to identify a certain repository)
 - name (the name / path of the repo, e.g. github.com/erikknave/go-code-oracle)
 - shortsummary (a short summary, a couple of sentences long of the repository)
 - summary (a long summary of the repository)

Module:
 - dbid  (used to identify a certain module)
 - name (the name / path of the module, e.g. github.com/erikknave/go-code-oracle)
 - shortsummary (a short summary, a couple of sentences long of the module)
 - summary (a long summary of the module)

Package:
 - dbid  (used to identify a certain package)
 - name (the name / path of the package, e.g. server, types or main)
 - repoPath (The import path to the package) e.g. github.com/erikknave/go-code-oracle/server
 - shortsummary (a short summary, a couple of sentences long of the package)
 - summary (a long summary of the package)

File:
 - dbid  (used to identify a certain file)
 - name (the file name, e.g. server.go, main.go)
 - repoPath (The import path to the file) e.g. github.com/erikknave/go-code-oracle/server/server.go
 - summary (a long summary of the file)

 Entity:
 - dbid  (used to identify a certain entity, such as a function or variable)
 - type (what type of entity it is, can be func, const, var, type)
 - name (the name, e.g. InitServer or main)
 - signature (the signature of the entity, e.g. func main())
 - string (the entire code content of the entity)
 - summary (a long summary of the file)

  FileCommit:
 - dbid  (used to identify a certain file commit)
 - authorName (The name of the author performing the code commit)
 - authorEmail (The email of the author performing the code commit)
 - commitDate (the date the commit was performed, e.g. "2024-07-08T14:45:39+01:00")
 - message (the message the author wrote to the commit e.g. Fixed certain bug )
 - patchString (the patchString of the commit)
 - summary (a summary of the commit)

 # cypher to get repo and packages
 MATCH (r:Repository {dbid: "5"})-[:HAS_MODULE]->(m:Module)-[:PART_OF_MODULE]-(p:Package)
RETURN r.name AS name, 
       r.shortsummary AS shortsummary, 
       COLLECT({name: p.name, shortsummary: p.shortsummary,importpath:p.repoPath}) AS packages
