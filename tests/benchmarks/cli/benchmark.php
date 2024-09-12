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
    $startTime = hrtime(true);

    // <insert PHP code here>

    // End time
    $endTime = hrtime(true);

    // Calculate the elapsed time in milliseconds
    $executionTime = ($endTime - $startTime);

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
echo "p50 - " . $p50 . " ns";
?>
