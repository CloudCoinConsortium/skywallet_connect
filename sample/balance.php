<?php


$raida_go = "/usr/bin/raida_go";
if (!file_exists($raida_go))
        die("Raida Go not found");

if (!is_executable($raida_go))
        die("Raida Go doesn't have exec permissions");


function getbalance() {
	 
	
	global $raida_go;
	
	$memo = base64_encode($meta); 
	$cmd =  "$raida_go balance";
	
	// Exec the binary
	$json = exec($cmd, $outarray, $error_code);
	if ($error_code != 0) {
			echo "Invalid response from raida_go: $error_code, Output $json";
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
