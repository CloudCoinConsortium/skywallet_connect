<?php
        $mywallet = $_GET['merchant_skywallet'];
	$customer_skywallet=$_GET['customer_skywallet'];
	$amount_due = $_GET['amount'];
	$memo = $_GET['guid'];

	$command = "/opt/raida_go";
	if (!file_exists($command))
  		die("Raida_go doesn't exist");

	if (!is_executable($command))
  		die("Raida_go is not executable for webserver user");

	$command = "$command view_receipt $memo $mywallet";
	echo "<br><b>The command is:</b> $command";

	$json_obj = exec($command, $outarray, $error_code);
	echo "<br><b>The response from raida_go:</b> <code>$json_obj</code><br>";

	$arr = json_decode($json_obj, true);
	echo "<br><b>raida_go verified that $customer_skywallet sent $merchant_skywallet this amount:</b>".intval($arr["amount_verified"]);
?>
