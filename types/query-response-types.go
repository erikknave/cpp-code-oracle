package types

type RepoQueryResponseRepository struct {
	Name         string `json:"name"`
	Dbid         string `json:"dbid"`
	Count        int    `json:"count"`
	ShortSummary string `json:"shortsummary"`
}

type RepoQueryResponsePackage struct {
	Name         string `json:"name"`
	ShortSummary string `json:"shortsummary"`
	ImportPath   string `json:"importpath"`
	Dbid         string `json:"dbid"`
}

type RepoQueryReponseResult struct {
	IsUsedByRepos []RepoQueryResponseRepository `json:"is_used_by_repos"`
	IsUsingRepos  []RepoQueryResponseRepository `json:"is_using_repos"`
	Name          string                        `json:"name"`
	ShortSummary  string                        `json:"shortsummary"`
	Summary       string                        `json:"summary"`
	Packages      []RepoQueryResponsePackage    `json:"packages"`
	Authors       []string                      `json:"authors"`
	LatestUpdate  string                        `json:"latestUpdate"`
	Dbid          string                        `json:"dbid"`
}

type PackageQueryResponsePackage struct {
	Name         string `json:"name"`
	Dbid         string `json:"dbid"`
	Count        int    `json:"count"`
	ShortSummary string `json:"shortsummary"`
	ImportPath   string `json:"importpath"`
	RepoDbid     string `json:"repodbid"`
	RepoName     string `json:"reponame"`
}

type PackageQueryResponseFile struct {
	Name       string `json:"name"`
	Summary    string `json:"summary"`
	ImportPath string `json:"importpath"`
	Dbid       string `json:"dbid"`
}

type PackageQueryReponseResult struct {
	IsUsedByPackages []PackageQueryResponsePackage `json:"is_used_by_packages"`
	IsUsingPackages  []PackageQueryResponsePackage `json:"is_using_packages"`
	Name             string                        `json:"name"`
	ShortSummary     string                        `json:"shortsummary"`
	Summary          string                        `json:"summary"`
	Files            []PackageQueryResponseFile    `json:"files"`
	Authors          []string                      `json:"authors"`
	LatestUpdate     string                        `json:"latestUpdate"`
	Dbid             string                        `json:"dbid"`
	RepoDbid         string                        `json:"repodbid"`
	RepoName         string                        `json:"reponame"`
	ImportPath       string                        `json:"importpath"`
}

type FileQueryResponseFile struct {
	Name       string `json:"name"`
	Dbid       string `json:"dbid"`
	Count      int    `json:"count"`
	Summary    string `json:"summary"`
	ImportPath string `json:"importpath"`
	RepoDbid   string `json:"repodbid"`
	RepoName   string `json:"reponame"`
}

type FileQueryResponseEntity struct {
	Name      string `json:"name"`
	Summary   string `json:"summary"`
	Signature string `json:"signature"`
	Dbid      string `json:"dbid"`
}

type FileQueryReponseResult struct {
	IsUsedByFiles     []FileQueryResponseFile   `json:"is_used_by_files"`
	IsUsingFiles      []FileQueryResponseFile   `json:"is_using_files"`
	Name              string                    `json:"name"`
	Summary           string                    `json:"summary"`
	Entities          []FileQueryResponseEntity `json:"entities"`
	Authors           []string                  `json:"authors"`
	LatestUpdate      string                    `json:"latestUpdate"`
	Dbid              string                    `json:"dbid"`
	RepoDbid          string                    `json:"repodbid"`
	RepoName          string                    `json:"reponame"`
	ImportPath        string                    `json:"importpath"`
	PackageDbid       string                    `json:"packagedbid"`
	PackageImportPath string                    `json:"packageimportpath"`
}

type EntityQueryResponseEntity struct {
	Name      string `json:"name"`
	Dbid      string `json:"dbid"`
	Count     int    `json:"count"`
	Signature string `json:"signature"`
	RepoDbid  string `json:"repodbid"`
	RepoName  string `json:"reponame"`
}

type EntityQueryResponseFile struct {
	Name      string `json:"name"`
	Summary   string `json:"summary"`
	Signature string `json:"signature"`
	Dbid      string `json:"dbid"`
}

type EntityQueryResponseResult struct {
	IsUsedByEntities    []EntityQueryResponseEntity `json:"is_used_by_entities"`
	IsUsingEntities     []EntityQueryResponseEntity `json:"is_using_entities"`
	Name                string                      `json:"name"`
	Summary             string                      `json:"summary"`
	Signature           string                      `json:"signature"`
	Authors             []string                    `json:"authors"`
	LatestUpdate        string                      `json:"latestUpdate"`
	Dbid                string                      `json:"dbid"`
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
	Dbid string `json:"dbid"`
}

type ListFilesResponseResult struct {
	Type       string `json:"type"`
	Dbid       string `json:"dbid"`
	ImportPath string `json:"importPath"`
	Files      []struct {
		ImportPath string `json:"importPath"`
		// Summary    string `json:"summary"`
		Dbid string `json:"dbid"`
		Name string `json:"name"`
	} `json:"files"`
}

type Stats struct {
	Repositories int `json:"repositories"`
	Directories  int `json:"directories"`
	Files        int `json:"files"`
	Codeblocks   int `json:"codeblocks"`
	// FileCommits      int                 `json:"fileCommits"`
	// Authors          int                 `json:"authors"`
	Relationships    int                 `json:"relationships"`
	MostDependedOn   StatsRepositoryInfo `json:"mostDependedOn"`
	MostDependencies StatsRepositoryInfo `json:"mostDependencies"`
}

type StatsRepositoryInfo struct {
	Name  string `json:"name"`
	Dbid  string `json:"dbid"`
	Count int    `json:"count"`
}
