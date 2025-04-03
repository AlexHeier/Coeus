package coeus

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	pg "github.com/lib/pq" // Import the PostgreSQL driver
)

// Struct for containing the database connection information
var sqlConfig struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

// ragfolder is the folder where the RAG files are stored
var ragfolder string = "./RAG"

// fileUpdateTime is a map that stores the last update time of each file
var fileUpdateTime = make(map[string]time.Time)

// rag is a boolean that indicates if RAG is enabled
var rag bool = false

// RAG config with default values
var closest int = 2        // closest is the number of closest results to return
var cs int = 300           // Chunk size for splitting the text into smaller chunks
var overlap float64 = 0.25 // Overlap ratio for splitting the text into smaller chunks
var mp float64 = 2.0       // Multiplier for the vector scaling

// db is the database connection
var db *sql.DB

// wordWeights is a map that stores the weights of each word in the database based on frequency within the RAG folder
// This is used to determine the importance of each word in the context of the RAG
var wordWeights map[string]float64
var totalWords int = 0 // Total number of words in the RAG folder

var mu sync.RWMutex // A mutex to handle concurrent access to the wordWeights map

var maxFrequency int = 0 // Maximum frequency of any word in the database

type vectorFrequency struct {
	Vector   []float64 // The vector representation of the word
	WordFreq int       // The frequency of the word in the database
}

/*
RAGConfig sets the configuration for the RAG (Retrieval-Augmented Generation) mode.

@param context: The number of closest results to use as context for the model. Default is 2.

@param chunkSize: The size of the chunks to split the text into. Default is 300.

@param overlapRatio: The ratio of overlap between chunks. Default is 0.25.

@param multiplier: The multiplier for the vector scaling. Default is 2.

@param folder: The folder where the RAG files are stored. Default is "./RAG".

@param error: An error if any of the fields are invalid.
*/
func RAGConfig(context, chunkSize int, overlapRatio, multiplier float64) error {
	if context < 1 {
		return fmt.Errorf("context must be greater than 0")
	}
	if chunkSize < 1 {
		return fmt.Errorf("chunk size must be greater than 0")
	}
	if overlapRatio < 0 || overlapRatio > 1 {
		return fmt.Errorf("overlap ratio must be between 0 and 1")
	}

	closest = context
	cs = chunkSize
	overlap = overlapRatio
	mp = multiplier
	return nil
}

/*
EnableRAG enables RAG (Retrieval-Augmented Generation) mode.
This mode allows the model to use external knowledge sources to improve its responses.

Recommended to use RAGConfig before Enable to remove racecondition between start of tokenization and config.

@param host: The host of the database.

@param dbname: The name of the database.

@param user: The user to connect to the database.

@param password: The password to connect to the database.

@param port: The port to connect to the database.

@return An error if any of the fields are empty or invalid.
*/
func EnableRAG(host, dbname, user, password string, port int, folder string) error {

	info, err := os.Stat(folder)
	if err != nil {
		return fmt.Errorf("folder does not exist")
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	rag = true

	if host == "" || dbname == "" || user == "" || password == "" {
		return fmt.Errorf("all fields need to be populated")
	}

	if port < 1000 || port > 65000 {
		return fmt.Errorf("port outside accepted boundry")
	}

	sqlConfig.host = host
	sqlConfig.dbname = dbname
	sqlConfig.user = user
	sqlConfig.password = password
	sqlConfig.port = port
	ragfolder = folder

	db, err = sql.Open("postgres", ragDBConnection())
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	err = db.QueryRow("SELECT word_frequency FROM embeddings ORDER BY word_frequency DESC LIMIT 1").Scan(&maxFrequency)
	if err != nil {
		log.Fatalf("Error fetching max frequency: %v", err)
	}

	go updateCheck()
	return nil
}

/*
ragDBConnection returns the connection string for the RAG database.
It is used to connect to the database and retrieve the data.
*/
func ragDBConnection() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", sqlConfig.host, sqlConfig.port, sqlConfig.user, sqlConfig.password, sqlConfig.dbname)
}

/*
GetRAG retrieves the closests chunks from the RAG database based of the user prompt.
*/
func GetRAG(userPrompt string) string {
	vector := chunkToVector(userPrompt)
	if vector == nil {
		return "Unable to use RAG"
	}

	query := `
	SELECT chunk 
	FROM rag
	ORDER BY vector <=> $1 
	LIMIT $2;
	`
	vecStr := "[" + strings.Trim(strings.Replace(fmt.Sprint(vector), " ", ",", -1), "[]") + "]"
	rows, err := db.Query(query, vecStr, closest)
	if err != nil {
		fmt.Println("\nGetRAG Error executing query:", err)
		return "Unable to use RAG"
	}
	defer rows.Close()

	var chunks []string
	for rows.Next() {

		var chunk string
		if err := rows.Scan(&chunk); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		chunks = append(chunks, chunk)
	}

	return strings.Join(chunks, "\n")
}

/*
updateCheck checks for updates in the RAG folder every 5 minutes.
If a file is updated, it updates the RAG.
*/
func updateCheck() {
	for {
		files, err := os.ReadDir(ragfolder)
		if err != nil {
			continue
		}
		for _, file := range files {
			filePath := ragfolder + "/" + file.Name()

			info, err := os.Stat(filePath)
			if err != nil {
				continue
			}

			updateTime := info.ModTime()
			if lastUpdate, exists := fileUpdateTime[filePath]; !exists || updateTime.After(lastUpdate) {
				go updateRAG(filePath)
			}
		}
		time.Sleep(15 * time.Minute)
	}
}

func fileToTokenAndFrequency(filePath string) ([]string, error) {
	if wordWeights == nil {
		wordWeights = make(map[string]float64)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	re := regexp.MustCompile(`[^\w\s]`)

	var rawChunks []string
	scanner := bufio.NewScanner(file)

	// Read and clean each line
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "\r", " ") // Remove carriage returns
		line = re.ReplaceAllString(line, "")       // Remove special characters
		line = strings.ToLower(line)               // Convert to lowercase
		line = strings.TrimSpace(line)             // Trim spaces

		for _, word := range strings.Fields(line) {
			mu.Lock()
			wordWeights[word]++
			totalWords++
			mu.Unlock()
		}

		if line != "" {
			rawChunks = append(rawChunks, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Process the chunks to ensure they are ≤ 200 characters with overlap
	finalChunks := splitChunks(rawChunks, cs, overlap)

	return finalChunks, nil
}

/*
splitChunks splits the chunks into smaller chunks of a specified maximum length.
It also ensures that the chunks overlap by a specified ratio.
*/
func splitChunks(chunks []string, maxLength int, overlapRatio float64) []string {
	var result []string
	overlap := int(float64(maxLength) * overlapRatio)

	for _, chunk := range chunks {
		if len(chunk) <= maxLength {
			result = append(result, chunk)
			continue
		}

		start := 0
		for start < len(chunk) {
			end := start + maxLength
			if end > len(chunk) {
				end = len(chunk)
			}
			result = append(result, chunk[start:end])

			// Stop if at the end
			if end == len(chunk) {
				break
			}

			// Move start forward with overlap
			start = end - overlap
		}
	}
	return result
}

func chunkToVector(chunk string) []float64 {
	query := `SELECT vector, word_frequency FROM embeddings WHERE word = ANY($1)`
	words := strings.Fields(chunk)

	rows, err := db.Query(query, pg.Array(words))
	if err != nil {
		fmt.Println("chunkToVector Error executing query:", err)
		return nil
	}
	defer rows.Close()

	var vf []vectorFrequency
	for rows.Next() {
		var v []float64
		var f int
		if err := rows.Scan(pg.Array(&v), &f); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		vf = append(vf, vectorFrequency{Vector: v, WordFreq: f})
	}

	return avrgVectorsLogScaled(vf)
}

func avrgVectorsLogScaled(vf []vectorFrequency) []float64 {
	if len(vf) == 0 {
		return []float64{} // Return empty slice if no vectors
	}

	vecSize := len(vf[0].Vector)
	sum := make([]float64, vecSize)

	// Add the vectors scaled logarithmically by their frequencies
	for _, v := range vf {
		if len(v.Vector) != vecSize {
			panic("all vectors must have the same length")
		}
		// Logarithmic scaling of frequency
		logScaledFrequency := math.Log(1 + float64(v.WordFreq)) // Adding 1 to avoid log(0)
		// Normalize with respect to the max frequency
		scale := logScaledFrequency / math.Log(1+float64(maxFrequency)) // Normalize by the max frequency

		for j, val := range v.Vector {
			sum[j] += val * scale * mp // Scale the vector by the log-scaled frequency and multiplier
		}
	}

	for i := range sum {
		sum[i] /= float64(len(vf)) // Average the vectors
	}

	return sum
}

func updateRAG(filePath string) {
	chunks, err := fileToTokenAndFrequency(filePath)
	if err != nil {
		fmt.Print("Error reading file:", err)
		return
	}

	query := `DELETE FROM rag WHERE file_location = $1;`
	_, err = db.Exec(query, filePath)
	if err != nil {
		fmt.Printf("error deleting old data: %v", err)
	}

	for _, chunk := range chunks {
		chunkVec := chunkToVector(chunk)
		if chunkVec == nil {
			continue
		}
		query := `INSERT INTO rag (file_location, chunk, vector) VALUES ($1, $2, $3);`
		vecStr := "[" + strings.Trim(strings.Replace(fmt.Sprint(chunkVec), " ", ",", -1), "[]") + "]"
		_, err = db.Exec(query, filePath, chunk, vecStr)
		if err != nil {
			fmt.Printf("error inserting data: %v", err)
			return
		}
	}
}

func getWordFrequency(word string) float64 {
	// Lock the wordWeights map for reading
	mu.RLock()
	defer mu.RUnlock()

	return wordWeights[word]
}

/*
CREATE TABLE rag (
  chunk TEXT PRIMARY KEY,  -- Text chunk as the unique identifier
  file_location TEXT,      -- Path of the file
  vector VECTOR(300)       -- 300-dimensional vector representing the chunk
);
*/
