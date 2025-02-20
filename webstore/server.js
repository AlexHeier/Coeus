const express = require("express");
const { Pool } = require("pg");
const cors = require("cors");

const app = express();
const port = 3000;

app.use(cors()); // Tillater forespørsler fra frontend
app.use(express.json()); // For JSON-data

// PostgreSQL tilkobling
const pool = new Pool({
  user: "coeus",
  host: "10.212.172.128",
  database: "coeus",
  password: "coeusPass101",
  port: 5432,
});

// Hent alle produkter
app.get("/products", async (req, res) => {
  try {
    const result = await pool.query("SELECT * FROM products.Product");
    res.json(result.rows);
  } catch (err) {
    console.error(err);
    res.status(500).send("Database error");
  }
});

// Hent alle produkter
app.get("/category", async (req, res) => {
    try {
      const result = await pool.query("SELECT * FROM products.category");
      res.json(result.rows);
    } catch (err) {
      console.error(err);
      res.status(500).send("Database error");
    }
  });

//henter produkter med produsent og kategori
app.get("/products/full", async (req, res) => {
    try {
      const result = await pool.query(`
        SELECT p.ID, p.Name AS ProductName, p.Description, p.Price, 
               c.Name AS Category, pr.Name AS Producer, p.ImageURL, p.quantity
        FROM products.Product p
        LEFT JOIN products.Category c ON p.CategoryID = c.ID
        LEFT JOIN products.Producer pr ON p.ProducerID = pr.ID
      `);
      res.json(result.rows);
    } catch (err) {
      console.error(err);
      res.status(500).send("Database error");
    }
  });
  

app.listen(port, () => {
  console.log(`Server kjører på http://localhost:${port}`);
});
