package primitive

const SIZE = 2048
// [][][0] - R, [][][1] - G, [][][2] - B
var CurrentMemo [SIZE][SIZE][3]int64
var TargetMemo [SIZE][SIZE][3]int64
var CurrentSquaredMemo [SIZE][SIZE][3]int64
var CurrentTargetMemo [SIZE][SIZE][3]int64


func CreateCurrentMemoizations(model *Model) {
	idx := 0
	w := model.Current.Bounds().Size().X
	h := model.Current.Bounds().Size().Y
	for i:=0; i<w; i++ {
		for j:=0; j<h; j++ {
			R := model.Current.Pix[idx]
			G := model.Current.Pix[idx+1]
			B := model.Current.Pix[idx+2]

			CurrentMemo[i][j][0] = int64(R)
			CurrentMemo[i][j][1] = int64(G)
			CurrentMemo[i][j][2] = int64(B)

			CurrentSquaredMemo[i][j][0] = int64(R)*int64(R)
			CurrentSquaredMemo[i][j][1] = int64(G)*int64(G)
			CurrentSquaredMemo[i][j][2] = int64(B)*int64(B)

			R1 := TargetMemo[i][j][0]
			G1 := TargetMemo[i][j][1]
			B1 := TargetMemo[i][j][2]

			if j>0 {
				R1 -= TargetMemo[i][j-1][0]
				G1 -= TargetMemo[i][j-1][1]
				B1 -= TargetMemo[i][j-1][2]
			}

			CurrentTargetMemo[i][j][0] = int64(R) * R1
			CurrentTargetMemo[i][j][1] = int64(G) * G1
			CurrentTargetMemo[i][j][2] = int64(B) * B1

			if j > 0 {
				CurrentMemo[i][j][0] += CurrentMemo[i][j-1][0]
				CurrentMemo[i][j][1] += CurrentMemo[i][j-1][1]
				CurrentMemo[i][j][2] += CurrentMemo[i][j-1][2]

				CurrentSquaredMemo[i][j][0] += CurrentSquaredMemo[i][j-1][0]
				CurrentSquaredMemo[i][j][1] += CurrentSquaredMemo[i][j-1][1]
				CurrentSquaredMemo[i][j][2] += CurrentSquaredMemo[i][j-1][2]

				CurrentTargetMemo[i][j][0] += CurrentTargetMemo[i][j-1][0]
				CurrentTargetMemo[i][j][1] += CurrentTargetMemo[i][j-1][1]
				CurrentTargetMemo[i][j][2] += CurrentTargetMemo[i][j-1][2]
			}
			idx += 4
		}
	}
}

func CreateTargetMemoizations(model *Model) {
	idx := 0
	w := model.Target.Bounds().Size().X
	h := model.Target.Bounds().Size().Y
	for i:=0; i<w; i++ {
		for j:=0; j<h; j++ {
			R := model.Target.Pix[idx]
			G := model.Target.Pix[idx+1]
			B := model.Target.Pix[idx+2]
			TargetMemo[i][j][0] = int64(R)
			TargetMemo[i][j][1] = int64(G)
			TargetMemo[i][j][2] = int64(B)
			if j > 0 {
				TargetMemo[i][j][0] += TargetMemo[i][j-1][0]
				TargetMemo[i][j][1] += TargetMemo[i][j-1][1]
				TargetMemo[i][j][2] += TargetMemo[i][j-1][2]
			}
			idx += 4
		}
	}
}