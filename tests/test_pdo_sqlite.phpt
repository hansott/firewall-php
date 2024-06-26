--TEST--
Test SQLite database operations

--INI--
extension=aikido.so
aikido.log_level=2 

--FILE--
<?php
try {
    $dbFile = 'my_database.sqlite';

    $pdo = new PDO('sqlite:my_database.sqlite');

    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

    $pdo->exec("CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY, 
                name TEXT, 
                email TEXT)");

    $pdo->exec("INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com')");

    $result = $pdo->query('SELECT * FROM users');

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
--EXPECTF--
[AIKIDO][WARN][GO] AIKIDO_TOKEN not found in env variables!
[AIKIDO][INFO][GO] Got PDO query: SELECT * FROM users
ID: %d
Name: John Doe
Email: john@example.com