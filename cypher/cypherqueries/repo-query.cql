// Query 1: Matching and collecting repositories using or being used by the repository with dbid "{{.DBID}}"
MATCH (r1:REPOSITORY { dbid: {{.DBID }} })<-[:IS_PART_OF]-(d1:DIRECTORY)<-[:IS_LOCATED_IN]-(f1:FILE)<-[:IS_DEFINED_IN]-(cd:CODEBLOCK)-[:calls|has_type]->(cd2:CODEBLOCK)-[:IS_DEFINED_IN]->(f2:FILE)-[:IS_LOCATED_IN]->(d2:DIRECTORY)-[: IS_PART_OF ]->(r2:REPOSITORY)
WHERE r2.dbid <> {{.DBID }} AND r2.name <> "DUMMY_REPO_NAME_NOT_EXIST"
WITH r2.name AS r2_name, r2.dbid AS r2_dbid, COUNT(*) AS r2_count
WITH r2_name, r2_dbid, r2_count
 ORDER BY r2_count DESC
WITH COLLECT( DISTINCT { name: r2_name, dbid: r2_dbid, count: r2_count }) AS is_used_by_repos
OPTIONAL MATCH (r3:REPOSITORY { dbid: {{.DBID }} })<-[:IS_PART_OF]->(d3:DIRECTORY)<-[:IS_LOCATED_IN]-(f3:FILE)<-[:IS_DEFINED_IN ]-(cd3:CODEBLOCK)-[:calls|has_type]->(cd:CODEBLOCK)-[:IS_DEFINED_IN]->(f:FILE)-[:IS_LOCATED_IN]->(d:DIRECTORY)-[: IS_PART_OF ]->(r:REPOSITORY)
WHERE r3.dbid <> {{.DBID }} AND r3.name <> "DUMMY_REPO_NAME_NOT_EXIST"
WITH is_used_by_repos, r3.name AS r3_name, r3.dbid AS r3_dbid, COUNT(*) AS r3_count
WITH is_used_by_repos, r3_name, r3_dbid, r3_count
 ORDER BY r3_count DESC
WITH is_used_by_repos, COLLECT( DISTINCT { name: r3_name, dbid: r3_dbid, count: r3_count }) AS is_using_repos

// Query 2: Matching packages of the repository with dbid "{{.DBID}}"
OPTIONAL MATCH (r:REPOSITORY { dbid: {{.DBID }} })<-[:IS_PART_OF]-(d:DIRECTORY)
WITH is_used_by_repos, is_using_repos, r, COLLECT({ name: d.name, shortsummary: d.summary, importpath: d.importPath, dbid:d.dbid }) AS directories
// Query 3: Matching file commits affecting the repository with dbid {{.DBID}}
OPTIONAL MATCH (r)<-[:IS_PART_OF]-(d:DIRECTORY)<-[:IS_LOCATED_IN]-(f:FILE)<-[:COMMIT_BELONGS_TO]-(fc:FILECOMMIT)
WITH is_used_by_repos, is_using_repos, directories, r.name AS name, r.dbid AS dbid, r.summary AS shortsummary, r.summary AS summary, COLLECT( DISTINCT fc.authorName) AS authors, MAX(fc.commitDate) AS latestUpdate
// Combining all results
RETURN COLLECT({
  is_used_by_repos: is_used_by_repos,
  is_using_repos: is_using_repos,
  name: name,
  shortsummary: shortsummary,
  summary: summary,
  authors: authors,
  latestUpdate: latestUpdate,
  directories: directories,
  dbid: dbid
  }) AS result
