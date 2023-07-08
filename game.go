package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/eiannone/keyboard"
)

type GameSig int

const (
	charWidth   = 6    // " 4096 "
	defaultSize = 4    // 4 X 4
	maxNum      = 2048 // game finished

	UpSig    GameSig = 1
	DownSig  GameSig = 2
	LeftSig  GameSig = 3
	RightSig GameSig = 4
)

func (sig GameSig) String() string {
	switch sig {
	case UpSig:
		return "TOP"
	case DownSig:
		return "DOWN"
	case LeftSig:
		return "LEFT"
	case RightSig:
		return "RIGHT"
	default:
		return ""
	}
}

type Game struct {
	x    int
	data [][]int
	nums []int

	max  int
	step int

	finished int32
}

func New(x ...int) *Game {
	size := defaultSize
	if len(x) > 0 && x[0] > 2 {
		size = x[0]
	}

	game := new(Game)
	game.x = size
	game.data = make([][]int, game.x)
	for i := 0; i < game.x; i++ {
		game.data[i] = make([]int, game.x)
	}
	game.nums = []int{2, 4}

	rand.Seed(time.Now().UnixNano())
	game.rand()
	return game
}

func (game *Game) Run() {
	game.Refresh()

	kl, err := newKeyboartListener()
	if err != nil {
		fmt.Printf("[ERROR] keyboard listener error:%v", err)
		return
	}
	defer kl.close()

	for {
		if atomic.LoadInt32(&game.finished) == 1 {
			return
		}
		if len(game.fillList()) == 0 {
			game.Shutdown("finished(no pos to filled)")
			break
		}
		if game.max == maxNum {
			game.Shutdown("success")
			return
		}

		sig, err := kl.get()
		if err != nil {
			fmt.Printf("[ERROR] keyboard get error:%v", err)
			return
		}
		if len(sig.String()) == 0 {
			fmt.Printf("[ERROR] empty sig ")
			return
		}

		if changed := game.Action(sig); changed {
			game.rand()
		}

		game.Refresh(sig.String())
	}
}

func (game *Game) Action(sig GameSig) bool {
	getxy := func(i, j int) (int, int) {
		switch sig {
		case UpSig:
			// y not change, x increase
			return j, i
		case DownSig:
			// y not change, x decrease
			return game.x - 1 - j, i
		case LeftSig:
			// x not change, y increase
			return i, j
		case RightSig:
			// x not change, y decrease
			return i, game.x - 1 - j
		default:
			return -1, -1
		}
	}
	if len(sig.String()) == 0 {
		return false
	}

	var res bool
	for i := 0; i < game.x; i++ {
		arr := make([]int, 0, game.x)
		for j := 0; j < game.x; j++ {
			x, y := getxy(i, j)
			arr = append(arr, game.data[x][y])
		}

		b1 := move(arr)
		b2 := merge(arr)
		b3 := move(arr)
		res = res || b1 || b2 || b3
		for j := 0; j < game.x; j++ {
			x, y := getxy(i, j)
			game.data[x][y] = arr[j]
		}
	}

	if !res {
		return false
	}

	game.step++
	return res
}

func move(arr []int) bool {
	var step int
	for i := 0; i < len(arr); i++ {
		if arr[i] == 0 {
			step = step + 1
			continue
		}
		arr[i-step] = arr[i]
		if step > 0 {
			arr[i] = 0
		}
	}
	return step > 0
}

func merge(arr []int) bool {
	var step int
	for i := 0; i < len(arr)-1; i++ {
		if arr[i] == arr[i+1] {
			arr[i] = arr[i] + arr[i+1]
			arr[i+1] = 0
			step = step + 1
			continue
		}
	}

	return step > 0
}

func (game *Game) Shutdown(reason string) {
	atomic.StoreInt32(&game.finished, 1)
	game.Refresh()
	fmt.Printf("\nshutdown reason:%v.\n", reason)
}

func (game *Game) fillList() [][]int {
	fillList := make([][]int, 0)
	for i := 0; i < len(game.data); i++ {
		for j := 0; j < len(game.data); j++ {
			game.max = maxInt(game.max, game.data[i][j])
			if game.data[i][j] == 0 {
				fillList = append(fillList, []int{i, j})
			}
		}
	}
	return fillList
}

func (game *Game) rand() bool {
	fillList := game.fillList()
	if len(fillList) == 0 {
		return false
	}

	pos := rand.Intn(len(fillList))
	num := game.nums[rand.Intn(len(game.nums))]

	i, j := fillList[pos][0], fillList[pos][1]
	game.data[i][j] = num
	return true
}

func (game *Game) Refresh(direction ...string) {
	fmt.Printf("\033[H\033[2J")

	strlist := make([]string, 0)
	head := strings.Repeat("+"+strings.Repeat("-", charWidth), game.x) + "+"
	for i := 0; i < game.x; i++ {
		var line string
		for j := 0; j < game.x; j++ {
			line = line + "|" + num2str(game.data[i][j], charWidth)
			game.max = maxInt(game.max, game.data[i][j])
		}
		line = line + "|"
		strlist = append(strlist, head, line)
	}
	strlist = append(strlist, head)
	titles := []string{
		fmt.Sprintf("step: %v", game.step),
		fmt.Sprintf("max_num: %v", game.max),
	}
	if len(direction) > 0 && len(direction[0]) > 0 {
		titles = append(titles, fmt.Sprintf("last_direction: %v", direction[0]))
	}
	tails := []string{
		fmt.Sprintf("please input your direction: "),
	}
	fmt.Printf(strings.Join(titles, "\n") + "\n")
	fmt.Printf(strings.Join(strlist, "\n") + "\n")
	fmt.Printf(strings.Join(tails, "\n"))
}

func num2str(i int, width int) string {
	v := strconv.Itoa(i)
	left := width - len(v)
	return strings.Repeat(" ", left-(left-left/2)) + v + strings.Repeat(" ", left-left/2)
}
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type keyboartListener struct {
	ignoreNoDirEvent bool
}

func newKeyboartListener() (*keyboartListener, error) {
	if err := keyboard.Open(); err != nil {
		return nil, err
	}
	return &keyboartListener{}, nil
}
func (kl *keyboartListener) get() (GameSig, error) {
	ch, key, err := keyboard.GetKey()
	if err != nil {
		return GameSig(0), err
	}

	switch key {
	case keyboard.KeyCtrlW, keyboard.KeyCtrlI, keyboard.KeyArrowUp:
		return UpSig, nil
	case keyboard.KeyCtrlS, keyboard.KeyCtrlK, keyboard.KeyArrowDown:
		return DownSig, nil
	case keyboard.KeyCtrlA, keyboard.KeyCtrlJ, keyboard.KeyArrowLeft:
		return LeftSig, nil
	case keyboard.KeyCtrlD, keyboard.KeyCtrlL, keyboard.KeyArrowRight:
		return RightSig, nil
	}

	switch ch {
	case '8':
		return UpSig, nil
	case '2':
		return DownSig, nil
	case '4':
		return LeftSig, nil
	case '6':
		return RightSig, nil
	}

	return GameSig(0), nil
}

func (kl *keyboartListener) close() error {
	return keyboard.Close()
}
