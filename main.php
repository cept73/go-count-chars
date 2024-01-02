<?php

$start = microtime(true);

$result = [];

$fp = fopen('data.txt', 'rt');
while ($buffer = fgets($fp, 4096)) {
    for ($i = 0; $i < strlen($buffer); $i ++) {
        $c = $buffer[$i];
        if (!isset($result[$c])) {
            $result[$c] = 0;
        }
        $result[$c] ++;
    }
}
fclose($fp);

$duration = microtime(true) - $start;

print json_encode($result)
    . PHP_EOL
    . sprintf("%0.3fs", $duration)
    . PHP_EOL;
