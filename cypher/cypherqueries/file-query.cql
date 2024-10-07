MATCH (f:FILE { dbid: {{.DBID}} })-[:IS_LOCATED_IN]->(d:DIRECTORY)-[:IS_PART_OF]->(r:REPOSITORY)
WITH f, d, r
OPTIONAL MATCH (f)<-[:IS_DEFINED_IN]-(cb:CODEBLOCK)
WITH f, d, r, COLLECT({
  name: cb.name,
  summary: cb.summary,
  signature: cb.signature,
  importpath: cb.importPath,
  dbid: cb.dbid,
  entitytype: cb.entityType
}) AS codeblocks
WITH f, d, r, codeblocks
OPTIONAL MATCH (f)<-[:COMMIT_BELONGS_TO]-(fc:FILECOMMIT)
WITH f, d, r, codeblocks, COLLECT(DISTINCT fc.authorName) AS authors, MAX(fc.commitDate) AS latestUpdate
OPTIONAL MATCH (f)<-[:IS_DEFINED_IN]-(cd1:CODEBLOCK)-[:calls|has_type]->(cd2:CODEBLOCK)-[:IS_DEFINED_IN]->(f2:FILE)
WHERE f2.dbid <> {{.DBID}}
WITH f, d, r, codeblocks, authors, latestUpdate, f2, COUNT(*) AS count_used
ORDER BY count_used DESC
WITH f, d, r, codeblocks, authors, latestUpdate,
     COLLECT(DISTINCT {
       name: f2.name,
       dbid: f2.dbid,
       importpath: f2.importPath,
       count: count_used
     }) AS is_used_by_files
OPTIONAL MATCH (f3:FILE)<-[:IS_DEFINED_IN]-(cd3:CODEBLOCK)-[:calls|has_type]->(cd1:CODEBLOCK)-[:IS_DEFINED_IN]->(f)
WHERE f3.dbid <> {{.DBID}}
WITH f, d, r, codeblocks, authors, latestUpdate, is_used_by_files, f3, COUNT(*) AS count_using
ORDER BY count_using DESC
WITH f, d, r, codeblocks, authors, latestUpdate, is_used_by_files,
     COLLECT(DISTINCT {
       name: f3.name,
       dbid: f3.dbid,
       importpath: f3.importPath,
       count: count_using
     }) AS is_using_files
RETURN {
  is_used_by_files: is_used_by_files,
  is_using_files: is_using_files,
  name: f.fileName,
  shortsummary: f.summary,
  summary: f.summary,
  authors: authors,
  latestUpdate: latestUpdate,
  codeblocks: codeblocks,
  dbid: f.dbid,
  directorydbid: d.dbid,
  directoryname: d.name,
  repodbid: r.dbid,
  reponame: r.name,
  importpath: f.importPath
} AS result