<?php


$skywallet_connect = "/usr/bin/skywallet_connect";
if (!file_exists($skywallet_connect))
        die("skywallet_connect not found");

if (!is_executable($skywallet_connect))
        die("skywallet_connect doesn't have exec permissions");


function getbalance() {
	 
	
	global $skywallet_connect;
	
	$memo = base64_encode($meta); 
	$cmd =  "$skywallet_connect balance";
	
	// Exec the binary
	$json = exec($cmd, $outarray, $error_code);
	if ($error_code != 0) {
			echo "Invalid response from skywallet_connect: $error_code, Output $json";
			return 1;
	}
	
	$arr = json_decode($json, true);
	if (!$arr) {
			echo "Failed to decode json: $json";
			return 1;
	}
	print_r($arr);
	print_r($json);
	die();
}

	$result = getbalance(); 
