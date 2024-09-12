<?php
// Number of warmup iterations
$warmupIterations = 1000;

// Number of actual iterations
$iterations = 10000;

// Array to store execution times
$executionTimes = [];

// Warmup phase
for ($i = 0; $i < $warmupIterations; $i++) {
    passthru("time");
}

// Loop through the iterations
for ($i = 0; $i < $iterations; $i++) {
    // Start time
    $startTime = microtime(true);

    passthru("time");

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
echo "P50 (Median) execution time: " . $p50 . " ms\n";
?>
