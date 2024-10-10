MATCH (c:CONTAINER { dbid: {{.DBID }} })
OPTIONAL MATCH (child:CONTAINER)-[:BELONGS_TO]->(c)
WITH c, COLLECT(child) AS children

WITH c, children, [c] + children AS allContainers
UNWIND allContainers AS container
WITH c, children, COLLECT(container.dbid) AS containerDbids

OPTIONAL MATCH (container)<-[:IMPLEMENTS]-(cb:CODEBLOCK)
WHERE container.dbid IN containerDbids
WITH c, children, containerDbids, COLLECT(DISTINCT cb) AS cbs

OPTIONAL MATCH (cb)-[:IS_DEFINED_IN]->(:FILE)-[:IS_LOCATED_IN]->(:DIRECTORY)-[:IS_PART_OF]->(r:REPOSITORY)
WITH c, children, containerDbids, cbs, HEAD(COLLECT(r)) AS r

OPTIONAL MATCH (cb)-[:IS_DEFINED_IN]->(:FILE)-[:COMMIT_BELONGS_TO]->(fc:FILECOMMIT)
WITH c, children, containerDbids, cbs, r, COLLECT(DISTINCT fc.authorName) AS authors, MAX(fc.commitDate) AS latestUpdate

OPTIONAL MATCH (cb)-[:calls|has_type]->(:CODEBLOCK)-[:IMPLEMENTS]->(c2:CONTAINER)
WHERE NOT c2.dbid IN containerDbids
WITH c, children, containerDbids, cbs, r, authors, latestUpdate, c2, COUNT(*) AS count_used
ORDER BY count_used DESC
WITH c, children, containerDbids, cbs, r, authors, latestUpdate,
     COLLECT(DISTINCT { name: c2.name, dbid: c2.dbid, type: c2.containerType, count: count_used, signature: c2.signature }) AS is_used_by_containers

OPTIONAL MATCH (c3:CONTAINER)<-[:IMPLEMENTS]-(:CODEBLOCK)-[:calls|has_type]->(cb)
WHERE NOT c3.dbid IN containerDbids
WITH c, children, containerDbids, cbs, r, authors, latestUpdate, is_used_by_containers, c3, COUNT(*) AS count_using
ORDER BY count_using DESC
WITH c, children, containerDbids, cbs, r, authors, latestUpdate, is_used_by_containers,
     COLLECT(DISTINCT { name: c3.name, dbid: c3.dbid, type: c3.containerType, count: count_using, signature: c3.signature }) AS is_using_containers

OPTIONAL MATCH (c)-[:BELONGS_TO]->(parent:CONTAINER)
WITH c, children, r, authors, latestUpdate, is_used_by_containers, is_using_containers, parent,
     [child IN children | { name: child.name, dbid: child.dbid, type: child.containerType, signature: child.signature, count: 0, summary: child.summary }] AS child_containers

RETURN {
  is_used_by_containers: is_used_by_containers,
  is_using_containers: is_using_containers,
  parentContainer: CASE WHEN parent IS NOT NULL THEN { name: parent.name, dbid: parent.dbid, type: parent.containerType, count: 0, signature: parent.signature, summary: parent.summary } ELSE null END,
  childContainers: child_containers,
  name: c.name,
  summary: c.summary,
  signature: c.signature,
  type: c.containerType,
  authors: authors,
  latestUpdate: latestUpdate,
  dbid: c.dbid,
  repodbid: r.dbid,
  reponame: r.name,
  reposhortsummary: r.summary
} AS result
