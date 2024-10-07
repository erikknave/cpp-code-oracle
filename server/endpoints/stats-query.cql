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
WITH repositoryCount, directoryCount, fileCount, codeblockCount,containerCount, count(fc) as fileCommitCount, COUNT(DISTINCT fc.authorName) AS authorCount
MATCH ()-[r]->()
WITH repositoryCount, directoryCount, fileCount, codeblockCount, containerCount, fileCommitCount, authorCount, COUNT(r) AS relationshipCount

// Subquery for most depended on repository
CALL {
    MATCH (d1:DIRECTORY)-[:IS_LOCATED_IN]-(f1:FILE)-[:IS_DEFINED_IN]-(cd1:CODEBLOCK)-[:calls|has_type]->(cd2:CODEBLOCK)-[:IS_DEFINED_IN]-(f2:FILE)-[:IS_LOCATED_IN]-(d2:DIRECTORY)
    WHERE f1 <> f2
    WITH d2, COUNT(DISTINCT d1) AS dependedOnByCount
    ORDER BY dependedOnByCount DESC
    LIMIT 1
    RETURN d2.name AS mostDependedOnName, d2.dbid AS mostDependedOnDbid, dependedOnByCount AS mostDependedOnCount
}

// Subquery for repository with most dependencies
CALL {
    MATCH (d1:DIRECTORY)<-[:IS_LOCATED_IN]-(f1:FILE)<-[:IS_DEFINED_IN]-(cd1:CODEBLOCK)-[:calls|has_type]->(cd2:CODEBLOCK)-[:IS_DEFINED_IN]->(f2:FILE)-[:IS_LOCATED_IN]->(d2:DIRECTORY)
    WHERE f1 <> f2
    WITH d1, COUNT(DISTINCT d2) AS dependenciesByCount
    ORDER BY dependenciesByCount DESC
    LIMIT 1
    RETURN d1.name AS mostDependenciesName, d1.dbid AS mostDependenciesDbid, dependenciesByCount AS mostDependenciesCount
}

WITH repositoryCount, directoryCount, fileCount, codeblockCount, containerCount, relationshipCount, fileCommitCount, authorCount, mostDependedOnName, mostDependedOnDbid, mostDependedOnCount, mostDependenciesName, mostDependenciesDbid, mostDependenciesCount


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
  mostDependencies: {name: mostDependenciesName, dbid: mostDependenciesDbid, count: mostDependenciesCount}
} AS results