--TEST--
Test curl initialization and setting options

--INI--
extension=curl.so
extension=aikido.so
aikido.log_level=2

--SKIPIF--
<?php if (version_compare(PHP_VERSION, '8.0.0', '<')) die('skip PHP 8.0 or later required'); ?>

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
[AIKIDO][WARN][GO] AIKIDO_TOKEN not found in env variables!
[AIKIDO][INFO][GO] Got domain: example.com
[AIKIDO][INFO][GO] Got domain: www.facebook.com
[AIKIDO][INFO][GO] Got domain: www.google.com
[AIKIDO][INFO][GO] Got domain: example2.com
[AIKIDO][INFO][GO] Got domain: en.wikipedia.org