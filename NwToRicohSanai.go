package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"golang.org/x/text/width"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func logWrite(logstr string, err error) {
	if err != nil {
		log.Printf("%s: %s", logstr, err)
	}
}

func main() {

	//ログファイル準備
	logfile, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	failOnError(err)
	defer logfile.Close()

	log.SetOutput(logfile)

	// 入力ファイル準備
	flag.Parse()
	infile, err := os.Open(flag.Arg(0))
	failOnError(err)
	defer infile.Close()

	// 書き込みファイル準備
	outfile, err := os.Create("./リコー三愛グループ健康保険組合健診データ" + time.Now().Format("20060102") + ".csv")
	failOnError(err)
	defer outfile.Close()

	// reader writerの準備
	reader := csv.NewReader(transform.NewReader(infile, japanese.ShiftJIS.NewDecoder()))
	reader.Comma = '\t'
	writer := csv.NewWriter(transform.NewWriter(outfile, japanese.ShiftJIS.NewEncoder()))
	writer.Comma = ','
	writer.UseCRLF = true

	// メイン処理をスタート
	log.Print("Start\r\n")

	// タイトル行をよみだす
	_, err = reader.Read()
	failOnError(err)

	// タイトル行を書きだす
	writer.Write(titleWrite())

	for {
		items, err := reader.Read() // １行読みだす
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}

		logstr := items[20] + " " + items[7] // ログ用　受診番号 氏名
		str := ""
		strCd := ""
		strName := ""

		var writeItems []string

		// CSVフォーマットVer
		writeItems = append(writeItems, "RB_Ver.1.0")

		// 提出先
		writeItems = append(writeItems, "BIO(RICOH)")

		// データ作成者
		writeItems = append(writeItems, "医療法人社団　松英会")

		// データ作成日
		writeItems = append(writeItems, time.Now().Format("2006/01/02"))

		// データ提出日
		writeItems = append(writeItems, time.Now().Format("2006/01/02"))

		// データ登録完了区分
		writeItems = append(writeItems, "1")

		// 登録未完了の連絡内容
		writeItems = append(writeItems, "")

		// 団体コード
		writeItems = append(writeItems, "RICOH")

		// 団体コード名称
		logWrite(logstr, requireChk(items[3], "所属名1"))
		writeItems = append(writeItems, items[3]) // ←所属名1

		// 事業所コード
		if items[0] != "04019001" { // （株）リコーは所属２をチェックしない
			logWrite(logstr, requireChk(items[4], "所属cd2"))
		}
		writeItems = append(writeItems, items[4]) // ←所属cd2

		// 事業所名称
		if items[0] != "04019001" { // （株）リコーは所属２をチェックしない
			logWrite(logstr, requireChk(items[5], "所属名2"))
		}
		writeItems = append(writeItems, items[5]) // ←所属名2

		// 個人ID
		if items[0] != "04019001" { // （株）リコーは個人IDをチェックしない
			str, err = kojinIdChk(items[6])
			logWrite(logstr, err)
			writeItems = append(writeItems, str) // ←社員No
		} else {
			writeItems = append(writeItems, items[6])
		}

		// 漢字氏名
		logWrite(logstr, requireChk(items[7], "漢字氏名"))
		writeItems = append(writeItems, items[7])

		// カナ氏名
		logWrite(logstr, requireChk(items[8], "カナ氏名"))
		writeItems = append(writeItems, items[8])

		// 生年月日
		str, err = waToSeireki(items[9])
		logWrite(logstr, err)
		logWrite(logstr, requireChk(items[9], "生年月日"))
		writeItems = append(writeItems, str)

		// 性別
		str, err = seiConv(items[10])
		logWrite(logstr, err)
		logWrite(logstr, requireChk(items[10], "性別"))
		writeItems = append(writeItems, str)

		// 保険者番号
		if items[0] != "04019001" { // （株）リコーは保健者番号をチェックしない
			logWrite(logstr, requireChk(items[12], "保険者番号"))
		}
		writeItems = append(writeItems, items[12])

		// 保険証記号
		if items[0] != "04019001" { // （株）リコーは保険証記号をチェックしない
			logWrite(logstr, requireChk(items[13], "保険証記号"))
		}
		writeItems = append(writeItems, items[13])

		// 保険証番号
		if items[0] != "04019001" { // （株）リコーは保険証番号をチェックしない
			logWrite(logstr, requireChk(items[14], "保険証番号"))
		}
		writeItems = append(writeItems, items[14])

		// 続柄
		writeItems = append(writeItems, "")

		// 予備
		writeItems = append(writeItems, "") // 予備
		writeItems = append(writeItems, "") // 予備

		// 受診券整理番号
		writeItems = append(writeItems, items[15])

		// 受診券有効期限
		writeItems = append(writeItems, items[16])

		// コースコード
		// コース名称
		strCd, strName, err = coursedConv(items[17], items[18], items[11])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// 受診日
		str, err = jdayConv(items[19])
		logWrite(logstr, err)
		logWrite(logstr, requireChk(items[19], "受診日"))
		writeItems = append(writeItems, str)

		// 施設/巡回区分
		str, err = sisetsuConv(items[21])
		logWrite(logstr, err)
		logWrite(logstr, requireChk(items[21], "施設/巡回区分"))
		writeItems = append(writeItems, str)

		// 健診機関コード
		writeItems = append(writeItems, "")

		// 健診機関名称
		writeItems = append(writeItems, "医療法人社団　松英会　馬込中央診療所")

		// [Met]特定健診機関番号
		writeItems = append(writeItems, "1311131242")

		// [Met]健診実施医師名
		writeItems = append(writeItems, "寺門　節雄")

		// 予備
		writeItems = append(writeItems, "") // 予備
		writeItems = append(writeItems, "") // 予備

		// 産業医判定区分
		writeItems = append(writeItems, "")

		// 就労区分
		writeItems = append(writeItems, "")

		// 産業医コメント
		writeItems = append(writeItems, "")

		// 伝達事項有無
		writeItems = append(writeItems, "")

		// 伝達内容
		writeItems = append(writeItems, "")

		// 診察判定区分コード
		// 診察判定区分名称
		strCd, strName, err = hanteiCdConv(items[479])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 診察所見
		writeItems = append(writeItems, limitStr(joinStr3(items[22], items[23], items[24]), 100))

		// 自覚症状など
		writeItems = append(writeItems, limitStr(joinStr5(items[25], items[26], items[27], items[28], items[29]), 100))

		// 既往歴の処理
		kiou := []string{items[30], items[33], items[36], items[39], items[42], items[45], items[48], items[51], items[54], items[57]}
		tenki := []string{items[32], items[35], items[38], items[41], items[44], items[47], items[50], items[53], items[56], items[59]}
		chiryoFlag, kiouFlag := tenkiConv(kiou, tenki)
		chiryoName, kiouName := kiouConv(kiou, tenki)

		// 治療中疾病有無区分
		writeItems = append(writeItems, chiryoFlag)

		// 治療中疾病名（文字）
		writeItems = append(writeItems, limitStr(chiryoName, 100))

		// 既往疾病有無区分
		writeItems = append(writeItems, kiouFlag)

		// 既往疾病名
		writeItems = append(writeItems, limitStr(kiouName, 100))

		// 総合判定区分コード
		// 総合判定区分名称
		strCd, strName, err = hanteiCdConv(items[314])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// 総合判定コメント
		var sogo [][2]string
		sogo = append(sogo, [2]string{items[314], items[316]}) // 総合判定
		sogo = append(sogo, [2]string{items[317], items[319]}) // BMI
		sogo = append(sogo, [2]string{items[320], items[322]}) // 体脂肪測定
		sogo = append(sogo, [2]string{items[323], items[325]}) // 聴力
		sogo = append(sogo, [2]string{items[326], items[328]}) // 視力
		sogo = append(sogo, [2]string{items[329], items[331]}) // 肺機能
		sogo = append(sogo, [2]string{items[332], items[334]}) // 肺年齢判定
		sogo = append(sogo, [2]string{items[335], items[337]}) // 血圧
		sogo = append(sogo, [2]string{items[338], items[340]}) // 尿糖
		sogo = append(sogo, [2]string{items[341], items[343]}) // 蛋白
		sogo = append(sogo, [2]string{items[344], items[346]}) // ウロビリ
		sogo = append(sogo, [2]string{items[347], items[349]}) // 潜血
		sogo = append(sogo, [2]string{items[350], items[352]}) // 尿比重
		sogo = append(sogo, [2]string{items[353], items[355]}) // 尿PH
		sogo = append(sogo, [2]string{items[356], items[358]}) // 尿沈渣まとめ
		sogo = append(sogo, [2]string{items[359], items[361]}) // 胸部X線
		sogo = append(sogo, [2]string{items[362], items[364]}) // 喀痰
		sogo = append(sogo, [2]string{items[365], items[367]}) // 心電図
		sogo = append(sogo, [2]string{items[368], items[370]}) // 貧血
		sogo = append(sogo, [2]string{items[371], items[373]}) // 血小板
		sogo = append(sogo, [2]string{items[374], items[376]}) // 白血球
		sogo = append(sogo, [2]string{items[377], items[379]}) // 白血球像
		sogo = append(sogo, [2]string{items[380], items[382]}) // 肝機能
		sogo = append(sogo, [2]string{items[383], items[385]}) // 膵機能
		sogo = append(sogo, [2]string{items[386], items[388]}) // 血中脂質
		sogo = append(sogo, [2]string{items[389], items[391]}) // 腎機能
		sogo = append(sogo, [2]string{items[392], items[394]}) // 腎機能コメント
		sogo = append(sogo, [2]string{items[395], items[397]}) // 血清尿酸
		sogo = append(sogo, [2]string{items[398], items[400]}) // 糖代謝
		sogo = append(sogo, [2]string{items[401], items[403]}) // 電解質
		sogo = append(sogo, [2]string{items[404], items[406]}) // 眼底
		sogo = append(sogo, [2]string{items[407], items[409]}) // 眼圧
		sogo = append(sogo, [2]string{items[410], items[412]}) // 胃部X線
		sogo = append(sogo, [2]string{items[413], items[415]}) // 胃内視鏡
		sogo = append(sogo, [2]string{items[416], items[418]}) // 胃内視生検
		sogo = append(sogo, [2]string{items[419], items[421]}) // 腹部エコー
		sogo = append(sogo, [2]string{items[422], items[424]}) // 便
		sogo = append(sogo, [2]string{items[425], items[427]}) // 便虫卵
		sogo = append(sogo, [2]string{items[428], items[430]}) // CRP
		sogo = append(sogo, [2]string{items[431], items[433]}) // リウマチ
		sogo = append(sogo, [2]string{items[434], items[436]}) // ピロリ菌
		sogo = append(sogo, [2]string{items[437], items[439]}) // PG検査
		sogo = append(sogo, [2]string{items[440], items[442]}) // 腫瘍マーカー
		sogo = append(sogo, [2]string{items[443], items[445]}) // 甲状腺
		sogo = append(sogo, [2]string{items[446], items[448]}) // 梅毒
		sogo = append(sogo, [2]string{items[449], items[451]}) // BNP
		sogo = append(sogo, [2]string{items[452], items[454]}) // 乳腺超音波
		sogo = append(sogo, [2]string{items[455], items[457]}) // マンモグラフィー
		sogo = append(sogo, [2]string{items[458], items[460]}) // 婦人内診察
		sogo = append(sogo, [2]string{items[461], items[463]}) // 子宮細胞診
		sogo = append(sogo, [2]string{items[464], items[466]}) // 骨密度
		sogo = append(sogo, [2]string{items[467], items[469]}) // 心エコー
		sogo = append(sogo, [2]string{items[470], items[472]}) // 血圧脈波
		sogo = append(sogo, [2]string{items[473], items[475]}) // 頸動脈エコー
		sogo = append(sogo, [2]string{items[476], items[478]}) // 甲状腺エコー
		sogo = append(sogo, [2]string{items[479], items[481]}) // 内科診察
		sogo = append(sogo, [2]string{items[482], items[484]}) // 腹部CT
		sogo = append(sogo, [2]string{items[485], items[487]}) // 治療中

		str, err = sogoConv(sogo)
		logWrite(logstr, err)
		writeItems = append(writeItems, limitStr(str, 1200))

		// 予備
		writeItems = append(writeItems, "") // 予備
		writeItems = append(writeItems, "") // 予備
		writeItems = append(writeItems, "") // 予備①(1)
		writeItems = append(writeItems, "") // 予備②(1)
		writeItems = append(writeItems, "") // 予備③(1)
		writeItems = append(writeItems, "") // 予備①(2)
		writeItems = append(writeItems, "") // 予備②(2)
		writeItems = append(writeItems, "") // 予備③(2)
		writeItems = append(writeItems, "") // 予備①(3)
		writeItems = append(writeItems, "") // 予備②(3)
		writeItems = append(writeItems, "") // 予備③(3)
		writeItems = append(writeItems, "") // 予備①(4)
		writeItems = append(writeItems, "") // 予備②(4)
		writeItems = append(writeItems, "") // 予備③(4)
		writeItems = append(writeItems, "") // 予備①(5)
		writeItems = append(writeItems, "") // 予備②(5)
		writeItems = append(writeItems, "") // 予備③(5)
		writeItems = append(writeItems, "") // 予備①(6)
		writeItems = append(writeItems, "") // 予備②(6)
		writeItems = append(writeItems, "") // 予備③(6)
		writeItems = append(writeItems, "") // 予備①(7)
		writeItems = append(writeItems, "") // 予備②(7)
		writeItems = append(writeItems, "") // 予備③(7)
		writeItems = append(writeItems, "") // 予備①(8)
		writeItems = append(writeItems, "") // 予備②(8)
		writeItems = append(writeItems, "") // 予備③(8)
		writeItems = append(writeItems, "") // 予備①(9)
		writeItems = append(writeItems, "") // 予備②(9)
		writeItems = append(writeItems, "") // 予備③(9)
		writeItems = append(writeItems, "") // 予備①(10)
		writeItems = append(writeItems, "") // 予備②(10)
		writeItems = append(writeItems, "") // 予備③(10)
		writeItems = append(writeItems, "") // 予備①(11)
		writeItems = append(writeItems, "") // 予備②(11)
		writeItems = append(writeItems, "") // 予備③(11)
		writeItems = append(writeItems, "") // 予備①(12)
		writeItems = append(writeItems, "") // 予備②(12)
		writeItems = append(writeItems, "") // 予備③(12)
		writeItems = append(writeItems, "") // 予備①(13)
		writeItems = append(writeItems, "") // 予備②(13)
		writeItems = append(writeItems, "") // 予備③(13)
		writeItems = append(writeItems, "") // 予備①(14)
		writeItems = append(writeItems, "") // 予備②(14)
		writeItems = append(writeItems, "") // 予備③(14)
		writeItems = append(writeItems, "") // 予備①(15)
		writeItems = append(writeItems, "") // 予備②(15)
		writeItems = append(writeItems, "") // 予備③(15)
		writeItems = append(writeItems, "") // 予備①(16)
		writeItems = append(writeItems, "") // 予備②(16)
		writeItems = append(writeItems, "") // 予備③(16)
		writeItems = append(writeItems, "") // 予備①(17)
		writeItems = append(writeItems, "") // 予備②(17)
		writeItems = append(writeItems, "") // 予備③(17)
		writeItems = append(writeItems, "") // 予備①(18)
		writeItems = append(writeItems, "") // 予備②(18)
		writeItems = append(writeItems, "") // 予備③(18)
		writeItems = append(writeItems, "") // 予備①(19)
		writeItems = append(writeItems, "") // 予備②(19)
		writeItems = append(writeItems, "") // 予備③(19)
		writeItems = append(writeItems, "") // 予備①(20)
		writeItems = append(writeItems, "") // 予備②(20)
		writeItems = append(writeItems, "") // 予備③(20)
		writeItems = append(writeItems, "") // 予備①(21)
		writeItems = append(writeItems, "") // 予備②(21)
		writeItems = append(writeItems, "") // 予備③(21)
		writeItems = append(writeItems, "") // 予備①(22)
		writeItems = append(writeItems, "") // 予備②(22)
		writeItems = append(writeItems, "") // 予備③(22)
		writeItems = append(writeItems, "") // 予備①(23)
		writeItems = append(writeItems, "") // 予備②(23)
		writeItems = append(writeItems, "") // 予備③(23)
		writeItems = append(writeItems, "") // 予備①(24)
		writeItems = append(writeItems, "") // 予備②(24)
		writeItems = append(writeItems, "") // 予備③(24)
		writeItems = append(writeItems, "") // 予備
		writeItems = append(writeItems, "") // 予備

		// その他判定区分コード
		writeItems = append(writeItems, "")

		// その他判定区分名称
		writeItems = append(writeItems, "")

		// その他データ内容
		writeItems = append(writeItems, "")

		// カンマ位置(131)
		writeItems = append(writeItems, "131")

		// 身長
		writeItems = append(writeItems, items[60])

		// 体重
		writeItems = append(writeItems, items[61])

		// BMI
		writeItems = append(writeItems, items[62])

		// 腹囲
		writeItems = append(writeItems, items[63])

		// 体脂肪率
		writeItems = append(writeItems, items[64])

		// 内臓脂肪面積
		writeItems = append(writeItems, "")

		// 5m視力裸眼右
		// 　データ属性
		strName, strCd = eyeConv(items[65])
		writeItems = append(writeItems, strName)
		writeItems = append(writeItems, strCd)

		// 5m視力裸眼左
		// 　データ属性
		strName, strCd = eyeConv(items[66])
		writeItems = append(writeItems, strName)
		writeItems = append(writeItems, strCd)

		// 5m視力矯正右
		// 　データ属性
		strName, strCd = eyeConv(items[67])
		writeItems = append(writeItems, strName)
		writeItems = append(writeItems, strCd)

		// 5m視力矯正左
		// 　データ属性
		strName, strCd = eyeConv(items[68])
		writeItems = append(writeItems, strName)
		writeItems = append(writeItems, strCd)

		// 近点視力裸眼右
		// 　データ属性
		strName, strCd = eyeConv(items[69])
		writeItems = append(writeItems, strName)
		writeItems = append(writeItems, strCd)

		// 近点視力裸眼左
		// 　データ属性
		strName, strCd = eyeConv(items[70])
		writeItems = append(writeItems, strName)
		writeItems = append(writeItems, strCd)

		// 近点視力矯正右
		// 　データ属性
		strName, strCd = eyeConv(items[71])
		writeItems = append(writeItems, strName)
		writeItems = append(writeItems, strCd)

		// 近点視力矯正左
		// 　データ属性
		strName, strCd = eyeConv(items[72])
		writeItems = append(writeItems, strName)
		writeItems = append(writeItems, strCd)

		// 視力矯正区分
		writeItems = append(writeItems, eyeKubun(items[67], items[68], items[71], items[72]))

		// 聴力右1K所見区分
		str, err = ear1kHantei(items[73])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 聴力右1K(dB)
		writeItems = append(writeItems, earConv(items[79]))

		// 聴力左1K所見区分
		str, err = ear1kHantei(items[74])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 聴力左1K(dB)
		writeItems = append(writeItems, earConv(items[80]))

		// 聴力右4K所見区分
		str, err = ear4kHantei(items[75], items[76])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 聴力右4K(dB)
		writeItems = append(writeItems, earConv(items[81]+items[82]))

		// 聴力左4K所見区分
		str, err = ear4kHantei(items[77], items[78])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 聴力左4K(dB)
		writeItems = append(writeItems, earConv(items[83]+items[84]))

		// 聴力会話法
		str, err = earKaiwa(items[85], items[323])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 聴力所見（文字）
		writeItems = append(writeItems, items[85])

		// 収縮期血圧（報告値）
		// 拡張期血圧（報告値）
		if times, err := ketsuatuTimes(items[86], items[87], items[88], items[89]); err != nil {
			logWrite(logstr, err)
			writeItems = append(writeItems, items[90])
			writeItems = append(writeItems, items[91])
		} else if times == 1 {
			writeItems = append(writeItems, items[90])
			writeItems = append(writeItems, items[91])
		} else {
			writeItems = append(writeItems, items[92])
			writeItems = append(writeItems, items[93])
		}

		// 収縮期血圧1回目
		writeItems = append(writeItems, items[90])

		// 収縮期血圧1回目
		writeItems = append(writeItems, items[91])

		// 収縮期血圧2回目
		writeItems = append(writeItems, items[92])

		// 収縮期血圧2回目
		writeItems = append(writeItems, items[93])

		// 脈拍数
		writeItems = append(writeItems, "")

		// 心電図実施区分
		writeItems = append(writeItems, "")

		// 心電図未実施理由
		writeItems = append(writeItems, "")

		// 心電図判定区分コード
		// 心電図判定区分名称
		strCd, strName, err = hanteiCdConv(items[365])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 心電図所見（文字）
		str = limitStr(joinStr5(items[94], items[95], items[96], items[97], items[98]), 256)
		writeItems = append(writeItems, str)

		// 心拍数
		writeItems = append(writeItems, items[99])

		// [Met]心電図所見有無
		str, err = syokenUmu(items[365])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]心電図対象者
		writeItems = append(writeItems, taisyo(items[365]))

		// [Met]心電図実施理由
		writeItems = append(writeItems, "")

		// 胸部X線実施区分
		writeItems = append(writeItems, "")

		// 胸部X線未実施理由
		writeItems = append(writeItems, "")

		// 胸部X線撮影区分
		writeItems = append(writeItems, satsuei(items[100], items[101]))

		// 胸部X線判定区分コード
		// 胸部X線判定区分名称
		strCd, strName, err = hanteiCdConv(items[359])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 胸部X線部位・所見（文字）
		str = limitStr(joinStr5(items[103], items[104], items[105], items[106], items[107]), 240)
		writeItems = append(writeItems, str)

		// 心胸比
		writeItems = append(writeItems, "")

		// [Met]胸部X線所見有無
		str, err = syokenUmu(items[359])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 胸部CT実施区分
		writeItems = append(writeItems, "")

		// 胸部CT未実施理由
		writeItems = append(writeItems, "")

		// 胸部CT判定区分コード
		// 胸部CT判定区分名称
		if items[108] != "" {
			strCd, strName, err = hanteiCdConv(items[482])
			logWrite(logstr, err)
			writeItems = append(writeItems, strCd)
			writeItems = append(writeItems, strName)
		} else {
			writeItems = append(writeItems, "")
			writeItems = append(writeItems, "")
		}

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 胸部CT部位・所見（文字）
		if items[108] != "" {
			str = limitStr(joinStr4(strings.TrimSpace(items[109]), strings.TrimSpace(items[110]), strings.TrimSpace(items[111]), strings.TrimSpace(items[112])), 240)
			writeItems = append(writeItems, str)
		} else {
			writeItems = append(writeItems, "")
		}

		// 喀痰実施区分
		writeItems = append(writeItems, "")

		// 喀痰未実施理由
		writeItems = append(writeItems, "")

		// 喀痰判定区分コード
		// 喀痰判定区分名称
		// 喀痰細胞診結果
		strCd, strName, str, err = kakutanConv(items[113])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)
		writeItems = append(writeItems, str)

		// 喀痰細胞診所見（文字）
		writeItems = append(writeItems, "")

		// 《予備》喀痰（抗酸菌）
		writeItems = append(writeItems, "")

		// 《予備》喀痰培養（ガフキー）
		writeItems = append(writeItems, "")

		// 肺活量
		writeItems = append(writeItems, items[114])

		// １秒量
		writeItems = append(writeItems, items[115])

		// 努力肺活量
		writeItems = append(writeItems, items[116])

		// １秒率
		writeItems = append(writeItems, items[117])

		// ％肺活量
		writeItems = append(writeItems, items[118])

		// ％１秒量
		writeItems = append(writeItems, items[119])

		// 肺機能換気障害区分
		writeItems = append(writeItems, "")

		// 眼底実施区分
		writeItems = append(writeItems, "")

		// 眼底未実施理由
		writeItems = append(writeItems, "")

		// 眼底判定区分
		// 眼底判定区分名称
		strCd, strName, err = hanteiCdConv(items[404])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// 眼底右シェイエ
		str = scheieConv(items[122], items[120])
		writeItems = append(writeItems, str)

		// 眼底左シェイエ
		str = scheieConv(items[123], items[124])
		writeItems = append(writeItems, str)

		// 予備（眼底）
		writeItems = append(writeItems, "")

		// 予備（眼底）
		writeItems = append(writeItems, "")

		// 眼底右Scott
		str, err = scottConv(items[126])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 眼底左Scott
		str, err = scottConv(items[127])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 眼底右KW
		str, err = kwConv(items[124])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 眼底左KW
		str, err = kwConv(items[125])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 眼底右Wong-Mitchell
		writeItems = append(writeItems, "")

		// 眼底左Wong-Mitchell
		writeItems = append(writeItems, "")

		// 眼底右Davis
		writeItems = append(writeItems, "")

		// 眼底左Davis
		writeItems = append(writeItems, "")

		// 眼底右その他所見（文字）
		str = limitStr(joinStr5(strings.TrimSpace(items[128]), strings.TrimSpace(items[129]), strings.TrimSpace(items[130]), strings.TrimSpace(items[131]), strings.TrimSpace(items[132])), 256)
		writeItems = append(writeItems, str)

		// 眼底左その他所見（文字）
		writeItems = append(writeItems, "")

		// [Met]眼底検査（対象者）
		writeItems = append(writeItems, taisyo(items[404]))

		// [Met]眼底検査（実施理由）
		writeItems = append(writeItems, "")

		// 予備
		writeItems = append(writeItems, "") // 予備

		// 眼圧右
		writeItems = append(writeItems, items[133])

		// 眼圧左
		writeItems = append(writeItems, items[134])

		// 腹部超音波実施区分
		writeItems = append(writeItems, "")

		// 腹部超音波未実施理由
		writeItems = append(writeItems, "")

		// 腹部超音波判定区分コード
		// 腹部超音波判定区分名称
		strCd, strName, err = hanteiCdConv(items[419])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 腹部超音波部位・所見（文字）
		str = limitStr(joinStr7(items[136], items[137], items[138], items[139], items[140], items[141], items[142]), 240)
		writeItems = append(writeItems, str)

		// 尿糖定性
		str, err = teiseiConv(items[143])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 尿蛋白定性
		str, err = teiseiConv(items[144])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 尿潜血定性
		str, err = teiseiConv(items[145])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 尿ウロビリノーゲン定性
		str, err = teiseiConv(items[146])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 尿比重
		writeItems = append(writeItems, items[147])

		// 尿pH
		writeItems = append(writeItems, items[148])

		// 尿沈渣判定区分コード
		// 尿沈渣判定区分名
		strCd, strName, err = hanteiCdConv(items[356])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// 尿沈渣赤血球
		writeItems = append(writeItems, items[149])

		// 尿沈渣白血球
		writeItems = append(writeItems, items[150])

		// 尿沈渣扁平上皮
		writeItems = append(writeItems, items[151])

		// 尿沈渣顆粒円柱
		writeItems = append(writeItems, items[152])

		// 尿沈渣ガラス円柱
		writeItems = append(writeItems, items[153])

		// 尿沈渣細菌
		// 尿沈渣その他
		saikin, sonota := nyoChinsaConv(items[154], items[155], items[156])
		writeItems = append(writeItems, saikin)
		writeItems = append(writeItems, sonota)

		// 赤血球数
		str, err = numChk(items[157])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 血色素量
		str, err = numChk(items[158])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// ヘマトクリット
		str, err = numChk(items[159])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 白血球数
		str, err = numChk(items[160])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 血小板数
		str, err = numChk(items[161])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// MCV
		str, err = numChk(items[162])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// MCH
		str, err = numChk(items[163])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// MCHC
		str, err = numChk(items[164])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]貧血検査（実施理由）
		writeItems = append(writeItems, "")

		// 血液像判定区分コード
		// 血液像判定区分名称
		strCd, strName, err = hanteiCdConv(items[377])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// 好中球(Neut)
		str, err = numChk(items[165])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 棹状核球(Stab)
		str, err = numChk(items[166])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 分葉核球(Seg)
		str, err = numChk(items[167])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 好酸球(Eosino)
		str, err = numChk(items[168])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 好塩基球(Baso)
		str, err = numChk(items[169])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// リンパ球(Lympho)
		str, err = numChk(items[170])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 単球(Mono)
		str, err = numChk(items[171])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 異形リンパ球(A-Lympho)
		writeItems = append(writeItems, "")

		// 骨髄球(Myelo)
		writeItems = append(writeItems, "")

		// 後骨髄球(Meta)
		writeItems = append(writeItems, "")

		// 白血球分画その他
		writeItems = append(writeItems, "")

		// その他の内容
		writeItems = append(writeItems, joinStr(items[172], items[173]))

		// 血清鉄
		str, err = numChk(items[174])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// フェリチン
		str, err = numChk(items[175])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 血液型ABO
		str, err = aboConv(items[176])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 血液型Rh
		str, err = rhConv(items[177])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 食後時間区分
		eatTime, err := eatTimeConv(items[203], items[178])
		logWrite(logstr, err)
		writeItems = append(writeItems, eatTime)

		// 生理区分
		str, err = seiriConv(items[179])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 妊娠区分
		str, err = ninshinConv(items[180], items[181])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 乳び
		str = nyubiConv(items[182], items[183])
		writeItems = append(writeItems, str)

		// 溶血
		str = yoketsuConv(items[182], items[183])
		writeItems = append(writeItems, str)

		// 血清総蛋白
		str, err = numChk(items[184])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 血清アルブミン
		str, err = numChk(items[185])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// A/G比
		str, err = numChk(items[186])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 尿中アルブミン
		writeItems = append(writeItems, "")

		// AST(GOT)
		str, err = numChk(items[187])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// ALT(GPT)
		str, err = numChk(items[188])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// γ-GTP
		str, err = numChk(items[189])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// ALP
		str, err = numChk(items[190])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// LDH
		str, err = numChk(items[191])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// コリンエステラーゼ
		str, err = numChk(items[192])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// LAP
		str, err = numChk(items[193])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 総ビリルビン
		str, err = numChk(items[194])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 直接ビリルビン
		str, err = numChk(items[195])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// CPK
		str, err = numChk(items[196])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　レベル区分
		writeItems = append(writeItems, "")

		// BNP
		str, err = numChk(items[197])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　レベル区分
		writeItems = append(writeItems, "")

		// 総コレステロール
		str, err = numChk(items[198])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// HDLコレステロール
		str, err = numChk(items[199])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// LDLコレステロール
		str, err = numChk(items[200])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 中性脂肪
		str, err = numChk(items[201])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// non-HDLコレステロール
		str, err = numChk(items[202])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 空腹時血糖
		// 随時血糖
		kufuku, zuiji := tohConv(items[203], eatTime)
		writeItems = append(writeItems, kufuku)
		writeItems = append(writeItems, zuiji)

		// HbA1c(NGSP)
		str, err = numChk(items[204])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 膵機能判定区分コード
		// 膵機能判定区分名称
		strCd, strName, err = hanteiCdConv(items[383])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// 血清アミラーゼ
		str, err = numChk(items[205])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　レベル区分
		writeItems = append(writeItems, "")

		// 膵アミラーゼ
		writeItems = append(writeItems, "")

		// 　レベル区分
		writeItems = append(writeItems, "")

		// 尿酸
		str, err = numChk(items[206])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 尿素窒素
		str, err = numChk(items[207])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 血清クレアチニン
		str, err = numChk(items[208])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// eGFR
		str, err = numChk(items[209])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]血清クレアチニン対象
		writeItems = append(writeItems, taisyo(items[208]))

		// [Met]血清クレアチニン実施理由
		writeItems = append(writeItems, "")

		// ナトリウム
		str, err = numChk(items[210])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// カリウム
		str, err = numChk(items[211])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// クロール
		str, err = numChk(items[212])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// カルシウム
		str, err = numChk(items[213])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// マグネシウム
		writeItems = append(writeItems, "")

		// 無機リン
		str, err = numChk(items[214])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// カンマ位置(331)
		writeItems = append(writeItems, "331")

		// 肝炎判定区分コード
		writeItems = append(writeItems, "")

		// 肝炎判定区分名称
		writeItems = append(writeItems, "")

		// HBs抗原定性
		str, err = teiseiConv(items[215])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// HBs抗体定性
		str, err = teiseiConv(items[217])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// HCV抗体定性
		str, err = teiseiConv(items[219])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// HBs抗原定量
		str, err = numChk(items[216])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　HBs抗原定量　陰・陽区分
		writeItems = append(writeItems, "")

		// HBs抗体定量
		str, err = numChk(items[218])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　HBs抗体定量　陰・陽区分
		writeItems = append(writeItems, "")

		// HCV抗体定量
		str, err = numChk(items[220])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　HCV抗体定量　陰・陽区分
		writeItems = append(writeItems, "")

		// CRP定性
		writeItems = append(writeItems, "")

		// CRP定量
		str, err = numChk(items[221])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　CRP定量　陰・陽区分
		writeItems = append(writeItems, "")

		// 高感度CRP
		writeItems = append(writeItems, "")

		// 　高感度CRP定量　陰・陽区分
		writeItems = append(writeItems, "")

		// RA(RF)定性
		writeItems = append(writeItems, "")

		// RF定量
		str, err = numChk(items[222])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　RF定量　陰・陽区分
		writeItems = append(writeItems, "")

		// 梅毒　総　陰・陽区分
		writeItems = append(writeItems, "")

		// 梅毒反応(TPHA)　定性
		str, err = teiseiConv(items[223])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 梅毒反応(TPHA)　定量
		writeItems = append(writeItems, "")

		// 　TPHA定量　陰・陽区分
		writeItems = append(writeItems, "")

		// 梅毒反応(RPR)　定性
		str, err = teiseiConv(items[224])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 梅毒反応(ガラス板)　定性
		writeItems = append(writeItems, "")

		// PSA定性
		writeItems = append(writeItems, "")

		// PSA定量
		str, err = numChk(items[225])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　PSA定量　陰・陽区分
		str, err = PSAconv(str)
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// CA125
		str, err = numChk(items[226])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　CA125　陰・陽区分
		str, err = CA125conv(str)
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// CA19_9
		str, err = numChk(items[227])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　CA19_9　陰・陽区分
		str, err = CA199conv(str)
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// CEA
		str, err = numChk(items[228])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　CEA　陰・陽区分
		str, err = CEAconv(str)
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// AFP
		str, err = numChk(items[229])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　AFP　陰・陽区分
		str, err = AFPconv(str)
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// シフラ
		str, err = numChk(items[230])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　シフラ　陰・陽区分
		str, err = SifuraConv(str)
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// TSH
		str, err = numChk(items[231])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　レベル区分
		writeItems = append(writeItems, "")

		// T3
		writeItems = append(writeItems, "")

		// 　レベル区分
		writeItems = append(writeItems, "")

		// T4
		writeItems = append(writeItems, "")

		// 　レベル区分
		writeItems = append(writeItems, "")

		// FT3
		str, err = numChk(items[232])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　レベル区分
		writeItems = append(writeItems, "")

		// FT4
		str, err = numChk(items[233])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 　レベル区分
		writeItems = append(writeItems, "")

		// 便中卵定性
		str, err = teiseiConv(items[234])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 便中卵所見
		writeItems = append(writeItems, "")

		// カンマ位置(382)
		writeItems = append(writeItems, "382")

		// 胃部X線実施区分
		writeItems = append(writeItems, "")

		// 胃部X線未実施理由
		writeItems = append(writeItems, "")

		// 胃部X線判定区分コード
		// 胃部X線判定区分名称
		strCd, strName, err = hanteiCdConv(items[410])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 胃部X線撮影区分
		writeItems = append(writeItems, satsuei(items[235], items[236]))

		// 胃部X線部位・所見（文字）
		str = limitStr(joinStr5(items[238], items[239], items[240], items[241], items[242]), 240)
		writeItems = append(writeItems, str)

		// 胃カメラ実施区分
		writeItems = append(writeItems, "")

		// 胃カメラ未実施理由
		writeItems = append(writeItems, "")

		// 胃カメラ判定区分コード
		// 胃カメラ判定区分名称
		strCd, strName, err = hanteiCdConv(items[413])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 胃部内視鏡部位・所見（文字）
		str = limitStr(joinStr5(items[243], items[244], items[245], items[246], items[247]), 240)
		writeItems = append(writeItems, str)

		// 胃部内視鏡組織検査実施区分
		writeItems = append(writeItems, "")

		// 胃部内視鏡組織・生検所見
		str = limitStr(joinStr(items[248], items[249]), 240)
		writeItems = append(writeItems, str)

		// PG・ピロリ判定区分コード
		// PG・ピロリ判定区分名称
		str, err = hantiHeavy(items[434], items[437])
		logWrite(logstr, err)

		strCd, strName, err = hanteiCdConv(str)
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// ABC検診判定分類
		str, err = iabcConv(items[255])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// PGⅠ
		str, err = numChk(items[250])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// PGⅡ
		str, err = numChk(items[251])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// PGⅠ/Ⅱ比
		str, err = numChk(items[252])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// PG比　陰・陽区分
		writeItems = append(writeItems, "")

		// ピロリIgG抗体定量
		str, err = numChk(items[254])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// ピロリIgG抗体定量　陰・陽区分
		str, err = teiseiConv(items[253])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 尿中ピロリ菌抗体定性
		writeItems = append(writeItems, "")

		// 呼気ピロリ菌抗体定性
		writeItems = append(writeItems, "")

		// PGに関する所見
		writeItems = append(writeItems, "")

		// 大腸内視鏡実施区分
		writeItems = append(writeItems, "")

		// 大腸内視鏡未実施理由
		writeItems = append(writeItems, "")

		// 大腸内視鏡判定区分コード
		writeItems = append(writeItems, "")

		// 大腸内視鏡判定区分名称
		writeItems = append(writeItems, "")

		// （予備）留意所見有無区
		writeItems = append(writeItems, "")

		// 大腸内視鏡部位・所見（文字）
		writeItems = append(writeItems, "")

		// 直腸診実施区分
		writeItems = append(writeItems, "")

		// 直腸診未実施区分
		writeItems = append(writeItems, "")

		// 直腸診判定区分コー
		writeItems = append(writeItems, "")

		// 直腸診判定区分名称
		writeItems = append(writeItems, "")

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 直腸診部位・所見（文字）
		writeItems = append(writeItems, "")

		// 便潜血実施区分
		writeItems = append(writeItems, "")

		// 便潜血未実施理由
		writeItems = append(writeItems, "")

		// 便潜血判定区分コード
		// 便潜血判定区分名称
		strCd, strName, err = hanteiCdConv(items[422])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// 便潜血１回目（定性）
		str, err = teiseiConv(items[256])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 便潜血２回目（定性）
		str, err = teiseiConv(items[257])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 便潜血１回目定量
		writeItems = append(writeItems, "")

		// 　１回目定量　陰・陽区
		writeItems = append(writeItems, "")

		// 便潜血２回目定量
		writeItems = append(writeItems, "")

		// 　２回目定量　陰・陽区分
		writeItems = append(writeItems, "")

		// カンマ位置(432)
		writeItems = append(writeItems, "432")

		// 乳がん総判定区分コード
		// 乳がん総判定区分名称
		str, err = hantiHeavy(items[452], items[455]) // 乳房エコーとマンモの判定
		logWrite(logstr, err)

		strCd, strName, err = hanteiCdConv(str)
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 乳がん総合所見（文字）
		writeItems = append(writeItems, "")

		// 乳房視触診（文字）
		writeItems = append(writeItems, "")

		// 乳腺エコー実施区分
		writeItems = append(writeItems, "")

		// 乳腺エコー未実施理由
		writeItems = append(writeItems, "")

		// 乳腺エコー判定区分コード
		// 乳腺エコー判定区分名称
		strCd, strName, err = hanteiCdConv(items[452])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 乳腺エコー所見（文字）
		str = limitStr(joinStr3(items[258], items[259], items[260]), 240)
		writeItems = append(writeItems, str)

		// マンモ実施区分
		writeItems = append(writeItems, "")

		// マンモ未実施理由
		writeItems = append(writeItems, "")

		// マンモ判定区分コード
		// マンモ判定区分名称
		strCd, strName, err = hanteiCdConv(items[455])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// マンモ撮影方向
		str = mmgSatsuei(items[261], items[262])
		writeItems = append(writeItems, str)

		// マンモ所見（文字）
		str = limitStr(joinStr3(items[263], items[264], items[265]), 240)
		writeItems = append(writeItems, str)

		// 子宮頸部細胞診実施区分
		writeItems = append(writeItems, "")

		// 子宮頸部細胞診未実施区分
		writeItems = append(writeItems, "")

		// 子宮頸部細胞診判定区分コード
		// 子宮頸部細胞診判定区分名称
		str, err = hantiHeavy(items[458], items[461]) // 婦人科内診と子宮細胞診の判定
		logWrite(logstr, err)

		strCd, strName, err = hanteiCdConv(str)
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 子宮内診所見（文字）
		str = limitStr(joinStr3(items[268], items[269], items[270]), 240)
		writeItems = append(writeItems, str)

		// 子宮頸部細胞診（ベセスダ）
		strCd, err = vesesudaConv(items[266])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)

		// 子宮頸部細胞診（日母分類）
		strCd, err = nichimoConv(items[267])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)

		// 子宮頸部細胞診結果
		writeItems = append(writeItems, "")

		// HPV
		writeItems = append(writeItems, "")

		// 子宮超音波実施区分
		writeItems = append(writeItems, "")

		// 子宮超音波未実施理由
		writeItems = append(writeItems, "")

		// 子宮超音波判定区分コード
		writeItems = append(writeItems, "")

		// 子宮超音波判定区分名称
		writeItems = append(writeItems, "")

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 子宮超音波所見（文字）
		writeItems = append(writeItems, "")

		// 骨密度(BMD)
		str, err = numChk(items[274])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// YAM
		writeItems = append(writeItems, "")

		// 同性年代平均値比
		writeItems = append(writeItems, "")

		// 骨密度検査その他
		writeItems = append(writeItems, "")

		// 心臓超音波実施区分
		writeItems = append(writeItems, "")

		// 心臓超音波未実施理由
		writeItems = append(writeItems, "")

		// 心臓超音波判定区分コード
		// 心臓超音波判定区分名称
		strCd, strName, err = hanteiCdConv(items[467])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// 心臓超音波所見（文字）
		str = limitStr(joinStr4(items[275], items[276], items[277], items[278]), 240)
		writeItems = append(writeItems, str)

		// ABI 右
		str, err = numChk(items[279])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// ABI 左
		str, err = numChk(items[280])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// PWV 右
		writeItems = append(writeItems, "")

		// PWV 左
		writeItems = append(writeItems, "")

		// CAVI 右
		str, err = numChk(items[281])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// CAVI 左
		str, err = numChk(items[282])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// 脳ドック実施区分
		writeItems = append(writeItems, "")

		// 脳ドック検査種別
		writeItems = append(writeItems, "")

		// 脳ドック総判定区分コード
		writeItems = append(writeItems, "")

		// 脳ドック総判定区分名称
		writeItems = append(writeItems, "")

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 脳ドック所見（文字）
		writeItems = append(writeItems, "")

		// 頸動脈超音波実施区分
		writeItems = append(writeItems, "")

		// 頸動脈超音波判定区分コード
		// 頸動脈超音波判定区分名称
		strCd, strName, err = hanteiCdConv(items[473])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 頸動脈超音波所見（文字）
		str = limitStr(joinStr3(items[283], items[284], items[285]), 240)
		writeItems = append(writeItems, str)

		// 甲状腺超音波実施区分
		writeItems = append(writeItems, "")

		// 甲状腺超音波判定区分コード
		// 甲状腺超音波判定区分名称
		strCd, strName, err = hanteiCdConv(items[476])
		logWrite(logstr, err)
		writeItems = append(writeItems, strCd)
		writeItems = append(writeItems, strName)

		// （予備）留意所見有無区分
		writeItems = append(writeItems, "")

		// 甲状腺超音波部位所見（文字）
		str = limitStr(joinStr4(items[286], items[287], items[288], items[289]), 240)
		writeItems = append(writeItems, str)

		// [Met]既往歴有無
		// [Met]具体的な既往歴
		kiou1 := kiouJoin(items[30], items[31], items[32])
		kiou2 := kiouJoin(items[33], items[34], items[35])
		kiou3 := kiouJoin(items[36], items[37], items[38])
		kiou4 := kiouJoin(items[39], items[40], items[41])
		kiou5 := kiouJoin(items[42], items[43], items[44])
		kiou6 := kiouJoin(items[45], items[46], items[47])
		kiou7 := kiouJoin(items[48], items[49], items[50])
		kiou8 := kiouJoin(items[51], items[52], items[53])
		kiou9 := kiouJoin(items[54], items[55], items[56])
		kiou10 := kiouJoin(items[57], items[58], items[59])
		str = limitStr(joinStr10(kiou1, kiou2, kiou3, kiou4, kiou5, kiou6, kiou7, kiou8, kiou9, kiou10), 256)
		writeItems = append(writeItems, umuConv(str))
		writeItems = append(writeItems, str)

		// [Met]自覚症状の有無
		// [Met]具体的な自覚症状
		str = limitStr(joinStr5(items[25], items[26], items[27], items[28], items[29]), 256)
		writeItems = append(writeItems, jikakuUmu(str))
		writeItems = append(writeItems, str)

		// [Met]他覚症状の有無
		// [Met]具体的な他覚症状
		str = limitStr(joinStr3(items[22], items[23], items[24]), 256)
		writeItems = append(writeItems, takakuUmu(str))
		writeItems = append(writeItems, str)

		// [Met]高血圧（服薬有無）
		str, err = yesNoConv(items[290])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]高血圧（薬剤名）
		writeItems = append(writeItems, "")

		// [Met]高血圧（服薬理由）
		writeItems = append(writeItems, "")

		// [Met]糖尿病（服薬有無）
		str, err = yesNoConv(items[291])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]糖尿病（薬剤名）
		writeItems = append(writeItems, "")

		// [Met]糖尿病（服薬理由）
		writeItems = append(writeItems, "")

		// [Met]脂質（服薬有無）
		str, err = yesNoConv(items[292])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]脂質（薬剤名）
		writeItems = append(writeItems, "")

		// [Met]脂質（服薬理由）
		writeItems = append(writeItems, "")

		// [Met]既往歴１（脳血管有無）
		str, err = yesNoConv(items[293])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]既往歴２（心血管有無）
		str, err = yesNoConv(items[294])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]既往歴３（腎不全・人口透析有無）
		str, err = yesNoConv(items[295])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]貧血既往有無
		str, err = yesNoConv(items[296])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]習慣的喫煙
		str, err = yesNoConv(items[297])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]喫煙本数／日
		writeItems = append(writeItems, "")

		// [Met]喫煙期間（年）
		writeItems = append(writeItems, "")

		// [Met]20歳からの体重変化
		str, err = yesNoConv(items[298])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]30分以上の運動習慣
		str, err = yesNoConv(items[299])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]歩行又は身体活動
		str, err = yesNoConv(items[300])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]歩行速度
		str, err = yesNoConv(items[301])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]咀嚼
		str, err = sosyakuConv(items[302])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]食べ方１（早食い等）
		str, err = eat1Conv(items[303])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]食べ方２（就寝前）
		str, err = yesNoConv(items[304])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]食べ方３（間食）
		str, err = eat3Conv(items[305])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]食習慣（朝食）
		str, err = yesNoConv(items[306])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]飲酒習慣
		str, err = sakeConv(items[307])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]飲酒量
		str, err = sakeryoConv(items[308])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]睡眠
		str, err = yesNoConv(items[309])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]生活習慣の改善意志
		str, err = seikatsuConv(items[310])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]保健指導の希望
		str, err = yesNoConv(items[311])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]保健指導レベル
		str, err = hokenConv(items[312])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]メタボリックシンドローム判定
		str, err = metaboConv(items[313])
		logWrite(logstr, err)
		writeItems = append(writeItems, str)

		// [Met]医師の診断（特定健診）
		writeItems = append(writeItems, items[315])

		// 初回面接実施
		writeItems = append(writeItems, "")

		// 初回面接補足内容
		writeItems = append(writeItems, "")

		// 情報提供の方法
		writeItems = append(writeItems, "")

		// カンマ位置(540)
		writeItems = append(writeItems, "540")

		writer.Write(writeItems) // 1行書き出す
	}

	writer.Flush()
	log.Print("Finesh !\r\n")
}

func titleWrite() []string {
	var title []string
	title = append(title, "CSVフォーマットVer")
	title = append(title, "提出先")
	title = append(title, "データ作成者")
	title = append(title, "データ作成日")
	title = append(title, "データ提出日")
	title = append(title, "データ登録完了区分")
	title = append(title, "登録未完了の連絡内容")
	title = append(title, "団体コード")
	title = append(title, "団体コード名称")
	title = append(title, "事業所コード")
	title = append(title, "事業所名称")
	title = append(title, "個人ID")
	title = append(title, "漢字氏名")
	title = append(title, "カナ氏名")
	title = append(title, "生年月日")
	title = append(title, "性別")
	title = append(title, "保険者番号")
	title = append(title, "保険証記号")
	title = append(title, "保険証番号")
	title = append(title, "続柄")
	title = append(title, "予備")
	title = append(title, "予備")
	title = append(title, "受診券整理番号")
	title = append(title, "受診券有効期限")
	title = append(title, "コースコード")
	title = append(title, "コース名称")
	title = append(title, "受診日")
	title = append(title, "施設/巡回区分")
	title = append(title, "健診機関コード")
	title = append(title, "健診機関名称")
	title = append(title, "[Met]特定健診機関番号")
	title = append(title, "[Met]健診実施医師名")
	title = append(title, "予備")
	title = append(title, "予備")
	title = append(title, "産業医判定区分")
	title = append(title, "就労区分")
	title = append(title, "産業医コメント")
	title = append(title, "伝達事項有無")
	title = append(title, "伝達内容")
	title = append(title, "診察判定区分コード")
	title = append(title, "診察判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "診察所見")
	title = append(title, "自覚症状など")
	title = append(title, "治療中疾病有無区分")
	title = append(title, "治療中疾病名（文字）")
	title = append(title, "既往疾病有無区分")
	title = append(title, "既往疾病名")
	title = append(title, "総合判定区分コード")
	title = append(title, "総合判定区分名称")
	title = append(title, "総合判定コメント")
	title = append(title, "予備")
	title = append(title, "予備")
	title = append(title, "予備①(1)")
	title = append(title, "予備②(1)")
	title = append(title, "予備③(1)")
	title = append(title, "予備①(2)")
	title = append(title, "予備②(2)")
	title = append(title, "予備③(2)")
	title = append(title, "予備①(3)")
	title = append(title, "予備②(3)")
	title = append(title, "予備③(3)")
	title = append(title, "予備①(4)")
	title = append(title, "予備②(4)")
	title = append(title, "予備③(4)")
	title = append(title, "予備①(5)")
	title = append(title, "予備②(5)")
	title = append(title, "予備③(5)")
	title = append(title, "予備①(6)")
	title = append(title, "予備②(6)")
	title = append(title, "予備③(6)")
	title = append(title, "予備①(7)")
	title = append(title, "予備②(7)")
	title = append(title, "予備③(7)")
	title = append(title, "予備①(8)")
	title = append(title, "予備②(8)")
	title = append(title, "予備③(8)")
	title = append(title, "予備①(9)")
	title = append(title, "予備②(9)")
	title = append(title, "予備③(9)")
	title = append(title, "予備①(10)")
	title = append(title, "予備②(10)")
	title = append(title, "予備③(10)")
	title = append(title, "予備①(11)")
	title = append(title, "予備②(11)")
	title = append(title, "予備③(11)")
	title = append(title, "予備①(12)")
	title = append(title, "予備②(12)")
	title = append(title, "予備③(12)")
	title = append(title, "予備①(13)")
	title = append(title, "予備②(13)")
	title = append(title, "予備③(13)")
	title = append(title, "予備①(14)")
	title = append(title, "予備②(14)")
	title = append(title, "予備③(14)")
	title = append(title, "予備①(15)")
	title = append(title, "予備②(15)")
	title = append(title, "予備③(15)")
	title = append(title, "予備①(16)")
	title = append(title, "予備②(16)")
	title = append(title, "予備③(16)")
	title = append(title, "予備①(17)")
	title = append(title, "予備②(17)")
	title = append(title, "予備③(17)")
	title = append(title, "予備①(18)")
	title = append(title, "予備②(18)")
	title = append(title, "予備③(18)")
	title = append(title, "予備①(19)")
	title = append(title, "予備②(19)")
	title = append(title, "予備③(19)")
	title = append(title, "予備①(20)")
	title = append(title, "予備②(20)")
	title = append(title, "予備③(20)")
	title = append(title, "予備①(21)")
	title = append(title, "予備②(21)")
	title = append(title, "予備③(21)")
	title = append(title, "予備①(22)")
	title = append(title, "予備②(22)")
	title = append(title, "予備③(22)")
	title = append(title, "予備①(23)")
	title = append(title, "予備②(23)")
	title = append(title, "予備③(23)")
	title = append(title, "予備①(24)")
	title = append(title, "予備②(24)")
	title = append(title, "予備③(24)")
	title = append(title, "予備")
	title = append(title, "予備")
	title = append(title, "その他判定区分コード")
	title = append(title, "その他判定区分名称")
	title = append(title, "その他データ内容")
	title = append(title, "カンマ位置(131)")
	title = append(title, "身長")
	title = append(title, "体重")
	title = append(title, "BMI")
	title = append(title, "腹囲")
	title = append(title, "体脂肪率")
	title = append(title, "内臓脂肪面積")
	title = append(title, "5m視力裸眼右")
	title = append(title, "　データ属性")
	title = append(title, "5m視力裸眼左")
	title = append(title, "　データ属性")
	title = append(title, "5m視力矯正右")
	title = append(title, "　データ属性")
	title = append(title, "5m視力矯正左")
	title = append(title, "　データ属性")
	title = append(title, "近点視力裸眼右")
	title = append(title, "　データ属性")
	title = append(title, "近点視力裸眼左")
	title = append(title, "　データ属性")
	title = append(title, "近点視力矯正右")
	title = append(title, "　データ属性")
	title = append(title, "近点視力矯正左")
	title = append(title, "　データ属性")
	title = append(title, "視力矯正区分")
	title = append(title, "聴力右1K所見区分")
	title = append(title, "聴力右1K(dB)")
	title = append(title, "聴力左1K所見区分")
	title = append(title, "聴力左1K(dB)")
	title = append(title, "聴力右4K所見区分")
	title = append(title, "聴力右4K(dB)")
	title = append(title, "聴力左4K所見区分")
	title = append(title, "聴力左4K(dB)")
	title = append(title, "聴力会話法")
	title = append(title, "聴力所見（文字）")
	title = append(title, "収縮期血圧（報告値）")
	title = append(title, "拡張期血圧（報告値）")
	title = append(title, "収縮期血圧1回目")
	title = append(title, "拡張期血圧1回目")
	title = append(title, "収縮期血圧2回目")
	title = append(title, "拡張期血圧2回目")
	title = append(title, "脈拍数")
	title = append(title, "心電図実施区分")
	title = append(title, "心電図未実施理由")
	title = append(title, "心電図判定区分コード")
	title = append(title, "心電図判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "心電図所見（文字）")
	title = append(title, "心拍数")
	title = append(title, "[Met]心電図所見有無")
	title = append(title, "[Met]心電図対象者")
	title = append(title, "[Met]心電図実施理由")
	title = append(title, "胸部X線実施区分")
	title = append(title, "胸部X線未実施理由")
	title = append(title, "胸部X線撮影区分")
	title = append(title, "胸部X線判定区分コード")
	title = append(title, "胸部X線判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "胸部X線部位・所見（文字）")
	title = append(title, "心胸比")
	title = append(title, "[Met]胸部X線所見有無")
	title = append(title, "胸部CT実施区分")
	title = append(title, "胸部CT未実施理由")
	title = append(title, "胸部CT判定区分コード")
	title = append(title, "胸部CT判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "胸部CT部位・所見（文字）")
	title = append(title, "喀痰実施区分")
	title = append(title, "喀痰未実施理由")
	title = append(title, "喀痰判定区分コード")
	title = append(title, "喀痰判定区分名称")
	title = append(title, "喀痰細胞診結果")
	title = append(title, "喀痰細胞診所見（文字）")
	title = append(title, "《予備》喀痰（抗酸菌）")
	title = append(title, "《予備》喀痰培養（ガフキー）")
	title = append(title, "肺活量")
	title = append(title, "１秒量")
	title = append(title, "努力肺活量")
	title = append(title, "１秒率")
	title = append(title, "％肺活量")
	title = append(title, "％１秒量")
	title = append(title, "肺機能換気障害区分")
	title = append(title, "眼底実施区分")
	title = append(title, "眼底未実施理由")
	title = append(title, "眼底判定区分")
	title = append(title, "眼底判定区分名称")
	title = append(title, "眼底右シェイエ")
	title = append(title, "眼底左シェイエ")
	title = append(title, "予備（眼底）")
	title = append(title, "予備（眼底）")
	title = append(title, "眼底右Scott")
	title = append(title, "眼底左Scott")
	title = append(title, "眼底右KW")
	title = append(title, "眼底左KW")
	title = append(title, "眼底右Wong-Mitchell")
	title = append(title, "眼底左Wong-Mitchell")
	title = append(title, "眼底右Davis")
	title = append(title, "眼底左Davis")
	title = append(title, "眼底右その他所見（文字）")
	title = append(title, "眼底左その他所見（文字）")
	title = append(title, "[Met]眼底検査（対象者）")
	title = append(title, "[Met]眼底検査（実施理由）")
	title = append(title, "予備")
	title = append(title, "眼圧右")
	title = append(title, "眼圧左")
	title = append(title, "腹部超音波実施区分")
	title = append(title, "腹部超音波未実施理由")
	title = append(title, "腹部超音波判定区分コード")
	title = append(title, "腹部超音波判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "腹部超音波部位・所見（文字）")
	title = append(title, "尿糖定性")
	title = append(title, "尿蛋白定性")
	title = append(title, "尿潜血定性")
	title = append(title, "尿ウロビリノーゲン定性")
	title = append(title, "尿比重")
	title = append(title, "尿pH")
	title = append(title, "尿沈渣判定区分コード")
	title = append(title, "尿沈渣判定区分名称")
	title = append(title, "尿沈渣赤血球")
	title = append(title, "尿沈渣白血球")
	title = append(title, "尿沈渣扁平上皮")
	title = append(title, "尿沈渣顆粒円柱")
	title = append(title, "尿沈渣ガラス円柱")
	title = append(title, "尿沈渣細菌")
	title = append(title, "尿沈渣その他")
	title = append(title, "赤血球数")
	title = append(title, "血色素量")
	title = append(title, "ヘマトクリット")
	title = append(title, "白血球数")
	title = append(title, "血小板数")
	title = append(title, "MCV")
	title = append(title, "MCH")
	title = append(title, "MCHC")
	title = append(title, "[Met]貧血検査（実施理由）")
	title = append(title, "血液像判定区分コード")
	title = append(title, "血液像判定区分名称")
	title = append(title, "好中球(Neut)")
	title = append(title, "棹状核球(Stab)")
	title = append(title, "分葉核球(Seg)")
	title = append(title, "好酸球(Eosino)")
	title = append(title, "好塩基球(Baso)")
	title = append(title, "リンパ球(Lympho)")
	title = append(title, "単球(Mono)")
	title = append(title, "異形リンパ球(A-Lympho)")
	title = append(title, "骨髄球(Myelo)")
	title = append(title, "後骨髄球(Meta)")
	title = append(title, "白血球分画その他")
	title = append(title, "その他の内容")
	title = append(title, "血清鉄")
	title = append(title, "フェリチン")
	title = append(title, "血液型ABO")
	title = append(title, "血液型Rh")
	title = append(title, "食後時間区分")
	title = append(title, "生理区分")
	title = append(title, "妊娠区分")
	title = append(title, "乳び")
	title = append(title, "溶血")
	title = append(title, "血清総蛋白")
	title = append(title, "血清アルブミン")
	title = append(title, "A/G比")
	title = append(title, "尿中アルブミン")
	title = append(title, "AST(GOT)")
	title = append(title, "ALT(GPT)")
	title = append(title, "γ-GTP")
	title = append(title, "ALP")
	title = append(title, "LDH")
	title = append(title, "コリンエステラーゼ")
	title = append(title, "LAP")
	title = append(title, "総ビリルビン")
	title = append(title, "直接ビリルビン")
	title = append(title, "CPK")
	title = append(title, "　レベル区分")
	title = append(title, "BNP")
	title = append(title, "　レベル区分")
	title = append(title, "総コレステロール")
	title = append(title, "HDLコレステロール")
	title = append(title, "LDLコレステロール")
	title = append(title, "中性脂肪")
	title = append(title, "non-HDLコレステロール")
	title = append(title, "空腹時血糖")
	title = append(title, "随時血糖")
	title = append(title, "HbA1c(NGSP)")
	title = append(title, "膵機能判定区分コード")
	title = append(title, "膵機能判定区分名称")
	title = append(title, "血清アミラーゼ")
	title = append(title, "　レベル区分")
	title = append(title, "膵アミラーゼ")
	title = append(title, "　レベル区分")
	title = append(title, "尿酸")
	title = append(title, "尿素窒素")
	title = append(title, "血清クレアチニン")
	title = append(title, "eGFR")
	title = append(title, "[Met]血清クレアチニン対象")
	title = append(title, "[Met]血清クレアチニン実施理由")
	title = append(title, "ナトリウム")
	title = append(title, "カリウム")
	title = append(title, "クロール")
	title = append(title, "カルシウム")
	title = append(title, "マグネシウム")
	title = append(title, "無機リン")
	title = append(title, "カンマ位置(331)")
	title = append(title, "肝炎判定区分コード")
	title = append(title, "肝炎判定区分名称")
	title = append(title, "HBs抗原定性")
	title = append(title, "HBs抗体定性")
	title = append(title, "HCV抗体定性")
	title = append(title, "HBs抗原定量")
	title = append(title, "　HBs抗原定量　陰・陽区分")
	title = append(title, "HBs抗体定量")
	title = append(title, "　HBs抗体定量　陰・陽区分")
	title = append(title, "HCV抗体定量")
	title = append(title, "　HCV抗体定量　陰・陽区分")
	title = append(title, "CRP定性")
	title = append(title, "CRP定量")
	title = append(title, "　CRP定量　陰・陽区分")
	title = append(title, "高感度CRP")
	title = append(title, "　高感度CRP定量　陰・陽区分")
	title = append(title, "RA(RF)定性")
	title = append(title, "RF定量")
	title = append(title, "　RF定量　陰・陽区分")
	title = append(title, "梅毒　総　陰・陽区分")
	title = append(title, "梅毒反応(TPHA)　定性")
	title = append(title, "梅毒反応(TPHA)　定量")
	title = append(title, "　TPHA定量　陰・陽区分")
	title = append(title, "梅毒反応(RPR)　定性")
	title = append(title, "梅毒反応(ガラス板)　定性")
	title = append(title, "PSA定性")
	title = append(title, "PSA定量")
	title = append(title, "　PSA定量　陰・陽区分")
	title = append(title, "CA125")
	title = append(title, "　CA125　陰・陽区分")
	title = append(title, "CA19_9")
	title = append(title, "　CA19_9　陰・陽区分")
	title = append(title, "CEA")
	title = append(title, "　CEA　陰・陽区分")
	title = append(title, "AFP")
	title = append(title, "　AFP　陰・陽区分")
	title = append(title, "シフラ")
	title = append(title, "　シフラ　陰・陽区分")
	title = append(title, "TSH")
	title = append(title, "　レベル区分")
	title = append(title, "T3")
	title = append(title, "　レベル区分")
	title = append(title, "T4")
	title = append(title, "　レベル区分")
	title = append(title, "FT3")
	title = append(title, "　レベル区分")
	title = append(title, "FT4")
	title = append(title, "　レベル区分")
	title = append(title, "便中卵定性")
	title = append(title, "便中卵所見")
	title = append(title, "カンマ位置(382)")
	title = append(title, "胃部X線実施区分")
	title = append(title, "胃部X線未実施理由")
	title = append(title, "胃部X線判定区分コード")
	title = append(title, "胃部X線判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "胃部X線撮影区分")
	title = append(title, "胃部X線部位・所見（文字）")
	title = append(title, "胃カメラ実施区分")
	title = append(title, "胃カメラ未実施理由")
	title = append(title, "胃カメラ判定区分コード")
	title = append(title, "胃カメラ判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "胃部内視鏡部位・所見（文字）")
	title = append(title, "胃部内視鏡組織検査実施区分")
	title = append(title, "胃部内視鏡組織・生検所見")
	title = append(title, "PG・ピロリ判定区分コード")
	title = append(title, "PG・ピロリ判定区分名称")
	title = append(title, "ABC検診判定分類")
	title = append(title, "PGⅠ")
	title = append(title, "PGⅡ")
	title = append(title, "PGⅠ/Ⅱ比")
	title = append(title, "PG比　陰・陽区分")
	title = append(title, "ピロリIgG抗体定量")
	title = append(title, "ピロリIgG抗体定量　陰・陽区分")
	title = append(title, "尿中ピロリ菌抗体定性")
	title = append(title, "呼気ピロリ菌抗体定性")
	title = append(title, "PGに関する所見")
	title = append(title, "大腸内視鏡実施区分")
	title = append(title, "大腸内視鏡未実施理由")
	title = append(title, "大腸内視鏡判定区分コード")
	title = append(title, "大腸内視鏡判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "大腸内視鏡部位・所見（文字）")
	title = append(title, "直腸診実施区分")
	title = append(title, "直腸診未実施区分")
	title = append(title, "直腸診判定区分コード")
	title = append(title, "直腸診判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "直腸診部位・所見（文字）")
	title = append(title, "便潜血実施区分")
	title = append(title, "便潜血未実施理由")
	title = append(title, "便潜血判定区分コード")
	title = append(title, "便潜血判定区分名称")
	title = append(title, "便潜血１回目（定性）")
	title = append(title, "便潜血２回目（定性）")
	title = append(title, "便潜血１回目定量")
	title = append(title, "　１回目定量　陰・陽区分")
	title = append(title, "便潜血２回目定量")
	title = append(title, "　２回目定量　陰・陽区分")
	title = append(title, "カンマ位置(432)")
	title = append(title, "乳がん総判定区分コード")
	title = append(title, "乳がん総判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "乳がん総合所見（文字）")
	title = append(title, "乳房視触診（文字）")
	title = append(title, "乳腺エコー実施区分")
	title = append(title, "乳腺エコー未実施理由")
	title = append(title, "乳腺エコー判定区分コード")
	title = append(title, "乳腺エコー判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "乳腺エコー所見（文字）")
	title = append(title, "マンモ実施区分")
	title = append(title, "マンモ未実施理由")
	title = append(title, "マンモ判定区分コード")
	title = append(title, "マンモ判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "マンモ撮影方向")
	title = append(title, "マンモ所見（文字）")
	title = append(title, "子宮頸部細胞診実施区分")
	title = append(title, "子宮頸部細胞診未実施区分")
	title = append(title, "子宮頸部細胞診判定区分コード")
	title = append(title, "子宮頸部細胞診判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "子宮内診所見（文字）")
	title = append(title, "子宮頸部細胞診（ベセスダ）")
	title = append(title, "子宮頸部細胞診（日母分類）")
	title = append(title, "子宮頸部細胞診結果")
	title = append(title, "HPV")
	title = append(title, "子宮超音波実施区分")
	title = append(title, "子宮超音波未実施理由")
	title = append(title, "子宮超音波判定区分コード")
	title = append(title, "子宮超音波判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "子宮超音波所見（文字）")
	title = append(title, "骨密度(BMD)")
	title = append(title, "YAM")
	title = append(title, "同性年代平均値比")
	title = append(title, "骨密度検査その他")
	title = append(title, "心臓超音波実施区分")
	title = append(title, "心臓超音波未実施理由")
	title = append(title, "心臓超音波判定区分コード")
	title = append(title, "心臓超音波判定区分名称")
	title = append(title, "心臓超音波所見（文字）")
	title = append(title, "ABI 右")
	title = append(title, "ABI 左")
	title = append(title, "PWV 右")
	title = append(title, "PWV 左")
	title = append(title, "CAVI 右")
	title = append(title, "CAVI 左")
	title = append(title, "脳ドック実施区分")
	title = append(title, "脳ドック検査種別")
	title = append(title, "脳ドック総判定区分コード")
	title = append(title, "脳ドック総判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "脳ドック所見（文字）")
	title = append(title, "頸動脈超音波実施区分")
	title = append(title, "頸動脈超音波判定区分コード")
	title = append(title, "頸動脈超音波判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "頸動脈超音波所見（文字）")
	title = append(title, "甲状腺超音波実施区分")
	title = append(title, "甲状腺超音波判定区分コード")
	title = append(title, "甲状腺超音波判定区分名称")
	title = append(title, "（予備）留意所見有無区分")
	title = append(title, "甲状腺超音波部位所見（文字）")
	title = append(title, "[Met]既往歴有無")
	title = append(title, "[Met]具体的な既往歴")
	title = append(title, "[Met]自覚症状の有無")
	title = append(title, "[Met]具体的な自覚症状")
	title = append(title, "[Met]他覚症状の有無")
	title = append(title, "[Met]具体的な他覚症状")
	title = append(title, "[Met]高血圧（服薬有無）")
	title = append(title, "[Met]高血圧（薬剤名）")
	title = append(title, "[Met]高血圧（服薬理由）")
	title = append(title, "[Met]糖尿病（服薬有無）")
	title = append(title, "[Met]糖尿病（薬剤名）")
	title = append(title, "[Met]糖尿病（服薬理由）")
	title = append(title, "[Met]脂質（服薬有無）")
	title = append(title, "[Met]脂質（薬剤名）")
	title = append(title, "[Met]脂質（服薬理由）")
	title = append(title, "[Met]既往歴１（脳血管有無）")
	title = append(title, "[Met]既往歴２（心血管有無）")
	title = append(title, "[Met]既往歴３（腎不全・人口透析有無）")
	title = append(title, "[Met]貧血既往有無")
	title = append(title, "[Met]習慣的喫煙")
	title = append(title, "[Met]喫煙本数／日")
	title = append(title, "[Met]喫煙期間（年）")
	title = append(title, "[Met]20歳からの体重変化")
	title = append(title, "[Met]30分以上の運動習慣")
	title = append(title, "[Met]歩行又は身体活動")
	title = append(title, "[Met]歩行速度")
	title = append(title, "[Met]咀嚼")
	title = append(title, "[Met]食べ方１（早食い等）")
	title = append(title, "[Met]食べ方２（就寝前）")
	title = append(title, "[Met]食べ方３（間食）")
	title = append(title, "[Met]食習慣（朝食）")
	title = append(title, "[Met]飲酒習慣")
	title = append(title, "[Met]飲酒量")
	title = append(title, "[Met]睡眠")
	title = append(title, "[Met]生活習慣の改善意志")
	title = append(title, "[Met]保健指導の希望")
	title = append(title, "[Met]保健指導レベル")
	title = append(title, "[Met]メタボリックシンドローム判定")
	title = append(title, "[Met]医師の診断（特定健診）")
	title = append(title, "初回面接実施")
	title = append(title, "初回面接補足内容")
	title = append(title, "情報提供の方法")
	title = append(title, "カンマ位置(540)")

	return title
}

func requireChk(str string, reqName string) error {
	// 必須項目が空欄か確認する

	if str == "" {
		return fmt.Errorf("必須項目[%s]が空欄です。", reqName)
	} else {
		return nil
	}

}

func kojinIdChk(kojinId string) (string, error) {
	// 個人IDの確認

	if kojinId == "" {
		return "", fmt.Errorf("個人IDに値がありません[%s]", kojinId)
	}

	if kojinId[0:1] != "K" {
		return kojinId, fmt.Errorf("個人IDの先頭が[K]ではありません。[%s]", kojinId)
	}

	return kojinId, nil
}

func waToSeireki(wareki string) (string, error) {
	// 和暦を西暦に変換する

	seireki := ""
	flag := false

	if len(wareki) == 9 {
		wa := wareki[0:1]
		yearNum, err := strconv.Atoi(wareki[1:3])
		logWrite("生年月日エラー", err)
		month := wareki[4:6]
		day := wareki[7:9]

		switch wa {
		case "M":
			yearNum = 1900 + yearNum - 33
		case "T":
			yearNum = 1900 + yearNum + 11
		case "S":
			yearNum = 1900 + yearNum + 25
		case "H":
			yearNum = 1900 + yearNum + 88
		case "R":
			yearNum = 1900 + yearNum + 118
		default:
			flag = true
			break
		}

		seireki = strconv.Itoa(yearNum) + "/" + month + "/" + day

	} else {
		// 和暦ではない
		flag = true
	}

	if flag {
		return wareki, fmt.Errorf("生年月日変換エラー[%s]", wareki)
	} else {
		return seireki, nil
	}

}

func seiConv(sei string) (string, error) {
	// 性別を変換する

	str := ""
	flag := false
	switch sei {
	case "":
		str = ""
	case "男":
		str = "1"
	case "女":
		str = "2"
	default:
		flag = true
	}

	if flag {
		return sei, fmt.Errorf("性別変換エラー[%s]", sei)
	} else {
		return str, nil
	}

}

func jdayConv(jday string) (string, error) {
	// 受診日を変換する

	str := strings.Replace(jday, "-", "/", -1)
	if strings.Contains(str, "/") {
		return str, nil
	} else {
		return jday, fmt.Errorf("受診日変換エラー[%s]", jday)
	}

}

func coursedConv(cd string, name string, age string) (string, string, error) {
	// コースコードとコース名を変換する
	cdFlag := false
	nameFlag := false
	courseFlag := false
	ricohCd := ""
	ricohName := ""
	ageNum, err := strconv.Atoi(age)
	logWrite("年齢エラー", err)

	switch cd {
	case "98009001000001":
		if name == "リコー_人間ドック" {
			courseFlag = true
		} else {
			nameFlag = true
		}

	case "98009001000002":
		if name == "リコー_ミニドック" {
			courseFlag = true
		} else {
			nameFlag = true
		}

	case "98009001000011":
		if name == "リコー_総合Ａ" {
			if ageNum == 35 {
				ricohCd = "31"
				ricohName = "総合健診A(35歳)"
			} else if ageNum == 40 || ageNum == 45 || ageNum == 50 || ageNum == 55 || ageNum == 60 || ageNum == 65 || ageNum == 70 {
				ricohCd = "32"
				ricohName = "総合健診A(節目年齢)"
			} else {
				courseFlag = true
			}
		} else {
			nameFlag = true
		}

	case "98009001000012":
		if name == "リコー_総合Ｂ" {
			if ageNum >= 36 && ageNum%5 != 0 {
				ricohCd = "33"
				ricohName = "総合健診B"
			} else {
				courseFlag = true
			}
		} else {
			nameFlag = true
		}

	case "98009001000013":
		if name == "リコー_事業主Ａ" {
			if ageNum <= 34 {
				ricohCd = "21"
				ricohName = "定期健診(34歳以下)"
			} else {
				courseFlag = true
			}
		} else {
			nameFlag = true
		}

	case "98009001000014":
		if name == "リコー_事業主Ｂ" {
			courseFlag = true
		} else {
			nameFlag = true
		}

	case "98009001000015":
		if name == "リコー_家族健診" {
			courseFlag = true
		} else {
			nameFlag = true
		}

	case "98009001000016":
		if name == "リコー_婦人科" {
			courseFlag = true
		} else {
			nameFlag = true
		}

	case "98009001000017":
		if name == "リコー_基本(ｽﾏｲﾙ)健診" || name == "リコー_基本(ｽﾏｲﾙ）健診" {
			ricohCd = "60"
			ricohName = "スマイル健診"
		} else {
			nameFlag = true
		}

	case "98009001000018":
		if name == "リコー_海外赴任時" {
			if ageNum <= 35 {
				ricohCd = "41"
				ricohName = "海外赴任時(35歳以下)"
			} else {
				ricohCd = "42"
				ricohName = "海外赴任時(36歳以上)"
			}
		} else {
			nameFlag = true
		}

	case "98009001000019":
		if name == "リコー_海外一時帰国" {
			if ageNum <= 34 {
				ricohCd = "45"
				ricohName = "海外一時帰国(34歳以下)"
			} else if ageNum == 35 || ageNum == 40 || ageNum == 45 || ageNum == 50 || ageNum == 55 || ageNum == 60 || ageNum == 65 || ageNum == 70 {
				ricohCd = "46"
				ricohName = "海外一時帰国(節目年齢)"
			} else {
				ricohCd = "47"
				ricohName = "海外一時帰国(節目年齢以外)"
			}
		} else {
			nameFlag = true
		}

	case "98009001000020":
		if name == "リコー_海外完全帰国" {
			ricohCd = "49"
			ricohName = "完全帰国時(全年齢)"
		} else {
			nameFlag = true
		}

	case "98009001000021":
		if name == "リコー_定期健診" {
			if ageNum <= 34 {
				ricohCd = "21"
				ricohName = "定期健診(34歳以下)"
			}
		} else {
			nameFlag = true
		}

	case "98009001000023":
		if name == "リコー_海外赴任時(被扶養配偶者)" {
			ricohCd = "51"
			ricohName = "海外赴任時(全年齢)"
		} else {
			nameFlag = true
		}

	case "98009001000024":
		if name == "リコー_海外一時帰国（被扶養配偶者）" {
			ricohCd = "52"
			ricohName = "海外一時帰国(全年齢)"
		} else {
			nameFlag = true
		}

	case "98009001000025":
		if name == "リコー_海外完全帰国（被扶養配偶者）" {
			ricohCd = "53"
			ricohName = "完全帰国時(全年齢)"
		} else {
			nameFlag = true
		}

	case "04019001000001":
		if name == "リコー定期" {
			ricohCd = "21"
			ricohName = "定期健診(34歳以下)"
		} else {
			nameFlag = true
		}

	case "04019001000002":
		if name == "リコー入社" {
			ricohCd = "11"
			ricohName = "雇入れ時健診"
		} else {
			nameFlag = true
		}


	default:
		cdFlag = true
	}

	if cdFlag {
		return ricohCd, ricohName, fmt.Errorf("コース変換エラー(%s_%s)変換プログラムのコースコードを確認してください。", cd, name)
	} else if nameFlag {
		return ricohCd, ricohName, fmt.Errorf("コース変換エラー(%s_%s)変換プログラムのコース名を確認してください。", cd, name)
	} else if courseFlag {
		return ricohCd, ricohName, fmt.Errorf("コース変換エラー(%s_%s)変換プログラムのコース登録の仕様を確認してください。", cd, name)
	}
	return ricohCd, ricohName, nil
}

func sisetsuConv(sisetsu string) (string, error) {
	// 施設/巡回区分を変換する

	str := ""
	flag := false
	switch sisetsu {
	case "":
		str = ""
	case "所内":
		str = "1"
	case "巡回":
		str = "2"
	default:
		flag = true
	}

	if flag {
		return sisetsu, fmt.Errorf("施設/巡回変換エラー[%s]", sisetsu)
	} else {
		return str, nil
	}
}

func hanteiCdConv(hanteiCd string) (string, string, error) {
	// 判定区分コードを変換する

	cd := ""
	name := ""
	flag := false
	switch hanteiCd {
	case "":
		cd = ""
		name = ""
	case "Ａ":
		cd = "1"
		name = "異常なし"
	case "Ｂ":
		cd = "2"
		name = "軽度異常"
	case "Ｃ":
		cd = "3"
		name = "要経過観察"
	case "Ｄ":
		cd = "5"
		name = "要医療（要精検・要治療）"
	case "Ｅ":
		cd = "5"
		name = "要医療（要精検・要治療）"
	case "Ｆ":
		cd = "5"
		name = "要医療（要精検・要治療）"
	case "Ｇ":
		cd = "7"
		name = "治療中"
	case "Ｈ":
		cd = "9"
		name = "判定不能または再検"
	default:
		flag = true
	}

	if flag {
		return cd, name, fmt.Errorf("判定区分変換エラー[%s]", hanteiCd)
	} else {
		return cd, name, nil
	}

}

func kanaConv(kana string) string {
	// 半角カナを全角にする
	str := kana

	str = strings.Replace(str, "ｶﾞ", "ガ", -1)
	str = strings.Replace(str, "ｷﾞ", "ギ", -1)
	str = strings.Replace(str, "ｸﾞ", "グ", -1)
	str = strings.Replace(str, "ｹﾞ", "ゲ", -1)
	str = strings.Replace(str, "ｺﾞ", "ゴ", -1)

	str = strings.Replace(str, "ｻﾞ", "ザ", -1)
	str = strings.Replace(str, "ｼﾞ", "ジ", -1)
	str = strings.Replace(str, "ｽﾞ", "ズ", -1)
	str = strings.Replace(str, "ｾﾞ", "ゼ", -1)
	str = strings.Replace(str, "ｿﾞ", "ゾ", -1)

	str = strings.Replace(str, "ﾀﾞ", "ダ", -1)
	str = strings.Replace(str, "ﾁﾞ", "ヂ", -1)
	str = strings.Replace(str, "ﾂﾞ", "ヅ", -1)
	str = strings.Replace(str, "ﾃﾞ", "デ", -1)
	str = strings.Replace(str, "ﾄﾞ", "ド", -1)

	str = strings.Replace(str, "ﾊﾞ", "バ", -1)
	str = strings.Replace(str, "ﾋﾞ", "ビ", -1)
	str = strings.Replace(str, "ﾌﾞ", "ブ", -1)
	str = strings.Replace(str, "ﾍﾞ", "ベ", -1)
	str = strings.Replace(str, "ﾎﾞ", "ボ", -1)

	str = strings.Replace(str, "ﾊﾟ", "パ", -1)
	str = strings.Replace(str, "ﾋﾟ", "ピ", -1)
	str = strings.Replace(str, "ﾌﾟ", "プ", -1)
	str = strings.Replace(str, "ﾍﾟ", "ペ", -1)
	str = strings.Replace(str, "ﾎﾟ", "ポ", -1)

	str = strings.Replace(str, "ｳﾞ", "ヴ", -1)

	str = strings.Replace(str, "ｱ", "ア", -1)
	str = strings.Replace(str, "ｲ", "イ", -1)
	str = strings.Replace(str, "ｳ", "ウ", -1)
	str = strings.Replace(str, "ｴ", "エ", -1)
	str = strings.Replace(str, "ｵ", "オ", -1)

	str = strings.Replace(str, "ｧ", "ァ", -1)
	str = strings.Replace(str, "ｨ", "ィ", -1)
	str = strings.Replace(str, "ｩ", "ゥ", -1)
	str = strings.Replace(str, "ｴ", "ェ", -1)
	str = strings.Replace(str, "ｵ", "ォ", -1)

	str = strings.Replace(str, "ｶ", "カ", -1)
	str = strings.Replace(str, "ｷ", "キ", -1)
	str = strings.Replace(str, "ｸ", "ク", -1)
	str = strings.Replace(str, "ｹ", "ケ", -1)
	str = strings.Replace(str, "ｺ", "コ", -1)

	str = strings.Replace(str, "ｻ", "サ", -1)
	str = strings.Replace(str, "ｼ", "シ", -1)
	str = strings.Replace(str, "ｽ", "ス", -1)
	str = strings.Replace(str, "ｾ", "セ", -1)
	str = strings.Replace(str, "ｿ", "ソ", -1)

	str = strings.Replace(str, "ﾀ", "タ", -1)
	str = strings.Replace(str, "ﾁ", "チ", -1)
	str = strings.Replace(str, "ﾂ", "ツ", -1)
	str = strings.Replace(str, "ﾃ", "テ", -1)
	str = strings.Replace(str, "ﾄ", "ト", -1)

	str = strings.Replace(str, "ﾅ", "ナ", -1)
	str = strings.Replace(str, "ﾆ", "ニ", -1)
	str = strings.Replace(str, "ﾇ", "ヌ", -1)
	str = strings.Replace(str, "ﾈ", "ネ", -1)
	str = strings.Replace(str, "ﾉ", "ノ", -1)

	str = strings.Replace(str, "ﾊ", "ハ", -1)
	str = strings.Replace(str, "ﾋ", "ヒ", -1)
	str = strings.Replace(str, "ﾌ", "フ", -1)
	str = strings.Replace(str, "ﾍ", "ヘ", -1)
	str = strings.Replace(str, "ﾎ", "ホ", -1)

	str = strings.Replace(str, "ﾏ", "マ", -1)
	str = strings.Replace(str, "ﾐ", "ミ", -1)
	str = strings.Replace(str, "ﾑ", "ム", -1)
	str = strings.Replace(str, "ﾒ", "メ", -1)
	str = strings.Replace(str, "ﾓ", "モ", -1)

	str = strings.Replace(str, "ﾔ", "ヤ", -1)
	str = strings.Replace(str, "ﾕ", "ユ", -1)
	str = strings.Replace(str, "ﾖ", "ヨ", -1)

	str = strings.Replace(str, "ｬ", "ャ", -1)
	str = strings.Replace(str, "ｭ", "ュ", -1)
	str = strings.Replace(str, "ｮ", "ョ", -1)

	str = strings.Replace(str, "ﾗ", "ラ", -1)
	str = strings.Replace(str, "ﾘ", "リ", -1)
	str = strings.Replace(str, "ﾙ", "ル", -1)
	str = strings.Replace(str, "ﾚ", "レ", -1)
	str = strings.Replace(str, "ﾛ", "ロ", -1)

	str = strings.Replace(str, "ﾜ", "ワ", -1)
	str = strings.Replace(str, "ｦ", "ヲ", -1)
	str = strings.Replace(str, "ﾝ", "ン", -1)

	str = strings.Replace(str, "ｰ", "ー", -1)

	return str

}

func limitStr(str string, limit int) string {
	// 文字列の最大をlimitで指定されたバイト数(shifJISのバイト数)で返す
	// utf8とshiftJISでは文字種によるバイト数が異なる
	// utf8:全角3バイト、半角カナ3バイト
	// shift-Jis:全角2バイト、半角カナ1バイト
	// 半角カナを全角カナに変換することで文字種による違いをなくし
	// 全角文字3バイトを2バイトとして文字数を計算する

	utf8str := kanaConv(str)   // 半角カタカナ -> カタカナを全角
	runeStr := []rune(utf8str) // 文字列を一文字ずつのSliceにする

	countLen := 0
	countStr := ""
	for _, v := range runeStr {
		s := string(v)
		if len(s) == 3 { //1文字が3バイトなら2バイトとして計算
			countLen = countLen + 2
		} else {
			countLen = countLen + len(s)
		}

		if countLen > limit {
			break
		} else {
			countStr = countStr + s
		}
	}

	return countStr

}

func joinStr(str1 string, str2 string) string {
	// 2つの文字列を結合する

	str := ""

	if str1 == "" {
		str = str2
	} else {
		if str2 == "" {
			str = str1
		} else {
			str = str1 + " " + str2
		}
	}

	return str
}

func joinStr3(str1 string, str2 string, str3 string) string {
	// 3つの文字列を結合する

	str := ""
	str = joinStr(str1, str2)
	str = joinStr(str, str3)

	return str

}

func joinStr4(str1 string, str2 string, str3 string, str4 string) string {
	// 5つの文字列を結合する

	str := ""
	str = joinStr(str1, str2)
	str = joinStr(str, str3)
	str = joinStr(str, str4)

	return str

}

func joinStr5(str1 string, str2 string, str3 string, str4 string, str5 string) string {
	// 5つの文字列を結合する

	str := ""
	str = joinStr(str1, str2)
	str = joinStr(str, str3)
	str = joinStr(str, str4)
	str = joinStr(str, str5)

	return str

}

func joinStr7(str1 string, str2 string, str3 string, str4 string, str5 string, str6 string, str7 string) string {
	// 7つの文字列を結合する

	str := ""
	str = joinStr(str1, str2)
	str = joinStr(str, str3)
	str = joinStr(str, str4)
	str = joinStr(str, str5)
	str = joinStr(str, str6)
	str = joinStr(str, str7)

	return str

}

func joinStr10(str1 string, str2 string, str3 string, str4 string, str5 string, str6 string, str7 string, str8 string, str9 string, str10 string) string {
	// 10個の文字列を結合する

	str := ""
	str = joinStr(str1, str2)
	str = joinStr(str, str3)
	str = joinStr(str, str4)
	str = joinStr(str, str5)
	str = joinStr(str, str6)
	str = joinStr(str, str7)
	str = joinStr(str, str8)
	str = joinStr(str, str9)
	str = joinStr(str, str10)

	return str

}

func chiryoChk(tenki string) bool {
	// 治療中の項目があればTrueを返す

	flag := false
	if tenki == "内服治療中" || tenki == "管理中" {
		flag = true
	}

	return flag

}

func tenkiConv(kiou []string, tenki []string) (string, string) {
	// 既往歴を確認。治療中と既往があればそれぞれ"1"を返す

	chiryoStr := ""
	kiouStr := ""
	for i, v := range kiou {
		if v != "" {
			if chiryoChk(tenki[i]) {
				chiryoStr = "1"
			} else {
				kiouStr = "1"
			}
		}

		if chiryoStr == "1" && kiouStr == "1" {
			break
		}
	}

	return chiryoStr, kiouStr

}

func kiouConv(kiou []string, tenki []string) (string, string) {
	//既往歴から、治療中の病名と既往の病名に分けて返す

	chiryoStr := ""
	kiouStr := ""
	for i, v := range kiou {
		if v != "" {
			if chiryoChk(tenki[i]) {
				chiryoStr = joinStr(chiryoStr, v)
			} else {
				kiouStr = joinStr(kiouStr, v)
			}
		}
	}

	return chiryoStr, kiouStr
}

func kiouJoin(kiou string, age string, tenki string) string {
	// 病名、年齢、転帰をつなげて返す

	str := ""
	if kiou == "" {
		return ""
	} else {
		str = kiou
	}

	if age != "" {
		str = str + " " + age + "才"
	}

	if tenki != "" {
		str = str + " " + tenki
	}

	return str
}

func umuConv(str string) string {
	// 値をみて有無(1:特記すべきことあり 2:特記すべきことなし)を返す

	if str != "" {
		return "1"
	} else {
		return "2"
	}
}

func eyeConv(eye string) (string, string) {
	// 視力の値とデータ属性(1:未満)を返す

	value := eye
	code := ""
	if strings.Index(eye, "↓") != -1 {
		value = strings.Replace(eye, "↓", "", -1)
		code = "1"
	}

	return value, code
}

func eyeKubun(kyoseiR string, kyoseiL string, kintenR string, kintenL string) string {
	// 視力矯正区分（1:眼鏡 2:CL 3:不明）を返す

	if kyoseiR == "" && kyoseiL == "" && kintenR == "" && kintenL == "" {
		return ""
	} else {
		return "3"
	}
}

func yusyokenChk(hantei string) (bool, error) {
	flag := false
	errFlag := false
	switch hantei {
	case "C", "D", "E", "F", "G", "Ｃ", "Ｄ", "Ｅ", "Ｆ", "Ｇ":
		flag = true
	case "", "A", "B", "Ａ", "Ｂ":
		flag = false
	default:
		errFlag = true
	}

	if errFlag {
		return flag, fmt.Errorf("有所見判定変換エラー[%s]", hantei)
	} else {
		return flag, nil
	}
}

func yusyokenKubun(flag bool) string {
	// 有所見チェック yusyokenChk() の戻り値から、有所見区分(1:所見なし 2:所見あり)を返す

	if flag {
		return "2"
	} else {
		return "1"
	}

}

func ear1kHantei(hantei string) (string, error) {
	// 聴力1000Hzの所見区分(1:所見なし 2:所見あり)を返す

	if hantei == "" {
		return "", nil
	}

	if flag, err := yusyokenChk(hantei); err != nil {
		return hantei, fmt.Errorf("聴力1000Hz判定変換エラー[%s]", hantei)
	} else {
		return yusyokenKubun(flag), nil
	}

}

func ear4kHantei(hantei1 string, hantei2 string) (string, error) {
	// 聴力4000Hzの所見区分(1:所見なし 2:所見あり)を返す

	if hantei1 != "" {
		if flag, err := yusyokenChk(hantei1); err != nil {
			return hantei1, fmt.Errorf("聴力4000Hz判定変換エラー[%s]", hantei1)
		} else {
			return yusyokenKubun(flag), nil
		}
	}

	if hantei2 != "" {
		if flag, err := yusyokenChk(hantei2); err != nil {
			return hantei2, fmt.Errorf("聴力4000Hz判定変換エラー[%s]", hantei2)
		} else {
			return yusyokenKubun(flag), nil
		}
	}

	return "", nil

}

func earKaiwa(syoken string, hantei string) (string, error) {
	// 聴力会話法の所見区分(1:所見なし 2:所見あり)を返す

	if syoken == "" {
		return "", nil
	}

	if flag, err := yusyokenChk(hantei); err != nil {
		return hantei, fmt.Errorf("聴力会話法判定変換エラー[%s]", hantei)
	} else {
		return yusyokenKubun(flag), nil
	}

}

func earConv(ear string) string {
	// 聴力の値が * なら空欄にする

	if ear == "*" || ear == "＊" {
		return ""
	} else {
		return ear
	}

}

func hanteiRank(hantei string) (int, error) {
	// ABC判定の重さ（ランキング）を返す
	rank := 0
	flag := false
	switch hantei {
	case "":
		rank = 0
	case "A", "Ａ":
		rank = 1
	case "B", "Ｂ":
		rank = 2
	case "C", "Ｃ":
		rank = 3
	case "D", "Ｄ":
		rank = 5
	case "E", "Ｅ":
		rank = 6
	case "F", "Ｆ":
		rank = 7
	case "G", "Ｇ":
		rank = 4
	default:
		flag = true
	}

	if flag {
		return rank, fmt.Errorf("判定ランク変換エラー[%s]", hantei)
	} else {
		return rank, nil
	}

}

func hantiHeavy(hantei1 string, hantei2 string) (string, error) {
	// 重い方の判定を返す

	var err error
	rank1 := 0
	rank2 := 0

	rank1, err = hanteiRank(hantei1)
	if err != nil {
		return hantei1, err
	}

	rank2, err = hanteiRank(hantei2)
	if err != nil {
		return hantei2, err
	}

	if rank1 < rank2 {
		return hantei2, nil
	} else {
		return hantei1, nil
	}

}

func ketsuatuTimes(hantei1H string, hantei1L string, hantei2H string, hantei2L string) (int, error) {
	// 血圧1回目と血圧2回目をどちらを報告値とするか決める(1:1回目, 2:2回目)

	if hantei2H == "" && hantei2L == "" {
		return 1, nil
	}

	rank1H, err := hanteiRank(hantei1H)
	if err != nil {
		return 1, err
	}

	rank1L, err := hanteiRank(hantei1L)
	if err != nil {
		return 1, err
	}

	rank2H, err := hanteiRank(hantei2H)
	if err != nil {
		return 1, err
	}

	rank2L, err := hanteiRank(hantei2L)
	if err != nil {
		return 1, err
	}

	count1 := 0
	count2 := 0

	if rank1H != rank2H {
		if rank1H > rank2H {
			count1++
		} else {
			count2++
		}
	}

	if rank1L != rank2L {
		if rank1L > rank2L {
			count1++
		} else {
			count2++
		}
	}

	if count1 >= count2 {
		return 2, nil
	} else {
		return 1, nil
	}

}

func syokenUmu(hantei string) (string, error) {
	// 所見有無(1:異常所見あり 2:所見なし)を返す

	umu := ""
	flag := false
	switch hantei {
	case "":
		umu = ""
	case "A", "Ａ":
		umu = "2"
	case "B", "Ｂ":
		umu = "2"
	case "C", "Ｃ":
		umu = "1"
	case "D", "Ｄ":
		umu = "1"
	case "E", "Ｅ":
		umu = "1"
	case "F", "Ｆ":
		umu = "1"
	case "G", "Ｇ":
		umu = "1"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("所見有無変換エラー[%s]", hantei)
	} else {
		return umu, nil
	}

}

func taisyo(hantei string) string {
	// 実施対象(hanteiに値がある)なら"0"を返す

	if hantei == "" {
		return ""
	} else {
		return "0"
	}

}

func satsuei(kan string, cyoku string) string {
	// 撮影方法(1:直接 2:間接 3:デジタル)を返す

	if kan == "" && cyoku == "" {
		return ""
	}

	if cyoku != "" {
		return "3"
	}

	if kan != "" {
		return "2"
	}

	return ""

}

func kakutanConv(kakutan string) (string, string, string, error) {
	// 喀痰細胞診結果(1:A取り直し 2:B異常なし 3:C要注意 4:Dがんの疑い 5:Fがん)を返す

	cd := ""
	name := ""
	str := ""
	flag := false
	switch kakutan {
	case "":
		cd = ""
		name = ""
		str = ""
	case "Ⅰ":
		cd = "1"
		name = "異常なし"
		str = "2"
	case "Ⅱ":
		cd = "2"
		name = "軽度異常"
		str = "2"
	case "Ⅲ":
		cd = "5"
		name = "要医療（要精検・要治療）"
		str = "3"
	case "Ⅳ":
		cd = "5"
		name = "要医療（要精検・要治療）"
		str = "4"
	case "Ⅴ":
		cd = "5"
		name = "要医療（要精検・要治療）"
		str = "5"
	case "判定不能":
		cd = "9"
		name = "判定不能または再検"
		str = "1"
	default:
		flag = true
	}

	if flag {
		return "", "", "", fmt.Errorf("喀痰変換エラー[%s]", kakutan)
	} else {
		return cd, name, str, nil
	}

}

func scheieConv(S string, H string) string {
	// シェイエ分類を返す

	if S == "" && H == "" {
		return ""
	}

	return "S:" + width.Narrow.String(S) + "/H:" + width.Narrow.String(H)
}

func scottConv(scott string) (string, error) {
	// scott分類を変換する

	str := ""
	flag := false
	switch scott {
	case "":
		str = ""
	case "０":
		str = ""
	case "Ⅰ":
		str = "Ⅰ(a)"
	case "Ⅰａ":
		str = "Ⅰ(a)"
	case "Ⅰｂ":
		str = "Ⅰ(b)"
	case "Ⅱ":
		str = "Ⅱ"
	case "Ⅲ":
		str = "Ⅲ(a)"
	case "Ⅲａ":
		str = "Ⅲ(a)"
	case "Ⅲｂ":
		str = "Ⅲ(b)"
	case "Ⅳ":
		str = "Ⅳ"
	case "Ⅴ":
		str = "Ⅴ(a)"
	case "Ⅴａ":
		str = "Ⅴ(a)"
	case "Ⅴｂ":
		str = "Ⅴ(b)"
	case "Ⅵ":
		str = "Ⅵ"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("Scott分類変換エラー[%s]", scott)
	} else {
		return str, nil
	}

}

func kwConv(kw string) (string, error) {
	// KW分類を変換する

	str := ""
	flag := false
	switch kw {
	case "":
		str = ""
	case "０":
		str = "０"
	case "Ⅰ":
		str = "Ⅰ"
	case "Ⅱ":
		str = "Ⅱ"
	case "Ⅱａ":
		str = "Ⅱ(a)"
	case "Ⅱｂ":
		str = "Ⅱ(b)"
	case "Ⅲ":
		str = "Ⅲ"
	case "Ⅳ":
		str = "Ⅳ"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("KW分類変換エラー[%s]", kw)
	} else {
		return str, nil
	}

}

func teiseiConv(teisei string) (string, error) {
	// 定性(1:(-), 2:(+-), 3:(+), 4:(2+), 5:(3+), 6:(4+), 7:(5+) )を変換する

	str := ""
	flag := false
	switch teisei {
	case "":
		str = ""
	case "-", "－":
		str = "1"
	case "+-":
		str = "2"
	case "+", "＋":
		str = "3"
	case "2+":
		str = "4"
	case "3+":
		str = "5"
	case "4+":
		str = "6"
	case "5+":
		str = "7"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("定性変換エラー[%s]", teisei)
	} else {
		return str, nil
	}

}

func nyoChinsaConv(sonota1 string, sonota2 string, sonota3 string) (string, string) {
	// 尿沈渣その他１・２・３　より　細菌を分けて返す

	if sonota1 == "" {
		return "", ""
	}

	saikin := ""
	if saikinChk(sonota1) {
		saikin = "細菌"
		sonota1 = ""
	} else if saikinChk(sonota2) {
		saikin = "細菌"
		sonota2 = ""
	} else if saikinChk(sonota3) {
		saikin = "細菌"
		sonota3 = ""
	}

	return saikin, joinStr3(sonota1, sonota2, sonota3)

}

func saikinChk(str string) bool {
	// 細菌なら trueを返す

	switch str {
	case "細菌", "ｻｲｷﾝ":
		return true
	default:
		return false
	}

}

func aboConv(abo string) (string, error) {
	// 血液型ABO(1:A 2:B 3:AB 4:O)を変換する

	str := ""
	flag := false
	switch abo {
	case "":
		str = ""
	case "Ａ型":
		str = "1"
	case "Ｂ型":
		str = "2"
	case "ＡＢ型":
		str = "3"
	case "Ｏ型":
		str = "4"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("血液型ABO変換エラー[%s]", abo)
	} else {
		return str, nil
	}

}

func rhConv(rh string) (string, error) {
	// 血液型Rh(1:+ 2:-)を変換する

	str := ""
	flag := false
	switch rh {
	case "":
		str = ""
	case "（＋）":
		str = "1"
	case "（－）":
		str = "2"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("血液型Rh変換エラー[%s]", rh)
	} else {
		return str, nil
	}

}

func numChk(str string) (string, error) {
	// 数値に変換できるか確認し、余計な文字列を削除する

	if str == "" {
		return "", nil
	}

	_, errInt := strconv.Atoi(str)
	_, errFloat := strconv.ParseFloat(str, 64)

	if errInt != nil && errFloat != nil {
		switch true {
		case strings.Index(str, "未満") >= 0:
			return strings.Replace(str, "未満", "", -1), nil
		case strings.Index(str, "以上") >= 0:
			return strings.Replace(str, "以上", "", -1), nil
		case strings.Index(str, "以下") >= 0:
			return strings.Replace(str, "以下", "", -1), nil
		default:
			return str, fmt.Errorf("値に文字が含まれています。[%s]", str)
		}
	} else {
		return str, nil
	}

}

func eatTimeConv(toh string, eatTime string) (string, error) {
	// 食後時間区分(2:10時間以上 3:食後3.5時間以上10時間未満 4: 食後3.5時間未満
	// 空腹時血糖結果値がある場合は「2」を登録

	if eatTime == "" && toh == "" {
		return "", nil
	}

	if eatTime == "" && toh != "" {
		return "2", nil //食後時間が空欄で血糖値があるなら 2:10時間以上
	}

	eatTimeFloat, err := strconv.ParseFloat(eatTime, 64)
	if err == nil {
		switch {
		case eatTimeFloat < 3.5:
			return "4", nil
		case eatTimeFloat >= 3.5 && eatTimeFloat < 10.0:
			return "3", nil
		case eatTimeFloat >= 10.0:
			return "2", nil
		default:
			return "", fmt.Errorf("食後時間の値が不正です。[%s]", eatTime)
		}
	} else {
		return eatTime, err
	}

}

func seiriConv(seiri string) (string, error) {
	// 生理区分(1: 生理中)を返す

	str := ""
	flag := false
	switch seiri {
	case "":
		str = ""
	case "はい":
		str = "1"
	case "いいえ":
		str = ""
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("生理区分変換エラー[%s]", seiri)
	} else {
		return str, nil
	}

}

func ninshinConv(ninshin string, utagai string) (string, error) {
	// 妊娠区分(1: 妊娠中)を返す

	if ninshin == "" && utagai == "" {
		return "", nil
	}

	if ninshin == "はい" || utagai == "はい" {
		return "1", nil
	}

	if ninshin == "いいえ" || utagai == "いいえ" {
		return "", nil
	}

	if ninshin != "" {
		return "", fmt.Errorf("妊娠区分変換エラー[%s]", ninshin)
	}

	if utagai != "" {
		return "", fmt.Errorf("妊娠区分変換エラー[%s]", utagai)
	}

	return "", fmt.Errorf("妊娠区分変換エラー[%s][%s]", ninshin, utagai)
}

func nyubiConv(cmt1 string, cmt2 string) string {
	// 乳び(1:1+ 2:2+ 3:3+)を返す

	if cmt1 == "乳び血清" || cmt2 == "乳び血清" {
		return "1"
	}

	if cmt1 == "乳糜検体" || cmt2 == "乳糜検体" {
		return "1"
	}

	return ""

}

func yoketsuConv(cmt1 string, cmt2 string) string {
	// 溶血(1:1+ 2:2+ 3:3+)を返す

	if cmt1 == "強溶血血清" || cmt2 == "強溶血血清" {
		return "2"
	}

	if cmt1 == "溶血血清" || cmt2 == "溶血血清" {
		return "1"
	}

	if cmt1 == "溶血検体" || cmt2 == "溶血溶血" {
		return "1"
	}

	return ""

}

func tohConv(toh string, eatTime string) (string, string) {
	// 空腹時血糖と随時血糖を返す
	// 空腹時血糖は食後10時間以上であること

	if toh == "" {
		return "", ""
	}

	switch eatTime {
	case "2": // 10時間以上
		return toh, ""
	case "3": // 食後3.5時間以上10時間未満
		return "", toh
	case "4": // 食後3.5時間未満
		return "", toh // 随時血糖ではないが一応値を返す
	default:
		return toh, "" // 空腹時血糖ではないが一応値を返す
	}

}

func iabcConv(iabc string) (string, error) {
	// 胃ABC検診判定分類を(1:A群 2:B群 3:C群 4:D群 5:判定不能 6:E群)返す

	str := ""
	flag := false
	switch iabc {
	case "":
		str = ""
	case "A群":
		str = "1"
	case "B群":
		str = "2"
	case "C群":
		str = "3"
	case "D群":
		str = "4"
	case "E群":
		str = "6"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("胃ABC分類変換エラー[%s]", iabc)
	} else {
		return str, nil
	}

}

func mmgSatsuei(one string, two string) string {
	// マンモ撮影方法(1:1方向 2:2方向)を返す

	if one == "" && two == "" {
		return ""
	}

	if one != "" {
		return "1"
	}

	if two != "" {
		return "2"
	}

	return ""

}

func vesesudaConv(vesesuda string) (string, error) {
	// ベセスダ分類を変換する

	str := ""
	flag := false
	switch vesesuda {
	case "":
		str = ""
	case "NILM":
		str = "1"
	case "ASC-US":
		str = "2"
	case "ASC-H":
		str = "3"
	case "LSIL":
		str = "4"
	case "HSIL":
		str = "5"
	case "SCC":
		str = "6"
	case "AGC":
		str = "7"
	case "AIS":
		str = "8"
	case "Adeno.ca":
		str = "9"
	case "Other":
		str = "10"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("ベセスダ分類変換エラー[%s]", vesesuda)
	} else {
		return str, nil
	}

}

func nichimoConv(nichimo string) (string, error) {
	// 日母分類を変換する

	str := ""
	flag := false
	switch nichimo {
	case "":
		str = ""
	case "Ⅰ":
		str = "1"
	case "Ⅱ":
		str = "2"
	case "Ⅲａ":
		str = "3"
	case "Ⅲｂ":
		str = "4"
	case "Ⅳ":
		str = "5"
	case "Ⅴ":
		str = "6"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("日母クラス分類変換エラー[%s]", nichimo)
	} else {
		return str, nil
	}

}

func jikakuUmu(jikaku string) string {
	// 自覚症状有無(1:特記すべきことあり 2: 特記すべきことなし)を返す

	if jikaku == "" || jikaku == "特になし" {
		return "2"
	} else {
		return "1"
	}

}

func takakuUmu(takaku string) string {
	// 他覚症状有無(1:特記すべきことあり 2: 特記すべきことなし)を返す

	if takaku == "" || takaku == "異常なし" {
		return "2"
	} else {
		return "1"
	}

}

func yesNoConv(yesNo string) (string, error) {
	// はい・いいえ(1:はい 2:いいえ)を変換する

	str := ""
	flag := false
	switch yesNo {
	case "":
		str = ""
	case "はい":
		str = "1"
	case "いいえ":
		str = "2"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("はい・いいえ変換エラー[%s]", yesNo)
	} else {
		return str, nil
	}

}

func sosyakuConv(sosyaku string) (string, error) {
	// 咀嚼(1:何でも 2:かみにくい 3:ほどんどかめない)を変換する

	str := ""
	flag := false
	switch sosyaku {
	case "":
		str = ""
	case "何でも":
		str = "1"
	case "かみにくい":
		str = "2"
	case "ほとんどかめない":
		str = "3"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("咀嚼変換エラー[%s]", sosyaku)
	} else {
		return str, nil
	}

}

func eat1Conv(eat string) (string, error) {
	// 食べ方1早食い等(1:速い 2:ふつう 3:遅い)を変換する

	str := ""
	flag := false
	switch eat {
	case "":
		str = ""
	case "速い":
		str = "1"
	case "普通":
		str = "2"
	case "遅い":
		str = "3"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("食べ方1(早食い等)変換エラー[%s]", eat)
	} else {
		return str, nil
	}

}

func eat3Conv(eat string) (string, error) {
	// 食べ方3間食(1:毎日 2:時々 3:ほとんど摂取しない)を変換する

	str := ""
	flag := false
	switch eat {
	case "":
		str = ""
	case "毎日":
		str = "1"
	case "時々":
		str = "2"
	case "ほとんど摂取しない":
		str = "3"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("食べ方3(間食)変換エラー[%s]", eat)
	} else {
		return str, nil
	}

}

func sakeConv(sake string) (string, error) {
	// 飲酒習慣(1:毎日 2:時々 3:ほとんど飲まない)を変換する

	str := ""
	flag := false
	switch sake {
	case "":
		str = ""
	case "毎日":
		str = "1"
	case "時々":
		str = "2"
	case "飲まない":
		str = "3"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("飲酒習慣変換エラー[%s]", sake)
	} else {
		return str, nil
	}

}

func sakeryoConv(sakeryo string) (string, error) {
	// 飲酒量(1:1号未満 2:1～2合未満 3:2～3合 4:3合以上)を変換する

	str := ""
	flag := false
	switch sakeryo {
	case "":
		str = ""
	case "１合未満":
		str = "1"
	case "１～２合未満":
		str = "2"
	case "２～３合未満":
		str = "3"
	case "３合以上":
		str = "4"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("飲酒量変換エラー[%s]", sakeryo)
	} else {
		return str, nil
	}

}

func seikatsuConv(seikatsu string) (string, error) {
	// 生活習慣の改善意志(1:意思なし 2:意志あり(6カ月いない) 3:意志あり(近いうち) 4:取組済み(6カ月未満) 5:取組済み(6カ月以上))を変換する

	str := ""
	flag := false
	switch seikatsu {
	case "":
		str = ""
	case "しない":
		str = "1"
	case "思う":
		str = "2"
	case "始めた":
		str = "3"
	case "６ヶ月経過":
		str = "4"
	case "６ヶ月以上":
		str = "5"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("生活習慣の改善意志変換エラー[%s]", seikatsu)
	} else {
		return str, nil
	}

}

func hokenConv(hoken string) (string, error) {
	// 保健指導レベル(1:積極的支援 2:動機付け支援 3:なし 4:判定不能)を変換する

	str := ""
	flag := false
	switch hoken {
	case "":
		str = ""
	case "積極的支援レベル":
		str = "1"
	case "動機づけ支援レベル":
		str = "2"
	case "情報提供レベル":
		str = "3"
	case "判定不能":
		str = "4"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("保健指導レベル変換エラー[%s]", hoken)
	} else {
		return str, nil
	}

}

func metaboConv(metabo string) (string, error) {
	// メタボリックシンドローム判(1:基準該当 2:予備軍該当 3:非該当 4:判定不能)を変換する

	str := ""
	flag := false
	switch metabo {
	case "":
		str = ""
	case "基準該当":
		str = "1"
	case "予備群該当":
		str = "2"
	case "非該当":
		str = "3"
	case "判定不能":
		str = "4"
	default:
		flag = true
	}

	if flag {
		return "", fmt.Errorf("メタボリックシンドローム判定変換エラー[%s]", metabo)
	} else {
		return str, nil
	}

}

func sogoConv(sogoStr [][2]string) (string, error) {

	str := ""
	strA := ""
	strB := ""
	strC := ""
	strD := ""
	strE := ""
	strF := ""
	strG := ""
	strH := ""
	flag := false
	for _, v := range sogoStr {
		switch v[0] {
		case "":

		case "Ａ":
			strA = joinStr(strA, v[1])
		case "Ｂ":
			strB = joinStr(strB, v[1])
		case "Ｃ":
			strC = joinStr(strC, v[1])
		case "Ｄ":
			strD = joinStr(strD, v[1])
		case "Ｅ":
			strE = joinStr(strE, v[1])
		case "Ｆ":
			strF = joinStr(strF, v[1])
		case "Ｇ":
			strG = joinStr(strG, v[1])
		case "Ｈ":
			strH = joinStr(strH, v[1])
		default:
			flag = true
			str = v[1]
			break
		}
	}

	if flag {
		return "", fmt.Errorf("総合判定変換エラー[%s]", str)
	} else {
		str = joinStr(str, strF)
		str = joinStr(str, strE)
		str = joinStr(str, strD)
		str = joinStr(str, strG)
		str = joinStr(str, strH)
		str = joinStr(str, strC)
		str = joinStr(str, strB)
		str = joinStr(str, strA)
		return str, nil
	}

}

func PSAconv(PSA string) (string, error) {
	// PSAの陰・陽区分(1:-(陰性) 3:+(陽性))を返す

	if PSA == "" {
		return "", nil
	}

	PSAfloat, err := strconv.ParseFloat(PSA, 64)
	if err != nil {
		return PSA, fmt.Errorf("PSA数値変換エラー[%s]", PSA)
	}

	if PSAfloat > 4.00 {
		return "3", nil
	}

	return "1", nil

}

func CA125conv(CA125 string) (string, error) {
	// CA125の陰・陽区分(1:-(陰性) 3:+(陽性))を返す

	if CA125 == "" {
		return "", nil
	}

	CA125float, err := strconv.ParseFloat(CA125, 64)
	if err != nil {
		return CA125, fmt.Errorf("CA125数値変換エラー[%s]", CA125)
	}

	if CA125float > 35.0 {
		return "3", nil
	}

	return "1", nil

}

func CA199conv(CA199 string) (string, error) {
	// CA19-9の陰・陽区分(1:-(陰性) 3:+(陽性))を返す

	if CA199 == "" {
		return "", nil
	}

	CA199float, err := strconv.ParseFloat(CA199, 64)
	if err != nil {
		return CA199, fmt.Errorf("CA199数値変換エラー[%s]", CA199)
	}

	if CA199float > 37.0 {
		return "3", nil
	}

	return "1", nil

}

func CEAconv(CEA string) (string, error) {
	// CEAの陰・陽区分(1:-(陰性) 3:+(陽性))を返す

	if CEA == "" {
		return "", nil
	}

	CEAfloat, err := strconv.ParseFloat(CEA, 64)
	if err != nil {
		return CEA, fmt.Errorf("CEA数値変換エラー[%s]", CEA)
	}

	if CEAfloat > 5.0 {
		return "3", nil
	}

	return "1", nil

}

func AFPconv(AFP string) (string, error) {
	// AFPの陰・陽区分(1:-(陰性) 3:+(陽性))を返す

	if AFP == "" {
		return "", nil
	}

	AFPfloat, err := strconv.ParseFloat(AFP, 64)
	if err != nil {
		return AFP, fmt.Errorf("AFP数値変換エラー[%s]", AFP)
	}

	if AFPfloat > 10.0 {
		return "3", nil
	}

	return "1", nil

}

func SifuraConv(Sifura string) (string, error) {
	// Sifuraの陰・陽区分(1:-(陰性) 3:+(陽性))を返す

	if Sifura == "" {
		return "", nil
	}

	SifuraFloat, err := strconv.ParseFloat(Sifura, 64)
	if err != nil {
		return Sifura, fmt.Errorf("Sifura数値変換エラー[%s]", Sifura)
	}

	if SifuraFloat > 3.5 {
		return "3", nil
	}

	return "1", nil

}
