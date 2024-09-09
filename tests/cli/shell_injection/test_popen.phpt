--TEST--
Test PHP shell injection (popen)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = 'www.example`whoami`.com';

$handle = popen('binary --domain www.example`whoami`.com', 'r');
while (!feof($handle)) {
    $buffer = fgets($handle);
    echo $buffer;
}
pclose($handle);

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Shell injection detected in.*