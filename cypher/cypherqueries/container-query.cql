MATCH (c:CONTAINER { dbid: {{.DBID }} })
OPTIONAL MATCH (child:CONTAINER)-[:BELONGS_TO]->(c)
WITH c, COLLECT(child) + [c] AS allContainers

// Process Repository Information
CALL {
  WITH allContainers
  UNWIND allContainers AS container
  MATCH (container)<-[:IMPLEMENTS]-(cb:CODEBLOCK)
  MATCH (cb)-[:IS_DEFINED_IN]->(:FILE)-[:IS_LOCATED_IN]->(:DIRECTORY)-[:IS_PART_OF]->(r:REPOSITORY)
  RETURN r
  LIMIT 1
}
WITH c, allContainers, r

// Process Authors and Latest Update
CALL {
  WITH allContainers
  UNWIND allContainers AS container
  MATCH (container)<-[:IMPLEMENTS]-(cb:CODEBLOCK)
  MATCH (cb)-[:IS_DEFINED_IN]->(:FILE)-[:COMMIT_BELONGS_TO]->(fc:FILECOMMIT)
  RETURN COLLECT(DISTINCT fc.authorName) AS authors, MAX(fc.commitDate) AS latestUpdate
}
WITH c, allContainers, r, authors, latestUpdate

// Process is_used_by_containers
CALL {
  WITH allContainers
  UNWIND allContainers AS container
  MATCH (container)<-[:IMPLEMENTS]-(cb:CODEBLOCK)
  MATCH (cb)-[:calls|has_type]->(:CODEBLOCK)-[:IMPLEMENTS]->(c2:CONTAINER)
  WHERE NOT c2 IN allContainers
  WITH c2, COUNT(*) AS count_used
  RETURN COLLECT(DISTINCT {
    name: c2.name,
    dbid: c2.dbid,
    type: c2.containerType,
    count: count_used,
    signature: c2.signature
  }) AS is_used_by_containers
}
WITH c, allContainers, r, authors, latestUpdate, is_used_by_containers

// Process is_using_containers
CALL {
  WITH allContainers
  UNWIND allContainers AS container
  MATCH (container)<-[:IMPLEMENTS]-(:CODEBLOCK)-[:calls|has_type]->(cb:CODEBLOCK)
  MATCH (cb)-[:IMPLEMENTS]->(c3:CONTAINER)
  WHERE NOT c3 IN allContainers
  WITH c3, COUNT(*) AS count_using
  RETURN COLLECT(DISTINCT {
    name: c3.name,
    dbid: c3.dbid,
    type: c3.containerType,
    count: count_using,
    signature: c3.signature
  }) AS is_using_containers
}
WITH c, allContainers, r, authors, latestUpdate, is_used_by_containers, is_using_containers

// Process Parent and Grandparent Containers
OPTIONAL MATCH (c)-[:BELONGS_TO]->(parent:CONTAINER)
OPTIONAL MATCH (parent)-[:BELONGS_TO]->(grandParent:CONTAINER)
WITH c, allContainers, r, authors, latestUpdate, is_used_by_containers, is_using_containers, parent, grandParent,
     [child IN allContainers WHERE child <> c | {
       name: child.name,
       dbid: child.dbid,
       type: child.containerType,
       signature: child.signature,
       count: 0,
       summary: child.summary
     }] AS child_containers

RETURN {
  is_used_by_containers: is_used_by_containers,
  is_using_containers: is_using_containers,
  parentContainer: CASE WHEN parent IS NOT NULL THEN {
    name: parent.name,
    dbid: parent.dbid,
    type: parent.containerType,
    count: 0,
    signature: parent.signature,
    summary: parent.summary
  } ELSE null END,
  grandParentContainer: CASE WHEN grandParent IS NOT NULL THEN {
    name: grandParent.name,
    dbid: grandParent.dbid,
    type: grandParent.containerType,
    count: 0,
    signature: grandParent.signature,
    summary: grandParent.summary
  } ELSE null END,
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
