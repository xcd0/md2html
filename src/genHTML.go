package main

import ( // {{{

	"encoding/base64"

	//"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	//"github.com/google/go-github/github"
	//"github.com/russross/blackfriday"

	gfm "github.com/shurcooL/github_flavored_markdown"
	"github.com/xcd0/go-nkf"
) // }}}

type Fileinfo struct { // {{{
	Md       string
	Apath    string     // 入力mdファイルの絶対パス
	Dpath    string     // 入力mdファイルのあるディレクトリのパス
	Filename string     // 入力mdファイルのファイル名
	Basename string     // 入力mdファイルのベースネーム 拡張子抜きの名前
	Ext      string     // 入力mdファイルの拡張子
	Csspath  string     // 入力Cssファイルの出力先パス
	Htmlpath string     // 生成されるhtmlファイルの出力先パス
	Flavor   string     // 生成に用いるmarkdownの方言
	Html     string     // 生成したhtml本体が入る
	Pdfpath  string     // 生成されるpdfファイルの出力先パス
	Rep      [][]string // 置換対象文字列と置換文字列
	RImgPath []string   // 置換対象画像ファイルへの相対パス
} // }}}

func Makehtml(fi *Fileinfo) string { // {{{

	header := Makeheader(*fi, "")
	body := Makebody(fi.Filename, fi.RImgPath, "")
	footer := Makefooter()

	fi.Html = header + body + footer

	return fi.Html
} // }}}

func Argparse(arg string) Fileinfo { // {{{

	fi := Fileinfo{}

	// 絶対パスを得る
	fi.Apath, _ = filepath.Abs(arg)
	// ファイルパスをディレクトリパスとファイル名に分割する
	fi.Dpath, fi.Filename = filepath.Split(fi.Apath)
	// 拡張子を得る
	fi.Ext = filepath.Ext(fi.Filename)
	// 拡張子なしの名前を得る
	fi.Basename = fi.Filename[:len(fi.Filename)-len(fi.Ext)]
	// 出力するhtmlのパスを得る
	fi.Htmlpath = fi.Dpath + fi.Basename + ".html"
	// 入力Cssのパスを得る
	fi.Csspath = fi.Dpath + "markdown.css"

	return fi
} // }}}

func Makeheader(fi Fileinfo, csspath string) string { // {{{

	header_css := CreateMinifiedCss(csspath)

	// katex新しいの行内で動かなかった
	/*
		<!-- オンラインの時用 -->
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.12.0/dist/katex.min.css" integrity="sha384-AfEj0r4/OFrOo5t7NnNe46zW/tFgW6x/bCJG8FqQCEo3+Aro6EYUG4+cU+KJWu/X" crossorigin="anonymous">
		<script defer src="https://cdn.jsdelivr.net/npm/katex@0.12.0/dist/katex.min.js" integrity="sha384-g7c+Jr9ZivxKLnZTDUhnkOnsh30B4H0rpLUpJ4jAIKs4fnJI+sEnkvrMWph2EDg4" crossorigin="anonymous"></script>
		<script defer src="https://cdn.jsdelivr.net/npm/katex@0.12.0/dist/contrib/auto-render.min.js" integrity="sha384-mll67QQFJfxn0IYznZYonOWZ644AWYC+Pt2cHqMaRhXVrursRwvLnLaebdGIlYNa" crossorigin="anonymous" onload="renderMathInElement(document.body);"></script>

		<!-- オフラインの時用 katexフォルダがある前提 -->
		<link rel="stylesheet" href="katex/katex.min.css">
		<script src="katex/katex.min.js"></script>
		<script src="katex/auto-render.min.js" onload="renderMathInElement(document.body);"></script>
	*/

	header1 := `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<style type="text/css"><!--
`
	header2 := `
--></style>
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.7.1/katex.min.css">
	<script src="https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.7.1/katex.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.7.1/contrib/auto-render.min.js"></script>
</head>
<body>
`
	return header1 + header_css + header2
} // }}}

func CreateMinifiedCss(csspath string) string { // {{{

	_, err := os.Stat(csspath)
	if err != nil {
		// cssファイルがない // デフォルトのCSSを使う
		// fmt.Printf("File : %v is not found.", csspath)
		// fmt.Println("デフォルトのcssファイルを使用します。")
		// fmt.Println(err)

		css := `
/* build in css*/
body{font-family:Helvetica,arial,sans-serif;font-size:14px;line-height:1.8;padding:30px;background-color:#fff;color:#333}body>:first-child{margin-top:0!important}body>:last-child{margin-bottom:0!important}a{color:#4183c4;text-decoration:none}a.absent{color:#c00}a.anchor{display:block;padding-left:30px;margin-left:-30px;cursor:pointer;position:absolute;top:0;left:0;bottom:0x}h1,h2,h3,h4,h5,h6{margin:20px 0 10px;padding:0;font-weight:700;-webkit-font-smoothing:antialiased;cursor:text;position:relative}h1:first-child,h1:first-child+h2,h2:first-child,h3:first-child,h4:first-child,h5:first-child,h6:first-child{margin-top:0;padding-top:0}h1:hover a.anchor,h2:hover a.anchor,h3:hover a.anchor,h4:hover a.anchor,h5:hover a.anchor,h6:hover a.anchor{text-decoration:none}h1 code,h1 tt,h2 code,h2 tt,h3 code,h3 tt,h4 code,h4 tt,h5 code,h5 tt,h6 code,h6 tt{font-size:inherit}h1{font-size:34px;margin-bottom:40px;padding-bottom:0}h1,h2{color:#000}h2{font-size:30px;border-bottom:2px solid #ccc}h3{font-size:24px;border-bottom:1px solid #ddd}h4{font-size:20px}h5{font-size:18px}h6{font-size:16px;color:#777}blockquote,dl,li,ol,p,pre,table,ul{margin:10px 0}hr{border:0 0 0;height:4px;padding:0}a:first-child h1,a:first-child h2,a:first-child h3,a:first-child h4,a:first-child h5,a:first-child h6,bo dy>h5:first-child,body>h1:first-child,body>h1:first-child+h2,body>h2:first-child,body>h3:first-child,body>h4:first-child,body>h6:first-child{margin-top:0;padding-top:0}h1 p,h2 p,h3 p,h4 p,h5 p,h6 p{margin-top:0}li p.first{display:inline-block}ol,ul{padding-left:30px}ol:first-child,ul:first-child{margin-top:0}dl,dl dt{padding:0}dl dt{font-size:14px;font-weight:700;font-weight:1400;font-style:italic;margin:15px 0 5px}dl dt:first-child{padding:0}dl dt>:first-child{margin-top:0}dl dt>:last-child{margin-bottom:0}dl dd{margin:0 0 15px;padding:0 15px}dl dd>:first-child{margin-top:0}dl dd>:last-child{margin-bottom:0}blockquote{border-left:4px solid #ddd;padding:0 15px;color:#777}blockquote>:first-child{margin-top:0}blockquote>:last-child{margin-bottom:0}table{padding:0;border-spacing:2px;border-collapse:collapse;width:80%;margin:auto}table,td,th{border:1px solid #ccc}td,th{padding:0;margin:0}table tr{background-color:#fff;border-top:1px solid #c6cbd1;margin:0;padding:0}table tr:nth-child(2n){background-color:#f6f8fa}table tr th{font-weight:700}table tr td,table tr th{border:1px solid #ccc;text-align:center;margin:0;padding:6px 13px}table tr td:first-child,table tr th:first-child{margin-top:0}img{max-width:100%}code,tt{margin:0 2px;padding:0 5px;white-space:nowrap;border:1px solid #eaeaea;background-color:#f8f8f8;border-radius:3px}pre code{margin:0;padding:0;white-space:pre;border:0;background:0 0}.highlight pre,pre{border:1px solid #ccc;font-size:13px;line-height:19px;overflow:auto;padding:6px 10px;border-radius:3px}pre code,pre tt{background-color:transparent;border:0}.main-content{max-width:50pc;margin:auto;padding-bottom:50px}hr{border:0!important;color:#fff;height:4px}.page_num{border:0;position:absolute;right:10;bottom:10}
`
		return css
	}

	minifiedCss := Minify(csspath)
	//fmt.Println(minifiedCss)

	//fmt.Println(outputCssPath)

	return string(minifiedCss)
} // }}}

func ReadMd(path string) string { // {{{

	bytemd, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ファイル%vが読み込めません\n", path)
		log.Println(err)
		panic(err)
		return ""
	}

	// ファイルの文字コード変換
	charset, err := nkf.CharDet(bytemd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "文字コード変換に失敗しました\nutf8を使用してください\n")
		log.Println(err)
		panic(err)
		return ""
	}

	stringmd, err := nkf.ToUtf8(string(bytemd), charset)

	stringmd = convNewline(stringmd, "\n")

	return stringmd
} // }}}

func convNewline(str, nlcode string) string { // {{{
	return strings.NewReplacer(
		"\r\n", nlcode,
		"\r", nlcode,
		"\n", nlcode,
	).Replace(str)
} // }}}

func readMd(fi *Fileinfo) { // {{{
	if fi.Ext != ".md" {
		fmt.Println("拡張子が.mdではありません")
		fmt.Fprintf(
			os.Stderr,
			"%s は拡張子が.mdではありません。\n"+
				"拡張子が.mdのマークダウンのファイルを指定してください。\n",
			fi.Filename)
		panic("終了します")
	}

	fi.Md = ReadMd(fi.Apath)
} // }}}

func ReplaceImg(outputList []string, html string) string { // {{{
	output := "" // これが出力される

	// ここで一行ずつ処理
	lines := strings.Split(html, "\n")

	for _, line := range lines { // 一行ずつ
		//for j, line := range lines { // 一行ずつ
		// <img src=が含まれるはずなので前もって判別しておく
		if strings.Contains(line, "<img src=") == false {
			// <img src=が含まれていない
			output += line + "\n"
		} else {
			// <img src=が含まれる

			// リストにある画像のパスが含まれるか調べる
			for _, path := range outputList {
				// もしoutputListにおいて先にマッチしていたら
				// パスが\で区切られている場合を考えてすべて/にする
				line = strings.Replace(line, "\\", "/", -1)
				path = strings.Replace(path, "\\", "/", -1)
				// チェック
				if strings.Contains(line, path) == false {
					// 違うので次のpathに
					continue
				}
				// マッチ

				//fmt.Printf("---- line : %v , match : %v ----\n%v\n", j, path, line)

				// lineに含まれるパスの前後を切り出す
				// 例
				// <li><p><img src="./img/build_on_win10.gif" alt=""></p></li>
				//                 ↑                       ↑
				// このダブルクォーテーションの位置を調べる
				// ./が混じらないよう 一文字づつ前方に調べて"を探す

				// <img>にリンクが張ってあればそれを消す
				tmpLinkNum := strings.Index(line, "<a")
				if tmpLinkNum >= 0 {
					tmpLinkNumPost := -1
					for i := tmpLinkNum; i < len(line); i++ {
						if line[i] == '>' {
							// 見つかったのでそれをpostDQとする
							tmpLinkNumPost = i + 1
							break
						}
					}
					if tmpLinkNumPost < 0 {
						panic(1)
					} else {

						//fmt.Println("--->>\n" + line)
						line = line[:tmpLinkNum] + line[tmpLinkNumPost:]
						//fmt.Println("---<<\n" + line)
					}
					tmpLinkNum = strings.Index(line, "</a>")

					if tmpLinkNum >= 0 {
						//fmt.Println("--->>\n" + line)
						line = line[:tmpLinkNum] + line[tmpLinkNum+len("</a>"):]
						//fmt.Println("---<<\n" + line)
					}
				}

				// <imgより前にマッチしないようにする
				preDQ := -1
				tmpNum := strings.Index(line, "<img")
				tmpLine := line[tmpNum:]
				//fmt.Println("---tmpline\n" + tmpLine)

				for i := strings.Index(tmpLine, path); i >= 0; i-- {
					if tmpLine[i] == '"' {
						// 見つかったのでそれをpreDQとする
						preDQ = tmpNum + i + 1
						//fmt.Println("---match!\n" + line[:preDQ])
						break
					}
				}
				if preDQ < 0 {
					// ダブルクォーテーションが見つからなかった
					fmt.Println("マッチしましたが前にダブルクォーテーションがありませんでした")
					output += line + "\n"
					break
				}

				postDQ := -1
				for i := preDQ; i < len(line); i++ {
					if line[i] == '"' {
						// 見つかったのでそれをpostDQとする
						postDQ = i
						break
					}
				}
				if postDQ < 0 {
					// ダブルクォーテーションが見つからなかった
					fmt.Println("マッチしましたが後ろにダブルクォーテーションがありませんでした")
					/*
						fmt.Println("ori " + line)
						fmt.Println("cut pre  :" + line[preDQ:])
						fmt.Println("cut post :" + line[:])
					*/
					output += line + "\n"
					break
				}

				//
				/*
					fmt.Println("<<\n" + line)
					fmt.Printf(">> pre : %v '%v', post : %v '%v'\n", preDQ, line[preDQ], postDQ, line[postDQ])
					fmt.Printf(">> %v\n", line[preDQ:postDQ])

					code := "@@@@@@"
					/ */
				///*
				// 文字列が含まれているかどうか調べる
				// base64でエンコードしたデータに置き換える
				base64code := EncodeBase64(path)
				// ファイルの拡張子ごとにヘッダをつける
				// gif,png,jpg,jpegのみ
				ext := filepath.Ext(path) // 一致した画像ファイルのパスから拡張子を調べる
				code := ""
				if ext == ".gif" {
					code = "data:image/gif;base64," + base64code
				} else if ext == ".png" {
					code = "data:image/png;base64," + base64code
				} else if ext == ".jpg" || ext == ".jpeg" {
					code = "data:image/jpeg;base64," + base64code
				}

				// */
				// 前後を切り出す
				pre := line[:preDQ]
				post := line[postDQ:]

				// くっつけて上書きする
				line = pre + code + post

				//fmt.Println(">>\n" + line)

				// 置換したらその行は終わる
				break

			}
			// すべてのpathをチェックした
			//fmt.Println(line)
			// 出力する
			output += line + "\n"
		}
	}

	return output

} // }}}

func ReplaceImg4mdPre(outputList []string, stringmd string) string { // {{{1
	output := ""

	// ここで一行ずつ処理
	lines := strings.Split(stringmd, "\n")

	for _, line := range lines { // 一行ずつ
		if strings.Contains(line, "![") == false {
			// ![が含まれていない
			output += line + "\n"
		} else {
			// lineに含まれるパスの前後を切り出す
			// 例
			// ![ hoge hoge ](./img/fuga) ふがふが ![ hoge hoge ](./img/fuga) にゃあ ![ hoge hoge ](./img/fuga) ほげ
			//               ↑         ↑
			// この()の位置を調べる

			// 一行に複数の画像がある場合
			// ![]()にある程度マッチする正規表現
			// バックスラッシュを蹴るのはどうやったらいいのかわからない...
			r := regexp.MustCompile(`!\[[^!\[\]\n\r\\\(\)\*{}&$#@]*\]([^!\n\r\[\]\*{}&%$#@]*)`)
			ret := r.FindAllStringSubmatch(line, -1)
			// これでこの行にある！の数がわかる len(ret)

			// ![]()が複数個あるとして、
			// それはこの時点でbase64に変換する

			// それ以外の部分
			// 空白とか文字とかを切り出す
			tmp := line
			lineout := ""
			//fmt.Printf(">> ret : %v\n", ret)
			for _, g := range ret {
				tmpPath := fmt.Sprintf("%v", g[1][1:len(g[1])-1])
				base64code := EncodeBase64(tmpPath)
				// ファイルの拡張子ごとにヘッダをつける
				// gif,png,jpg,jpegのみ
				ext := filepath.Ext(tmpPath) // 一致した画像ファイルのパスから拡張子を調べる
				code := ""
				if ext == ".gif" {
					code = "data:image/gif;base64," + base64code
				} else if ext == ".png" {
					code = "data:image/png;base64," + base64code
				} else if ext == ".jpg" || ext == ".jpeg" {
					code = "data:image/jpeg;base64," + base64code
				}

				index := strings.Index(tmp, "!")
				if index < 0 {
					// !がない
					// 残りの部分をくっつける
					lineout += tmp
					break
				} else if index > 0 {
					// 複数あったときの行頭以外
					// 前後を切り取る
					pre := tmp[:index]
					post := tmp[index+len(g[0]):]
					lineout += pre
					lineout += g[0][:len(g[0])-len(tmpPath+")")] // ![aaa ]( ←ここまで
					lineout += code
					lineout += ")"
					if len(post) > 0 {
						tmp = post
					} else {
						break
					}
				} else {
					// 行頭
					lineout += g[0][:len(g[0])-len(tmpPath+")")] // ![aaa ]( ←ここまで
					lineout += code                              // エンコードしたもの
					lineout += ")"
					// これで1つ目の切り出しは終了
					// )以降の文字列をさらに処理していく
					tmp = tmp[len(g[0]):]
					// これでtmpに![]()以降が入る
					if len(tmp) <= 0 {
						// 後ろがないので終了
						break
					}

				}
			}
			output += lineout + "\n"
			/* // {{{

			for _, path := range outputList {
			// リストにある画像のパスが含まれるか調べる
					マッチ不要な気がしている
				// もしoutputListにおいて先にマッチしていたら
				// パスが\で区切られている場合を考えてすべて/にする
				line = strings.Replace(line, "\\", "/", -1)
				path = strings.Replace(path, "\\", "/", -1)
				// チェック
				if strings.Contains(line, path) == false {
					// 違うので次のpathに
					continue
				}
				// マッチ

				// * /
				//fmt.Printf("---- line : %v , match : %v ----\n%v\n", j, path, line)

				// ]より前にマッチしないようにする
				preDQ := -1
				tmpNum := strings.Index(line, "]")
				tmpLine := line[tmpNum:]
				//fmt.Println("---tmpline\n" + tmpLine)

				for i := strings.Index(tmpLine, path); i >= 0; i-- {
					if tmpLine[i] == '"' {
						// 見つかったのでそれをpreDQとする
						preDQ = tmpNum + i + 1
						//fmt.Println("---match!\n" + line[:preDQ])
						break
					}
				}
				if preDQ < 0 {
					// ダブルクォーテーションが見つからなかった
					fmt.Println("マッチしましたが前にダブルクォーテーションがありませんでした")
					output += line + "\n"
					break
				}

				postDQ := -1
				for i := preDQ; i < len(line); i++ {
					if line[i] == '"' {
						// 見つかったのでそれをpostDQとする
						postDQ = i
						break
					}
				}
				if postDQ < 0 {
					// ダブルクォーテーションが見つからなかった
					fmt.Println("マッチしましたが後ろにダブルクォーテーションがありませんでした")
					/*
						fmt.Println("ori " + line)
						fmt.Println("cut pre  :" + line[preDQ:])
						fmt.Println("cut post :" + line[:])
					// * /
					output += line + "\n"
					break
				}

				// / *
					fmt.Println("<<\n" + line)
					fmt.Printf(">> pre : %v '%v', post : %v '%v'\n", preDQ, line[preDQ], postDQ, line[postDQ])
					fmt.Printf(">> %v\n", line[preDQ:postDQ])

					code := "@@@@@@"
					// * /
				/// *
				// 文字列が含まれているかどうか調べる
				// base64でエンコードしたデータに置き換える
				base64code := EncodeBase64(path)
				// ファイルの拡張子ごとにヘッダをつける
				// gif,png,jpg,jpegのみ
				ext := filepath.Ext(path) // 一致した画像ファイルのパスから拡張子を調べる
				code := ""
				if ext == ".gif" {
					code = "data:image/gif;base64," + base64code
				} else if ext == ".png" {
					code = "data:image/png;base64," + base64code
				} else if ext == ".jpg" || ext == ".jpeg" {
					code = "data:image/jpeg;base64," + base64code
				}

				// * /
				// 前後を切り出す
				pre := line[:preDQ]
				post := line[postDQ:]

				// くっつけて上書きする
				line = pre + code + post

				//fmt.Println(">>\n" + line)

				// 置換したらその行は終わる
				break
			}
			// すべてのpathをチェックした
			//fmt.Println(line)
			// 出力する
			*/ // }}}
		}
	}

	return output

} // }}}1

/*
func gitHubAPI(md []byte) ([]byte, error) { // {{{
	client := github.NewClient(nil)
	opt := &github.MarkdownOptions{Mode: "gfm", Context: "google/go-github"}
	body, _, err := client.Markdown(context.Background(), string(md), opt)
	return []byte(body), err
} // }}}
*/

func shurcooL_GFM(md []byte) ([]byte, error) { // {{{

	bytehtml := gfm.Markdown(md)

	return bytehtml, nil
}

// }}}

func EncodeBase64(str string) string { // {{{

	file, err := os.Open(str)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	fi, _ := file.Stat() //FileInfo interface
	size := fi.Size()    //ファイルサイズ

	data := make([]byte, size)
	file.Read(data)

	return base64.StdEncoding.EncodeToString(data)
} // }}}
