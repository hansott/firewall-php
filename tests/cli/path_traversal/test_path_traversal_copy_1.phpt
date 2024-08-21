--TEST--
Test path traversal (copy)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = '../file';

$file = '../file/test.txt';
$dest = 'test2.txt';
    
copy($file, $dest);
    

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Path traversal detected.*
