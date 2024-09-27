<?php
    
\aikido\set_user("12345", "Tudor");

// Read the raw POST body
$requestBody = file_get_contents('php://input');

// Decode the JSON data to an associative array
$data = json_decode($requestBody, true);

if (isset($data['url'])) {
    file_get_contents($data['url']);
    echo "Got URL content!";
} else {
    echo "Field 'url' is not present in the JSON data.";
}

?>
