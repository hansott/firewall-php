<?php

$requestBody = file_get_contents('php://input');
$data = json_decode($requestBody, true);

if (extension_loaded('aikido')) {
    \aikido\set_user($data['userId'], $data['userName']);

    $decision = \aikido\should_block_request();

    // If the rate limit is exceeded, return a 429 status code
    if ($decision->block && $decision->type == "ratelimited") {
        http_response_code(429);
        echo "Rate limit exceeded";
        exit();
    }
}

// Continue handling the request
echo "Request successful!";

?>
