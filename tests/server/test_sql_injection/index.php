<?php
    
\aikido\set_user("12345", "Tudor");

try {
    $dbFile = 'unsafe_database.sqlite';

    $pdo = new PDO('sqlite:unsafe_database.sqlite');
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
    $pdo->exec("CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY, 
                name TEXT, 
                email TEXT)");

    $pdo->exec("INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com')");

    // Read the raw POST body
    $requestBody = file_get_contents('php://input');

    // Decode the JSON data to an associative array
    $data = json_decode($requestBody, true);

    // Vulnerable query
    $result = $pdo->query("SELECT * FROM users WHERE id = " . $data['userId']);

    echo "Query executed!";
} catch (PDOException $e) {
    echo "Connection failed: " . $e->getMessage();
}

// Close the database connection
$pdo = null;

// Delete the database file
if (file_exists($dbFile)) {
    unlink($dbFile);
}
?>
