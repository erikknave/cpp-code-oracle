package types

type RepoQueryResponseRepository struct {
	Name         string `json:"name"`
	Dbid         int    `json:"dbid"`
	Count        int    `json:"count"`
	ShortSummary string `json:"shortsummary"`
}

type RepoQueryResponseDirectory struct {
	Name         string `json:"name"`
	ShortSummary string `json:"shortsummary"`
	ImportPath   string `json:"importpath"`
	Dbid         int    `json:"dbid"`
}

type RepoQueryReponseResult struct {
	IsUsedByRepos []RepoQueryResponseRepository `json:"is_used_by_repos"`
	IsUsingRepos  []RepoQueryResponseRepository `json:"is_using_repos"`
	Name          string                        `json:"name"`
	ShortSummary  string                        `json:"shortsummary"`
	Summary       string                        `json:"summary"`
	Packages      []RepoQueryResponseDirectory  `json:"directories"`
	Authors       []string                      `json:"authors"`
	LatestUpdate  string                        `json:"latestUpdate"`
	Dbid          int                           `json:"dbid"`
}

type DirectoryQueryResponseDirectory struct {
	Name         string `json:"name"`
	Dbid         int    `json:"dbid"`
	Count        int    `json:"count"`
	ShortSummary string `json:"shortsummary"`
	ImportPath   string `json:"importpath"`
	RepoDbid     int    `json:"repodbid"`
	RepoName     string `json:"reponame"`
}

type DirectoryQueryResponseFile struct {
	Name       string `json:"name"`
	Summary    string `json:"summary"`
	ImportPath string `json:"importpath"`
	Dbid       int    `json:"dbid"`
}

type DirectoryQueryReponseResult struct {
	IsUsedByDirectories []DirectoryQueryResponseDirectory `json:"is_used_by_directories"`
	IsUsingDirectories  []DirectoryQueryResponseDirectory `json:"is_using_directories"`
	Name                string                            `json:"name"`
	ShortSummary        string                            `json:"shortsummary"`
	Summary             string                            `json:"summary"`
	Files               []DirectoryQueryResponseFile      `json:"files"`
	Authors             []string                          `json:"authors"`
	LatestUpdate        string                            `json:"latestUpdate"`
	Dbid                int                               `json:"dbid"`
	RepoDbid            int                               `json:"repodbid"`
	RepoName            string                            `json:"reponame"`
	ImportPath          string                            `json:"importpath"`
}

type FileQueryResponseFile struct {
	Name       string `json:"name"`
	Dbid       int    `json:"dbid"`
	Count      int    `json:"count"`
	Summary    string `json:"summary"`
	ImportPath string `json:"importpath"`
	RepoDbid   string `json:"repodbid"`
	RepoName   string `json:"reponame"`
}

type FileQueryResponseCodeblock struct {
	Name       string `json:"name"`
	Summary    string `json:"summary"`
	Signature  string `json:"signature"`
	EntityType string `json:"entitytype"`
	Dbid       int    `json:"dbid"`
}

type FileQueryReponseResult struct {
	IsUsedByFiles []FileQueryResponseFile      `json:"is_used_by_files"`
	IsUsingFiles  []FileQueryResponseFile      `json:"is_using_files"`
	Name          string                       `json:"name"`
	Summary       string                       `json:"summary"`
	Codeblocks    []FileQueryResponseCodeblock `json:"codeblocks"`
	Authors       []string                     `json:"authors"`
	LatestUpdate  string                       `json:"latestUpdate"`
	Dbid          int                          `json:"dbid"`
	RepoDbid      int                          `json:"repodbid"`
	RepoName      string                       `json:"reponame"`
	ImportPath    string                       `json:"importpath"`
	DirectoryDbid int                          `json:"directorydbid"`
	DirectoryName string                       `json:"directoryname"`
}

type EntityQueryResponseEntity struct {
	Name      string `json:"name"`
	Dbid      int    `json:"dbid"`
	Count     int    `json:"count"`
	Signature string `json:"signature"`
	RepoDbid  int    `json:"repodbid"`
	RepoName  string `json:"reponame"`
}

type EntityQueryResponseFile struct {
	Name      string `json:"name"`
	Summary   string `json:"summary"`
	Signature string `json:"signature"`
	Dbid      int    `json:"dbid"`
}

type EntityQueryResponseResult struct {
	IsUsedByEntities    []EntityQueryResponseEntity `json:"is_used_by_entities"`
	IsUsingEntities     []EntityQueryResponseEntity `json:"is_using_entities"`
	Name                string                      `json:"name"`
	Summary             string                      `json:"summary"`
	Signature           string                      `json:"signature"`
	Authors             []string                    `json:"authors"`
	LatestUpdate        string                      `json:"latestUpdate"`
	Dbid                int                         `json:"dbid"`
	RepoDbid            string                      `json:"repodbid"`
	RepoName            string                      `json:"reponame"`
	RepoShortSummary    string                      `json:"reposhortsummary"`
	ImportPath          string                      `json:"importpath"`
	PackageDbid         string                      `json:"packagedbid"`
	PackageImportPath   string                      `json:"packageimportpath"`
	PackageShortSummary string                      `json:"packageshortsummary"`
	FileDbid            string                      `json:"filedbid"`
	FileName            string                      `json:"filename"`
	FileSummary         string                      `json:"filesummary"`
}

type RepoListQueryResult struct {
	Name string `json:"name"`
	Dbid int    `json:"dbid"`
}

type ListFilesResponseResult struct {
	Type       string `json:"type"`
	Dbid       int    `json:"dbid"`
	ImportPath string `json:"importPath"`
	Files      []struct {
		ImportPath string `json:"importPath"`
		// Summary    string `json:"summary"`
		Dbid string `json:"dbid"`
		Name string `json:"name"`
	} `json:"files"`
}

type Stats struct {
	Repositories     int                 `json:"repositories"`
	Directories      int                 `json:"directories"`
	Files            int                 `json:"files"`
	Codeblocks       int                 `json:"codeblocks"`
	FileCommits      int                 `json:"fileCommits"`
	Authors          int                 `json:"authors"`
	Relationships    int                 `json:"relationships"`
	MostDependedOn   StatsRepositoryInfo `json:"mostDependedOn"`
	MostDependencies StatsRepositoryInfo `json:"mostDependencies"`
}

type StatsRepositoryInfo struct {
	Name  string `json:"name"`
	Dbid  int    `json:"dbid"`
	Count int    `json:"count"`
}
