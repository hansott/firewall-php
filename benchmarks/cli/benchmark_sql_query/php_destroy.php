$pdo = null;

if (file_exists($dbFile)) {
    unlink($dbFile);
}