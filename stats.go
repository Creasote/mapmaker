package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

/*
Stats to record:
Score

DPS (10s avg)
Spawns created (total)
Spawns created per second (10s average? or just multiple of spawners)
Spawns alive

Target health
Target armour



TODO: mini-graph dps
*/

// Points constants for scoring
const (
	pointsGotToGoal     = 1
	pointsSuccessfulHit = 1
	pointsManualSpawn   = 1 // Mobs spawned via clicking
	pointsAutoSpawn     = 2 // Mobs spawned via spawner
)

const (
	scoreboardWidth  = 256
	scoreboardHeight = 128
	scoreboardLocX   = 1200
	scoreboardLocY   = 20
)

func initScoreboard() {

	//	sb := createScoreboard()
	go statsPSUpdate()
	//go updateScoreboard()

}

func updateScoreboard() {
	scoreboardText := make([]string, 6)
	var scorePS float64
	var dmgPS float64
	//statsPSUpdate()
	//for {
	scoreTally := 0
	for _, s := range sps[1:11] {
		scoreTally += s
	}
	scorePS = float64(scoreTally) / 10
	dmgTally := 0.0
	for _, d := range dps[1:11] {
		dmgTally += d
	}
	dmgPS = dmgTally / 10

	scoreboardText[0] = fmt.Sprintf("Score: %d", score)
	scoreboardText[1] = fmt.Sprintf("(per second: %0.1f)", scorePS)
	scoreboardText[2] = fmt.Sprintf("DPS: %0.1f", dmgPS)
	scoreboardText[3] = fmt.Sprintf("Spawns alive: %d", len(spawn_list))
	scoreboardText[4] = fmt.Sprintf("Target health: %0.2f", target_list[0].health)
	scoreboardText[5] = fmt.Sprintf("Target armour: %0.2f", target_list[0].armour)

	sb := createScoreboard()
	for ind, txt := range scoreboardText {
		text.Draw(sb, txt, mplusNormalFont, 1, 14*(ind+2), color.White)
	}
	img_scoreboard = sb
	//time.Sleep(100 * time.Millisecond)
	//}
}

func createScoreboard() *ebiten.Image {
	sb := ebiten.NewImage(scoreboardWidth, scoreboardHeight)
	sb.Fill(color.Black)
	text.Draw(sb, "Scoreboard", mplusNormalFont, 1, 14, color.White)
	return sb
}

func statsPSUpdate() {
	for {
		for i := len(dps) - 1; i > 0; i-- {
			dps[i] = dps[i-1]
		}
		dps[0] = 0

		// Do the shift, which locks in the final value of the last second (sps[0] now becomes locked at sps[1], which can then be written out
		// to the the tally)
		for i := len(sps) - 1; i > 0; i-- {
			sps[i] = sps[i-1]
		}
		sps[0] = 0
		time.Sleep(1000 * time.Millisecond)
	}
}

func drawScoreboard(screen *ebiten.Image) {
	updateScoreboard()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(scoreboardLocX), float64(scoreboardLocY))
	screen.DrawImage(img_scoreboard, op)

}

func addScore(s int) {
	score += s
	sps[0] += s
}
