<?php
 ini_set('display_errors', '1');
ini_set('display_startup_errors', '1');
error_reporting(E_ALL); 


$raida_go = "/usr/bin/raida_go";
if (!file_exists($raida_go))
        die("Raida Go not found");

if (!is_executable($raida_go))
        die("Raida Go doesn't have exec permissions");

$data['status'] = 0;
function verify_payment($guid, $amount, $skywallet,$meta) {
        global $raida_go;

        $cmd =  "$raida_go view_receipt $guid $skywallet";

        // Exec the binary
        $json = exec($cmd, $outarray, $error_code);
        if ($error_code != 0) {
            // echo "Invalid response from raida_go: $error_code, Output $json";
            // return 1;
			return  $data['status'] = 1;
        }

        $arr = json_decode($json, true);
        if (!$arr) {
			// echo "Failed to decode json: $json";
			return  $data['status'] = 1;
        }

        if (!isset($arr['amount_verified']) || !isset($arr['status'])) {
			// echo "Corrupted response: $json";
			return  $data['status'] = 1;
        }

        if ($arr['status'] != "success") {
			// echo "Invalid status in response: $json";
			return  $data['status'] = 1;
        }

        // Return Failed here if amount doesn't much. It means that transaction didn't happen
        $amountVerified = $arr['amount_verified'];
	
       /*  if ($amountVerified != $amount) {
                //echo "Invalid amount: $amountVerified, expected: $amount";
               return  $data['status'] = 2;
        } */
		$data['status'] = 1;
		$data['amount'] = $amountVerified;
		$data['mywallet'] = $skywallet;
		$data['guid'] = $guid;
		$mta = explode('=',base64_decode($meta));
		$mta3 = str_replace('"','',$mta[2]);
		$mta3 = str_replace('Meta: ','',$mta3);
		$mta3 = str_replace(' ','',$mta3);
		
		$data['meta'] = $mta3;
        return $data;
}

isset($_GET['merchant_skywallet'])? $mywallet = $_GET['merchant_skywallet']: $mywallet = '';
isset($_GET['amount'])? $amount = $_GET['amount']: $amount = 0;
isset($_GET['guid'])? $guid = $_GET['guid']: $guid = '';
isset($_GET['meta'])? $meta = $_GET['meta']: $meta = '';

$result = verify_payment($guid, $amount, $mywallet,$meta);

// echo '<pre>';
// print_r($result);
// die();


ob_start();
echo '<pre>';				
print_r($_GET);		
echo '</pre>';	
$out1 = ob_get_contents();
ob_end_clean();
$file = fopen("debug/".$guid.".txt", "w");
fwrite($file, $out1); 
fclose($file);   

if($result['status'] == 1){	
	$result['key'] = '3c8dea116b964d529faf34b21bbaf595';
	
	$url = 'https://yourcompany.com/lht/receivecc';
	$ch = curl_init();
	curl_setopt($ch, CURLOPT_URL, $url);	
	curl_setopt($ch, CURLOPT_FOLLOWLOCATION, 1);
	curl_setopt($ch,CURLOPT_USERAGENT,'Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.13) Gecko/20080311 Firefox/2.0.0.13');
	//curl_setopt($ch, CURLOPT_USERPWD, $wallet_username . ":" . $wallet_password);
	curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
	curl_setopt($ch, CURLOPT_HTTPHEADER, array('Content-Type: application/json')); 
	//curl_setopt($ch, CURLOPT_POST, count($request));
	curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($result));
	curl_setopt($ch, CURLOPT_ENCODING, 'gzip,deflate');
	curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
	$return = curl_exec($ch);		
	curl_close($ch);
	echo $return;
}

