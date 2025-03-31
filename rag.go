package coeus

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
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
const ragfolder string = "./rag"

// fileUpdateTime is a map that stores the last update time of each file
var fileUpdateTime = make(map[string]time.Time)

// rag is a boolean that indicates if RAG is enabled
var rag bool = false

// closest is the number of closest results to return
var closest int = 1

// db is the database connection
var db *sql.DB

/*
EnableRAG enables RAG (Retrieval-Augmented Generation) mode.
This mode allows the model to use external knowledge sources to improve its responses.

@param host: The host of the database.

@param dbname: The name of the database.

@param user: The user to connect to the database.

@param password: The password to connect to the database.

@param port: The port to connect to the database.

@param number: The number of closest results to return.

@return An error if any of the fields are empty or invalid.
*/
func EnableRAG(host, dbname, user, password string, port, number int) error {
	closest = number
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

	var err error
	db, err = sql.Open("postgres", ragDBConnection())
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
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
getRAG retrieves the closests chunk from the RAG database based on the user prompt.
*/
func getRAG(userPrompt string) string {
	vector := chunkToVector(userPrompt)
	if vector == nil {
		return "Unable to use RAG"
	}

	query := `
	SELECT chunk 
	FROM rag
	ORDER BY vector <-> $1 
	LIMIT $2;
	`
	rows, err := db.Query(query, pg.Array(vector), closest)
	if err != nil {
		fmt.Println("Error executing query:", err)
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

	fmt.Printf("\n\nChunks: %v\n\n", chunks)

	return strings.Join(chunks, "\n")
}

/*
updateCheck checks for updates in the RAG folder every 5 minutes.
If a file is updated, it updates the RAG.
*/
func updateCheck() {
	for {
		time.Sleep(15 * time.Minute)
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
	}
}

func fileToToken(filePath string) ([]string, error) {
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

		if line != "" {
			rawChunks = append(rawChunks, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Process the chunks to ensure they are ≤ 200 characters with overlap
	finalChunks := splitChunks(rawChunks, 200, 0.3)

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
	query := `SELECT vector FROM embeddings WHERE word = ANY($1)`
	words := strings.Fields(chunk)

	rows, err := db.Query(query, words)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil
	}
	defer rows.Close()

	var vector [][]float64
	for rows.Next() {
		var v []float64
		if err := rows.Scan(pg.Array(&v)); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		vector = append(vector, v)
	}

	vec := sumVectors(vector)
	return vec
}

func sumVectors(vectors [][]float64) []float64 {
	if len(vectors) == 0 {
		return nil
	}

	sum := make([]float64, len(vectors[0]))

	for _, vector := range vectors {
		for i, val := range vector {
			sum[i] += val
		}
	}

	return sum
}

func updateRAG(filePath string) {
	chunks, err := fileToToken(filePath)
	if err != nil {
		fmt.Print("Error reading file:", err)
		return
	}

	query := `DELETE FROM rag WHERE file_location = $1;`
	_, err = db.Exec(query, filePath)
	if err != nil {
		fmt.Printf("error deleting old data: %v", err)
		return
	}

	for _, chunk := range chunks {
		chunkVec := chunkToVector(chunk)
		if chunkVec == nil {
			continue
		}
		query := `INSERT INTO rag (file_location, chunk, vector) VALUES ($1, $2, $3);`
		_, err = db.Exec(query, filePath, chunk, pg.Array(chunkVec))
		if err != nil {
			fmt.Printf("error inserting data: %v", err)
			return
		}
		fmt.Printf("\n\nInserted chunk: %s\n", chunk)
	}
}

/*
CREATE TABLE rag (
  chunk TEXT PRIMARY KEY,  -- Text chunk as the unique identifier
  file_location TEXT,      -- Path of the file
  vector VECTOR(300)       -- 300-dimensional vector representing the chunk
);
*/
