<?php
    
\aikido\set_user("12345", "Tudor");

// Read the raw POST body
$requestBody = file_get_contents('php://input');

// Decode the JSON data to an associative array
$data = json_decode($requestBody, true);

if (isset($data['url'])) {
    $ch1 = curl_init($data['url']);
    curl_setopt($ch1, CURLOPT_RETURNTRANSFER, true);
    curl_exec($ch1);
    curl_close($ch1);
    echo "Got URL content!";
} else {
    echo "Field 'url' is not present in the JSON data.";
}

?>
