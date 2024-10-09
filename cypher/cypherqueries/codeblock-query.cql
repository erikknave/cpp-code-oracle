MATCH (cb:CODEBLOCK { dbid: {{.DBID}} })
OPTIONAL MATCH (cb)-[:IS_DEFINED_IN]->(f:FILE)
OPTIONAL MATCH (f)-[:IS_LOCATED_IN]->(d:DIRECTORY)
OPTIONAL MATCH (d)-[:IS_PART_OF]->(r:REPOSITORY)
OPTIONAL MATCH (cb)-[:IS_DEFINED_IN]->(f)<-[:COMMIT_BELONGS_TO]-(fc:FILECOMMIT)
WITH cb, f, d, r, COLLECT(DISTINCT fc.authorName) AS authors, MAX(fc.commitDate) AS latestUpdate

OPTIONAL MATCH (cb)-[:calls|has_type]->(cb2:CODEBLOCK)
OPTIONAL MATCH (cb2)-[:IS_DEFINED_IN]->(f2:FILE)
OPTIONAL MATCH (f2)-[:IS_LOCATED_IN]->(d2:DIRECTORY)
OPTIONAL MATCH (d2)-[:IS_PART_OF]->(r2:REPOSITORY)
WHERE cb2.dbid <> {{.DBID}}
WITH cb, f, d, r, authors, latestUpdate, cb2, r2, COUNT(*) AS count_used
ORDER BY count_used DESC
WITH cb, f, d, r, authors, latestUpdate, 
     COLLECT(DISTINCT {
       name: cb2.name,
       type: cb2.entityType,
       dbid: cb2.dbid,
       signature: cb2.signature,
       repodbid: r2.dbid,
       reponame: r2.name,
       count: count_used
     }) AS is_used_by_codeblocks

OPTIONAL MATCH (cb3:CODEBLOCK)-[:calls|has_type]->(cb)
OPTIONAL MATCH (cb3)-[:IS_DEFINED_IN]->(f3:FILE)
OPTIONAL MATCH (f3)-[:IS_LOCATED_IN]->(d3:DIRECTORY)
OPTIONAL MATCH (d3)-[:IS_PART_OF]->(r3:REPOSITORY)
WHERE cb3.dbid <> {{.DBID}}
WITH cb, f, d, r, authors, latestUpdate, is_used_by_codeblocks, cb3, r3, COUNT(*) AS count_using
ORDER BY count_using DESC
WITH cb, f, d, r, authors, latestUpdate, is_used_by_codeblocks, 
     COLLECT(DISTINCT {
       name: cb3.name,
       type: cb3.entityType,
       dbid: cb3.dbid,
       signature: cb3.signature,
       repodbid: r3.dbid,
       reponame: r3.name,
       count: count_using
     }) AS is_using_codeblocks

RETURN {
  is_used_by_codeblocks: is_used_by_codeblocks,
  is_using_codeblocks: is_using_codeblocks,
  name: cb.name,
  type: cb.entityType,
  summary: cb.summary,
  signature: cb.signature,
  authors: authors,
  latestUpdate: latestUpdate,
  startoffset: cb.startOffset,
  endoffset: cb.endOffset,
  dbid: cb.dbid,
  repodbid: r.dbid,
  reponame: r.name,
  reposhortsummary: r.summary,
  importpath: f.importPath,
  directorydbid: d.dbid,
  directoryimportpath: d.importPath,
  directoryshortsummary: d.summary,
  filedbid: f.dbid,
  filename: f.fileName,
  filesummary: f.summary
} AS result
