# README

## About

This is the official Wails Vue template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.


## TODO

1. Windows、Linux、Mac、标签、集合和指令六个按钮，每次点击其中一个都出发一次查询。
2. 点击查询的时候，带上所有选项。
3. 

系统只作为选项
1. 使用搜索框
    1.1 标签、集合和指令三选一，如选择标签时，点击搜索查询名称相关的所有标签，默认显示所有标签相关的指令。
2. 不使用搜索框时，默认显示所有标签、集合和指令。