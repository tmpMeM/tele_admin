package tele_service

import (
	"fmt"

	"github.com/AnnonaOrg/osenv"
	"github.com/umfaka/umfaka_core/internal/log"
	"github.com/umfaka/umfaka_core/internal/service"
	tele "gopkg.in/telebot.v3"
)

// 检查屏蔽名单 如果在黑名单或消息符合屏蔽条件，就返回错误，否则返回 nil
func CheckBlockedUser(c tele.Context) error {
	if isFromGroup := c.Message().FromGroup(); !isFromGroup {
		return nil
	}

	var userID, groupID int64
	if service.IsPenetrationShielding() {
		if sender := c.Message().OriginalSender; sender != nil {
			log.Debugf("OriginalSender: %+v", sender)
			userID = sender.ID
		} else {
			return nil
		}
	} else {
		if sender := c.Message().Sender; sender != nil {
			log.Debugf("Sender: %+v", sender)
			userID = sender.ID
		} else {
			return nil
		}
	}
	groupID = c.Message().Chat.ID
	log.Debugf("userID: %d,groupID: %d", userID, groupID)
	if osenv.IsBotManagerID(userID) || IsChatAdmin(c, userID) {
		return nil
	}

	// if isForwarded := c.Message().IsForwarded(); isForwarded {
	// 	return fmt.Errorf("Message Sender(%d,%d) IsForwarded", userID, groupID)
	// }

	if count, err := service.GetCountBlockedUserByUserIDAndGroupID(userID, groupID); count > 0 {
		return fmt.Errorf("Message Sender(%d,%d) is Blocked", userID, groupID)
	} else if err != nil {
		log.Errorf("GetCountBlockedUserByUserIDAndGroupID(%d,%d): %v", userID, groupID, err)
	}
	return nil
}
