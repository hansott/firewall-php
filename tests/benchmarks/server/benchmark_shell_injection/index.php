<?php
// Number of iterations
$iterations = 10000;

// Array to store execution times
$executionTimes = [];

// Loop through the iterations
for ($i = 0; $i < $iterations; $i++) {
    // Start time
    $startTime = microtime(true);

    // Your code here (e.g., some computation or function call)
    passthru("time");

    // End time
    $endTime = microtime(true);

    // Calculate the elapsed time in milliseconds
    $executionTime = ($endTime - $startTime) * 1000;

    // Store execution time in array
    $executionTimes[] = $executionTime;

    // Output the result for each iteration
    //echo "Iteration " . ($i + 1) . " - Execution time: " . $executionTime . " ms\n";
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
