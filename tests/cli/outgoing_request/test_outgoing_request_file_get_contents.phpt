--TEST--
Test outgoing request (file_get_contents)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

--FILE--
<?php
    
file_get_contents("http://www.example.com");

?>

--EXPECT--
[AIKIDO][INFO] [BEFORE] Got domain: www.example.com
[AIKIDO][INFO] [AFTER] Got domain: www.example.com port: 80