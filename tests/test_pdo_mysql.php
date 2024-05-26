<?php
$host = '127.0.0.1';
$db   = 'testdb';
$user = 'username';  // Replace with your MySQL username
$pass = 'password';  // Replace with your MySQL password
$charset = 'utf8mb4';

$dsn = "mysql:host=$host;dbname=$db;charset=$charset";
$options = [
    PDO::ATTR_ERRMODE            => PDO::ERRMODE_EXCEPTION,
    PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
    PDO::ATTR_EMULATE_PREPARES   => false,
    PDO::ATTR_TIMEOUT            => 30,
    PDO::ATTR_PERSISTENT         => true,  // Enable persistent connection
];

try {
    $pdo = new PDO($dsn, $user, $pass, $options);
    echo "Connected successfully!<br>";

    // Simple query
    $sql = 'SELECT id, name, email FROM users';
    //$stmt = $pdo->query($sql);

} catch (\PDOException $e) {
    throw new \PDOException($e->getMessage(), (int)$e->getCode());
}
?>
