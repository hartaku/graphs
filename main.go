..package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func findDotExe() (string, error) {
	var dotExePath string
	root := "C:\\"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil // Пропустить ошибки доступа
			}
			return err
		}
		if info.Name() == "dot.exe" {
			dotExePath = path
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if dotExePath == "" {
		return "", fmt.Errorf("файл dot.exe не найден")
	}

	return strings.Replace(dotExePath, `\dot.exe`, "", 1), nil
}

// Функция для добавления пути к dot.exe в переменную окружения PATH
func addToPath(path string) error {
	cmd := exec.Command("setx", "PATH", fmt.Sprintf("%s;%s", path, os.Getenv("PATH")))
	return cmd.Run()
}

func isVerticesConnected(adjMatrix [][]float64, vertex1 int, vertex2 int) bool {
	if vertex1 < 0 || vertex1 >= len(adjMatrix) || vertex2 < 0 || vertex2 >= len(adjMatrix) {
		return false
	}

	visited := make([]bool, len(adjMatrix))

	// Вызываем dfs, начиная с vertex1
	if dfstarget(adjMatrix, visited, vertex1, vertex2) {
		return true
	}

	// Если не удалось найти связь из vertex1, сбрасываем состояние visited и повторяем поиск с vertex2
	visited = make([]bool, len(adjMatrix))

	if dfstarget(adjMatrix, visited, vertex2, vertex1) {
		return true
	}

	return false
}

func dfstarget(adjMatrix [][]float64, visited []bool, current int, target int) bool {
	visited[current] = true

	if current == target {
		return true
	}

	for next := 0; next < len(adjMatrix); next++ {
		if adjMatrix[current][next] > 0 && !visited[next] {
			if dfstarget(adjMatrix, visited, next, target) {
				return true
			}
		}
	}

	return false
}

func printPath(paths []float64, current float64) []float64 {
	var path []float64

	if paths[int(current)] == -1 {
		path = append(path, current)
	} else {
		path = printPath(paths, paths[int(current)])
		path = append(path, current)
	}

	return path
}

func getPath(paths []float64, current float64) [][]float64 {
	path := printPath(paths, current)
	matrixSize := len(paths)
	matrixPath := make([][]float64, matrixSize)
	for i := 0; i < matrixSize; i++ {
		matrixPath[i] = make([]float64, matrixSize)
	}
	for k := 0; k < len(path)-1; k++ {
		matrixPath[int(path[k])][int(path[k+1])] = 1
		matrixPath[int(path[k+1])][int(path[k])] = 1
	}
	return matrixPath
}

func IsGraphConnected(graph [][]float64) bool {
	visited := make([]bool, len(graph))
	dfs(0, graph, visited)
	for _, v := range visited {
		if !v {
			return false
		}
	}
	return true
}

func dfs(node int, graph [][]float64, visited []bool) {
	visited[node] = true
	for i := 0; i < len(graph[node]); i++ {
		if graph[node][i] > 0 && !visited[i] {
			dfs(i, graph, visited)
		}
	}
}

func MinLenthCom(a [][]float64) [][]float64 {
	for k, v := range a {
		for k1, _ := range v {
			if a[k][k1] != 0 {
				a[k1][k] = a[k][k1]
			}

		}
	}
	size := len(a)
	visited := map[int]bool{0: true}
	visitedar := []int{0}
	res := make([][]float64, len(a))
	for i := range res {
		res[i] = make([]float64, len(a))
	}
	for len(visitedar) < size {

		min := math.MaxFloat64
		var i, j int
		for _, v := range visitedar {
			for k, v1 := range a[v] {
				if v1 < min && visited[k] == false && v1 != 0 {
					i = v
					j = k
					min = v1
				}
			}

		}
		res[i][j] = min
		visited[j] = true
		visitedar = append(visitedar, j)
	}
	return res
}

func dijkstra(GR [][]float64, V, st int) (distance []float64, path []float64) {
	distance = make([]float64, V)
	path = make([]float64, V)
	visited := make([]bool, V)

	for i := 0; i < V; i++ {
		distance[i] = math.MaxInt32
		visited[i] = false
	}

	distance[st] = 0
	path[st] = -1

	for count := 0; count < V-1; count++ {
		min := math.MaxFloat32
		var u int

		for i := 0; i < V; i++ {
			if !visited[i] && distance[i] <= min {
				min = distance[i]
				u = i
			}
		}

		visited[u] = true

		for i := 0; i < V; i++ {
			if !visited[i] && GR[u][i] != 0 && distance[u] != math.MaxInt32 && distance[u]+GR[u][i] < distance[i] {
				distance[i] = distance[u] + GR[u][i]
				path[i] = float64(u)
			}
		}
	}

	return distance, path
}

func main() {
	type Num struct {
		num int
	}
	setnum := func(n Num) int {
		return n.num
	}
	path := os.Getenv("PATH")

	if !strings.Contains(path, `GRAPHS\bin`) {
		// Добавление пути к dot.exe в переменную окружения PATH
		pathenv, _ := findDotExe()
		err := addToPath(pathenv)
		if err != nil {

			return
		}
		if err := addToPath(pathenv); err != nil {

			return
		}
	}

	//core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowFlags(core.Qt__Window | core.Qt__WindowFullscreenButtonHint)

	// Максимизируем окно
	window.ShowMaximized()
	window.SetMinimumSize2(600, 400)
	window.SetWindowTitle("GraphBox")

	layout := widgets.NewQVBoxLayout()

	startlayout := widgets.NewQVBoxLayout()

	title := widgets.NewQLabel2("Расчет параметров системы управления в имитационной модели\n", nil, 0)
	font := title.Font()
	font.SetPointSize(16)
	title.SetFont(font)

	nInput := widgets.NewQLineEdit(nil)
	font = nInput.Font()
	font.SetPointSize(13)
	nInput.SetFont(font)
	nInput.SetPlaceholderText("Введите количество вершин:")
	nInput.SetFixedHeight(30)
	nInput.SetFixedWidth(500)
	startlayout.AddWidget(title, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	startlayout.AddWidget(nInput, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)

	createTableBtn := widgets.NewQPushButton2("Ввести матрицу смежности графа управления с весами рёбер:", nil)
	createTableBtn.SetFixedSize2(500, 30)
	font = createTableBtn.Font()
	font.SetPointSize(13)
	createTableBtn.SetFont(font)
	startlayout.AddWidget(createTableBtn, 0, core.Qt__AlignHCenter)
	tableWidget := widgets.NewQTableWidget(nil)
	//tableWidget.SetFixedSize2(800, 800)
	//tableWidget.SetMinimumSize2(300, 200)
	tableWidget.SetHidden(true)
	layout.AddLayout(startlayout, 0)
	layout.AddWidget(tableWidget, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	layout.SetAlignment2(startlayout, core.Qt__AlignHCenter|core.Qt__AlignTop)
	Hlayout := widgets.NewQHBoxLayout()
	chooseBtn := widgets.NewQPushButton2("Выберите действие", nil)
	font = chooseBtn.Font()
	font.SetPointSize(13)
	chooseBtn.SetFont(font)
	chooseBtn.SetHidden(true)
	layout.AddWidget(chooseBtn, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	checkboxlayout := widgets.NewQHBoxLayout()
	showgraphbtn := widgets.NewQPushButton2("Построить граф управления", nil)
	font = showgraphbtn.Font()
	font.SetPointSize(13)
	showgraphbtn.SetFont(font)
	showgraphbtn.SetHidden(true)
	Hlayout.AddWidget(showgraphbtn, 0, 0)
	ALGD := widgets.NewQPushButton2("Выполнить Алгоритм Дейкстры", nil)
	ALGD.SetToolTip("Алгори́тм Де́йкстры — алгоритм на графах,\n находит кратчайшие пути от одной из вершин графа до всех остальных.")
	font = ALGD.Font()
	font.SetPointSize(13)
	ALGD.SetFont(font)
	ALGD.SetHidden(true)
	Hlayout.AddWidget(ALGD, 0, 0)
	KSMD := widgets.NewQPushButton2("Построить КСМД", nil)
	KSMD.SetToolTip("Коммуникационная сеть минимальной длины - \n это совокупность дуг сети,\n имеющая минимальную суммарную длину\n и обеспечивабщая достижение всех узлов сети")
	KSMD.SetHidden(true)
	font = KSMD.Font()
	font.SetPointSize(13)
	KSMD.SetFont(font)
	Hlayout.AddWidget(KSMD, 0, 0)
	layout.AddLayout(Hlayout, 0)
	layout.SetAlignment2(Hlayout, core.Qt__AlignHCenter|core.Qt__AlignTop)
	power := widgets.NewQPushButton2("Расчитать степень вершины графа управления", nil)
	power.SetToolTip("Степень вершины графа — количество рёбер\n графа G, инцидентных вершине x.")
	font = power.Font()
	font.SetPointSize(13)
	power.SetFont(font)
	power.SetHidden(true)
	Hlayout.AddWidget(power, 0, 0)
	checkbox := widgets.NewQCheckBox2("до всех вершин", window)
	font = checkbox.Font()
	font.SetPointSize(13)
	checkbox.SetFont(font)
	checkbox.SetHidden(true)
	checkboxlayout.AddWidget(checkbox, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	checkbox2 := widgets.NewQCheckBox2("до конкретной вершины", window)
	font = checkbox2.Font()
	font.SetPointSize(13)
	checkbox2.SetFont(font)
	checkbox2.SetHidden(true)
	checkboxlayout.AddWidget(checkbox2, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	layout.AddLayout(checkboxlayout, 0)
	layout.SetAlignment2(checkboxlayout, core.Qt__AlignHCenter|core.Qt__AlignTop)
	verInput := widgets.NewQLineEdit(nil)
	verInput.SetPlaceholderText("Введите начальную вершину:")
	font = verInput.Font()
	font.SetPointSize(13)
	verInput.SetFont(font)
	verInput.SetHidden(true)
	layout.AddWidget(verInput, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	verInput2 := widgets.NewQLineEdit(nil)
	verInput2.SetPlaceholderText("Введите конечную вершину:")
	font = verInput2.Font()
	font.SetPointSize(13)
	verInput2.SetFont(font)
	verInput2.SetHidden(true)
	layout.AddWidget(verInput2, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	calculatebtn := widgets.NewQPushButton2("Расчитать", nil)
	calculatebtn.SetHidden(true)
	font = calculatebtn.Font()
	font.SetPointSize(13)
	calculatebtn.SetFont(font)
	layout.AddWidget(calculatebtn, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	var n int
	scrollArea := widgets.NewQScrollArea(nil)
	scrollArea1 := widgets.NewQScrollArea(nil)
	label2 := widgets.NewQLabel(nil, 0)

	layout.SetAlignment(nInput, core.Qt__AlignHCenter|core.Qt__AlignTop)
	layout.SetAlignment(tableWidget, core.Qt__AlignHCenter|core.Qt__AlignTop)
	checkbox.ConnectToggled(func(checked bool) {
		verInput.SetVisible(checked)
		if checked {
			checkbox2.SetChecked(false)
			verInput2.SetVisible(false)
		}
	})
	checkbox2.ConnectToggled(func(checked bool) {
		verInput2.SetVisible(checked)
		verInput.SetVisible(checked)
		if checked {
			checkbox.SetChecked(false)
			verInput.SetVisible(true)
		}
	})

	power.ConnectClicked(func(checked bool) {
		checkbox.Hide()
		checkbox2.Hide()
		label2.Hide()
		scrollArea.SetHidden(true)
		calculatebtn.SetHidden(true)
		verInput2.SetHidden(true)
		label2.Hide()
		var values [][]int
		for i := 0; i < tableWidget.RowCount(); i++ {
			var row []int
			for j := 0; j < tableWidget.ColumnCount(); j++ {
				item := tableWidget.Item(i, j)
				if item != nil {
					n, _ := strconv.Atoi(item.Text())
					row = append(row, n)

				} else {
					row = append(row, 0)
				}
			}
			values = append(values, row)
		}
		for k, v := range values {
			for k1, _ := range v {
				if values[k][k1] != 0 {
					values[k1][k] = values[k][k1]
				}
			}
		}
		str := ""
		for vershina := 1; vershina <= n; vershina++ {
			count := 0
			for _, v := range values[vershina-1] {
				if v != 0 {
					count++
				}
			}
			str += fmt.Sprintf("Степень %d-ой вершины равна %d \n", vershina, count)
		}

		label2 = widgets.NewQLabel2(str, nil, 0)
		font = label2.Font()
		font.SetPointSize(13)
		label2.SetFont(font)
		layout.AddWidget(label2, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	})
	errlabel := widgets.NewQLabel(window, 0)

	showgraphbtn.ConnectClicked(func(checked bool) {
		checkbox.Hide()
		checkbox2.Hide()
		label2.Hide()
		verInput.Hide()
		verInput2.Hide()
		scrollArea.Hide()
		scrollArea1.Hide()
		tableWidget.SetHidden(false)
		calculatebtn.SetHidden(true)
		createTableBtn.SetHidden(false)
		var values [][]float64
		for i := 0; i < tableWidget.RowCount(); i++ {
			var row []float64
			for j := 0; j < tableWidget.ColumnCount(); j++ {
				item := tableWidget.Item(i, j)
				if item != nil {
					n, _ := strconv.ParseFloat(item.Text(), 64)
					row = append(row, n)

				} else {
					row = append(row, 0)
				}
			}
			values = append(values, row)
		}
		for k, v := range values {
			for k1, _ := range v {
				if values[k][k1] != 0 {
					values[k1][k] = values[k][k1]
				}
			}
		}
		g := graph.New(setnum)
		for i := 0; i < n; i++ {
			_ = g.AddVertex(Num{num: i + 1})
		}
		for i := 0; i < len(values); i++ {
			for j := 0; j < len(values); j++ {
				if values[i][j] != 0 {
					_ = g.AddEdge(i+1, j+1, graph.EdgeAttribute("label", fmt.Sprintf("%.2f", values[i][j])))
				}
			}
		}
		file, _ := os.Create(`image\my-graph.gv`)
		_ = draw.DOT(g, file)
		cmd := exec.Command("cmd", "/C", "dot", "-Tpng", "-O", `image\my-graph.gv`)
		_, err := cmd.Output()
		if err != nil {

			return
		}
		contentWidget := widgets.NewQWidget(nil, 0)
		vlayout := widgets.NewQVBoxLayout2(contentWidget)
		font := label2.Font()
		font.SetPointSize(13)
		label2.SetFont(font)
		layout.AddWidget(label2, 0, 0)
		lbl := widgets.NewQLabel2("Граф", nil, 0)
		font = lbl.Font()
		font.SetPointSize(13)
		lbl.SetFont(font)
		vlayout.AddWidget(lbl, 0, core.Qt__AlignHCenter)
		imageView := widgets.NewQLabel(nil, 0)
		pixmap := gui.NewQPixmap3(`image\my-graph.gv.png`, "", core.Qt__AutoColor)
		imageView.SetPixmap(pixmap)
		vlayout.AddWidget(imageView, 0, core.Qt__AlignHCenter)
		scrollArea = widgets.NewQScrollArea(nil)
		scrollArea.SetWidgetResizable(true)
		scrollArea.SetWidget(contentWidget)
		scrollArea.SetMinimumSize2(500, 500)
		scrollArea.SetFixedHeight(imageView.Height() + lbl.Height() + 100)
		layout.AddWidget(scrollArea, 1, 0)
		layout.SetAlignment(scrollArea, core.Qt__AlignHCenter|core.Qt__AlignTop)
	})

	KSMD.ConnectClicked(func(checked bool) {
		calculatebtn.Hide()
		checkbox.Hide()
		checkbox2.Hide()
		scrollArea.Hide()
		scrollArea1.Hide()
		calculatebtn.Hide()
		verInput.Hide()
		verInput2.Hide()
		scrollArea.Hide()
		scrollArea1.Hide()
		label2.Hide()
		var values [][]float64
		for i := 0; i < tableWidget.RowCount(); i++ {
			var row []float64
			for j := 0; j < tableWidget.ColumnCount(); j++ {
				item := tableWidget.Item(i, j)
				if item != nil {
					n, _ := strconv.ParseFloat(item.Text(), 64)

					row = append(row, n)

				} else {
					row = append(row, 0)
				}
			}
			values = append(values, row)
		}
		if IsGraphConnected(values) {
			com := MinLenthCom(values)
			var sum float64
			for i := 0; i < tableWidget.RowCount(); i++ {
				for j := 0; j < tableWidget.ColumnCount(); j++ {
					sum += com[i][j]
				}
			}

			g := graph.New(setnum)
			for i := 0; i < n; i++ {
				_ = g.AddVertex(Num{num: i + 1})
			}
			for i := 0; i < len(com); i++ {
				for j := 0; j < len(com); j++ {
					if com[i][j] != 0 {
						_ = g.AddEdge(i+1, j+1, graph.EdgeAttribute("label", fmt.Sprintf("%.2f", com[i][j])))
					}
				}
			}
			file, _ := os.Create(`image\my-graph.gv`)
			_ = draw.DOT(g, file)
			cmd := exec.Command("cmd", "/C", "dot", "-Tpng", "-O", `image\my-graph.gv`)
			_, err := cmd.Output()
			if err != nil {

				return
			}
			label2 = widgets.NewQLabel2(fmt.Sprintf("Длина коммуникационной сети равна %.2f", sum), nil, 0)
			contentWidget := widgets.NewQWidget(nil, 0)
			vlayout := widgets.NewQVBoxLayout2(contentWidget)
			font := label2.Font()
			font.SetPointSize(13)
			label2.SetFont(font)
			layout.AddWidget(label2, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
			lbl := widgets.NewQLabel2("КМСД", nil, 0)
			font = lbl.Font()
			font.SetPointSize(13)
			lbl.SetFont(font)
			vlayout.AddWidget(lbl, 0, core.Qt__AlignHCenter)
			imageView := widgets.NewQLabel(nil, 0)
			pixmap := gui.NewQPixmap3(`image\my-graph.gv.png`, "", core.Qt__AutoColor)
			imageView.SetPixmap(pixmap)
			vlayout.AddWidget(imageView, 0, core.Qt__AlignHCenter)
			scrollArea = widgets.NewQScrollArea(nil)
			scrollArea.SetWidgetResizable(true)
			scrollArea.SetWidget(contentWidget)
			scrollArea.SetMinimumSize2(500, 500)
			layout.AddWidget(scrollArea, 1, 0)
			layout.SetAlignment(scrollArea, core.Qt__AlignHCenter|core.Qt__AlignTop)
		} else {
			label2 = widgets.NewQLabel2("Граф несвязный, невозможно построить коммуникационную сеть", nil, 0)
			font := label2.Font()
			font.SetPointSize(13)
			label2.SetFont(font)
			layout.AddWidget(label2, 0, 0)
		}

	})
	ALGD.ConnectClicked(func(checked bool) {
		label2.Hide()
		scrollArea.Hide()
		errlabel.SetHidden(true)

		verInput.SetHidden(true)
		calculatebtn.SetHidden(false)
		checkbox.SetHidden(false)
		checkbox2.SetHidden(false)
	})
	createTableBtn.ConnectClicked(func(checked bool) {
		if nInput.Text() == "" || nInput.Text() == "0" {
			return
		}
		n, _ = strconv.Atoi(nInput.Text())
		if n == 0 {
			return
		}
		if n <= 20 {
			tableWidget.SetFixedSize2(n*110, n*40)
		} else {
			tableWidget.SetFixedSize2(16*110, 20*40)
		}
		tableWidget.Clear()
		tableWidget.SetRowCount(int(n))
		tableWidget.SetColumnCount(int(n))
		tableWidget.SetHidden(false)
		chooseBtn.SetHidden(false)

		//createTableBtn.SetHidden(true)
	})
	chooseBtn.ConnectClicked(func(checked bool) {
		if ALGD.IsHidden() {
			KSMD.SetHidden(false)
			ALGD.SetHidden(false)
			showgraphbtn.SetHidden(false)
			power.SetHidden(false)
		} else {
			KSMD.SetHidden(true)
			ALGD.SetHidden(true)
			showgraphbtn.SetHidden(true)
			power.SetHidden(true)
			label2.Hide()
			checkbox.Hide()
			checkbox2.Hide()
			scrollArea.Hide()
			scrollArea1.Hide()
			verInput.Hide()
			verInput2.Hide()
		}
	})

	calculatebtn.ConnectClicked(func(checked bool) {
		label2.Hide()
		scrollArea.Hide()
		scrollArea1.Hide()
		tableWidget.SetHidden(false)
		calculatebtn.SetHidden(false)
		createTableBtn.SetHidden(false)
		var values [][]float64
		for i := 0; i < tableWidget.RowCount(); i++ {
			var row []float64
			for j := 0; j < tableWidget.ColumnCount(); j++ {
				item := tableWidget.Item(i, j)
				if item != nil {
					n, _ := strconv.ParseFloat(item.Text(), 64)
					row = append(row, n)

				} else {
					row = append(row, 0)
				}
			}
			values = append(values, row)
		}
		for k, v := range values {
			for k1, _ := range v {
				if values[k][k1] != 0 {
					values[k1][k] = values[k][k1]
				}
			}
		}
		var dijkstrarez string
		vershina := verInput.Text()
		endvershina := verInput2.Text()
		matrixpath := make([][]float64, len(values))
		for k, _ := range matrixpath {
			matrixpath[k] = make([]float64, len(values))
		}
		if vershina == "" {
			// for start := 0; start < len(values); start++ {
			// 	distance := dijkstra(values, len(values), start)
			// 	m := start + 1
			// 	dijkstrarez += fmt.Sprintf("Оптимальные пути из %d вершины до остальных:\n", m)
			// 	for i := 0; i < len(values); i++ {
			// 		if distance[i] != math.MaxInt32 {
			// 			if m != i+1 {
			// 				dijkstrarez += fmt.Sprintf("%d > %d = %d\n", m, i+1, distance[i])
			// 			}
			// 		} else {
			// 			dijkstrarez += fmt.Sprintf("%d > %d = маршрут недоступен\n", m, i+1)
			// 		}
			// 	}
			// }
		} else if vershina != "" {
			if !checkbox2.IsChecked() {
				vershinaint, _ := strconv.Atoi(vershina)
				if vershinaint > n {
					return
				}
				distance, _ := dijkstra(values, len(values), vershinaint-1)

				m := vershinaint
				dijkstrarez += fmt.Sprintf("Оптимальные пути из %d-го динамического звена до остальных:\n", m)
				for i := 0; i < len(values); i++ {
					if distance[i] != math.MaxInt32 {
						if m != i+1 {
							dijkstrarez += fmt.Sprintf("%d -> %d = %f(усл. ед.)\n", m, i+1, distance[i])
						}
					} else {
						dijkstrarez += fmt.Sprintf("%d -> %d = маршрут недоступен\n", m, i+1)
					}
				}
			} else if checkbox2.IsChecked() && endvershina == "" {

			} else if checkbox2.IsChecked() && endvershina != "" {
				vershinaint, _ := strconv.Atoi(vershina)
				vershinaint2, _ := strconv.Atoi(endvershina)
				if vershinaint > n || vershinaint2 > n {
					return
				}
				if isVerticesConnected(values, vershinaint-1, vershinaint2-1) {

					distance, path := dijkstra(values, len(values), vershinaint-1)
					matrixpath = getPath(path, float64(vershinaint2-1))
					m := vershinaint
					dijkstrarez += fmt.Sprintf("Оптимальные пути из %d-го динамического звена до %d вершины:\n", m, vershinaint2)

					if distance[vershinaint2-1] != math.MaxInt32 {
						if m != vershinaint2 {
							dijkstrarez += fmt.Sprintf("%d -> %d = %f(усл. ед.)\n", m, vershinaint2, distance[vershinaint2-1])
						}
					} else {
						dijkstrarez += fmt.Sprintf("%d -> %d = маршрут недоступен\n", m, vershinaint2)
					}
				} else {

					distance, _ := dijkstra(values, len(values), vershinaint-1)
					m := vershinaint
					dijkstrarez += fmt.Sprintf("Оптимальные пути из %d-го динамического звена до остальных:\n", m)
					if distance[vershinaint2-1] != math.MaxInt32 {
						if m != vershinaint2 {
							dijkstrarez += fmt.Sprintf("%d -> %d = %f(усл. ед.)\n", m, vershinaint2, distance[vershinaint2-1])
						}
					} else {
						dijkstrarez += fmt.Sprintf("%d -> %d = маршрут недоступен\n", m, vershinaint2)
					}
				}

			}
		}

		g := graph.New(setnum)
		for i := 0; i < n; i++ {
			_ = g.AddVertex(Num{num: i + 1})
		}
		for i := 0; i < len(values); i++ {
			for j := 0; j < len(values); j++ {
				if values[i][j] != 0 {
					if matrixpath[i][j] != 1 {
						_ = g.AddEdge(i+1, j+1, graph.EdgeAttribute("label", fmt.Sprintf("%.2f", values[i][j])))
					} else {
						_ = g.AddEdge(i+1, j+1, graph.EdgeAttribute("color", "green"), graph.EdgeAttribute("label", fmt.Sprintf("%.2f", values[i][j])))
					}
				}
			}
		}
		file, _ := os.Create(`image\my-graph.gv`)
		_ = draw.DOT(g, file)
		cmd := exec.Command("cmd", "/C", "dot", "-Tpng", "-O", `image\my-graph.gv`)
		_, err := cmd.Output()
		if err != nil {

			return
		}
		label2 = widgets.NewQLabel2(dijkstrarez, nil, 0)
		contentWidget := widgets.NewQWidget(nil, 0)
		vlayout := widgets.NewQVBoxLayout2(contentWidget)
		font := label2.Font()
		font.SetPointSize(13)
		label2.SetFont(font)
		layout.AddWidget(label2, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
		lbl := widgets.NewQLabel2("Граф", nil, 0)
		font = lbl.Font()
		font.SetPointSize(13)
		lbl.SetFont(font)
		vlayout.AddWidget(lbl, 0, core.Qt__AlignHCenter)
		imageView := widgets.NewQLabel(nil, 0)
		pixmap := gui.NewQPixmap3(`image\my-graph.gv.png`, "", core.Qt__AutoColor)
		imageView.SetPixmap(pixmap)
		vlayout.AddWidget(imageView, 0, core.Qt__AlignHCenter)
		scrollArea = widgets.NewQScrollArea(nil)
		scrollArea.SetWidgetResizable(true)
		scrollArea.SetWidget(contentWidget)
		scrollArea.SetMinimumSize2(500, 400)
		scrollArea.SetFixedHeight(imageView.Height() + lbl.Height() + 100)
		layout.AddWidget(scrollArea, 1, 0)
		layout.SetAlignment(scrollArea, core.Qt__AlignHCenter|core.Qt__AlignTop)
	})
	scroll := widgets.NewQScrollArea(window)
	content := widgets.NewQWidget(nil, 0)
	content.SetLayout(layout)
	scroll.SetWidget(content)
	scroll.SetWidgetResizable(true)
	scroll.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOn)
	scroll.SetFixedSize(window.Size())
	createTableBtn.Move2(createTableBtn.X(), 500)
	mainlayout := widgets.NewQVBoxLayout()
	mainlayout.AddWidget(scroll, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(mainlayout)
	window.SetCentralWidget(widget)

	window.Show()

	app.Exec()
}
