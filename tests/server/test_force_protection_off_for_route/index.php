<?php
    
\aikido\set_user("12345", "Tudor");

$requestBody = file_get_contents('php://input');

$data = json_decode($requestBody, true);

fopen($data['folder'] . '/file', 'r');
passthru('binary --domain www.example' .  $data['command'] . '.com');
file_get_contents($data['url']);

?>
