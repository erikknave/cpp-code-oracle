package chromaclient

import (
	"log"
	"os"

	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

var FileStore chroma.Store
var RepositoryStore chroma.Store
var ModuleStore chroma.Store
var PackageStore chroma.Store
var EntityStore chroma.Store

func Init() {
	var err error
	hostMode := os.Getenv("AI_HOST_MODE")

	var embedder *embeddings.EmbedderImpl

	if hostMode == "azure" {
		llm, err := openai.New(
			openai.WithAPIType(openai.APITypeAzure),
			openai.WithModel(os.Getenv("EMBEDDING_MODEL")),
			openai.WithEmbeddingModel(os.Getenv("EMBEDDING_MODEL")),
			openai.WithBaseURL(os.Getenv("AZURE_BASE_URL")),
			openai.WithToken(os.Getenv("OPENAI_API_KEY")),
		)
		if err != nil {
			log.Fatal(err)
		}

		embedder, err = embeddings.NewEmbedder(llm)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		llm, err := openai.New(
			openai.WithAPIType(openai.APITypeOpenAI),
			// openai.WithModel(os.Getenv("EMBEDDING_MODEL")),
			// openai.WithEmbeddingModel(os.Getenv("EMBEDDING_MODEL")),
			openai.WithBaseURL(os.Getenv("AZURE_BASE_URL")),
			openai.WithToken(os.Getenv("OPENAI_API_KEY")),
		)
		if err != nil {
			log.Fatal(err)
		}

		embedder, err = embeddings.NewEmbedder(llm)
		if err != nil {
			log.Fatal(err)
		}

	}

	FileStore, err = chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace("File"),
	)
	if err != nil {
		log.Fatalf("new: %v\n", err)
	}
	ModuleStore, err = chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace("Module"),
	)
	if err != nil {
		log.Fatalf("new: %v\n", err)
	}
	RepositoryStore, err = chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace("Repository"),
	)
	if err != nil {
		log.Fatalf("new: %v\n", err)
	}

	PackageStore, err = chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace("Package"),
	)
	if err != nil {
		log.Fatalf("new: %v\n", err)
	}
	EntityStore, err = chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace("Entity"),
	)
	if err != nil {
		log.Fatalf("new: %v\n", err)
	}
}
