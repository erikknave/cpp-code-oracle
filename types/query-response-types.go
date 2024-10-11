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
	Directories   []RepoQueryResponseDirectory  `json:"directories"`
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

type CodeblockQueryResponseCodeblock struct {
	Name          string `json:"name"`
	Dbid          int    `json:"dbid"`
	Count         int    `json:"count"`
	CodeblockType string `json:"type"`
	Signature     string `json:"signature"`
	RepoDbid      int    `json:"repodbid"`
	RepoName      string `json:"reponame"`
}

type EntityQueryResponseFile struct {
	Name      string `json:"name"`
	Summary   string `json:"summary"`
	Signature string `json:"signature"`
	Dbid      int    `json:"dbid"`
}

type CodeblockQueryResponseResult struct {
	IsUsedByCodeblocks    []CodeblockQueryResponseCodeblock `json:"is_used_by_codeblocks"`
	IsUsingCodeblocks     []CodeblockQueryResponseCodeblock `json:"is_using_codeblocks"`
	Name                  string                            `json:"name"`
	Summary               string                            `json:"summary"`
	CodeblockType         string                            `json:"type"`
	Signature             string                            `json:"signature"`
	Authors               []string                          `json:"authors"`
	LatestUpdate          string                            `json:"latestUpdate"`
	Dbid                  int                               `json:"dbid"`
	RepoDbid              int                               `json:"repodbid"`
	RepoName              string                            `json:"reponame"`
	RepoShortSummary      string                            `json:"reposhortsummary"`
	ImportPath            string                            `json:"importpath"`
	DirectoryDbid         int                               `json:"directorydbid"`
	DirectoryImportPath   string                            `json:"directoryimportpath"`
	DirectoryShortSummary string                            `json:"directoryshortsummary"`
	FileDbid              int                               `json:"filedbid"`
	FileName              string                            `json:"filename"`
	FileSummary           string                            `json:"filesummary"`
	Content               string                            `json:"content"`
	StartOffSet           int                               `json:"startoffset"`
	EndOffSet             int                               `json:"endoffset"`
}

type ContainerQueryResponseResult struct {
	IsUsedByContainers   []ContainerQueryResponseContainer `json:"is_used_by_containers"`
	IsUsingContainers    []ContainerQueryResponseContainer `json:"is_using_containers"`
	ParentContainer      ContainerQueryResponseContainer   `json:"parentContainer"`
	GrandParentContainer ContainerQueryResponseContainer   `json:"grandParentContainer"`
	ChildContainers      []ContainerQueryResponseContainer `json:"childContainers"`
	Name                 string                            `json:"name"`
	Summary              string                            `json:"summary"`
	Signature            string                            `json:"signature"`
	ContainerType        string                            `json:"type"`
	Authors              []string                          `json:"authors"`
	LatestUpdate         string                            `json:"latestUpdate"`
	Dbid                 int                               `json:"dbid"`
	RepoDbid             int                               `json:"repodbid"`
	RepoName             string                            `json:"reponame"`
	RepoShortSummary     string                            `json:"reposhortsummary"`
}

type ContainerQueryResponseContainer struct {
	Name          string `json:"name"`
	Signature     string `json:"signature"`
	Summary       string `json:"summary"`
	Dbid          int    `json:"dbid"`
	Count         int    `json:"count"`
	ContainerType string `json:"type"`
}

type ContainerQueryResponseCodeblock struct {
	Name          string `json:"name"`
	Signature     string `json:"signature"`
	Dbid          int    `json:"dbid"`
	CodeblockType string `json:"type"`
	RepoDbid      int    `json:"repodbid"`
	RepoName      string `json:"reponame"`
}

type ContainerAgentQueryResponseResult struct {
	Codeblocks    []ContainerAgentQueryResponseCodeblock `json:"codeblocks"`
	Name          string                                 `json:"name"`
	Summary       string                                 `json:"summary"`
	Signature     string                                 `json:"signature"`
	ContainerType string                                 `json:"type"`
	Dbid          int                                    `json:"dbid"`
	RepoDbid      int                                    `json:"repodbid"`
	RepoName      string                                 `json:"reponame"`
	RepoSummary   string                                 `json:"reposhortsummary"`
}

type ContainerAgentQueryResponseCodeblock struct {
	Signature      string `json:"signature"`
	Dbid           int    `json:"dbid"`
	FileImportPath string `json:"fileimportpath"`
	FileDbid       int    `json:"filedbid"`
	FileSummary    string `json:"filesummary"`
}

type RepoListQueryResult struct {
	Name    string `json:"name"`
	Dbid    int    `json:"dbid"`
	Summary string `json:"summary"`
}

type ListFilesResponseResult struct {
	Type       string `json:"type"`
	Dbid       int    `json:"dbid"`
	ImportPath string `json:"importPath"`
	Files      []struct {
		ImportPath string `json:"importPath"`
		// Summary    string `json:"summary"`
		Dbid int    `json:"dbid"`
		Name string `json:"name"`
	} `json:"files"`
}

type Stats struct {
	Repositories              int                 `json:"repositories"`
	Directories               int                 `json:"directories"`
	Files                     int                 `json:"files"`
	Codeblocks                int                 `json:"codeblocks"`
	Containers                int                 `json:"containers"`
	FileCommits               int                 `json:"fileCommits"`
	Authors                   int                 `json:"authors"`
	Relationships             int                 `json:"relationships"`
	MostDependedOn            StatsRepositoryInfo `json:"mostDependedOn"`
	MostDependencies          StatsRepositoryInfo `json:"mostDependencies"`
	MostDependedOnContainer   StatsContainerInfo  `json:"mostDependedOnContainer"`
	MostDependenciesContainer StatsContainerInfo  `json:"mostDependenciesContainer"`
}

type StatsRepositoryInfo struct {
	Name  string `json:"name"`
	Dbid  int    `json:"dbid"`
	Count int    `json:"count"`
}

type StatsContainerInfo struct {
	Name          string `json:"name"`
	Signature     string `json:"signature"`
	ContainerType string `json:"type"`
	Dbid          int    `json:"dbid"`
	Count         int    `json:"count"`
}
