--TEST--
Test path traversal (move_uploaded_file)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = '../file';

$file = '../file/test.txt';
$dest = 'test2.txt';
    
move_uploaded_file($file, $dest);
    

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Path traversal detected.*
