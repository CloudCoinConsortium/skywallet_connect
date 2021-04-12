<?php
/*
Allows you to command your Skywallet transfer your coins to another Skywallet account. 

Open Source
Free to use with all errors, flaws, bugs etc. 
Last Changed: April 10, 2021
*/

ini_set('display_errors', '1');
ini_set('display_startup_errors', '1');
error_reporting(E_ALL);

$skywallet_connect = "/usr/bin/skywallet_connect";
if (!file_exists($skywallet_connect))
        die("skywallet_connect not found");

if (!is_executable($skywallet_connect))
        die("skywallet_connect doesn't have exec permissions");


	function getbalance($meta ,$merchant , $amount ) {
		 
		
		global $skywallet_connect;
		
		$memo = base64_encode($meta); 
		
		$cmd =  "$skywallet_connect transfer $amount $merchant $memo";
		
		// Exec the binary
		$json = exec($cmd, $outarray, $error_code);
		
		if ($error_code == 0) {
			$res['trx'] = time();			
			$res['status'] = 1;			
			$res['message'] = 'CloudCoins sent';
			return $res;	
		}
		if ($error_code != 0) {
			$res['status'] = 0;			
			$res['message'] = 'Amount, To, Memoparameters required: /usr/bin/skywallet_connect transfer 250 destination.skywallet.cc memo';			
			//echo "Invalid response from skywallet_connect: $error_code, Output $json";
			//return 1;
			return $res;
		}	
		/* $arr = json_decode($json, true);
		if (!$arr) {
				echo "Failed to decode json: $json";
				
				return 1;
		} */
		
	} 
	$out = @file_get_contents('php://input');	
	$event_json = json_decode('['.$out.']');
	isset($event_json[0]->key) ? $key = $event_json[0]->key : $key= '';	
		
	if(!empty($key)){
		if(md5($key) == 'The 32 digit hexidecimal hash of your key/password goes here'){		
			isset($event_json[0]->meta) ? $meta = $event_json[0]->meta : $meta= '';	
			isset($event_json[0]->merchant) ? $merchant = $event_json[0]->merchant : $merchant= '';	
			isset($event_json[0]->amount) ? $amount = $event_json[0]->amount : $amount= 0;
			if(!empty($merchant)){				
				if(!empty($meta)){	
					if(!empty($meta) && !empty($merchant) && !empty($amount)){	
						$data = getbalance($meta ,$merchant , $amount ); 			
					}else{
						$data['status'] = 0;
						$data['message'] = 'unauthorized access';
					}	
				}else{
					$data['status'] = 0;
					$data['message'] = 'meta or memo id is required';
					
				}	
			}else{
				$data['status'] = 0;
				$data['message'] = 'merchant address is required';
				
			}			
		}else{
			$data['status'] = 0;
			$data['message'] = 'unauthorized access';
		}
	}else{
		$data['status'] = 0;
		$data['message'] = 'auth key key is required';
	}

	echo json_encode($data);
