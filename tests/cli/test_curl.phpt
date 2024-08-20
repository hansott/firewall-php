--TEST--
Test curl exec function

--ENV--
AIKIDO_LOG_LEVEL=INFO

--FILE--
<?php
$ch1 = curl_init("https://example.com/");
curl_setopt($ch1, CURLOPT_RETURNTRANSFER, true);
curl_exec($ch1);
curl_close($ch1);


$ch2 = curl_init("https://httpbin.org/get");
$queryParams = http_build_query([
    'param1' => 'value1',
    'param2' => 'value2'
]);
curl_setopt($ch2, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch2, CURLOPT_URL, "https://httpbin.org/get?" . $queryParams);
curl_exec($ch2);
curl_close($ch2);

$ch3 = curl_init();
$options = [
CURLOPT_URL => "https://facebook.com",
CURLOPT_RETURNTRANSFER => true,
CURLOPT_HEADER => false,
];
curl_setopt_array($ch3, $options);
curl_exec($ch3);
curl_close($ch3);

?>

--EXPECT--
[AIKIDO][INFO] [BEFORE] Got domain: example.com
[AIKIDO][INFO] [AFTER] Got domain: example.com port: 443
[AIKIDO][INFO] [BEFORE] Got domain: httpbin.org
[AIKIDO][INFO] [AFTER] Got domain: httpbin.org port: 443
[AIKIDO][INFO] [BEFORE] Got domain: facebook.com
[AIKIDO][INFO] [AFTER] Got domain: facebook.com port: 443
