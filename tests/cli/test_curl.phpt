--TEST--
Test curl initialization and setting options

--ENV--
AIKIDO_LOG_LEVEL=INFO

--FILE--
<?php
$ch1 = curl_init("https://example.com/");
$ch2 = curl_init("http://www.facebook.com");
$ch3 = curl_init("https://www.google.com/search?q=test");
$ch4 = curl_copy_handle($ch1);
curl_setopt($ch2, CURLOPT_URL, "https://example2.com");
$ch5 = curl_copy_handle($ch1);
curl_setopt($ch2, CURLOPT_URL, "https://en.wikipedia.org/wiki/Runtime_library");
?>

--EXPECT--
[AIKIDO][INFO] Got domain: example.com
[AIKIDO][INFO] Got domain: www.facebook.com
[AIKIDO][INFO] Got domain: www.google.com
[AIKIDO][INFO] Got domain: example2.com
[AIKIDO][INFO] Got domain: en.wikipedia.org