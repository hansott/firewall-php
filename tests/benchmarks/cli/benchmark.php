<?php
// Number of warmup iterations
$warmupIterations = 1000;

// Number of actual iterations
$iterations = 10000;

// Array to store execution times
$executionTimes = [];

// Warmup phase
for ($i = 0; $i < $warmupIterations; $i++) {
    // <insert PHP code here>
}

// Loop through the iterations
for ($i = 0; $i < $iterations; $i++) {
    // Start time
    $startTime = microtime(true);

    // <insert PHP code here>

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
?>
