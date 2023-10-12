package checkers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	BOARD_DIM = 8
	RED       = "red"
	BLACK     = "black"
	ROW_SEP   = "|"
)

type Player struct {
	Color string
}

type Piece struct {
	Player Player
	King   bool
}

var PieceStrings = map[Player]string{
	RED_PLAYER:   "r",
	BLACK_PLAYER: "b",
	NO_PLAYER:    "*",
}

var NO_PIECE = Piece{NO_PLAYER, false}

var StringPieces = map[string]Piece{
	"r": Piece{RED_PLAYER, false},
	"b": Piece{BLACK_PLAYER, false},
	"R": Piece{RED_PLAYER, true},
	"B": Piece{BLACK_PLAYER, true},
	"*": NO_PIECE,
}

type Pos struct {
	X int
	Y int
}

var NO_POS = Pos{-1, -1}

var BLACK_PLAYER = Player{BLACK}
var RED_PLAYER = Player{RED}
var NO_PLAYER = Player{
	Color: "NO_PLAYER",
}

var Players = map[string]Player{
	RED:   RED_PLAYER,
	BLACK: BLACK_PLAYER,
}

var Opponents = map[Player]Player{
	BLACK_PLAYER: RED_PLAYER,
	RED_PLAYER:   BLACK_PLAYER,
}

// solidity의 2중 3중 mapping과 동일하게 동작
var Usable = map[Pos]bool{}
var Moves = map[Player]map[Pos]map[Pos]bool{}
var Jumps = map[Player]map[Pos]map[Pos]Pos{}
var KingMoves = map[Pos]map[Pos]bool{}
var KingJumps = map[Pos]map[Pos]Pos{}

// 상대 말을 잡을 때, 상대 말을 뛰어넘으므로 (이전 + 이후 위치) / 2는 잡히는 말의 위치
func Capture(src, dst Pos) Pos {
	return Pos{(src.X + dst.X) / 2, (src.Y + dst.Y) / 2}
}

func init() {

 	// 01010101
 	// 10101010
  	// 01010101
  	// 10101010
  	// 01010101
  	// 10101010
  	// 01010101
  	// 10101010
	// Initialize usable spaces
	for y := 0; y < BOARD_DIM; y++ {
		for x := (y + 1) % 2; x < BOARD_DIM; x += 2 {
			Usable[Pos{X: x, Y: y}] = true
		}
	}

  	// Go에서 중첩된 맵에 접근, 할당을 시도하면 nil map 오류가 발생
  	// 사용 전에 외부 맵 각각에 내부 맵을 할당해야 함
  	// var Moves = map[Player]map[Pos]map[Pos]bool{}
  	// 즉, Moves[RED], Moves[BLACK] 할당
	for _, p := range Players {
		Moves[p] = map[Pos]map[Pos]bool{}
		Jumps[p] = map[Pos]map[Pos]Pos{}
	}

  
	for pos := range Usable {
		// KingMoves = map[Pos]map[Pos]bool{
    	// KingMoves[true], KingMoves[false] 할당
    	KingMoves[pos] = map[Pos]bool{}
		KingJumps[pos] = map[Pos]Pos{}
    	// 방향 좌(-1), 우(1)
		var directions = []int{1, -1}
		for i, player := range []Player{BLACK_PLAYER, RED_PLAYER} {
      		// Moves[Red][true], Moves[Red][false], Moves[BLACK][true], Moves[BLACK][false]
			Moves[player][pos] = map[Pos]bool{}
			Jumps[player][pos] = map[Pos]Pos{}
			// 이동(대각선으로만 이동)
      		movOff := 1
      		// 잡기
			jmpOff := 2
	      	for _, direction := range directions {
				// 상하좌우 대각선 mov
	        	mov := Pos{pos.X + (movOff * direction), pos.Y + (movOff * directions[i])}
				// mov 위치가 움직일 수 있는 위치면
	        	if Usable[mov] {
	          		// 움직일 수 있는 위치 true로 변경
					Moves[player][pos][mov] = true
	          		// 킹은 방향에 상관없이 자유롭게 움직이므로 방향(player)은 관계없음
					KingMoves[pos][mov] = true
				}

		        // 눈 목자로 잡는 위치들 mov
				jmp := Pos{pos.X + (jmpOff * direction), pos.Y + (jmpOff * directions[i])}
				// 움직이는게 가능한 경우
		        if Usable[jmp] {
		          	// capturePos의 좌표를 찾은 다음
					capturePos := Capture(pos, jmp)
		          	// 점프했을 때 잡히는 적의 말의 위치 저장
					Jumps[player][pos][jmp] = capturePos
					KingJumps[pos][jmp] = capturePos
				}
			}
		}
	}
}

type Game struct {
	Pieces map[Pos]Piece
	Turn   Player
}

// 새 게임을 위한 인스턴스 생성
func New() *Game {
	// 말 mapping 초기화
	pieces := make(map[Pos]Piece)
	// 게임 인스턴스 초기화
	game := &Game{pieces, BLACK_PLAYER}
	// 말 생성
	game.addInitialPieces()
	return game
}

func (game *Game) addInitialPieces() {
	for pos := range Usable {
		// 흑은 위 3열
		// 1, 2, 3
		if pos.Y >= 0 && pos.Y < 3 {
			game.Pieces[pos] = Piece{BLACK_PLAYER, false}
		}
		
		// 백은 아래 3열
		// 6, 7, 8
		if pos.Y >= BOARD_DIM-3 && pos.Y < BOARD_DIM {
			game.Pieces[pos] = Piece{RED_PLAYER, false}
		}
	}
}

// 해당 위치에 말이 있는지 체크
func (game *Game) PieceAt(pos Pos) bool {
	_, ok := game.Pieces[pos]
	return ok
}

// 해당 플레이어의 턴인지 체크
func (game *Game) TurnIs(player Player) bool {
	return game.Turn == player
}

// 승패 결정
func (game *Game) Winner() Player {
	red_count := 0
	black_count := 0
	// 양 쪽의 말의 개수 세기
	for _, piece := range game.Pieces {
		switch {
		case piece.Player == BLACK_PLAYER:
			black_count += 1
		case piece.Player == RED_PLAYER:
			red_count += 1
		}
	}
	// 검정이 1개 이상이고 빨강이 없으면 검정 승리
	if black_count > 0 && red_count <= 0 {
		return BLACK_PLAYER
	// 반대의 경우 빨강 승리
	} else if red_count > 0 && black_count <= 0 {
		return RED_PLAYER
	}
	// 아닌 경우 아무도 이기지 않았음
	return NO_PLAYER
}

// 움직일 수 있는 위치인지 체크
func (game *Game) ValidMove(src, dst Pos) bool {
	// 기존 위치에 말이 없거나 목표 위치에 이미 말이 있다면 
	if !game.PieceAt(src) || game.PieceAt(dst) {
		return false
	}
	piece := game.Pieces[src]
	if (!piece.King && Moves[piece.Player][src][dst]) || (piece.King && KingMoves[src][dst]) {
		return !game.playerHasJump(piece.Player)
	}
	return game.ValidJump(src, dst)
}

// 잡을 수 있는 위치인지 체크
func (game *Game) ValidJump(src, dst Pos) bool {
	if !game.PieceAt(src) || game.PieceAt(dst) {
		return false
	}
	piece := game.Pieces[src]
	if !piece.King {
		capLoc, jumpOk := Jumps[piece.Player][src][dst]
		return jumpOk && game.PieceAt(capLoc) && game.Pieces[capLoc].Player == Opponents[piece.Player]
	} else {
		capLoc, kingJumpOk := KingJumps[src][dst]
		return kingJumpOk && game.PieceAt(capLoc) && game.Pieces[capLoc].Player == Opponents[piece.Player]
	}
}

func (game *Game) kingPiece(dst Pos) {
	if !game.PieceAt(dst) {
		return
	}
	piece := game.Pieces[dst]
	if (dst.Y == 0 && piece.Player == RED_PLAYER) ||
		(dst.Y == BOARD_DIM-1 && piece.Player == BLACK_PLAYER) {
		piece.King = true
		game.Pieces[dst] = piece
	}
}

func (game *Game) updateTurn(dst Pos, jumped bool) {
	opponent := Opponents[game.Turn]
	if (!jumped || !game.jumpPossibleFrom(dst)) && game.playerHasMove(opponent) {
		game.Turn = opponent
	}
}

func (game *Game) jumpPossibleFrom(src Pos) bool {
	if !game.PieceAt(src) {
		return false
	}
	piece := game.Pieces[src]
	if !piece.King {
		// enumerate all player jumps and return true if one is valid
		for dst := range Jumps[piece.Player][src] {
			if game.ValidJump(src, dst) {
				return true
			}
		}
	} else {
		// enumerate all king jumps and return true if one is valid
		for dst := range KingJumps[src] {
			if game.ValidJump(src, dst) {
				return true
			}
		}
	}
	return false
}

func (game *Game) movePossibleFrom(src Pos) bool {
	if !game.PieceAt(src) {
		return false
	}
	piece := game.Pieces[src]
	if !piece.King {
		for dst := range Moves[piece.Player][src] {
			if game.ValidMove(src, dst) {
				return true
			}
		}
	} else {
		for dst := range KingMoves[src] {
			if game.ValidMove(src, dst) {
				return true
			}
		}
	}
	return false
}

func (game *Game) playerHasMove(player Player) bool {
	for loc, piece := range game.Pieces {
		if piece.Player == player && (game.movePossibleFrom(loc) || game.jumpPossibleFrom(loc)) {
			return true
		}
	}
	return false
}

func (game *Game) playerHasJump(player Player) bool {
	for loc, piece := range game.Pieces {
		if piece.Player == player && game.jumpPossibleFrom(loc) {
			return true
		}
	}
	return false
}

func (game *Game) Move(src, dst Pos) (captured Pos, err error) {
	captured = NO_POS
	err = nil
	if !game.PieceAt(src) {
		return NO_POS, errors.New(fmt.Sprintf("No piece at source position: %v", src))
	}
	if game.PieceAt(dst) {
		return NO_POS, errors.New(fmt.Sprintf("Already piece at destination position: %v", dst))
	}
	if !game.TurnIs(game.Pieces[src].Player) {
		return NO_POS, errors.New(fmt.Sprintf("Not %v's turn", game.Pieces[src].Player))
	}
	if !game.ValidMove(src, dst) {
		return NO_POS, errors.New(fmt.Sprintf("Invalid move: %v to %v", src, dst))
	}
	if game.ValidJump(src, dst) {
		game.Pieces[dst] = game.Pieces[src]
		delete(game.Pieces, src)
		captured = Capture(src, dst)
		delete(game.Pieces, captured)
	} else {
		game.Pieces[dst] = game.Pieces[src]
		delete(game.Pieces, src)
	}
	game.updateTurn(dst, captured != NO_POS)
	game.kingPiece(dst)
	return
}

func (game *Game) String() string {
	var buf bytes.Buffer
	for y := 0; y < BOARD_DIM; y++ {
		for x := 0; x < BOARD_DIM; x++ {
			pos := Pos{x, y}
			if game.PieceAt(pos) {
				piece := game.Pieces[pos]
				val := PieceStrings[piece.Player]
				if piece.King {
					val = strings.ToUpper(val)
				}
				buf.WriteString(val)
			} else {
				buf.WriteString(PieceStrings[NO_PLAYER])
			}
		}
		if y < (BOARD_DIM - 1) {
			buf.WriteString(ROW_SEP)
		}
	}
	return buf.String()
}

func ParsePiece(s string) (Piece, bool) {
	piece, ok := StringPieces[s]
	return piece, ok
}

func Parse(s string) (*Game, error) {
	if len(s) != BOARD_DIM*BOARD_DIM+(BOARD_DIM-1) {
		return nil, errors.New(fmt.Sprintf("invalid board string: %v", s))
	}
	pieces := make(map[Pos]Piece)
	result := &Game{pieces, BLACK_PLAYER}
	for y, row := range strings.Split(s, ROW_SEP) {
		for x, c := range strings.Split(row, "") {
			if x >= BOARD_DIM || y >= BOARD_DIM {
				return nil, errors.New(fmt.Sprintf("invalid board, piece out of bounds: %v, %v", x, y))
			}
			if piece, ok := ParsePiece(c); !ok {
				return nil, errors.New(fmt.Sprintf("invalid board, invalid piece at %v, %v", x, y))
			} else if piece != NO_PIECE {
				result.Pieces[Pos{x, y}] = piece
			}
		}
	}
	return result, nil
}
