# Gs
 win中没有grep，那么我自己造一个，原本是用gf工具来提取关键词的，但是window上没有grep，就手搓了一个
 In Windows, there's no grep, so I decided to make one myself. Originally, I was using the gf tool to extract keywords, but since Windows lacks grep, I ended up crafting my own.
 
## 安装-install
```
go install github.com/soryecker/Gs@latest
```

## 运行-run

`Gs -f <input_file> -o <output_file> -r <regular_expression_file>`

```
Gs:Window上的grep工具 @PuffDog

  -f string
        Path to the file containing URLs to scan.
        要扫描的url的文件路径
  -o string
        Path to the output file for results
        结果输出文件的路径
  -r string
        Path to the file containing regular expressions
        正则表达式的文件的路径
```
