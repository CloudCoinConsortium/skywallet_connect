<?php
        $mywallet = $_GET['merchant_skywallet'];
	$amount_due = $_GET['amount'];
	$guid = $_GET['guid'];
      	$customer_skywallet=$_GET['customer_skywallet'];


	$command = "/opt/raida_go view_receipt $guid $mywallet";
	echo "<br><b>The command is:</b> $command";
	$json_obj = exec($command, $outarray, $error_code);

	echo "<br><b>The response:</b> <code>$json_obj</code><br>";

	$arr = json_decode($json_obj, true);

	echo "<br><b>$customer_skywallet sent:</b>".intval($arr["amount_verified"]);



?>
