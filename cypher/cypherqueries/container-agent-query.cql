MATCH (c:CONTAINER { dbid: {{.DBID}} })
MATCH (c)<-[:IMPLEMENTS]-(cb:CODEBLOCK)
MATCH (cb)-[:IS_DEFINED_IN]->(f:FILE)
MATCH (f)-[:IS_LOCATED_IN]->(:DIRECTORY)-[:IS_PART_OF]->(r:REPOSITORY)
WITH c, r, COLLECT({
  signature: cb.signature,
  dbid: cb.dbid,
  fileimportpath: f.importPath,
  filedbid: f.dbid,
  filesummary: f.summary
}) AS codeblocks
RETURN {
  codeblocks: codeblocks,
  name: c.name,
  summary: c.summary,
  signature: c.signature,
  type: c.containerType,
  dbid: c.dbid,
  repodbid: r.dbid,
  reponame: r.name,
  reposhortsummary: r.summary
} AS result
