<?php

$ch = curl_init("https://example.com/");
$ch2 = curl_init();

curl_setopt($ch2, CURLOPT_URL, "https://google.com/");
curl_setopt($ch2, CURLOPT_HEADER, 0);

?>
