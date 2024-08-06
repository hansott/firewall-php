<?php
    http_response_code(429);
    header('Content-Type: text/plain');
    echo 'Test PHP script';
    exit;
?>