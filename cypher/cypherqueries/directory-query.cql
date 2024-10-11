// Match the directory with dbid {{.DBID }} and its repository
MATCH (d:DIRECTORY { dbid: {{.DBID }} })-[:IS_PART_OF]->(r:REPOSITORY)
WITH d, r

// Get files in the directory
optional MATCH (d)<-[:IS_LOCATED_IN]-(f:FILE)
WITH d, r, COLLECT({
  name: f.name, 
  summary: f.summary, 
  importpath: f.importPath, 
  dbid: f.dbid 
}) AS files
WITH d, r, files

// Get commit authors and latest commit date
OPTIONAL MATCH (d)<-[:IS_LOCATED_IN]-(f)<-[:COMMIT_BELONGS_TO]-(fc:FILECOMMIT)
WITH d, r, files, 
     COLLECT(DISTINCT fc.authorName) AS authors, 
     MAX(fc.commitDate) AS latestUpdate

// Get directories used by the directory with dbid {{.DBID }}
OPTIONAL MATCH (d)<-[:IS_LOCATED_IN]-(f1:FILE)
<-[:IS_DEFINED_IN]-(cd1:CODEBLOCK)
-[:calls|has_type]->(cd2:CODEBLOCK)
-[:IS_DEFINED_IN]->(f2:FILE)
-[:IS_LOCATED_IN]->(d2:DIRECTORY)
WHERE d2.dbid <> {{.DBID }}
WITH d, r, files, authors, latestUpdate, d2, COUNT(*) AS count_used
ORDER BY count_used DESC
WITH d, r, files, authors, latestUpdate, 
     COLLECT(DISTINCT { 
       name: d2.name, 
       dbid: d2.dbid, 
       importpath: d2.importPath, 
       count: count_used 
     }) AS is_used_by_directories

// Get directories that use the directory with dbid {{.DBID }}
OPTIONAL MATCH (d3:DIRECTORY)
<-[:IS_LOCATED_IN]-(f3:FILE)
<-[:IS_DEFINED_IN]-(cd3:CODEBLOCK)
-[:calls|has_type]->(cd1:CODEBLOCK)
-[:IS_DEFINED_IN]->(f1:FILE)
-[:IS_LOCATED_IN]->(d)
WHERE d3.dbid <> {{.DBID }}
WITH d, r, files, authors, latestUpdate, is_used_by_directories, d3, COUNT(*) AS count_using
ORDER BY count_using DESC
WITH d, r, files, authors, latestUpdate, is_used_by_directories, 
     COLLECT(DISTINCT { 
       name: d3.name, 
       dbid: d3.dbid, 
       importpath: d3.importPath, 
       count: count_using 
     }) AS is_using_directories

// Return the collected results
RETURN {
  is_used_by_directories: is_used_by_directories,
  is_using_directories: is_using_directories,
  name: d.name,
  shortsummary: d.summary,
  summary: d.summary,
  authors: authors,
  latestUpdate: latestUpdate,
  files: files,
  dbid: d.dbid,
  repodbid: r.dbid,
  reponame: r.name,
  importpath: d.importPath
} AS result
  