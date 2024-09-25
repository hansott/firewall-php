<?php
    $ch1 = curl_init("https://example.com/");
    curl_setopt($ch1, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch1, CURLOPT_TIMEOUT_MS, 1);
?>