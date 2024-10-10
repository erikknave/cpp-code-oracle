MATCH (r:REPOSITORY)
WITH COUNT(r) AS repositoryCount
MATCH (d:DIRECTORY)
WITH repositoryCount, COUNT(d) AS directoryCount
MATCH (f:FILE)
WITH repositoryCount, directoryCount, COUNT(f) AS fileCount
MATCH (cb:CODEBLOCK)
WITH repositoryCount, directoryCount, fileCount, COUNT(cb) AS codeblockCount
MATCH (c:CONTAINER)
WITH repositoryCount, directoryCount, fileCount, codeblockCount, COUNT(c) AS containerCount
MATCH (fc:FILECOMMIT)
WITH repositoryCount, directoryCount, fileCount, codeblockCount, containerCount, COUNT(fc) AS fileCommitCount, COUNT(DISTINCT fc.authorName) AS authorCount
MATCH ()-[r]->()
WITH repositoryCount, directoryCount, fileCount, codeblockCount, containerCount, fileCommitCount, authorCount, COUNT(r) AS relationshipCount

CALL {
    MATCH (d1:DIRECTORY)-[:IS_LOCATED_IN]-(f1:FILE)-[:IS_DEFINED_IN]-(cd1:CODEBLOCK)-[:calls|has_type]->(cd2:CODEBLOCK)-[:IS_DEFINED_IN]-(f2:FILE)-[:IS_LOCATED_IN]-(d2:DIRECTORY)
    WHERE f1 <> f2
    WITH d2, COUNT(DISTINCT d1) AS dependedOnByCount
    ORDER BY dependedOnByCount DESC
    LIMIT 1
    RETURN d2.name AS mostDependedOnName, d2.dbid AS mostDependedOnDbid, dependedOnByCount AS mostDependedOnCount
}

CALL {
    MATCH (d1:DIRECTORY)<-[:IS_LOCATED_IN]-(f1:FILE)<-[:IS_DEFINED_IN]-(cd1:CODEBLOCK)-[:calls|has_type]->(cd2:CODEBLOCK)-[:IS_DEFINED_IN]->(f2:FILE)-[:IS_LOCATED_IN]->(d2:DIRECTORY)
    WHERE f1 <> f2
    WITH d1, COUNT(DISTINCT d2) AS dependenciesByCount
    ORDER BY dependenciesByCount DESC
    LIMIT 1
    RETURN d1.name AS mostDependenciesName, d1.dbid AS mostDependenciesDbid, dependenciesByCount AS mostDependenciesCount
}

CALL {
    MATCH (c1:CONTAINER)<-[:IMPLEMENTS]-(cd1:CODEBLOCK)-[:calls|has_type]->(cd2:CODEBLOCK)-[:IMPLEMENTS]->(c2:CONTAINER)
    WHERE c1 <> c2
    WITH c2, COUNT(DISTINCT c1) AS dependedOnByCount
    ORDER BY dependedOnByCount DESC
    LIMIT 1
    RETURN c2.name AS mostDependedOnContainerName, c2.dbid AS mostDependedOnContainerDbid, dependedOnByCount AS mostDependedOnContainerCount, c2.containerType as mostDependedOnContainerType, c2.signature as mostDependedOnContainerSignature
}

CALL {
    MATCH (c1:CONTAINER)<-[:IMPLEMENTS]-(cd1:CODEBLOCK)-[:calls|has_type]->(cd2:CODEBLOCK)-[:IMPLEMENTS]->(c2:CONTAINER)
    WHERE c1 <> c2 AND NOT c1.name = "GLOBAL_NAMESPACE_ENTITY" AND NOT c1.name = "stingray" AND NOT c1.containerType = "variable"
    WITH c1, COUNT(DISTINCT c2) AS dependenciesByCount
    ORDER BY dependenciesByCount DESC
    LIMIT 1
    RETURN c1.name AS mostDependenciesContainerName, c1.dbid AS mostDependenciesContainerDbid, dependenciesByCount AS mostDependenciesContainerCount, c1.containerType AS mostDependenciesContainerType, c1.signature as mostDependenciesContainerSignature
}

WITH repositoryCount, directoryCount, fileCount, codeblockCount, containerCount, relationshipCount, fileCommitCount, authorCount, 
mostDependedOnName, mostDependedOnDbid, mostDependedOnCount, mostDependenciesName, mostDependenciesDbid, mostDependenciesCount, 
mostDependedOnContainerName, mostDependedOnContainerDbid, mostDependedOnContainerCount, mostDependedOnContainerType, mostDependedOnContainerSignature, mostDependenciesContainerName, mostDependenciesContainerDbid, mostDependenciesContainerCount, mostDependenciesContainerType, mostDependenciesContainerSignature

RETURN {
  repositories: repositoryCount,
  directories: directoryCount,
  files: fileCount,
  codeblocks: codeblockCount,
  containers: containerCount,
  fileCommits: fileCommitCount,
  authors: authorCount,
  relationships: relationshipCount,
  mostDependedOn: {name: mostDependedOnName, dbid: mostDependedOnDbid, count: mostDependedOnCount},
  mostDependencies: {name: mostDependenciesName, dbid: mostDependenciesDbid, count: mostDependenciesCount},
  mostDependedOnContainer: {name: mostDependedOnContainerName, dbid: mostDependedOnContainerDbid, count: mostDependedOnContainerCount, type: mostDependedOnContainerType, signature: mostDependedOnContainerSignature},
  mostDependenciesContainer: {name: mostDependenciesContainerName, dbid: mostDependenciesContainerDbid, count: mostDependenciesContainerCount, type: mostDependenciesContainerType, signature: mostDependenciesContainerSignature}
} AS results
