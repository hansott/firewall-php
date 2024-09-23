
$dbFile = 'unsafe_database.sqlite';

$pdo = new PDO('sqlite:unsafe_database.sqlite');
$pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
$pdo->exec("CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY, 
            name TEXT, 
            email TEXT)");

$pdo->exec("INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com')");
