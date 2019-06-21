<?php
/**
 * php数组转文本格式，key=》val
 * arr.php
 * @auther: lingyi
 * email: zhenpsm@163.com
 * Date: 19-5-7
 * Time: 下午1:42
 */

set_time_limit(0);

define("BASEPATH", "beluga");


// 测试代码
// $conf = [
//     "index" => "123",
//     "index2" => [
//         "index3" => [
//             "index4" => 1234,
//             "n1" => 21,
//             "n2" => [
//                 "ce1" => 1
//             ]
//         ],
//         "index5" => "eweww",
//         "ceshi" => 123,
//         "hoale" => '111'
//     ],
//     "nima" => [
//         "serq" => 123,
//         "serq1" => [
//             "dsa" => 11222,
//             "hata" => [
//                 2, 22, 33, 11111, 22331
//             ]
//         ]
//     ],
//     "2342" => 234,
//     "putong" => [
//         1, 2, 3, 4, 5, 6
//     ]
// ];

// arrToText($arr_arr);
// print_r($text);

$conf = [
    'sub_list' => array(
        '1' => array(
            'detail' => array(
                '1' => array(
                    'method_name' => 'invite_gold',
                ),
            ),
        ),
    )
];

$toole_obj = new Toole();
$toole_obj->arrToText($conf);

print_r($toole_obj->text);

class Toole
{
    public $text;

    public function __construct()
    {
        $this->text = "";
    }

    public function arrToText($arr)
    {
        if ($this->text) {
            dir("我曹数据不为空了");
        }
        foreach ($arr as $key => $val) {
            if (is_array($val)) {
                if (!$this->is_assoc($val)) {
                    $this->text .= $key . "=" . json_encode($val, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES) . "\n";
                } else {
                    $this->digui($val, $key);
                }
            } else {
                $this->text .= $key . "=" . $val . "\n";
            }
        }
    }

    public function digui($arr, $k = "")
    {
        $tmp = $k;

        foreach ($arr as $key => $val) {
            if (is_array($val)) {
                if (!$this->is_assoc($val)) {
                    $this->text .= $k . "." . $key . "=" . json_encode($val, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES) . "\n";;
                } else {
                    $k .= "." . $key;
                    $this->digui($val, $k);
                    $k = $tmp;
                }
            } else {
                $this->text .= $k . "." . $key . "=" . $val . "\n";
            }
        }
    }

    // 判断数组是否为普通数组
    public function is_assoc($arr)
    {
        return array_keys($arr) !== range(0, count($arr) - 1);
    }
}

// // 文件所在目录
// $read_file_path = $argv[1];

// // 文件写入目录
// $wirt_file_path = $argv[2];

// if (!is_dir($wirt_file_path)) {
//     mkdir($wirt_file_path);
// }
// if (!is_file($read_file_path)) {
//     die("配置文件不存在！\n");
// }

// // 获取文件名
// $file_path = explode("/", $read_file_path);
// $file_name = $wirt_file_path . $file_path[count($file_path) - 1];

// $toole_obj = new Toole();
// $configs = require($read_file_path);
// if($configs === 1) {
//     $configs = $config;
// }
// $toole_obj->arrToText($configs);

// file_put_contents($file_name, $toole_obj->text);
// unset($toole_obj);
// print_r("写入完成\n");