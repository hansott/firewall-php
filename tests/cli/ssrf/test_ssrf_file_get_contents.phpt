--TEST--
Test path ssrf (file_get_contents)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = 'http://127.0.0.1:8081';

$file = 'http://127.0.0.1:8081';
    
file_get_contents($file);

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Aikido firewall has blocked a server-side request forgery.*
