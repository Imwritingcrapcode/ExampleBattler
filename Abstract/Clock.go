package Abstract

import (
	"log"
	"strconv"
	"time"
)

//Clock for the Client
type Clock struct {
	Client        *ClientChannels
	timeLeft      int
	tickingPeriod int
	isTicking     bool
	skippedATurn  bool
	dcedTimes     int
	dcLeft        int
	ticker        *time.Ticker
}

func (self *Clock) TellTheTime() {
	if self.Client.Opponent == self.Client {
		//solo game
		if self.timeLeft <= ACTUALTURNLENGTH {
			msg := "Time:" + strconv.Itoa(self.timeLeft)
			self.Client.TimeOutput <- msg
		}
	} else {
		//2 ppl game
		var msg string
		if self.Client.Opponent.Clock.State() {
			self.Client.Opponent.Clock.TellTheTime()
		} else {
			if self.timeLeft <= ACTUALTURNLENGTH {
				msg = "Time:" + strconv.Itoa(self.timeLeft)
				self.Client.TimeOutput <- msg
				self.Client.Opponent.TimeOutput <- msg
			}
		}
	}
}

func (self *Clock) State() bool {
	return self.isTicking
}

func (self *Clock) StartTicking() {
	if self.isTicking {
		log.Println("The clock is already ticking!!")
	}
	self.ticker = time.NewTicker(1 * time.Second)
	self.isTicking = true
	self.timeLeft = TURNLENGTH
	for {
		<-self.ticker.C
		if self.isTicking {
			self.timeLeft -= 1
			if (ACTUALTURNLENGTH-self.timeLeft)%TICKEVERY == 0 {
				self.TellTheTime()
			}
			if self.timeLeft <= 0 {
				self.TimedOut()
				break
			}
		} else {
			break
		}
	}
}

func (self *Clock) StartTickingBot() {
	if self.isTicking {
		log.Println("[BOTCLOCK] The clock is already ticking!!")
	}
	self.ticker = time.NewTicker(1 * time.Second)
	self.isTicking = true
	self.timeLeft = TURNLENGTH
	for {
		<-self.ticker.C
		if self.isTicking {
			self.timeLeft -= 1
			if (ACTUALTURNLENGTH-self.timeLeft)%TICKEVERY == 0 {
				self.TellTheTime()
			}
			if self.timeLeft <= 0 {
				self.Stop()
				self.Client.HasGivenUp <- false
				break
			}
		} else {
			break
		}
	}
}

func (self *Clock) Stop() {
	if self.isTicking {
		self.isTicking = false
		self.timeLeft = 0
	} else {
		log.Panic("Clock is not ticking!!")
	}

}

func (self *Clock) TimedOut() {
	if self.skippedATurn {
		self.Client.GiveUp()
	} else {
		self.Client.HasGivenUp <- false
		self.skippedATurn = true
	}
	if self.isTicking {
		self.Stop()
	}
}
func (self *Clock) Disconnected() {
	log.Println("[INGAME]", self.Client.UserID, "disconnected!")
	if self.dcedTimes > 2 {
		//dced again, gg
		self.Client.GiveUp()
		log.Println("[INGAME]", self.Client.UserID, "lost the game due to dcing 3 times!")
	} else {
		if self.dcedTimes == 0 {
			self.dcLeft = MAXDCTIME
		}
		self.dcedTimes += 1
		//start counting
		timeLeft := self.dcLeft
		timer := time.NewTimer(time.Duration(self.dcLeft) * time.Second)
		tickingThing := time.NewTicker(TICKEVERY * time.Second)
		for {
			select {
			case <-tickingThing.C:
				if self.Client.Opponent != self.Client {
					msg := "Time:" + strconv.Itoa(timeLeft)
					self.Client.Opponent.TimeOutput <- msg
				}
				self.dcLeft -= TICKEVERY
			case message, stillOpen := <-self.Client.Time:
				if stillOpen {
					log.Println("[CLOCK]", self.Client.UserID, "got a reconnect message!")
					if message {
						self.TellTheTime()
						timer.Stop()
						tickingThing.Stop()
						return
					} else {
						log.Fatal("[CLOCK]", self.Client.UserID, "has not sent reconnected tho he should, received", message, "instead")
					}
				} else {
					timer.Stop()
					tickingThing.Stop()
					return
				}
			case <-timer.C:
				log.Println("[CLOCK]", self.Client.UserID, "lost the game due to dcing for too long!")
				//ran out of time, gg
				self.Client.GiveUp()
				timer.Stop()
				tickingThing.Stop()
				return
			}
		}
	}
}
