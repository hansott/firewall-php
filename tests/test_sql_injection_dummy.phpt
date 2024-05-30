--TEST--
Test SQLite database operations

--INI--
extension=aikido.so
aikido.log_level=1
aikido.blocking=1

--FILE--
<?php
try {
    $dbFile = 'unsafe_database.sqlite';

    $pdo = new PDO('sqlite:unsafe_database.sqlite');
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

    $pdo->exec("CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY, 
                name TEXT, 
                email TEXT)");

    $pdo->exec("INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com')");

    // Simulate user input
    $unsafeInput = "1 OR 1=1";

    // Vulnerable query
    $result = $pdo->query("SELECT * FROM users WHERE id = $unsafeInput");

    foreach ($result as $row) {
        echo "ID: " . $row['id'] . "\n";
        echo "Name: " . $row['name'] . "\n";
        echo "Email: " . $row['email'] . "\n\n";
    }
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

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Sql injection detected.*