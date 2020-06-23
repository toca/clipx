# clipx
clipboard extension console base application

![2020-06-23](https://user-images.githubusercontent.com/3282299/85348462-a3408f00-b536-11ea-9ec3-a7eea5d451ca.png)

* cmd.exe から起動
* ctrl x 2 => ウインドウを表示
* 選択して貼り付け (マウスでも可能)

* ログは std err へ出力する

### option
* -s クリップボードの内容を保存するファイルの path を指定

### memo
* clipx.exe -s .\cb_history.dat 2> log.txt
