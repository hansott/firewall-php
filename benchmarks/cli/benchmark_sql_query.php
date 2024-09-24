<?php

$dbFile = 'unsafe_database.sqlite';

$pdo = new PDO('sqlite:unsafe_database.sqlite');
$pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
$pdo->exec("CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY, 
            name TEXT, 
            email TEXT)");

$pdo->exec("INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com')");

// Number of warmup iterations
$warmupIterations = 1000;

// Number of actual iterations
$iterations = 10000;

// Array to store execution times
$executionTimes = [];

// Warmup phase
for ($i = 0; $i < $warmupIterations; $i++) {
    $pdo->query("SELECT * FROM users WHERE id = 3");
}

// Loop through the iterations
for ($i = 0; $i < $iterations; $i++) {
    // Start time
    $startTime = microtime(true);

    $pdo->query("SELECT * FROM users WHERE id = 3");

    // End time
    $endTime = microtime(true);

    // Calculate the elapsed time in milliseconds
    $executionTime = ($endTime - $startTime) * 1000;

    // Store execution time in array
    $executionTimes[] = $executionTime;
}

// Function to calculate the median (P50)
function calculateMedian($arr) {
    sort($arr); // Sort the array
    $count = count($arr);
    $middle = floor(($count - 1) / 2);
    
    if ($count % 2) {
        // If odd, return the middle element
        return $arr[$middle];
    } else {
        // If even, return the average of the middle elements
        return ($arr[$middle] + $arr[$middle + 1]) / 2.0;
    }
}

// Calculate and print the median (P50)
$p50 = calculateMedian($executionTimes);

if (isset($argv[1])) {
    $filename = $argv[1];
    $file = fopen($filename, 'w');
    if ($file) {
        fwrite($file, "p50 - " . $p50 . " ms\n");
        fclose($file);
    } else {
        echo "Unable to open the file: $filename";
    }
} else {
    echo "Please provide a filename as the first argument.";
}

$pdo = null;

if (file_exists($dbFile)) {
    unlink($dbFile);
}

?>
