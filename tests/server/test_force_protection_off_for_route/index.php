<?php
    
\aikido\set_user("12345", "Tudor");

$requestBody = file_get_contents('php://input');

$data = json_decode($requestBody, true);

fopen($data['folder'] . '/file', 'r');
passthru('binary --domain www.example' .  $data['command'] . '.com');
file_get_contents($data['url']);

try {
    $dbFile = 'unsafe_database_force_protection_off.sqlite';

    $pdo = new PDO('sqlite:unsafe_database.sqlite');
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
    $pdo->exec("CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY, 
                name TEXT, 
                email TEXT)");

    $pdo->exec("INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com')");

    // Vulnerable query
    $result = $pdo->query("SELECT * FROM users WHERE id = " . $data['userId']);

} catch (PDOException $e) {

}

// Close the database connection
$pdo = null;

// Delete the database file
if (file_exists($dbFile)) {
    unlink($dbFile);
}

?>
