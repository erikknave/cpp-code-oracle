match (c:CONTAINER)<-[:IMPLEMENTS]-(cd:CODEBLOCK)-[:calls|has_type]->(cd2:CODEBLOCK)-[:IMPLEMENTS]->(c2:CONTAINER) 
return c,cd,cd2,c2 limit 40
