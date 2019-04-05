# Expgraphy
Expgraphy 解析一个数学表达式，输出复合数学规范的表达式图片

## 用法

[example](https://github.com/x6uderek/expgraphy/blob/master/example.go)

```go
var input = flag.String("exp", "", "math expression")
var output = flag.String("out", "./out.png", "output file")
var fontSize = flag.Float64("font", 36.0, "font size in float64")

func main() {
	flag.Parse()
	if *input == "" {
		log.Fatal("empty input")
	}
	expr, err := exp.Parse(*input) //解析表达式
	if err != nil {
		log.Fatalf("parse err: %v",err)
	}
	gra := expr.GetGraphy() //获取绘图工具
	size := gra.Measure(*fontSize) //测量表达式尺寸
	img := image.NewRGBA(image.Rect(0,0, int(exp.FixedToFloat(size.X)*2), int(exp.FixedToFloat(size.Y)*2))) //绘图板尺寸是图片的两倍
	startPoint := fixed.Point26_6{X:size.X/2, Y:size.Y*3/2}  //绘图起点在左下角
	gra.Draw(img, startPoint) //绘制

	f,err := os.Create(*output)
	if err!=nil {
		log.Fatalf("output file: %v", err)
	}
	defer f.Close()
	png.Encode(f, img) //保存文件
}
```

## 效果

输入·expgraphy -exp "(x+y)*(x-y)/(y+1)" -out out1.png -font 30·

输出
![geometry](https://github.com/x6uderek/expgraphy/blob/master/output/out1.png)

输入·expgraphy -exp "pow(x+y,x)" -out out2.png -font 30·

输出
![geometry](https://github.com/x6uderek/expgraphy/blob/master/output/out2.png)

输入·expgraphy -exp "sqrt(x*x+y*y)" -out out3.png -font 30·

输出
![geometry](https://github.com/x6uderek/expgraphy/blob/master/output/out3.png)