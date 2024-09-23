--TEST--
Test SQLite database operations

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

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
    $_SERVER['HTTP_USER'] = $unsafeInput;

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
.*Fatal error: Uncaught Exception: Aikido firewall has blocked an SQL injection.*