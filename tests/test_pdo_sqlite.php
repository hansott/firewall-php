<?php
try {
    // Create (connect to) SQLite database in file
    $pdo = new PDO('sqlite:my_database.sqlite');

    // Set errormode to exceptions
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

    // Create a table
    $pdo->exec("CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY, 
                name TEXT, 
                email TEXT)");

    // Insert a row of data
    $pdo->exec("INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com')");

    // Query the database and fetch data
    $result = $pdo->query('SELECT * FROM users');

    foreach ($result as $row) {
        echo "ID: " . $row['id'] . "<br>";
        echo "Name: " . $row['name'] . "<br>";
        echo "Email: " . $row['email'] . "<br><br>";
    }
} catch (PDOException $e) {
    echo "Connection failed: " . $e->getMessage();
}
?>