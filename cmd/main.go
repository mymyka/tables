package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"tables/internal/builder"
	"tables/internal/parser"
	"tables/internal/writer"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var (
	dbConnectionString string
	outputPath         string
)

var rootCmd = &cobra.Command{
	Use:   "datatypes",
	Short: "Generate Go types from PostgreSQL database schema",
	Long: `A CLI tool that connects to a PostgreSQL database, reads the schema,
and generates Go type definitions for each table with proper type mappings.`,
	Run: func(cmd *cobra.Command, args []string) {
		if dbConnectionString == "" {
			log.Fatal("Database connection string is required. Use --db flag or set DB_CONNECTION_STRING environment variable.")
		}

		if outputPath == "" {
			log.Fatal("Output path is required. Use --output flag.")
		}

		generateTypes()
	},
}

func init() {
	// Add flags
	rootCmd.Flags().StringVarP(&dbConnectionString, "db", "d", "", "PostgreSQL connection string (required)")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output directory path (required)")

	// Mark flags as required
	rootCmd.MarkFlagRequired("db")
	rootCmd.MarkFlagRequired("output")

	// Allow environment variable for db connection
	if envDB := os.Getenv("DB_CONNECTION_STRING"); envDB != "" {
		dbConnectionString = envDB
	}
}

func generateTypes() {
	fmt.Printf("Connecting to database...\n")

	// Connect to database
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Printf("Connected successfully!\n")
	fmt.Printf("Parsing database schema...\n")

	inspector := parser.NewSchemaParser(db)

	tables, err := inspector.GetTables()
	if err != nil {
		log.Fatal("Failed to get tables:", err)
	}

	fmt.Printf("Found %d tables\n", len(tables))
	fmt.Printf("Generating Go types...\n")

	block := builder.Build(tables)

	fmt.Printf("Writing files to %s...\n", outputPath)

	err = writer.Write(outputPath, block)
	if err != nil {
		log.Fatal("Failed to write files:", err)
	}

	fmt.Printf("Successfully generated types for %d tables!\n", len(tables))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
